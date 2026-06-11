import { defineStore } from "pinia"
import { ref, computed } from "vue"
import { useAuthStore } from "@/stores/auth"
import { authFetch } from "@/utils/api"

const API = "/api"
const PAGE_SIZE = 20

export interface NoteReaction { [emoji: string]: string[] }

export interface Note {
  id: string
  content: string
  createdAt: number
  updatedAt: number
  pinned: boolean
  tags: string[]
  username: string
  avatar?: string
  nickname?: string
  reactions?: NoteReaction
}

function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 8)
}

export const useNotesStore = defineStore("notes", () => {
  const notes = ref<Note[]>([])
  const loaded = ref(false)
  const loadingMore = ref(false)
  const hasMore = ref(true)
  const total = ref(0)
  const page = ref(0)

  // Current filter state (synced with the page)
  const searchQuery = ref("")
  const selectedTag = ref("")
  const selectedDay = ref("")

  const allTags = computed(() => {
    const tagCount = new Map<string, number>()
    for (const n of notes.value) {
      if (n.tags && Array.isArray(n.tags)) for (const t of n.tags) tagCount.set(t, (tagCount.get(t) || 0) + 1)
    }
    return [...tagCount.entries()].sort((a, b) => b[1] - a[1])
  })

  async function fetchNotes(reset = true) {
    try {
      if (reset) {
        page.value = 0
        hasMore.value = true
        loadingMore.value = false
      }

      const params = new URLSearchParams()
      params.set("limit", String(PAGE_SIZE))
      params.set("offset", String(page.value * PAGE_SIZE))

      // Pass filter params to backend
      if (searchQuery.value.trim()) params.set("search", searchQuery.value.trim())
      if (selectedTag.value) params.set("tag", selectedTag.value)
      if (selectedDay.value) params.set("date", selectedDay.value)

      const auth = useAuthStore()
      if (auth.isLoggedIn && auth.userName) params.set("username", auth.userName)

      const url = `${API}/notes?${params.toString()}`
      loadingMore.value = !reset

      const res = await authFetch(url)
      if (res.ok) {
        const data = await res.json()
        if (data && typeof data === "object" && "notes" in data) {
          // Paginated response
          if (reset) {
            notes.value = data.notes
          } else {
            // Append, avoiding duplicates
            const existingIds = new Set(notes.value.map(n => n.id))
            for (const n of data.notes) {
              if (!existingIds.has(n.id)) {
                notes.value.push(n)
                existingIds.add(n.id)
              }
            }
          }
          total.value = data.total || 0
          hasMore.value = notes.value.length < total.value
          page.value++
        } else if (Array.isArray(data)) {
          // Legacy fallback (unpaginated)
          notes.value = data
          total.value = data.length
          hasMore.value = false
        }
        loaded.value = true
      }
    } catch {
      console.warn("Failed to fetch notes from server")
      // Fallback for offline: show cached empty state
      loaded.value = true
    } finally {
      loadingMore.value = false
    }
  }

  async function addNote(content: string, tags: string[] = [], username: string = "") {
    const auth = useAuthStore()
    const note: Note = {
      id: generateId(), content, createdAt: Date.now(), updatedAt: Date.now(),
      pinned: false, tags, username,
      avatar: auth.userAvatar || undefined,
      nickname: auth.userNickname || undefined,
    }
    try {
      const res = await authFetch(`${API}/notes`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify(note),
      })
      if (res.ok) {
        // Prepend and re-fetch to keep pagination consistent
        await fetchNotes(true)
      }
    } catch { console.warn("Failed to create note") }
  }

  async function updateNote(id: string, content: string, tags?: string[], username?: string) {
    const note = notes.value.find(m => m.id === id)
    if (!note) return
    const updatedAt = Date.now()
    try {
      const res = await authFetch(`${API}/notes/${id}`, {
        method: "PUT", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content, tags: tags ?? note.tags, updatedAt, username: username || useAuthStore().userName || "" }),
      })
      if (res.ok) {
        const auth = useAuthStore()
        note.content = content
        note.updatedAt = updatedAt
        if (tags !== undefined) note.tags = tags
        note.avatar = auth.userAvatar || undefined
        note.nickname = auth.userNickname || undefined
      }
    } catch { console.warn("Failed to update note") }
  }

  async function deleteNote(id: string, username?: string) {
    try {
      const auth = useAuthStore()
      const res = await authFetch(`${API}/notes/${id}?username=${encodeURIComponent(username || auth.userName || "")}`, { method: "DELETE" })
      if (res.ok) {
        notes.value = notes.value.filter(n => n.id !== id)
        total.value = Math.max(0, total.value - 1)
        hasMore.value = notes.value.length < total.value
      }
    } catch { console.warn("Failed to delete note") }
  }

  async function togglePin(id: string) {
    try {
      const res = await authFetch(`${API}/notes/${id}/pin`, { method: "PATCH" })
      if (res.ok) { await fetchNotes(true) }
    } catch { console.warn("Failed to toggle pin") }
  }

  async function reactToNote(id: string, emoji: string, uid?: string) {
    if (!uid) uid = useAuthStore().userName || ""
    try {
      const res = await authFetch(`${API}/notes/${id}/react`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ emoji, username: uid }),
      })
      if (res.ok) {
        const note = notes.value.find(n => n.id === id)
        if (note) {
          if (!note.reactions) note.reactions = {}
          if (!note.reactions[emoji]) note.reactions[emoji] = []
          if (!note.reactions[emoji].includes(uid)) note.reactions[emoji].push(uid)
        }
      }
    } catch { console.warn("store action failed") }
  }

  async function removeReaction(id: string, emoji: string, uid?: string) {
    if (!uid) uid = useAuthStore().userName || ""
    try {
      const res = await authFetch(`${API}/notes/${id}/react`, {
        method: "DELETE", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ emoji, username: uid }),
      })
      if (res.ok) {
        const note = notes.value.find(n => n.id === id)
        if (note && note.reactions && note.reactions[emoji]) {
          note.reactions[emoji] = note.reactions[emoji].filter(u => u !== uid)
          if (note.reactions[emoji].length === 0) delete note.reactions[emoji]
        }
      }
    } catch { console.warn("store action failed") }
  }

  return {
    notes, loaded, loadingMore, hasMore, total,
    searchQuery, selectedTag, selectedDay, allTags,
    fetchNotes, addNote, updateNote, deleteNote,
    togglePin, reactToNote, removeReaction,
  }
})
