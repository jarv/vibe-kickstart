# Build stage
FROM node:22-alpine AS js-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY src/ src/
COPY public/ public/
COPY build.js .
RUN npm run build

# Go build stage
FROM golang:1.25-alpine AS go-builder
WORKDIR /app
COPY vibekickstart/go.mod vibekickstart/go.sum ./
RUN go mod download
COPY vibekickstart/*.go ./
COPY vibekickstart/tmpl/ tmpl/
COPY --from=js-builder /app/dist/ dist/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -o vibekickstart

# Final stage - scratch for minimal image
FROM scratch
COPY --from=go-builder /app/vibekickstart /vibekickstart
EXPOSE 8910
ENTRYPOINT ["/vibekickstart"]
CMD ["-addr", "0.0.0.0:8910"]
