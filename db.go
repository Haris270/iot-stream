package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitSchema(dbPool *pgxpool.Pool) {
	query := `
	CREATE TABLE IF NOT EXISTS sensor_data(
		sensorID		INTEGER,
		temperature		DOUBLE PRECISION,
		rpm				INTEGER,
		time			TIMESTAMPTZ NOT NULL
	);`

	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dbPool.Exec(startupCtx, query)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database schema: %v", err))
	}

	fmt.Println("Database successfully created!")
}

func InsertTelemetry(dbPool *pgxpool.Pool, data SensorData) {
	query := `
	INSERT INTO sensor_data (sensorID, temperature, rpm, time)
	VALUES ($1, $2, $3, $4);
	`
	insertCtx, insertCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer insertCancel()

	_, err := dbPool.Exec(insertCtx, query, data.Device_id, data.Temperature, data.Rpm, data.Time_stamp)
	if err != nil {
		log.Printf("Error in Insertion for Sensor %d\n", data.Device_id)
		return
	}
	fmt.Printf("Successfull insertion for Sensor %d\n", data.Device_id)
}
