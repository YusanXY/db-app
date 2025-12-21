<template>
  <div class="comment-item">
    <div class="comment-header">
      <el-avatar :src="comment.user.avatar_url" :size="40">
        {{ comment.user.nickname?.[0] || comment.user.username[0] }}
      </el-avatar>
      <div class="comment-info">
        <span class="username">{{ comment.user.nickname || comment.user.username }}</span>
        <span class="time">{{ formatTime(comment.created_at) }}</span>
      </div>
    </div>
    <div class="comment-content" v-html="comment.content_html || comment.content"></div>
    <div class="comment-actions">
      <el-button
        :type="comment.is_liked ? 'primary' : 'default'"
        size="small"
        @click="handleLike"
      >
        👍 {{ comment.like_count || 0 }}
      </el-button>
      <el-button
        v-if="userStore.isLoggedIn"
        type="text"
        size="small"
        @click="handleReply"
      >
        回复
      </el-button>
      <el-button
        v-if="canDelete"
        type="text"
        size="small"
        danger
        @click="handleDelete"
      >
        删除
      </el-button>
    </div>
    <div v-if="comment.replies && comment.replies.length > 0" class="replies">
      <CommentItem
        v-for="reply in comment.replies"
        :key="reply.id"
        :comment="reply"
        @reply="handleReply"
        @delete="handleDelete"
        @like="handleLike"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import type { Comment } from '@/api/types'

const props = defineProps<{
  comment: Comment
}>()

const emit = defineEmits<{
  reply: [comment: Comment]
  delete: [comment: Comment]
  like: [comment: Comment]
}>()

const userStore = useUserStore()

const canDelete = computed(() => {
  return userStore.isLoggedIn && (
    userStore.user?.id === props.comment.user.id ||
    userStore.user?.role === 'admin'
  )
})

function formatTime(time: string) {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}

function handleReply() {
  emit('reply', props.comment)
}

function handleDelete() {
  emit('delete', props.comment)
}

function handleLike() {
  emit('like', props.comment)
}
</script>

<style scoped>
.comment-item {
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.comment-info {
  margin-left: 10px;
  flex: 1;
}

.username {
  font-weight: bold;
  margin-right: 10px;
}

.time {
  color: #999;
  font-size: 12px;
}

.comment-content {
  margin: 10px 0;
  line-height: 1.6;
}

.comment-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.replies {
  margin-top: 15px;
  padding-left: 20px;
  border-left: 2px solid #eee;
}
</style>

