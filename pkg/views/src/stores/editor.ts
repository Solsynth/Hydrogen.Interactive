import { defineStore } from "pinia"
import { reactive, ref } from "vue"

export const useEditor = defineStore("editor", () => {
  const done = ref(false)

  const show = reactive({
    moment: false,
    article: false,
    comment: false
  })

  const related = reactive<{ comment_to: any; reply_to: any; repost_to: any }>({
    comment_to: null,
    reply_to: null,
    repost_to: null
  })

  return { show, related, done }
})
