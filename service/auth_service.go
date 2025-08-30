package service

import (
	"errors"
	"go_chat/global"
	"go_chat/model"

	"github.com/sirupsen/logrus"
)

type AuthService struct{}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	UserName string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 用户注册
func (s *AuthService) Register(req RegisterRequest) (*UserResponse, error) {
	db := global.GetMySQLClient()
	if db == nil {
		return nil, errors.New("数据库连接不可用")
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Where("user_name = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 创建新用户
	user := model.User{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   model.Active,
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		logrus.Error("密码加密失败:", err)
		return nil, errors.New("密码处理失败")
	}

	// 保存用户到数据库
	if err := db.Create(&user).Error; err != nil {
		logrus.Error("用户注册失败:", err)
		return nil, errors.New("注册失败，请重试")
	}

	logrus.Info("用户注册成功:", user.UserName)

	// 返回用户信息
	return &UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req LoginRequest) (*UserResponse, error) {
	db := global.GetMySQLClient()
	if db == nil {
		return nil, errors.New("数据库连接不可用")
	}

	// 查找用户
	var user model.User
	if err := db.Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != model.Active {
		return nil, errors.New("账户已被禁用")
	}

	logrus.Info("用户登录成功:", user.UserName)

	// 返回用户信息
	return &UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *AuthService) GetUserByID(userID uint) (*UserResponse, error) {
	db := global.GetMySQLClient()
	if db == nil {
		return nil, errors.New("数据库连接不可用")
	}

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, nil
}
