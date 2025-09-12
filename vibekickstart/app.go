package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	oneYearCacheControl string = "public, max-age=31536000"
)

var (
	addr      = flag.String("addr", "localhost:8910", "http service address")
	jsonLog   = flag.Bool("json", false, "use json logs")
	templates = template.Must(template.ParseFS(tmplFiles, "tmpl/*.tmpl"))
	cacheBust = time.Now().Format("20060102150405")

	//go:embed dist/*
	distFiles embed.FS

	//go:embed tmpl/*.tmpl
	tmplFiles embed.FS

	logger    *slog.Logger
	cm        *ConnectionManager
	counter   int
	counterMu sync.RWMutex

	startTime = time.Now()
)

func setupLogging() {
	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	var slogHandler slog.Handler
	if *jsonLog {
		slogHandler = slog.NewJSONHandler(os.Stdout, slogOpts)
	} else {
		slogHandler = &MultilineHandler{Writer: os.Stdout}
	}
	logger = slog.New(slogHandler)
}

func main() {
	flag.Parse()
	setupLogging()

	// web server
	mux := http.NewServeMux()

	// ws connection manager
	cm = NewConnectionManager()

	// File server for static assets
	distFS, _ := fs.Sub(distFiles, "dist")
	mux.Handle("GET /static/", cacheControlMiddleware(
		http.StripPrefix("/static/", http.FileServer(http.FS(distFS))),
		oneYearCacheControl,
	))

	// WebSocket endpoint
	mux.HandleFunc("GET /ws", handleWebSocket)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := struct{ CacheBust string }{CacheBust: cacheBust}
		if err := templates.ExecuteTemplate(w, "index.html.tmpl", d); err != nil {
			logger.Error("Error executing template", "err", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	})

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}

	// Start counter increment goroutine
	go startCounterIncrement()

	logger.Info("Server started", "addr", "http://"+*addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Unable to setup listener", "err", err)
		os.Exit(1)
	}
}

type CounterMessage struct {
	Type    string `json:"type"`
	Counter int    `json:"counter"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		logger.Error("Failed to accept WebSocket connection", "err", err)
		return
	}
	defer conn.CloseNow()

	ctx := context.Background()
	clientID := r.Header.Get("X-Forwarded-For")
	if clientID == "" {
		clientID = r.RemoteAddr
	}

	cm.Add("counter", conn)
	defer cm.Remove("counter", conn)

	// Send current counter value to new client
	currentCounter := getCounter()
	welcomeMsg := CounterMessage{
		Type:    "update",
		Counter: currentCounter,
	}
	welcomeData, _ := json.Marshal(welcomeMsg)
	conn.Write(ctx, websocket.MessageText, welcomeData)

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			return
		}

		var msg CounterMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Error("Failed to parse message", "err", err)
			continue
		}

		if msg.Type == "reset" {
			resetCounter()
			logger.Info("Counter reset by client", "client", clientID)
		}
	}
}

func getCounter() int {
	counterMu.RLock()
	defer counterMu.RUnlock()
	return counter
}

func resetCounter() {
	counterMu.Lock()
	counter = 0
	counterMu.Unlock()
	broadcastCounter()
}

func incrementCounter() {
	counterMu.Lock()
	counter++
	counterMu.Unlock()
	broadcastCounter()
}

func broadcastCounter() {
	currentCounter := getCounter()
	msg := CounterMessage{
		Type:    "update",
		Counter: currentCounter,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal counter message", "err", err)
		return
	}

	ctx := context.Background()
	cm.BroadcastAll(ctx, data)
}

func startCounterIncrement() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		incrementCounter()
	}
}

func cacheControlMiddleware(next http.Handler, cacheControl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cacheControl)
		next.ServeHTTP(w, r)
	})
}
