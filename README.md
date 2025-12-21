# 百科Web应用

一个类似萌娘百科的在线百科平台，支持用户创建、编辑、管理百科条目，并提供完善的社区功能。

## 技术栈

### 前端
- Vue 3.3+ (Composition API)
- TypeScript
- Vite 5+
- Element Plus
- Pinia
- Vue Router 4+

### 后端
- Go 1.21+
- Gin 1.9+
- GORM 1.25+
- PostgreSQL 15+
- Redis 7+
- JWT认证

## 项目结构

```
dbapp/
├── doc/                    # 文档目录
├── frontend/               # 前端项目
├── backend/                # 后端项目
├── docker/                 # Docker相关文件
├── nginx/                  # Nginx配置
└── README.md
```

## 快速开始

### 使用Docker Compose（推荐）

1. 克隆项目
```bash
git clone <repository-url>
cd dbapp
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑.env文件，修改配置
```

3. 启动服务
```bash
cd docker
docker-compose up -d
```

4. 访问应用
- 前端: http://localhost
- 后端API: http://localhost:8080

### 本地开发

#### 后端开发

1. 安装依赖
```bash
cd backend
go mod download
```

2. 配置数据库
- 创建PostgreSQL数据库
- 修改 `config/config.yaml` 中的数据库配置

3. 运行服务
```bash
go run cmd/api/main.go
```

#### 前端开发

1. 安装依赖
```bash
cd frontend
npm install
```

2. 运行开发服务器
```bash
npm run dev
```

## 功能特性

- ✅ 用户注册、登录、JWT认证
- ✅ 文章创建、编辑、删除
- ✅ 文章分类和标签
- ✅ 评论系统
- ✅ 点赞功能
- ✅ 全文搜索
- ✅ 文件上传

## API文档

API接口文档请参考 [doc/03-API接口设计.md](./doc/03-API接口设计.md)

## 文档

完整的设计文档位于 `doc/` 目录：
- [项目总体设计](./doc/00-项目总体设计.md)
- [数据库设计](./doc/01-数据库设计.md)
- [系统架构设计](./doc/02-系统架构设计.md)
- [API接口设计](./doc/03-API接口设计.md)
- [前端架构设计](./doc/04-前端架构设计.md)
- [后端架构设计](./doc/05-后端架构设计.md)
- [容器化部署指南](./doc/06-容器化部署指南.md)

## 开发计划

- [x] 项目设计和文档
- [x] 后端基础框架
- [ ] 前端项目实现
- [ ] 完整功能实现
- [ ] 测试和优化

## 许可证

MIT License

