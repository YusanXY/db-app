<template>
  <div class="tag-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>标签管理</h2>
          <el-button v-if="userStore.isLoggedIn" type="primary" @click="openCreateDialog">
            创建标签
          </el-button>
        </div>
      </template>
      <el-input
        v-model="keyword"
        placeholder="搜索标签"
        style="margin-bottom: 20px;"
        clearable
        @input="loadTags"
      />
      <el-table :data="tags" style="width: 100%">
        <el-table-column prop="name" label="名称">
          <template #default="{ row }">
            <el-tag v-if="row.color" :style="{ backgroundColor: row.color, borderColor: row.color, color: '#fff' }">
              {{ row.name }}
            </el-tag>
            <span v-else>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="slug" label="标识" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="article_count" label="文章数" width="100" />
        <el-table-column label="操作" width="150" v-if="userStore.isLoggedIn">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="isEdit ? '编辑标签' : '创建标签'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="标识">
          <el-input v-model="form.slug" placeholder="留空自动生成" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="标签描述（可选）" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-color-picker v-model="form.color" />
          <span style="margin-left: 10px; color: #999;">{{ form.color || '未设置' }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTagList, createTag, updateTag, deleteTag as deleteTagApi } from '@/api/tag'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Tag } from '@/api/types'

const userStore = useUserStore()
const tags = ref<Tag[]>([])
const keyword = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const loading = ref(false)
const form = ref({
  name: '',
  slug: '',
  description: '',
  color: ''
})

async function loadTags() {
  try {
    tags.value = await getTagList({ keyword: keyword.value })
  } catch (error: any) {
    ElMessage.error(error.message || '加载标签失败')
  }
}

function resetForm() {
  form.value = { name: '', slug: '', description: '', color: '' }
  editingId.value = null
  isEdit.value = false
}

function openCreateDialog() {
  resetForm()
  showDialog.value = true
}

function openEditDialog(tag: Tag) {
  isEdit.value = true
  editingId.value = tag.id
  form.value = {
    name: tag.name,
    slug: tag.slug,
    description: tag.description || '',
    color: tag.color || ''
  }
  showDialog.value = true
}

async function handleSubmit() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  loading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await updateTag(editingId.value, {
        name: form.value.name,
        description: form.value.description,
        color: form.value.color
      })
      ElMessage.success('更新成功')
    } else {
      await createTag(form.value)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    resetForm()
    loadTags()
  } catch (error: any) {
    ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
  } finally {
    loading.value = false
  }
}

async function handleDelete(tag: Tag) {
  try {
    await ElMessageBox.confirm(`确定要删除标签"${tag.name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteTagApi(tag.id)
    ElMessage.success('删除成功')
    loadTags()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tag-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

