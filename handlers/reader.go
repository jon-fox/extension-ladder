package handlers

import "strings"

// Reader mode CSS and JS that gets injected into proxied pages
const readerModeSnippet = `
<!-- Extension Ladder Reader Mode -->
<style id="el-reader-style">
#el-reader-toggle {
    position: fixed;
    bottom: 20px;
    right: 20px;
    z-index: 999999;a
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 50%;
    width: 44px;
    height: 44px;
    cursor: pointer;
    font-size: 20px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
    line-height: 1;
}
#el-reader-toggle:hover {
    background: #2563eb;
}
</style>
<button id="el-reader-toggle" title="Toggle Reader Mode">📖</button>
<script>
(function() {
    var active = false;
    var btn = document.getElementById('el-reader-toggle');
    var bar = document.getElementById('el-reader-bar');
    var originalHTML = '';

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
        // Try to find the article title
        var title = '';
        var titleEl = document.querySelector('h1');
        if (titleEl) title = titleEl.innerHTML;

        // Try to find article content using selectors
        var content = null;
        for (var i = 0; i < SELECTORS.length; i++) {
            var el = document.querySelector(SELECTORS[i]);
            if (el && el.textContent.trim().length > 200) {
                content = el.innerHTML;
                break;
            }
        }

        // Fallback: grab all paragraphs
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

        // Collect images from the article area
        var images = document.querySelectorAll('article img, main img, .article-body img, [role="article"] img, figure img');
        var imgHTML = '';
        images.forEach(function(img) {
            if (img.src && img.naturalWidth > 100) {
                imgHTML += '<img src="' + img.src + '" alt="' + (img.alt || '') + '">\n';
            }
        });

        return { title: title, content: content, images: imgHTML };
    }

    function toggle() {
        active = !active;
        if (active) {
            // Save original page
            originalHTML = document.documentElement.innerHTML;

            // Extract article
            var article = extractArticle();

            // Replace the entire page
            document.documentElement.innerHTML = '<!DOCTYPE html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Reader Mode — Extension Ladder</title>' +
                '<style>' +
                'body{max-width:680px;margin:40px auto;padding:0 20px;font-family:Georgia,"Times New Roman",serif;font-size:19px;line-height:1.7;color:#1c1917;background:#fafaf9;}' +
                'h1{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:1.8em;line-height:1.3;margin-bottom:0.5em;color:#0c0a09;}' +
                'img{max-width:100%;height:auto;display:block;margin:1em auto;border-radius:4px;}' +
                'a{color:#2563eb;}' +
                'p{margin:1em 0;}' +
                '#el-bar{position:fixed;top:0;left:0;right:0;background:#1e293b;color:#e2e8f0;padding:8px 16px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:13px;display:flex;align-items:center;justify-content:space-between;z-index:9999;box-shadow:0 2px 8px rgba(0,0,0,0.2);}' +
                '#el-bar a{color:#93c5fd;text-decoration:none;}' +
                '#el-bar button{background:#3b82f6;color:white;border:none;border-radius:4px;padding:4px 12px;cursor:pointer;font-size:12px;}' +
                '#el-bar button:hover{background:#2563eb;}' +
                '@media(prefers-color-scheme:dark){body{background:#1a1a2e;color:#e2e8f0;}h1{color:#f1f5f9;}a{color:#93c5fd;}}' +
                '</style></head><body>' +
                '<div id="el-bar"><span>📖 <a href="/">Extension Ladder</a></span><button id="el-exit">Exit Reader Mode</button></div>' +
                '<div style="margin-top:50px;">' +
                (article.title ? '<h1>' + article.title + '</h1>' : '') +
                (article.content || '<p>Could not extract article content.</p>') +
                '</div></body>';

            // Re-attach exit handler
            document.getElementById('el-exit').addEventListener('click', function() {
                active = false;
                document.documentElement.innerHTML = originalHTML;
                // Re-attach the toggle button listener
                setTimeout(function() {
                    var newBtn = document.getElementById('el-reader-toggle');
                    var newBar = document.getElementById('el-reader-bar');
                    if (newBtn) newBtn.addEventListener('click', toggle);
                }, 100);
            });
        } else {
            document.documentElement.innerHTML = originalHTML;
            setTimeout(function() {
                var newBtn = document.getElementById('el-reader-toggle');
                if (newBtn) newBtn.addEventListener('click', toggle);
            }, 100);
        }
    }
    btn.addEventListener('click', toggle);
    document.addEventListener('keydown', function(e) {
        if (e.key === 'r' && !e.ctrlKey && !e.metaKey && !e.altKey &&
            e.target.tagName !== 'INPUT' && e.target.tagName !== 'TEXTAREA') {
            toggle();
        }
    });
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
