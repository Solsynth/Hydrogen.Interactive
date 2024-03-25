<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn v-bind="props" icon="mdi-dots-vertical" variant="text" size="x-small" />
    </template>

    <v-list density="compact" lines="one">
      <v-list-item disabled append-icon="mdi-flag" title="Report" />
      <v-list-item v-if="isOwned" append-icon="mdi-pencil" title="Edit" @click="editPost" />
      <v-list-item v-if="isOwned" append-icon="mdi-delete" title="Delete" @click="deletePost" />
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor"
import { useUserinfo } from "@/stores/userinfo"
import { computed } from "vue"

const id = useUserinfo()
const editor = useEditor()

const props = defineProps<{ item: any }>()

const isOwned = computed(() => props.item?.author_id === id.userinfo.data.id)

function editPost() {
  editor.related.edit_to = JSON.parse(JSON.stringify(props.item))
  // eslint-disable-next-line
  if (editor.show.hasOwnProperty(props.item.model_type)) {
    // @ts-ignore
    editor.show[props.item.model_type] = true
  }
  if (props.item.model_type === "comment") {
    editor.related.comment_to = JSON.parse(JSON.stringify(props.item))
  }
}

function deletePost() {
  editor.related.delete_to = JSON.parse(JSON.stringify(props.item))
  editor.related.delete_to.model_type = props.item.model_type + "s"
  editor.show.delete = true
}
</script>
