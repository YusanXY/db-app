<template>
  <div class="article-create">
    <el-container>
      <el-main>
        <el-card>
          <template #header>
            <h2>创建文章</h2>
          </template>
          <el-form :model="form" :rules="rules" ref="formRef">
            <el-form-item label="标题" prop="title">
              <el-input v-model="form.title" />
            </el-form-item>
            <el-form-item label="封面图片">
              <div class="cover-image-upload">
                <el-upload
                  class="cover-uploader"
                  :action="uploadAction"
                  :headers="uploadHeaders"
                  :show-file-list="false"
                  :on-success="handleCoverSuccess"
                  :on-error="handleCoverError"
                  :before-upload="beforeCoverUpload"
                  accept="image/*"
                >
                  <img v-if="form.cover_image_url" :src="form.cover_image_url" class="cover-image" />
                  <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
                </el-upload>
                <div v-if="form.cover_image_url" class="cover-actions">
                  <el-button type="danger" size="small" @click="removeCover">删除</el-button>
                </div>
              </div>
            </el-form-item>
            <el-form-item label="内容" prop="content">
              <MarkdownEditor v-model="form.content" :height="500" />
            </el-form-item>
            <el-form-item label="摘要" prop="summary">
              <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="文章摘要，将显示在文章列表中" />
            </el-form-item>
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio label="draft">草稿</el-radio>
                <el-radio label="published">发布</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="loading">发布</el-button>
              <el-button @click="handleSaveDraft" :loading="loading">保存草稿</el-button>
              <el-button @click="$router.back()">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { createArticle } from '@/api/article'
import { uploadFile } from '@/api/file'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useUserStore } from '@/stores/user'
import request from '@/api/request'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover_image_url: '/uploads/default/cover-default.jpg', // 默认封面图片
  status: 'published' // 默认发布
})

const uploadAction = computed(() => {
  // 使用相对路径，由代理转发
  return '/api/v1/files/upload'
})

const uploadHeaders = computed(() => {
  const token = userStore.token
  return token ? { Authorization: `Bearer ${token}` } : {}
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const formRef = ref()
const loading = ref(false)
const uploading = ref(false)

function beforeCoverUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB!')
    return false
  }
  uploading.value = true
  return true
}

async function handleCoverSuccess(response: any) {
  uploading.value = false
  if (response.code === 200 && response.data) {
    form.value.cover_image_url = response.data.url
    ElMessage.success('封面图片上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

function handleCoverError() {
  uploading.value = false
  ElMessage.error('封面图片上传失败')
}

function removeCover() {
  form.value.cover_image_url = '/uploads/default/cover-default.jpg' // 恢复默认图片
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await createArticle({
          ...form.value,
          status: 'published'
        })
        ElMessage.success('文章已发布')
        router.push('/')
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
      } finally {
        loading.value = false
      }
    }
  })
}

async function handleSaveDraft() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await createArticle({
          ...form.value,
          status: 'draft'
        })
        ElMessage.success('草稿已保存')
        router.push('/')
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.article-create {
  min-height: 100vh;
  padding: 20px;
}

:deep(.el-form-item__content) {
  width: 100%;
  line-height: normal;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

:deep(.el-form-item__label) {
  width: auto;
  padding-right: 12px;
}

.cover-image-upload {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.cover-uploader {
  :deep(.el-upload) {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: all 0.3s;
    width: 300px;
    height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      border-color: #409eff;
    }
  }
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  background-color: #f5f5f5;
}

.cover-image[src=""] {
  display: none;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

.cover-actions {
  display: flex;
  gap: 10px;
}
</style>

