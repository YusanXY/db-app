# Markdown编辑器和文章编辑功能实现文档

## 1. 功能概述

根据项目总体设计文档3.2节（文章创建编辑流程），实现了以下功能：
1. **Markdown编辑器**：使用vditor实现所见即所得的Markdown编辑
2. **文章编辑功能**：用户可以编辑自己创建的文章

## 2. 技术实现

### 2.1 Markdown编辑器组件

**文件位置**：`frontend/src/components/MarkdownEditor.vue`

**功能特性**：
- 使用vditor 3.11.2作为Markdown编辑器
- 支持所见即所得（WYSIWYG）模式
- 支持代码高亮
- 支持图片上传（配置了上传接口）
- 支持全屏编辑
- 支持预览模式
- 支持大纲视图

**主要配置**：
```typescript
{
  height: 500, // 编辑器高度
  mode: 'sv', // 所见即所得模式
  toolbar: [...], // 工具栏配置
  upload: {...} // 图片上传配置
}
```

### 2.2 Markdown渲染组件

**文件位置**：`frontend/src/components/MarkdownViewer.vue`

**功能特性**：
- 使用marked库解析Markdown
- 使用highlight.js实现代码高亮
- 支持GitHub风格的Markdown（GFM）
- 自定义样式，美观易读

**支持的Markdown语法**：
- 标题（H1-H6）
- 段落
- 列表（有序、无序、任务列表）
- 引用
- 代码块和行内代码
- 表格
- 链接
- 图片
- 水平线

### 2.3 文章创建页面

**文件位置**：`frontend/src/views/Article/Create.vue`

**更新内容**：
- 将普通textarea替换为MarkdownEditor组件
- 添加状态选择（草稿/发布）
- 添加"保存草稿"功能
- 优化表单布局

### 2.4 文章编辑页面

**文件位置**：`frontend/src/views/Article/Edit.vue`

**功能特性**：
- 加载现有文章内容
- 权限检查（只能编辑自己的文章）
- 使用MarkdownEditor编辑内容
- 支持保存草稿和发布
- 编辑后跳转回文章详情页

### 2.5 文章详情页更新

**文件位置**：`frontend/src/views/Article/Detail.vue`

**更新内容**：
- 添加"编辑"按钮（仅文章作者可见）
- 使用MarkdownViewer渲染Markdown内容
- 优化内容显示样式

## 3. 路由配置

**文件位置**：`frontend/src/router/index.ts`

**新增路由**：
```typescript
{
  path: '/article/:id/edit',
  name: 'ArticleEdit',
  component: () => import('@/views/Article/Edit.vue'),
  meta: { requiresAuth: true }
}
```

## 4. 依赖包

### 前端新增依赖
- `vditor@^3.11.2` - Markdown编辑器
- `marked@^11.1.1` - Markdown解析（已存在）
- `highlight.js@^11.9.0` - 代码高亮（已存在）

## 5. 使用流程

### 5.1 创建文章
1. 用户登录后，点击"创建文章"
2. 填写标题
3. 使用Markdown编辑器编写内容
4. 填写摘要（可选）
5. 选择状态（草稿/发布）
6. 点击"发布"或"保存草稿"

### 5.2 编辑文章
1. 在文章详情页，点击"编辑"按钮（仅作者可见）
2. 系统加载文章内容到编辑器
3. 修改内容
4. 点击"保存"或"保存草稿"
5. 自动跳转回文章详情页

### 5.3 查看文章
1. 文章详情页自动使用MarkdownViewer渲染Markdown内容
2. 代码块自动高亮
3. 表格、列表等格式正确显示

## 6. 权限控制

### 6.1 编辑权限
- 只有文章作者可以编辑自己的文章
- 编辑页面会检查权限，无权限会跳转回详情页
- 编辑按钮仅对文章作者显示

### 6.2 访问控制
- 编辑页面需要登录（`requiresAuth: true`）
- 未登录用户访问编辑页面会跳转到登录页

## 7. API接口

### 7.1 获取文章详情
```
GET /api/v1/articles/:id
```

### 7.2 更新文章
```
PUT /api/v1/articles/:id
Body: {
  title?: string,
  content?: string,
  summary?: string,
  status?: 'draft' | 'published'
}
```

## 8. 样式说明

### 8.1 Markdown编辑器样式
- 使用vditor默认样式
- 编辑器高度：500px
- 响应式布局

### 8.2 Markdown渲染样式
- GitHub风格的Markdown样式
- 代码块使用深色主题高亮
- 表格、引用等元素有适当的间距和边框
- 链接有悬停效果

## 9. 注意事项

### 9.1 图片上传
- 图片上传接口配置为：`/api/v1/files/upload`
- 目前为占位配置，需要实现文件上传功能后生效

### 9.2 内容存储
- Markdown原始内容存储在`content`字段
- HTML渲染内容可存储在`content_html`字段（可选）
- 前端使用MarkdownViewer实时渲染

### 9.3 编辑器模式
- 当前使用所见即所得模式（`mode: 'sv'`）
- 可以切换到其他模式（如分屏模式）

## 10. 后续优化建议

1. **图片上传功能**：实现完整的图片上传和管理功能
2. **版本控制**：实现文章版本历史记录（如设计文档3.2节提到的）
3. **自动保存**：实现草稿自动保存功能
4. **预览功能**：在编辑器中添加实时预览
5. **工具栏扩展**：根据需求添加更多Markdown语法支持
6. **内容验证**：添加Markdown内容格式验证

## 11. 测试建议

1. **创建文章测试**：
   - 测试各种Markdown语法
   - 测试代码高亮
   - 测试保存草稿和发布

2. **编辑文章测试**：
   - 测试权限控制
   - 测试内容加载
   - 测试更新功能

3. **显示测试**：
   - 测试Markdown渲染效果
   - 测试不同浏览器的兼容性
   - 测试响应式布局

