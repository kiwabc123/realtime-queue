<template>
  <div class="container">
    <h2 class="title">ล้างคิว</h2>
    <div class="queue-display">{{ queueNumber }}</div>
    <div class="button-group">
      <button @click="clearQueue">ล้างคิว</button>
      <button @click="goBack">กลับหน้ารับบัตรคิว</button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useQueueStore } from '../store/queueStore'
const router = useRouter()
const queueStore = useQueueStore({ skipAutoFetch: true })
const queueNumber = queueStore.currentQueue
function clearQueue() {
  queueStore.clearQueueWS()
}
function goBack() {
  router.push('/')
}
</script>
<style scoped>
.title {
  text-align: center;
  font-size: 2.2rem;
  color: #0d47a1;
  margin-bottom: 32px;
  font-weight: bold;
  letter-spacing: 1px;
}
.button-group {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}
.container {
  max-width: 400px;
  margin: 40px auto;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 16px #b3c6e7;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
}
.queue-display {
  font-size: 3rem;
  margin: 40px auto;
  padding: 32px 48px;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 4px 16px #90caf9;
  color: #0d47a1;
  text-align: center;
  max-width: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
}
button {
  margin: 8px;
  padding: 10px 24px;
  border-radius: 6px;
  background: #1976d2;
  color: #fff;
  border: none;
  font-size: 1rem;
  box-shadow: 0 1px 4px #90caf9;
  cursor: pointer;
  transition: background 0.2s;
}
button:hover {
  background: #1565c0;
}
</style>
