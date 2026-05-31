package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"api/src"

	"github.com/joho/godotenv"
)

func main() {

	// Initialize Slog.
	logFile, err := os.OpenFile("tmp/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil { panic(err) }; defer logFile.Close()
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func (groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { return slog.Attr {} }
			return a
		},
	})
	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(slog.NewMultiHandler(stdoutHandler, fileHandler))
	slog.SetDefault(logger)
	slog.Info("Initialized Logger")

	// Initialize Env Config.
	if err := godotenv.Load("../.env"); err != nil { slog.Error("Failed to initialize Env Config") }
	config := src.LoadConfig()
	slog.Info("Initialized Config")

	// Initialize the App.
	app, err := src.NewApp(config)
	if err != nil { slog.Error("Failed to initialize App"); panic(err) }
	defer app.Close()
	app.InitSchema()
	slog.Info("Initialized App")

	// Initialize Payout.
	ctx, cancel := context.WithCancel(context.Background()); defer cancel()
	interval, err := time.ParseDuration(config.PAYOUT_INTERVAL + "s")
	if err != nil { slog.Error("bad PAYOUT_INTERVAL", "err", err); panic(err) }
	go src.RunPayoutLoop(ctx, app, interval, src.PayoutOptions {
		MinPending: 0,
	})

	// Initialize Router.
	mux := http.NewServeMux()
	// Middleware.
	var handler http.Handler = mux
	handler = src.RequestIDMiddleware(handler)
	handler = src.LoggerMiddleware(handler)
	handler = src.RecovererMiddleware(handler)

	// Routes.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello, World!") })
	mux.HandleFunc("POST /message", src.HandleRequest(app, src.HandleEndpointMessage))

	// Start Server.
	addr := ":" + config.PORT
	slog.Info("Initializing Server", "Address", addr, "PAYOUT_INTERVAL", config.PAYOUT_INTERVAL)
	if err := http.ListenAndServe(addr, handler); err != nil { slog.Error("Server Error", "err", err) }

}
