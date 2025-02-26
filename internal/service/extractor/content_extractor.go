package extractor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/pemistahl/lingua-go"
)

type ContentExtractor interface {
	Extract(ctx context.Context, url string) (*ExtractedContent, error)
}

type ExtractedContent struct {
	Title       string
	Content     string
	Language    string
	SiteName    string
	Author      string
	Excerpt     string
	ImageURL    string
	PublishDate string
}

type contentExtractor struct {
	languageDetector lingua.LanguageDetector
	httpClient       *http.Client
}

func NewContentExtractor() (ContentExtractor, error) {
	languages := []lingua.Language{
		lingua.English,
		lingua.Indonesian,
		lingua.Spanish,
		lingua.French,
		lingua.German,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &contentExtractor{
		languageDetector: detector,
		httpClient:       client,
	}, nil
}

func (ce *contentExtractor) Extract(ctx context.Context, source string) (*ExtractedContent, error) {
	if source == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// Validate URL format
	urlRes, err := url.ParseRequestURI(source)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to mimic a regular browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	start := time.Now()
	resp, err := ce.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch webpage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("webpage returned status code: %d", resp.StatusCode)
	}

	// Read the body once and create a bytes.Reader for readability
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create a new reader from the body bytes
	article, err := readability.FromReader(bytes.NewReader(body), urlRes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webpage content: %w", err)
	}
	log.Printf("Scraping process completed in %s", time.Since(start))

	// Detect language
	langStart := time.Now()
	language := ce.detectLanguage(article.TextContent)
	log.Printf("Language detection completed in %s", time.Since(langStart))

	// Handle nil PublishedTime
	publishDate := ""
	if article.PublishedTime != nil {
		publishDate = article.PublishedTime.Format(time.RFC3339)
	}

	result := &ExtractedContent{
		Title:       article.Title,
		Content:     article.TextContent,
		Language:    language,
		SiteName:    article.SiteName,
		Author:      article.Byline,
		Excerpt:     article.Excerpt,
		ImageURL:    article.Image,
		PublishDate: publishDate,
	}

	return result, nil
}

func (ce *contentExtractor) detectLanguage(text string) string {
	language, ableToDetect := ce.languageDetector.DetectLanguageOf(text)
	if !ableToDetect {
		return ""
	} else {
		return language.String()
	}
}
