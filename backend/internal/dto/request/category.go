package request

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Slug        string  `json:"slug" binding:"omitempty,min=1,max=100"`
	Description string  `json:"description"`
	ParentID    *uint64 `json:"parent_id"`
	IconURL     string  `json:"icon_url"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=1,max=100"`
	Description string  `json:"description"`
	ParentID    *uint64 `json:"parent_id"`
	IconURL     string  `json:"icon_url"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

type ListCategoryRequest struct {
	ParentID  *uint64 `form:"parent_id"`
	IsActive  *bool   `form:"is_active"`
	Tree      bool    `form:"tree"` // 是否返回树形结构
}

