# VibeKickstart

A template project for my personal "vibe coding"

**[view the demo](https://vk.jarv.org)**

## Features

- **Go Backend** - Fast HTTP server with WebSocket support
- **JavaScript Frontend** - Modern ES6 modules (plain JS, no typescript) with esbuild
- **Single binary** - compiles everything into a **single, self-contained binary** with no external dependencies

## Why a kickstart template?

I found myself starting with the same project scaffolding for multiple projects, a Go backend with frontend assets compiled into a single binary.
This kickstart template gives you just that, and not much more.
Most of the time I am using Websockets so this includes minimal Websocket support.

## Projects that started with this template

_If you use this and have a project to add, please feel free to send an MR_

- [stardewar.com](https://stardewar.com)
- [flyemoji.com](https://flyemoji.com)

## Dependencies

There aren't many dependencies for this project but those that are required are updated automatically with [Renovate](https://github.com/renovatebot/renovate).
Everything here should have latest versions.

## Testing

To ensure that nothing is broken for when dependencies are updated, there is a [playwrite](https://github.com/microsoft/playwright) integration test that is run on every change that ensures that the root page loads and that we can send and receive messages using the backend.

## Quick Start

### Compile and run local

1. Install [mise](https://mise.jdx.dev/)
2. `mise install`
3. `mise run watch`
4. Visit <http://localhost:8910>

### Run with Docker

```bash
# Build Docker image
docker build -t vibekickstart .

# Run container
docker run -p 8910:8910 vibekickstart

# Or run in background
docker run -d -p 8910:8910 --name vibekickstart vibekickstart
```

## Demo site

- A button is displayed in the center of the screen
- The button shows a count that increments over time
- When any user clicks the button, it broadcasts a reset to all connected clients
- All clients see the counter reset in real-time

## Architecture

### Backend (Go)

- HTTP server with embedded static files
- WebSocket connection management with broadcasting
- Template rendering with cache-busting
- Custom pretty-printed logging
- Compiled to single static binary

### Frontend (JavaScript)

- ReconnectingWebSocket for automatic websocket reconnection
- ES6 modules bundled with esbuild
- Modern browser targets (Chrome 58+, Firefox 57+, Safari 11+, Edge 16+)
- Assets embedded in Go binary at build time

### Deployment

The app can be deployed as a Docker container or as a single binary on a VM.

## Development

### Vibe Coding

- The entry point is `src/main.js` for the frontend JS code.
- For the server, routes are defined in `app.go`.
- If you require assets (images, etc) place them in `public/`.

### Commands

- `mise run watch` - Watch for changes and rebuild automatically
- `mise run run` - Build assets and run server once
- `mise run build` - Build production assets and Go binary
- `mise run lint` - Run linting and formatting
- `mise run test` - Run integration tests

## License

MIT
