<template>
  <v-card title="Delete a realm" :loading="loading">
    <template #text>
      You are deleting a realm
      <b>{{ realms.related.delete_to?.name }}</b> <br />
      All posts belonging to this domain will be deleted and never appear again. Are you confirm?
    </template>
    <template #actions>
      <div class="w-full flex justify-end">
        <v-btn color="grey-darken-3" @click="realms.show.delete = false">Not really</v-btn>
        <v-btn color="error" :disabled="loading" @click="deletePost">Yes</v-btn>
      </div>
    </template>
  </v-card>

  <v-snackbar v-model="success" :timeout="3000">The realm has been deleted.</v-snackbar>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useRealms } from "@/stores/realms"
import { getAtk } from "@/stores/userinfo"
import { useRoute, useRouter } from "vue-router"
import { ref } from "vue"

const route = useRoute()
const router = useRouter()
const realms = useRealms()

const emits = defineEmits(["relist"])

const error = ref<string | null>(null)
const success = ref(false)
const loading = ref(false)

async function deletePost() {
  const target = realms.related.delete_to
  const url = `/api/realms/${target.id}`

  loading.value = true
  const res = await request(url, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${getAtk()}` }
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    success.value = true
    realms.show.delete = false
    realms.related.delete_to = null
    emits("relist")
    if (route.name?.toString()?.startsWith("realm")) {
      router.push({ name: "explore" })
    }
  }
  loading.value = false
}
</script>
