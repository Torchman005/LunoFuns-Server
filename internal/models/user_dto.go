package models

import (
    "time"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=32" validate:"alphanum"`
    Password string `json:"password" binding:"required,min=6,max=32"`
    Nickname string `json:"nickname" binding:"required,min=1,max=32"`
    Email    string `json:"email" binding:"required,email"`
    Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required"` // 用户名或邮箱
    Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
    Nickname string `json:"nickname" binding:"omitempty,min=1,max=32"`
    Avatar   string `json:"avatar" binding:"omitempty,url"`
    Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
    ID          uint64     `json:"id"`
    Username    string     `json:"username"`
    Nickname    string     `json:"nickname"`
    Avatar      string     `json:"avatar"`
    Email       string     `json:"email"`
    Phone       string     `json:"phone,omitempty"`
    Status      int8       `json:"status"`
    Role        int8       `json:"role"`
    LastLoginAt *time.Time `json:"last_login_at,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应对象
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:          u.ID,
        Username:    u.Username,
        Nickname:    u.Nickname,
        Avatar:      u.Avatar,
        Email:       u.Email,
        Phone:       u.Phone,
        Status:      u.Status,
        Role:        u.Role,
        LastLoginAt: u.LastLoginAt,
        CreatedAt:   u.CreatedAt,
        UpdatedAt:   u.UpdatedAt,
    }
}

// UserListResponse 用户列表响应
type UserListResponse struct {
    Total int64           `json:"total"`
    Users []*UserResponse `json:"users"`
}