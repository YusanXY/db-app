package response

import "time"

type CommentResponse struct {
	ID          uint64            `json:"id"`
	ArticleID   uint64            `json:"article_id"`
	Content     string            `json:"content"`
	ContentHTML string            `json:"content_html"`
	User        UserResponse      `json:"user"`
	ParentID    *uint64           `json:"parent_id"`
	Parent      *CommentResponse  `json:"parent,omitempty"`
	LikeCount   int               `json:"like_count"`
	ReplyCount  int               `json:"reply_count"`
	IsLiked     bool              `json:"is_liked"`
	Replies     []CommentResponse `json:"replies,omitempty"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type CommentListResponse struct {
	Items      []CommentResponse `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

