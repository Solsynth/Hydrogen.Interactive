<template>
  <v-dialog v-model="editor.show.moment" class="max-w-[540px]">
    <moment-editor />
  </v-dialog>
  <v-dialog v-model="editor.show.comment" class="max-w-[540px]">
    <comment-editor />
  </v-dialog>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor"
import MomentEditor from "@/components/publish/MomentEditor.vue"
import CommentEditor from "@/components/publish/CommentEditor.vue";

const editor = useEditor()
</script>
