<template>
  <div class="comment-list">
    <div v-if="comments.length === 0 && !loading" class="empty">
      <el-empty description="暂无评论" />
    </div>
    <el-skeleton v-if="loading" :rows="3" animated />
    <div v-for="comment in comments" :key="comment.id" class="comment-item">
      <CommentItem :comment="comment" @reply="handleReply" @delete="handleDelete" @like="handleLike" />
    </div>
    <el-pagination
      v-if="pagination.total > 0"
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.page_size"
      :total="pagination.total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      @size-change="loadComments"
      @current-change="loadComments"
      style="margin-top: 20px;"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCommentList, deleteComment, toggleCommentLike } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Comment, Pagination } from '@/api/types'
import CommentItem from './CommentItem.vue'

const props = defineProps<{
  articleId: number
}>()

const emit = defineEmits<{
  reply: [comment: Comment]
}>()

const comments = ref<Comment[]>([])
const loading = ref(false)
const pagination = ref<Pagination>({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0
})

async function loadComments() {
  loading.value = true
  try {
    const result = await getCommentList(props.articleId, {
      page: pagination.value.page,
      page_size: pagination.value.page_size
    })
    comments.value = result.items || []
    pagination.value = result.pagination || pagination.value
  } catch (error: any) {
    ElMessage.error(error.message || '加载评论失败')
  } finally {
    loading.value = false
  }
}

function handleReply(comment: Comment) {
  emit('reply', comment)
}

async function handleDelete(comment: Comment) {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteComment(comment.id)
    ElMessage.success('删除成功')
    loadComments()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

async function handleLike(comment: Comment) {
  try {
    const result = await toggleCommentLike(comment.id)
    comment.is_liked = result.is_liked
    if (result.is_liked) {
      comment.like_count++
    } else {
      comment.like_count = Math.max(0, comment.like_count - 1)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

onMounted(() => {
  loadComments()
})

defineExpose({
  loadComments
})
</script>

<style scoped>
.comment-list {
  margin-top: 20px;
}

.comment-item {
  margin-bottom: 20px;
}

.empty {
  padding: 40px 0;
  text-align: center;
}
</style>

