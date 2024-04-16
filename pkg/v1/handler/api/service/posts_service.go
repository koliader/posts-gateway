package service

type CreatePostReq struct {
	Owner string `json:"owner" binding:"required,email"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}
