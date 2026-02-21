package handlers

import "strings"

// Reader mode CSS and JS that gets injected into proxied pages.
// Defaults to reader view on load; user can switch to raw view.
const readerModeSnippet = `
<!-- Extension Ladder Reader Mode -->
<script>
(function() {
    var inReader = true;
    var originalHTML = document.documentElement.innerHTML;

    // Article content selectors in priority order
    var SELECTORS = [
        'article',
        '[role="article"]',
        '[itemprop="articleBody"]',
        '.article-body',
        '.article-content',
        '.story-body',
        '.story-content',
        '.post-content',
        '.entry-content',
        '.content-body',
        '.article__body',
        '.article__content',
        'main',
        '#article-body',
        '#story-body',
        '#content',
        '.td-post-content',
        '.article_content'
    ];

    function extractArticle() {
        var title = '';
        var titleEl = document.querySelector('h1');
        if (titleEl) title = titleEl.innerHTML;

        var content = null;
        for (var i = 0; i < SELECTORS.length; i++) {
            var el = document.querySelector(SELECTORS[i]);
            if (el && el.textContent.trim().length > 200) {
                content = el.innerHTML;
                break;
            }
        }

        if (!content) {
            var paragraphs = document.querySelectorAll('p');
            var collected = [];
            paragraphs.forEach(function(p) {
                if (p.textContent.trim().length > 40) {
                    collected.push(p.outerHTML);
                }
            });
            content = collected.join('\n');
        }

        return { title: title, content: content };
    }

    function buildReaderPage(article) {
        return '<html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">' +
            '<title>Reader — Extension Ladder</title>' +
            '<link rel="icon" type="image/x-icon" href="/favicon.ico">' +
            '<style>' +
            '*{box-sizing:border-box;}' +
            'body{max-width:680px;margin:40px auto;padding:0 20px;font-family:Georgia,"Times New Roman",serif;font-size:19px;line-height:1.7;color:#1c1917;background:#fafaf9;}' +
            'h1{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:1.8em;line-height:1.3;margin-bottom:0.5em;color:#0c0a09;}' +
            'h2,h3{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;line-height:1.3;margin-top:1.5em;}' +
            'img{max-width:100%;height:auto;display:block;margin:1em auto;border-radius:4px;}' +
            'figure{margin:1em 0;padding:0;}' +
            'figcaption{font-size:0.85em;color:#666;text-align:center;margin-top:0.5em;}' +
            'a{color:#2563eb;}' +
            'p{margin:1em 0;}' +
            'blockquote{border-left:3px solid #d1d5db;margin:1em 0;padding:0.5em 1em;color:#4b5563;}' +
            'ul,ol{margin:1em 0;padding-left:1.5em;}' +
            'li{margin:0.3em 0;}' +
            'pre,code{font-size:0.9em;background:#f3f4f6;border-radius:3px;padding:0.2em 0.4em;}' +
            'pre{padding:1em;overflow-x:auto;}' +
            '#el-bar{position:fixed;top:0;left:0;right:0;background:#1e293b;color:#e2e8f0;padding:8px 16px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:13px;display:flex;align-items:center;justify-content:space-between;z-index:9999;box-shadow:0 2px 8px rgba(0,0,0,0.2);}' +
            '#el-bar a{color:#93c5fd;text-decoration:none;}' +
            '#el-bar button{background:transparent;color:#e2e8f0;border:1px solid #475569;border-radius:4px;padding:4px 12px;cursor:pointer;font-size:12px;transition:all 0.2s;}' +
            '#el-bar button:hover{background:#334155;border-color:#94a3b8;}' +
            '@media(prefers-color-scheme:dark){body{background:#1a1a2e;color:#e2e8f0;}h1,h2,h3{color:#f1f5f9;}a{color:#93c5fd;}blockquote{border-color:#475569;color:#9ca3af;}pre,code{background:#1e293b;}}' +
            '</style></head><body>' +
            '<div id="el-bar">' +
            '<span style="display:flex;align-items:center;gap:8px;">📖 <a href="/">Extension Ladder</a></span>' +
            '<div style="display:flex;gap:8px;align-items:center;">' +
            '<span style="font-size:11px;opacity:0.5;">Press R for raw view</span>' +
            '<button id="el-raw">View Raw Page</button>' +
            '</div></div>' +
            '<div style="margin-top:50px;">' +
            (article.title ? '<h1>' + article.title + '</h1>' : '') +
            (article.content || '<p>Could not extract article content.</p>') +
            '</div></body></html>';
    }

    function showReader() {
        var article = extractArticle();
        inReader = true;
        document.documentElement.innerHTML = buildReaderPage(article);
        document.getElementById('el-raw').addEventListener('click', showRaw);
    }

    function showRaw() {
        inReader = false;
        document.documentElement.innerHTML = originalHTML;
        // Inject a floating "Reader" button on the raw page
        var floatBtn = document.createElement('button');
        floatBtn.id = 'el-reader-toggle';
        floatBtn.title = 'Switch to Reader Mode';
        floatBtn.textContent = '📖';
        floatBtn.style.cssText = 'position:fixed;bottom:20px;right:20px;z-index:999999;background:#3b82f6;color:white;border:none;border-radius:50%;width:44px;height:44px;cursor:pointer;font-size:20px;box-shadow:0 2px 8px rgba(0,0,0,0.3);display:flex;align-items:center;justify-content:center;';
        floatBtn.addEventListener('click', function() {
            originalHTML = document.documentElement.innerHTML;
            showReader();
        });
        document.body.appendChild(floatBtn);
    }

    document.addEventListener('keydown', function(e) {
        if (e.key === 'r' && !e.ctrlKey && !e.metaKey && !e.altKey &&
            e.target.tagName !== 'INPUT' && e.target.tagName !== 'TEXTAREA') {
            if (inReader) { showRaw(); } else { originalHTML = document.documentElement.innerHTML; showReader(); }
        }
    });

    // Auto-activate reader mode on page load
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', showReader);
    } else {
        showReader();
    }
})();
</script>
<!-- End Extension Ladder Reader Mode -->
`

// injectReaderMode injects the reader mode toggle button and styles into HTML pages
func injectReaderMode(body string) string {
	// Only inject into HTML pages (check for </body> tag)
	closingBody := "</body>"
	idx := strings.LastIndex(strings.ToLower(body), closingBody)
	if idx == -1 {
		// No </body> found, try appending at the end
		return body + readerModeSnippet
	}

	// Find the actual case-insensitive position
	lowerBody := strings.ToLower(body)
	idx = strings.LastIndex(lowerBody, closingBody)

	return body[:idx] + readerModeSnippet + body[idx:]
}
