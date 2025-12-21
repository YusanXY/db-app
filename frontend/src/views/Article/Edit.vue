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
              <el-button type="primary" @click="handleSubmit" :loading="submitting">保存</el-button>
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
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getArticleDetail, updateArticle } from '@/api/article'
import { ElMessage } from 'element-plus'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const articleId = ref(0)
const loading = ref(false)
const submitting = ref(false)
const formRef = ref()

const form = ref({
  title: '',
  content: '',
  summary: '',
  status: 'draft'
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

onMounted(async () => {
  articleId.value = parseInt(route.params.id as string)
  await loadArticle()
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
      content: article.content,
      summary: article.summary || '',
      status: article.status
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
</style>

