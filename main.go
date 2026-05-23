package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/webview/webview_go"
)

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}
	defer closeDB()

	mux := setupRoutes()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	defer listener.Close()

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		log.Printf("Serving on port %d\n", listener.Addr().(*net.TCPAddr).Port)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Sistem Pendukung Keputusan – Pemilihan Penyedia Layanan Cloud")
	w.SetSize(900, 600, webview.HintMin)
	w.Navigate("http://" + listener.Addr().String())
	w.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
