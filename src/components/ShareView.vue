<script setup lang="ts">
import { ref, onMounted } from "vue"
import type { Note } from "@/stores/notes"
import MarkdownPreview from "./MarkdownPreview.vue"

const note = ref<Note | null>(null)
const loading = ref(true)
const error = ref("")

async function loadSharedNote() {
  const token = window.location.pathname.replace("/share/", "")
  if (!token) {
    error.value = "无效的分享链接"
    loading.value = false
    return
  }
  try {
    const res = await fetch(`/api/share/${token}`)
    if (res.ok) {
      note.value = await res.json()
    } else {
      const data = await res.json()
      error.value = data.error || "笔记不存在或分享链接已失效"
    }
  } catch {
    error.value = "无法加载笔记"
  }
  loading.value = false
}

onMounted(loadSharedNote)

function timeAgo(ts: number) {
  const diff = Date.now() - ts
  const seconds = Math.floor(diff / 1000)
  if (seconds < 60) return "刚刚"
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes} 分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days} 天前`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months} 个月前`
  return `${Math.floor(months / 12)} 年前`
}

function displayName(memo: Note) {
  return memo.nickname?.trim() || memo.username || "匿名"
}

function isImage(val?: string) {
  return val?.startsWith("/uploads/") || val?.startsWith("http")
}
</script>

<template>
  <div class="share-page">
    <div class="share-header">
      <a href="/" class="share-home-link">
        <svg width="20" height="20" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" fill="currentColor">
          <path d="M659.655431 521.588015q23.970037-6.71161 46.022472-13.423221 19.17603-5.752809 39.310861-11.505618t33.558052-10.546816l-13.423221 50.816479q-5.752809 21.093633-10.546816 31.640449-9.588015 25.88764-22.531835 47.940075t-24.449438 38.35206q-13.423221 19.17603-27.805243 35.475655l-117.932584 35.475655 96.838951 17.258427q-19.17603 16.299625-41.228464 33.558052-19.17603 14.382022-43.625468 30.202247-51.29588 29.243446-59.925094 13.902622-12.821397-55.131086-17.258427-101.632958t-1.917603-79.580525q2.876405-34.516854 17.258427-86.771537t41.707865-115.535581q14.382022-35.475655 22.531835-56.569288-23.490637 28.764045-46.981274 55.610486t-30.681647 36.913857q-11.026217 13.902622-9.588015 6.71161-4.314607-15.340824-46.022472-74.307116-52.734084-73.827716-87.730339-107.865169t-100.674157-64.718927q13.423221 33.558052 21.093633 65.677729t9.588015 52.254682q1.438202 14.861424-0.4794 29.722846-13.423221-18.217228-45.06367-39.310861-22.531835-15.340824-35.955056-16.299625t-40.749064 13.423221q-23.011236 12.342022-53.213484 48.419476t-56.569288 73.348316q-29.722846 40.269663-63.760299 99.236955t-60.883895 115.056181q-29.722846 58.965093-36.434457 83.894932t8.149813 36.434457h0.958801q59.445693-16.299625 132.314607-60.404495t132.314607-103.550562q56.569288-56.569288 109.303371-129.917603t78.142322-144.778652q-3.835206 33.558052-22.531835 97.318353t-45.06367 117.932584-51.775281 111.220973-42.666667 80.539326q21.093633 1.917603 42.187266-0.958801 14.861424-1.917603 28.285246-1.917603t29.722846 2.876404 27.325843 6.71161l18.217228 4.314607q-3.835206 11.026217-12.821397 35.475655t-20.134831 54.172284-20.614232 59.445693-9.588015 43.625468q3.835206 1.438202 23.011236-0.479401 22.531835-2.876405 56.569288-13.423221t71.430712-25.40824q30.202247-11.505618 48.898876-19.65543l8.149813 19.17603q-7.191011 12.342022-36.434457 43.625468-27.805243 29.722846-47.460674 45.06367t-30.202247 19.17603q3.835206 9.588015 32.11985 21.093633 21.093633 8.629213 36.434457 11.026217-6.23221 11.026217-28.284644 25.40824-14.382022 8.149813-23.011236 11.505618 11.026217 15.340824 35.475655 27.805243t52.734084 20.614232q28.764045 8.149813 54.172284 8.149813t43.146067-5.273409q13.423221-10.067416 7.670412-16.299625z" />
        </svg>
        <span>碎碎 SuiSui</span>
      </a>
    </div>

    <main class="share-main">
      <div v-if="loading" class="d-flex justify-center py-16">
        <v-progress-circular indeterminate color="primary" />
      </div>
      <div v-else-if="error" class="share-error">
        <v-icon size="48" color="error" class="mb-3">mdi-link-variant-off</v-icon>
        <h2>分享链接无效</h2>
        <p class="text-medium-emphasis">{{ error }}</p>
        <v-btn variant="flat" color="primary" href="/" class="mt-4 rounded-pill">返回首页</v-btn>
      </div>
      <div v-else-if="note" class="share-note-card">
        <div class="share-note-header">
          <div class="d-flex align-center ga-2">
            <div v-if="isImage(note.avatar)" class="share-avatar">
              <img :src="note.avatar" alt="" width="32" height="32" style="border-radius:8px;object-fit:cover" />
            </div>
            <div v-else class="share-avatar-fallback">{{ displayName(note).charAt(0).toUpperCase() }}</div>
            <div>
              <div class="share-author-name">{{ displayName(note) }}</div>
              <div class="share-time">{{ timeAgo(note.createdAt) }}</div>
            </div>
          </div>
        </div>
        <div class="share-content">
          <MarkdownPreview :content="note.content" />
        </div>
        <div v-if="note.tags && note.tags.length" class="d-flex flex-wrap ga-1 mt-3">
          <v-chip v-for="tag in note.tags" :key="tag" size="x-small" variant="tonal" color="primary">
            #{{ tag }}
          </v-chip>
        </div>
      </div>
    </main>

    <footer class="share-footer">
      <span>来自 <a href="/">碎碎 SuiSui</a> 的分享</span>
    </footer>
  </div>
</template>

<style scoped>
.share-page {
  min-height: 100vh;
  background: rgb(var(--v-theme-background));
  display: flex;
  flex-direction: column;
}
.share-header {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}
.share-home-link {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: rgb(var(--v-theme-on-surface));
  font-weight: 600;
  font-size: 0.95rem;
  opacity: 0.7;
  transition: opacity 0.2s;
}
.share-home-link:hover { opacity: 1; }
.share-main {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 40px 16px;
}
.share-error {
  text-align: center;
  padding: 48px 16px;
}
.share-error h2 { font-size: 1.2rem; margin-bottom: 8px; }
.share-note-card {
  width: 100%;
  max-width: 680px;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 16px;
  padding: 24px;
}
.share-note-header {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}
.share-avatar-fallback {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgb(var(--v-theme-primary));
  color: #fff;
  font-size: 0.8rem;
  font-weight: 600;
}
.share-author-name { font-weight: 600; font-size: 0.95rem; }
.share-time { font-size: 0.75rem; color: rgba(var(--v-theme-on-surface), 0.45); }
.share-content {
  line-height: 1.7;
  font-size: 0.95rem;
}
.share-footer {
  text-align: center;
  padding: 16px;
  font-size: 0.75rem;
  color: rgba(var(--v-theme-on-surface), 0.35);
}
.share-footer a { color: inherit; }
</style>
