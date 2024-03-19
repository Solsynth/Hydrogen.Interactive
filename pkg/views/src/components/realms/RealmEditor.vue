<template>
  <v-card title="Organize a realm" prepend-icon="mdi-account-multiple" :loading="loading">
    <v-form @submit.prevent="submit">
      <v-card-text>
        <v-text-field label="Name" variant="outlined" density="comfortable" v-model="data.name" />
        <v-textarea label="Description" variant="outlined" density="comfortable" v-model="data.description" />
        <v-select
          label="Realm type"
          item-title="label"
          item-value="value"
          variant="outlined"
          density="comfortable"
          :items="realmTypeOptions"
          v-model="data.realm_type"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn type="reset" color="grey-darken-3" @click="realms.show.editor = false">Cancel</v-btn>
        <v-btn type="submit" :disabled="loading">Save</v-btn>
      </v-card-actions>
    </v-form>
  </v-card>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { getAtk } from "@/stores/userinfo"
import { useRealms } from "@/stores/realms"

const emits = defineEmits(["relist"])

const realms = useRealms()

const realmTypeOptions = [
  { label: "Public Realm", value: 0 },
  { label: "Restricted Realm", value: 1 },
  { label: "Private Realm", value: 2 }
]

const error = ref<null | string>(null)
const loading = ref(false)

const data = ref({
  name: "",
  description: "",
  realm_type: 0
})

async function submit(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const payload = data.value
  if (!payload.name) return

  const url = realms.related.edit_to ? `/api/realms/${realms.related.edit_to?.id}` : "/api/realms"
  const method = realms.related.edit_to ? "PUT" : "POST"

  loading.value = true
  const res = await fetch(url, {
    method: method,
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    emits("relist")
    form.reset()
    realms.done = true
    realms.show.editor = false
    realms.related.edit_to = null
  }
  loading.value = false
}

watch(
  realms.related,
  (val) => {
    if (val.edit_to) {
      data.value = JSON.parse(JSON.stringify(val.edit_to))
    }
  },
  { immediate: true }
)
</script>
