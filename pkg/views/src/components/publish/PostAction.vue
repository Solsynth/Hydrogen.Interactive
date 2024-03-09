<template>
  <v-dialog v-model="editor.show.comment" class="max-w-[540px]" eager>
    <comment-editor />
  </v-dialog>
  <v-dialog v-model="editor.show.moment" class="max-w-[540px]" eager>
    <moment-editor />
  </v-dialog>
  <v-dialog v-model="editor.show.article" transition="dialog-bottom-transition" fullscreen eager>
    <article-editor />
  </v-dialog>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor"
import MomentEditor from "@/components/publish/MomentEditor.vue"
import CommentEditor from "@/components/publish/CommentEditor.vue";
import ArticleEditor from "@/components/publish/ArticleEditor.vue";

const editor = useEditor()
</script>
