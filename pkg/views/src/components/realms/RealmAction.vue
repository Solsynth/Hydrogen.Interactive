<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn v-bind="props" icon="mdi-dots-vertical" variant="text" size="x-small" />
    </template>

    <v-list density="compact" lines="one">
      <v-list-item disabled append-icon="mdi-flag" title="Report" />
      <v-list-item v-if="isOwned" append-icon="mdi-pencil" title="Edit" @click="editRealm" />
      <v-list-item v-if="isOwned" append-icon="mdi-delete" title="Delete" @click="deleteRealm" />
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { useRealms } from "@/stores/realms";
import { useUserinfo } from "@/stores/userinfo";
import { computed } from "vue"

const id = useUserinfo()
const realms = useRealms()

const props = defineProps<{ item: any }>()

const isOwned = computed(() => props.item?.account_id === id.userinfo.data.id)

function editRealm() {
  realms.related.edit_to = props.item
  realms.show.editor = true
}

function deleteRealm() {
  realms.related.delete_to = props.item
  realms.show.delete = true
}
</script>
