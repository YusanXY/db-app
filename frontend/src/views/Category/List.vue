<template>
  <div class="category-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>分类管理</h2>
          <el-button v-if="userStore.isLoggedIn" type="primary" @click="openCreateDialog">
            创建分类
          </el-button>
        </div>
      </template>
      <el-table :data="flatCategoriesWithLevel" style="width: 100%" row-key="id">
        <el-table-column prop="name" label="名称">
          <template #default="{ row }">
            <span :style="{ paddingLeft: row.level * 20 + 'px' }">
              <span v-if="row.level > 0">└─ </span>{{ row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="slug" label="标识" width="150" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="article_count" label="文章数" width="100" />
        <el-table-column prop="is_active" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" v-if="userStore.isLoggedIn">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="isEdit ? '编辑分类' : '创建分类'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="标识">
          <el-input v-model="form.slug" placeholder="留空自动生成" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="分类描述（可选）" />
        </el-form-item>
        <el-form-item label="父分类">
          <el-select v-model="form.parent_id" clearable placeholder="选择父分类（可选）" style="width: 100%">
            <el-option
              v-for="cat in availableParentCategories"
              :key="cat.id"
              :label="cat.displayName"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" v-if="isEdit">
          <el-switch v-model="form.is_active" active-text="启用" inactive-text="禁用" />
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
import { ref, onMounted, computed } from 'vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory as deleteCategoryApi } from '@/api/category'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Category } from '@/api/types'

const userStore = useUserStore()
const categories = ref<Category[]>([])
const showDialog = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const loading = ref(false)
const form = ref({
  name: '',
  slug: '',
  description: '',
  parent_id: undefined as number | undefined,
  is_active: true
})

// 扁平化分类列表（带层级信息）
const flatCategoriesWithLevel = computed(() => {
  const flatten = (cats: Category[], level: number = 0): (Category & { level: number })[] => {
    const result: (Category & { level: number })[] = []
    cats.forEach(cat => {
      result.push({ ...cat, level })
      if (cat.children && cat.children.length > 0) {
        result.push(...flatten(cat.children, level + 1))
      }
    })
    return result
  }
  return flatten(categories.value)
})

// 可选的父分类（排除自己及其子分类）
const availableParentCategories = computed(() => {
  const result: { id: number; displayName: string }[] = []
  const flatten = (cats: Category[], prefix: string = '') => {
    cats.forEach(cat => {
      // 编辑时排除自己
      if (isEdit.value && editingId.value === cat.id) return
      result.push({ id: cat.id, displayName: prefix + cat.name })
      if (cat.children && cat.children.length > 0) {
        flatten(cat.children, prefix + cat.name + ' / ')
      }
    })
  }
  flatten(categories.value)
  return result
})

async function loadCategories() {
  try {
    categories.value = await getCategoryList({ tree: true })
  } catch (error: any) {
    ElMessage.error(error.message || '加载分类失败')
  }
}

function resetForm() {
  form.value = { name: '', slug: '', description: '', parent_id: undefined, is_active: true }
  editingId.value = null
  isEdit.value = false
}

function openCreateDialog() {
  resetForm()
  showDialog.value = true
}

function openEditDialog(category: Category) {
  isEdit.value = true
  editingId.value = category.id
  form.value = {
    name: category.name,
    slug: category.slug,
    description: category.description || '',
    parent_id: category.parent_id,
    is_active: category.is_active !== false
  }
  showDialog.value = true
}

async function handleSubmit() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  loading.value = true
  try {
    if (isEdit.value && editingId.value) {
      await updateCategory(editingId.value, {
        name: form.value.name,
        description: form.value.description,
        parent_id: form.value.parent_id,
        is_active: form.value.is_active
      })
      ElMessage.success('更新成功')
    } else {
      await createCategory({
        name: form.value.name,
        slug: form.value.slug,
        description: form.value.description,
        parent_id: form.value.parent_id
      })
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    resetForm()
    loadCategories()
  } catch (error: any) {
    ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
  } finally {
    loading.value = false
  }
}

async function handleDelete(category: Category) {
  try {
    await ElMessageBox.confirm(`确定要删除分类"${category.name}"吗？`, '提示', {
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

