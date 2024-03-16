<template>
  <div class="flex gap-[8px] my-[8px]">
    <v-chip
      v-for="[k, v] in Object.entries(props.reactions)"
      :color="pickColor()"
      :size="props.size"
      @click="reactPost(k, emojis[k].attitude)"
    >
      <div class="ms-2">{{ v }}</div>
      <template #prepend>{{ emojis[k].icon }}</template>
    </v-chip>

    <v-menu v-if="!props.readonly" location="bottom center">
      <template v-slot:activator="{ props: binding }">
        <v-chip v-if="id.userinfo.isLoggedIn" v-bind="binding" :size="props.size" prepend-icon="mdi-emoticon-plus">
          React
        </v-chip>
      </template>

      <v-list density="compact" lines="one">
        <v-list-item v-for="[k, v] in Object.entries(emojis)" @click="reactPost(k, v.attitude)">
          <v-list-item-title class="font-mono">:{{ k }}:</v-list-item-title>
          <template #prepend>
            <div class="me-3">{{ v.icon }}</div>
          </template>
        </v-list-item>
      </v-list>
    </v-menu>

    <v-snackbar v-model="status.added" :timeout="3000">Your react has been added into post.</v-snackbar>
    <v-snackbar v-model="status.removed" :timeout="3000">Your react has been removed from post.</v-snackbar>

    <!-- @vue-ignore -->
    <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { getAtk, useUserinfo } from "@/stores/userinfo"
import { reactive, ref } from "vue"

const id = useUserinfo()

const emits = defineEmits(["update"])
const props = defineProps<{
  size?: string
  readonly?: boolean
  model: any
  item: any
  reactions: { [id: string]: number }
}>()

const emojis: { [id: string]: { icon: string; attitude: number } } = {
  thumb_up: { icon: "👍", attitude: 1 },
  clap: { icon: "👏", attitude: 1 }
}

function pickColor(): string {
  const colors = ["blue", "green", "purple"]
  const randomIndex = Math.floor(Math.random() * colors.length)
  return colors[randomIndex]
}

const status = reactive({ added: false, removed: false })
const error = ref<string | null>(null)

async function reactPost(symbol: string, attitude: number) {
  const res = await request(`/api/p/${props.model}/${props.item?.id}/react`, {
    method: "POST",
    headers: { Authorization: `Bearer ${getAtk()}`, "Content-Type": "application/json" },
    body: JSON.stringify({ symbol, attitude })
  })
  if (res.status === 201) {
    status.added = true
    emits("update", symbol, 1)
  } else if (res.status === 204) {
    status.removed = true
    emits("update", symbol, -1)
  } else {
    error.value = await res.text()
  }
}
</script>
