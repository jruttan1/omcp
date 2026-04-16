# omcp

Point omcp at an OpenAPI spec and get a fully working MCP server. One binary, no runtime, instant startup.

## Install

```
brew tap jruttan1/omcp
brew install omcp
```

## Usage

**1. Run setup**
```
omcp init
```

You'll be prompted for:
- **API spec file** — path to your OpenAPI YAML file
- **Base URL** — where your server is running (`https://api.stripe.com`, `http://localhost:8080`, etc.)
- **API key** — optional, used as a Bearer token if your API requires auth

Then pick which endpoints to expose as tools. Your config is saved to `~/.omcp/config.json`.

**2. Add omcp to your MCP client**

```json
{
  "mcpServers": {
    "myapi": {
      "command": "omcp"
    }
  }
}
```

When your client starts, it spawns `omcp` which reads the saved config and registers your selected endpoints as tools — ready to call.

## What it supports

- Local servers (`http://localhost:8080`) or live APIs (`https://api.stripe.com`)
- Bearer token / API key auth
- Path params, query params, and request bodies
- Any MCP-compatible client: Claude Desktop, Cursor, Cline, Zed, and more
