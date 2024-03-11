<template>
  <div>
    <section v-if="!props.contentOnly" class="mb-2">
      <h1 class="text-lg font-bold">{{ props.item?.title }}</h1>
      <div class="text-sm">{{ props.item?.description }}</div>
    </section>

    <div v-if="props.brief" class="mt-2">
      <v-btn
        append-icon="mdi-arrow-right"
        variant="tonal"
        size="small"
        rounded="sm"
        :to="{ name: 'posts.details.articles', params: { alias: props.item?.alias ?? 'not-found' } }"
      >
        Read more
      </v-btn>
    </div>

    <div v-else>
      <article class="prose max-w-none" v-html="parseContent(props.item?.content ?? '')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import dompurify from "dompurify"
import { parse } from "marked"

const props = defineProps<{ item: any; brief?: boolean; contentOnly?: boolean }>()

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string)
}
</script>
