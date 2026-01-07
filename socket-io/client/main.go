package main

import (
	"context"
	"fmt"
	"log"
	"time"

	socketio "github.com/maldikhan/go.socket.io/socket.io/v5/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := socketio.NewClient(
		socketio.WithRawURL("http://localhost:8080/socket.io/"),
		//socketio.WithLogger(&socketio.DefaultLogger{}),
	)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	client.On("eventname", func(data []byte) {
		fmt.Printf("Received message: %s\n", string(data))
	})

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}

	if err := client.Emit("hello", "world"); err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	<-ctx.Done()
}
