<script setup lang="ts">
import { ref } from "vue"
import { useAuthStore } from "@/stores/auth"

const auth = useAuthStore()
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ "update:modelValue": [value: boolean] }>()

const tab = ref<"upload" | "url">("upload")
const urlInput = ref("")
const uploading = ref(false)
const previewUrl = ref(auth.userAppIcon || "")
const fileInput = ref<HTMLInputElement | null>(null)
const activeMethod = ref<"upload" | "url" | null>(null)

async function onFileSelected(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return
  const file = input.files[0]
  if (file.size > 5 * 1024 * 1024) { alert("图片大小不能超过 5MB"); return }
  uploading.value = true
  const formData = new FormData()
  formData.append("avatar", file)
  try {
    const res = await fetch("/api/auth/avatar/upload", { method: "POST", body: formData })
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

function removeIcon() {
  previewUrl.value = ""
  urlInput.value = ""
  activeMethod.value = null
}

async function save() {
  await auth.updateAppIcon(previewUrl.value)
  emit('update:modelValue', false)
}

function isImage(val: string) {
  return val.startsWith("/uploads/") || val.startsWith("http")
}
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', )"
    max-width="420" scrollable persistent transition="dialog-bottom-transition">
    <v-card class="rounded-xl" rounded="xl">
      <v-card-title class="d-flex align-center pa-4 pb-0">
        <v-icon start color="primary">mdi-image-outline</v-icon>
        <span class="text-h6 font-weight-medium">设置工具栏图标</span>
        <v-spacer />
        <v-btn icon="mdi-close" size="small" variant="text" @click.stop="emit('update:modelValue', false)" />
      </v-card-title>

      <!-- Preview -->
      <div class="d-flex justify-center pa-4 pb-2">
        <div style="position:relative">
          <v-img v-if="isImage(previewUrl)" :src="previewUrl" width="80" height="80" class="rounded-lg" />
          <v-icon v-else size="80" color="primary" class="icon-placeholder">mdi-pencil-box-multiple</v-icon>
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
            上传图片
          </v-tab>
          <v-tab value="url">
            <v-icon start size="small">mdi-link</v-icon>
            图片链接
          </v-tab>
        </v-tabs>

        <!-- Upload tab -->
        <div v-if="tab === 'upload'" class="text-center">
          <input ref="fileInput" type="file" accept="image/*" hidden @change="onFileSelected" />
          <v-card variant="outlined" class="upload-zone rounded-lg pa-6 mb-3" @click="fileInput?.click()">
            <v-icon size="48" color="primary" class="mb-2" style="opacity: 0.5">mdi-cloud-upload-outline</v-icon>
            <p class="text-body-2 font-weight-medium">点击选择图片</p>
            <p class="text-caption text-medium-emphasis mt-1">JPG / PNG / GIF / WebP，最大5MB</p>
          </v-card>
          <v-progress-linear v-if="uploading" indeterminate color="primary" class="rounded-pill" />
          <p v-if="activeMethod === 'url'" class="text-caption text-warning mt-2">切换到该方式将替换链接图标</p>
        </div>

        <!-- URL tab -->
        <div v-if="tab === 'url'">
          <p class="text-body-2 text-medium-emphasis mb-3">输入图片链接，支持任意网络图片地址</p>
          <div class="d-flex ga-2 mb-3">
            <v-text-field v-model="urlInput" variant="outlined" density="compact" hide-details
              placeholder="https://example.com/icon.png" class="flex-1" @keyup.enter="useUrl" />
          </div>
          <v-btn variant="tonal" color="primary" class="rounded-pill" block @click="useUrl">预览</v-btn>
          <p v-if="activeMethod === 'upload'" class="text-caption text-warning mt-2">切换到该方式将替换上传的图标</p>
        </div>
      </v-card-text>

      <v-card-actions class="pa-4 pt-0 d-flex ga-2">
        <v-btn variant="text" class="flex-1 rounded-pill" @click="emit('update:modelValue', false)">取消</v-btn>
        <v-btn v-if="previewUrl && previewUrl !== auth.userAppIcon"
          variant="text" color="error" class="rounded-pill" @click="removeIcon">清除</v-btn>
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
}
</style>
