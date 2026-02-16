# Weather Notificator

Sends weather forecast notifications to Discord for a specified city.

## Installation

```json
{
  "mcpServers": {
    "weather_notificator": {
      "command": "/home/user/devbox/pkg/bin/mcp/linux_amd64/devbox-mcp-tools",
      "args": [
        "weather_notificator"
      ],
      "env": {
        "OPENWEATHER_API_KEY": "YOUR_VALUE",
        "DISCORD_WEBHOOK_URL": "YOUR_VALUE"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

## License

Follows this repository.
