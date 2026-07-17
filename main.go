package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	client := connectMQTT(os.Getenv("MQTT_BROKER"), "telemetry-engine")

	dsn := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Database Connection: %v", err))
	}

	// Ping the database and form connection
	if err := pingPool(pool); err != nil {
		panic(fmt.Sprintf("Database failed to become ready: %v", err))
	}

	fmt.Println("Database connection successful!")

	defer pool.Close()
	InitSchema(pool)

	app := &App{DB: pool}

	token := client.Subscribe("telemetry/sensor/+", 1, app.HandleTelemetry)
	if token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Failed to Subscribe to MQTT topic: %v", token.Error()))
	}

	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

AppLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt Signal received. Shutting down gracefully.")
			break AppLoop

		default:
			for i := 0; i < 100; i++ {
				wg.Add(1)

				go sendTelemetry(ctx, i, &wg, client)

			}
			time.Sleep(1 * time.Second)
		}

	}

	fmt.Println("Waiting for remaining telemetry routines to finish...")
	wg.Wait()
	fmt.Println("Clean shutdown complete.")

}
