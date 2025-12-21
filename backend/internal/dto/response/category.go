package response

import "time"

type CategoryResponse struct {
	ID           uint64             `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Description  string             `json:"description"`
	ParentID     *uint64            `json:"parent_id"`
	Parent       *CategoryResponse  `json:"parent,omitempty"`
	IconURL      string             `json:"icon_url"`
	SortOrder    int                `json:"sort_order"`
	ArticleCount int                `json:"article_count"`
	IsActive     bool               `json:"is_active"`
	Children     []CategoryResponse `json:"children,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

