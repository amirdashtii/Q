# Q App Documentation

## Overview

This super app is designed using a microservices architecture. Each service is independent, allowing for modular development, deployment, and scaling. The initial setup includes an authentication service, a flight provider service, and a flight ticket booking service. Other services like payment handling are integrated as well, with more planned to be added over time.

## Architecture

- **Microservices**: Independent services with dedicated responsibilities.
- **Communication**: Services communicate via HTTP APIs.
- **Data Management**: Each service manages its database to ensure decoupling and independent scaling.
- **Caching and Session Management**: Redis is used for caching and session management.

## Services

1. **Authentication Service**: Manages user registration, login, and authorization.
2. **Flight Provider Service**: Generates random flight data for the next 30 days, manages flight capacity, and allows retrieving flight details by ID.
3. **Flight Ticket Service**: Manages flight bookings, cancellations, and ticketing, allowing users to search for available flights, book tickets, and manage passengers.
4. **Payment Service**: Integrates with payment gateways to handle ticket payments, offering functionality to pay for tickets and verify successful transactions.

## Prerequisites

- Docker
- Go 1.19 or higher
- PostgreSQL (or any preferred database)
- Redis
- [Optional] Kubernetes (for deployment and scaling)

## Setting Up the Project

Each service in this project is independent and has its own configuration and environment variables. Below is the general process for setting up any service.

### Clone the Repository:
```
git clone https://github.com/amirdashtii/Q.git
cd Q
```

### Navigate to the Specific Service Directory:
Each service has a folder, containing its configuration files including `.env` and `docker-compose.yml`.

For example, for the authentication service:
```
cd auth-service
```

For the flight provider service:
```
cd flight-provider-service
```

For the flight ticket service:
```
cd flight-ticket-service
```

### Environment Configuration:
Each service includes a `.env-sample` file that contains the required environment variables. You must rename this file to `.env` and fill in the required variables.

```
mv .env-sample .env
```

Update the variables according to your environment.

### Running the Service:
Each service can be run independently using Docker Compose:

```bash
docker-compose up --build
```

## Inter-Service Communication
- **API Calls**: Services expose RESTful APIs that other services can consume.

## Deployment
Services can be deployed independently using Docker or Kubernetes. Each service has its own Dockerfile and deployment configuration.
