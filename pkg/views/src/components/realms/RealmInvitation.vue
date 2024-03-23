<template>
  <v-card prepend-icon="mdi-account-plus" title="Invite someone">
    <v-form @submit.prevent="inviteMember">
      <v-card-text>
        <v-text-field
          label="Username"
          variant="outlined"
          density="comfortable"
          hint="Require username not the nickname"
          v-model="targetName"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn type="reset" color="grey-darken-3" @click="emits('close')">Cancel</v-btn>
        <v-btn type="submit" :disabled="loading">Invite</v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"

const props = defineProps<{item: any}>()
const emits = defineEmits(["close", "error", "relist"])

const loading = ref(false)

const targetName = ref("")

async function inviteMember(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement

  loading.value = true
  const res = await request(`/api/realms/${props.item?.id}/invite`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify({
        account_name: targetName.value
    })
  })
  if (res.status !== 200) {
    emits("error", await res.text())
  } else {
    form.reset()
    emits("relist")
  }
  loading.value = false
}
</script>
