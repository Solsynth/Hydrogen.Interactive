<template>
  <div>
    <article class="prose prose-moment" v-html="parseContent(props.item?.content ?? '')" />

    <div v-if="props.brief" class="my-1">
      <v-btn
        append-icon="mdi-arrow-right"
        variant="tonal"
        size="x-small"
        rounded="sm"
        :to="{ name: 'posts.details.moments', params: { alias: props.item?.alias ?? 'not-found' } }"
      >
        More
      </v-btn>
    </div>
  </div>
</template>

<script setup lang="ts">
import dompurify from "dompurify"
import { parse } from "marked"

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
