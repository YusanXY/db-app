<template>
  <div class="article-detail">
    <el-container>
      <el-main>
        <el-card v-if="article">
          <template #header>
            <h1>{{ article.title }}</h1>
          </template>
          <div v-html="article.content_html || article.content"></div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleDetail } from '@/api/article'
import type { Article } from '@/api/types'

const route = useRoute()
const article = ref<Article | null>(null)

onMounted(async () => {
  const id = parseInt(route.params.id as string)
  try {
    article.value = await getArticleDetail(id)
  } catch (error) {
    console.error('加载文章失败', error)
  }
})
</script>

<style scoped>
.article-detail {
  min-height: 100vh;
}
</style>

