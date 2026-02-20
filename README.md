# Email Verifier

A robust, high-performance email verification service built with Go. This service provides a comprehensive suite of checks to validate email addresses, ensuring higher deliverability rates and cleaner contact lists.

## Features

- **Syntax Validation**: Ensures email addresses conform to standard formats.
- **MX Record Lookup**: Verifies that the domain has valid Mail Exchange records.
- **SMTP Handshake**: Connects to the mail server to verify existence without sending an actual email.
- **Disposable Email Detection**: Identifies and flags temporary or disposable email domains.
- **Catch-All Detection**: Detects if a domain accepts all incoming emails (catch-all configuration).
- **Rate Limiting**: Built-in protection against abuse.
- **Worker Pool Architecture**: Efficiently handles concurrent verification requests.
- **Dockerized**: Ready for deployment with Docker and Docker Compose.

## Prerequisites

- [Go](https://go.dev/) 1.22+
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) or [Just](https://github.com/casey/just) (optional, but recommended for development workflow)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/thelazydo/email-verifier.git
cd email-verifier
```

### 2. Setup Configuration

Initialize the project configuration by copying the example environment file.

Using Make:
```bash
make setup
```

Using Just:
```bash
just setup
```

Or manually:
```bash
cp .env.schema .env
```

### 3. Configuration Variables

Configure the `.env` file with your specific settings:
NOTE: Change localhost to the name of your database service when running docker-compose

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:password@localhost:5432/email_verifier?sslmode=disable` |
| `PORT` | Service port | `8080` |
| `SMTP_DOMAIN` | Domain used for EHLO handshake | `localhost` |
| `SMTP_ADDRESS` | Email address used in MAIL FROM | `verify@example.com` |
| `RATE_LIMIT` | Requests per second limit | `5` |
| `WORKERS_COUNT` | Number of concurrent workers | `5` |

## Usage

### Run Locally (Development)

Ensure you have a PostgreSQL instance running or update the `DATABASE_URL` to point to a valid database.

```bash
make run
# or
just run
```

### Run with Docker

This will spin up the API and a PostgreSQL database in containers.

```bash
make up
# or
just up
```

To stop the containers:
```bash
make down
# or
just down
```

### Build Binary

Compile the application into a static binary located in `./bin/api`.

```bash
make build
# or
just build
```

## API Endpoints

### 1. Verify Email

**Endpoint:** `POST /verify`

**Request Body:**
```json
{
  "email": "test@example.com"
  "webhook_url": "https://localhost:8081/aikomail" // (optional)
}
```

**Response:**
```json
{
    "message": "job received successfully",
    "success": true,
    "data": {
        "job_id": "0f021d27-fb08-4e2d-a77f-7c928ab8dec4"
    }
}
```

### 2. Service Status

**Endpoint:** `GET /status/:job_id`

**Response:**
```json
{
    "message": "job fetched successfully",
    "success": true,
    "data": {
        "id": "ca166dfa-77e0-44de-8eea-ed563186dde3",
        "email": "dev@smartx.ng",
        "domain": "smartx.ng",
        "status": "processing",
        "result": "valid", // "error", "unknown", "invalid"
        "webhook_url": "https://localhost:8081/aikomail",
    }
}
```


### 2. Service Status

**Endpoint:** `GET /`

**Response:**
```json
{
    "message": "Email Verifier API",
    "success": true,
    "data": {
        "/status": "Allowed methods - [GET]",
        "/verify": "Allowed methods - [POST]",
    }
}
```

## Testing

Run the test suite including race condition detection.

```bash
make test
# or
just test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
