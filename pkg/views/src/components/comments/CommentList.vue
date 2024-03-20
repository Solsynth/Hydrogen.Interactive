<template>
  <div v-if="loading" class="text-center flex items-center justify-center">
    <v-progress-circular indeterminate />
  </div>

  <div v-else class="flex flex-col gap-5 mt-3">
    <div v-for="(item, idx) in props.comments" class="text-sm">
      <post-item :item="item" @update:item="(val) => updateItem(idx, val)" />
    </div>
  </div>

  <v-divider class="mt-2 mb-3 border-opacity-50 mx-[-1rem]" />

  <v-btn block prepend-icon="mdi-pencil" variant="plain" :disabled="!id.userinfo.isLoggedIn" @click="leaveComment">
    Leave your comment
  </v-btn>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { reactive, ref, watch } from "vue"
import { useEditor } from "@/stores/editor"
import { useUserinfo } from "@/stores/userinfo"
import PostItem from "@/components/posts/PostItem.vue"

const id = useUserinfo()
const editor = useEditor()

const props = defineProps<{
  comments: any[]
  model: string
  dataset: string
  alias: any
  item: any
}>()
const emits = defineEmits(["update:comments"])

const loading = ref(false)
const error = ref<string | null>(null)

const pagination = reactive({ page: 0, pageSize: 10, total: 0 })

async function readComments() {
  loading.value = true
  const res = await request(
    `/api/p/${props.dataset}/${props.alias}/comments?` +
      new URLSearchParams({
        take: pagination.pageSize.toString(),
        offset: (pagination.page * pagination.pageSize).toString()
      })
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    error.value = null
    const data = await res.json()
    pagination.total = data["total"]
    emits("update:comments", data["data"])
  }
  loading.value = false
}

readComments()

function updateItem(idx: number, data: any) {
  const comments = JSON.parse(JSON.stringify(props.comments))
  comments[idx] = data
  emits("update:comments", comments)
}

watch(editor, (val) => {
  if (val.done) {
    readComments().then(() => (val.done = false))
  }
})

function leaveComment() {
  editor.related.comment_to = JSON.parse(JSON.stringify(props.item))
  editor.related.comment_to.model_type = props.dataset
  editor.show.comment = true
}
</script>
