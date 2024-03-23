<template>
  <div>
    <v-list density="comfortable" lines="one">
      <v-list-item v-for="item in members" :title="item.account.nick">
        <template #prepend>
          <v-avatar
            color="grey-lighten-2"
            icon="mdi-account-circle"
            class="rounded-card me-2"
            size="small"
            :image="item?.account.avatar"
          />
        </template>
        <template #subtitle>@{{ item.account.name }}</template>
      </v-list-item>
    </v-list>

    <div v-if="isOwned">
      <v-divider class="mt-2 mb-3 border-opacity-50 mx-[-1rem]" />

      <div class="px-3">
        <v-dialog class="max-w-[540px]">
          <template #activator="{ props }">
            <v-btn v-bind="props" block prepend-icon="mdi-account-plus" variant="plain"> Invite someone </v-btn>
          </template>

          <template #default="{ isActive }">
            <v-card prepend-icon="mdi-account-plus" title="Invite someone">
              <v-form @submit.prevent="inviteMember">
                <v-card-text>
                  <v-text-field
                    label="Username"
                    variant="outlined"
                    density="comfortable"
                    hint="Require username not the nickname"
                    v-model="data.account_name"
                  />
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>

                  <v-btn type="reset" color="grey-darken-3" @click="isActive.value = false">Cancel</v-btn>
                  <v-btn type="submit" :disabled="loading">Invite</v-btn>
                </v-card-actions>
              </v-form>
            </v-card>
          </template>
        </v-dialog>
      </div>
    </div>

    <!-- @vue-ignore -->
    <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { request } from "@/scripts/request"
import { getAtk, useUserinfo } from "@/stores/userinfo"
import { computed } from "vue"

const id = useUserinfo()

const props = defineProps<{ item: any }>()

const data = ref<any>({
  account_name: ""
})

const members = ref<any[]>([])

const isOwned = computed(() => {
  return id.userinfo.data?.id === props.item?.account_id
})

const loading = ref(false)
const error = ref<string | null>(null)

watch(
  () => props.item,
  (val) => {
    if (val?.id) {
      listMembers(val.id)
    }
  },
  { deep: true, immediate: true }
)

async function listMembers(id: number) {
  loading.value = true
  const res = await request(`/api/realms/${id}/members`)
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    error.value = null
    members.value = await res.json()
  }
  loading.value = false
}

async function inviteMember(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement
  const payload = data.value

  loading.value = true
  const res = await request(`/api/realms/${props.item?.id}/invite`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    form.reset()
    await listMembers(props.item?.id)
  }
  loading.value = false
}
</script>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
