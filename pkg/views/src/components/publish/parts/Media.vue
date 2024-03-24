<template>
  <v-dialog
    eager
    class="max-w-[540px]"
    :model-value="props.show"
    @update:model-value="(val) => emits('update:show', val)"
  >
    <v-card title="Media management">
      <template #text>
        <v-file-input
          prepend-icon=""
          append-icon="mdi-upload"
          variant="solo-filled"
          label="File Picker"
          v-model="picked"
          :accept="['image/*', 'video/*', 'audio/*']"
          :loading="props.uploading"
          @click:append="upload()"
        />

        <h2 class="px-2 mb-1">Media list</h2>
        <v-card variant="tonal">
          <v-list>
            <v-list-item v-for="(item, idx) in props.value" :title="getFileName(item)">
              <template #subtitle> {{ getFileType(item) }} · {{ formatBytes(item.filesize) }} </template>
              <template #append>
                <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="dispose(idx)" />
              </template>
            </v-list-item>
          </v-list>
        </v-card>
      </template>
      <template #actions>
        <v-btn class="ms-auto" text="Ok" @click="emits('update:show', false)"></v-btn>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import { ref } from "vue"

const props = defineProps<{ show: boolean; uploading: boolean; value: any[] }>()
const emits = defineEmits(["update:show", "update:uploading", "update:value"])

const picked = ref<any[]>([])

const error = ref<string | null>(null)

async function upload(file?: any) {
  if (props.uploading) return

  const data = new FormData()
  if (!file) {
    if (!picked.value) return
    data.set("attachment", picked.value[0])
  } else {
    data.set("attachment", file)
  }

  data.set("hashcode", await calculateHashCode(picked.value[0]))

  emits("update:uploading", true)
  const res = await request("/api/attachments", {
    method: "POST",
    headers: { Authorization: `Bearer ${getAtk()}` },
    body: data
  })
  let meta: any
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    meta = await res.json()
    emits("update:value", props.value.concat([meta.info]))
    picked.value = []
  }
  emits("update:uploading", false)
  return meta
}

async function dispose(idx: number) {
  const media = JSON.parse(JSON.stringify(props.value))
  const item = media.splice(idx)[0]
  emits("update:value", media)

  const res = await request(`/api/attachments/${item.id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${getAtk()}` }
  })
  if (res.status !== 200) {
    error.value = await res.text()
  }
}

defineExpose({ upload, dispose })

async function calculateHashCode(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = async () => {
      const buffer = reader.result as ArrayBuffer
      const hashBuffer = await crypto.subtle.digest("SHA-256", buffer)
      const hashArray = Array.from(new Uint8Array(hashBuffer))
      const hashHex = hashArray.map((byte) => byte.toString(16).padStart(2, "0")).join("")
      resolve(hashHex)
    }
    reader.onerror = () => {
      reject(reader.error)
    }
    reader.readAsArrayBuffer(file)
  })
}

function getFileName(item: any) {
  return item.filename.replace(/\.[^/.]+$/, "")
}

function getFileType(item: any) {
  switch (item.type) {
    case 1:
      return "Photo"
    case 2:
      return "Video"
    case 3:
      return "Audio"
    default:
      return "Others"
  }
}

function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return "0 Bytes"

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ["Bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"]

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}
</script>
