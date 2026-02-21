package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"extension-ladder/pkg/ruleset"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// rewriteHtml
// ---------------------------------------------------------------------------

func TestRewriteHtml_ImagesRewritten(t *testing.T) {
	body := []byte(`<img src="/image.jpg">`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, `/https://example.com/image.jpg`)
}

func TestRewriteHtml_ScriptsRewritten(t *testing.T) {
	body := []byte(`<script src="/app.js"></script>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, `src="/https://example.com/app.js"`)
	// Verify it doesn't produce the old 'script=' attribute bug
	assert.NotContains(t, result, `script="/https://example.com/app.js"`)
}

func TestRewriteHtml_RelativeHrefsRewritten(t *testing.T) {
	body := []byte(`<a href="/about">About</a>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, `href="/https://example.com/about"`)
}

func TestRewriteHtml_AlreadyRewrittenHrefsSkipped(t *testing.T) {
	body := []byte(`<a href="/https://example.com/about">About</a>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	// Should not double-rewrite
	assert.NotContains(t, result, `/https://example.com/https://example.com/`)
}

func TestRewriteHtml_CSSUrlsRewritten(t *testing.T) {
	body := []byte(`<div style="background-image: url('/bg.jpg')"></div>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, `url('/https://example.com/bg.jpg')`)
}

func TestRewriteHtml_CSSUrlsUnquotedRewritten(t *testing.T) {
	body := []byte(`<div style="background-image: url(/bg.jpg)"></div>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, `url(/https://example.com/bg.jpg)`)
}

func TestRewriteHtml_InjectsReaderMode(t *testing.T) {
	body := []byte(`<html><body><p>Hello</p></body></html>`)
	u := &url.URL{Host: "example.com"}
	result := rewriteHtml(body, u, ruleset.Rule{})
	assert.Contains(t, result, "Extension Ladder Reader Mode")
}

// ---------------------------------------------------------------------------
// isBotDetected
// ---------------------------------------------------------------------------

func TestIsBotDetected_DefaultPatterns(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{"captcha page", "<html><body>Please complete the captcha</body></html>", true},
		{"access denied", "<html><body>Access Denied - please verify</body></html>", true},
		{"security check", "<html><body>Security Check required</body></html>", true},
		{"challenge platform", "<html><body><script src='challenge-platform'></script></body></html>", true},
		{"clean page", "<html><body><p>This is a normal article about cooking</p></body></html>", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBotDetected(tt.body, ruleset.Rule{})
			assert.Equal(t, tt.expected, result, "body: %s", tt.body)
		})
	}
}

func TestIsBotDetected_CustomPatterns(t *testing.T) {
	rule := ruleset.Rule{
		BotDetectionPatterns: []string{"perimeterx", "px-captcha"},
	}
	assert.True(t, isBotDetected("<html>PerimeterX challenge</html>", rule))
	assert.True(t, isBotDetected("<html><div id='px-captcha'></div></html>", rule))
	assert.False(t, isBotDetected("<html><p>Normal page</p></html>", rule))
}

func TestIsBotDetected_LargePageSkipped(t *testing.T) {
	// Pages >50KB are considered real content, even if they contain monitoring scripts
	largeBody := strings.Repeat("x", 50001) + " challenge-platform "
	assert.False(t, isBotDetected(largeBody, ruleset.Rule{}))
}

func TestIsBotDetected_SmallPageWithPattern(t *testing.T) {
	smallBody := "<html>challenge-platform</html>"
	assert.True(t, isBotDetected(smallBody, ruleset.Rule{}))
}

func TestIsBotDetected_ExactThreshold(t *testing.T) {
	// Exactly 50000 bytes should still be checked
	body := strings.Repeat("x", 49990) + " captcha "
	assert.True(t, isBotDetected(body, ruleset.Rule{}))
}

// ---------------------------------------------------------------------------
// determineFallbacks
// ---------------------------------------------------------------------------

func TestDetermineFallbacks(t *testing.T) {
	tests := []struct {
		strategy string
		expected []string
	}{
		{"headless+archive", []string{"headless", "archive"}},
		{"archive+headless", []string{"archive", "headless"}},
		{"", []string{"headless", "archive"}},
		{"unknown", []string{"headless", "archive"}},
	}

	for _, tt := range tests {
		t.Run("strategy_"+tt.strategy, func(t *testing.T) {
			result := determineFallbacks(tt.strategy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ---------------------------------------------------------------------------
// stripScriptTags
// ---------------------------------------------------------------------------

func TestStripScriptTags_InlineScript(t *testing.T) {
	input := `<html><body><script>alert('hi')</script><p>Content</p></body></html>`
	result := stripScriptTags(input)
	assert.NotContains(t, result, "<script")
	assert.NotContains(t, result, "alert")
	assert.Contains(t, result, "<p>Content</p>")
}

func TestStripScriptTags_ExternalScript(t *testing.T) {
	input := `<html><script src="/app.js"></script><body><p>Content</p></body></html>`
	result := stripScriptTags(input)
	assert.NotContains(t, result, "<script")
	assert.NotContains(t, result, "app.js")
	assert.Contains(t, result, "<p>Content</p>")
}

func TestStripScriptTags_MultilineScript(t *testing.T) {
	input := `<html><body>
<script type="text/javascript">
  var x = 1;
  var y = 2;
  console.log(x + y);
</script>
<p>Content</p></body></html>`
	result := stripScriptTags(input)
	assert.NotContains(t, result, "<script")
	assert.NotContains(t, result, "console.log")
	assert.Contains(t, result, "<p>Content</p>")
}

func TestStripScriptTags_SelfClosing(t *testing.T) {
	input := `<html><body><script src="/app.js"/><p>Content</p></body></html>`
	result := stripScriptTags(input)
	assert.NotContains(t, result, "<script")
	assert.Contains(t, result, "<p>Content</p>")
}

func TestStripScriptTags_MultipleScripts(t *testing.T) {
	input := `<html><head><script src="/a.js"></script><script src="/b.js"></script></head><body><script>var x=1;</script><p>Keep</p></body></html>`
	result := stripScriptTags(input)
	assert.NotContains(t, result, "<script")
	assert.Contains(t, result, "<p>Keep</p>")
}

func TestStripScriptTags_NoScripts(t *testing.T) {
	input := `<html><body><p>No scripts here</p></body></html>`
	result := stripScriptTags(input)
	assert.Equal(t, input, result)
}

// ---------------------------------------------------------------------------
// modifyURL
// ---------------------------------------------------------------------------

func TestModifyURL_NoRules(t *testing.T) {
	result, err := modifyURL("https://example.com/article", ruleset.Rule{})
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/article", result)
}

func TestModifyURL_DomainRewrite(t *testing.T) {
	rule := ruleset.Rule{}
	rule.URLMods.Domain = []ruleset.Regex{
		{Match: "www\\.example\\.com", Replace: "example.com"},
	}
	result, err := modifyURL("https://www.example.com/path", rule)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/path", result)
}

func TestModifyURL_PathRewrite(t *testing.T) {
	rule := ruleset.Rule{}
	rule.URLMods.Path = []ruleset.Regex{
		{Match: "/amp/", Replace: "/"},
	}
	result, err := modifyURL("https://example.com/amp/article", rule)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/article", result)
}

func TestModifyURL_QueryAdd(t *testing.T) {
	rule := ruleset.Rule{}
	rule.URLMods.Query = []ruleset.KV{
		{Key: "utm_source", Value: ""},  // delete
		{Key: "format", Value: "clean"}, // add
	}
	result, err := modifyURL("https://example.com/page?utm_source=fb&id=1", rule)
	assert.NoError(t, err)
	parsed, _ := url.Parse(result)
	assert.Equal(t, "", parsed.Query().Get("utm_source"))
	assert.Equal(t, "clean", parsed.Query().Get("format"))
	assert.Equal(t, "1", parsed.Query().Get("id"))
}

func TestModifyURL_GoogleCache(t *testing.T) {
	rule := ruleset.Rule{GoogleCache: true}
	result, err := modifyURL("https://example.com/article", rule)
	assert.NoError(t, err)
	assert.Contains(t, result, "webcache.googleusercontent.com")
	assert.Contains(t, result, "cache:https://example.com/article")
}

func TestModifyURL_InvalidURL(t *testing.T) {
	_, err := modifyURL("://bad-url", ruleset.Rule{})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// StringInSlice
// ---------------------------------------------------------------------------

func TestStringInSlice(t *testing.T) {
	list := []string{"/article", "/blog", "/news"}

	assert.True(t, StringInSlice("/article", list))
	assert.True(t, StringInSlice("/article/123", list)) // prefix match
	assert.True(t, StringInSlice("/blog", list))
	assert.False(t, StringInSlice("/about", list))
	assert.False(t, StringInSlice("", list))
}

func TestStringInSlice_EmptyList(t *testing.T) {
	assert.False(t, StringInSlice("/anything", []string{}))
}

// ---------------------------------------------------------------------------
// injectReaderMode
// ---------------------------------------------------------------------------

func TestInjectReaderMode_WithBody(t *testing.T) {
	body := `<html><body><p>Article</p></body></html>`
	result := injectReaderMode(body)
	assert.Contains(t, result, "Extension Ladder Reader Mode")
	// Snippet should be before </body>
	readerIdx := strings.Index(result, "Extension Ladder Reader Mode")
	bodyIdx := strings.LastIndex(strings.ToLower(result), "</body>")
	assert.Less(t, readerIdx, bodyIdx)
}

func TestInjectReaderMode_WithoutBody(t *testing.T) {
	body := `<p>Fragment without body tags</p>`
	result := injectReaderMode(body)
	assert.Contains(t, result, "Extension Ladder Reader Mode")
}

func TestInjectReaderMode_CaseInsensitive(t *testing.T) {
	body := `<html><BODY><p>Content</p></BODY></html>`
	result := injectReaderMode(body)
	assert.Contains(t, result, "Extension Ladder Reader Mode")
}

// ---------------------------------------------------------------------------
// getenv
// ---------------------------------------------------------------------------

func TestGetenv_Fallback(t *testing.T) {
	result := getenv("DEFINITELY_NOT_SET_12345", "default_val")
	assert.Equal(t, "default_val", result)
}

func TestGetenv_Set(t *testing.T) {
	t.Setenv("TEST_GETENV_KEY", "custom_val")
	result := getenv("TEST_GETENV_KEY", "default_val")
	assert.Equal(t, "custom_val", result)
}

// ---------------------------------------------------------------------------
// applyRules
// ---------------------------------------------------------------------------

func TestApplyRules_RegexRules(t *testing.T) {
	// Temporarily populate the global rulesSet so applyRules doesn't short-circuit
	origRulesSet := rulesSet
	rulesSet = ruleset.RuleSet{{Domain: "test.com"}}
	defer func() { rulesSet = origRulesSet }()

	rule := ruleset.Rule{
		RegexRules: []ruleset.Regex{
			{Match: `class="paywall"`, Replace: `class=""`},
			{Match: `display:\s*none`, Replace: `display: block`},
		},
	}

	body := `<div class="paywall" style="display: none"><p>Hidden content</p></div>`
	result := applyRules(body, rule)
	assert.Contains(t, result, `class=""`)
	assert.Contains(t, result, `display: block`)
	assert.NotContains(t, result, `class="paywall"`)
}

func TestApplyRules_NoRuleset(t *testing.T) {
	// With empty global rulesSet, should return body unchanged
	origRulesSet := rulesSet
	rulesSet = ruleset.RuleSet{}
	defer func() { rulesSet = origRulesSet }()

	body := `<div class="paywall"><p>Content</p></div>`
	result := applyRules(body, ruleset.Rule{})
	assert.Equal(t, body, result)
}

// ---------------------------------------------------------------------------
// ProxySite - integration test with local mock server
// ---------------------------------------------------------------------------

func TestProxySite_MockServer(t *testing.T) {
	// Stand up a local HTTP server instead of hitting external URLs
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><p>Mock article content for testing</p></body></html>`))
	}))
	defer mockServer.Close()

	app := fiber.New()
	app.Get("/*", ProxySite(""))

	req := httptest.NewRequest("GET", "/"+mockServer.URL, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
