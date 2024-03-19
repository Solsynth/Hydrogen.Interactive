<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="timeline flex-grow-1 mt-[-16px]">
      <post-list v-model:posts="posts" :loader="readMore" />
    </div>

    <div class="aside sticky top-0 w-full h-fit md:min-w-[280px] md:max-w-[320px] max-md:order-first">
      <v-card title="Realm Info" :loading="loading">
        <template #title>
          <div class="flex justify-between">
            <span>Realm Info</span>

            <realm-action :item="metadata" />
          </div>
        </template>
        <template #text>
          <div>
            <h2 class="font-medium">Name</h2>
            <p>{{ metadata?.name }}</p>

            <h2 class="font-medium mt-2">Description</h2>
            <div v-html="parseContent(metadata?.description ?? '')"></div>
          </div>
        </template>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import { request } from "@/scripts/request"
import { useRealms } from "@/stores/realms"
import { useRoute } from "vue-router"
import { parse } from "marked"
import dompurify from "dompurify"
import PostList from "@/components/posts/PostList.vue"
import RealmAction from "@/components/realms/RealmAction.vue"

const route = useRoute()
const realms = useRealms()

const loading = ref(false)
const error = ref<string | null>(null)
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const metadata = ref<any>(null)
const posts = ref<any[]>([])

async function readMetadata() {
  loading.value = true
  const res = await request(`/api/realms/${route.params.realmId}`)
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    error.value = null
    metadata.value = await res.json()
  }
  loading.value = false
}

async function readPosts() {
  const res = await request(
    `/api/feed?` +
      new URLSearchParams({
        take: pagination.pageSize.toString(),
        offset: ((pagination.page - 1) * pagination.pageSize).toString(),
        realmId: route.params.realmId as string
      })
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    error.value = null
    const data = await res.json()
    pagination.total = data["count"]
    posts.value.push(...(data["data"] ?? []))
  }
}

async function readMore({ done }: any) {
  // Reach the end of data
  if (pagination.total <= pagination.page * pagination.pageSize) {
    done("empty")
    return
  }

  pagination.page++
  await readPosts()

  if (error.value != null) done("error")
  else {
    if (pagination.total > 0) done("ok")
    else done("empty")
  }
}

watch(
  () => route.params.realmId,
  () => {
    posts.value = []
    readMetadata()
    readPosts()
  },
  { immediate: true }
)

watch(realms, (val) => {
  if (val.done) {
    readMetadata().then(() => (realms.done = false))
  }
})

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string)
}
</script>
