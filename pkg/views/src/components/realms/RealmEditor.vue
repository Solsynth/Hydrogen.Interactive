<template>
  <v-dialog :model-value="props.show" @update:model-value="(val) => emits('update:show', val)" class="max-w-[540px]">
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

          <v-btn type="reset" color="grey-darken-3" @click="emits('update:show', false)">Cancel</v-btn>
          <v-btn type="submit" :disabled="loading">Save</v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { getAtk } from "@/stores/userinfo"

const props = defineProps<{ show: boolean }>()
const emits = defineEmits(["update:show", "relist"])

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

  loading.value = true
  const res = await fetch("/api/realms", {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    emits("relist")
    form.reset()
    emits("update:show", false)
  }
  loading.value = false
}
</script>
