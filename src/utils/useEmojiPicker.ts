import { ref, watch } from "vue"

export function useEmojiPicker() {
  const showEmojiPicker = ref(false)
  const emojiCategories = ref<{ id: number; icon: string; list: string[] }[]>([])
  const activeEmojiCat = ref(0)

  const groupLabels: Record<number, string> = {
    0: "😊", 1: "🤝", 3: "🐻", 4: "🍔", 5: "🏠", 6: "⚽", 7: "💡", 8: "❤️", 9: "🚩",
  }

  async function loadEmojiData() {
    if (emojiCategories.value.length) return
    const raw = (await import("emojibase-data/zh/compact.json")).default
    const cats = [0, 1, 3, 4, 5, 6, 7, 8, 9].map((g) => ({
      id: g,
      icon: groupLabels[g] || "?",
      list: [] as string[],
    }))
    for (const e of raw) {
      if (e.group === undefined || e.group === 2) continue
      const cat = cats.find((c) => c.id === e.group)
      if (cat && e.unicode) cat.list.push(e.unicode)
    }
    emojiCategories.value = cats
  }

  watch(showEmojiPicker, (v) => {
    if (v) loadEmojiData()
  })

  return { showEmojiPicker, emojiCategories, activeEmojiCat }
}
