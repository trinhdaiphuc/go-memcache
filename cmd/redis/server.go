package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/trinhdaiphuc/go-memcache/gomap"
	"github.com/trinhdaiphuc/go-memcache/hashmap"
	"github.com/trinhdaiphuc/go-memcache/internal/handler"
	"github.com/trinhdaiphuc/go-memcache/resp"
)

func main() {
	server := NewServer()

	server.ListenAndServe()
}

type Server struct {
	listener   net.Listener
	handler    handler.Map
	mapString  gomap.Map[string, string]
	hashString hashmap.HashMap[string, string]
	quit       chan os.Signal
}

func NewServer() *Server {
	s := &Server{
		quit:       make(chan os.Signal, 1),
		handler:    handler.NewMap(),
		mapString:  gomap.NewMap[string, string](),
		hashString: hashmap.NewHashMap[string, string](),
	}

	signal.Notify(s.quit, os.Interrupt)
	return s
}

func (s *Server) ListenAndServe() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	s.listener = l

	go s.run()

	fmt.Println("Server is running on port 6379")

	<-s.quit

	fmt.Println("Server is shutting down")
	err = s.listener.Close()
	if err != nil {
		fmt.Println("Error closing listener: ", err.Error())
	}

	fmt.Println("Server stopped")
}

func (s *Server) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go s.processConnection(conn)
	}
}

func (s *Server) processConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)

		respCmd := resp.NewRESPCommand()
		err := respCmd.Read(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
				break
			}
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}

		// Write back to the client
		if respCmd.Array != nil {
			args := respCmd.Array.Expressions
			if len(args) == 0 {
				conn.Write([]byte(resp.NewErrorExpression("No command provided").Serialize()))
				continue
			}
			cmd := args[0].Value().(string)
			h, ok := s.handler[cmd]
			if !ok {
				conn.Write([]byte(resp.NewErrorExpression("Unknown command").Serialize()))
				continue
			}
			ctx := handler.NewContext(s.mapString, s.hashString)
			result := h.Handle(ctx, args[1:])
			conn.Write([]byte(result.Serialize()))
			continue
		}
	}
}
