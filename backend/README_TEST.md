# 测试文档

本文档说明如何运行和编写项目的单元测试。

## 测试结构

项目使用Go标准测试框架，测试文件以`_test.go`结尾。

```
backend/
├── internal/
│   ├── repository/
│   │   ├── user_repository_test.go
│   │   └── article_repository_test.go
│   ├── service/
│   │   ├── user_service_test.go
│   │   └── article_service_test.go
│   ├── handler/
│   │   ├── auth_handler_test.go
│   │   └── article_handler_test.go
│   ├── middleware/
│   │   └── auth_test.go
│   └── test/
│       └── test_helper.go  # 测试辅助工具
└── pkg/
    └── utils/
        └── jwt_test.go
```

## 运行测试

### 运行所有测试

```bash
go test ./...
```

### 运行测试并显示详细信息

```bash
go test ./... -v
```

### 运行特定包的测试

```bash
go test ./internal/service -v
go test ./internal/repository -v
```

### 运行特定测试函数

```bash
go test ./internal/service -run TestUserService_Register -v
```

### 生成测试覆盖率报告

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 使用Makefile

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-cover

# 竞态检测
make test-race
```

## 测试辅助工具

测试使用SQLite内存数据库，无需真实数据库连接。

### 设置测试数据库

```go
db := test.SetupTestDB(t)
defer test.TeardownTestDB(db)
```

### 创建测试数据

```go
// 创建测试用户
user := test.CreateTestUser(db, "username", "email@example.com")

// 创建测试文章
article := test.CreateTestArticle(db, userID, "文章标题")
```

## 测试覆盖范围

### Repository层测试

- ✅ UserRepository: Create, GetByID, GetByUsername, GetByEmail, Update, UpdateLastLogin
- ✅ ArticleRepository: Create, GetByID, GetBySlug, List, ListWithPagination, ListWithFilter, Update, IncrementViewCount, Delete

### Service层测试

- ✅ UserService: Register, Register_DuplicateUsername, Register_DuplicateEmail, Login, Login_InvalidCredentials, GetByID, GetByID_NotFound
- ✅ ArticleService: Create, GetByID, GetByID_NotFound, List, ListWithPagination, Update, Update_Unauthorized, Delete, Delete_Unauthorized

### Handler层测试

- ✅ AuthHandler: Register, Register_InvalidData, Login, Login_InvalidCredentials
- ✅ ArticleHandler: GetArticleList, GetArticleDetail, CreateArticle, CreateArticle_Unauthorized

### Middleware测试

- ✅ AuthMiddleware: ValidToken, NoToken, InvalidToken, InvalidFormat

### 工具函数测试

- ✅ JWT: GenerateJWT, ParseJWT, ParseJWT_InvalidToken, ParseJWT_WrongSecret, JWT_Expiration

## 编写新测试

### 测试函数命名

测试函数必须以`Test`开头，后跟被测试的函数名：

```go
func TestUserService_Register(t *testing.T) {
    // 测试代码
}
```

### 测试结构

```go
func TestFunctionName(t *testing.T) {
    // 1. 设置测试环境
    db := test.SetupTestDB(t)
    defer test.TeardownTestDB(db)
    
    // 2. 准备测试数据
    user := test.CreateTestUser(db, "testuser", "test@example.com")
    
    // 3. 执行被测试的函数
    result, err := service.DoSomething(user.ID)
    
    // 4. 验证结果
    if err != nil {
        t.Fatalf("执行失败: %v", err)
    }
    
    if result.ExpectedField != "expected_value" {
        t.Errorf("期望值 expected_value, 得到 %s", result.ExpectedField)
    }
}
```

### 使用断言库

项目使用`testify/assert`进行断言：

```go
import "github.com/stretchr/testify/assert"

func TestExample(t *testing.T) {
    result := DoSomething()
    assert.Equal(t, "expected", result)
    assert.NoError(t, err)
    assert.NotNil(t, obj)
}
```

## 测试最佳实践

1. **独立性**: 每个测试应该独立，不依赖其他测试的执行顺序
2. **可重复性**: 测试应该可以重复运行，结果一致
3. **快速执行**: 使用内存数据库，避免IO操作
4. **清晰命名**: 测试函数名应该清楚说明测试的内容
5. **完整覆盖**: 测试正常流程和异常流程
6. **清理资源**: 使用defer确保测试后清理资源

## 持续集成

在CI/CD流程中运行测试：

```yaml
# .github/workflows/test.yml
- name: Run tests
  run: go test ./... -v -coverprofile=coverage.out

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.out
```

## 依赖

测试需要以下依赖（已在go.mod中）：

- `github.com/stretchr/testify` - 断言库
- `gorm.io/driver/sqlite` - SQLite驱动（用于测试）

安装依赖：

```bash
go mod tidy
```

