# IoT Telemetry Pipeline
<img width="1127" height="553" alt="Screenshot 2026-07-15 at 12 01 18 AM" src="https://github.com/user-attachments/assets/642ba39e-b2c8-4bca-b737-b6b1255590d2" />

## Overview
This is a containerized telemetry pipeline that simulates and visualizes industrial sensors in real time. THe project involves high throughput data generation, time series data aggregation and automated provisioning of the data through visual dashboard.

## Tech Stack
* **Data Generation (Go):** A custom Go application that simulates concurrent IoT sensor data (temperature, RPM) and subscribes to the Broker, consumes payload and writes to the database.
* **Message Broker (MQTT):** Utilizes **Eclipse Mosquitto** to set up a Publish/Subscribe (Pub/Sub) architecture where the messages containing sensor data are routed through MQTT broker on port `1883`.
* **Database (TimescaleDB):** A PostgreSQL engine optimized with TimescaleDB extensions that handles and stores high frequency write loads and enables time-based aggregation of data.
* **Visualization (Grafana):** Utilizes Grafana to create dashboard that retrieves data from the database and displays real time metrics.
* **Infrastructure (Docker):** Orchestrates all services, internal networks, volumes and environment variable configurations
---
<p align="center">
<img width="392" height="695" alt="image" src="https://github.com/user-attachments/assets/0168ded7-0bc6-43c0-89ab-5a38cecf371a" />
</p>

---

## Getting Started

### Prerequisites
* [Docker](https://docs.docker.com/get-started/get-docker/) and Docker Compose installed on your machine.
* Git

## Installation & Setup

**1. Clone the Repository**
```bash
git clone https://github.com/Haris270/iot-stream.git
cd iot-stream
```
**2. Configure the Environment**<br>
Create a .env file and add these environment variables:
```
POSTGRES_USER=admin
POSTGRES_PASSWORD=secret_password
POSTGRES_DB=telemetry_db

GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin
```
**3. Deploy the Pipeline**
```
docker compose up -d
```
**4. Clean Up (Teardown)**
```
docker compose down -v
```
