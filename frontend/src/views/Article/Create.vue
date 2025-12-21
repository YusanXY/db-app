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
            <el-form-item label="内容" prop="content">
              <el-input v-model="form.content" type="textarea" :rows="10" />
            </el-form-item>
            <el-form-item label="摘要" prop="summary">
              <el-input v-model="form.summary" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
              <el-button @click="$router.back()">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createArticle } from '@/api/article'
import { ElMessage } from 'element-plus'

const router = useRouter()

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

const formRef = ref()
const loading = ref(false)

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await createArticle(form.value)
        ElMessage.success('创建成功')
        router.push('/')
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
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
}
</style>

