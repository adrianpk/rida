# Rida

Rida is a personal challenge to simulate a simple scooter sharing system. It tracks trips, locations, and availability of scooters, and lets users find and follow them around the city. It also includes a basic simulation of client behavior to keep things moving. Not a real service, just a fun way to explore ideas.

## Overview

This repository contains the implementation, installation instructions, deployment and API usage details, as well as relevant technical information.


## Installation

- Clone the repository:
  ```sh
  git clone git@github.com:adrianpk/rida.git
  cd rida
  ```
- Run the application (this will build if needed):
  ```sh
  make run
  ```

## Docker Usage

You can build and run the application using Docker Compose:

```sh
make run-docker
```

### Sanity Check

The following command builds the Docker image, runs the container, and checks the `/healthz` endpoint to ensure the service starts correctly:

```sh
make test-docker
```

## API

- **GET /api/v1/scooters**: Search for scooters by area and status.
- **POST /api/v1/events**: Report scooter events (start, end, location updates).
- **GET /healthz**: Health check

Authentication is performed via the `X-API-Key` header.

## Project Structure

- `main.go`: Entry point.
- `internal/`: Business logic, simulated client, repo, API.
- `deployment/`: Dockerfile and docker-compose.
- `docs/`: Documentation and requirements.

## Testing & Quality

- Run tests:
  ```sh
  make test
  ```
- Lint and format:
  ```sh
  make check
  ```
