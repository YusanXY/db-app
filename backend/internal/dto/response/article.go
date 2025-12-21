package response

import "time"

type ArticleResponse struct {
	ID            uint64         `json:"id"`
	Title         string         `json:"title"`
	Slug          string         `json:"slug"`
	Content       string         `json:"content,omitempty"`
	ContentHTML   string         `json:"content_html,omitempty"`
	Summary       string         `json:"summary"`
	CoverImageURL string         `json:"cover_image_url"`
	Author        *UserResponse  `json:"author"`
	Editor        *UserResponse  `json:"editor,omitempty"`
	Categories    []CategoryResponse `json:"categories"`
	Tags          []TagResponse   `json:"tags"`
	ViewCount     int            `json:"view_count"`
	LikeCount     int            `json:"like_count"`
	CommentCount  int            `json:"comment_count"`
	IsFeatured    bool           `json:"is_featured"`
	IsLiked       bool           `json:"is_liked,omitempty"`
	Status        string         `json:"status"`
	PublishedAt   *time.Time     `json:"published_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ArticleListResponse struct {
	Items      []*ArticleResponse `json:"items"`
	Pagination Pagination         `json:"pagination"`
}


