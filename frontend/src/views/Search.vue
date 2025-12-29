<template>
  <div class="search-page">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1 @click="$router.push('/')" style="cursor: pointer;">百科Web应用</h1>
          <div class="header-search">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索文章..."
              size="large"
              clearable
              @keyup.enter="handleSearch"
              style="width: 400px;"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
              <template #append>
                <el-button @click="handleSearch">搜索</el-button>
              </template>
            </el-input>
          </div>
          <div class="header-actions">
            <el-button v-if="!userStore.isLoggedIn" @click="$router.push('/login')">登录</el-button>
            <el-button v-if="userStore.isLoggedIn" @click="$router.push('/')">首页</el-button>
          </div>
        </div>
      </el-header>

      <el-main>
        <el-row :gutter="20">
          <!-- 左侧筛选 -->
          <el-col :span="5">
            <el-card class="filter-card">
              <template #header>
                <div class="filter-header">
                  <span>筛选条件</span>
                  <el-button type="primary" link size="small" @click="clearFilters">清除</el-button>
                </div>
              </template>

              <!-- 分类筛选 -->
              <div class="filter-section">
                <h4>分类</h4>
                <el-radio-group v-model="filters.category_id" @change="handleSearch">
                  <el-radio :label="0">全部</el-radio>
                  <el-radio v-for="cat in categories" :key="cat.id" :label="cat.id">
                    {{ cat.name }} ({{ cat.article_count || 0 }})
                  </el-radio>
                </el-radio-group>
              </div>

              <!-- 标签筛选 -->
              <div class="filter-section">
                <h4>热门标签</h4>
                <div class="tag-list">
                  <el-tag
                    v-for="tag in popularTags"
                    :key="tag.id"
                    :class="{ 'tag-selected': filters.tag_id === tag.id }"
                    :style="getTagStyle(tag)"
                    size="small"
                    @click="toggleTag(tag.id)"
                    style="cursor: pointer; margin: 3px;"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
              </div>

              <!-- 排序 -->
              <div class="filter-section">
                <h4>排序方式</h4>
                <el-radio-group v-model="filters.sort" @change="handleSearch">
                  <el-radio label="created_at">最新发布</el-radio>
                  <el-radio label="view_count">最多浏览</el-radio>
                  <el-radio label="like_count">最多点赞</el-radio>
                  <el-radio label="comment_count">最多评论</el-radio>
                </el-radio-group>
              </div>
            </el-card>
          </el-col>

          <!-- 右侧结果 -->
          <el-col :span="19">
            <el-card>
              <template #header>
                <div class="result-header">
                  <span v-if="searchKeyword">
                    搜索 "<strong>{{ searchKeyword }}</strong>" 的结果
                  </span>
                  <span v-else>全部文章</span>
                  <span class="result-count">共 {{ pagination.total }} 条结果</span>
                </div>
              </template>

              <el-skeleton v-if="loading" :rows="10" animated />

              <el-empty v-else-if="articles.length === 0" description="没有找到相关文章">
                <el-button type="primary" @click="clearFilters">清除筛选条件</el-button>
              </el-empty>

              <div v-else class="article-list">
                <div v-for="article in articles" :key="article.id" class="article-item">
                  <div class="article-cover" v-if="article.cover_image_url">
                    <img 
                      :src="article.cover_image_url" 
                      :alt="article.title"
                      @click="$router.push(`/article/${article.id}`)"
                    />
                  </div>
                  <div class="article-info">
                    <h3 class="article-title">
                      <el-link 
                        type="primary" 
                        @click="$router.push(`/article/${article.id}`)"
                      >
                        <span v-html="highlightKeyword(article.title)"></span>
                      </el-link>
                    </h3>
                    <p class="article-summary">
                      <span v-html="highlightKeyword(article.summary || '暂无摘要')"></span>
                    </p>
                    <div class="article-tags" v-if="article.categories?.length || article.tags?.length">
                      <el-tag 
                        v-for="cat in article.categories" 
                        :key="'cat-' + cat.id" 
                        size="small"
                        @click="filterByCategory(cat.id)"
                        style="cursor: pointer; margin-right: 5px;"
                      >
                        {{ cat.name }}
                      </el-tag>
                      <el-tag 
                        v-for="tag in article.tags" 
                        :key="'tag-' + tag.id"
                        size="small"
                        :style="tag.color ? { backgroundColor: tag.color, borderColor: tag.color, color: '#fff' } : {}"
                        :type="tag.color ? undefined : 'info'"
                        @click="filterByTag(tag.id)"
                        style="cursor: pointer; margin-right: 5px;"
                      >
                        {{ tag.name }}
                      </el-tag>
                    </div>
                    <div class="article-meta">
                      <span>{{ article.author?.nickname || article.author?.username }}</span>
                      <span>{{ formatDate(article.created_at) }}</span>
                      <span>👁 {{ article.view_count || 0 }}</span>
                      <span>👍 {{ article.like_count || 0 }}</span>
                      <span>💬 {{ article.comment_count || 0 }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <el-pagination
                v-if="pagination.total > 0"
                v-model:current-page="pagination.page"
                v-model:page-size="pagination.page_size"
                :total="pagination.total"
                :page-sizes="[10, 20, 50]"
                layout="total, sizes, prev, pager, next"
                @size-change="handleSearch"
                @current-change="handleSearch"
                style="margin-top: 20px; justify-content: center;"
              />
            </el-card>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { getArticleList } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { Article, Category, Tag, Pagination } from '@/api/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const searchKeyword = ref('')
const loading = ref(false)
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const popularTags = ref<Tag[]>([])

const pagination = ref<Pagination>({
  page: 1,
  page_size: 10,
  total: 0,
  total_pages: 0
})

const filters = ref({
  category_id: 0,
  tag_id: 0,
  sort: 'created_at',
  order: 'desc'
})

// 从 URL 读取参数
function loadFromQuery() {
  const query = route.query
  searchKeyword.value = (query.q as string) || ''
  filters.value.category_id = parseInt(query.category as string) || 0
  filters.value.tag_id = parseInt(query.tag as string) || 0
  filters.value.sort = (query.sort as string) || 'created_at'
  pagination.value.page = parseInt(query.page as string) || 1
}

// 更新 URL
function updateQuery() {
  const query: Record<string, string> = {}
  if (searchKeyword.value) query.q = searchKeyword.value
  if (filters.value.category_id) query.category = String(filters.value.category_id)
  if (filters.value.tag_id) query.tag = String(filters.value.tag_id)
  if (filters.value.sort !== 'created_at') query.sort = filters.value.sort
  if (pagination.value.page > 1) query.page = String(pagination.value.page)
  
  router.replace({ query })
}

async function handleSearch() {
  loading.value = true
  updateQuery()
  
  try {
    const params: any = {
      page: pagination.value.page,
      page_size: pagination.value.page_size,
      sort: filters.value.sort,
      order: filters.value.order,
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    if (filters.value.category_id) {
      params.category_id = filters.value.category_id
    }
    if (filters.value.tag_id) {
      params.tag_id = filters.value.tag_id
    }

    const result = await getArticleList(params)
    articles.value = result.items || []
    pagination.value = result.pagination || pagination.value
  } catch (error: any) {
    ElMessage.error(error.message || '搜索失败')
    articles.value = []
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    categories.value = await getCategoryList({})
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

async function loadTags() {
  try {
    const tags = await getTagList({ sort: 'article_count', order: 'desc', limit: 20 })
    popularTags.value = tags
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

function clearFilters() {
  searchKeyword.value = ''
  filters.value = {
    category_id: 0,
    tag_id: 0,
    sort: 'created_at',
    order: 'desc'
  }
  pagination.value.page = 1
  handleSearch()
}

function toggleTag(tagId: number) {
  if (filters.value.tag_id === tagId) {
    filters.value.tag_id = 0
  } else {
    filters.value.tag_id = tagId
  }
  pagination.value.page = 1
  handleSearch()
}

function filterByCategory(categoryId: number) {
  filters.value.category_id = categoryId
  pagination.value.page = 1
  handleSearch()
}

function filterByTag(tagId: number) {
  filters.value.tag_id = tagId
  pagination.value.page = 1
  handleSearch()
}

function getTagStyle(tag: Tag) {
  const isSelected = filters.value.tag_id === tag.id
  if (tag.color) {
    return {
      backgroundColor: isSelected ? tag.color : 'transparent',
      borderColor: tag.color,
      color: isSelected ? '#fff' : tag.color,
    }
  }
  return {}
}

function highlightKeyword(text: string): string {
  if (!searchKeyword.value || !text) return text
  const regex = new RegExp(`(${searchKeyword.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 监听路由变化
watch(() => route.query, () => {
  if (route.name === 'Search') {
    loadFromQuery()
    handleSearch()
  }
}, { immediate: false })

onMounted(async () => {
  loadFromQuery()
  await Promise.all([loadCategories(), loadTags()])
  handleSearch()
})
</script>

<style scoped>
.search-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.el-header {
  background-color: #409eff;
  color: white;
  height: 70px;
  padding: 0 20px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.header-content h1 {
  margin: 0;
  font-size: 24px;
  white-space: nowrap;
}

.header-search {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 40px;
}

.header-actions {
  white-space: nowrap;
}

.el-main {
  padding: 20px;
}

.filter-card {
  position: sticky;
  top: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-section {
  margin-bottom: 20px;
}

.filter-section h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #666;
}

.filter-section :deep(.el-radio-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-section :deep(.el-radio) {
  margin-right: 0;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
}

.tag-selected {
  font-weight: bold;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-count {
  color: #999;
  font-size: 14px;
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-item {
  display: flex;
  gap: 15px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.article-item:last-child {
  border-bottom: none;
}

.article-cover {
  flex-shrink: 0;
  width: 180px;
  height: 120px;
  overflow: hidden;
  border-radius: 8px;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s;
}

.article-cover img:hover {
  transform: scale(1.05);
}

.article-info {
  flex: 1;
  min-width: 0;
}

.article-title {
  margin: 0 0 10px 0;
  font-size: 18px;
}

.article-title :deep(mark) {
  background-color: #fff3cd;
  padding: 0 2px;
}

.article-summary {
  color: #666;
  margin: 0 0 10px 0;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-summary :deep(mark) {
  background-color: #fff3cd;
  padding: 0 2px;
}

.article-tags {
  margin-bottom: 10px;
}

.article-meta {
  color: #999;
  font-size: 13px;
  display: flex;
  gap: 15px;
}
</style>
