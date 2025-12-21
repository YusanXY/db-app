<template>
  <div class="category-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>分类管理</h2>
          <el-button v-if="userStore.isLoggedIn" type="primary" @click="showCreateDialog = true">
            创建分类
          </el-button>
        </div>
      </template>
      <el-tree
        :data="categories"
        :props="{ children: 'children', label: 'name' }"
        default-expand-all
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <span>{{ data.name }}</span>
            <div class="actions">
              <el-button type="text" size="small" @click="editCategory(data)">编辑</el-button>
              <el-button type="text" size="small" danger @click="deleteCategory(data)">删除</el-button>
            </div>
          </div>
        </template>
      </el-tree>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="创建分类" width="500px">
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
        <el-form-item label="父分类">
          <el-select v-model="form.parent_id" clearable placeholder="选择父分类">
            <el-option
              v-for="cat in flatCategories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
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
import { ref, onMounted, computed } from 'vue'
import { getCategoryList, createCategory, deleteCategory as deleteCategoryApi } from '@/api/category'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Category } from '@/api/types'

const userStore = useUserStore()
const categories = ref<Category[]>([])
const showCreateDialog = ref(false)
const loading = ref(false)
const form = ref({
  name: '',
  slug: '',
  description: '',
  parent_id: undefined as number | undefined
})

const flatCategories = computed(() => {
  const flatten = (cats: Category[]): Category[] => {
    const result: Category[] = []
    cats.forEach(cat => {
      result.push(cat)
      if (cat.children) {
        result.push(...flatten(cat.children))
      }
    })
    return result
  }
  return flatten(categories.value)
})

async function loadCategories() {
  try {
    categories.value = await getCategoryList({ tree: true })
  } catch (error: any) {
    ElMessage.error(error.message || '加载分类失败')
  }
}

async function handleCreate() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  loading.value = true
  try {
    await createCategory(form.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    form.value = { name: '', slug: '', description: '', parent_id: undefined }
    loadCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
  } finally {
    loading.value = false
  }
}

function editCategory(category: Category) {
  ElMessage.info('编辑功能待实现')
}

async function deleteCategory(category: Category) {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteCategoryApi(category.id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.category-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
  padding-right: 8px;
}

.actions {
  display: flex;
  gap: 10px;
}
</style>

