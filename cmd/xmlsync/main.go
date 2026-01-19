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

// XMLProperty represents a property in the XML feed
type XMLProperty struct {
	ID          string `xml:"id,attr"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Price       string `xml:"price"`
	Location    string `xml:"location"`
	Country     string `xml:"country"`
	Type        string `xml:"type"`
	Bedrooms    string `xml:"bedrooms"`
	Bathrooms   string `xml:"bathrooms"`
	Area        string `xml:"area"`
	YearBuilt   string `xml:"yearbuilt"`
	Features    string `xml:"features"`
	Images      []string `xml:"images>image"`
	Photos      []string `xml:"photos>photo"`
	Pictures    []string `xml:"pictures>picture"`
	Picture     []string `xml:"picture"`
	ImageURL    string   `xml:"image_url"`
	Image       string   `xml:"image"`
}

// XMLProperties wraps a collection of properties
type XMLProperties struct {
	Properties []XMLProperty `xml:"property"`
}

// XMLClientDetails holds client metadata
type XMLClientDetails struct {
	ClientName string `xml:"clientName"`
}

// XMLClient wraps client info with their properties
type XMLClient struct {
	ClientDetails XMLClientDetails `xml:"clientDetails"`
	Properties    XMLProperties    `xml:"properties"`
}

// XMLClients wraps multiple clients
type XMLClients struct {
	Clients []XMLClient `xml:"client"`
}

// XMLDocument is the root feed structure
type XMLDocument struct {
	Clients XMLClients `xml:"clients"`
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
			listing := Listing{
				ID:          prop.ID,
				Title:       prop.Title,
				Description: prop.Description,
				Price:       formatPrice(prop.Price),
				Location:    prop.Location,
				Country:     prop.Country,
				ListingType: prop.Type,
				Bedrooms:    prop.Bedrooms,
				Bathrooms:   prop.Bathrooms,
				Area:        prop.Area,
				YearBuilt:   prop.YearBuilt,
				Features:    parseFeatures(prop.Features),
				Images:      aggregateImages(prop),
			}

			// Fallback ID if missing; sanitize spaces
			if strings.TrimSpace(listing.ID) == "" {
				listing.ID = strings.TrimSpace(listing.Title)
			}
			listing.ID = strings.ReplaceAll(listing.ID, " ", "-")

			listings = append(listings, listing)
		}
	}

	return listings, nil
}

func formatPrice(price string) string {
	if price == "" {
		return "Contact for pricing"
	}
	// Basic price formatting (add $ if not present)
	if !strings.HasPrefix(price, "$") && !strings.HasPrefix(price, "€") && !strings.HasPrefix(price, "£") {
		return "$" + price
	}
	return price
}

func parseFeatures(featuresStr string) []string {
	if featuresStr == "" {
		return []string{}
	}
	// Split by comma, semicolon, or newline
	features := strings.FieldsFunc(featuresStr, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n'
	})

	var cleaned []string
	for _, f := range features {
		trimmed := strings.TrimSpace(f)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func aggregateImages(prop XMLProperty) []string {
	var images []string

	addAll := func(list []string) {
		for _, img := range list {
			for _, token := range strings.Fields(img) {
				trimmed := strings.TrimSpace(token)
				if trimmed != "" {
					images = append(images, trimmed)
				}
			}
		}
	}

	addIf := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" {
			for _, token := range strings.Fields(s) {
				trimmed := strings.TrimSpace(token)
				if trimmed != "" {
					images = append(images, trimmed)
				}
			}
		}
	}

	addAll(prop.Images)
	addAll(prop.Photos)
	addAll(prop.Pictures)
	addAll(prop.Picture)
	addIf(prop.ImageURL)
	addIf(prop.Image)

	// Deduplicate while preserving order
	seen := make(map[string]bool)
	var unique []string
	for _, img := range images {
		if !seen[img] {
			unique = append(unique, img)
			seen[img] = true
		}
	}

	return unique
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

	// YAML frontmatter
	frontmatter := fmt.Sprintf(`title: "%s"
description: "%s"
id: "%s"
price: "%s"
location: "%s"
country: "%s"
listingtype: "%s"
bedrooms: "%s"
bathrooms: "%s"
area: "%s"
yearbuilt: "%s"
date: %d
draft: false
`,
		escapeYAML(listing.Title),
		escapeYAML(listing.Description),
		listing.ID,
		escapeYAML(listing.Price),
		escapeYAML(listing.Location),
		escapeYAML(listing.Country),
		escapeYAML(listing.ListingType),
		listing.Bedrooms,
		listing.Bathrooms,
		escapeYAML(listing.Area),
		listing.YearBuilt,
		time.Now().Unix(),
	)

	_, err = io.WriteString(file, frontmatter)
	if err != nil {
		return err
	}

	if len(listing.Images) > 0 {
		_, err = io.WriteString(file, "images:\n")
		if err != nil {
			return err
		}
		for _, img := range listing.Images {
			_, err = io.WriteString(file, fmt.Sprintf("  - \"%s\"\n", escapeYAML(img)))
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
			_, err = io.WriteString(file, fmt.Sprintf("  - \"%s\"\n", escapeYAML(feature)))
			if err != nil {
				return err
			}
		}
	}

	// Tags/taxonomies
	_, err = io.WriteString(file, fmt.Sprintf("countries:\n  - \"%s\"\n", escapeYAML(listing.Country)))
	if err != nil {
		return err
	}

	if listing.Location != "" {
		_, err = io.WriteString(file, fmt.Sprintf("locations:\n  - \"%s\"\n", escapeYAML(listing.Location)))
		if err != nil {
			return err
		}
	}

	if listing.ListingType != "" {
		_, err = io.WriteString(file, fmt.Sprintf("types:\n  - \"%s\"\n", escapeYAML(listing.ListingType)))
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
		content = "Premium property listing with curated details and investment potential."
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
