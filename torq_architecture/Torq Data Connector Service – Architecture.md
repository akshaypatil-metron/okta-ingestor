### Torq Data Connector Service – Architecture Documentation





This document explains the architecture and working of the Torq Data Connector Service in a clear and structured way. The service is responsible for ingesting security events from third-party vendors and forwarding them to the Torq platform via webhooks.





**1. Purpose of the Service**



The Torq Data Connector Service ingests security events from third-party vendors such as Chronicle, CrowdStrike, and Splunk, and forwards those events to the Torq platform.

It supports two processing models:

• Polling

• Streaming



**2. High-Level Architecture Overview**



**Main components involved in the architecture:**



PubSub (Scheduler / Trigger)

Integration Service (Fetch secrets and configuration) .Integration Service fetches API keys, URLs, configuration via gRPC.

Registry and Factory Pattern

Connector (Vendor-specific logic)

Event Service

Webhook Dispatcher

Redis (State management, locks, dedupe)

Postgres (Offset storage for streaming)

gRPC Service (Update/Delete integration handling) In the Torq architecture, it is the "telephone line" between the Connector Service and the Integration Service.





**3. End-to-End Flow**





PubSub triggers the connector service with an Integration ID.

Integration Service fetches API keys, URLs, and configuration using the Integration ID.

Registry retrieves the correct connector using the Factory pattern.

Connector's Consume() function is executed.

Events are fetched from the vendor (Polling or Streaming).

Events are deduplicated.

Events are sent to Torq via Webhook.

After successful webhook response (200 OK), state is updated in Redis/Postgres.



**4. Polling Model (Example: Splunk)**



Polling is used when the vendor does not support streaming APIs.



**Polling Flow:**



PubSub triggers every 5 minutes.

Last timestamp is retrieved from Redis.

Vendor API is called using (last\_timestamp → current time).

Events are received.

Each event is checked for duplication.

Each event is sent to webhook individually.

If webhook returns 200 OK, Redis timestamp is updated.



**State Storage in Polling:**

• Redis stores last timestamp.

• Redis stores dedupe keys.



**5. Streaming Model (Example: Chronicle / CrowdStrike)**



Streaming is preferred when the vendor provides a real-time streaming API.



**Streaming Flow:**

PubSub triggers every 1 minute.

Redis distributed lock is acquired (2-minute lock).

Stream connection is opened.



**For each incoming event:**

&nbsp;  - Perform dedupe check.

&nbsp;  - Send event to webhook.

&nbsp;  - Save offset after successful webhook response.

After 2 minutes, stream is intentionally cancelled.

Next PubSub trigger restarts stream from last saved offset.



**State Storage in Streaming:**

• Redis stores distributed lock and dedupe keys.

• Postgres stores event offset.



**6. Reliability Principles**



**The system follows two critical rules:**

1\. Never lose events.

2\. Never send duplicate events.

Offset or timestamp is updated only after the webhook successfully processes the event. 

If the service crashes mid-way, events are replayed from the last saved offset, and dedupe logic prevents duplicate delivery.



**7. Integration Update and Delete Handling**



When a user updates integration configuration, a gRPC call is sent to the connector service.

A Redis flag is created for that integration ID.

Streaming watcher detects the flag.

Stream is cancelled immediately.

Next PubSub trigger starts stream with updated configuration.



**8. Observability and Monitoring**



Grafana dashboards for metrics.

Google Cloud logs for error and info logs.

Distributed tracing for debugging.

Metrics for stream health and processing visibility.



