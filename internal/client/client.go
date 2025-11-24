package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

const TextMessageType int = 1

type Client struct {
	Conn *websocket.Conn
}

func StartClient() error {
	url := "ws://localhost:9000/ws/connect"
	fmt.Println("Connecting to", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println("Connected successfully!")
	client := &Client{
		Conn: conn,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	readDone := make(chan struct{})
	scanDone := make(chan struct{})
	errChan := make(chan error, 2)

	go func() {
		defer close(readDone)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := client.Conn.ReadMessage()
				if err != nil {
					fmt.Printf("connection closed or read error: %v\n", err)
					select {
					case errChan <- fmt.Errorf("failed to read message: %w", err):
					default:
					}
					cancel()
					return
				}
				fmt.Println(string(message))
			}
		}
	}()

	go func() {
		defer close(scanDone)

		lines := make(chan string)

		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				text := scanner.Text()
				select {
				case lines <- text:
				case <-ctx.Done():
					return
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case errChan <- fmt.Errorf("error occur while scanner.Scan(): %w", err):
				default:
				}
			}
			close(lines)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-lines:
				if !ok {
					cancel()
					return
				}
				if err := conn.WriteMessage(TextMessageType, []byte(line)); err != nil {
					select {
					case errChan <- fmt.Errorf("failed to write message: %w", err):
					default:
					}
					cancel()
					return
				}
			}
		}
	}()

	select {
	case <-quit:
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		cancel()
	case <-readDone:
		cancel()
	case <-scanDone:
		cancel()
	}

	<-readDone
	<-scanDone

	select {
	case err = <-errChan:
	default:
	}

	if err != nil {
		fmt.Printf("client exiting due to error: %v\n", err)
		return err
	}

	fmt.Println("client exiting normally.")
	return nil
}
