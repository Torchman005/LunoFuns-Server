package models

type UploadVideoRequest struct {
	Title       string `json:"title" binding:"required,max=256"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url" binding:"required,url"`
	VideoURL    string `json:"video_url" binding:"required,url"`
	Duration    uint   `json:"duration" binding:"required"`
	Size        uint64 `json:"size" binding:"required"`
	Format      string `json:"format" binding:"required"`
	Width       uint16 `json:"width"`
	Height      uint16 `json:"height"`
	CategoryID  uint   `json:"category_id"`
	Tags        string `json:"tags"`
}

type VideoListRequest struct {
	Page       int  `form:"page" binding:"min=1"`
	PageSize   int  `form:"page_size" binding:"min=1,max=100"`
	CategoryID uint `form:"category_id"`
}

type VideoListResponse struct {
	Total int64   `json:"total"`
	Items []Video `json:"items"`
}

