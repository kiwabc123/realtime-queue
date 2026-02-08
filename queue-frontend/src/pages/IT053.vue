<template>
  <div class="container">
    <h2>ล้างคิว</h2>
    <div class="queue-display">{{ queueNumber }}</div>
    <button @click="clearQueue">ล้างคิว</button>
    <button @click="goBack">กลับหน้ารับบัตรคิว</button>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQueueStore } from '../store/queueStore'
useQueueStore({ skipAutoFetch: true })
const router = useRouter()
const queueNumber = ref('A0')
let ws: WebSocket | null = null

function connectWs() {
  const wsUrl = 'ws://' + window.location.hostname + ':8080/ws'
  ws = new WebSocket(wsUrl)
  ws.onopen = () => {
    ws?.send('get')
  }
  ws.onmessage = (evt) => {
    try {
      const data = JSON.parse(evt.data)
      if (typeof data.queue === 'string') {
        queueNumber.value = data.queue
      }
    } catch {}
  }
  ws.onclose = () => {
    setTimeout(connectWs, 1000)
  }
}

function clearQueue() {
  ws?.send('clear')
}
function goBack() {
  router.push('/')
}
onMounted(connectWs)
</script>
<style scoped>
.container { max-width: 400px; margin: 40px auto; padding: 24px; border-radius: 8px; box-shadow: 0 2px 8px #ccc; background: #fff; }
.queue-display { font-size: 2rem; margin: 16px 0; }
button { margin: 8px; padding: 8px 16px; }
</style>
