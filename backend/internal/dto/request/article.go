package request

type CreateArticleRequest struct {
	Title         string   `json:"title" binding:"required,min=1,max=500"`
	Content       string   `json:"content" binding:"required"`
	Summary       string   `json:"summary"`
	CoverImageURL string   `json:"cover_image_url"`
	CategoryIDs   []uint64 `json:"category_ids"`
	TagIDs        []uint64 `json:"tag_ids"`
	Status        string   `json:"status" binding:"oneof=draft published"`
}

type UpdateArticleRequest struct {
	Title         string   `json:"title" binding:"min=1,max=500"`
	Content       string   `json:"content"`
	Summary       string   `json:"summary"`
	CoverImageURL string   `json:"cover_image_url"`
	CategoryIDs   []uint64 `json:"category_ids"`
	TagIDs        []uint64 `json:"tag_ids"`
	Status        string   `json:"status" binding:"oneof=draft published archived"`
}

type ListArticleRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	CategoryID uint64 `form:"category_id"`
	TagID      uint64 `form:"tag_id"`
	AuthorID   uint64 `form:"author_id"`
	Status     string `form:"status"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort"`
	Order      string `form:"order"`
}

