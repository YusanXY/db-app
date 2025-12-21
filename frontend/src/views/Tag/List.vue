<template>
  <div class="tag-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>标签管理</h2>
          <el-button v-if="userStore.isLoggedIn" type="primary" @click="showCreateDialog = true">
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
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="slug" label="标识" />
        <el-table-column prop="article_count" label="文章数" width="100" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="text" size="small" @click="editTag(row)">编辑</el-button>
            <el-button type="text" size="small" danger @click="deleteTag(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="创建标签" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="标识">
          <el-input v-model="form.slug" placeholder="留空自动生成" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-input v-model="form.color" placeholder="CSS颜色值，如 #409EFF" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="loading">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTagList, createTag, deleteTag as deleteTagApi } from '@/api/tag'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Tag } from '@/api/types'

const userStore = useUserStore()
const tags = ref<Tag[]>([])
const keyword = ref('')
const showCreateDialog = ref(false)
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

async function handleCreate() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  loading.value = true
  try {
    await createTag(form.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    form.value = { name: '', slug: '', description: '', color: '' }
    loadTags()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
  } finally {
    loading.value = false
  }
}

function editTag(tag: Tag) {
  ElMessage.info('编辑功能待实现')
}

async function deleteTag(tag: Tag) {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', {
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

