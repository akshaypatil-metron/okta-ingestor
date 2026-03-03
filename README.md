
# Okta System Logs Ingestor

![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Coverage](https://img.shields.io/badge/coverage-90%25-brightgreen)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## Description
The **Okta System Logs Ingestor** is a robust, stateful Go service designed to continuously fetch System Logs from Okta via API pagination, store them securely in MongoDB, and forward them to an external Webhook. 

Built specifically for reliable API integration with cybersecurity tools and SIEM/SOAR platforms, this service guarantees zero data loss, exact-once processing, and strict adherence to UTC timestamps for accurate security event correlation.

**Key Features:**
* **Stateful Pagination:** Tracks the latest Okta cursor (`next_url`) in MongoDB to resume exactly where it left off after a restart or failure.
* **Resilient Polling Loop:** Designed to handle Okta API Rate Limiting (HTTP 429) gracefully.
* **Zero Duplicates:** Uses MongoDB `_id` mapping to the Okta log `uuid`, silently ignoring `11000` duplicate key errors.
* **Webhook Forwarding:** Streams validated log batches in real-time to external endpoints.
* **Graceful Shutdown & Health Checks:** Exposes a `/health` endpoint for orchestrators and safely finalizes database connections during `SIGINT`/`SIGTERM`.

## Visuals
### Architecture & Ingestion Flow
Below is the system's execution flow, detailing the pagination loop, deduplication strategy, and webhook retry mechanisms.

![Okta System Logs Ingestion Flow](./out/sequence/sequence.png)

## Installation

### Requirements
* **Go** (v1.19 or higher)
* **MongoDB** (Running locally on port `27017` or a remote cluster)
* **Okta API Token** (Requires `System Logs: Read` permissions)

### Steps
1. **Clone the repository:**
   ```bash
   git clone <your-repository-url>
   cd okta-ingestor
   ```

2. **Download dependencies:**
   ```bash
   go mod tidy
   ```


3. **Configure the environment:**
Open `internal/config/config.go` and update your target environment details. **Note: The `StartDate` must be in strict UTC format (Zulu time).**
   ```go
   OktaDomain:   "[https://trial-5020092.okta.com](https://trial-5020092.okta.com)",
   OktaToken:    "YOUR_OKTA_API_TOKEN", 
   MongoURI:     "mongodb://localhost:27017",
   WebhookURL:   "[https://my-okta-logs.free.beeceptor.com](https://my-okta-logs.free.beeceptor.com)",
   PollInterval: 30 * time.Second,
   BatchSize:    100,
   StartDate:    "2026-02-01T00:00:00.000Z", 
   ```



## Usage

Start the ingestor application from the root directory:

```bash
go run main.go

```

**Expected Console Output:**

```text
2026/03/02 10:41:00 Starting health check server on :8080
2026/03/02 10:41:00 Starting scheduler. Polling every 30s...
2026/03/02 10:41:00 Asking Okta for next batch...
2026/03/02 10:41:01 Successfully fetched 100 logs. Latest Timestamp: 2026-02-01T12:00:00.000Z
2026/03/02 10:41:01 Saved to MongoDB successfully.
...
2026/03/02 10:41:44 No new logs found. Caught up to present time!

```

**Health Check Verification:**

```bash
curl http://localhost:8080/health
# {"status": "ok"}
```

## Support

For issues regarding API rate limits, database connectivity, or Okta token expiration, please check the console logs first. If the issue persists, open a ticket in the repository's issue tracker with the attached error logs and timestamps.

## Roadmap

* [ ] Migrate configuration from hardcoded Go structs to `.env` or `.yaml` file loading (e.g., using Viper).
* [ ] Implement exponential backoff for Webhook `Send()` failures.
* [ ] Add Prometheus metrics to track ingestion lag and rate-limit hits.
* [ ] Containerize the application using Docker and provide a `docker-compose.yml` for simplified local testing.

## Contributing

Contributions are welcome, especially those that improve the ingestor's resilience or data parsing capabilities.

1. Fork the repository and create your feature branch.
2. Ensure you write and pass **Unit Test Cases (UTCs)** for any new logic.
3. Verify your changes maintain a high test coverage threshold:
   ```bash
   # Run isolated Unit Tests
   go test -short ./...

   # Run Integration Tests (Requires local MongoDB instance)
   go test -v -tags=integration -coverprofile coverage.out ./...
   go tool cover -html coverage.out

   ```


4. Commit your changes and open a Pull Request.

## Authors and acknowledgment

* Built by the core engineering team for the Torq Integration Project.
* Special thanks to the architecture team for the flow designs.

## License

This project is licensed under the MIT License. See the `LICENSE` file for full details.

## Project status

**Active Development.** The core polling loop, MongoDB state management, and Okta API pagination logic are fully operational. Current focus is expanding test coverage and environment variable integration.
