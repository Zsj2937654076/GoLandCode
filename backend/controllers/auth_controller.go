package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"student-management/middleware"
	"student-management/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthController 处理认证相关的端点
type AuthController struct {
	DB *sql.DB
}

// NewAuthController 创建新的 AuthController
func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{DB: db}
}

// LoginRequest 表示登录表单数据
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 表示登录成功后的响应
type LoginResponse struct {
	Token  string           `json:"token"`
	User   models.UserProfile `json:"user"`
}

// PasswordChangeRequest 表示修改密码的表单数据
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Login 处理 POST /api/auth/login 用户认证
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析请求体错误: %v", err)
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	log.Printf("用户登录尝试: %s", req.Username)
	log.Printf("请求体: %+v", req)

	// 验证必填字段
	if req.Username == "" || req.Password == "" {
		log.Printf("缺少必填字段: username=%s, password=%s", req.Username, req.Password)
		http.Error(w, "用户名和密码为必填项", http.StatusBadRequest)
		return
	}

	// 通过用户名查找用户
	user, err := models.GetUserByUsername(c.DB, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("用户不存在: %s", req.Username)
			http.Error(w, "无效的凭据", http.StatusUnauthorized)
		} else {
			log.Printf("数据库错误: %v", err)
			http.Error(w, "认证失败", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("找到用户: %+v", user)

	// 验证密码
	if !models.CheckPasswordHash(req.Password, user.Password) {
		log.Printf("用户密码无效: %s", req.Username)
		http.Error(w, "无效的凭据", http.StatusUnauthorized)
		return
	}

	log.Printf("用户密码验证成功: %s", req.Username)

	// 生成 JWT 令牌
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := middleware.Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		http.Error(w, "创建令牌失败", http.StatusInternalServerError)
		return
	}

	// 获取用户资料信息（不含密码）
	profile, err := models.GetUserProfile(c.DB, user.ID)
	if err != nil {
		http.Error(w, "获取用户资料失败", http.StatusInternalServerError)
		return
	}

	// 发送响应
	response := LoginResponse{
		Token: tokenString,
		User:  profile,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Profile 处理 GET /api/auth/profile 获取当前用户资料
func (c *AuthController) Profile(w http.ResponseWriter, r *http.Request) {
	// 从 JWT 声明中获取用户 ID
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		http.Error(w, "未授权", http.StatusUnauthorized)
		return
	}

	// 获取用户资料
	profile, err := models.GetUserProfile(c.DB, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在", http.StatusNotFound)
		} else {
			http.Error(w, "获取用户资料失败", http.StatusInternalServerError)
		}
		return
	}

	// 发送响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// ChangePassword 处理 POST /api/auth/change-password 更新用户密码
func (c *AuthController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 从 JWT 声明中获取用户 ID
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		http.Error(w, "未授权", http.StatusUnauthorized)
		return
	}

	// 解析请求体
	var req PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 验证请求
	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "旧密码和新密码为必填项", http.StatusBadRequest)
		return
	}

	// 获取用户以验证旧密码
	user, err := models.GetUserByID(c.DB, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在", http.StatusNotFound)
		} else {
			http.Error(w, "获取用户失败", http.StatusInternalServerError)
		}
		return
	}

	// 验证旧密码
	if !models.CheckPasswordHash(req.OldPassword, user.Password) {
		http.Error(w, "当前密码不正确", http.StatusBadRequest)
		return
	}

	// 修改密码
	err = models.ChangePassword(c.DB, claims.UserID, req.NewPassword)
	if err != nil {
		http.Error(w, "更新密码失败", http.StatusInternalServerError)
		return
	}

	// 发送成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "密码更新成功"})
}

// Logout 处理 POST /api/auth/logout（注意：JWT 令牌是无状态的，这主要是客户端操作）
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// 在实际应用中，您可能想要在服务器端将令牌加入黑名单
	
	// 发送成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "退出登录成功"})
} 