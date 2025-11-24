package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type BroadcastServer struct {
	HttpServer *http.Server
	Clients    map[*websocket.Conn]bool
	mu         sync.Mutex
}

func StartBCServer(addr string) error {
	bcServer := &BroadcastServer{
		Clients: make(map[*websocket.Conn]bool),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", bcServer.handleConnect)
	bcServer.HttpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		if err := bcServer.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("bcServer ListenAndServe error: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Printf("shutting down all clients...\n")
	bcServer.closeAllClients()

	fmt.Printf("closing server...\n")
	return bcServer.HttpServer.Shutdown(ctx)
}

var upgrader = websocket.Upgrader{}

func (s *BroadcastServer) handleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("failed to upgrade to WebSocket: %v\n", err)
		return
	}
	fmt.Printf("New client connected from %s\n", r.RemoteAddr)

	s.mu.Lock()
	s.Clients[conn] = true
	s.mu.Unlock()

	s.handleClient(conn)
}

func (s *BroadcastServer) handleClient(conn *websocket.Conn) {
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			s.mu.Lock()
			delete(s.Clients, conn)
			s.mu.Unlock()
			break
		}
		fmt.Println(string(message))
		go s.broadcastMessage(mt, message, conn)
	}
}

func (s *BroadcastServer) broadcastMessage(mt int, message []byte, exceptedConn *websocket.Conn) {
	s.mu.Lock()
	// conns 记录除 exceptedConn 外的所有 conn
	var conns []*websocket.Conn
	for conn, connected := range s.Clients {
		if !connected || conn == exceptedConn {
			continue
		}
		conns = append(conns, conn)
	}
	s.mu.Unlock()
	for _, conn := range conns {
		if err := conn.WriteMessage(mt, message); err != nil {
			fmt.Printf("failed to write meeage to client: %v\n", err)
			s.mu.Lock()
			delete(s.Clients, conn)
			s.mu.Unlock()
		}
	}
}

func (s *BroadcastServer) closeAllClients() {
	s.mu.Lock()
	for conn := range s.Clients {
		conn.Close()
	}
	s.mu.Unlock()
}
