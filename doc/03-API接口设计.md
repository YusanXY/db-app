# 百科Web应用 - API接口设计文档

## 1. API设计规范

### 1.1 基础规范
- **协议**: HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8
- **时间格式**: ISO 8601 (如: 2024-01-01T00:00:00Z)

### 1.2 URL规范
- 使用RESTful风格
- 使用复数名词: `/api/v1/articles`
- 使用小写字母和连字符: `/api/v1/article-categories`
- 版本号: `/api/v1/`

### 1.3 HTTP方法
- **GET**: 查询资源
- **POST**: 创建资源
- **PUT**: 完整更新资源
- **PATCH**: 部分更新资源
- **DELETE**: 删除资源

### 1.4 状态码
- **200**: 成功
- **201**: 创建成功
- **400**: 请求参数错误
- **401**: 未认证
- **403**: 无权限
- **404**: 资源不存在
- **422**: 验证失败
- **429**: 请求过于频繁
- **500**: 服务器错误

### 1.5 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "errors": [
    {
      "field": "username",
      "message": "用户名已存在"
    }
  ]
}
```

### 1.6 分页格式
```json
{
  "code": 200,
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

## 2. 认证接口

### 2.1 用户注册
**POST** `/api/v1/auth/register`

**请求体**:
```json
{
  "username": "string (必填, 3-50字符)",
  "email": "string (必填, 邮箱格式)",
  "password": "string (必填, 8-50字符)",
  "nickname": "string (可选)"
}
```

**响应**:
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户"
    }
  }
}
```

### 2.2 用户登录
**POST** `/api/v1/auth/login`

**请求体**:
```json
{
  "username": "string (必填)",
  "password": "string (必填)"
}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "token": "jwt_token_string",
    "refresh_token": "refresh_token_string",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户",
      "avatar_url": "https://...",
      "role": "user"
    }
  }
}
```

### 2.3 刷新Token
**POST** `/api/v1/auth/refresh`

**请求头**:
```
Authorization: Bearer {refresh_token}
```

**响应**: 同登录接口

### 2.4 用户登出
**POST** `/api/v1/auth/logout`

**请求头**:
```
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 200,
  "message": "登出成功"
}
```

### 2.5 获取当前用户信息
**GET** `/api/v1/auth/me`

**请求头**:
```
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar_url": "https://...",
    "bio": "个人简介",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## 3. 用户接口

### 3.1 获取用户列表
**GET** `/api/v1/users`

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20, 最大: 100)
- `role`: 角色筛选
- `status`: 状态筛选
- `keyword`: 关键词搜索

**响应**: 分页用户列表

### 3.2 获取用户详情
**GET** `/api/v1/users/:id`

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar_url": "https://...",
    "bio": "个人简介",
    "role": "user",
    "article_count": 10,
    "follower_count": 5,
    "following_count": 3,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3.3 更新用户信息
**PATCH** `/api/v1/users/:id`

**请求头**:
```
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "nickname": "string (可选)",
  "avatar_url": "string (可选)",
  "bio": "string (可选)"
}
```

**响应**: 更新后的用户信息

### 3.4 修改密码
**POST** `/api/v1/users/:id/password`

**请求体**:
```json
{
  "old_password": "string (必填)",
  "new_password": "string (必填, 8-50字符)"
}
```

### 3.5 关注用户
**POST** `/api/v1/users/:id/follow`

**响应**:
```json
{
  "code": 200,
  "message": "关注成功"
}
```

### 3.6 取消关注
**DELETE** `/api/v1/users/:id/follow`

### 3.7 获取用户文章列表
**GET** `/api/v1/users/:id/articles`

**查询参数**: 同文章列表接口

## 4. 文章接口

### 4.1 获取文章列表
**GET** `/api/v1/articles`

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `category_id`: 分类ID
- `tag_id`: 标签ID
- `author_id`: 作者ID
- `status`: 状态 (published/draft/archived)
- `keyword`: 搜索关键词
- `sort`: 排序 (created_at/view_count/like_count)
- `order`: 排序方向 (asc/desc)

**响应**:
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "文章标题",
        "slug": "article-slug",
        "summary": "文章摘要",
        "cover_image_url": "https://...",
        "author": {
          "id": 1,
          "username": "testuser",
          "nickname": "测试用户",
          "avatar_url": "https://..."
        },
        "categories": [
          {"id": 1, "name": "分类1", "slug": "category-1"}
        ],
        "tags": [
          {"id": 1, "name": "标签1", "slug": "tag-1"}
        ],
        "view_count": 100,
        "like_count": 10,
        "comment_count": 5,
        "is_featured": false,
        "published_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### 4.2 获取文章详情
**GET** `/api/v1/articles/:id`

**查询参数**:
- `version`: 版本号 (可选，获取历史版本)

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "文章标题",
    "slug": "article-slug",
    "content": "Markdown内容",
    "content_html": "<p>HTML内容</p>",
    "summary": "文章摘要",
    "cover_image_url": "https://...",
    "author": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户"
    },
    "editor": {
      "id": 2,
      "username": "editor",
      "nickname": "编辑者"
    },
    "categories": [],
    "tags": [],
    "view_count": 100,
    "like_count": 10,
    "comment_count": 5,
    "edit_count": 3,
    "is_featured": false,
    "is_locked": false,
    "status": "published",
    "published_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "is_liked": false  // 当前用户是否点赞
  }
}
```

### 4.3 创建文章
**POST** `/api/v1/articles`

**请求头**:
```
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "title": "string (必填, 1-500字符)",
  "content": "string (必填, Markdown格式)",
  "summary": "string (可选)",
  "cover_image_url": "string (可选)",
  "category_ids": [1, 2],
  "tag_ids": [1, 2],
  "status": "draft|published (默认: draft)"
}
```

**响应**: 创建的文章信息

### 4.4 更新文章
**PUT** `/api/v1/articles/:id`

**请求体**: 同创建文章

**响应**: 更新后的文章信息

### 4.5 删除文章
**DELETE** `/api/v1/articles/:id`

**响应**:
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 4.6 点赞/取消点赞文章
**POST** `/api/v1/articles/:id/like`
**DELETE** `/api/v1/articles/:id/like`

### 4.7 获取文章版本列表
**GET** `/api/v1/articles/:id/versions`

**响应**: 版本列表

### 4.8 获取文章版本详情
**GET** `/api/v1/articles/:id/versions/:version`

**响应**: 版本详情

### 4.9 回滚到指定版本
**POST** `/api/v1/articles/:id/versions/:version/rollback`

**请求体**:
```json
{
  "reason": "string (可选, 回滚原因)"
}
```

## 5. 分类接口

### 5.1 获取分类列表
**GET** `/api/v1/categories`

**查询参数**:
- `parent_id`: 父分类ID (获取子分类)
- `is_active`: 是否启用

**响应**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "分类名称",
      "slug": "category-slug",
      "description": "分类描述",
      "parent_id": null,
      "icon_url": "https://...",
      "article_count": 10,
      "children": [
        // 子分类
      ]
    }
  ]
}
```

### 5.2 获取分类详情
**GET** `/api/v1/categories/:id`

### 5.3 创建分类
**POST** `/api/v1/categories`

**请求体**:
```json
{
  "name": "string (必填)",
  "slug": "string (必填, 唯一)",
  "description": "string (可选)",
  "parent_id": "integer (可选)",
  "icon_url": "string (可选)",
  "sort_order": "integer (可选, 默认: 0)"
}
```

### 5.4 更新分类
**PUT** `/api/v1/categories/:id`

### 5.5 删除分类
**DELETE** `/api/v1/categories/:id`

### 5.6 获取分类下的文章
**GET** `/api/v1/categories/:id/articles`

## 6. 标签接口

### 6.1 获取标签列表
**GET** `/api/v1/tags`

**查询参数**:
- `keyword`: 关键词搜索
- `sort`: 排序 (name/article_count)

**响应**: 标签列表

### 6.2 获取标签详情
**GET** `/api/v1/tags/:id`

### 6.3 创建标签
**POST** `/api/v1/tags`

**请求体**:
```json
{
  "name": "string (必填, 唯一)",
  "slug": "string (必填, 唯一)",
  "description": "string (可选)",
  "color": "string (可选, CSS颜色值)"
}
```

### 6.4 获取标签下的文章
**GET** `/api/v1/tags/:id/articles`

## 7. 评论接口

### 7.1 获取评论列表
**GET** `/api/v1/articles/:article_id/comments`

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `parent_id`: 父评论ID (获取回复)

**响应**:
```json
{
  "code": 200,
  "data": {
    "items": [
      {
        "id": 1,
        "content": "评论内容",
        "content_html": "<p>评论内容</p>",
        "user": {
          "id": 1,
          "username": "testuser",
          "nickname": "测试用户",
          "avatar_url": "https://..."
        },
        "parent_id": null,
        "like_count": 5,
        "reply_count": 2,
        "is_liked": false,
        "replies": [
          // 回复列表
        ],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {}
  }
}
```

### 7.2 创建评论
**POST** `/api/v1/articles/:article_id/comments`

**请求头**:
```
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "content": "string (必填, 1-5000字符)",
  "parent_id": "integer (可选, 回复的评论ID)"
}
```

**响应**: 创建的评论信息

### 7.3 更新评论
**PUT** `/api/v1/comments/:id`

**请求体**:
```json
{
  "content": "string (必填)"
}
```

### 7.4 删除评论
**DELETE** `/api/v1/comments/:id`

### 7.5 点赞/取消点赞评论
**POST** `/api/v1/comments/:id/like`
**DELETE** `/api/v1/comments/:id/like`

## 8. 文件接口

### 8.1 上传文件
**POST** `/api/v1/files/upload`

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:
- `file`: 文件 (必填，FormData格式)

**文件限制**:
- 支持的文件类型：jpg, jpeg, png, gif, webp（可在配置中修改）
- 最大文件大小：10MB（可在配置中修改）

**响应**:
```json
{
  "code": 200,
  "data": {
    "url": "/uploads/1734789123456_filename.jpg",
    "name": "original.jpg",
    "size": 102400
  }
}
```

**错误响应**:
- `400`: 文件大小超限、文件类型不支持、未选择文件
- `401`: 未认证
- `500`: 服务器错误（创建目录失败、保存文件失败等）

**说明**:
- 文件保存在服务器本地文件系统（`./uploads`目录）
- 文件名格式：`{timestamp}_{原始文件名}`
- 返回的URL为相对路径，前端通过Nginx代理访问
- 上传的文件会自动通过静态文件服务提供访问（`/uploads/*`）

### 8.2 获取文件列表
**GET** `/api/v1/files`

**查询参数**:
- `file_type`: 文件类型
- `uploader_id`: 上传者ID
- `page`: 页码

### 8.3 删除文件
**DELETE** `/api/v1/files/:id`

## 9. 搜索接口

### 9.1 搜索文章
**GET** `/api/v1/search/articles`

**查询参数**:
- `q`: 搜索关键词 (必填)
- `page`: 页码
- `page_size`: 每页数量
- `category_id`: 分类筛选
- `tag_id`: 标签筛选

**响应**: 搜索结果列表

### 9.2 搜索建议
**GET** `/api/v1/search/suggestions`

**查询参数**:
- `q`: 搜索关键词

**响应**:
```json
{
  "code": 200,
  "data": {
    "keywords": ["关键词1", "关键词2"],
    "articles": [
      {"id": 1, "title": "文章标题"}
    ]
  }
}
```

### 9.3 热门搜索
**GET** `/api/v1/search/hot`

**响应**: 热门搜索关键词列表

## 10. 通知接口

### 10.1 获取通知列表
**GET** `/api/v1/notifications`

**请求头**:
```
Authorization: Bearer {token}
```

**查询参数**:
- `page`: 页码
- `type`: 通知类型
- `is_read`: 是否已读

**响应**: 通知列表

### 10.2 标记通知为已读
**PATCH** `/api/v1/notifications/:id/read`

### 10.3 标记所有通知为已读
**POST** `/api/v1/notifications/read-all`

### 10.4 获取未读通知数量
**GET** `/api/v1/notifications/unread-count`

**响应**:
```json
{
  "code": 200,
  "data": {
    "count": 5
  }
}
```

## 11. 管理接口

### 11.1 审核内容
**POST** `/api/v1/admin/moderate`

**请求体**:
```json
{
  "target_type": "article|comment|user",
  "target_id": 1,
  "action": "approve|reject|delete|ban|unban",
  "reason": "string (可选)"
}
```

### 11.2 获取审核日志
**GET** `/api/v1/admin/moderation-logs`

**查询参数**:
- `target_type`: 目标类型
- `action`: 操作类型
- `moderator_id`: 审核员ID
- `page`: 页码

### 11.3 获取统计数据
**GET** `/api/v1/admin/statistics`

**响应**:
```json
{
  "code": 200,
  "data": {
    "user_count": 1000,
    "article_count": 500,
    "comment_count": 2000,
    "today_views": 10000
  }
}
```

## 12. 错误码定义

| 错误码 | 说明 |
|--------|------|
| 1000 | 参数错误 |
| 1001 | 用户名已存在 |
| 1002 | 邮箱已存在 |
| 1003 | 用户名或密码错误 |
| 1004 | Token无效或过期 |
| 1005 | 无权限 |
| 2001 | 文章不存在 |
| 2002 | 文章已锁定 |
| 2003 | 文章状态不允许此操作 |
| 3001 | 评论不存在 |
| 3002 | 不能回复自己的评论 |
| 4001 | 文件类型不支持 |
| 4002 | 文件大小超限 |
| 5001 | 分类不存在 |
| 5002 | 标签不存在 |

## 13. 接口限流

### 13.1 限流规则
- 普通接口: 100次/分钟
- 登录接口: 5次/分钟
- 上传接口: 10次/分钟
- 搜索接口: 60次/分钟

### 13.2 限流响应
```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "retry_after": 60
}
```

## 14. WebSocket接口（可选）

### 14.1 实时通知
**WS** `/ws/notifications`

**连接认证**: 通过Token认证

**消息格式**:
```json
{
  "type": "notification",
  "data": {
    "id": 1,
    "type": "comment",
    "title": "新评论",
    "content": "有人评论了你的文章"
  }
}
```

