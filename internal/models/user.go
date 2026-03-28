package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;comment:用户ID，自增主键" json:"id"`
    Username    string         `gorm:"type:varchar(32);not null;uniqueIndex:uk_username;comment:用户名，登录账号" json:"username"`
    Password    string         `gorm:"type:varchar(128);not null;comment:密码，bcrypt加密存储" json:"-"` // json:"-" 确保密码不会序列化返回
    Nickname    string         `gorm:"type:varchar(32);not null;comment:昵称，展示名称" json:"nickname"`
    Avatar      string         `gorm:"type:varchar(512);default:'';comment:头像URL，存储对象存储地址" json:"avatar"`
    Email       string         `gorm:"type:varchar(128);not null;uniqueIndex:uk_email;comment:邮箱，用于找回密码" json:"email"`
    Phone       string         `gorm:"type:varchar(20);default:'';comment:手机号" json:"phone"`
    Status      int8           `gorm:"type:tinyint;not null;default:1;index:idx_status;comment:状态：0-禁用，1-正常" json:"status"`
    Role        int8           `gorm:"type:tinyint;not null;default:0;comment:角色：0-普通用户，1-VIP，2-管理员" json:"role"`
    LastLoginAt *time.Time     `gorm:"type:datetime;comment:最后登录时间" json:"last_login_at,omitempty"`
    CreatedAt   time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at;comment:创建时间" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
    return "user"
}

// 状态常量
const (
    UserStatusDisabled = 0 // 禁用
    UserStatusEnabled  = 1 // 正常
)

// 角色常量
const (
    UserRoleNormal = 0 // 普通用户
    UserRoleVIP    = 1 // VIP用户
    UserRoleAdmin  = 2 // 管理员
)

// 辅助方法：判断用户状态
func (u *User) IsEnabled() bool {
    return u.Status == UserStatusEnabled
}

// 辅助方法：判断用户角色
func (u *User) IsAdmin() bool {
    return u.Role == UserRoleAdmin
}

func (u *User) IsVIP() bool {
    return u.Role == UserRoleVIP
}

// 辅助方法：禁用用户
func (u *User) Disable() {
    u.Status = UserStatusDisabled
}

// 辅助方法：启用用户
func (u *User) Enable() {
    u.Status = UserStatusEnabled
}

// 辅助方法：设置为管理员
func (u *User) SetAdmin() {
    u.Role = UserRoleAdmin
}

// 辅助方法：设置为VIP
func (u *User) SetVIP() {
    u.Role = UserRoleVIP
}