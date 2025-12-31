# API Documentation

## Table of Contents

- [How to Use?](#how-to-use)

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

    - [List Posts Inside a Topic](#list-posts-inside-a-topic)

    - [Create Post Inside a Topic](#create-post-inside-a-topic)

    - [Get Post Details](#get-post-details)

    - [Update Post](#update-post)

    - [Delete Post](#delete-post)

- [Comments](#comments)

    - [List Comments on a Post](#list-comments-on-a-post)

    - [Create Comment on a Post](#create-comment-on-a-post)

    - [Get Comment Details](#get-comment-details)

    - [Update Comment](#update-comment)

    - [Delete Comment](#delete-comment)

## How to Use?

## Authentication & Users

### Register User

- Endpoint: `POST {ROOT}/auth/register`

### Login

- Endpoint: `POST {ROOT}/auth/login`

### Refresh Access Token

- Endpoint: `POST {ROOT}/auth/refresh`

### Get User Public Profile

- Endpoint: `GET {ROOT}/users/{user_id}`

### Get Current User Profile

- Endpoint: `GET {ROOT}/users/me`

### Update User Profile

- Endpoint: `PATCH {ROOT}/users/me`

## Topics

### Create Topic

- Endpoint: `POST {ROOT}/topics`

### List Topics

- Endpoint: `GET {ROOT}/topics`

### Get Topic Details

- Endpoint: `GET {ROOT}/topics/{topic_id}`

### Update Topic

- Endpoint: `PATCH {ROOT}/topics/{topic_id}`

### Delete Topic

- Endpoint: `DELETE {ROOT}/topics/{topic_id}`

## Posts

### Create Post

- Endpoint: `POST {ROOT}/topics/{topic_id}/posts`

### List Posts

- Endpoint: `GET {ROOT}/posts`

### Get Post Details

- Endpoint: `GET {ROOT}/posts/{post_id}`

### Update Post

- Endpoint: `PATCH {ROOT}/posts/{post_id}`

### Delete Post

- Endpoint: `DELETE {ROOT}/posts/{post_id}`

## Comments

### Create Comment

- Endpoint: `POST {ROOT}/posts/{post_id}/comments`

### List Comments

- Endpoint: `GET {ROOT}/comments`

### Get Comment Details

- Endpoint: `GET {ROOT}/comments/{comment_id}`

### Update Comment

- Endpoint: `PATCH {ROOT}/comments/{comment_id}`

### Delete Comment

- Endpoint: `DELETE {ROOT}/comments/{comment_id}`