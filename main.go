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

func sendTelemetry(ctx context.Context, id int, wg *sync.WaitGroup, client MQTT.Client) {
	defer wg.Done()
	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Sensor %d shutting down\n", id)
			return

		default:
			var data SensorData = SensorData{
				id,
				(rand.Float64() * 25) + 25,
				rand.Intn(25) + 50,
				time.Now()}

			topic := fmt.Sprintf("factory/building/device/%d", id)

			marshalledByte, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Marshal Error")
				continue
			}

			client.Publish(topic, 0, false, marshalledByte)

			time.Sleep(1 * time.Second)
		}

	}

}

func main() {
	opts := MQTT.NewClientOptions()        //creates a ClientOptionsType with default values e.g Port:1883, KeepAlive:30
	opts.AddBroker("tcp://localhost:1883") // the MQTT Broker (the central hub in a publish/subscribe messaging system, responsible for receiving message from publishers, filtering them by topic and delivering them to subscribers)
	opts.SetClientID("telemetry-engine")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Mosquitto Broker connected Successfully!")

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go sendTelemetry(ctx, i, &wg, client) //pass the MQTT client to the sendTelemetry func, so it can call client.Publish())
	}

	wg.Wait()

}
