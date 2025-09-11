# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

VibeKickstart is a hybrid web application with:

- **Go backend** (`vibekickstart/` directory) - HTTP server with WebSocket support for real-time communication
- **JavaScript frontend** (`src/` directory) - ES6 modules bundled with esbuild
- **Static assets** (`public/` directory) - CSS and other static files

## Development Commands

### Building Assets

- `npm run build-dev` - Development build with sourcemaps and no minification
- `npm run build-prod` - Production build with minification

### Using Mise (Recommended)

The project uses mise for task management and dependencies:

- `mise run watch` - Watch for changes and rebuild automatically
- `mise run run` - Build assets (dev) and run the Go server
- `mise run build` - Build production assets and compile Go binary
- `mise run lint` - Run ESLint and Prettier on JavaScript code

### Manual Go Commands

- `cd vibekickstart && go run .` - Run the Go server directly
- `cd vibekickstart && go build -o vibekickstart` - Build Go binary

### Linting and Formatting

- `npx eslint src/` - Lint JavaScript
- `npx prettier src/ --write` - Format JavaScript

## Architecture

### Backend (Go)

- **Main server** (`app.go`) - HTTP server with embedded static files, template rendering, and cache busting
- **WebSocket manager** (`wsconn.go`) - Connection pooling and broadcasting for real-time features
- **Custom logging** (`multiline.go`) - Pretty-printed log formatting

The Go server embeds both the `dist/` directory (built assets) and `tmpl/` templates using `//go:embed`.

### Frontend (JavaScript)  

- **Entry point** (`src/main.js`) - Uses ReconnectingWebSocket for WebSocket connections
- **Build system** - esbuild with separate dev/prod configurations
- **Browser targets** - Modern browsers (Chrome 58+, Firefox 57+, Safari 11+, Edge 16+)

### Static Assets

- CSS and other assets in `public/` are copied to `dist/` during build
- The Go server serves static files from `/static/` with long-term caching

## Key Implementation Details

- Templates use Go's `html/template` with cache-busting via timestamp
- WebSocket connections are managed per-name with broadcasting capabilities  
- ESBuild handles module bundling and supports both development and production modes
- Static files are served with aggressive caching (1 year) via cache-busting query parameters
