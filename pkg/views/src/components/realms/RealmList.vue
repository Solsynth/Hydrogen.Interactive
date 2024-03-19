<template>
  <v-list density="comfortable">
    <v-list-subheader>
      Realms
      <v-badge color="warning" content="Alpha" inline />
    </v-list-subheader>

    <v-list-item
      v-for="item in realms.available"
      exact
      prepend-icon="mdi-account-multiple"
      :to="{ name: 'realms.page', params: { realmId: item.id } }"
      :title="item.name"
    />

    <v-divider v-if="realms.available.length > 0" class="border-opacity-75 my-2" />

    <v-list-item
      prepend-icon="mdi-plus"
      title="Create a realm"
      :disabled="!id.userinfo.isLoggedIn"
      @click="createRealm"
    />
  </v-list>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"
import { useRealms } from "@/stores/realms"

const id = useUserinfo()
const realms = useRealms()

function createRealm() {
  realms.related.edit_to = null
  realms.related.delete_to = null
  realms.show.editor = true
}
</script>
