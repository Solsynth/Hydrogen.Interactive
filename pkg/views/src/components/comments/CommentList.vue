<template>
  <div v-if="loading" class="text-center flex items-center justify-center">
    <v-progress-circular indeterminate />
  </div>

  <div v-else class="flex flex-col gap-2 mt-3">
    <div v-for="(item, idx) in props.comments" class="text-sm">
      <post-item :item="item" @update:item="val => updateItem(idx, val)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { reactive, ref } from "vue"
import PostItem from "@/components/posts/PostItem.vue"

const props = defineProps<{
  comments: any[]
  model: any
  alias: any
}>()
const emits = defineEmits(["update:comments"])

const loading = ref(false)
const error = ref<string | null>(null)

const pagination = reactive({ page: 0, pageSize: 10, total: 0 })

async function readComments() {
  loading.value = true
  const res = await request(
    `/api/p/${props.model}/${props.alias}/comments?` +
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
  const comments = JSON.parse(JSON.stringify(props.comments));
  comments[idx] = data;
  emits("update:comments", comments);
}
</script>
