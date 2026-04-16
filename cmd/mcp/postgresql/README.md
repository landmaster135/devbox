# PostgreSQL

Provides PostgreSQL tools for read-only queries, schema inspection, table listing, and table dump operations.

## Installation

```json
{
  "mcpServers": {
    "postgresql": {
      "command": "/home/user/devbox/pkg/bin/mcp/linux_amd64/devbox-mcp-tools",
      "args": [
        "postgresql"
      ],
      "env": {
        "POSTGRESQL_DATABASE_URL": "YOUR_URL"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

## License

Follows this repository.
