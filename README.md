# URL Shortener

A fully containerized URL shortening service built using Golang, PostgreSQL, Redis, and the Echo web framework. This service allows users to shorten URLs, store them in PostgreSQL, and cache frequent lookups in Redis. The entire system is orchestrated using Docker Compose, making deployment simple and consistent.

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [Prometheus and Grafana Metrics](#prometheus-and-grafana-metrics)
- [License](#license)

## Features
- Shortens URLs and stores them in PostgreSQL.
- Frequently accessed URLs are cached in Redis.
- Exposes a list of shortened URLs via API.
- Provides Prometheus metrics for monitoring HTTP requests and durations.
- Grafana for visualizing metrics and creating dashboards.
- Fully containerized and runs with a single `docker-compose` command.

## Project Structure
```
├── Dockerfile                     # Dockerfile to build the URL shortener service
├── cmd
│   └── server
│       └── main.go                # Main entry point for the application
├── data                           # Directory for storing Prometheus data
│   └── ...                        
├── docker-compose.yml             # Docker Compose file to orchestrate services
├── go.mod                         # Go module definition
├── go.sum                         # Go dependencies lock file
├── internal
│   ├── DB
│   │   └── db.go                  # Database (PostgreSQL, Redis) initialization
│   ├── handlers
│   │   └── handlers.go            # API handlers for shortening and redirecting URLs
│   └── util
│       └── base62.go              # Utility for Base62 encoding
├── prometheus.yml                 # Prometheus configuration
└── testClient                     # Test client for generating and sending requests
    └── ... 
```

## Installation

### Prerequisites
- Docker
- Docker Compose

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/url-shortener.git
   cd url-shortener
   ```

2. Create a `.env` file in the project root (if not already present) and configure the required environment variables:
   ```bash
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=Amineyk85
   DB_NAME=urlshortener
   REDIS_HOST=redis
   REDIS_PORT=6379
   REDIS_PASSWORD=amin
   ```

3. Ensure that Docker and Docker Compose are installed on your system.

## Running the Application

1. Simply run the following command:
   ```bash
   docker-compose up -d
   ```

2. The application will automatically:
   - Spin up the Golang service.
   - Set up PostgreSQL and Redis containers.
   - Expose the service on `http://localhost:1323`.

3. Access Prometheus metrics at `http://localhost:9090` and Grafana at `http://localhost:3000` (default credentials are admin/admin).

## API Endpoints

### 1. **Shorten a URL**
   - **Endpoint**: `/shorten`
   - **Method**: `POST`
   - **Description**: Shortens a given URL.
   - **Request**: 
     - Query Parameter: `url`
   - **Example**:
     ```bash
     curl -X POST "http://localhost:1323/shorten?url=https://example.com"
     ```

### 2. **Redirect to Original URL**
   - **Endpoint**: `/:shortenedURL`
   - **Method**: `GET`
   - **Description**: Redirects to the original URL using the shortened version.
   - **Example**:
     ```bash
     curl -X GET "http://localhost:1323/{shortenedURL}"
     ```

### 3. **List All URLs**
   - **Endpoint**: `/urls`
   - **Method**: `GET`
   - **Description**: Lists all shortened URLs and their original mappings.
   - **Example**:
     ```bash
     curl -X GET "http://localhost:1323/urls"
     ```

### 4. **Prometheus Metrics**
   - **Endpoint**: `/metrics`
   - **Method**: `GET`
   - **Description**: Exposes Prometheus metrics for monitoring HTTP requests and durations.
   - **Example**:
     ```bash
     curl -X GET "http://localhost:1323/metrics"
     ```

## Prometheus and Grafana Metrics

Prometheus metrics are exposed via the `/metrics` endpoint. The following metrics are collected:
- `http_requests_total`: Total number of HTTP requests, labeled by path and method.
- `http_request_duration_seconds`: Duration of HTTP requests, labeled by path and method.

Grafana is included in the setup for visualizing metrics and creating dashboards. You can log in with the default credentials (admin/admin) to create and manage your dashboards.

To scrape the metrics, you can configure Prometheus using the provided `prometheus.yml` configuration file.


