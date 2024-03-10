<template>
  <v-card title="Record a moment" :loading="loading">
    <v-form @submit.prevent="postMoment">
      <v-card-text>
        <v-textarea required hide-details name="content" variant="outlined" label="What's happened?!" />

        <div class="flex mt-1">
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
          <v-tooltip text="Categories" location="start">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                type="button"
                variant="text"
                icon="mdi-shape"
                size="small"
                @click="dialogs.categories = true"
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

  <planned-publish v-model:show="dialogs.plan" v-model:value="extras.publishedAt" />
  <media v-model:show="dialogs.media" v-model:uploading="uploading" v-model:value="extras.attachments" />

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
import { reactive, ref } from "vue"
import PlannedPublish from "@/components/publish/parts/PlannedPublish.vue"
import Media from "@/components/publish/parts/Media.vue"

const editor = useEditor()

const dialogs = reactive({
  plan: false,
  categories: false,
  media: false
})

const extras = reactive({
  publishedAt: null,
  attachments: []
})

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)
const uploading = ref(false)

async function postMoment(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const data: any = Object.fromEntries(new FormData(form))
  if (!data.hasOwnProperty("content")) return
  if (!extras.publishedAt) data["published_at"] = new Date().toISOString()
  else data["published_at"] = extras.publishedAt

  data["attachments"] = extras.attachments

  loading.value = true
  const res = await request("/api/p/moments", {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(data)
  })
  if (res.status === 200) {
    form.reset()
    success.value = true
    editor.show.moment = false
  } else {
    error.value = await res.text()
  }
  loading.value = false
}
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