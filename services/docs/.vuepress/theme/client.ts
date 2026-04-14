import { defineClientConfig } from 'vuepress/client'
import Layout from './components/layouts/Layout.vue'
import HomeLayout from './components/layouts/HomeLayout.vue'
import NotFound from './components/layouts/NotFound.vue'
import './styles/main.css'

export default defineClientConfig({
  layouts: {
    Layout,
    HomeLayout,
    NotFound,
  },
})
