<script setup lang="ts">
import { ref, onMounted, watch } from "vue"
import { useAuthStore } from "@/stores/auth"
import { useSettingsStore } from "@/stores/settings"
import AvatarPicker from "@/components/AvatarPicker.vue"
import AppIconPicker from "@/components/AppIconPicker.vue"
import FaviconPicker from "@/components/FaviconPicker.vue"

const API = "/api"
const auth = useAuthStore()
const settings = useSettingsStore()
const emit = defineEmits<{ back: [] }>()

const tab = ref("overview")
const stats = ref<null | any>(null)
const users = ref<any[]>([])
const loading = ref(false)
const deleting = ref<null | number>(null)
const nickError = ref("")
const oldPwd = ref("")
const newPwd = ref("")
const snackbar = ref(false)
const snackMsg = ref("")
const nickInput = ref(auth.userNickname)
const showAvatarPicker = ref(false)
const showAppIconPicker = ref(false)
const showFaviconPicker = ref(false)
const allowRegister = ref(true)
const siteTitle = ref("")
const siteIcp = ref("")

onMounted(() => {
  if (auth.userRole !== "admin") tab.value = "profile"
  else { loadData(); loadSettings() }
  nickInput.value = auth.userNickname
})

watch(() => auth.userRole, (val) => { if (val !== "admin") tab.value = "profile" })

async function loadData() {
  loading.value = true
  await Promise.all([loadStats(), loadUsers()])
  loading.value = false
}
async function loadStats() {
  try { const r = await fetch(API + "/admin/stats"); if (r.ok) stats.value = await r.json() } catch {}
}
async function loadUsers() {
  try { const r = await fetch(API + "/admin/users"); if (r.ok) users.value = await r.json() } catch {}
}
async function loadSettings() {
  await settings.load()
  siteTitle.value = settings.siteTitle
  siteIcp.value = settings.siteIcp
  document.title = settings.siteTitle || "Mengji"
  allowRegister.value = settings.allowRegister
}
async function saveSiteTitle() {
  await settings.save("site_title", siteTitle.value.trim())
  siteTitle.value = siteTitle.value.trim()
  settings.applyTitle()
  snackMsg.value = "???????"; snackbar.value = true
}
async function saveSiteIcp() {
  await settings.save("site_icp", siteIcp.value.trim())
  siteIcp.value = siteIcp.value.trim()
  snackMsg.value = "??????"; snackbar.value = true
}
async function toggleRegister(val: boolean) {
  await settings.save("allow_register", val ? "true" : "false")
  settings.allowRegister = val
  snackMsg.value = val ? "?????" : "?????"; snackbar.value = true
}
async function deleteUser(id: number) {
  if (!confirm("ȷ��ɾ����")) return
  deleting.value = id
  try { await fetch(API + "/admin/users/" + id, { method: "DELETE" }); await loadData() } catch {}
  deleting.value = null
}
async function saveNickname() {
  nickError.value = ""
  if (!nickInput.value.trim()) return
  const err = await auth.updateNickname(nickInput.value)
  if (err) { nickError.value = err; return }
  snackMsg.value = "�ǳ��ѱ���"; snackbar.value = true
}
async function savePassword() {
  if (!oldPwd.value || !newPwd.value || newPwd.value.length < 4) return
  try {
    const res = await fetch(API + "/auth/password", {
      method: "PATCH", headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: auth.userName, oldPassword: oldPwd.value, newPassword: newPwd.value })
    })
    const result = await res.json()
    if (result.error) return
    oldPwd.value = ""; newPwd.value = ""
    snackMsg.value = "�������޸�"; snackbar.value = true
  } catch {}
}
function formatDate(ts: number) { return new Date(ts).toLocaleString("zh-CN") }
</script>

<template>
  <v-container fluid class="pa-6 admin-container" style="max-width:900px">
    <v-snackbar v-model="snackbar" :timeout="2000" location="top right" color="success" variant="tonal">
      {{ snackMsg }}
    </v-snackbar>
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" size="small" class="mr-2" @click="emit('back')" />
      <div>
        <h1 class="text-h4 font-weight-bold mb-1">��̨����</h1>
        <p class="text-body-2 text-medium-emphasis">�����û��뱸��¼</p>
      </div>
      <v-spacer />
      <v-btn prepend-icon="mdi-refresh" variant="text" size="small" :loading="loading" @click="loadData">ˢ��</v-btn>
    </div>

    <v-tabs v-model="tab" color="primary" class="mb-4">
      <v-tab value="overview" v-if="auth.isAdmin"><v-icon start size="small">mdi-view-dashboard</v-icon>����</v-tab>
      <v-tab value="system" v-if="auth.isAdmin"><v-icon start size="small">mdi-cog</v-icon>ϵͳ����</v-tab>
      <v-tab value="users" v-if="auth.isAdmin"><v-icon start size="small">mdi-account-group</v-icon>�û�����</v-tab>
      <v-tab value="profile"><v-icon start size="small">mdi-account</v-icon>��������</v-tab>
    </v-tabs>

    <template v-if='tab === "overview" && auth.isAdmin'>
      <v-card variant="outlined" class="rounded-xl pa-6 mb-4 stat-card">
        <h3 class="text-subtitle-1 font-weight-medium mb-4">��վ����</h3>
        <div class="d-flex align-center justify-space-between py-3">
          <div class="d-flex align-center ga-3">
            <v-icon color="primary">mdi-account</v-icon>
            <span class="text-body-2">�û�����</span>
          </div>
          <span class="text-h5 font-weight-bold">{{ stats?.totalUsers || 0 }}</span>
        </div>
        <v-divider />
        <div class="d-flex align-center justify-space-between py-3">
          <div class="d-flex align-center ga-3">
            <v-icon color="primary">mdi-pencil-box-multiple</v-icon>
            <span class="text-body-2">����¼����</span>
          </div>
          <span class="text-h5 font-weight-bold">{{ stats?.totalNotes || 0 }}</span>
        </div>
      </v-card>
    </template>

    <template v-if='tab === "system" && auth.isAdmin'>
      <v-card variant="outlined" class="rounded-xl pa-6 mb-4 stat-card">
        <h3 class="text-subtitle-1 font-weight-medium mb-4">ϵͳ����</h3>
        <div class="d-flex flex-column ga-4">
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center ga-3">
              <v-icon color="primary">mdi-web</v-icon>
              <span class="text-body-2">��վ����</span>
            </div>
            <div class="d-flex align-center ga-2" style="flex:1;max-width:400px">
              <v-text-field v-model="siteTitle" variant="outlined" hide-details density="compact" placeholder="��վ����" style="width:100%" @keyup.enter="saveSiteTitle" />
              <v-btn size="small" variant="tonal" color="primary" @click="saveSiteTitle">����</v-btn>
            </div>
          </div>
          <v-divider />
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center ga-3">
              <v-icon color="primary">mdi-account-plus</v-icon>
              <span class="text-body-2">�������û�ע��</span>
            </div>
            <v-switch v-model="allowRegister" hide-details density="compact" @update:model-value="toggleRegister" color="primary" />
          </div>
          <v-divider />
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center ga-3">
              <v-icon color="primary">mdi-certificate-outline</v-icon>
              <span class="text-body-2">������</span>
            </div>
            <div class="d-flex align-center ga-2" style="flex:1;max-width:400px">
              <v-text-field v-model="siteIcp" variant="outlined" hide-details density="compact" placeholder="��ICP��xxxxxxxx��" style="width:100%" @keyup.enter="saveSiteIcp" />
              <v-btn size="small" variant="tonal" color="primary" @click="saveSiteIcp">����</v-btn>
            </div>
          </div>
          <v-divider />
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center ga-3">
              <v-icon color="primary">mdi-apps</v-icon>
              <span class="text-body-2">������ͼ��</span>
            </div>
            <div class="d-flex align-center ga-2">
              <v-btn size="small" variant="tonal" color="primary" @click="showAppIconPicker = true">�޸�</v-btn>
            </div>
          </div>
          <v-divider />
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center ga-3">
              <v-icon color="primary">mdi-image-multiple</v-icon>
              <span class="text-body-2">��վͼ�� (Favicon)</span>
            </div>
            <div class="d-flex align-center ga-2">
              <v-btn size="small" variant="tonal" color="primary" @click="showFaviconPicker = true">����</v-btn>
            </div>
          </div>
        </div>
      </v-card>
      <AppIconPicker v-model="showAppIconPicker" />
      <FaviconPicker v-model="showFaviconPicker" />
    </template>

    <template v-if='tab === "users" && auth.isAdmin'>
      <v-card variant="outlined" class="rounded-xl stat-card">
        <v-list lines="two" bg-color="transparent">
          <v-list-item v-for="user in users" :key="user.id">
            <template #prepend>
              <v-avatar color="primary" variant="tonal">
                <v-img v-if="user.avatar && (user.avatar.startsWith('/uploads/') || user.avatar.startsWith('http'))" :src="user.avatar" alt="" cover />
                <span v-else class="font-weight-medium">{{ (user.nickname || user.username).charAt(0).toUpperCase() }}</span>
              </v-avatar>
            </template>
            <v-list-item-title>{{ user.nickname || user.username }}</v-list-item-title>
            <v-list-item-subtitle>@{{ user.username }} - {{ formatDate(user.createdAt) }} - {{ user.memoCount }}�� -
              <v-chip size="x-small" :color="user.role === 'admin' ? 'primary' : 'default'" variant="tonal">
                {{ user.role === 'admin' ? '����Ա' : '��ͨ�û�' }}
              </v-chip>
            </v-list-item-subtitle>
            <template #append>
              <v-btn v-if="user.username !== auth.userName" icon="mdi-delete" size="small" variant="text" color="error" :loading="deleting === user.id" @click="deleteUser(user.id)" />
            </template>
          </v-list-item>
        </v-list>
      </v-card>
    </template>

    <div class="text-center text-caption text-medium-emphasis pt-4">v1.0.0</div>
  </v-container>
</template>

<style scoped>
.stat-card { border: 1px solid rgba(var(--v-theme-on-surface), 0.12); }
@media (max-width: 768px) {
  .admin-container { padding: 12px !important; }
  .admin-container :deep(.v-tabs) { flex-wrap: nowrap; overflow-x: auto; }
  .admin-container :deep(.v-tab) { min-width: auto; padding: 0 12px; font-size: 0.8rem; }
}
</style>


