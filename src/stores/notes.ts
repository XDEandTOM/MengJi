import { defineStore } from "pinia"
import { ref } from "vue"
import { useAuthStore } from "@/stores/auth"

const API = "/api"

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
}

function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 8)
}

export const useNotesStore = defineStore("notes", () => {
  const notes = ref<Note[]>([])
  const loaded = ref(false)

  async function fetchNotes() {
    try {
      const res = await fetch(`${API}/notes`)
      if (res.ok) { notes.value = await res.json(); loaded.value = true }
    } catch { console.warn("Failed to fetch notes from server") }
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
      const res = await fetch(`${API}/notes`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify(note),
      })
      if (res.ok) notes.value.unshift(note)
    } catch { console.warn("Failed to create note") }
  }

  async function updateNote(id: string, content: string, tags?: string[]) {
    const note = notes.value.find(m => m.id === id)
    if (!note) return
    const updatedAt = Date.now()
    try {
      const res = await fetch(`${API}/notes/${id}`, {
        method: "PUT", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content, tags: tags ?? note.tags, updatedAt }),
      })
      if (res.ok) { note.content = content; note.updatedAt = updatedAt; if (tags) note.tags = tags }
    } catch { console.warn("Failed to update note") }
  }

  async function deleteNote(id: string) {
    try {
      const res = await fetch(`${API}/notes/${id}`, { method: "DELETE" })
      if (res.ok) notes.value = notes.value.filter(n => n.id !== id)
    } catch { console.warn("Failed to delete note") }
  }

  async function togglePin(id: string) {
    const note = notes.value.find(m => m.id === id)
    if (!note) return
    try {
      const res = await fetch(`${API}/notes/${id}/pin`, { method: "PATCH" })
      if (res.ok) note.pinned = !note.pinned
    } catch { console.warn("Failed to toggle pin") }
  }

return { notes, loaded, fetchNotes, addNote, updateNote, deleteNote, togglePin }
})


