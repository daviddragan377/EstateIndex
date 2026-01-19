package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Listing represents a property listing
type Listing struct {
	ID          string
	Title       string
	Description string
	Price       string
	Location    string
	Country     string
	ListingType string
	Bedrooms    string
	Bathrooms   string
	Area        string
	YearBuilt   string
	Features    []string
	Images      []string
}

// XMLAddress holds property address information
type XMLAddress struct {
	Street    string `xml:"street"`
	Location  string `xml:"location"`
	Region    string `xml:"region"`
	Country   string `xml:"country"`
	Latitude  string `xml:"latitude"`
	Longitude string `xml:"longitude"`
}

// XMLPrice holds pricing information
type XMLPrice struct {
	Price    string `xml:"price"`
	Currency string `xml:"currency"`
}

// XMLDescriptionContent holds detailed description
type XMLDescriptionContent struct {
	PropertyType  string `xml:"propertyType"`
	Title         string `xml:"title"`
	Description   string `xml:"description"`
	Bedrooms      string `xml:"bedrooms"`
	FullBathrooms string `xml:"fullBathrooms"`
	HalfBathrooms string `xml:"halfBathrooms"`
	YearBuilt     string `xml:"yearBuilt"`
	Heating       string `xml:"heating"`
	Elevator      string `xml:"elevator"`
	SwimmingPool  string `xml:"swimmingPool"`
	Furnishings   string `xml:"furnishings"`
	Features      struct {
		Feature []string `xml:"Feature"`
	} `xml:"Features"`
	FloorSize struct {
		Size  string `xml:"floorSize"`
		Units string `xml:"floorSizeUnits"`
	} `xml:"FloorSize"`
}

// XMLImage holds individual image data
type XMLImage struct {
	URL string `xml:"image"`
}

// XMLImages holds all images
type XMLImages struct {
	Images []XMLImage `xml:"image"`
}

// XMLProperty represents a single property in the XML feed
type XMLProperty struct {
	PropertyID  string                `xml:"propertyid"`
	LastUpdate  string                `xml:"lastUpdateDate"`
	Category    string                `xml:"category"`
	Address     XMLAddress            `xml:"Address"`
	Price       XMLPrice              `xml:"Price"`
	Description XMLDescriptionContent `xml:"Description"`
	Images      XMLImages             `xml:"images"`
}

// XMLProperties wraps a collection of properties
type XMLProperties struct {
	Properties []XMLProperty `xml:"Property"`
}

// XMLClientDetails holds client metadata
type XMLClientDetails struct {
	ClientName string `xml:"clientName"`
}

// XMLClient wraps client info with their properties
type XMLClient struct {
	ClientDetails XMLClientDetails `xml:"ClientDetails"`
	Properties    XMLProperties    `xml:"properties"`
}

// XMLClients wraps multiple clients
type XMLClients struct {
	Clients []XMLClient `xml:"Client"`
}

// XMLDocument is the root feed structure
type XMLDocument struct {
	Clients XMLClients `xml:"Clients"`
}

func main() {
	feedURL := flag.String("feed", "https://www.xml2u.com/Xml/International%20Property%20Alerts_3968/7212_Default.xml", "XML feed URL")
	contentDir := flag.String("content", "./content/listings", "Content directory")
	dryRun := flag.Bool("dry-run", false, "Dry run mode")
	flag.Parse()

	// If a local xml_feed.txt exists, prefer its contents as the feed URL
	if data, err := os.ReadFile("xml_feed.txt"); err == nil {
		candidate := strings.TrimSpace(string(data))
		if candidate != "" {
			*feedURL = candidate
		}
	}

	log.Println("Estate Index XML Sync Tool")
	log.Println("==========================")
	log.Printf("Feed URL: %s\n", *feedURL)
	log.Printf("Content Dir: %s\n", *contentDir)
	if *dryRun {
		log.Println("Mode: DRY RUN (no files will be written)")
	}
	log.Println()

	// Fetch and parse XML
	listings, err := fetchAndParseFeed(*feedURL)
	if err != nil {
		log.Fatalf("Error fetching feed: %v", err)
	}

	log.Printf("Fetched %d listings from feed\n\n", len(listings))

	// Get existing listings
	existingFiles, err := getExistingListings(*contentDir)
	if err != nil {
		log.Fatalf("Error reading existing listings: %v", err)
	}

	log.Printf("Found %d existing listing files\n\n", len(existingFiles))

	// Track changes
	var added, updated, removed int
	newListingIDs := make(map[string]bool)

	// Write/update listings
	for _, listing := range listings {
		newListingIDs[listing.ID] = true
		filename := filepath.Join(*contentDir, fmt.Sprintf("%s.md", listing.ID))

		if _, exists := existingFiles[listing.ID]; exists {
			updated++
			log.Printf("[UPDATE] %s: %s\n", listing.ID, listing.Title)
		} else {
			added++
			log.Printf("[ADD] %s: %s\n", listing.ID, listing.Title)
		}

		if !*dryRun {
			if err := writeListingFile(filename, listing); err != nil {
				log.Printf("Error writing file %s: %v\n", filename, err)
			}
		}
	}

	// Remove listings no longer in feed
	for id := range existingFiles {
		if !newListingIDs[id] {
			removed++
			filename := filepath.Join(*contentDir, fmt.Sprintf("%s.md", id))
			log.Printf("[REMOVE] %s\n", id)

			if !*dryRun {
				if err := os.Remove(filename); err != nil {
					log.Printf("Error removing file %s: %v\n", filename, err)
				}
			}
		}
	}

	// Summary
	log.Println()
	log.Println("==========================")
	log.Printf("Summary:\n")
	log.Printf("  Added:   %d\n", added)
	log.Printf("  Updated: %d\n", updated)
	log.Printf("  Removed: %d\n", removed)
	log.Printf("  Total:   %d listings\n", len(listings))
	if *dryRun {
		log.Println("\nDRY RUN - No files were written")
	}
}

func fetchAndParseFeed(feedURL string) ([]Listing, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var doc XMLDocument
	if err := xml.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	var listings []Listing

	// Iterate through clients and their properties
	for _, client := range doc.Clients.Clients {
		for _, prop := range client.Properties.Properties {
			// Extract area from floor size if available
			area := ""
			if prop.Description.FloorSize.Size != "" {
				area = prop.Description.FloorSize.Size + " " + prop.Description.FloorSize.Units
			}

			// Extract features
			var features []string
			for _, f := range prop.Description.Features.Feature {
				trimmed := strings.TrimSpace(f)
				if trimmed != "" {
					features = append(features, trimmed)
				}
			}

			// Extract images
			var images []string
			for _, img := range prop.Images.Images {
				if strings.TrimSpace(img.URL) != "" {
					images = append(images, strings.TrimSpace(img.URL))
				}
			}

			// Parse description text (remove CDATA markers if present)
			description := strings.TrimSpace(prop.Description.Description)
			description = strings.ReplaceAll(description, "<![CDATA[", "")
			description = strings.ReplaceAll(description, "]]>", "")
			description = strings.TrimSpace(description)

			listing := Listing{
				ID:          strings.TrimSpace(prop.PropertyID),
				Title:       strings.TrimSpace(prop.Description.Title),
				Description: description,
				Price:       formatPrice(prop.Price.Price, prop.Price.Currency),
				Location:    strings.TrimSpace(prop.Address.Location),
				Country:     strings.TrimSpace(prop.Address.Country),
				ListingType: strings.TrimSpace(prop.Description.PropertyType),
				Bedrooms:    strings.TrimSpace(prop.Description.Bedrooms),
				Bathrooms:   strings.TrimSpace(prop.Description.FullBathrooms),
				Area:        area,
				YearBuilt:   strings.TrimSpace(prop.Description.YearBuilt),
				Features:    features,
				Images:      images,
			}

			// Generate slug ID if missing
			if listing.ID == "" {
				listing.ID = strings.ToLower(strings.ReplaceAll(listing.Title, " ", "-"))
			}

			listings = append(listings, listing)
		}
	}

	return listings, nil
}

func formatPrice(price string, currency string) string {
	if price == "" {
		return "Contact for pricing"
	}

	currencySymbol := "$"
	if currency == "€" || currency == "EUR" {
		currencySymbol = "€"
	} else if currency == "£" || currency == "GBP" {
		currencySymbol = "£"
	}

	// If price already has a symbol, don't add another
	if strings.HasPrefix(price, "$") || strings.HasPrefix(price, "€") || strings.HasPrefix(price, "£") {
		return price
	}

	return currencySymbol + price
}

func getExistingListings(contentDir string) (map[string]bool, error) {
	existing := make(map[string]bool)

	entries, err := os.ReadDir(contentDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist yet, create it
			if err := os.MkdirAll(contentDir, 0755); err != nil {
				return nil, err
			}
			return existing, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			id := strings.TrimSuffix(entry.Name(), ".md")
			existing[id] = true
		}
	}

	return existing, nil
}

func writeListingFile(filename string, listing Listing) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write frontmatter
	_, err = io.WriteString(file, "---\n")
	if err != nil {
		return err
	}

	// Helper to safely format optional strings
	formatOpt := func(v string) string {
		if v == "" {
			return "null"
		}
		return fmt.Sprintf(`"%s"`, escapeYAML(v))
	}

	// YAML frontmatter
	frontmatter := fmt.Sprintf(`title: %s
description: %s
id: "%s"
price: %s
location: %s
country: %s
listingtype: %s
bedrooms: %s
bathrooms: %s
area: %s
yearbuilt: %s
date: %d
draft: false
`,
		formatOpt(listing.Title),
		formatOpt(listing.Description),
		listing.ID,
		formatOpt(listing.Price),
		formatOpt(listing.Location),
		formatOpt(listing.Country),
		formatOpt(listing.ListingType),
		formatOpt(listing.Bedrooms),
		formatOpt(listing.Bathrooms),
		formatOpt(listing.Area),
		formatOpt(listing.YearBuilt),
		time.Now().Unix(),
	)

	_, err = io.WriteString(file, frontmatter)
	if err != nil {
		return err
	}

	// Images array in frontmatter
	if len(listing.Images) > 0 {
		_, err = io.WriteString(file, "images:\n")
		if err != nil {
			return err
		}
		for _, img := range listing.Images {
			_, err = io.WriteString(file, fmt.Sprintf("  - %s\n", formatOpt(img)))
			if err != nil {
				return err
			}
		}
	}

	// Features array in frontmatter
	if len(listing.Features) > 0 {
		_, err = io.WriteString(file, "features:\n")
		if err != nil {
			return err
		}
		for _, feature := range listing.Features {
			_, err = io.WriteString(file, fmt.Sprintf("  - %s\n", formatOpt(feature)))
			if err != nil {
				return err
			}
		}
	}

	// Tags/taxonomies
	if listing.Country != "" {
		_, err = io.WriteString(file, fmt.Sprintf("countries:\n  - %s\n", formatOpt(listing.Country)))
		if err != nil {
			return err
		}
	}

	if listing.Location != "" {
		_, err = io.WriteString(file, fmt.Sprintf("locations:\n  - %s\n", formatOpt(listing.Location)))
		if err != nil {
			return err
		}
	}

	if listing.ListingType != "" {
		_, err = io.WriteString(file, fmt.Sprintf("types:\n  - %s\n", formatOpt(listing.ListingType)))
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(file, "---\n\n")
	if err != nil {
		return err
	}

	// Write content
	content := listing.Description
	if content == "" {
		content = "Premium property listing."
	}

	_, err = io.WriteString(file, content)
	return err
}

func escapeYAML(s string) string {
	// Remove problematic characters and escape quotes
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
