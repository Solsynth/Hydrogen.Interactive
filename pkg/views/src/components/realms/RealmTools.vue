<template>
  <v-dialog v-model="realms.show.editor" class="max-w-[540px]">
    <realm-editor @relist="realms.list" />
  </v-dialog>
</template>

<script setup lang="ts">
import { useRealms } from "@/stores/realms"
import RealmEditor from "@/components/realms/RealmEditor.vue"

const realms = useRealms()
</script>
