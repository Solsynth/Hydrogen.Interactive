import { reactive, ref } from "vue"
import { defineStore } from "pinia"
import { checkLoggedIn, getAtk } from "@/stores/userinfo"
import { request } from "@/scripts/request"

export const useRealms = defineStore("realms", () => {
  const done = ref(false)

  const show = reactive({
    editor: false,
    delete: false
  })

  const related_to = reactive<{ edit_to: any; delete_to: any }>({
    edit_to: null,
    delete_to: null
  })

  const available = ref<any[]>([])

  async function list() {
    if (!checkLoggedIn()) return

    const res = await request("/api/realms/me/available", {
      headers: { Authorization: `Bearer ${getAtk()}` }
    })
    if (res.status !== 200) {
      throw new Error(await res.text())
    } else {
      available.value = await res.json()
    }
  }

  list().then(() => console.log("[STARTUP HOOK] Fetch available realm successes."))

  return { done, show, related: related_to, available, list }
})
