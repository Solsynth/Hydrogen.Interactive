<template>
  <v-card title="Record a moment" :loading="loading">
    <v-form @submit.prevent="postMoment">
      <v-card-text>
        <v-alert v-if="editor.related.edit_to" class="mb-5" type="info" variant="tonal">
          You are editing a post with alias <b class="font-mono">{{ editor.related.edit_to?.alias }}</b>
        </v-alert>

        <v-textarea
          required
          persistent-counter
          variant="outlined"
          label="What's happened?!"
          counter="1024"
          v-model="data.content"
        />

        <div class="flex mt-[-18px]">
          <v-tooltip text="Planned publish" location="start">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                type="button"
                variant="text"
                icon="mdi-calendar"
                size="small"
                @click="dialogs.plan = true"
              />
            </template>
          </v-tooltip>
          <v-tooltip text="Media" location="start">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                type="button"
                variant="text"
                icon="mdi-camera"
                size="small"
                @click="dialogs.media = true"
              />
            </template>
          </v-tooltip>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn type="reset" color="grey-darken-3" @click="editor.show.moment = false">Cancel</v-btn>
        <v-btn type="submit" :disabled="loading">Publish</v-btn>
      </v-card-actions>
    </v-form>
  </v-card>

  <planned-publish v-model:show="dialogs.plan" v-model:value="data.published_at" />
  <media v-model:show="dialogs.media" v-model:uploading="uploading" v-model:value="data.attachments" />

  <v-snackbar v-model="success" :timeout="3000">Your post has been published.</v-snackbar>
  <v-snackbar v-model="uploading" :timeout="-1">
    Uploading your media, please stand by...
    <v-progress-linear class="snackbar-progress" indeterminate />
  </v-snackbar>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useEditor } from "@/stores/editor"
import { getAtk } from "@/stores/userinfo"
import { reactive, ref, watch } from "vue"
import { useRouter } from "vue-router"
import PlannedPublish from "@/components/publish/parts/PlannedPublish.vue"
import Media from "@/components/publish/parts/Media.vue"

const editor = useEditor()

const dialogs = reactive({
  plan: false,
  media: false
})

const data = ref<any>({
  content: "",
  published_at: null,
  attachments: []
})

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)
const uploading = ref(false)

const router = useRouter()

async function postMoment(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const payload = data.value
  if (!payload.content) return
  if (!payload.published_at) payload.published_at = new Date().toISOString()

  const url = editor.related.edit_to ? `/api/p/moments/${editor.related.edit_to?.id}` : "/api/p/moments"
  const method = editor.related.edit_to ? "PUT" : "POST"

  loading.value = true
  const res = await request(url, {
    method: method,
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  })
  if (res.status === 200) {
    form.reset()
    const data = await res.json()
    success.value = true
    editor.show.moment = false
    router.push({ name: "posts.details.moments", params: { alias: data.alias } })
  } else {
    error.value = await res.text()
  }
  loading.value = false
}

watch(editor.related, (val) => {
  if (val.edit_to && val.edit_to.model_type === "moment") {
    data.value = val.edit_to
  }
})
</script>

<style>
.snackbar-progress {
  margin-left: -16px;
  margin-right: -16px;
  margin-bottom: -14px;
  margin-top: 12px;
  width: calc(100% + 64px);
}
</style>
