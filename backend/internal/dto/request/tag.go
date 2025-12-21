package request

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Slug        string `json:"slug" binding:"required,min=1,max=50"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateTagRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type ListTagRequest struct {
	Keyword string `form:"keyword"`
	Sort    string `form:"sort"`    // name, article_count
	Order   string `form:"order"`   // asc, desc
	Limit   int    `form:"limit"`
}

