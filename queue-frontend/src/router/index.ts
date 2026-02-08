import { createRouter, createWebHistory } from 'vue-router'
import IT051 from '../pages/IT051.vue'
import IT052 from '../pages/IT052.vue'
import IT053 from '../pages/IT053.vue'

const routes = [
  { path: '/', name: 'IT051', component: IT051 },
  { path: '/queue', name: 'IT052', component: IT052 },
  { path: '/clear', name: 'IT053', component: IT053 },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
