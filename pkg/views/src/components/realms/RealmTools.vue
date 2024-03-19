<template>
  <v-dialog v-model="realms.show.editor" class="max-w-[540px]">
    <realm-editor @relist="realms.list" />
  </v-dialog>
  <v-dialog v-model="realms.show.delete" class="max-w-[540px]">
    <realm-deletion @relist="realms.list" />
  </v-dialog>
</template>

<script setup lang="ts">
import { useRealms } from "@/stores/realms"
import RealmEditor from "@/components/realms/RealmEditor.vue"
import RealmDeletion from "./RealmDeletion.vue"

const realms = useRealms()
</script>
