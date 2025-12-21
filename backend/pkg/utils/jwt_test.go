package utils

import (
	"dbapp/internal/config"
	"os"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	// 设置测试配置
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret-key",
			ExpiresIn: 3600,
		},
	}

	token, err := GenerateJWT(1, "testuser", "user")
	if err != nil {
		t.Fatalf("生成JWT失败: %v", err)
	}

	if token == "" {
		t.Error("Token应该被生成")
	}
}

func TestParseJWT(t *testing.T) {
	// 设置测试配置
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret-key",
			ExpiresIn: 3600,
		},
	}

	// 生成token
	token, err := GenerateJWT(1, "testuser", "user")
	if err != nil {
		t.Fatalf("生成JWT失败: %v", err)
	}

	// 解析token
	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("解析JWT失败: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("期望UserID 1, 得到 %d", claims.UserID)
	}

	if claims.Username != "testuser" {
		t.Errorf("期望Username testuser, 得到 %s", claims.Username)
	}

	if claims.Role != "user" {
		t.Errorf("期望Role user, 得到 %s", claims.Role)
	}
}

func TestParseJWT_InvalidToken(t *testing.T) {
	// 设置测试配置
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret-key",
			ExpiresIn: 3600,
		},
	}

	// 解析无效token
	_, err := ParseJWT("invalid.token.here")
	if err == nil {
		t.Error("应该返回解析失败的错误")
	}
}

func TestParseJWT_WrongSecret(t *testing.T) {
	// 使用一个密钥生成token
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "secret-key-1",
			ExpiresIn: 3600,
		},
	}

	token, err := GenerateJWT(1, "testuser", "user")
	if err != nil {
		t.Fatalf("生成JWT失败: %v", err)
	}

	// 使用另一个密钥解析token
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "secret-key-2",
			ExpiresIn: 3600,
		},
	}

	_, err = ParseJWT(token)
	if err == nil {
		t.Error("应该返回签名验证失败的错误")
	}
}

func TestJWT_Expiration(t *testing.T) {
	// 设置测试配置（很短的过期时间）
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret-key",
			ExpiresIn: 1, // 1秒
		},
	}

	token, err := GenerateJWT(1, "testuser", "user")
	if err != nil {
		t.Fatalf("生成JWT失败: %v", err)
	}

	// 立即解析应该成功
	_, err = ParseJWT(token)
	if err != nil {
		t.Fatalf("解析JWT失败: %v", err)
	}

	// 注意：实际测试过期需要等待，这里只是演示结构
	// 在实际测试中，可以使用mock时间或等待
}

