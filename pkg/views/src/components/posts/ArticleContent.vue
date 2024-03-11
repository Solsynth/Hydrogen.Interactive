<template>
  <component
    :is="props.brief ? RouterLink : 'div'"
    :to="{ name: 'posts.details.articles', params: { alias: props.item?.alias } }"
  >
    <section v-if="!props.contentOnly" class="mb-2">
      <h1 class="text-lg font-bold">{{ props.item?.title }}</h1>
      <div class="text-sm">{{ props.item?.description }}</div>
    </section>

    <div v-else>
      <article class="prose max-w-none" v-html="parseContent(props.item?.content ?? '')" />
    </div>
  </component>
</template>

<script setup lang="ts">
import dompurify from "dompurify"
import { parse } from "marked"
import { RouterLink } from "vue-router"

const props = defineProps<{ item: any; brief?: boolean; contentOnly?: boolean }>()

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string)
}
</script>
