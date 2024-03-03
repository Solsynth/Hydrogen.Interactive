import { defineStore } from "pinia";
import { reactive, ref } from "vue";

export const useEditor = defineStore("editor", () => {
  const show = reactive({
    moment: false,
    article: false,
  });

  return { show };
});