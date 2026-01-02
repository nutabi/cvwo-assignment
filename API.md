# API Documentation

## Table of Contents

- [How to Use?](#how-to-use)
  - [Base URL](#base-url)
  - [Authentication](#authentication)
  - [Request Format](#request-format)
  - [Response Format](#response-format)
  - [Error Handling](#error-handling)
  - [Pagination](#pagination)

- [Authentication & Users](#authentication--users)
  - [User Registration](#user-registration)
  - [Login](#login)
  - [Refresh Token](#refresh-token)
  - [Current User Profile](#current-user-profile)
  - [User Public Profile](#user-public-profile)
  - [Update User Profile](#update-user-profile)

- [Topics](#topics)
  - [List Topics](#list-topics)
  - [Create Topic](#create-topic)
  - [Get Topic Details](#get-topic-details)
  - [Update Topic](#update-topic)
  - [Delete Topic](#delete-topic)

- [Posts](#posts)
  - [List Posts](#list-posts)
  - [Create Post Inside a Topic](#create-post-inside-a-topic)
  - [Get Post Details](#get-post-details)
  - [Update Post](#update-post)
  - [Delete Post](#delete-post)

- [Comments](#comments)
  - [List Comments](#list-comments)
  - [Create Comment on a Post](#create-comment-on-a-post)
  - [Get Comment Details](#get-comment-details)
  - [Update Comment](#update-comment)
  - [Delete Comment](#delete-comment)

---

## How to Use?

### Base URL

All API endpoints are prefixed with `/v1`:

```
http://localhost:8080/v1
```

### Authentication

The API uses **JWT (JSON Web Tokens)** for authentication. After logging in, you'll receive an access token that must be included in subsequent requests.

#### Including the Token

Include the JWT token in one of two ways:

1. **Authorization Header** (recommended):
   ```
   Authorization: Bearer <your-token>
   ```

2. **Cookie**:
   The API automatically sets an `httpOnly` cookie named `jwt` upon login.

#### Protected Endpoints

Endpoints marked with 🔒 require authentication. Unauthenticated requests will receive a `401 Unauthorized` response.

### Request Format

The API accepts both **JSON** and **form-encoded** data:

- **Content-Type**: `application/json` or `application/x-www-form-urlencoded`
- **Accept**: `application/json`

**JSON Example:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Form Example:**
```
username=john_doe&email=john@example.com&password=SecurePass123!
```

### Response Format

All responses are in JSON format.

**Success Response:**
```json
{
  "user_id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2026-01-01T12:00:00Z"
}
```

**Error Response:**
```json
{
  "error": "username taken"
}
```

### Error Handling

The API uses standard HTTP status codes:

| Status Code | Meaning |
|------------|---------|
| `200 OK` | Request successful |
| `201 Created` | Resource created successfully |
| `204 No Content` | Request successful, no response body |
| `400 Bad Request` | Invalid request body or missing required fields |
| `401 Unauthorized` | Missing or invalid authentication token |
| `403 Forbidden` | Authenticated but not authorized (e.g., trying to edit someone else's post) |
| `404 Not Found` | Resource not found |
| `409 Conflict` | Resource already exists (e.g., username/email taken) |
| `422 Unprocessable Entity` | Validation failed (e.g., invalid email format) |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Server error |

**Rate Limiting:**
- Authentication endpoints: 5 requests per minute per IP
- General endpoints: 100 requests per minute per IP

When rate limited, the response includes:
```json
{
  "error": "Too many authentication attempts. Please try again later.",
  "retry_after": 45.5
}
```

### Pagination

List endpoints support pagination via query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | Maximum number of items to return |
| `offset` | integer | 0 | Number of items to skip |

**Example:**
```
GET /v1/posts?limit=10&offset=20
```

---

---

## Authentication & Users

### User Registration

Register a new user account.

- **Endpoint**: `POST /v1/auth/register`
- **Authentication**: Not required
- **Rate Limit**: 5 requests per minute

**Request Body:**

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `username` | string | Yes | 3-20 characters, alphanumeric, must start with a letter |
| `email` | string | Yes | Valid email format |
| `password` | string | Yes | 8-32 characters, must contain uppercase, lowercase, and digit |

**Example Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Success Response** (`201 Created`):
```json
{
  "user_id": 1,
  "created_at": "2026-01-02T14:32:43.085572+07:00",
  "updated_at": "2026-01-02T14:32:43.085572+07:00",
  "username": "john_doe",
  "email": "john@example.com",
  "avatar_url": null,
  "bio": null
}
```

**Error Responses:**
- `409 Conflict`: Username or email already exists
- `422 Unprocessable Entity`: Validation failed (invalid username, email, or password format)

---

### Login

Authenticate and receive access tokens.

- **Endpoint**: `POST /v1/auth/login`
- **Authentication**: Not required
- **Rate Limit**: 5 requests per minute

**Request Body:**

| Field | Type | Required |
|-------|------|----------|
| `username` | string | Yes |
| `password` | string | Yes |

**Example Request:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Success Response** (`200 OK`):
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 300,
  "refresh_token": "ySe6IUytVS5Skf5GMOQJKGmXT8SFNUztZWhRVtR2Liw=",
  "token_type": "Bearer"
}
```

The response also sets an `httpOnly` cookie named `jwt` containing the access token.

**Token Details:**
- **expires_in**: 300 seconds (5 minutes)
- **Refreshable for**: 7 days using the refresh_token

**Error Responses:**
- `401 Unauthorized`: Invalid username or password
  ```json
  {"error": "incorrect Username or Password"}
  ```

---

### Refresh Token

Refresh an expired access token.

- **Endpoint**: `POST /v1/auth/refresh`
- **Authentication**: Required (expired token acceptable)
- **Rate Limit**: 5 requests per minute

**Headers:**
```
Authorization: Bearer <expired-or-valid-token>
```

**Success Response** (`200 OK`):
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 300,
  "refresh_token": "ySe6IUytVS5Skf5GMOQJKGmXT8SFNUztZWhRVtR2Liw=",
  "token_type": "Bearer"
}
```

**Error Responses:**
- `401 Unauthorized`: Token is not refreshable (older than 7 days) or invalid

---

### Get Current User Profile

Get the authenticated user's full profile (includes private information).

- **Endpoint**: `GET /v1/users/me`
- **Authentication**: 🔒 Required

**Success Response** (`200 OK`):
```json
{
  "user_id": 1,
  "created_at": "2026-01-02T14:32:43.085572+07:00",
  "updated_at": "2026-01-02T14:32:43.085572+07:00",
  "username": "john_doe",
  "email": "john@example.com",
  "avatar_url": "https://example.com/avatar.jpg",
  "bio": "Software developer and forum enthusiast"
}
```

**Error Responses:**
- `401 Unauthorized`: Not authenticated

---

### Get User Public Profile

Get a user's public profile (excludes email and updated_at).

- **Endpoint**: `GET /v1/users/{user_id}`
- **Authentication**: Not required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `user_id` | integer | User ID |

**Example Request:**
```
GET /v1/users/1
```

**Success Response** (`200 OK`):
```json
{
  "user_id": 1,
  "created_at": "2026-01-02T14:32:43.085572+07:00",
  "username": "john_doe",
  "avatar_url": "https://example.com/avatar.jpg",
  "bio": "Software developer and forum enthusiast"
}
```

**Error Responses:**
- `404 Not Found`: User does not exist
  ```json
  {"error": "user not found"}
  ```

---

### Update User Profile

Update the authenticated user's profile.

- **Endpoint**: `PATCH /v1/users/me`
- **Authentication**: 🔒 Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `avatar_url` | string | No | URL to profile picture |
| `bio` | string | No | User biography |

At least one field must be provided.

**Example Request:**
```json
{
  "avatar_url": "https://example.com/new-avatar.jpg",
  "bio": "Updated bio text"
}
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `422 Unprocessable Entity`: No fields provided

---

---

## Topics

### List Topics

Get a list of topics with optional filtering.

- **Endpoint**: `GET /v1/topics`
- **Authentication**: Not required

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | Number of topics to return |
| `offset` | integer | 0 | Number of topics to skip |
| `user_id` | integer | - | Filter by topic author |

**Example Request:**
```
GET /v1/topics?limit=10&offset=0&user_id=1
```

**Success Response** (`200 OK`):
```json
[
  {
    "topic_id": 1,
    "created_at": "2026-01-02T14:33:18.768276+07:00",
    "updated_at": "2026-01-02T14:33:18.768276+07:00",
    "name": "General Discussion",
    "description": "A place for general discussions",
    "author": {
      "user_id": 1,
      "created_at": "2026-01-02T14:32:43.085572+07:00",
      "username": "john_doe",
      "avatar_url": null,
      "bio": null
    }
  }
]
```

---

### Create Topic

Create a new discussion topic.

- **Endpoint**: `POST /v1/topics`
- **Authentication**: 🔒 Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Topic name |
| `description` | string | No | Topic description |

**Example Request:**
```json
{
  "title": "General Discussion",
  "description": "A place for general discussions"
}
```

**Success Response** (`201 Created`):
```json
{
  "topic_id": 1,
  "created_at": "2026-01-02T14:33:18.768276+07:00",
  "updated_at": "2026-01-02T14:33:18.768276+07:00",
  "name": "General Discussion",
  "description": "A place for general discussions",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Not authenticated

---

### Get Topic Details

Get details of a specific topic.

- **Endpoint**: `GET /v1/topics/{topic_id}`
- **Authentication**: Not required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `topic_id` | integer | Topic ID |

**Example Request:**
```
GET /v1/topics/1
```

**Success Response** (`200 OK`):
```json
{
  "topic_id": 1,
  "created_at": "2026-01-02T14:33:18.768276+07:00",
  "updated_at": "2026-01-02T14:33:18.768276+07:00",
  "name": "General Discussion",
  "description": "A place for general discussions",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  }
}
```

**Error Responses:**
- `404 Not Found`: Topic does not exist
  ```json
  {"error": "topic not found"}
  ```

---

### Update Topic

Update a topic's details. Only the topic author can update it.

- **Endpoint**: `PATCH /v1/topics/{topic_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `topic_id` | integer | Topic ID |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | No | Updated topic name |
| `description` | string | No | Updated description |

At least one field must be provided.

**Example Request:**
```json
{
  "title": "Updated Topic Name",
  "description": "Updated description"
}
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the topic author
- `404 Not Found`: Topic does not exist
- `422 Unprocessable Entity`: No fields provided

---

### Delete Topic

Delete a topic. Only the topic author can delete it. All posts and comments in the topic will be cascade deleted.

- **Endpoint**: `DELETE /v1/topics/{topic_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `topic_id` | integer | Topic ID |

**Example Request:**
```
DELETE /v1/topics/1
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the topic author
- `404 Not Found`: Topic does not exist

---

---

## Posts

### List Posts

Get a list of posts with optional filtering.

- **Endpoint**: `GET /v1/posts`
- **Authentication**: Not required

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | Number of posts to return |
| `offset` | integer | 0 | Number of posts to skip |
| `topic_id` | integer | - | Filter by topic |
| `user_id` | integer | - | Filter by post author |

**Example Request:**
```
GET /v1/posts?limit=10&offset=0&topic_id=1
```

**Success Response** (`200 OK`):
```json
[
  {
    "post_id": 1,
    "created_at": "2026-01-02T14:33:36.48214+07:00",
    "updated_at": "2026-01-02T14:33:36.48214+07:00",
    "title": "Welcome to the forum",
    "content": "This is my first post",
    "author": {
      "user_id": 1,
      "created_at": "2026-01-02T14:32:43.085572+07:00",
      "username": "john_doe",
      "avatar_url": null,
      "bio": null
    },
    "topic_id": 1
  }
]
```

---

### Create Post Inside a Topic

Create a new post within a topic.

- **Endpoint**: `POST /v1/topics/{topic_id}/posts`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `topic_id` | integer | Topic ID where the post will be created |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Post title |
| `content` | string | No | Post content |

**Example Request:**
```json
{
  "title": "Welcome to the forum",
  "content": "This is my first post"
}
```

**Success Response** (`201 Created`):
```json
{
  "post_id": 1,
  "created_at": "2026-01-02T14:33:36.48214+07:00",
  "updated_at": "2026-01-02T14:33:36.48214+07:00",
  "title": "Welcome to the forum",
  "content": "This is my first post",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  },
  "topic_id": 1
}
```

**Error Responses:**
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: Topic does not exist

---

### Get Post Details

Get details of a specific post.

- **Endpoint**: `GET /v1/posts/{post_id}`
- **Authentication**: Not required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `post_id` | integer | Post ID |

**Example Request:**
```
GET /v1/posts/1
```

**Success Response** (`200 OK`):
```json
{
  "post_id": 1,
  "created_at": "2026-01-02T14:33:36.48214+07:00",
  "updated_at": "2026-01-02T14:33:36.48214+07:00",
  "title": "Welcome to the forum",
  "content": "This is my first post",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  },
  "topic_id": 1
}
```

**Error Responses:**
- `404 Not Found`: Post does not exist
  ```json
  {"error": "post not found"}
  ```

---

### Update Post

Update a post's details. Only the post author can update it.

- **Endpoint**: `PATCH /v1/posts/{post_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `post_id` | integer | Post ID |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | No | Updated post title |
| `content` | string | No | Updated post content |

At least one field must be provided.

**Example Request:**
```json
{
  "title": "Updated Title",
  "content": "Updated content"
}
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the post author
- `404 Not Found`: Post does not exist
- `422 Unprocessable Entity`: No fields provided

---

### Delete Post

Delete a post. Only the post author can delete it. All comments on the post will be cascade deleted.

- **Endpoint**: `DELETE /v1/posts/{post_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `post_id` | integer | Post ID |

**Example Request:**
```
DELETE /v1/posts/1
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the post author
- `404 Not Found`: Post does not exist

---

---

## Comments

### List Comments

Get a list of comments with optional filtering.

- **Endpoint**: `GET /v1/comments`
- **Authentication**: Not required

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | Number of comments to return |
| `offset` | integer | 0 | Number of comments to skip |
| `post_id` | integer | - | Filter by post |
| `user_id` | integer | - | Filter by comment author |

**Example Request:**
```
GET /v1/comments?limit=10&offset=0&post_id=1
```

**Success Response** (`200 OK`):
```json
[
  {
    "comment_id": 1,
    "created_at": "2026-01-02T14:33:45.900668+07:00",
    "updated_at": "2026-01-02T14:33:45.900668+07:00",
    "content": "Great post!",
    "author": {
      "user_id": 1,
      "created_at": "2026-01-02T14:32:43.085572+07:00",
      "username": "john_doe",
      "avatar_url": null,
      "bio": null
    },
    "post_id": 1
  }
]
```

---

### Create Comment on a Post

Create a new comment on a post.

- **Endpoint**: `POST /v1/posts/{post_id}/comments`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `post_id` | integer | Post ID where the comment will be created |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `content` | string | Yes | Comment content |

**Example Request:**
```json
{
  "content": "Great post!"
}
```

**Success Response** (`201 Created`):
```json
{
  "comment_id": 1,
  "created_at": "2026-01-02T14:33:45.900668+07:00",
  "updated_at": "2026-01-02T14:33:45.900668+07:00",
  "content": "Great post!",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  },
  "post_id": 1
}
```

**Error Responses:**
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: Post does not exist

---

### Get Comment Details

Get details of a specific comment.

- **Endpoint**: `GET /v1/comments/{comment_id}`
- **Authentication**: Not required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `comment_id` | integer | Comment ID |

**Example Request:**
```
GET /v1/comments/1
```

**Success Response** (`200 OK`):
```json
{
  "comment_id": 1,
  "created_at": "2026-01-02T14:33:45.900668+07:00",
  "updated_at": "2026-01-02T14:33:45.900668+07:00",
  "content": "Great post!",
  "author": {
    "user_id": 1,
    "created_at": "2026-01-02T14:32:43.085572+07:00",
    "username": "john_doe",
    "avatar_url": null,
    "bio": null
  },
  "post_id": 1
}
```

**Error Responses:**
- `404 Not Found`: Comment does not exist
  ```json
  {"error": "comment not found"}
  ```

---

### Update Comment

Update a comment's content. Only the comment author can update it.

- **Endpoint**: `PATCH /v1/comments/{comment_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `comment_id` | integer | Comment ID |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `content` | string | Yes | Updated comment content |

**Example Request:**
```json
{
  "content": "Updated comment text"
}
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the comment author
- `404 Not Found`: Comment does not exist

---

### Delete Comment

Delete a comment. Only the comment author can delete it.

- **Endpoint**: `DELETE /v1/comments/{comment_id}`
- **Authentication**: 🔒 Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `comment_id` | integer | Comment ID |

**Example Request:**
```
DELETE /v1/comments/1
```

**Success Response** (`204 No Content`):
No response body.

**Error Responses:**
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: User is not the comment author
- `404 Not Found`: Comment does not exist