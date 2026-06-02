package main

import (
	"encoding/json"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB *pgxpool.Pool
}

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

func (a *App) HandleTelemetry(client MQTT.Client, msg MQTT.Message) {
	msg_payload := msg.Payload()

	var data SensorData
	err := json.Unmarshal(msg_payload, &data)
	if err != nil {
		log.Printf("Error Unmarshalling data: %v\n", err)
		return
	}
	InsertTelemetry(a.DB, data)

}
