<template>
  <div class="home">
    <el-container>
      <el-header>
        <h1>百科Web应用</h1>
        <div class="header-search">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章..."
            clearable
            @keyup.enter="handleSearchSubmit"
            style="width: 300px;"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearchSubmit" style="margin-left: 10px;">搜索</el-button>
        </div>
        <div class="header-actions">
          <el-button v-if="!userStore.isLoggedIn" @click="$router.push('/login')">登录</el-button>
          <el-button v-if="!userStore.isLoggedIn" @click="$router.push('/register')">注册</el-button>
          <el-button v-if="userStore.isLoggedIn" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <el-card>
          <h2>欢迎使用百科Web应用</h2>
          <p>这是一个类似萌娘百科的在线百科平台</p>
          <el-divider />
          <el-row :gutter="20" style="margin-bottom: 20px;">
            <el-col :span="8">
              <el-card shadow="hover">
                <h3>创建文章</h3>
                <p>创建和编辑百科条目</p>
                <el-button type="primary" @click="$router.push('/article/new')">创建文章</el-button>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover">
                <h3>浏览文章</h3>
                <p>查看和搜索百科内容</p>
                <el-button type="primary" @click="$router.push('/search')">搜索文章</el-button>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover">
                <h3>分类标签</h3>
                <p>管理分类和标签</p>
                <el-button type="primary" @click="$router.push('/categories')">分类管理</el-button>
                <el-button type="primary" @click="$router.push('/tags')" style="margin-left: 10px;">标签管理</el-button>
              </el-card>
            </el-col>
          </el-row>
          <el-divider />
          <h3>最新文章</h3>
          <el-empty v-if="articles.length === 0 && !loading" description="暂无文章" />
          <el-skeleton v-if="loading" :rows="5" animated />
          <el-list v-else>
            <el-list-item v-for="article in articles" :key="article.id" style="border-bottom: 1px solid #eee; padding: 15px 0;">
              <template #default>
                <div style="width: 100%; display: flex; gap: 15px;">
                  <div v-if="article.cover_image_url" style="flex-shrink: 0; width: 200px; height: 120px; overflow: hidden; border-radius: 4px;">
                    <img :src="article.cover_image_url" :alt="article.title" style="width: 100%; height: 100%; object-fit: cover; cursor: pointer;" @click="$router.push(`/article/${article.id}`)" />
                  </div>
                  <div style="flex: 1; min-width: 0;">
                    <h4 style="margin: 0 0 10px 0;">
                      <el-link :href="`/article/${article.id}`" type="primary" @click.prevent="$router.push(`/article/${article.id}`)">
                        {{ article.title }}
                      </el-link>
                      <el-tag v-if="article.status === 'draft'" size="small" type="info" style="margin-left: 10px;">草稿</el-tag>
                      <el-tag v-if="article.status === 'published'" size="small" type="success" style="margin-left: 10px;">已发布</el-tag>
                    </h4>
                    <p style="color: #666; margin: 5px 0;">{{ article.summary || '暂无摘要' }}</p>
                    <div v-if="(article.categories && article.categories.length > 0) || (article.tags && article.tags.length > 0)" style="margin: 8px 0;">
                      <template v-if="article.categories && article.categories.length > 0">
                        <el-tag v-for="cat in article.categories" :key="'cat-' + cat.id" size="small" style="margin-right: 5px;">
                          {{ cat.name }}
                        </el-tag>
                      </template>
                      <template v-if="article.tags && article.tags.length > 0">
                        <el-tag 
                          v-for="tag in article.tags" 
                          :key="'tag-' + tag.id" 
                          size="small"
                          :style="tag.color ? { backgroundColor: tag.color, borderColor: tag.color, color: '#fff', marginRight: '5px' } : { marginRight: '5px' }"
                          :type="tag.color ? undefined : 'info'"
                        >
                          {{ tag.name }}
                        </el-tag>
                      </template>
                    </div>
                    <div style="color: #999; font-size: 12px;">
                      <span>作者：{{ article.author?.nickname || article.author?.username || '未知' }}</span>
                      <span style="margin-left: 20px;">发布时间：{{ formatDate(article.created_at) }}</span>
                      <span style="margin-left: 20px;">浏览：{{ article.view_count || 0 }}</span>
                    </div>
                  </div>
                </div>
              </template>
            </el-list-item>
          </el-list>
          <el-pagination
            v-if="pagination.total > 0"
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadArticles"
            @current-change="loadArticles"
            style="margin-top: 20px; justify-content: center;"
          />
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getArticleList } from '@/api/article'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import type { Article, Pagination } from '@/api/types'

const router = useRouter()
const userStore = useUserStore()

const articles = ref<Article[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const pagination = ref<Pagination>({
  page: 1,
  page_size: 10,
  total: 0,
  total_pages: 0
})

function handleLogout() {
  userStore.logout()
  ElMessage.success('已退出登录')
}

function handleSearchSubmit() {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/search', query: { q: searchKeyword.value.trim() } })
  } else {
    router.push('/search')
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadArticles() {
  loading.value = true
  try {
    const result = await getArticleList({ 
      page: pagination.value.page, 
      page_size: pagination.value.page_size 
    })
    articles.value = result.items || []
    pagination.value = result.pagination || pagination.value
  } catch (error: any) {
    ElMessage.error(error.message || '加载文章失败')
    articles.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.el-header {
  background-color: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 70px;
}

.el-header h1 {
  margin: 0;
  font-size: 24px;
  white-space: nowrap;
  cursor: pointer;
}

.header-search {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: center;
  padding: 0 40px;
}

.header-actions {
  display: flex;
  gap: 10px;
  white-space: nowrap;
}

.el-main {
  padding: 20px;
}
</style>

