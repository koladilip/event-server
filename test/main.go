package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func publishMessages(ctx context.Context, userId string) error {
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down")
			return nil
		default:
			fmt.Println("Sending events to", userId)
		}
		values := map[string]string{"userId": userId,
			"payload": fmt.Sprintf("%s: %03d", userId, i)}
		json_data, err := json.Marshal(values)

		if err != nil {
			return err
		}

		_, err = http.Post("http://localhost:8080/publish", "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return publishMessages(gCtx, "user1")
	})
	g.Go(func() error {
		return publishMessages(gCtx, "user2")
	})

	g.Wait()
}
