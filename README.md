## Havuz TC Reverse Proxy
A serverless app that verifies Havuz users' license tokens and
reverse proxies requests to [Vekilio's Tunnel Collector](https://github.com/Vekilio/tunnel-collector).

> Note: This app is intended to work in serverless environments, having a single entrypoint (`Handler` func) without a `main` func. By passing `-tags "withMain"` flag to `go build`, it can be compiled with a `main` func that runs an HTTP server listening at `ADDR` env var (`:8090` if empty).
