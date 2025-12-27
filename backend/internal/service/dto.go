package service

import "time"

type UserProfile struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	AvatarUrl *string    `json:"avatar_url"`
	Bio       *string    `json:"bio"`
}

type TopicInfo struct {
	TopicID     uint        `json:"topic_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Author      UserProfile `json:"author"`
	Posts       []PostInfo  `json:"posts,omitempty"`
}

type PostInfo struct {
	PostID    uint          `json:"post_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Author    *UserProfile  `json:"author,omitempty"`
	Topic     *TopicInfo    `json:"topic,omitempty"`
	Comments  []CommentInfo `json:"comments,omitempty"`
}

type CommentInfo struct {
	CommentID uint         `json:"comment_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Content   string       `json:"content"`
	Author    *UserProfile `json:"author,omitempty"`
	Post      *PostInfo    `json:"post,omitempty"`
}
