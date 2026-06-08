<script setup lang="ts">
import { ref } from "vue"

const API = "/api"
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ "update:modelValue": [value: boolean] }>()

const tab = ref<"upload" | "url">("upload")
const urlInput = ref("")
const uploading = ref(false)
const previewUrl = ref("")
const fileInput = ref<HTMLInputElement | null>(null)
const activeMethod = ref<"upload" | "url" | null>(null)

async function onFileSelected(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return
  const file = input.files[0]
  if (file.size > 2 * 1024 * 1024) { alert("剧剧墖大小不能超过 2MB"); return }
  uploading.value = true
  const formData = new FormData()
  formData.append("image", file)
  try {
    const res = await fetch("/api/notes/upload", { method: "POST", body: formData })
    const data = await res.json()
    if (data.success) {
      previewUrl.value = data.url
      urlInput.value = ""
      activeMethod.value = "upload"
    } else alert(data.error || "上传失败")
  } catch { alert("上传失败") }
  uploading.value = false
}

function useUrl() {
  const val = urlInput.value.trim()
  if (!val) return
  previewUrl.value = val
  activeMethod.value = "url"
}

function onTabChange(newTab: "upload" | "url") {
  if (newTab === "upload") urlInput.value = ""
  tab.value = newTab
}

function clearFavicon() {
  previewUrl.value = ""
  urlInput.value = ""
  activeMethod.value = null
}

async function save() {
  try {
    await fetch(API + "/settings", {
      method: "POST", headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ key: "site_favicon", value: previewUrl.value })
    })
    updateFavicon(previewUrl.value)
  } catch {}
  emit("update:modelValue", false)
}

function updateFavicon(url: string) {
  if (!url) return
  let link = document.querySelector('link[rel~=icon]')
  if (!link) { link = document.createElement('link'); link.rel = 'icon'; document.head.appendChild(link) }
  link.href = url
}

function isImage(val: string) {
  return val.startsWith("/uploads/") || val.startsWith("http")
}
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)"
    max-width="420" scrollable persistent transition="dialog-bottom-transition">
    <v-card class="rounded-xl" rounded="xl">
      <v-card-title class="d-flex align-center pa-4 pb-0">
        <v-icon start color="primary">mdi-image-multiple</v-icon>
        <span class="text-h6 font-weight-medium">设置网页片标</span>
        <v-spacer />
        <v-btn icon="mdi-close" size="small" variant="text" @click.stop="emit('update:modelValue', false)" />
      </v-card-title>

      <!-- Preview -->
      <div class="d-flex justify-center pa-4 pb-2">
        <div style="position:relative">
          <v-img v-if="isImage(previewUrl)" :src="previewUrl" width="80" height="80" class="rounded-lg" />
          <v-icon v-else size="80" color="primary" class="icon-placeholder">mdi-web</v-icon>
          <v-chip v-if="activeMethod" size="x-small" variant="flat" color="primary" class="method-badge"
            style="position:absolute;bottom:-4px;right:-4px">
            {{ activeMethod === 'upload' ? "已上传" : "链接" }}
          </v-chip>
        </div>
      </div>

      <v-card-text class="pa-4 pt-2">
        <v-tabs v-model="tab" color="primary" density="compact" class="mb-4" @update:model-value="onTabChange">
          <v-tab value="upload">
            <v-icon start size="small">mdi-cloud-upload</v-icon>
            上传片片
          </v-tab>
          <v-tab value="url">
            <v-icon start size="small">mdi-link</v-icon>
            片片链接
          </v-tab>
        </v-tabs>

        <!-- Upload tab -->
        <div v-if="tab === 'upload'" class="text-center">
          <input ref="fileInput" type="file" accept="image/*" hidden @change="onFileSelected" />
          <v-card variant="outlined" class="upload-zone rounded-lg pa-6 mb-3" @click="fileInput?.click()">
            <v-icon size="48" color="primary" class="mb-2" style="opacity: 0.5">mdi-cloud-upload-outline</v-icon>
            <p class="text-body-2 font-weight-medium">点击选择片片</p>
            <p class="text-caption text-medium-emphasis mt-1">JPG / PNG / GIF / WebP，最大2MB</p>
          </v-card>
          <v-progress-linear v-if="uploading" indeterminate color="primary" class="rounded-pill" />
          <p v-if="activeMethod === 'url'" class="text-caption text-warning mt-2">切换到该方式将替换链接片标</p>
        </div>

        <!-- URL tab -->
        <div v-if="tab === 'url'">
          <p class="text-body-2 text-medium-emphasis mb-3">输入片片链接，支持任意网络片片地址</p>
          <div class="d-flex ga-2 mb-3">
            <v-text-field v-model="urlInput" variant="outlined" density="compact" hide-details
              placeholder="https://example.com/favicon.png" class="flex-1" @keyup.enter="useUrl" />
          </div>
          <v-btn variant="tonal" color="primary" class="rounded-pill" block @click="useUrl">预览</v-btn>
          <p v-if="activeMethod === 'upload'" class="text-caption text-warning mt-2">切换到该方式将替换上传的片标</p>
        </div>
      </v-card-text>

      <v-card-actions class="pa-4 pt-0 d-flex ga-2">
        <v-btn variant="text" class="flex-1 rounded-pill" @click="emit('update:modelValue', false)">取消</v-btn>
        <v-btn v-if="previewUrl"
          variant="text" color="error" class="rounded-pill" @click="clearFavicon">清除</v-btn>
        <v-btn variant="flat" color="primary" class="flex-1 rounded-pill" @click="save">确认</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.flex-1 { flex: 1; }
.upload-zone { cursor: pointer; transition: all 0.2s; border: 2px dashed rgba(var(--v-theme-primary), 0.3) !important; }
.upload-zone:hover { border-color: rgb(var(--v-theme-primary)) !important; background: rgba(var(--v-theme-primary), 0.04); }

@media (max-width: 768px) {
  .picker-dialog :deep(.v-card) { margin: 12px; border-radius: 16px !important; }
}</style>








