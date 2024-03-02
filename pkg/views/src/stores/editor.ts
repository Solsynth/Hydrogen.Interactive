import { defineStore } from "pinia";
import { ref } from "vue";

export const useEditor = defineStore("editor", () => {
  const show = ref(false);

  return { show };
});