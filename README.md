# VibeKickstart

A template project for "vibe coding" - an opinionated starting point for a a hybrid web application that serves as a starting point for rapid prototyping and experimentation

## Features

- **Go Backend** - Fast HTTP server with WebSocket support for real-time communication
- **JavaScript Frontend** - Modern ES6 modules with automatic reconnection
- **Live Reload** - Watch mode for automatic rebuilding during development
- **Real-time Counter** - Demo WebSocket feature with synchronized counter across all connected clients

## Quick Start

### Prerequisites

- [mise](https://mise.jdx.dev/) (recommended) or:
  - Go 1.24+
  - Node.js 22+

### Development

```bash
# Install dependencies
mise install

# Start development server with watch mode
mise run watch

```

Visit <http://localhost:8750>

### Production

```bash
# Build production assets and binary
mise run build

# Run the binary
./vibekickstart/vibekickstart
```

### Docker

```bash
# Build Docker image
docker build -t vibekickstart .

# Run container
docker run -p 8750:8750 vibekickstart

# Or run in background
docker run -d -p 8750:8750 --name vibekickstart vibekickstart
```

## Demo Feature

The application includes a simple WebSocket counter demonstration:

- A button is displayed in the center of the screen
- The button shows a count that increments over time
- When any user clicks the button, it broadcasts a reset to all connected clients
- All clients see the counter reset in real-time

This demonstrates the real-time synchronization capabilities perfect for collaborative applications, games, or live updates.

## Architecture

### Single Binary Deployment

VibeKickstart compiles everything into a **single, self-contained binary** with no external dependencies:

- **Static Assets Embedded**: All CSS, JavaScript, and other assets are embedded using Go's `//go:embed` directive
- **Templates Embedded**: HTML templates are compiled directly into the binary
- **Zero Dependencies**: The final Docker image uses `scratch` as the base, containing only the binary
- **Tiny Footprint**: Final Docker image is typically under 20MB
- **Static Linking**: Binary is statically linked with no libc dependencies

This approach means you can deploy the application anywhere that can run a Linux binary - no runtime dependencies, no external files, just copy and run.

### Backend (Go)

- HTTP server with embedded static files
- WebSocket connection management with broadcasting
- Template rendering with cache-busting
- Custom pretty-printed logging
- Compiled to single static binary

### Frontend (JavaScript)

- ReconnectingWebSocket for automatic reconnection
- ES6 modules bundled with esbuild
- Modern browser targets (Chrome 58+, Firefox 57+, Safari 11+, Edge 16+)
- Assets embedded in Go binary at build time

### Build System

- **Multi-stage Docker build**: JavaScript assets built first, then embedded in Go binary
- **esbuild**: Fast JavaScript bundling and minification
- **mise**: Task management and tool versions for local development
- **Static compilation**: `CGO_ENABLED=0` ensures no C dependencies

## Development Commands

- `mise run watch` - Watch for changes and rebuild automatically
- `mise run run` - Build assets and run server once
- `mise run build` - Build production assets and Go binary
- `mise run lint` - Run linting and formatting
- `npm run build` - Build JavaScript assets only

## Project Structure

```
├── vibekickstart/          # Go backend
│   ├── app.go             # Main HTTP server
│   ├── wsconn.go          # WebSocket connection manager
│   ├── multiline.go       # Custom log formatting
│   └── tmpl/              # HTML templates
├── src/                   # JavaScript frontend
│   └── main.js           # Entry point
├── public/               # Static assets (CSS, images, etc.)
└── dist/                # Built assets (generated)
```

## License

MIT
