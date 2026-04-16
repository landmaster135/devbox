# Notion Sync

Provides tools to patch Markdown content and web clips into Notion pages.

## Installation

```json
{
  "mcpServers": {
    "notion_sync": {
      "command": "/home/user/devbox/pkg/bin/mcp/linux_amd64/devbox-mcp-tools",
      "args": [
        "notion_sync"
      ],
      "env": {
        "NOTION_INTEGRATION_TOKEN": "YOUR_VALUE",
        "NOTION_ENDPOINT_URL": "YOUR_VALUE",
        "NOTION_ENDPOINT_URL_TO_PATCH_WEB_CLIP": "YOUR_VALUE"
      },
      "disabled": false,
      "autoApprove": [],
      "timeout": 180000,
      "requestTimeout": 210000
    }
  }
}
```

## License

Follows this repository.
