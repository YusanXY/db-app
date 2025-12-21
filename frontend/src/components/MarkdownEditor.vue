<template>
  <div class="markdown-editor">
    <div ref="editorRef" :id="editorId" class="vditor-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const props = defineProps<{
  modelValue?: string
  height?: number
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// 生成唯一的编辑器ID（在组件创建时生成，确保ID稳定）
const editorId = ref(`vditor-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`)
const editorRef = ref<HTMLDivElement | null>(null)
let vditor: Vditor | null = null
let pasteHandler: ((e: ClipboardEvent) => void) | null = null

onMounted(async () => {
  await nextTick()
  
  // 等待DOM完全渲染
  await new Promise(resolve => setTimeout(resolve, 300))
  
  const container = editorRef.value || document.getElementById(editorId.value)
  if (!container) {
    console.error('MarkdownEditor: Container element not found', { editorId: editorId.value, editorRef: editorRef.value })
    return
  }
  
  // 添加粘贴图片事件监听
  pasteHandler = async (e: ClipboardEvent) => {
    const items = e.clipboardData?.items
    if (!items || !vditor) return
    
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      if (item.type.indexOf('image') !== -1) {
        e.preventDefault()
        e.stopPropagation()
        const file = item.getAsFile()
        if (file) {
          // 上传图片
          const formData = new FormData()
          formData.append('file', file)
          
          try {
            const token = localStorage.getItem('token')
            const response = await fetch('/api/v1/files/upload', {
              method: 'POST',
              headers: token ? { Authorization: `Bearer ${token}` } : {},
              body: formData
            })
            
            const result = await response.json()
            if (result.code === 200 && result.data) {
              // 插入图片到编辑器
              const imageUrl = result.data.url
              const markdown = `![${file.name}](${imageUrl})`
              vditor.insertValue(markdown)
            }
          } catch (error) {
            console.error('Paste image upload error:', error)
          }
        }
        break
      }
    }
  }
  
  // 监听粘贴事件
  container.addEventListener('paste', pasteHandler)
  
  try {
    const cacheId = editorId.value
    console.log('MarkdownEditor: Initializing with cache.id:', cacheId)
    
    vditor = new Vditor(container as HTMLElement, {
      height: props.height || 500,
      placeholder: props.placeholder || '请输入Markdown内容...',
      mode: 'sv', // 使用分屏模式，更稳定
      cache: {
        id: cacheId, // 添加必需的cache.id
        enable: true
      },
      toolbar: [
        'headings',
        'bold',
        'italic',
        'strike',
        '|',
        'line',
        'quote',
        'list',
        'ordered-list',
        'check',
        'outdent',
        'indent',
        '|',
        'code',
        'inline-code',
        'insert-before',
        'insert-after',
        '|',
        'link',
        'table',
        '|',
        'undo',
        'redo',
        '|',
        'fullscreen',
        'preview',
        'outline',
        '|',
        'help'
      ],
      value: props.modelValue || '',
      input: (value: string) => {
        emit('update:modelValue', value)
      },
      upload: {
        accept: 'image/*',
        url: '/api/v1/files/upload',
        fieldName: 'file',
        headers: (() => {
          // 获取token
          const token = localStorage.getItem('token')
          return token ? { Authorization: `Bearer ${token}` } : {}
        })(),
        format: (files: File[], responseText: string) => {
          try {
            const res = JSON.parse(responseText)
            if (res.code === 200 && res.data) {
              return JSON.stringify([res.data.url])
            }
          } catch (e) {
            console.error('Upload response parse error:', e)
          }
          return '[]'
        },
        linkToImgUrl: (url: string) => {
          return url
        }
      },
      // 支持粘贴图片
      paste: {
        enable: true
      }
    })
    console.log('MarkdownEditor: Vditor initialized successfully', vditor)
  } catch (error) {
    console.error('MarkdownEditor: Failed to initialize Vditor:', error)
  }
})

onBeforeUnmount(() => {
  // 移除粘贴事件监听
  if (editorRef.value && pasteHandler) {
    editorRef.value.removeEventListener('paste', pasteHandler)
    pasteHandler = null
  }
  
  if (vditor) {
    vditor.destroy()
    vditor = null
  }
})

watch(() => props.modelValue, (newValue) => {
  if (vditor && vditor.getValue() !== newValue) {
    vditor.setValue(newValue || '')
  }
})

defineExpose({
  getValue: () => vditor?.getValue() || '',
  setValue: (value: string) => vditor?.setValue(value || ''),
  focus: () => vditor?.focus(),
  blur: () => vditor?.blur()
})
</script>

<style>
/* 全局引入vditor样式，避免scoped影响 */
@import 'vditor/dist/index.css';
</style>

<style scoped>
.markdown-editor {
  width: 100%;
  min-height: 500px;
  position: relative;
}

.vditor-container {
  width: 100%;
  min-height: 500px;
  position: relative;
}

/* 确保vditor样式正确应用 */
.markdown-editor :deep(.vditor) {
  width: 100% !important;
  min-height: 500px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.markdown-editor :deep(.vditor-content) {
  min-height: 400px;
}

.markdown-editor :deep(.vditor-toolbar) {
  width: 100%;
}

.markdown-editor :deep(.vditor-preview) {
  width: 100%;
}

.markdown-editor :deep(.vditor-ir) {
  width: 100%;
}
</style>

