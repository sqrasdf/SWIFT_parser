# SWIFT Code API

## Description

This is a simple REST API for managing SWIFT codes, built with Go, Gin, and PostgreSQL.

## Prerequisites

- Docker
- Docker Compose

## Setup

1.  Clone the repository:

    ```
    git clone https://github.com/sqrasdf/SWIFT_parser/

    ```

2.  Enter project folder

    ```
    cd SWIFT_parser
    ```

## Running the Application

1.  Build and run the application using Docker Compose:

    ```
    docker-compose up --build
    ```

    This command will:

    - Build the Docker image for the application.
    - Start the PostgreSQL database in a separate container.
    - Run the application, connecting to the database.

2.  The API will be accessible at `http://localhost:8080`.

## API Endpoints

### 1. Get SWIFT Code Details

- **Endpoint:** `GET /v1/swift-codes/{swift-code}`
- **Description:** Retrieves details for a single SWIFT code (headquarters or branch).
- **Example:**

  ```
  curl http://localhost:8080/v1/swift-codes/AAISALTRXXX
  ```

### 2. Get SWIFT Codes by Country

- **Endpoint:** `GET /v1/swift-codes/country/{countryISO2code}`
- **Description:** Retrieves all SWIFT codes (headquarters and branches) for a specific country.
- **Example:**

  ```
  curl http://localhost:8080/v1/swift-codes/country/US
  ```

### 3. Add a New SWIFT Code

- **Endpoint:** `POST /v1/swift-codes`
- **Description:** Adds a new SWIFT code entry to the database.
- **Request Body Example:**

  ```
  {
    "swiftCode": "AAABBBCCXXX",
    "bankName": "New Bank",
    "countryISO2": "PL",
    "countryName": "POLAND",
    "address": "ul. Testowa",
    "isHeadquarter": true
  }
  ```

- **Example:**

  ```
  curl -X POST -H "Content-Type: application/json" -d '{
    "swiftCode": "AAABBBCCXXX",
    "bankName": "New Bank",
    "countryISO2": "PL",
    "countryName": "POLAND",
    "address": "ul. Testowa",
    "isHeadquarter": true
  }' http://localhost:8080/v1/swift-codes
  ```

### 4. Delete a SWIFT Code

- **Endpoint:** `DELETE /v1/swift-codes/{swift-code}`
- **Description:** Deletes a SWIFT code from the database.
- **Example:**

  ```
  curl -X DELETE http://localhost:8080/v1/swift-codes/CITIBGSFTRD
  ```

### Testing

To run tests, you can use the following command:

```
docker-compose -f docker-compose.test.yml up --build
```

This command will run all tests in the project. Running it for the first time may take a little while.

## Stopping and Cleaning Up

To stop the application, you can either press `CTRL+C` in the terminal or run the following command:

```
docker compose down
```
