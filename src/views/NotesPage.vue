<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from "vue"
import { useDisplay } from "vuetify"
import { useNotesStore } from "@/stores/notes"
import { useAuthStore } from "@/stores/auth"
import { useSettingsStore } from "@/stores/settings"
import NoteCard from "@/components/NoteCard.vue"
import Heatmap from "@/components/Heatmap.vue"
import SidePanel from "@/components/SidePanel.vue"
import ZoomOverlay from "@/components/ZoomOverlay.vue"

const props = defineProps<{ mobileHeatmap: boolean }>()
const emit = defineEmits<{ "close-heatmap": [] }>()

const store = useNotesStore()
const auth = useAuthStore()
const settings = useSettingsStore()
const { mobile } = useDisplay()
const searchQuery = ref("")
const selectedTag = ref("")

onMounted(async () => { await store.fetchNotes(); await settings.load() })

const filteredNotes = computed(() => {
  let list = store.notes
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(n => n.content.toLowerCase().includes(q) || n.tags.some(t => t.toLowerCase().includes(q)))
  }
  if (selectedTag.value) list = list.filter(n => n.tags.includes(selectedTag.value))
  return list
})

const inlineContent = ref("")
const inlineTagsInput = ref("")
const showInlineTags = ref(false)
const inlineUploading = ref(false)
const inlineTextarea = ref<HTMLTextAreaElement | null>(null)
const inlineFileInput = ref<HTMLInputElement | null>(null)
const uploadedImages = ref<string[]>([])
const editingNoteId = ref("")
const zoomedUpload = ref("")

function onInlineKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) submitInline()
}

async function submitInline() {
  if ((!inlineContent.value.trim() && !uploadedImages.value.length) || !auth.isLoggedIn) return
  const tags = inlineTagsInput.value.split(/[,，]/).map(t => t.trim()).filter(Boolean)
  let content = inlineContent.value
  for (const url of uploadedImages.value) content += "\n\n![](" + url + ")"
  if (editingNoteId.value) {
    await store.updateNote(editingNoteId.value, content.trim(), tags)
    editingNoteId.value = ""
  } else {
    await store.addNote(content.trim(), tags, auth.userName)
  }
  inlineContent.value = ""
  inlineTagsInput.value = ""
  uploadedImages.value = []
  showInlineTags.value = false
  nextTick(() => {
    const el = document.querySelector(".inline-textarea") as HTMLTextAreaElement
    if (el) el.style.height = ""
  })
}

function triggerInlineUpload() { inlineFileInput.value?.click() }

async function onInlineUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files || [])
  if (!files.length) return
  if (files.some(f => f.size > 10 * 1024 * 1024)) { alert("图片大小不能超过 10MB"); input.value = ""; return }
  inlineUploading.value = true
  for (const file of files) {
    const fd = new FormData()
    fd.append("image", file)
    try {
      const res = await fetch("/api/notes/upload", { method: "POST", body: fd })
      const data = await res.json()
      if (data.success) uploadedImages.value.push(data.url)
      else alert(data.error || "上传失败")
    } catch { alert("上传失败") }
  }
  inlineUploading.value = false
  input.value = ""
}

function autoGrowTextarea(e: Event) {
  const el = e.target as HTMLTextAreaElement
  el.style.height = "auto"
  el.style.height = el.scrollHeight + "px"
}

function handleEdit(memo: any) {
  const imgRegex = /!\[.*?\]\((.+?)\)/g
  const urls: string[] = []
  const text = memo.content.replace(imgRegex, (_m: string, url: string) => { urls.push(url); return "" })
  inlineContent.value = text.trim()
  uploadedImages.value = urls
  editingNoteId.value = memo.id
  showInlineTags.value = false
  nextTick(() => {
    const el = document.querySelector(".inline-textarea") as HTMLTextAreaElement
    if (el) { el.style.height = "auto"; el.style.height = el.scrollHeight + "px" }
  })
  document.querySelector(".inline-textarea")?.scrollIntoView({ behavior: "smooth" })
  ;(document.querySelector(".inline-textarea") as HTMLTextAreaElement)?.focus()
}
</script>

<template>
  <div class="notes-layout" :class="{ mobile: mobile }">
    <div class="side-col">
      <SidePanel v-model:search="searchQuery" v-model:tag="selectedTag" />
    </div>

    <div class="main-col">
      <v-dialog :model-value="mobileHeatmap" max-width="400" scrollable persistent transition="dialog-bottom-transition"
        @update:model-value="v => !v && emit('close-heatmap')">
        <v-card class="rounded-xl pa-4">
          <div class="d-flex align-center mb-3">
            <span class="font-weight-medium">活动日历</span>
            <v-spacer /><v-btn icon="mdi-close" size="small" variant="text" @click="emit('close-heatmap')" />
          </div>
          <SidePanel v-model:search="searchQuery" v-model:tag="selectedTag"
            :emitOnTag="true" @tag-click="emit('close-heatmap')" />
        </v-card>
      </v-dialog>

      <div v-if="auth.isLoggedIn" class="inline-editor mb-4">
        <div class="editor-box">
          <textarea ref="inlineTextarea" v-model="inlineContent" class="inline-textarea"
            placeholder="写点什么呢.." rows="1" @keydown="onInlineKeydown" @input="autoGrowTextarea"></textarea>
          <div v-if="uploadedImages.length" class="upload-preview-row">
            <div v-for="(img, ii) in uploadedImages" :key="ii" class="upload-preview-item">
              <img :src="img" class="upload-preview-img" @click.stop="zoomedUpload = img" />
              <button class="upload-remove-btn" @click="uploadedImages.splice(ii, 1)">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
          </div>
          <div class="editor-toolbar">
            <div class="d-flex align-center ga-1">
              <v-btn icon="mdi-image-plus" size="x-small" variant="text" class="tool-btn" :loading="inlineUploading" @click="triggerInlineUpload" />
              <input ref="inlineFileInput" type="file" accept="image/*" multiple hidden @change="onInlineUpload" />
              <v-btn :icon="showInlineTags ? 'mdi-tag-off' : 'mdi-tag-outline'" size="x-small" variant="text" class="tool-btn" @click="showInlineTags = !showInlineTags" />
            </div>
            <v-btn color="#1976D2" size="small" variant="flat" class="rounded-pill px-4 submit-btn" @click="submitInline">
              <v-icon start size="x-small">mdi-send</v-icon>{{ editingNoteId ? "更新" : "发布" }}
            </v-btn>
          </div>
          <v-expand-transition>
            <div v-if="showInlineTags" class="pa-3">
              <v-text-field v-model="inlineTagsInput" label="标签（逗号分隔）" variant="outlined" hide-details density="compact" placeholder="vue, memos, md" />
            </div>
          </v-expand-transition>
        </div>
      </div>

      <div v-if="!store.loaded" class="d-flex justify-center py-16">
        <v-progress-circular indeterminate color="primary" />
      </div>
      <template v-else>
        <div v-if="filteredNotes.length === 0" class="empty-state">
          <p class="empty-text">{{ (searchQuery || selectedTag) ? "没有找到匹配的备忘录" : "还没有备忘录" }}</p>
        </div>
        <div class="d-flex flex-column ga-3">
          <NoteCard v-for="note in filteredNotes" :key="note.id" :memo="note" :logged-in="auth.isLoggedIn" @edit="handleEdit" />
        </div>
      </template>
      <div v-if="settings.siteIcp" class="icp-text">{{ settings.siteIcp }}</div>
    </div>
  </div>

  <ZoomOverlay v-if="zoomedUpload" :src="zoomedUpload" @close="zoomedUpload = ''" />
</template>

<style scoped>
.notes-layout { display: flex; gap: 24px; padding: 24px; max-width: 1200px; margin: 0 auto; align-items: flex-start; }
.notes-layout.mobile { flex-direction: column; padding: 12px; gap: 12px; }
.side-col { width: 280px; flex-shrink: 0; position: sticky; top: 24px; align-self: flex-start; }
.notes-layout.mobile .side-col { display: none; }
.main-col { flex: 1; min-width: 0; }

.inline-editor { width: 100%; }
.editor-box {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 12px; overflow: hidden;
  transition: border-color 0.2s, box-shadow 0.2s;
  background: rgb(var(--v-theme-surface));
}
.editor-box:focus-within {
  border-color: rgba(var(--v-theme-primary), 0.3);
  box-shadow: 0 2px 12px rgba(var(--v-theme-primary), 0.06);
}
.inline-textarea {
  width: 100%; border: none; outline: none; resize: none;
  padding: 14px 16px 8px; font-size: 0.95rem; line-height: 1.6;
  font-family: inherit; background: transparent;
  color: rgb(var(--v-theme-on-surface)); min-height: 80px;
}
.inline-textarea::placeholder { color: rgba(var(--v-theme-on-surface), 0.35); }
.editor-toolbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 4px 8px 8px;
}
.tool-btn { opacity: 0.5; transition: opacity 0.2s; }
.tool-btn:hover { opacity: 1; }
.submit-btn { height: 30px; }

.upload-preview-row { display: flex; flex-wrap: wrap; gap: 8px; padding: 8px 8px 4px; }
.upload-preview-item {
  position: relative; display: inline-block;
  width: 72px; height: 72px; border-radius: 8px; overflow: hidden;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08); flex-shrink: 0;
}
.upload-preview-img { width: 100%; height: 100%; object-fit: cover; cursor: zoom-in; }
.upload-remove-btn {
  position: absolute; top: -4px; right: -4px; width: 20px; height: 20px;
  border-radius: 50%; border: none; background: rgb(var(--v-theme-surface));
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; color: rgba(var(--v-theme-on-surface), 0.6);
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
}

.empty-state {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 64px 0; color: rgba(var(--v-theme-on-surface), 0.5);
}
.empty-text { font-size: 1rem; font-weight: 500; }

.icp-text {
  text-align: center; font-size: 0.75rem; padding: 16px 0; opacity: 0.6;
  color: rgba(var(--v-theme-on-surface), 0.6);
}

@media (max-width: 768px) {
  .notes-layout.mobile { flex-direction: column; padding: 12px; gap: 8px; }
  .notes-layout.mobile .main-col { width: 100%; }
  .inline-textarea { min-height: 60px; padding: 12px 14px 8px; font-size: 0.9rem; }
}
</style>





