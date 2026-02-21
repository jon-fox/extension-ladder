package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// fetchFromArchive fetches a URL via the Wayback Machine (archive.org).
// It first tries the latest available snapshot.
func fetchFromArchive(targetURL string) (string, error) {
	// Use the Wayback Machine's "id_" modifier to get the original page content
	archiveURL := fmt.Sprintf("https://web.archive.org/web/2/%s", targetURL)

	log.Printf("INFO: Fetching from archive.org: %s", archiveURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", archiveURL, nil)
	if err != nil {
		return "", fmt.Errorf("archive.org request creation failed: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("archive.org fetch failed for '%s': %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("archive.org returned status %d for '%s'", resp.StatusCode, targetURL)
	}

	bodyB, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("archive.org read failed: %w", err)
	}

	body := string(bodyB)
	log.Printf("INFO: Archive.org fetch complete for: %s (%d bytes)", targetURL, len(body))

	return body, nil
}
