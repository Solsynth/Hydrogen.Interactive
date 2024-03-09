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

  <v-snackbar v-model="success" :timeout="3000">Your post has been published.</v-snackbar>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useEditor } from "@/stores/editor"
import { getAtk } from "@/stores/userinfo"
import { reactive, ref } from "vue"
import PlannedPublish from "@/components/publish/parts/PlannedPublish.vue"

const editor = useEditor()

const dialogs = reactive({
  plan: false,
  categories: false,
  media: false
})

const extras = reactive({
  publishedAt: null
})

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)

async function postMoment(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const data = new FormData(form)
  if (!data.has("content")) return
  if (!extras.publishedAt) data.set("published_at", new Date().toISOString())
  else data.set("published_at", extras.publishedAt)

  loading.value = true
  const res = await request("/api/p/moments", {
    method: "POST",
    headers: { Authorization: `Bearer ${getAtk()}` },
    body: data
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
