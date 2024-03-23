<template>
  <div>
    <v-list density="comfortable" lines="one">
      <v-list-item v-for="item in members" :title="item.account.nick">
        <template #subtitle>@{{ item.account.name }}</template>
        <template #prepend>
          <v-avatar
            color="grey-lighten-2"
            icon="mdi-account-circle"
            class="rounded-card me-2"
            size="small"
            :image="item?.account.avatar"
          />
        </template>
        <template #append>
          <v-btn
            icon="mdi-account-remove"
            size="x-small"
            color="error"
            variant="text"
            :disabled="!checkKickable(item)"
            @click="kickMember(item)"
          />
        </template>
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
            <realm-invitation
              :item="props.item"
              @relist="listMembers"
              @error="(val) => (error = val)"
              @close="isActive.value = false"
            />
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
import RealmInvitation from "@/components/realms/RealmInvitation.vue"

const id = useUserinfo()

const props = defineProps<{ item: any }>()

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

async function kickMember(item: any) {
  loading.value = true
  const res = await request(`/api/realms/${props.item?.id}/kick`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify({
        account_name: item.account.name
    })
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await listMembers(props.item?.id)
  }
  loading.value = false
}

function checkKickable(item: any) {
  if (item.account?.id === id.userinfo.data?.id) return false
  if (item.account?.id === props.item?.account_id) return false
  return true
}
</script>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
