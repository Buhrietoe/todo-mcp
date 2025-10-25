package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	socketPath := flag.String("socket", "/tmp/todo-mcp.sock", "The socket path to listen on.")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	server := mcp.NewServer(&mcp.Implementation{Name: "todo", Version: "1.0.0"}, &TodoServer{})
	if err := server.ListenAndServe(ctx, "unix", *socketPath); err != nil {
		log.Fatal(err)
	}
}
