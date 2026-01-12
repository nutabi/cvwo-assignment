# 1. Planning

## 1.1. User Requirements

These are the basic requirements. Hopefully, by the end of the assignment
period, I can add things like moderators role, orphaned posts/comments (like
Reddit), and email-based account recovery (if I ever figured out how to send
emails :D)

### 1.1.1. Authentication

- As a Visitor, I can create an account with a unique username and password,
so that I can participate in the forum.

- As a Registered User, I can sign in with my username and password, so that
I can access my account and create content.

- As a Registered User, I can sign out, so that I can securely end
my session.

### 1.1.2. Topics

- As a Registered User, I can create topics, so that I can
organize discussions.

- As a Visitor or Registered User, I can view topics, so that I can see
what others are discussing.

- As a Registered User, I can edit topics that I created, so that I can
update their titles and descriptions.

- As a Registered User, I can delete topics that I created, so that I can
remove discussions I no longer want to host.

### 1.1.3. Posts

- As a Registered User, I can create posts within a topic, so that I can
contribute to the discussion.

- As a Visitor or Registered User, I can view posts within a topic, so that
I can stay informed.

- As a Registered User, I can edit posts that I created, so that I
can correct mistakes and provide updates.

- As a Registered User, I can delete posts that I created, so that I
can remove content I no longer wish to display.

### 1.1.4. Comments

- As a Registered User, I can create comments on posts, so that I can
contribute to the conversation.

- As a Visitor or Registered User, I can view comments on posts, so that I
can read the perspectives of others.

- As a Registered User, I can edit comments that I created, so that I can
correct mistakes.

- As a Registered User, I can delete comments that I created, so that I can
remove content I no longer wish to display.

## 1.2. Tools

### 1.2.1. Backend

Per assignment requirement, the backend must be written in Go. Go has several
excellent web frameworks, but I have decided to try out Gin, reasons being:

- It is the most used one, so there is a good chance others have encountered
and resovled any issues I might have down the road.

- The documentation looks fantastic. This is especially important for someone
with no prior experience in Go.

- The latest commit on master branch is only days ago. While it means things
can be unstable from time to time, it also means there is an active community
behind it ready to answer questions.

Beside the web framework, I have also decided to use:

- `gin-jwt`: For stateless JWT-based authentication with automatic token refresh.

- `gorm`: For SQL-agnostic database ORM with migrations support.

- `argon2`: For secure password hashing using the Argon2id algorithm.

### 1.2.2. Database

SQLite has been my go-to database for all of my web apps, be it written in
Python, Rust, or now, Go. For this project, it's a perfect fit:

- **Serverless**: No separate database process to manage
- **Fast**: Excellent performance for read-heavy workloads like forums
- **Portable**: Single file database that's easy to backup and deploy
- **GORM Support**: Full ORM support with automatic migrations
- **Soft Deletes**: Built-in support via GORM's `DeletedAt` field

The database uses GORM's embedded `gorm.Model` which provides:
- Auto-incrementing ID (primary key)
- Timestamps (CreatedAt, UpdatedAt)
- Soft delete support (DeletedAt)

### 1.2.3. Frontend

Requirements call for TypeScript with React.js. 

Honestly, I despise writing web frontend code, so much so that I have resorted
to building mobile (iOS) apps instead of web apps for all of my toy backends.

I hope this is a good learning experience that can change my mind :D

### 1.2.4. Development

The standard tools are Git/GitHub for version control and Visual Studio Code
for code editor. I also installed LSP and linters for both Go and TypeScript.

### 1.2.5. Deployment

Requirement mentions AWS by name. However, I'll stick to DigitalOcean since I
have been using it for years now as a reverse proxy. I'm also using Cloudflare
to proxy traffic to the VPS, which in turn runs Nginx to allow SSL passthrough
to my home server.

The result is that SSL termination happens on my server, not on a $6/month
cloud VPS. My home server runs Caddy to reverse proxy HTTP traffic to different
services, hosted baremetal or as Docker containers. I'll use Docker for this
since it helps with Level 5 requirements.

### 1.2.6. LLMs

For someone with no prior experience in Go, JS/TS, React.js, PostgreSQL, using
LLMs is a must if I want to get proficient quickly. Currently, I'm using both
Claude and Gemini to help me learn and understand errors. AI disclosures will
be updated from time to time to reflect this usage.

## 1.3. API Error Handling

All API error responses follow a consistent structure with two fields:

- `error_code`: A machine-readable error code (e.g., `INVALID_INPUT`, `USER_NOT_FOUND`)
- `error_message`: A human-readable error message

Example error response:
```json
{
  "error_code": "INVALID_INPUT",
  "error_message": "invalid input"
}
```

### 1.3.1. Error Codes

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `UNAUTHORIZED` | 401 | User is not authenticated |
| `FORBIDDEN` | 403 | User lacks permission for the operation |
| `BAD_REQUEST` | 400 | Invalid request parameters |
| `INVALID_INPUT` | 400 | Invalid request body or input validation failed |
| `NOT_FOUND` | 404 | Generic resource not found |
| `USER_NOT_FOUND` | 404 | User does not exist |
| `TOPIC_NOT_FOUND` | 404 | Topic does not exist |
| `POST_NOT_FOUND` | 404 | Post does not exist |
| `COMMENT_NOT_FOUND` | 404 | Comment does not exist |
| `CONFLICT` | 409 | Generic conflict error |
| `USERNAME_TAKEN` | 409 | Username already exists |
| `EMAIL_IN_USE` | 409 | Email already registered |
| `TOPIC_TITLE_TAKEN` | 409 | Topic title already exists |
| `UNPROCESSABLE_ENTITY` | 422 | Request is well-formed but semantically incorrect |
| `NO_UPDATE_FIELDS` | 422 | No valid fields provided for update |
| `INTERNAL_SERVER_ERROR` | 500 | Internal server error |

## 1.4. Relational Data Models

The followings were considered during design:

- Cascading deletions simplify things greatly, i.e. there is no need to deal
with orphans if there are no orphans. When, or if, moderator roles are added,
changes will be made to ownership of topics, posts, and comments on deletion.

- Non-nullability constraints have meaning, i.e. `null` is not the same as "".

- Soft-deletion (with `is_deleted`) is used to allow archival and undoing,
except for users. This also allows re-usability of usernames.

### 1.4.1. `users`

| name       | type      | constraints      |
|------------|-----------|------------------|
| id         | integer   | primary key      |
| username   | string    | not null, unique |
| email      | string    | not null, unique |
| phc        | string    | not null (Argon2id hash) |
| avatar_url | string    | nullable         |
| bio        | text      | nullable         |
| created_at | timestamp | not null         |
| updated_at | timestamp | not null         |
| deleted_at | timestamp | nullable (soft delete) |

### 1.4.2. `topics`

| name        | type      | constraints                                       |
|-------------|-----------|---------------------------------------------------|
| id          | integer   | primary key                                       |
| name        | string    | not null, unique                                  |
| description | text      | nullable                                          |
| author_id   | integer   | foreign key users(id) cascade delete, not null    |
| created_at  | timestamp | not null                                          |
| updated_at  | timestamp | not null                                          |
| deleted_at  | timestamp | nullable (soft delete)                            |

### 1.4.3. `posts`

| name       | type      | constraints                                          |
|------------|-----------|------------------------------------------------------|
| id         | integer   | primary key                                          |
| title      | string    | not null                                             |
| content    | text      | nullable                                             |
| author_id  | integer   | foreign key users(id) cascade delete, not null       |
| topic_id   | integer   | foreign key topics(id) cascade delete, not null      |
| created_at | timestamp | not null                                             |
| updated_at | timestamp | not null                                             |
| deleted_at | timestamp | nullable (soft delete)                               |

### 1.4.4. `comments`

| name       | type      | constraints                                        |
|------------|-----------|----------------------------------------------------||
| id         | integer   | primary key                                        |
| content    | text      | not null                                           |
| author_id  | integer   | foreign key users(id) cascade delete, not null     |
| post_id    | integer   | foreign key posts(id) cascade delete, not null     |
| created_at | timestamp | not null                                           |
| updated_at | timestamp | not null                                           |
| deleted_at | timestamp | nullable (soft delete)                             |

## 1.5. Project Structures

The frontend and backend code will be organised into a single monorepo.
Since this is my first time using these languages/frameworks, I'll stick to the
structures as suggested by LLMs and modified to my preference and needs.

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/                 # Application setup and initialization
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP handlers/controllers
│   ├── middleware/          # HTTP middlewares
│   ├── model/               # Data models/structs
│   ├── repository/          # Database access layer
│   ├── service/             # Business logic layer
│   ├── test/                # Integration tests
│   └── utility/             # Utility functions
├── Dockerfile               # Dockerfile for backend
├── go.mod
└── go.sum
```

# 2. Installation

## 2.1. Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose (for containerized deployment)
- Git

## 2.2. Quick Start with Docker Compose

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd cvwo-assignment
   ```

2. Create environment file:
   ```bash
   cp example.env .env
   ```

3. Edit `.env` and update the `JWT_SECRET` to a secure random string:
   ```bash
   # Generate a random secret (on macOS/Linux):
   openssl rand -base64 32
   ```

4. Start the backend service:
   ```bash
   docker compose up -d
   ```

5. The backend API will be available at `http://localhost:8080`

6. View logs:
   ```bash
   docker compose logs -f backend
   ```

7. Stop the service:
   ```bash
   docker compose down
   ```

## 2.3. Manual Backend Setup (Development)

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the project root (copy from `example.env`):
   ```bash
   cp ../example.env ../.env
   ```

4. Update the `.env` file with appropriate values, especially:
   - `JWT_SECRET`: Set to a secure random string
   - `DEBUG`: Set to `true` for development
   - `DATABASE_URL`: Path to SQLite database file

5. Run the backend:
   ```bash
   go run cmd/api/main.go
   ```

6. The API will be available at `http://localhost:8080`

## 2.4. Running Tests

Run all tests in the backend:
```bash
cd backend
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```


# 3. AI Usage Disclosure

Regardless of whether AI tools were used and how much autonomy was given to
them, all AI-generated content was read, verified, and sometimes modified by me
before being included in the project. That being said, here are the details of AI
usage in this project.

## 3.1. Learning & Debugging

I used both Claude and Gemini extensively to help me learn Go and SQL. Specifically,
I have used them to:

- Understand Go syntax and semantics.

- Understand Gin framework and its ecosystem (third-party middlewares, ORM, etc.)

- Debug errors and issues encountered during development (Note: Agents were asked
to explain error messages and the approach to fixing them. Code was written by myself).

## 3.2. Code Quality Analysis

After every major feature was completed, I used AI agents to analyze the codebase
for code quality, security, and performance issues. The following prompt is used
to ensure that no code is generated by the agents:

```
Analyse [the feature] implementation. Pay close attention to consistency, potential
bugs, and error handling. Present the findings in order of decreasing severity.
For each issue found, provide a consise explanation and suggest improvements.
Do not generate any code.
```

I then decide if the suggestions are valid and relevant before implementing the
fixes myself.

## 3.3. Testing & Documentation

Agentic AI was used to autonomously generate code and non-code artifacts for
the following parts of the project:

- API documentation in [API.md](API.md).

- Unit tests in individual internal packages.

- Integration tests in `backend/internal/test/`.