<template>
  <div class="home">
    <el-container>
      <el-header>
        <h1>百科Web应用</h1>
        <div style="float: right;">
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
          <el-row :gutter="20">
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
                <el-button type="primary" @click="loadArticles">浏览文章</el-button>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover">
                <h3>社区功能</h3>
                <p>评论、点赞、关注</p>
                <el-button type="primary">了解更多</el-button>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import { getArticleList } from '@/api/article'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

function handleLogout() {
  userStore.logout()
  ElMessage.success('已退出登录')
}

async function loadArticles() {
  try {
    const result = await getArticleList({ page: 1, page_size: 10 })
    ElMessage.success(`加载了 ${result.items.length} 篇文章`)
  } catch (error) {
    ElMessage.error('加载文章失败')
  }
}
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
}

.el-header h1 {
  margin: 0;
  font-size: 24px;
}

.el-main {
  padding: 20px;
}
</style>

