# Room Booker

A comprehensive room booking system for office meeting rooms, built with Go, supporting PostgreSQL and SQLite, with Microsoft Graph integration for Outlook calendars.

## Features

- **Multi-Database Support**: PostgreSQL (production) and SQLite (development/testing)
- **Authentication**: OIDC with Microsoft accounts, optional password-based auth with Argon2 hashing
- **Authorization**: RBAC with User, Admin, Manager roles
- **Calendar Integration**: Microsoft Graph API for Outlook room calendars, fallback to local storage
- **Booking Management**: Create, update, cancel bookings with conflict detection
- **Availability Search**: Find available rooms by capacity, equipment, floor, time
- **Recurring Bookings**: Support for RRULE-based recurring events
- **Audit Logging**: Track all actions for compliance
- **Web UI**: Simple interface with FullCalendar.js
- **API**: OpenAPI 3.1 compliant REST API
- **Docker**: Containerized deployment
- **CI/CD**: Jenkins pipeline for automated testing and deployment

## Architecture

### Overview

The system follows a clean architecture with layers:

- **Transport**: HTTP handlers (Chi router)
- **Service/Usecase**: Business logic
- **Repository**: Data access layer
- **Provider**: External integrations (Graph, SMTP)

### Database Schema

Key tables:

- `users`: User accounts and roles
- `offices`: Office locations
- `floors`: Building floors
- `rooms`: Meeting rooms with equipment
- `bookings`: Room reservations
- `booking_participants`: Meeting attendees
- `booking_rules`: Office policies
- `holidays`: Non-working days
- `audit_logs`: Action history

### Conflict Prevention

- PostgreSQL: Uses `tsrange` with `EXCLUDE` constraint for atomic conflict detection
- SQLite: Application-level locking with pessimistic concurrency

## Quick Start

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- (Optional) PostgreSQL for production

### Local Development

1. Clone the repository:

   ```bash
   git clone <repo-url>
   cd roombooker
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment:

   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

4. Run migrations:

   ```bash
   make migrate-up
   make seed
   ```

5. Start the server:

   ```bash
   make dev
   ```

6. Open http://localhost:8080

### Docker

```bash
make docker-up
```

### Testing

```bash
make test
make test-integration
```

## Configuration

### Environment Variables

See `.env.example` for all options.

### Microsoft Graph Setup

1. Register an app in Azure AD:

   - Go to Azure Portal > App Registrations
   - Create new app, note Client ID and Tenant ID
   - Add redirect URI: `http://localhost:8080/auth/oidc/callback`
   - Generate client secret

2. Grant permissions:

   - Calendars.ReadWrite
   - User.Read
   - offline_access

3. For room calendars:

   - Create resource mailboxes in Exchange Admin Center
   - Note the email addresses for `graph_resource_id`

4. Update `.env`:
   ```
   GRAPH_CLIENT_ID=your-client-id
   GRAPH_CLIENT_SECRET=your-client-secret
   GRAPH_TENANT_ID=your-tenant-id
   ```

### Database Switching

- **SQLite** (default): `DATABASE_DRIVER=sqlite3`
- **PostgreSQL**: `DATABASE_DRIVER=postgres`

For PostgreSQL, install golang-migrate and run:

```bash
migrate -path migrations -database "postgres://user:pass@localhost/dbname?sslmode=disable" up
```

## API Documentation

The API follows OpenAPI 3.1 spec in `openapi.yaml`.

### Key Endpoints

- `POST /auth/login` - Password login
- `GET /auth/oidc/start` - OIDC login
- `GET /rooms` - List rooms
- `GET /availability` - Search availability
- `POST /rooms/{id}/bookings` - Create booking
- `GET /rooms/{id}/calendar` - Room calendar (JSON feed)
- `GET /rooms/{id}/calendar.ics` - ICS export

### Authentication

Use Bearer tokens in `Authorization` header.

## Development

### Project Structure

```
.
├── cmd/server/          # Main application
├── internal/
│   ├── auth/           # Authentication service
│   ├── config/         # Configuration
│   ├── http/handlers/  # HTTP handlers
│   ├── msgraph/        # Graph client
│   ├── repository/     # Data access
│   └── ...
├── migrations/         # DB migrations
├── web/                # Static files and templates
├── pkg/                # Shared utilities
├── scripts/            # Helper scripts
├── Dockerfile
├── docker-compose.yaml
├── Jenkinsfile
├── Makefile
└── openapi.yaml
```

### Adding New Features

1. Define API in `openapi.yaml`
2. Add repository methods
3. Implement service logic
4. Create HTTP handler
5. Update UI if needed

### Testing

- Unit tests: `*_test.go` files
- Integration tests: Use `-tags=integration`
- Coverage: `go test -cover`

## Deployment

### Docker Compose

```bash
docker-compose up -d
```

### Jenkins CI/CD

The `Jenkinsfile` provides:

- Code linting and formatting
- Unit and integration tests
- Docker image build and push
- Automated deployment

### Production Considerations

- Use PostgreSQL with connection pooling
- Enable HTTPS with proper certificates
- Configure rate limiting and CORS
- Set up monitoring and logging
- Use secrets management for credentials

## Security

- Passwords hashed with Argon2id
- JWT tokens with expiration
- CSRF protection (planned)
- Rate limiting middleware
- Input validation and sanitization
- Audit logging for compliance

## Performance

- Database indexes on frequently queried columns
- Connection pooling for PostgreSQL
- Efficient queries with proper joins
- Caching for static data (planned)

## Scalability

- Horizontal scaling with multiple instances
- Database sharding by office_id (future)
- Read replicas for high read loads
- Async processing for Graph sync

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Submit a pull request

## License

MIT License

## Roadmap

- [ ] TOTP 2FA
- [ ] Email notifications
- [ ] Mobile app
- [ ] Advanced reporting
- [ ] Multi-office support
- [ ] Calendar sync with Google
