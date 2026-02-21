<p align="center">
    <img src="assets/extension_ladder.png" width="100px">
</p>

<h1 align="center">Extension Ladder</h1>


### Why

Ads are everywhere and I'm tired of it

> **Disclaimer:** Don't use if interested in seeing ads

### How it works

It does a lot of the time, if it doesn't then please contribute to fix it

### Features
- [x] Bypass paywalls (Googlebot UA, headless Chrome, archive.org fallback chain)
- [x] Reader Mode — clean article view by default, press R for raw
- [x] Bot detection bypass (Cloudflare, PerimeterX)
- [x] Script stripping for JS-gated content
- [x] Domain-based rulesets
- [x] API & raw HTML endpoints
- [x] Docker, Helm, binary support
- [x] Basic Auth, no tracking
- [x] Removes most ads

### Reader Mode
Every proxied page opens in Reader Mode by default — a clean, distraction-free view with just the article text. No ads, no popups, no clutter. Press **R** or click **View Raw Page** to see the original site.

### Limitations
Not all sites work. Some block crawlers entirely, and archive.org won't have very recent articles. Reader Mode can often still extract the article content even when the page appears gated. If a site doesn't work, add a ruleset for it or contribute a fix.

## Installation

### Local Development
1) Install Go: `brew install go`
2) Install Chrome (required for headless browser fallback): `brew install --cask google-chrome`
3) Clone the repo: `git clone https://github.com/jon-fox/extension-ladder.git`
4) Run the start script:
```bash
cd extension-ladder
./local_start.sh
```
4) Open Browser (Default: http://localhost:8080)

> The `local_start.sh` script will automatically kill any existing process on port 8080 before starting.

### Binary
1) Download binary [here](https://github.com/jon-fox/extension-ladder/releases/latest)
2) Unpack and run the binary `./ladder -r https://t.ly/14PSf`
3) Open Browser (Default: http://localhost:8080)

### Docker
```bash
docker run -p 8080:8080 -d --env RULESET=https://t.ly/14PSf --name ladder ghcr.io/jon-fox/extension-ladder:latest
```

### Docker Compose
```bash
curl https://raw.githubusercontent.com/jon-fox/extension-ladder/main/docker-compose.yaml --output docker-compose.yaml
docker-compose up -d
```

### Helm
See [README.md](/helm-chart/README.md) in helm-chart sub-directory for more information.

## Usage

### Browser
1) Open Browser (Default: http://localhost:8080)
2) Enter URL
3) Press Enter

Or direct by appending the URL to the end of the proxy URL:
http://localhost:8080/https://www.example.com

Or create a bookmark with the following URL:
```javascript
javascript:window.location.href="http://localhost:8080/"+location.href
```

### API
```bash
curl -X GET "http://localhost:8080/api/https://www.example.com"
```

### RAW
http://localhost:8080/raw/https://www.example.com


### Running Ruleset
http://localhost:8080/ruleset

## Configuration

### Environment Variables

| Variable | Description | Value |
| --- | --- | --- |
| `PORT` | Port to listen on | `8080` |
| `PREFORK` | Spawn multiple server instances | `false` |
| `USER_AGENT` | User agent to emulate | `Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)` |
| `X_FORWARDED_FOR` | IP forwarder address | `66.249.66.1` |
| `USERPASS` | Enables Basic Auth, format `admin:123456` | `` |
| `LOG_URLS` | Log fetched URL's | `true` |
| `DISABLE_FORM` | Disables URL Form Frontpage | `false` |
| `FORM_PATH` | Path to custom Form HTML | `` |
| `RULESET` | Path or URL to a ruleset file, accepts local directories | `https://raw.githubusercontent.com/jon-fox/extension-ladder-rules/main/ruleset.yaml` or `/path/to/my/rules.yaml` or `/path/to/my/rules/` |
| `EXPOSE_RULESET` | Make your Ruleset available to other ladders | `true` |
| `ALLOWED_DOMAINS` | Comma separated list of allowed domains. Empty = no limitations | `` |
| `ALLOWED_DOMAINS_RULESET` | Allow Domains from Ruleset. false = no limitations | `false` |

`ALLOWED_DOMAINS` and `ALLOWED_DOMAINS_RULESET` are joined together. If both are empty, no limitations are applied.

### Ruleset

YAML-based rules per domain. Loaded from a file, directory, or URL on startup.

```yaml
- domain: example.com
  domains:
    - www.example.de
  headers:
    user-agent: Mozilla/5.0 ...
    x-forwarded-for: none
    referer: none
    cookie: privacy=1
  fetchStrategy: headless+archive  # default, headless, archive, headless+archive
  headlessWaitSeconds: 8
  stripScripts: true
  botDetectionPatterns:
    - "unusual activity"
  regexRules:
    - match: <div class="paywall">.*?</div>
      replace: ""
  injections:
    - position: head
      append: |
        <script>window.localStorage.clear();</script>
```

## Development

```bash
./local_start.sh
```

Or manually:
```bash
echo "dev" > handlers/VERSION
RULESET="./ruleset.yaml" go run cmd/main.go
```

### Optional: Live reloading development server with [cosmtrek/air](https://github.com/cosmtrek/air)

Install air according to the [installation instructions](https://github.com/cosmtrek/air#installation). 

Run a development server at http://localhost:8080:

```bash
air # or the path to air if you haven't added a path alias to your .bashrc or .zshrc
```

This project uses [pnpm](https://pnpm.io/) to build a stylesheet with the [Tailwind CSS](https://tailwindcss.com/) classes. For local development, if you modify styles in `form.html`, run `pnpm build` to generate a new stylesheet.
