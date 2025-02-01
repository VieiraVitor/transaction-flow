# 🚀 Transaction Flow API

## 📌 About the Project

The **Transaction Flow API** is an account and transaction management system that allows:
- Creating accounts associated with a **document number**.
- Retrieving information about an existing account.
- Registering financial transactions of different types (**compras à vista, compras parceladas, saques e pagamentos**).

## 🛠️ Technologies Used

This project uses the following technologies:

- **[Go](https://go.dev/doc/)** - Programming language.
- **[Chi Router](https://github.com/go-chi/chi)** - Lightweight and flexible routing for HTTP APIs.
- **[PostgreSQL](https://www.postgresql.org/docs/)** - Relational database.
- **[Docker](https://docs.docker.com/)** - Project containerization for easy setup.
- **[Slog](https://pkg.go.dev/log/slog)** - Library for structured logging.

---

## 📂 Project Structure

```plaintext
├── cmd/
│   ├── transaction-flow/        # Service entry point
│   │   ├── main.go              # HTTP server initialization
│
├── internal/
│   ├── api/                     # HTTP interface layer
│   │   ├── handlers/             # Handlers for HTTP requests
│   │   ├── middleware/           # Middlewares (logging, recovery)
│   │   ├── dto/                  # Request/response structures
│   │   ├── response/              # Error formatting
│   │
│   ├── application/             # Business logic (use cases)
│   │   ├── usecase/               # Use cases (Account, Transaction)
│   │
│   ├── domain/                  # Entity models and domain rules
│   │   ├── account.go
│   │   ├── transaction.go
│   │
│   ├── infra/                   # Infrastructure layer
│   │   ├── repository/            # Database access
│   │   ├── logger/                # Logging configuration
│   │
├── migrations/                  # SQL files for database creation and updates
│   ├── 001_init.up.sql
│
├── config/                      # Environment configuration
│   ├── config.go
│
├── docker-compose.yml            # Docker configuration for development
├── Dockerfile                    # Application container definition
├── Makefile                      # Useful automation commands
├── go.mod                        # Project dependencies
├── README.md                     # Project documentation
```

---

## 🚀 **How to Run the Project Locally**

### **📌 Prerequisites**

Before you begin, you will need to have installed:

- **[Docker](https://docs.docker.com/get-docker/)**
- **[Docker Compose](https://docs.docker.com/compose/install/)**

### **📌 Step by step**

1️⃣ **Clone the repository**
```bash
git clone https://github.com/your-user/transaction-flow.git
cd transaction-flow
```

2️⃣ **Start the project containers (API, database, and migrations)**
```bash
docker-compose up -d
```
📌 This will start the required services, including the database and application.


---

## 🔥 **API Endpoints**

### **📌 Create an Account**
📍 **POST** `/accounts`
```bash
curl -X POST http://localhost:8080/accounts \
     -H "Content-Type: application/json" \
     -d '{"document_number": "12345678900"}'
```
📌 **Response (201 Created)**
```json
{
  "id": 1
}
```

### **📌 Retrieve an Account**
📍 **GET** `/accounts/{id}`
```bash
curl -X GET http://localhost:8080/accounts/1
```
📌 **Response (200 OK)**
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

### **📌 Create a Transaction**
📍 **POST** `/transactions`
```bash
curl --X POST http://localhost:8080/transactions \
     -H "Content-Type: application/json" \
     -d '{
            "account_id": 1,
            "operation_type_id": 4,
            "amount": 123.45
        }'
```
📌 **Response (201 Created)**
```json
{
  "id": 10
}
```

---

## 📜 **License**
This project is distributed under the **MIT** license. See the [LICENSE](LICENSE) file for details.
