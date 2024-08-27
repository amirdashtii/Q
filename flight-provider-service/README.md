# Flight Provider Service Documentation

## Overview
The Flight Provider Service is designed to manage and provide flight-related data, including the generation of random flights for the next 30 days, retrieving flight details by ID, and managing flight capacity (for booking and cancellations).

## How It Works
- **Flight Generation**: When the service starts, it automatically generates random flight data for the next 30 days if they are not already present. Each flight includes details like departure and arrival times, airline information, and available seats.
- **Retrieve Flight by ID**: The service allows you to query specific flights by their unique ID and get detailed flight information.
- **Manage Flight Capacity**: You can book or cancel seats on flights, which updates the remaining seat capacity accordingly.

## How to Use the Service
The service exposes several API endpoints that you can interact with:

### API Endpoints
1. **Retrieve Flight by ID**
   - **Endpoint**: `/flights/:id`
   - **Method**: `GET`
   - **Description**: Retrieves the details of a specific flight based on its ID.

2. **Reserve Seats on a Flight**
   - **Endpoint**: `/flights/:id/reserve`
   - **Method**: `PATCH`
   - **Request Body**:
     ```json
     {
       "seats": 2
     }
     ```
   - **Description**: Books a specified number of seats and decreases the available capacity.

3. **Cancel Seats on a Flight**
   - **Endpoint**: `/flights/:id/cancel`
   - **Method**: `PATCH`
   - **Request Body**:
     ```json
     {
       "seats": 2
     }
     ```
   - **Description**: Cancels a specified number of seats and increases the available capacity.

## Setup and Running the Service
To get the service up and running:

### 1. Configure the Environment
Rename the `.env-sample` file to `.env` and update it with your environment-specific values:

```plaintext
PORT=<your-host-port>

# Postgres Config
POSTGRES_HOST=<your-postgres-host>
POSTGRES_PORT=<your-postgres-port>
POSTGRES_DB_NAME=<your-postgres-db-name>
POSTGRES_USER=<your-postgres-user>
POSTGRES_PASSWORD=<your-postgres-password>
POSTGRES_AUTH_DATA_PATH=</path/to/your/local/folder/>
```

### 2. Run the Service
Use Docker Compose to build and start the service

```plaintext
docker-compose up --build
```

This command will start both the flight provider service and the PostgreSQL database. The service will be accessible at the configured port (e.g., http://localhost:8080).

## Technologies Used
- Go (Golang) for the core logic.
- PostgreSQL as the relational database.
- GORM as the ORM for database interactions.
- Docker & Docker Compose for containerization and orchestration.
 
