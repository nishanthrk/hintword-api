# Technical Context: Hintword API

## Technology Stack

### Core Technologies

#### Backend Framework
- **Language**: Go 1.23.0 (locked for compatibility)
- **Framework**: Fiber v2.52.6 (high-performance HTTP framework)
- **ORM**: GORM v1.25.12 (Go ORM library)
- **Database Driver**: MySQL driver v1.5.7

#### Database & Storage
- **Primary Database**: MySQL 8.0+
- **Caching**: Redis v8.11.5
- **Cloud Storage**: AWS S3
- **Migration Tool**: Goose v3.24.2

#### Authentication & Security
- **JWT**: golang-jwt/jwt v4.5.2
- **OAuth**: Google OAuth 2.0 integration
- **Password Hashing**: bcrypt (via Go crypto package)
- **Validation**: go-playground/validator v10.26.0

#### AI & External Services
- **AI Integration**: OpenAI API
- **Email Service**: SendGrid v3.16.0
- **SMS Service**: MSG91 integration
- **Cloud Services**: AWS SDK v2

## Development Environment

### Prerequisites
- **Go Version**: 1.23.0+ (toolchain go1.23.8)
- **MySQL**: 8.0+ with InnoDB engine
- **Redis**: 6.0+ for caching and sessions
- **Docker**: For containerization and deployment
- **AWS CLI**: For cloud service integration

### Environment Setup

#### Local Development
```bash
# Clone repository
git clone <repository-url>
cd hintword-api

# Install dependencies
go mod download

# Setup environment variables
cp .env.example .env
# Edit .env with your configuration

# Run migrations
go run main.go migrate

# Start development server
go run main.go
```

#### Environment Variables
```env
# Application
ENV=development
APP_HOST=0.0.0.0
APP_PORT=9090

# Database
DB_HOST=localhost
DB_PORT=3306
DB_DRIVER=mysql
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=hintword

# JWT
JWT_ACCESS_SIGN_KEY=your_access_secret
JWT_REFRESH_SIGN_KEY=your_refresh_secret
JWT_ISSUER=hintword.com

# OAuth
GOOGLE_OAUTH_CLIENT_ID=your_client_id
GOOGLE_OAUTH_CLIENT_SECRET=your_client_secret
GOOGLE_OAUTH_REDIRECTION_URL=your_redirect_url

# AI Services
OPENAI_API_KEY=your_openai_key
OPENAI_API_URL=https://api.openai.com/v1

# AWS (Production)
AWS_REGION=ap-south-1
ENV_LOAD_METHOD=SSM
ENV_LOAD_PATH=/hintword/production
```

## Project Structure

### Directory Organization
```
hintword-api/
├── app/                          # Application code
│   ├── common/                   # Shared utilities
│   │   ├── constants/           # Application constants
│   │   ├── passwordutil/        # Password utilities
│   │   ├── utility/             # General utilities
│   │   └── validator/            # Custom validators
│   ├── configs/                 # Configuration management
│   │   ├── config.go           # Main configuration
│   │   ├── mysql.go            # Database config
│   │   └── tenant.go           # Multi-tenant config
│   ├── controllers/            # HTTP controllers
│   │   └── v1/                 # API version 1
│   │       ├── agent/          # AI agent endpoints
│   │       ├── auth/           # Authentication
│   │       ├── note/           # Note management
│   │       └── tab/            # Tab/collection management
│   ├── database/               # Database layer
│   │   ├── database.go        # Database connection
│   │   ├── mysql.go           # MySQL specific
│   │   └── redis.go           # Redis connection
│   ├── factory/               # Service factories
│   │   ├── cloud_storage/     # Storage abstraction
│   │   ├── email_provider/    # Email services
│   │   └── sms_provider/      # SMS services
│   ├── handler/               # Request handlers
│   ├── middlewares/           # HTTP middlewares
│   ├── models/                # Data models
│   ├── routes/                # Route definitions
│   └── services/              # Business logic
├── migrations/                 # Database migrations
├── memory-bank/               # Project documentation
├── scripts/                   # Deployment scripts
└── main.go                    # Application entry point
```

## Database Schema

### Core Tables

#### Users Table
```sql
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    google_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### Notes Table
```sql
CREATE TABLE notes (
    note_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    title VARCHAR(500),
    content JSON,
    folder VARCHAR(255),
    tags JSON,
    sequence INT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
```

#### Tabs/Collections Table
```sql
CREATE TABLE tabs (
    tab_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
```

## API Design

### RESTful API Structure
```
Base URL: /v1/

Authentication:
POST /v1/auth/google/user-info

Notes:
GET    /v1/note/list
POST   /v1/note/create

Tabs/Collections:
GET    /v1/tab/collection
POST   /v1/tab/collection
POST   /v1/tab/create

AI Agent:
POST   /v1/agent/completion
```

### Request/Response Format

#### Standard Response Format
```json
{
    "success": true,
    "data": { ... },
    "message": "Operation successful",
    "timestamp": "2024-01-01T00:00:00Z"
}
```

#### Error Response Format
```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid input data",
        "details": { ... }
    },
    "timestamp": "2024-01-01T00:00:00Z"
}
```

## Configuration Management

### Environment-Based Configuration

#### Local Development
- Uses `.env` file for configuration
- Environment variables loaded via `godotenv`
- Development-specific settings

#### Production Deployment
- AWS SSM Parameter Store integration
- Secure configuration management
- Environment-specific parameter paths

### Configuration Loading Strategy
```go
// Configuration loading based on environment
switch envMethod {
case constants.EnvLoadMethodLocal:
    cfg.LoadLocalConfig()
case constants.EnvLoadMethodSSM:
    cfg.LoadSSMConfig()
}
```

## Dependencies

### Core Dependencies
```go
require (
    github.com/gofiber/fiber/v2 v2.52.6          // HTTP framework
    github.com/gofiber/jwt/v2 v2.2.7              // JWT middleware
    github.com/gofiber/websocket/v2 v2.2.1        // WebSocket support
    github.com/golang-jwt/jwt/v4 v4.5.2          // JWT implementation
    github.com/go-playground/validator/v10 v10.26.0 // Validation
    github.com/go-redis/redis/v8 v8.11.5         // Redis client
    github.com/joho/godotenv v1.5.1              // Environment loading
    github.com/sirupsen/logrus v1.9.3            // Logging
    gorm.io/driver/mysql v1.5.7                 // MySQL driver
    gorm.io/gorm v1.25.12                       // ORM
)
```

### AWS Dependencies
```go
require (
    github.com/aws/aws-sdk-go-v2 v1.36.3         // AWS SDK v2
    github.com/aws/aws-sdk-go-v2/config v1.29.14 // AWS configuration
    github.com/aws/aws-sdk-go-v2/service/s3 v1.79.2 // S3 service
    github.com/aws/aws-sdk-go-v2/service/ssm v1.58.2 // SSM service
)
```

### External Service Dependencies
```go
require (
    github.com/sendgrid/sendgrid-go v3.16.0+incompatible // Email service
    github.com/google/uuid v1.6.0                // UUID generation
    github.com/guregu/null v4.0.0+incompatible   // Null handling
    github.com/pressly/goose/v3 v3.24.2          // Database migrations
)
```

## Development Tools

### Code Quality
- **Linting**: Built-in Go linting
- **Formatting**: `go fmt` for code formatting
- **Testing**: `go test` for unit and integration tests
- **Documentation**: Go doc for API documentation

### Database Management
- **Migrations**: Goose for database schema management
- **Connection Pooling**: GORM with connection pooling
- **Query Optimization**: GORM query optimization

### Monitoring & Logging
- **Logging**: Logrus for structured logging
- **Error Tracking**: Custom error handling and logging
- **Performance**: Fiber built-in performance monitoring

## Deployment Architecture

### Containerization
- **Docker**: Multi-stage builds for optimization
- **Base Image**: Alpine Linux for minimal size
- **Port**: 9090 (configurable via environment)

### Cloud Infrastructure
- **Container Registry**: AWS ECR
- **Compute**: AWS EC2 instances
- **Database**: MySQL on EC2 or RDS
- **Storage**: AWS S3 for file storage
- **Configuration**: AWS SSM Parameter Store

### CI/CD Pipeline
- **Source Control**: GitHub
- **CI/CD**: GitHub Actions
- **Build**: Docker container builds
- **Deploy**: Automated deployment to EC2
- **Branch Strategy**: `release/hintword` branch for production

## Performance Considerations

### Database Optimization
- **Connection Pooling**: Efficient database connections
- **Query Optimization**: GORM query optimization
- **Indexing**: Proper database indexing strategy
- **Caching**: Redis for frequently accessed data

### Application Performance
- **Fiber Framework**: High-performance HTTP framework
- **Goroutines**: Concurrent request processing
- **Memory Management**: Efficient memory usage
- **Response Caching**: Strategic caching implementation

### Scalability
- **Horizontal Scaling**: Stateless application design
- **Load Balancing**: Multiple instance support
- **Database Scaling**: Read replicas and connection pooling
- **Caching Strategy**: Distributed caching with Redis

## Security Considerations

### Authentication & Authorization
- **JWT Tokens**: Secure token-based authentication
- **OAuth Integration**: Google OAuth for user authentication
- **Role-Based Access**: User permission management
- **Token Refresh**: Secure token refresh mechanism

### Data Security
- **Input Validation**: Comprehensive input validation
- **SQL Injection Prevention**: GORM ORM protection
- **XSS Protection**: Input sanitization
- **HTTPS**: Secure communication protocols

### Infrastructure Security
- **Environment Variables**: Secure configuration management
- **AWS IAM**: Proper AWS permissions
- **Network Security**: VPC and security groups
- **Data Encryption**: Encryption at rest and in transit

---

*This document provides the technical foundation for development, deployment, and maintenance of the Hintword API.*
