package main

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func connectMQTT(brokerURL string, clientID string) MQTT.Client {

	opts := MQTT.NewClientOptions() //creates a ClientOptionsType with default values e.g Port:1883, KeepAlive:30
	opts.AddBroker(brokerURL)       // the MQTT Broker (the central hub in a publish/subscribe messaging system, responsible for receiving message from publishers, filtering them by topic and delivering them to subscribers)
	opts.SetClientID(clientID)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Mosquitto Broker connected Successfully!")
	return client

}
