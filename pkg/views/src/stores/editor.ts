import { defineStore } from "pinia"
import { reactive, ref } from "vue"
import { checkLoggedIn, getAtk } from "@/stores/userinfo"

export const useEditor = defineStore("editor", () => {
  const done = ref(false)

  const show = reactive({
    moment: false,
    article: false,
    comment: false,
    delete: false
  })

  const related = reactive<{
    edit_to: any
    comment_to: any
    reply_to: any
    repost_to: any
    delete_to: any
  }>({
    edit_to: null,
    comment_to: null,
    reply_to: null,
    repost_to: null,
    delete_to: null
  })

  const availableRealms = ref<any[]>([])

  async function listRealms() {
    if (!checkLoggedIn()) return

    const res = await fetch("/api/realms/me/available", {
      headers: { Authorization: `Bearer ${getAtk()}` }
    })
    if (res.status !== 200) {
      throw new Error(await res.text())
    } else {
      availableRealms.value = await res.json()
    }
  }

  listRealms().then(() => console.log("[STARTUP HOOK] Fetch available realm successes."))

  return { show, related, availableRealms, listRealms, done }
})
