<template>
  <component
    :is="props.brief ? RouterLink : 'div'"
    :to="{ name: 'posts.details.moments', params: { alias: props.item?.alias } }"
  >
    <article class="prose prose-moment" v-html="parseContent(props.item?.content ?? '')" />
  </component>
</template>

<script setup lang="ts">
import dompurify from "dompurify"
import { parse } from "marked"
import { RouterLink } from "vue-router"

const props = defineProps<{ item: any; brief?: boolean }>()

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string)
}
</script>

<style>
.prose.prose-moment p {
  margin: 0 !important;
}
</style>
