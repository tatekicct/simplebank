# Simple Bank

A secure bank management system built with Go, featuring RESTful APIs for user authentication, account management, and money transfers with PostgreSQL as the database and JWT/PASETO token-based authentication.

## 🏗️ Architecture Overview

```mermaid
graph TB
    Client[Client Applications] --> API[REST API Server<br/>Gin Framework]
    API --> Auth[Authentication Layer<br/>JWT/PASETO Tokens]
    Auth --> Middleware[Authorization Middleware<br/>Role-based Access Control]
    Middleware --> Store[Store Layer<br/>Transaction Management]
    Store --> DB[(PostgreSQL<br/>Database)]
    
    Store --> SQLC[SQLC Generated Code<br/>Type-safe SQL]
    SQLC --> Queries[SQL Queries]
    
    API --> Mock[Mock Store<br/>Testing]
    
    subgraph "Authentication & Authorization"
        Auth
        Middleware
        JWT[JWT Token Maker]
        PASETO[PASETO Token Maker]
        Auth --> JWT
        Auth --> PASETO
    end
    
    subgraph "Database Layer"
        DB
        Migration[Database Migrations]
        Users[(Users Table)]
        Accounts[(Accounts Table)]
        DB --> Users
        DB --> Accounts
    end
    
    subgraph "Business Logic"
        Store
        SQLC
        Queries
    end
    
    subgraph "API Layer"
        API
        Handlers[HTTP Handlers]
        Router[Router Configuration]
    end
    
    subgraph "Testing"
        Mock
        Tests[Unit Tests]
    end
```

## 📁 Project Structure

```
simplebank/
├── 📁 api/                    # HTTP API layer
│   ├── account.go             # Account handlers
│   ├── account_test.go        # Account handler tests
│   ├── main_test.go           # API test setup
│   ├── middleware.go          # Authentication middleware
│   ├── middleware_test.go     # Middleware tests
│   ├── server.go              # HTTP server setup
│   ├── transfer.go            # Transfer handlers
│   ├── transfer_test.go       # Transfer handler tests
│   ├── user.go                # User authentication handlers
│   ├── user_test.go           # User handler tests
│   └── validator.go           # Custom validators
│
├── 📁 db/                     # Database layer
│   ├── 📁 migration/          # Database migration files
│   │   ├── 000001_init_schema.up.sql    # Create tables
│   │   ├── 000001_init_schema.down.sql  # Drop tables
│   │   ├── 000002_add_users.up.sql      # Add users table
│   │   └── 000002_add_users.down.sql    # Remove users table
│   │
│   ├── 📁 mock/               # Mock implementations for testing
│   │   └── store.go           # Generated mock store
│   │
│   ├── 📁 query/              # SQL query files
│   │   ├── account.sql        # Account operations
│   │   ├── entry.sql          # Entry operations
│   │   ├── transfer.sql       # Transfer operations
│   │   └── user.sql           # User operations
│   │
│   └── 📁 sqlc/               # Generated SQLC code
│       ├── account.sql.go     # Account query implementations
│       ├── account_test.go    # Account query tests
│       ├── db.go              # Database connection
│       ├── entry.sql.go       # Entry query implementations
│       ├── error.go           # Custom errors
│       ├── main_test.go       # Database test setup
│       ├── models.go          # Database models
│       ├── querier.go         # Query interface
│       ├── store.go           # Transaction store
│       ├── store_test.go      # Store tests
│       ├── transfer.sql.go    # Transfer query implementations
│       ├── user.sql.go        # User query implementations
│       └── user_test.go       # User query tests
│
├── 📁 token/                  # Authentication tokens
│   ├── jwt_maker.go           # JWT token implementation
│   ├── jwt_maker_test.go      # JWT maker tests
│   ├── maker.go               # Token maker interface
│   ├── paseto_maker.go        # PASETO token implementation
│   ├── paseto_maker_test.go   # PASETO maker tests
│   └── payload.go             # Token payload structure
│
├── 📁 util/                   # Utility functions
│   ├── config.go              # Configuration management
│   ├── currency.go            # Currency validation
│   ├── password.go            # Password hashing utilities
│   ├── password_test.go       # Password tests
│   ├── random.go              # Random data generation
│   └── role.go                # User roles definition
│
├── 📁 .github/workflows/      # CI/CD pipelines
│   └── ci.yml                 # GitHub Actions workflow
│
├── 📄 app.env                 # Environment configuration
├── 📄 go.mod                  # Go module definition
├── 📄 go.sum                  # Go dependencies lock
├── 📄 main.go                 # Application entry point
├── 📄 Makefile                # Build automation
└── 📄 sqlc.yaml               # SQLC configuration
```

## 🏦 Database Schema

```mermaid
erDiagram
    USERS ||--o{ ACCOUNTS : "owns"
    ACCOUNTS ||--o{ ENTRIES : "has entries"
    ACCOUNTS ||--o{ TRANSFERS : "from account"
    ACCOUNTS ||--o{ TRANSFERS : "to account"
    TRANSFERS ||--o{ ENTRIES : "generates entries"

    USERS {
        varchar username PK
        varchar hashed_password
        varchar full_name
        varchar email "UNIQUE"
        timestamptz password_changed_at
        timestamptz created_at
    }

    ACCOUNTS {
        bigserial id PK
        varchar owner FK
        bigint balance
        varchar currency
        timestamptz created_at
        unique_constraint owner_currency_key "owner, currency"
    }

    ENTRIES {
        bigserial id PK
        bigint account_id FK
        bigint amount
        timestamptz created_at
    }

    TRANSFERS {
        bigserial id PK
        bigint from_account_id FK
        bigint to_account_id FK
        bigint amount
        timestamptz created_at
    }
```

## 🔄 Transaction Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Auth
    participant Store
    participant DB

    Client->>API: POST /users/login
    API->>DB: Verify credentials
    DB->>API: User data
    API->>Auth: Generate access token
    Auth->>API: JWT/PASETO token
    API->>Client: Access token + User info

    Note over Client: Store token for future requests

    Client->>API: POST /transfers (with Bearer token)
    API->>Auth: Validate token
    Auth->>API: Token payload
    API->>Store: TransferTx(params)
    
    Store->>DB: BEGIN TRANSACTION
    Store->>DB: CREATE transfer record
    Store->>DB: CREATE from_entry (-amount)
    Store->>DB: CREATE to_entry (+amount)
    Store->>DB: UPDATE from_account balance
    Store->>DB: UPDATE to_account balance
    
    alt Success
        Store->>DB: COMMIT
        Store->>API: TransferResult
        API->>Client: 200 OK
    else Error
        Store->>DB: ROLLBACK
        Store->>API: Error
        API->>Client: Error Response
    end
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 17
- Docker (optional)
- Make

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/tatekicct/simplebank.git
   cd simplebank
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup PostgreSQL database**
   ```bash
   # Using Docker
   make postgres
   make createdb
   
   # Run migrations
   make migrateup
   ```

4. **Generate SQLC code**
   ```bash
   make sqlc
   ```

5. **Generate mock files**
   ```bash
   make mock
   ```

6. **Run tests**
   ```bash
   make test
   ```

7. **Start the server**
   ```bash
   make server
   ```

## 🛠️ Technology Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.24 |
| **Web Framework** | Gin |
| **Database** | PostgreSQL 17 |
| **Database Driver** | pgx/v5 |
| **Code Generation** | SQLC |
| **Authentication** | JWT & PASETO |
| **Password Hashing** | bcrypt |
| **Testing** | testify, GoMock |
| **Configuration** | Viper |
| **Validation** | go-playground/validator |
| **CI/CD** | GitHub Actions |
| **Containerization** | Docker |

## 📋 Available Make Commands

```bash
# Database operations
make postgres          # Start PostgreSQL container
make createdb          # Create database
make deletedb          # Delete database
make migrateup         # Run migrations up
make migratedown       # Run migrations down
make migrateup1        # Run single migration up
make migratedown1      # Run single migration down
make new_migration     # Create new migration (usage: make new_migration name=add_table)

# Code generation
make sqlc             # Generate SQLC code
make mock             # Generate mock files

# Development
make server           # Start the server
make test             # Run tests with coverage
make build            # Build the application
```

## 🔐 API Endpoints

### Public Endpoints (No Authentication Required)
- `POST /users` - Register a new user
- `POST /users/login` - User login (returns access token)

### Protected Endpoints (Require Authentication)
#### Accounts
- `POST /accounts` - Create a new account
- `GET /accounts/:id` - Get account by ID (owner verification)
- `GET /accounts` - List user's accounts with pagination

#### Transfers
- `POST /transfers` - Transfer money between accounts (with ownership verification)

## ⚡ Key Features

### � **Secure Authentication**
- JWT and PASETO token-based authentication
- Secure password hashing with bcrypt
- Role-based access control (Depositor/Banker roles)
- Token-based session management with configurable expiration

### 🛡️ **Authorization & Security**
- Bearer token authentication middleware
- Account ownership verification
- Protected API endpoints
- Secure token validation and verification

### �🔒 **ACID Transactions**
- All money transfers are executed within database transactions
- Automatic rollback on any failure
- Deadlock prevention through consistent ordering

### 🧪 **Comprehensive Testing**
- Unit tests for all database operations
- API endpoint testing with authentication
- Token generation and validation testing
- Concurrent transaction testing
- Mock-based testing for isolation

### 🏗️ **Clean Architecture**
- Separation of concerns (API, Authentication, Business Logic, Database)
- Interface-based design for testability
- Type-safe database operations with SQLC
- Middleware pattern for cross-cutting concerns

### 🚀 **Production Ready**
- Environment-based configuration
- Dual token support (JWT/PASETO)
- Structured logging
- Health checks
- CI/CD pipeline with GitHub Actions

## 🧪 Testing Strategy

The project includes comprehensive tests:

- **Unit Tests**: Database operations and business logic
- **Integration Tests**: API endpoints with authentication
- **Authentication Tests**: Token generation, validation, and middleware
- **Concurrent Tests**: Race condition and deadlock testing
- **Mock Tests**: Isolated component testing

Run specific test suites:
```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test ./api/...
go test ./db/sqlc/...
```

## 🔧 Configuration

The application uses environment-based configuration via `app.env`:

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
ACCESS_TOKEN_DURATION=15m
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run with Docker Compose
docker-compose up --build
```

### Manual Deployment
```bash
# Build for production
go build -o main main.go

# Run the application
./main
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is created for educational purposes based on the original [Simple Bank repository](https://github.com/techschool/simplebank). Please refer to the original repository for license information.

## 📚 About This Repository

This repository is a study project created for learning purposes, following the tutorial and concepts from the original [techschool/simplebank](https://github.com/techschool/simplebank) repository. It demonstrates modern Go development practices including:

- Backend web development with Go
- User authentication and authorization
- JWT and PASETO token implementation
- Database design and management
- RESTful API development
- Security best practices
- Testing strategies
- CI/CD implementation

## 🙏 Acknowledgments

- [TechSchool](https://github.com/techschool) for the original Simple Bank tutorial
- [SQLC](https://sqlc.dev/) for type-safe SQL code generation
- [Gin](https://gin-gonic.com/) for the web framework
- [testify](https://github.com/stretchr/testify) for testing utilities
- [GoMock](https://github.com/uber-go/mock) for mock generation
- [JWT-Go](https://github.com/golang-jwt/jwt) for JWT implementation
- [PASETO](https://github.com/o1egl/paseto) for PASETO implementation
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing
