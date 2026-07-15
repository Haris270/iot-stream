package main

import (
	"context"
	"fmt"
	"log"
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

	// --- ADD THIS RESILIENT PING LOOP ---
	var pingErr error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		// Create a short 2-second context for each individual ping attempt
		pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr = pool.Ping(pingCtx)
		cancel()

		if pingErr == nil {
			fmt.Println("Database connection verified and ready!")
			break
		}

		log.Printf("Database not ready yet (Attempt %d/%d). Retrying in 2 seconds...", i, maxRetries)
		time.Sleep(2 * time.Second)
	}

	// If we exhausted all retries and still can't connect, THEN we panic
	if pingErr != nil {
		panic(fmt.Sprintf("Database failed to become ready: %v", pingErr))
	}

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
