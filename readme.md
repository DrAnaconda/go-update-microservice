# Introduction

Custom API server to handle auto-update feature based which provided by `github.com/inconshreveable/go-update`.
API built via `github.com/labstack/echo/v4`.

# Flow

1. Client asks metadata of product
2. Authorization & rate limiting
3. API provides last update date of products
4. Client compares last update
5. Client downloads update file in case it is never
6. Client stores metadata of downloaded file

# Integration

1. Host DB server on postgres
2. Create `config.json` in root folder and add config values. Check `Config` struct for more details.
3. Fill DB with products. Refer to `Product` struct and table for more details.
4. Launch API server
5. Get code from `updater-client` package and integrate it in the client. Configure client.

# Features

1. Flexible configuration
2. Rate limiting 
3. Basic authorization
4. Certificate pinning in client
