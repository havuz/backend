## Havuz TC Reverse Proxy

A serverless app that verifies Havuz users' license tokens and
reverse proxies requests to [Vekilio's Tunnel Collector](https://github.com/Vekilio/tunnel-collector).

> This app is intended to work in serverless environments, having a single entrypoint (`Handler` func) without a `main` func. By passing `-tags "with_main"` flag to `go build`, it can be compiled with a `main` func that runs an HTTP server listening at `ADDR` env var (`:8090` if empty).

----

### Configuration
#### Environment Variables
```json
// Google Sheets URL, exported as CSV. Must be public.
// Must have CSV headers corresponding to type User.
// Location headers are followed by HTTP client.
// Example: https://bit.ly/2IndTaa
"SHEET_URL": "",

// NaCl crypto_sign public key to verify and decrypt the digest.
// Must be encoded with base64std, not url variant.
"PUB_KEY": "",

// HTTP URL of TC for reverse proxying.
// User-Pass may be omitted. Only relative root dir ("/") will be hit.
// Example: https://usern:passw@collector.example.com
//          https://collector.another.com/another-path
"TC_URL": ""
```
