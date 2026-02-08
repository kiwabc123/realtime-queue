<template>
  <div class="container">
    <h2>รับบัตรคิว</h2>
    <button @click="goQueue">รับบัตรคิว</button>
    <button @click="goClear">ล้างคิว</button>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
function goQueue() {
  router.push('/queue')
}
function goClear() {
  router.push('/clear')
}
</script>
<style scoped>
.container { max-width: 400px; margin: 40px auto; padding: 24px; border-radius: 8px; box-shadow: 0 2px 8px #ccc; background: #fff; }
.queue-display { font-size: 2rem; margin: 16px 0; }
button { margin: 8px; padding: 8px 16px; }
</style>
