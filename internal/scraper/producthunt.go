package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type ProductHuntProduct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Votes       int    `json:"votes"`
}

func ScrapeProductHunt() ([]ProductHuntProduct, error) {
	// Create a new collector with proper configuration
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowURLRevisit(),
		colly.Async(true), // Enable async scraping
	)

	// Add extensions
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Set rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*producthunt.com*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	var allProducts []ProductHuntProduct
	var scrapeErrors []error

	// Debug logging
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s\n", r.URL)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping %s: %v\n", r.Request.URL, err)
		scrapeErrors = append(scrapeErrors, fmt.Errorf("error scraping %s: %w", r.Request.URL, err))
	})

	// Handle product cards with multiple selectors for robustness
	c.OnHTML("div[class*='styles_post'], div[class*='post-item'], article[class*='post']", func(e *colly.HTMLElement) {
		// Try different selectors for each field
		name := findText(e, []string{
			"h3[class*='styles_title']",
			"h3[class*='title']",
			"h3",
		})

		description := findText(e, []string{
			"p[class*='styles_description']",
			"p[class*='description']",
			"p",
		})

		url := findAttr(e, "a[href*='/posts/']", "href", []string{
			"a[href*='/posts/']",
			"a[class*='post-link']",
		})

		votes := findVotes(e, []string{
			"div[class*='styles_voteCount']",
			"div[class*='vote-count']",
			"span[class*='votes']",
		})

		// Only add if we have at least a name and URL
		if name != "" && url != "" {
			product := ProductHuntProduct{
				Name:        name,
				Description: description,
				URL:         makeAbsoluteURL(url),
				Category:    e.Request.Ctx.Get("category"),
				Votes:       votes,
			}
			allProducts = append(allProducts, product)
		}
	})

	// Visit each category
	categories := []string{"", "tech", "productivity", "ai"}
	for _, category := range categories {
		url := "https://www.producthunt.com"
		if category != "" {
			url += "/topics/" + category
		}

		ctx := colly.NewContext()
		ctx.Put("category", category)

		err := c.Request("GET", url, nil, ctx, nil)
		if err != nil {
			scrapeErrors = append(scrapeErrors, fmt.Errorf("error visiting %s: %w", url, err))
			continue
		}
	}

	// Wait for all requests to complete
	c.Wait()

	// If we have any errors but also some products, return both
	if len(scrapeErrors) > 0 && len(allProducts) > 0 {
		log.Printf("Completed scraping with %d errors and %d products\n", len(scrapeErrors), len(allProducts))
	}

	return allProducts, nil
}

// Helper function to find text using multiple selectors
func findText(e *colly.HTMLElement, selectors []string) string {
	for _, selector := range selectors {
		if text := e.ChildText(selector); text != "" {
			return strings.TrimSpace(text)
		}
	}
	return ""
}

// Helper function to find attribute using multiple selectors
func findAttr(e *colly.HTMLElement, selector, attr string, fallbackSelectors []string) string {
	if value := e.ChildAttr(selector, attr); value != "" {
		return value
	}
	for _, selector := range fallbackSelectors {
		if value := e.ChildAttr(selector, attr); value != "" {
			return value
		}
	}
	return ""
}

// Helper function to find votes using multiple selectors
func findVotes(e *colly.HTMLElement, selectors []string) int {
	for _, selector := range selectors {
		if votesStr := e.ChildText(selector); votesStr != "" {
			return parseVotes(votesStr)
		}
	}
	return 0
}

// Helper function to make URLs absolute
func makeAbsoluteURL(url string) string {
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "http") {
		return url
	}
	return "https://www.producthunt.com" + url
}

// parseVotes converts the vote count string to an integer
func parseVotes(votesStr string) int {
	votesStr = strings.TrimSpace(votesStr)
	votesStr = strings.ReplaceAll(votesStr, "k", "000")
	votesStr = strings.ReplaceAll(votesStr, "m", "000000")
	votesStr = strings.ReplaceAll(votesStr, ",", "")

	var votes int
	fmt.Sscanf(votesStr, "%d", &votes)
	return votes
}

// CompareIdeaWithProducts compares an idea with scraped products and returns similar ones
func CompareIdeaWithProducts(idea string, products []ProductHuntProduct) []ProductHuntProduct {
	var similarProducts []ProductHuntProduct
	idea = strings.ToLower(idea)

	for _, product := range products {
		// Simple similarity check based on name and description
		if strings.Contains(strings.ToLower(product.Name), idea) ||
			strings.Contains(strings.ToLower(product.Description), idea) {
			similarProducts = append(similarProducts, product)
		}
	}

	return similarProducts
}
