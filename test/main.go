package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koladilip/event-server/utils"
	"golang.org/x/sync/errgroup"
)

func publishMessages(userId string) error {
	for i := 0; i < 100; i++ {
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
		utils.WaitForRandomPeriod()
	}
	return nil
}

func main() {
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return publishMessages("user1")
	})
	g.Go(func() error {
		return publishMessages("user2")
	})
	g.Wait()
}
