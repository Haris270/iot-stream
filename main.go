package main

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type SensorData struct {
	Device_id   int       `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Rpm         int       `json:"rpm"`
	Time_stamp  time.Time `json:"time_stamp"`
}

func sendTelemetry(ctx context.Context, id int, wg *sync.WaitGroup) {
	var data SensorData = SensorData{
		id,
		(rand.Float64() * 25) + 25,
		rand.Intn(25) + 50,
		time.Now()}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go sendTelemetry(ctx, rand.Intn(100), &wg)
	}

	wg.Wait()

}
