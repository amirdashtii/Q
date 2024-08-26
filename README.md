# Super App Documentation

## Overview

This super app is designed using a microservices architecture. Each service is independent, allowing for modular development, deployment, and scaling. The initial setup includes an authentication service, and other services like ticket booking, payment, and more will be added over time.

## Architecture

- **Microservices**: Independent services with dedicated responsibilities.
- **Communication**: Services communicate via HTTP APIs 
- **Data Management**: Each service manages its own database to ensure decoupling and independent scaling.
- **Caching and Session Management**: Redis is used for caching and session management.

## Services

1. **Authentication Service**: Manages user registration, login, and authorization.

## Prerequisites

- Docker
- Go 1.19 or higher
- PostgreSQL (or any preferred database)
- Redis
- [Optional] Kubernetes (for deployment and scaling)


## Setting Up the Project

Each service in this project is independent and has its own configuration and environment variables. Below is the general process for setting up any service.

### Clone the Repository:
   ```bash
   git clone https://github.com/amirdashtii/Q.git
   cd Q
   ```

### Navigate to the Specific Service Directory:
Each service has its own folder, which contains its configuration files including `.env` and `docker-compose.yml`.

For example, for the authentication service:
```bash
cd auth-service
```

### Environment Configuration:
Each service includes a `.env-sample` file that contains the required environment variables. You need to rename this file to `.env` and fill in the required variables.

```bash
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
