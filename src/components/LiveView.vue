<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue"

const streamUrl = ref("")
const loading = ref(true)
let dp: any = null

onMounted(async () => {
  // Load DPlayer from CDN
  const link = document.createElement("link")
  link.rel = "stylesheet"
  link.href = "https://cdn.jsdelivr.net/npm/dplayer@1.26.0/dist/DPlayer.min.css"
  document.head.appendChild(link)

  const script = document.createElement("script")
  script.src = "https://cdn.jsdelivr.net/npm/dplayer@1.26.0/dist/DPlayer.min.js"
  script.onload = initPlayer
  document.head.appendChild(script)

  try {
    const r = await fetch("/api/live/config")
    const d = await r.json()
    streamUrl.value = d.streamUrl || ""
  } catch {}
  loading.value = false
})

async function initPlayer() {
  const Hls = (window as any).Hls
  const DPlayer = (window as any).DPlayer
  if (!DPlayer) return
  await nextTick()
  dp = new DPlayer({
    container: document.getElementById("dplayer"),
    video: { url: streamUrl.value, type: streamUrl.value ? "customHLS" : "" },
    autoplay: true,
    volume: 1,
    theme: "#1976D2",
  })
  if (Hls && streamUrl.value) {
    // DPlayer customType handles HLS via hls.js
    dp.video = { ...dp.video, type: "customHLS", customType: {
      customHLS: (v: HTMLVideoElement) => {
        const h = new Hls()
        h.loadSource(v.src)
        h.attachMedia(v)
      }
    }}
    dp.switchVideo({ url: streamUrl.value, type: "customHLS" })
  }
}

onBeforeUnmount(() => {
  if (dp) { try { dp.destroy() } catch {} }
})
</script>

<template>
  <div class="live-page">
    <div v-if="loading" class="loading-state">
      <v-progress-circular indeterminate color="primary" />
    </div>
    <div v-else-if="!streamUrl" class="loading-state">
      <v-icon size="48" color="grey-darken-1">mdi-video-off-outline</v-icon>
      <p class="text-body-2 text-medium-emphasis mt-3">直播流未配置</p>
    </div>
    <div v-else id="dplayer"></div>
  </div>
</template>

<style scoped>
.live-page { width: 100vw; height: 100vh; background: #000; overflow: hidden; }
.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: 8px; }
#dplayer { width: 100%; height: 100%; }
</style>
