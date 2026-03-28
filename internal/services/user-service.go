package services

import (
    "errors"
    "time"
    
    "gorm.io/gorm"
    
	"LunoFuns-Server/internal/models"
    "LunoFuns-Server/internal/utils"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// CreateUser 创建用户（注册）
func (s *UserService) CreateUser(req *models.RegisterRequest) (*models.User, error) {
    // 检查用户名是否存在
    var existUser models.User
    if err := s.db.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
        return nil, errors.New("用户名已存在")
    }
    
    // 检查邮箱是否存在
    if err := s.db.Where("email = ?", req.Email).First(&existUser).Error; err == nil {
        return nil, errors.New("邮箱已注册")
    }
    
    // 加密密码
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, errors.New("密码加密失败")
    }
    
    // 创建用户
    user := &models.User{
        Username: req.Username,
        Password: hashedPassword,
        Nickname: req.Nickname,
        Email:    req.Email,
        Phone:    req.Phone,
        Status:   models.UserStatusEnabled,
        Role:     models.UserRoleNormal,
    }
    
    if err := s.db.Create(user).Error; err != nil {
        return nil, errors.New("创建用户失败")
    }
    
    return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    err := s.db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := s.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint64) (*models.User, error) {
    var user models.User
    err := s.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    return &user, nil
}

// Login 用户登录
func (s *UserService) Login(req *models.LoginRequest) (*models.User, error) {
    var user models.User
    
    // 支持用户名或邮箱登录
    err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户名或密码错误")
        }
        return nil, err
    }
    
    // 检查用户状态
    if !user.IsEnabled() {
        return nil, errors.New("账号已被禁用，请联系管理员")
    }
    
    // 验证密码
    if !utils.CheckPasswordHash(req.Password, user.Password) {
        return nil, errors.New("用户名或密码错误")
    }
    
    // 更新最后登录时间
    now := time.Now()
    user.LastLoginAt = &now
    s.db.Model(&user).Update("last_login_at", now)
    
    return &user, nil
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID uint64, req *models.UpdateProfileRequest) (*models.User, error) {
    user, err := s.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    
    updates := make(map[string]interface{})
    
    if req.Nickname != "" {
        updates["nickname"] = req.Nickname
    }
    if req.Avatar != "" {
        updates["avatar"] = req.Avatar
    }
    if req.Phone != "" {
        updates["phone"] = req.Phone
    }
    
    if len(updates) > 0 {
        if err := s.db.Model(user).Updates(updates).Error; err != nil {
            return nil, errors.New("更新用户信息失败")
        }
    }
    
    return s.GetUserByID(userID)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint64, req *models.ChangePasswordRequest) error {
    user, err := s.GetUserByID(userID)
    if err != nil {
        return err
    }
    
    // 验证旧密码
    if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
        return errors.New("原密码错误")
    }
    
    // 加密新密码
    newHashedPassword, err := utils.HashPassword(req.NewPassword)
    if err != nil {
        return errors.New("密码加密失败")
    }
    
    // 更新密码
    if err := s.db.Model(user).Update("password", newHashedPassword).Error; err != nil {
        return errors.New("修改密码失败")
    }
    
    return nil
}

// GetUserList 获取用户列表（分页）
func (s *UserService) GetUserList(page, pageSize int, status *int8, role *int8) (*models.UserListResponse, error) {
    var users []models.User
    var total int64
    
    query := s.db.Model(&models.User{})
    
    // 条件筛选
    if status != nil {
        query = query.Where("status = ?", *status)
    }
    if role != nil {
        query = query.Where("role = ?", *role)
    }
    
    // 统计总数
    if err := query.Count(&total).Error; err != nil {
        return nil, err
    }
    
    // 分页查询
    offset := (page - 1) * pageSize
    if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
        return nil, err
    }
    
    // 转换为响应对象
    userResponses := make([]*models.UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = user.ToResponse()
    }
    
    return &models.UserListResponse{
        Total: total,
        Users: userResponses,
    }, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint64, status int8) error {
    if status != models.UserStatusEnabled && status != models.UserStatusDisabled {
        return errors.New("无效的状态值")
    }
    
    result := s.db.Model(&models.User{}).Where("id = ?", userID).Update("status", status)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("用户不存在")
    }
    
    return nil
}

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(userID uint64, role int8) error {
    if role < models.UserRoleNormal || role > models.UserRoleAdmin {
        return errors.New("无效的角色值")
    }
    
    result := s.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("用户不存在")
    }
    
    return nil
}