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

- As a Registered User who forgot my password, I can reset it by entering an
admin-provided recovery code, so that I can regain access to my account.

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

Beside the web framwork, I have also decided to use:

- `jwt`: For stateless authentication.

- `bun`: For SQL-agnostic database interface.

### 1.2.2. Database

SQLite has been my go-to database for all of my web apps, be it written in
Python, Rust, or now, Go. For this project, it seems like a no-brainer. It's
fast, serverless, and easy-to-use. However, since this is an assignment, I
have decided to try out something new, something that is most commonly used in
newer enterprise applications: PostgreSQL.

Furthermore, PostgreSQL has a lot of out-of-the-box features that I used to
implement manually in SQLite, such as timestamp, advanced triggers, etc.

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

## 1.3. Relational Data Models

The followings were considered during design:

- Cascading deletions simplify things greatly, i.e. there is no need to deal
with orphans if there are no orphans. When, or if, moderator roles are added,
changes will be made to ownership of topics, posts, and comments on deletion.

- Non-nullability constraints have meaning, i.e. `null` is not the same as "".

- Soft-deletion (with `is_deleted`) is used to allow archival and undoing,
except for users. This also allows re-usability of usernames.

### 1.3.1. `user`

| name          | type      | constraints      |
|---------------|-----------|------------------|
| user_id       | integer   | primary key      |
| username      | string    | not null, unique |
| password_hash | string    | not null         |
| created_at    | timestamp | not null         |
| updated_at    | timestamp | not null         |

### 1.3.2. `recovery`

| name        | type      | constraints                                        |
|-------------|-----------|----------------------------------------------------|
| recovery_id | integer   | primary key                                        |
| user_id     | integer   | foreign key user(user_id) delete cascade, not null |
| code        | string    | not null, unique                                   |
| created_at  | timestamp | not null                                           |
| used_at     | timestamp |                                                    |

### 1.3.3. `topic`

| name        | type      | constraints                                        |
|-------------|-----------|----------------------------------------------------|
| topic_id    | integer   | primary key                                        |
| owned_by    | integer   | foreign key user(user_id) delete cascade, not null |
| title       | string    | not null, unique                                   |
| description | string    |                                                    |
| is_deleted  | boolean   | not null, default false                            |
| created_at  | timestamp | not null                                           |
| updated_at  | timestamp | not null                                           |

### 1.3.4. `post`

| name       | type      | constraints                                          |
|------------|-----------|------------------------------------------------------|
| post_id    | integer   | primary key                                          |
| owned_by   | integer   | foreign key user(user_id) delete cascade, not null   |
| topic_id   | integer   | foreign key topic(topic_id) delete cascade, not null |
| title      | string    | not null                                             |
| body       | string    |                                                      |
| is_deleted | boolean   | not null, default false                              |
| created_at | timestamp | not null                                             |
| updated_at | timestamp | not null                                             |

### 1.3.5. `comment`

| name       | type      | constraints                                        |
|------------|-----------|----------------------------------------------------|
| comment_id | integer   | primary key                                        |
| owned_by   | integer   | foreign key user(user_id) delete cascade, not null |
| post_id    | integer   | foreign key post(post_id) delete cascade, not null |
| text       | string    | not null                                           |
| is_deleted | boolean   | not null, default false                            |
| created_at | timestamp | not null                                           |
| updated_at | timestamp | not null                                           |

## 1.4. Project Structures

The frontend and backend code will be organised into a single monorepo.
Since this is my first time using these languages/frameworks, I'll stick to the
structures as suggested by LLMs and modified to my preference and needs.

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers/controllers
│   ├── models/              # Data models/structs
│   ├── services/            # Business logic layer
│   └── repository/          # Database access layer
├── migrations/              # Database migrations
├── go.mod
└── go.sum

frontend/
├── public/
│   ├── index.html
│   └── favicon.ico
├── src/
│   ├── api/                   # API client & HTTP requests
│   │   ├── client.ts          # Axios/fetch config
│   │   └── userApi.ts         # User-related API calls
│   ├── components/            # Reusable UI components
│   │   ├── common/            # Generic components
│   │   │   ├── Button.tsx
│   │   │   └── Input.tsx
│   │   └── layout/            # Layout components
│   │       ├── Header.tsx
│   │       └── Sidebar.tsx
│   ├── features/              # Feature-based modules
│   │   └── users/
│   │       ├── UserList.tsx
│   │       ├── UserForm.tsx
│   │       └── userSlice.ts   # With Redux
│   ├── hooks/                 # Custom React hooks
│   │   └── useAuth.ts
│   ├── pages/                 # Page/route components
│   │   ├── HomePage.tsx
│   │   ├── UserPage.tsx
│   │   └── NotFoundPage.tsx
│   ├── services/              # Business logic
│   │   └── authService.ts
│   ├── store/                 # State management (with Redux)
│   │   └── store.ts
│   ├── types/                 # TypeScript interfaces
│   │   └── user.ts
│   ├── utils/                 # Utility functions
│   │   └── validators.ts
│   ├── App.tsx                # Main app component
│   ├── main.tsx               # Entry point
│   └── routes.tsx             # Route definitions
└── package.json
```

# 2. Executions

To be written...

# 3. AI Usage Disclosure

- Asked Claude, "What is the best web framework in Go?"

- Asked Claude, "What are some commonly used third-party utility packages for
web development in Go?"

- Asked Claude and Gemini, "How to organise a Go web app?"

- Asked Claude and Gemini, "How to organise a React.js app?"

- Asked Claude, "Review the following relational data models and suggest
improvements."