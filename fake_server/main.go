package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const duration = 15

var ipAddresses = generateIPAddresses()

func generateIPAddresses() []string {
	ips := make([]string, 100)
	for i := 1; i <= 100; i++ {
		ips[i-1] = fmt.Sprintf("192.168.1.%d", i)
	}
	return ips
}

type Payload struct {
	IP       string `json:"ip"`
	Duration int    `json:"duration"`
	Time     int    `json:"time"`
}

func sendPostRequest(url string, ip string) {
	payload := Payload{}
	payload.IP = ip
	payload.Duration = duration
	payload.Time = int(time.Now().UnixMilli())
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("IP: %s, Status Code: %d\n", ip, resp.StatusCode)
	bodyBytes := new(bytes.Buffer)
	bodyBytes.ReadFrom(resp.Body)
}

func main() {
	url := "http://localhost:8081/healthcheck"

	ticker := time.NewTicker(duration * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < 10; i++ {
				go func() {
					ip := ipAddresses[rand.Intn(len(ipAddresses))]
					sendPostRequest(url, ip)
				}()
			}
		}
	}
}
