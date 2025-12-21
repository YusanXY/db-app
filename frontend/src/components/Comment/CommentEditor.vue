<template>
  <div class="comment-editor">
    <el-input
      v-model="content"
      type="textarea"
      :rows="4"
      placeholder="写下你的评论..."
      :disabled="!userStore.isLoggedIn"
    />
    <div v-if="!userStore.isLoggedIn" class="login-tip">
      请先<el-link type="primary" @click="$router.push('/login')">登录</el-link>后发表评论
    </div>
    <div v-if="parentComment" class="reply-to">
      回复 <strong>{{ parentComment.user.nickname || parentComment.user.username }}</strong>
      <el-button type="text" size="small" @click="cancelReply">取消</el-button>
    </div>
    <div class="actions">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading" :disabled="!content.trim()">
        发表
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { createComment } from '@/api/comment'
import { ElMessage } from 'element-plus'
import type { Comment } from '@/api/types'

const props = defineProps<{
  articleId: number
  parentComment?: Comment | null
}>()

const emit = defineEmits<{
  success: []
  cancel: []
}>()

const userStore = useUserStore()
const content = ref('')
const loading = ref(false)

watch(() => props.parentComment, (newVal) => {
  if (newVal) {
    content.value = `@${newVal.user.nickname || newVal.user.username} `
  }
})

function handleSubmit() {
  if (!content.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  loading.value = true
  createComment(props.articleId, {
    content: content.value.trim(),
    parent_id: props.parentComment?.id
  })
    .then(() => {
      ElMessage.success('评论成功')
      content.value = ''
      emit('success')
    })
    .catch((error: any) => {
      ElMessage.error(error.message || '评论失败')
    })
    .finally(() => {
      loading.value = false
    })
}

function handleCancel() {
  content.value = ''
  emit('cancel')
}

function cancelReply() {
  emit('cancel')
}
</script>

<style scoped>
.comment-editor {
  margin-bottom: 20px;
}

.login-tip {
  margin-top: 10px;
  color: #999;
  font-size: 12px;
}

.reply-to {
  margin-top: 10px;
  padding: 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 12px;
}

.actions {
  margin-top: 10px;
  text-align: right;
}
</style>

