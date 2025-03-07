package models

import (
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// User 表示系统中的用户
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 密码永远不会在 JSON 中暴露
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfile 是用户信息的简化版本，不包含密码
type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// HashPassword 使用 bcrypt 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 比较密码和哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	log.Printf("正在验证密码哈希")
	log.Printf("输入的密码: %s", password)
	log.Printf("存储的哈希值: %s", hash)
	log.Printf("哈希值长度: %d", len(hash))
	log.Printf("密码长度: %d", len(password))
	
	// 检查哈希值是否以正确的前缀开头
	if !strings.HasPrefix(hash, "$2a$") {
		log.Printf("无效的哈希格式：不以 $2a$ 开头")
		return false
	}
	
	// 如果存在重复的 $2a$ 前缀，则移除
	if strings.HasPrefix(hash, "$2a$2a$") {
		hash = strings.TrimPrefix(hash, "$2a$")
		log.Printf("修复重复前缀，新的哈希值: %s", hash)
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("密码验证失败: %v", err)
		return false
	}
	log.Printf("密码验证成功")
	return true
}

// GetUserByUsername 通过用户名获取用户信息
func GetUserByUsername(db *sql.DB, username string) (User, error) {
	var user User
	query := `
		SELECT id, username, password, email, role, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

// GetUserByID 通过 ID 获取用户信息
func GetUserByID(db *sql.DB, id int64) (User, error) {
	var user User
	query := `
		SELECT id, username, password, email, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

// GetUserProfile 获取用户资料（不包含密码）
func GetUserProfile(db *sql.DB, id int64) (UserProfile, error) {
	var profile UserProfile
	query := `
		SELECT id, username, email, role
		FROM users
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&profile.ID, &profile.Username, &profile.Email, &profile.Role,
	)
	return profile, err
}

// ChangePassword 更新用户密码
func ChangePassword(db *sql.DB, userID int64, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	query := `
		UPDATE users
		SET password = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err = db.Exec(query, hashedPassword, userID)
	return err
} 