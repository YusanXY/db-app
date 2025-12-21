import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './styles/main.scss'

const app = createApp(App)

// 注册插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 图标按需导入，不全局注册，避免循环依赖
// 如需使用图标，在组件中单独导入：
// import { Edit, Delete } from '@element-plus/icons-vue'

app.mount('#app')

