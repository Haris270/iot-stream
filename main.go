package main

import (
	"context"

	"sync"
	"time"
)

func main() {

	client := connectMQTT("tcp://localhost:1883", "telemetry-engine")

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go sendTelemetry(ctx, i, &wg, client) //pass the MQTT client to the sendTelemetry func, so it can call client.Publish())
	}

	wg.Wait()

}
