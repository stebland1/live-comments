# Live Comments

A distributed live-comments backend written in Go, designed to explore realtime systems, service boundaries, pub/sub messaging, and layered application architecture.

The project consists of two separate HTTP services:

* **write-api** → accepts and persists comments
* **stream-api** → streams live comments to connected clients using Server-Sent Events (SSE)

Redis is used for realtime fan-out via Pub/Sub, while PostgreSQL acts as the durable source of truth.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Features](#features)
4. [Project Structure](#project-structure)
5. [How It Works](#how-it-works)
6. [What I Learned](#what-i-learned)
7. [Running the Project](#running-the-project)
8. [API Endpoints](#api-endpoints)
9. [Known Limitations](#known-limitations)
10. [Future Improvements](#future-improvements)

---

## Overview

This project was built to better understand how realtime systems are structured internally.

Rather than building a monolithic application, the system is split into two separate services:

* a **write service** responsible for validating and persisting comments
* a **stream service** responsible for realtime delivery

The services communicate indirectly through Redis Pub/Sub channels.

The architecture intentionally separates:

* persistence concerns
* transport concerns
* realtime messaging concerns
* application/business logic

to mirror patterns commonly used in larger distributed systems.

---

## Architecture

```text
Client
   │
   ├── POST /comment ───────────────► write-api
   │                                      │
   │                                      ├── Persist to PostgreSQL
   │                                      │
   │                                      └── Publish event to Redis
   │
   └── GET /stream/:videoID ───────► stream-api
                                          │
                                          └── Subscribe to Redis Pub/Sub
                                                     │
                                                     ▼
                                             Live SSE Stream
```

### Components

#### PostgreSQL

Used as the durable source of truth for persisted comments.

#### Redis Pub/Sub

Used for transient realtime message fan-out between services.

#### write-api

Handles:

* HTTP comment creation
* validation
* persistence
* publishing events to Redis

#### stream-api

Handles:

* SSE connections
* Redis subscriptions
* broadcasting live comments to connected clients

---

## Features

* Realtime comment streaming using SSE
* Redis Pub/Sub integration
* PostgreSQL persistence layer
* Separation between write and stream services
* Layered architecture with infrastructure boundaries
* Docker Compose setup for local development
* Simple, idiomatic Go project structure

---

## Project Structure

```text
cmd/
    stream-api/     Entry point for SSE streaming service
    write-api/      Entry point for comment write service

internal/
    comment/        Comment domain models + service logic
    stream/         Stream subscription service
    infra/
        postgres/  PostgreSQL repository implementations
        redis/     Redis publisher/subscriber implementations
    transport/http/
        handlers/  HTTP handlers
        routers/   Route registration
```

---

## How It Works

### Writing Comments

1. Client sends a request to `write-api`
2. Comment is validated and persisted to PostgreSQL
3. Comment event is published to Redis
4. All active subscribers receive the event

### Streaming Comments

1. Client opens an SSE connection to `stream-api`
2. Service subscribes to a Redis channel for the target video
3. Incoming Redis messages are forwarded directly to connected clients

---

## What I Learned

This project was primarily built for learning distributed backend concepts and realtime communication patterns.

Topics explored include:

### Server-Sent Events (SSE)

Implemented long-lived HTTP streaming connections using:

* `http.Flusher`
* streaming response bodies
* request context cancellation

### Redis Pub/Sub

Explored transient message broadcasting between independent services.

### Layered Architecture

Separated:

* transport layer
* domain/service layer
* infrastructure layer

to better understand dependency boundaries and interface-driven design in Go.

### Infrastructure Abstractions

Used interfaces to decouple:

* repositories
* publishers
* subscribers

from concrete Redis/Postgres implementations.

### Service Separation

Split read/stream responsibilities from write responsibilities to mirror common distributed system patterns.

---

## Running the Project

### Requirements

* Go
* Docker
* Docker Compose

---

### Start infrastructure

```bash
docker compose up
```

---

### Run the services

#### write-api

```bash
go run ./cmd/write-api
```

#### stream-api

```bash
go run ./cmd/stream-api
```

---

## API Endpoints

### Create Comment

```http
POST /comment
```

Example body:

```json
{
  "videoID": 1,
  "content": "hello world"
}
```

---

### Stream Comments

```http
GET /stream/{videoID}
```

Returns a live SSE stream.

Example:

```text
data: {"ID":3,"VideoID":1,"Content":"hello world"}
```

---

## Known Limitations

* No authentication or authorization
* No reconnect handling for SSE clients
* No message replay/history on reconnect
* Redis Pub/Sub messages are transient
* No rate limiting
* Minimal validation/error handling
* No horizontal scaling strategy for SSE connections yet

---

## Future Improvements

* WebSocket support
* Redis Streams or Kafka for durable event streaming
* Authentication and user accounts
* Connection metrics and observability
* Graceful subscriber cleanup
* Replay support for missed comments
* Kubernetes deployment
* Load testing and benchmarking
* Structured logging and tracing

---
