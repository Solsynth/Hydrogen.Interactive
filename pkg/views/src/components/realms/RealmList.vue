<template>
  <v-list density="comfortable">
    <v-list-subheader>
      Realms
      <v-badge color="warning" content="Alpha" inline />
    </v-list-subheader>

    <v-list-item
      v-for="item in realms"
      exact
      prepend-icon="mdi-account-multiple"
      :to="{ name: 'realms.details', params: { realmId: item.id } }"
      :title="item.name"
    />

    <v-divider v-if="realms.length > 0" class="border-opacity-75 my-2" />

    <v-list-item
      prepend-icon="mdi-plus"
      title="Create a realm"
      :disabled="!id.userinfo.isLoggedIn"
      @click="creating = true"
    />
  </v-list>

  <realm-editor v-model:show="creating" @relist="list" />

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { useUserinfo } from "@/stores/userinfo"
import { useEditor } from "@/stores/editor"
import RealmEditor from "@/components/realms/RealmEditor.vue"

const id = useUserinfo()
const editor = useEditor()

const realms = computed(() => editor.availableRealms)

const creating = ref(false)

const error = ref<string | null>(null)
const reverting = ref(false)

async function list() {
  reverting.value = true
  try {
    await editor.listRealms()
  } catch (err) {
    error.value = (err as Error).message
  }
  reverting.value = false
}
</script>
