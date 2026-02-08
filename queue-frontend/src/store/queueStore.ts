import { ref } from 'vue'

function getOrCreateClientId() {
  let id = localStorage.getItem('client_id')
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem('client_id', id)
  }
  return id
}
let clientId = getOrCreateClientId()

const currentQueue = ref('A0') // global queue (from WS)
const yourQueue = ref('A0')    // local queue (from REST)
let updateCallbacks: ((q: string) => void)[] = []
let wsInstance: WebSocket | null = null

function clearQueueWS() {
  if (wsInstance && wsInstance.readyState === WebSocket.OPEN) {
    wsInstance.send('clear')
  }
}

function connectWs() {
  const wsUrl = 'ws://' + window.location.hostname + ':8080/ws'
  const ws = new WebSocket(wsUrl)
  wsInstance = ws
  console.log('WS connecting:', wsUrl)

  ws.onopen = () => {
    console.log('WS connected')
  }

  ws.onmessage = (evt) => {
    console.log('WS message:', evt.data)
    try {
      const data = JSON.parse(evt.data)
      if (data && data.event === 'clear') {
        localStorage.removeItem('client_id')
        clientId = getOrCreateClientId()
        console.log('WS clear event: localStorage cleared & queue refreshed')
      }
      if (data && data.event === 'queue') {
        ws.send('get')
      }
      if (data && typeof data.queue === 'string') {
        currentQueue.value = data.queue
      }
    } catch (e) {
      console.warn('ws message parse error', e)
    }
  }

  ws.onerror = (err) => {
    console.error('WS error:', err)
  }

  ws.onclose = () => {
    console.warn('WS closed, reconnecting...')
    setTimeout(connectWs, 1000)
  }
}

async function fetchQueue() {
  console.log('fetchQueue called')
  const res = await fetch('http://' + window.location.hostname + ':8080/queue', {
    headers: { 'x-client-id': clientId }
  })
  const data = await res.json()
  if (typeof data.queue === 'string') {
    yourQueue.value = data.queue
    updateCallbacks.forEach(cb => cb(yourQueue.value))
  }
}

async function nextQueue() {
  const res = await fetch('http://' + window.location.hostname + ':8080/queue/next', {
    method: 'POST',
    headers: { 'x-client-id': clientId }
  })
  const data = await res.json()
  if (typeof data.queue === 'string') {
    yourQueue.value = data.queue
    updateCallbacks.forEach(cb => cb(yourQueue.value))
  }
}

async function clearQueue() {
  const res = await fetch('http://' + window.location.hostname + ':8080/queue/clear', {
    method: 'POST',
    headers: { 'x-client-id': clientId }
  })
  const data = await res.json()
  if (typeof data.queue === 'string') {
    yourQueue.value = data.queue
    updateCallbacks.forEach(cb => cb(yourQueue.value))
  }
}

function onUpdate(cb: (q: string) => void) {
  updateCallbacks.push(cb)
}

connectWs()

export function useQueueStore(options?: { skipAutoFetch?: boolean }) {
  if (!options?.skipAutoFetch) {
    fetchQueue()
  }
  return { currentQueue, yourQueue, nextQueue, clearQueue, clearQueueWS, onUpdate, fetchQueue }
}
