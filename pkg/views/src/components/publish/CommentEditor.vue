<template>
  <v-card title="Leave your comment" :loading="loading">
    <v-form @submit.prevent="postComment">
      <v-card-text>
        <v-textarea required hide-details name="content" variant="outlined" label="What do you want to say?" />

        <p class="px-2 mt-1 text-body-2 opacity-80">Your comment will leave below {{ postIdentifier }}</p>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn type="reset" color="grey-darken-3" @click="editor.show.comment = false">Cancel</v-btn>
        <v-btn type="submit" :disabled="loading">Publish</v-btn>
      </v-card-actions>
    </v-form>
  </v-card>

  <v-snackbar v-model="success" :timeout="3000">Your comment has been published.</v-snackbar>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useEditor } from "@/stores/editor"
import { getAtk } from "@/stores/userinfo"
import { computed, ref } from "vue"

const editor = useEditor()

const target = computed<any>(() => editor.related.comment_to)
const postIdentifier = computed(() => {
  if (editor.related.comment_to?.title) {
    return `${editor.related.comment_to.title}`
  } else {
    return `#${editor.related.comment_to?.alias}`
  }
})

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)

async function postComment(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const data = new FormData(form)
  if (!data.has("content")) return

  loading.value = true
  const res = await request(`/api/p/${target.value?.model_type}/${target.value?.alias}/comments`, {
    method: "POST",
    headers: { Authorization: `Bearer ${getAtk()}` },
    body: data
  })
  if (res.status === 200) {
    form.reset()
    success.value = true
    editor.show.comment = false
  } else {
    error.value = await res.text()
  }
  loading.value = false
  editor.done = true
}
</script>
