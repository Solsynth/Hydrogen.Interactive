<template>
  <v-card title="Delete a post" :loading="loading">
    <template #text>
      You are deleting a post with alias
      <b class="font-mono">{{ editor.related.delete_to?.alias }}</b>
      Are you confirm?
    </template>
    <template #actions>
      <div class="w-full flex justify-end">
        <v-btn color="grey-darken-3" @click="editor.show.delete = false">Not really</v-btn>
        <v-btn color="error" :disabled="loading" @click="deletePost">Yes</v-btn>
      </div>
    </template>
  </v-card>

  <v-snackbar v-model="success" :timeout="3000">The post has been deleted.</v-snackbar>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor"
import { getAtk } from "@/stores/userinfo"
import { ref } from "vue"

const editor = useEditor()

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)

async function deletePost() {
  const target = editor.related.delete_to
  const url = `/api/p/${target.model_type}/${target.id}`

  loading.value = true
  const res = await fetch(url, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${getAtk()}` }
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    success.value = true
    editor.show.delete = false
    editor.related.delete_to = null
  }
  loading.value = false
}
</script>
