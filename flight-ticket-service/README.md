# Flight Ticket Service Documentation

## Overview
The Flight Ticket Service is designed to manage flight bookings, cancellations, and retrieving flight details. The service offers functionalities for booking tickets, managing passengers, and handling payments.

## How It Works
- **Flight Management**: Users can search for available flights and book tickets. The service tracks seat availability and manages reservations.
- **Ticketing**: Users can reserve and cancel tickets for flights. It provides real-time updates on seat availability and booking status.
- **Passenger Management**: The system allows users to manage passenger profiles.
- **Payment Handling**: The service integrates with payment gateways to process ticket payments and confirm transactions.

## How to Use the Service

The following API endpoints are available:

### API Endpoints

#### Flights
1. **Retrieve Available Flights**
   - **Endpoint**: `/flights`
   - **Method**: `GET`
   - **Query Parameters**:
     - `source`: Flight source location (e.g., `Tehran`).
     - `destination`: Flight destination (e.g., `Mashhad`).
     - `departure_date`: Date of flight (in `YYYY-MM-DD` format).
     - `sort_by`, `order`, `filter_by`: Optional parameters for sorting and filtering.
   
   **Example Request**:
   ```
   GET /flights?source=Tehran&destination=Mashhad&departure_date=2024-09-01
   ```

2. **Retrieve Flight by ID**
   - **Endpoint**: `/flights/:id`
   - **Method**: `GET`
   - **Description**: Fetches flight details based on the flight ID.

#### Tickets
1. **Reserve Ticket**
   - **Endpoint**: `/tickets`
   - **Method**: `POST`
   - **Request Body**:
     ```
     {
       "flight_id": "12345",
       "passenger_ids": ["uuid1", "uuid2"]
     }
     ```
   - **Description**: Reserves tickets for passengers on a specific flight.

2. **Get Ticket by ID**
   - **Endpoint**: `/tickets/:id`
   - **Method**: `GET`
   - **Description**: Retrieves ticket details for a specific ticket ID.

3. **Get All Tickets**
   - **Endpoint**: `/tickets`
   - **Method**: `GET`
   - **Description**: Retrieves all tickets associated with the logged-in user.

4. **Cancel Ticket**
   - **Endpoint**: `/tickets/cancel/:id`
   - **Method**: `POST`
   - **Description**: Cancels a ticket by its ID.

#### Payments
1. **Pay for Ticket**
   - **Endpoint**: `/payment/pay`
   - **Method**: `POST`
   - **Request Body**:
     ```
     {
       "ticket_id": "12345",
       "payment_gateway": "saman"
     }
     ```
   - **Description**: Initiates payment for a ticket using the specified payment gateway.

2. **Confirm Payment Success**
   - **Endpoint**: `/payment/success`
   - **Method**: `POST`
   - **Description**: Confirms successful payment after receiving the payment receipt.

#### Passengers
1. **Create Passenger**
   - **Endpoint**: `/user/passengers`
   - **Method**: `POST`
   - **Description**: Adds a new passenger for the logged-in user.

2. **Get All Passengers**
   - **Endpoint**: `/user/passengers`
   - **Method**: `GET`
   - **Description**: Retrieves a list of all passengers associated with the logged-in user.

3. **Get Passenger by ID**
   - **Endpoint**: `/user/passengers/:id`
   - **Method**: `GET`
   - **Description**: Retrieves passenger details based on the passenger's ID.

4. **Update Passenger**
   - **Endpoint**: `/user/passengers/:id`
   - **Method**: `PATCH`
   - **Description**: Updates passenger details for a specific passenger ID.

5. **Delete Passenger**
   - **Endpoint**: `/user/passengers/:id`
   - **Method**: `DELETE`
   - **Description**: Deletes a passenger by their ID.

## Setup and Running the Service

### 1. Configure the Environment
Copy the `.env-sample` file to `.env` and update the values as per your setup:

```plaintext
PORT=<your-host-port>

# Database Configuration
DB_HOST=<your-db-host>
DB_PORT=<your-db-port>
DB_NAME=<your-db-name>
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
``` 

### 2. Run the Service

#### Using Docker
To build and start the service using Docker:

```plaintext
docker-compose up --build
```

The service will be available at the configured port (e.g., http://localhost:8080).

#### Running Locally
If running locally:

1. Install Go dependencies:
   ```plaintext
   go mod download
   ```
2. Run the application:
   ```plaintext
   go run cmd/main.go
   ```
