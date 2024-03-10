<template>
  <v-card :rounded="false">
    <v-form @submit.prevent="postArticle">
      <v-toolbar>
        <div class="article-toolbar">
          <div class="flex">
            <v-btn type="button" icon="mdi-close" @click="editor.show.article = false"></v-btn>
          </div>

          <v-toolbar-title>Write an article</v-toolbar-title>

          <div class="flex justify-end items-center">
            <v-tooltip text="Publish">
              <template #activator="{ props: binding }">
                <v-btn type="submit" icon="mdi-publish" v-bind="binding" :loading="loading" />
              </template>
            </v-tooltip>
          </div>
        </div>
      </v-toolbar>

      <v-card-text>
        <v-container class="article-container">
          <v-textarea
            required
            class="mb-3"
            variant="outlined"
            label="Content"
            hint="The content supports markdown syntax"
            v-model="data.content"
            @paste="pasteMedia"
          />

          <v-expansion-panels>
            <v-expansion-panel title="Brief describe">
              <template #text>
                <div class="mt-1">
                  <v-text-field
                    required
                    variant="solo-filled"
                    density="comfortable"
                    label="Title"
                    v-model="data.title"
                  />

                  <v-textarea
                    required
                    auto-grow
                    variant="solo-filled"
                    density="comfortable"
                    label="Description"
                    v-model="data.description"
                  />
                </div>
              </template>
            </v-expansion-panel>

            <v-expansion-panel title="Planned publish">
              <template #text>
                <div class="flex justify-between items-center">
                  <div>
                    <p class="text-xs">Your content will visible for public at</p>
                    <p class="text-lg font-medium">
                      {{ data.publishedAt ? new Date(data.publishedAt).toLocaleString() : new Date().toLocaleString() }}
                    </p>
                  </div>
                  <v-btn size="small" icon="mdi-pencil" variant="text" @click="dialogs.plan = true" />
                </div>
              </template>
            </v-expansion-panel>

            <v-expansion-panel title="Media">
              <template #text>
                <div class="flex justify-between items-center">
                  <div>
                    <p class="text-xs">This article attached</p>
                    <p class="text-lg font-medium">{{ data.attachments.length }} attachment(s)</p>
                  </div>
                  <v-btn size="small" icon="mdi-camera-plus" variant="text" @click="dialogs.media = true" />
                </div>
              </template>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-container>
      </v-card-text>
    </v-form>
  </v-card>

  <planned-publish v-model:show="dialogs.plan" v-model:value="data.publishedAt" />
  <media ref="media" v-model:show="dialogs.media" v-model:uploading="uploading" v-model:value="data.attachments" />

  <v-snackbar v-model="success" :timeout="3000">Your article has been published.</v-snackbar>
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
import { useRouter } from "vue-router"
import PlannedPublish from "@/components/publish/parts/PlannedPublish.vue"
import Media from "@/components/publish/parts/Media.vue"

const editor = useEditor()

const dialogs = reactive({
  plan: false,
  categories: false,
  media: false
})

const data = reactive<any>({
  title: "",
  content: "",
  description: "",
  publishedAt: null,
  attachments: []
})

const router = useRouter()

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)
const uploading = ref(false)

async function postArticle(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement

  if (uploading.value) return

  if (!data.content) return
  if (!data.title || !data.description) return
  if (!data.publishedAt) data.publishedAt = new Date().toISOString()

  loading.value = true
  const res = await request("/api/p/articles", {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(data)
  })
  if (res.status === 200) {
    const data = await res.json()
    form.reset()
    router.push({ name: "posts.details.articles", params: { alias: data.alias } })
    success.value = true
    editor.show.article = false
  } else {
    error.value = await res.text()
  }
  loading.value = false
}

const media = ref<any>(null)

function pasteMedia(evt: ClipboardEvent) {
  const files = evt.clipboardData?.files
  if (files) {
    Array.from(files).forEach((item) => {
      media.value.upload(item).then((meta: any) => {
        if (meta) {
          data.content += `\n![${item.name}](/api/attachments/o/${meta.info.file_id})`
        }
      })
    })
    evt.preventDefault()
  }
}
</script>

<style scoped>
.article-toolbar {
  display: grid;
  flex-grow: 1;
  align-items: center;
  margin: 0 16px;

  grid-template-columns: 1fr auto 1fr;
}

.article-container {
  max-width: 720px;
}

.snackbar-progress {
  margin-left: -16px;
  margin-right: -16px;
  margin-bottom: -14px;
  margin-top: 12px;
  width: calc(100% + 64px);
}
</style>
