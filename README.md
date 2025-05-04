# Workdone AI

An intelligent CRM platform for micro and small businesses, powered by AI and built with Go and Fiber.

## Overview

Workdone AI is a smart CRM platform designed to help small businesses manage leads, track agent performance, and improve customer engagement through AI-powered automation. The platform supports multiple communication channels including email, SMS, and WhatsApp, with a focus on intelligent automation and predictive analytics.

## Features

- **AI-Powered Lead Management**
  - Intelligent multi-channel lead capture
  - Automated lead tagging and categorization
  - Smart status tracking
  - AI-driven follow-up management

- **Smart Agent Tools**
  - Location intelligence
  - Automated visit logging
  - AI communication assistant
  - Predictive performance analytics

- **AI Integration Support**
  - Smart Gmail processing
  - Intelligent SMS handling
  - WhatsApp Business API
  - GitHub Inference API

## Tech Stack

- **Backend**: Go 1.21+, Fiber
- **Database**: MySQL 8.0+
- **Authentication**: JWT, OAuth 2.0
- **AI Integration**: GitHub Inference API

## Project Structure

```
workdone-api/
├── app/
│   ├── configs/         # Configuration management
│   ├── controllers/     # API endpoints
│   ├── database/        # Database connections
│   ├── factory/         # Factory patterns
│   ├── handler/         # Request handlers
│   ├── middlewares/     # HTTP middlewares
│   ├── models/          # Database models
│   ├── routes/          # Route definitions
│   ├── service/         # Business logic
│   └── common/          # Shared utilities
├── migrations/          # Database migrations
└── memory-bank/         # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- MySQL 8.0 or later
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/workdone-api.git
   cd workdone-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Run migrations:
   ```bash
   go run main.go migrate
   ```

5. Start the server:
   ```bash
   go run main.go
   ```

## API Documentation

API documentation is available in the `memory-bank` directory. Key documents include:

- `projectbrief.md`: Project overview and objectives
- `productContext.md`: Product requirements and user workflows
- `systemPatterns.md`: Technical architecture and patterns
- `techContext.md`: Development environment and setup
- `activeContext.md`: Current focus and next steps
- `progress.md`: Implementation status and milestones

## Development

### Code Style

- Follow Go standard formatting
- Use meaningful variable and function names
- Document complex logic
- Write tests for new features

### Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 