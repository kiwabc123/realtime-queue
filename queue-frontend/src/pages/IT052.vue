<template>
  <div class="container">
    <h2>หมายเลขคิวของคุณ</h2>
    <div class="queue-display">{{ queueNumber }}</div>
    <button @click="goBack">กลับหน้ารับบัตรคิว</button>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useQueueStore } from '../store/queueStore'
const router = useRouter()
const queueStore = useQueueStore()
const queueNumber = queueStore.yourQueue
function goBack() {
  router.push('/')
}
</script>
<style scoped>
.container { max-width: 400px; margin: 40px auto; padding: 24px; border-radius: 8px; box-shadow: 0 2px 8px #ccc; background: #fff; }
.queue-display { font-size: 2rem; margin: 16px 0; }
button { margin: 8px; padding: 8px 16px; }
</style>
