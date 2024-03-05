<template>
  <article class="prose prose-moment" v-html="parseContent(props.item.content)" />
</template>

<script setup lang="ts">
import dompurify from "dompurify";
import { parse } from "marked";

const props = defineProps<{ item: any }>();

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string);
}
</script>

<style>
.prose.prose-moment p {
  margin: 0 !important;
}
</style>
