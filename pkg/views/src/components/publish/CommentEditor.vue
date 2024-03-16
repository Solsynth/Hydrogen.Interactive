<template>
  <v-card title="Leave your comment" :loading="loading">
    <v-form @submit.prevent="postComment">
      <v-card-text>
        <v-alert v-if="editor.related.edit_to" class="mb-5" type="info" variant="tonal">
          You are editing a comment with alias <b class="font-mono">{{ editor.related.edit_to?.alias }}</b>
        </v-alert>

        <v-textarea required hide-details variant="outlined" label="What do you want to say?" v-model="data.content" />

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
import { computed, ref, watch } from "vue"

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

const data = ref<any>({
  content: ""
})

async function postComment(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const payload = data.value

  if (!payload.content) return

  const url = editor.related.edit_to
    ? `/api/p/comments/${editor.related.edit_to?.id}`
    : `/api/p/${target.value?.model_type}/${target.value?.alias}/comments`
  const method = editor.related.edit_to ? "PUT" : "POST"

  loading.value = true
  const res = await request(url, {
    method: method,
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
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

watch(editor.related, (val) => {
  if (val.edit_to && val.edit_to.model_type === "comment") {
    data.value = val.edit_to
  }
})
</script>
