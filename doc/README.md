# 百科Web应用 - 文档索引

## 文档概述

本文档目录包含了百科Web应用的完整设计和部署文档。项目参考萌娘百科（https://zh.moegirl.org.cn/）设计，使用Vue 3 + Go技术栈构建。

## 文档列表

### 1. [项目总体设计](./00-项目总体设计.md)
- 项目概述和背景
- 核心功能说明
- 技术选型
- 系统架构概览
- 核心业务流程
- 安全设计
- 性能优化策略
- 项目结构

### 2. [数据库设计](./01-数据库设计.md)
- 数据库选型（PostgreSQL + Redis）
- 核心数据表设计（14张表）
  - 用户表
  - 文章表
  - 文章版本表
  - 分类表
  - 标签表
  - 评论表
  - 点赞表
  - 文件表
  - 用户关注表
  - 通知表
  - 审核日志表
  - 搜索历史表
- 数据库关系图
- 索引和约束设计
- 性能优化策略
- 扩展性考虑

### 3. [系统架构设计](./02-系统架构设计.md)
- 架构模式和原则
- 系统分层架构
- 前端架构设计
- 后端架构设计
- 数据流设计
- 缓存策略
- 安全架构
- 性能优化架构
- 扩展性设计
- 监控和日志
- 部署架构

### 4. [API接口设计](./03-API接口设计.md)
- API设计规范
- 认证接口（注册、登录、刷新Token等）
- 用户接口（用户信息、关注等）
- 文章接口（CRUD、版本控制、点赞等）
- 分类接口
- 标签接口
- 评论接口
- 文件接口
- 搜索接口
- 通知接口
- 管理接口
- 错误码定义
- 接口限流规则

### 5. [前端架构设计](./04-前端架构设计.md)
- 技术栈（Vue 3 + TypeScript + Vite）
- 项目结构
- 核心模块设计
  - API模块
  - 状态管理（Pinia）
  - 路由设计
  - 组合式函数
  - 工具函数
- 组件设计示例
- 样式设计
- 性能优化
- 构建配置
- 开发规范

### 6. [后端架构设计](./05-后端架构设计.md)
- 技术栈（Go + Gin + GORM）
- 项目结构
- 核心模块设计
  - 配置管理
  - 数据模型
  - Repository层
  - Service层
  - Handler层
  - 中间件
  - 应用入口
- 错误处理
- 数据库迁移
- 测试策略

### 7. [容器化部署指南](./06-容器化部署指南.md)
- Docker和Docker Compose配置
- Dockerfile编写（前端、后端）
- docker-compose.yml配置
- Nginx配置
- 环境变量配置
- 部署步骤（开发环境、生产环境）
- 常用操作（服务管理、日志查看、数据备份）
- 监控和维护
- 故障排查
- 安全建议
- 扩展部署

## 快速开始

### 1. 阅读顺序建议
1. **项目总体设计** - 了解项目整体情况
2. **数据库设计** - 理解数据结构和关系
3. **系统架构设计** - 掌握系统架构
4. **API接口设计** - 了解接口规范
5. **前端/后端架构设计** - 深入技术实现
6. **容器化部署指南** - 学习部署方法

### 2. 开发准备
1. 阅读项目总体设计，了解项目需求
2. 根据数据库设计创建数据库表结构
3. 参考API接口设计实现后端接口
4. 参考前端架构设计实现前端页面
5. 使用容器化部署指南进行部署

### 3. 技术栈总结

#### 前端
- Vue 3.3+ (Composition API)
- TypeScript
- Vite 5+
- Element Plus
- Pinia
- Vue Router 4+
- Axios

#### 后端
- Go 1.21+
- Gin 1.9+
- GORM 1.25+
- PostgreSQL 15+
- Redis 7+
- JWT认证

#### 部署
- Docker
- Docker Compose
- Nginx

## 核心功能

### 用户系统
- ✅ 用户注册、登录
- ✅ JWT认证
- ✅ 个人资料管理
- ✅ 用户关注
- ✅ 角色权限管理

### 文章系统
- ✅ 文章创建、编辑、删除
- ✅ 版本控制（历史版本、回滚）
- ✅ Markdown编辑器
- ✅ 文章分类和标签
- ✅ 文章点赞、浏览统计

### 评论系统
- ✅ 多级评论（回复）
- ✅ 评论点赞
- ✅ 评论通知

### 内容管理
- ✅ 分类管理（多级分类）
- ✅ 标签管理
- ✅ 内容审核
- ✅ 审核日志

### 其他功能
- ✅ 文件上传
- ✅ 全文搜索
- ✅ 站内通知
- ✅ 搜索历史

## 数据库表结构

项目包含15张核心数据表：
1. users - 用户表
2. articles - 文章表
3. article_versions - 文章版本表
4. categories - 分类表
5. article_categories - 文章分类关联表
6. tags - 标签表
7. article_tags - 文章标签关联表
8. comments - 评论表
9. likes - 点赞表
10. files - 文件表
11. article_images - 文章图片表
12. user_follows - 用户关注表
13. notifications - 通知表
14. moderation_logs - 审核日志表
15. search_history - 搜索历史表

详细设计请参考[数据库设计文档](./01-数据库设计.md)。

## API接口

项目提供完整的RESTful API接口，包括：
- 认证接口（/api/v1/auth/*）
- 用户接口（/api/v1/users/*）
- 文章接口（/api/v1/articles/*）
- 分类接口（/api/v1/categories/*）
- 标签接口（/api/v1/tags/*）
- 评论接口（/api/v1/comments/*）
- 文件接口（/api/v1/files/*）
- 搜索接口（/api/v1/search/*）
- 通知接口（/api/v1/notifications/*）
- 管理接口（/api/v1/admin/*）

详细接口文档请参考[API接口设计文档](./03-API接口设计.md)。

## 部署方式

### 开发环境
使用Docker Compose快速启动开发环境：
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### 生产环境
使用Docker Compose部署生产环境：
```bash
docker-compose up -d
```

详细部署步骤请参考[容器化部署指南](./06-容器化部署指南.md)。

## 项目结构

```
dbapp/
├── doc/                    # 文档目录（本目录）
│   ├── 00-项目总体设计.md
│   ├── 01-数据库设计.md
│   ├── 02-系统架构设计.md
│   ├── 03-API接口设计.md
│   ├── 04-前端架构设计.md
│   ├── 05-后端架构设计.md
│   ├── 06-容器化部署指南.md
│   └── README.md
├── frontend/               # 前端项目（待实现）
├── backend/                # 后端项目（待实现）
├── docker/                 # Docker相关文件（待实现）
├── nginx/                  # Nginx配置（待实现）
└── README.md               # 项目主README
```

## 开发计划

### 阶段一：基础框架搭建 ✅
- [x] 项目设计文档
- [x] 数据库设计
- [ ] 项目初始化
- [ ] 基础API框架
- [ ] 前端项目搭建

### 阶段二：核心功能开发
- [ ] 用户系统
- [ ] 文章系统
- [ ] 分类标签系统
- [ ] 评论系统

### 阶段三：高级功能开发
- [ ] 搜索功能
- [ ] 通知系统
- [ ] 内容审核
- [ ] 文件管理

### 阶段四：优化和部署
- [ ] 性能优化
- [ ] 安全加固
- [ ] 容器化部署
- [ ] 文档完善

## 参考资料

- 萌娘百科: https://zh.moegirl.org.cn/
- Vue 3 文档: https://cn.vuejs.org/
- Go 官方文档: https://go.dev/doc/
- PostgreSQL 文档: https://www.postgresql.org/docs/
- Docker 文档: https://docs.docker.com/
- Gin 文档: https://gin-gonic.com/docs/
- GORM 文档: https://gorm.io/docs/

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证。

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。

---

**注意**: 本文档目录包含了项目的完整设计文档。在实际开发前，请仔细阅读相关文档，确保理解系统设计和实现要求。

