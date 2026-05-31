package main

import (
	"context"

	"fmt"

	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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
