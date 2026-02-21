package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Stealth JS to inject before page load — hides automation indicators
const stealthJS = `
// Override navigator.webdriver
Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

// Override chrome runtime
window.chrome = { runtime: {} };

// Override permissions query
const originalQuery = window.navigator.permissions.query;
window.navigator.permissions.query = (parameters) =>
    parameters.name === 'notifications'
        ? Promise.resolve({ state: Notification.permission })
        : originalQuery(parameters);

// Override plugins
Object.defineProperty(navigator, 'plugins', {
    get: () => [1, 2, 3, 4, 5],
});

// Override languages
Object.defineProperty(navigator, 'languages', {
    get: () => ['en-US', 'en'],
});

// Remove automation-specific properties
delete navigator.__proto__.webdriver;
`

// fetchWithHeadless uses a headless Chrome browser to fetch a URL,
// allowing JavaScript to execute and bot-detection challenges to be handled.
func fetchWithHeadless(targetURL string, waitSeconds int) (string, error) {
	if waitSeconds <= 0 {
		waitSeconds = 8
	}

	// Create a new Chrome context with stealth options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.Flag("disable-site-isolation-trials", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set an overall timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(waitSeconds+30)*time.Second)
	defer cancel()

	var body string

	log.Printf("INFO: Fetching with headless browser: %s (wait: %ds)", targetURL, waitSeconds)

	err := chromedp.Run(ctx,
		// Inject stealth scripts before navigation
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(stealthJS).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return fmt.Errorf("stealth JS exception: %v", exp)
			}
			return nil
		}),
		chromedp.Navigate(targetURL),
		// Wait for the page to settle
		chromedp.Sleep(time.Duration(waitSeconds)*time.Second),
		// Try to wait for body content to be non-empty
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &body),
	)
	if err != nil {
		return "", fmt.Errorf("headless fetch failed for '%s': %w", targetURL, err)
	}

	log.Printf("INFO: Headless fetch complete for: %s (%d bytes)", targetURL, len(body))
	return body, nil
}
