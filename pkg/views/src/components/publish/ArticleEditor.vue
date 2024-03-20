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
          <v-alert v-if="editor.related.edit_to" class="mb-5" type="info" variant="tonal">
            You are editing a post with alias <b class="font-mono">{{ editor.related.edit_to?.alias }}</b>
          </v-alert>

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
                    :loading="reverting"
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
                      {{ data.published_at ? new Date(data.published_at).toLocaleString() : new Date().toLocaleString() }}
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

            <v-expansion-panel title="Publish area">
              <template #text>
                <div class="flex justify-between items-center">
                  <div>
                    <p class="text-xs">This article will publish in</p>
                    <p class="text-lg font-medium">{{ currentRealm?.name ?? "No realm" }}</p>
                  </div>
                  <v-btn size="small" icon="mdi-account-group" variant="text" @click="dialogs.area = true" />
                </div>
              </template>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-container>
      </v-card-text>
    </v-form>
  </v-card>

  <planned-publish v-model:show="dialogs.plan" v-model:value="data.published_at" />
  <media ref="media" v-model:show="dialogs.media" v-model:uploading="uploading" v-model:value="data.attachments" />
  <publish-area v-model:show="dialogs.area" v-model:value="data.realm_id" />

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
import { useRealms } from "@/stores/realms";
import { computed, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router"
import PlannedPublish from "@/components/publish/parts/PlannedPublish.vue"
import Media from "@/components/publish/parts/Media.vue"
import PublishArea from "@/components/publish/parts/PublishArea.vue";

const route = useRoute()
const realms = useRealms()
const editor = useEditor()

const dialogs = reactive({
  plan: false,
  categories: false,
  media: false,
  area: false,
})

const data = ref<any>({
  title: "",
  content: "",
  description: "",
  realm_id: null,
  published_at: null,
  attachments: []
})

const currentRealm = computed(() => {
  if(data.value.realm_id) {
    return realms.available.find((e: any) => e.id === data.value.realm_id)
  } else {
    return null
  }
})

const router = useRouter()

const error = ref<string | null>(null)
const success = ref(false)
const reverting = ref(false)
const loading = ref(false)
const uploading = ref(false)

async function postArticle(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement

  if (uploading.value) return

  const payload = data.value
  console.log(payload)
  if (!payload.content) return
  if (!payload.title || !payload.description) return
  if (!payload.published_at) payload.published_at = new Date().toISOString()
  if (!payload.realm_id) payload.realm_id = undefined

  const url = editor.related.edit_to ? `/api/p/articles/${editor.related.edit_to?.id}` : "/api/p/articles"
  const method = editor.related.edit_to ? "PUT" : "POST"

  loading.value = true
  const res = await request(url, {
    method: method,
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
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
          data.value.content += `\n![${item.name}](/api/attachments/o/${meta.info.file_id})`
        }
      })
    })
    evt.preventDefault()
  }
}

watch(editor.related, (val) => {
  if (val.edit_to && val.edit_to.model_type === "article") {
    request(`/api/p/articles/${val.edit_to.alias}`).then(async (res) => {
      data.value = await res.json()
      data.value.attachments = data.value.attachments ?? []
    })
  }
})

watch(
  () => route.params.realmId,
  (val) => {
    if (val) {
      data.value.realm_id = parseInt(val as string)
    }
  },
  { deep: true, immediate: true }
)
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
