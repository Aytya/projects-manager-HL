# Go - Project Management Service

### REST API for managing tasks that include three entities: User, Task, and Project. The service support CRUD operations, published on GitHub, embedded in Render, and run using Docker Compose and Makefile.

## Getting Started

### Prerequisites:
    - Go 1.22 or later
    - Docker
### Installation:
1. Clone the repository:
   ```bash
   https://github.com/Aytya/projects-manager-HL
   ```
2. Navigate into the project directory:
   ```bash
    cd projects-manager-HL
   ```
3. Install dependencies:
   ```bash
    go get -u "github.com/swaggo/gin-swagger"
    go get -u "github.com/swaggo/files"
    go get -u "github.com/gin-gonic/gin"
    go get -u "github.com/lib/pq"
    go get -u "github.com/joho/godotenv"
    go get -u "github.com/spf13/viper"
   ```

##  Build and Run Locally:
### Build the application:
   ```bash
   make build
   ```
### Run the application:
   ```bash
   make run
   ```
### Stop the application:
   ```bash
   make down
   ```
### Make tests:
   ```bash
   make test
   ```

### Generate Swagger Documentation:
   ```bash
   swag init --md ./docs --parseInternal --parseDependency --parseDepth 2 -g cmd/main.go
   ```

## API Endpoints:
### All API endpoints can be accessed through swagger, but here is data for post requests
#### Create a New User:
   - URL: http://localhost:8080/users
  - Method: POST
  - Request Body:
 ```bash
    {
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "USER"
   }
 ```

#### Create a New Project:
- URL: http://localhost:8080/users
- Method: POST
- Request Body:
 ```bash
    {
      "title": "New Project 2",
      "description": "This is a easy project.",
      "manager": "b48d5792-a0b6-4720-9f10-1100c549e7bd"
    } 
 ```

#### Create a New Task:
- URL: http://localhost:8080/users
- Method: POST
- Request Body:
 ```bash
    {
      "title": "Write Unit Tests 1",
      "description": "Create unit tests for the authentication module.",
      "priority": "High",
      "status": "Not Started",
      "assignee": "219edf66-e5e3-488e-822d-9318ba1e2598",
      "project": "248b14f4-0bec-43c7-a846-cf728a313961",
      "finished_at": null
   }
 ```

### Swagger Documentation
- URL: http://localhost:8080/swagger/index.html#/
