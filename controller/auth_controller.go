package controller

import (
	"go_chat/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var authService = service.NewAuthService()

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param register body service.RegisterRequest true "注册信息"
// @Success 200 {object} Response "注册成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 409 {object} Response "用户已存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/auth/register [post]
func Register(c *gin.Context) {
	var req service.RegisterRequest

	// 绑定和验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error("参数绑定失败:", err)
		Error(c, http.StatusBadRequest, "请求参数格式错误: "+err.Error())
		return
	}

	// 调用服务层处理注册逻辑
	user, err := authService.Register(req)
	if err != nil {
		logrus.Error("用户注册失败:", err)
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已被注册" {
			Error(c, http.StatusConflict, err.Error())
		} else {
			Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Success(c, user, "注册成功")
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body service.LoginRequest true "登录信息"
// @Success 200 {object} Response "登录成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 401 {object} Response "用户名或密码错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req service.LoginRequest

	// 绑定和验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error("参数绑定失败:", err)
		Error(c, http.StatusBadRequest, "请求参数格式错误: "+err.Error())
		return
	}

	// 调用服务层处理登录逻辑
	user, err := authService.Login(req)
	if err != nil {
		logrus.Error("用户登录失败:", err)
		if err.Error() == "用户名或密码错误" || err.Error() == "账户已被禁用" {
			Error(c, http.StatusUnauthorized, err.Error())
		} else {
			Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Success(c, user, "登录成功")
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id query string true "用户ID"
// @Success 200 {object} Response "获取成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 404 {object} Response "用户不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/user/profile [get]
func GetProfile(c *gin.Context) {
	// 从查询参数获取用户ID
	userIDStr := c.Query("id")
	if userIDStr == "" {
		Error(c, http.StatusBadRequest, "用户ID不能为空")
		return
	}

	// 转换用户ID
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 调用服务层获取用户信息
	user, err := authService.GetUserByID(uint(userID))
	if err != nil {
		logrus.Error("获取用户信息失败:", err)
		if err.Error() == "用户不存在" {
			Error(c, http.StatusNotFound, err.Error())
		} else {
			Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Success(c, user, "获取用户信息成功")
}
