package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/koladilip/event-server/utils"
)

func publishMessages(userId string) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		values := map[string]string{"userId": userId,
			"payload": fmt.Sprintf("%s: %s", userId, time.Now().String())}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}

		_, err = http.Post("http://localhost:8080/publish", "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Println(err)
		}
		utils.WaitForRandomPeriod()
	}
}

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	go publishMessages("user1")
	go publishMessages("user2")
	wg.Wait()
}
