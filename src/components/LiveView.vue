<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue"
import ArtPlayer from "artplayer"

const API = "/api"
const streamUrl = ref("")
const loading = ref(true)
const error = ref("")
let player: ArtPlayer | null = null

onMounted(async () => {
  try {
    const r = await fetch(`${API}/live/config`)
    if (r.ok) {
      const data = await r.json()
      streamUrl.value = data.streamUrl || ""
    }
  } catch {
    error.value = "无法加载直播配置"
  }
  loading.value = false

  await nextTick()
  initPlayer()
})

onUnmounted(() => {
  if (player) {
    player.destroy()
    player = null
  }
})

async function initPlayer() {
  if (!streamUrl.value) return
  if (player) return

  const Hls = (await import("hls.js")).default

  player = new ArtPlayer({
    container: "#artplayer-container",
    url: streamUrl.value,
    volume: 1,
    isLive: true,
    autoSize: false,
    autoMini: false,
    screenshot: false,
    setting: false,
    playbackRate: false,
    pip: false,
    mutex: true,
    fullscreen: true,
    fullscreenWeb: true,
    customType: {
      m3u8: function (video: HTMLVideoElement, url: string) {
        if (Hls.isSupported()) {
          const hls = new Hls()
          hls.loadSource(url)
          hls.attachMedia(video)
        } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
          video.src = url
        }
      },
    },
  })

  // Enter fullscreen on first play
  player.on("play", () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen?.()
    }
  })
}
</script>

<template>
  <div class="live-page">
    <div v-if="loading" class="live-status">
      <p class="text-body-2 text-medium-emphasis mt-2">加载直播...</p>
    </div>
    <div v-else-if="!streamUrl" class="live-status">
      <v-icon size="48" color="rgba(255,255,255,0.3)" class="mb-2">mdi-video-off-outline</v-icon>
      <p class="text-body-2 text-medium-emphasis mt-3">直播流未配置</p>
      <p class="text-caption text-medium-emphasis mt-1">请管理员在后台配置直播流地址</p>
    </div>
    <div v-else-if="error" class="live-status">
      <p class="text-body-2 text-medium-emphasis">{{ error }}</p>
    </div>
    <div v-else id="artplayer-container" class="player-container"></div>
  </div>
</template>

<style scoped>
.live-page {
  width: 100vw;
  height: 100vh;
  background: #000;
  overflow: hidden;
}
.live-status {
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
}
.player-container {
  width: 100%;
  height: 100%;
}
:deep(.art-video-player) {
  width: 100% !important;
  height: 100% !important;
}
:deep(.art-video-player video) {
  width: 100% !important;
  height: 100% !important;
  object-fit: contain !important;
}
/* Remove ArtPlayer internal padding so controls span full width */
:deep(.art-controls) {
  width: 100% !important;
  padding: 0 !important;
}
:deep(.art-bottom) {
  padding: 0 !important;
}
</style>
