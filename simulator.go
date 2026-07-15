package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//struct containing the device data

type SensorData struct {
	Device_id   int       `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Rpm         int       `json:"rpm"`
	Time_stamp  time.Time `json:"time_stamp"`
}

func generateSensorPayload(id int) ([]byte, error) {
	baseTemp := 25.0 + float64(id%70)
	temp := baseTemp + (rand.Float64() * 15.0)
	var data SensorData = SensorData{
		id,
		temp,
		rand.Intn(25) + 50,
		time.Now()}

	return json.Marshal(data)
}

func sendTelemetry(ctx context.Context, id int, wg *sync.WaitGroup, client MQTT.Client) {
	defer wg.Done()
	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Sensor %d shutting down\n", id)
			return

		default:
			marshalledByte, err := generateSensorPayload(id)
			if err != nil {
				fmt.Println("Marshal Error")
				continue
			}

			topic := fmt.Sprintf("telemetry/sensor/%d", id)

			client.Publish(topic, 0, false, marshalledByte)

			// time.Sleep(1 * time.Second)
		}

	}

}
