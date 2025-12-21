<template>
  <div class="article-detail">
    <el-container>
      <el-main>
        <el-card v-if="article">
          <template #header>
            <div class="article-header">
              <h1>{{ article.title }}</h1>
              <div class="article-meta">
                <span>作者：{{ article.author?.nickname || article.author?.username }}</span>
                <span>发布时间：{{ formatDate(article.created_at) }}</span>
                <span>浏览：{{ article.view_count || 0 }}</span>
                <span>评论：{{ article.comment_count || 0 }}</span>
              </div>
              <div class="article-actions">
                <el-button
                  v-if="canEdit"
                  type="primary"
                  @click="handleEdit"
                >
                  编辑
                </el-button>
                <el-button
                  :type="article.is_liked ? 'primary' : 'default'"
                  @click="handleLike"
                >
                  👍 {{ article.like_count || 0 }}
                </el-button>
              </div>
            </div>
          </template>
          <div v-if="article.cover_image_url" class="article-cover">
            <img :src="article.cover_image_url" :alt="article.title" />
          </div>
          <div class="article-content-wrapper">
            <MarkdownViewer v-if="article.content && article.content.trim()" :content="article.content" />
            <div v-else-if="article.content_html && article.content_html.trim()" class="article-content" v-html="article.content_html"></div>
            <div v-else class="article-content-empty">暂无内容</div>
          </div>
          <div v-if="article.categories && article.categories.length > 0" class="article-categories">
            <el-tag v-for="cat in article.categories" :key="cat.id" style="margin-right: 10px;">
              {{ cat.name }}
            </el-tag>
          </div>
          <div v-if="article.tags && article.tags.length > 0" class="article-tags">
            <el-tag v-for="tag in article.tags" :key="tag.id" type="info" style="margin-right: 10px;">
              {{ tag.name }}
            </el-tag>
          </div>
        </el-card>

        <!-- 评论区域 -->
        <el-card style="margin-top: 20px;" v-if="article">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <h3 style="margin: 0;">评论 ({{ article.comment_count || 0 }})</h3>
            </div>
          </template>
          <div style="margin-bottom: 20px;">
            <CommentEditor
              :article-id="articleId"
              :parent-comment="replyToComment"
              @success="handleCommentSuccess"
              @cancel="replyToComment = null"
            />
          </div>
          <CommentList
            ref="commentListRef"
            :article-id="articleId"
            @reply="handleReply"
          />
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticleDetail, toggleArticleLike } from '@/api/article'
import { ElMessage } from 'element-plus'
import type { Article, Comment } from '@/api/types'
import CommentEditor from '@/components/Comment/CommentEditor.vue'
import CommentList from '@/components/Comment/CommentList.vue'
import MarkdownViewer from '@/components/MarkdownViewer.vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const article = ref<Article | null>(null)
const articleId = ref(0)
const replyToComment = ref<Comment | null>(null)
const commentListRef = ref<InstanceType<typeof CommentList> | null>(null)

const canEdit = computed(() => {
  return userStore.isLoggedIn && article.value && article.value.author?.id === userStore.user?.id
})

onMounted(async () => {
  articleId.value = parseInt(route.params.id as string)
  await loadArticle()
})

async function loadArticle() {
  try {
    const data = await getArticleDetail(articleId.value)
    article.value = data
    // 调试信息
    console.log('Article loaded:', {
      id: data.id,
      title: data.title,
      hasContent: !!data.content,
      hasContentHtml: !!data.content_html,
      contentLength: data.content?.length || 0,
      contentHtmlLength: data.content_html?.length || 0
    })
  } catch (error: any) {
    console.error('Load article error:', error)
    ElMessage.error(error.message || '加载文章失败')
  }
}

async function handleLike() {
  if (!article.value) return
  try {
    const result = await toggleArticleLike(article.value.id)
    article.value.is_liked = result.is_liked
    if (result.is_liked) {
      article.value.like_count++
    } else {
      article.value.like_count = Math.max(0, article.value.like_count - 1)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

function handleReply(comment: Comment) {
  replyToComment.value = comment
}

function handleCommentSuccess() {
  replyToComment.value = null
  if (commentListRef.value) {
    commentListRef.value.loadComments()
  }
  loadArticle() // 更新评论数
}

function handleEdit() {
  router.push(`/article/${articleId.value}/edit`)
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.article-detail {
  min-height: 100vh;
  padding: 20px;
}

.article-header h1 {
  margin: 0 0 10px 0;
}

.article-meta {
  color: #999;
  font-size: 14px;
  margin-bottom: 10px;
}

.article-meta span {
  margin-right: 15px;
}

.article-actions {
  margin-top: 10px;
}

.article-cover {
  margin: 20px 0;
  width: 100%;
  max-height: 400px;
  overflow: hidden;
  border-radius: 8px;
}

.article-cover img {
  width: 100%;
  height: auto;
  object-fit: cover;
  display: block;
  background-color: #f5f5f5;
}

.article-cover img[src=""] {
  display: none;
}

.article-content-wrapper {
  margin: 20px 0;
  min-height: 100px;
}

.article-content {
  line-height: 1.8;
}

.article-content-empty {
  padding: 40px;
  text-align: center;
  color: #999;
}

.article-categories,
.article-tags {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}
</style>

