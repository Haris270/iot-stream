package main

import (
	"encoding/json"
	"testing"
)

func TestGenerateSensorPayload(t *testing.T) {
	var test_ID int = 42

	encoded_payload, err := generateSensorPayload(test_ID)

	if err != nil {
		t.Fatalf("Test Failed: Expected no error, got %v", err)
	}

	var result SensorData
	err = json.Unmarshal(encoded_payload, &result)

	if err != nil {
		t.Fatalf("Test Failed: Expected no error, got %v", err)
	}

	if result.Device_id != test_ID {
		t.Fatalf("Test Failed: Expected Sensor ID: %d, instead got: %d", test_ID, result.Device_id)
	}

	if result.Temperature < 25.0 || result.Temperature >= 110.0 {
		t.Errorf("Expected Temperature between 25 and 50, got %f", result.Temperature)
	}

	if result.Rpm < 50 || result.Rpm >= 75 {
		t.Errorf("Expected RPM between 50 and 75, got %d", result.Rpm)
	}
}
