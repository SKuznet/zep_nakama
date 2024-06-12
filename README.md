# Project Documentation

## Requirements
- Go 1.20
- Docker and Docker Compose

## Installation Steps
1. Install Go 1.20 from the official [Go website](https://golang.org/dl/).
2. Install Docker and Docker Compose by following the instructions on the [Docker website](https://docs.docker.com/get-docker/).

## Setup
1. Run the following commands to tidy and vendor Go modules:
    ```bash
    go mod tidy
    go mod vendor
    ```

## Running the Application
1. Start the application using Docker Compose:
    ```bash
    docker-compose up
    ```
2. Ensure that the health check indicates the application is running successfully.

## Testing Endpoints
1. Access the API Explorer at [http://127.0.0.1:7351/#/apiexplorer?endpoint=checksum](http://127.0.0.1:7351/#/apiexplorer?endpoint=checksum).
2. Use the following credentials to log in:
    - Username: admin
    - Password: password
3. Example Payload
```json
{
  "type": "core",
  "version": "1.0.0",
  "hash": null
}
