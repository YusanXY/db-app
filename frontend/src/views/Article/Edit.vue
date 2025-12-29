<template>
  <div class="article-edit">
    <el-container>
      <el-main>
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <h2>编辑文章</h2>
              <el-button @click="$router.back()">返回</el-button>
            </div>
          </template>
          <el-form :model="form" :rules="rules" ref="formRef" v-loading="loading">
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
                  name="file"
                  accept="image/*"
                >
                  <img v-if="form.cover_image_url" :src="form.cover_image_url" class="cover-image" />
                  <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
                </el-upload>
                <div v-if="form.cover_image_url && form.cover_image_url !== '/uploads/default/cover-default.jpg'" class="cover-actions">
                  <el-button type="danger" size="small" @click="removeCover">删除封面</el-button>
                </div>
              </div>
            </el-form-item>
            <el-form-item label="内容" prop="content">
              <MarkdownEditor v-model="form.content" :height="500" />
            </el-form-item>
            <el-form-item label="摘要" prop="summary">
              <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="文章摘要，将显示在文章列表中" />
            </el-form-item>
            <el-form-item label="分类">
              <el-select v-model="form.category_ids" multiple placeholder="选择分类" style="width: 100%">
                <el-option
                  v-for="cat in flatCategories"
                  :key="cat.id"
                  :label="cat.displayName"
                  :value="cat.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="标签">
              <el-select v-model="form.tag_ids" multiple filterable placeholder="选择标签" style="width: 100%">
                <el-option
                  v-for="tag in tags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.id"
                >
                  <span>{{ tag.name }}</span>
                  <el-tag v-if="tag.color" :style="{ backgroundColor: tag.color, marginLeft: '10px' }" size="small">
                    &nbsp;
                  </el-tag>
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio label="draft">草稿</el-radio>
                <el-radio label="published">发布</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="submitting">保存并发布</el-button>
              <el-button @click="handleSaveDraft" :loading="submitting">保存草稿</el-button>
              <el-button @click="$router.back()">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getArticleDetail, updateArticle } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useUserStore } from '@/stores/user'
import type { Category, Tag } from '@/api/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const articleId = ref(0)
const loading = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const formRef = ref()
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover_image_url: '/uploads/default/cover-default.jpg',
  status: 'draft',
  category_ids: [] as number[],
  tag_ids: [] as number[]
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const uploadAction = computed(() => '/api/v1/files/upload')

const uploadHeaders = computed(() => {
  const token = userStore.token
  return token ? { Authorization: `Bearer ${token}` } : {}
})

// 扁平化分类列表用于选择
const flatCategories = computed(() => {
  const flatten = (cats: Category[], prefix: string = ''): { id: number; displayName: string }[] => {
    const result: { id: number; displayName: string }[] = []
    cats.forEach(cat => {
      result.push({ id: cat.id, displayName: prefix + cat.name })
      if (cat.children && cat.children.length > 0) {
        result.push(...flatten(cat.children, prefix + cat.name + ' / '))
      }
    })
    return result
  }
  return flatten(categories.value)
})

async function loadCategories() {
  try {
    categories.value = await getCategoryList({ tree: true })
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

async function loadTags() {
  try {
    tags.value = await getTagList()
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

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

function handleCoverSuccess(response: any) {
  uploading.value = false
  if (response && response.code === 200 && response.data) {
    form.value.cover_image_url = response.data.url
    ElMessage.success('封面图片上传成功')
  } else {
    ElMessage.error(response?.message || '上传失败')
  }
}

function handleCoverError(error: any) {
  uploading.value = false
  let errorMsg = '封面图片上传失败'
  if (error?.message) {
    try {
      const parsed = JSON.parse(error.message)
      errorMsg = parsed.message || errorMsg
    } catch {
      errorMsg = error.message || errorMsg
    }
  }
  ElMessage.error(errorMsg)
}

function removeCover() {
  form.value.cover_image_url = '/uploads/default/cover-default.jpg'
}

onMounted(async () => {
  articleId.value = parseInt(route.params.id as string)
  await Promise.all([loadArticle(), loadCategories(), loadTags()])
})

async function loadArticle() {
  loading.value = true
  try {
    const article = await getArticleDetail(articleId.value)
    
    // 检查权限
    if (article.author?.id !== userStore.user?.id) {
      ElMessage.error('您没有权限编辑此文章')
      router.push(`/article/${articleId.value}`)
      return
    }

    form.value = {
      title: article.title,
      content: article.content || '',
      summary: article.summary || '',
      cover_image_url: article.cover_image_url || '/uploads/default/cover-default.jpg',
      status: article.status,
      category_ids: article.categories?.map(c => c.id) || [],
      tag_ids: article.tags?.map(t => t.id) || []
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载文章失败')
    router.push('/')
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        await updateArticle(articleId.value, {
          ...form.value,
          status: 'published'
        })
        ElMessage.success('文章已发布')
        router.push(`/article/${articleId.value}`)
      } catch (error: any) {
        ElMessage.error(error.message || '更新失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

async function handleSaveDraft() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        await updateArticle(articleId.value, {
          ...form.value,
          status: 'draft'
        })
        ElMessage.success('草稿已保存')
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        submitting.value = false
      }
    }
  })
}
</script>

<style scoped>
.article-edit {
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

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

.cover-actions {
  display: flex;
  gap: 10px;
}
</style>
