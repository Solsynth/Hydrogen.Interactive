import { request } from "@/scripts/request"
import { defineStore } from "pinia"
import { ref } from "vue"

export const useWellKnown = defineStore("well-known", () => {
  const wellKnown = ref({})

  async function readWellKnown() {
    const res = await request("/.well-known")
    wellKnown.value = await res.json()
  }

  return { wellKnown, readWellKnown }
})
