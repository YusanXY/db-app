package request

type CreateCommentRequest struct {
	ArticleID uint64  `json:"article_id" binding:"-"` // 从路径参数获取，不需要在请求体中验证
	Content   string  `json:"content" binding:"required,min=1,max=5000"`
	ParentID  *uint64 `json:"parent_id"` // 回复的评论ID，如果是顶级评论则为nil
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

type ListCommentRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

