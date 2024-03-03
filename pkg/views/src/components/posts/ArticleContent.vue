<template>
  <div>
    <section v-if="!props.contentOnly" class="mb-2">
      <h1 class="text-lg font-bold">{{ props.item?.title }}</h1>
      <div class="text-sm">{{ props.item?.description }}</div>
    </section>

    <div v-if="props.brief">
      <router-link
        :to="{ name: 'posts.details', params: { postType: 'articles', alias: props.item?.alias ?? 'not-found' } }"
        append-icon="mdi-arrow-right"
        class="link underline text-primary font-medium"
      >
        Read more...
      </router-link>
    </div>

    <div v-else>
      <article class="prose max-w-none" v-html="parseContent(props.item?.content ?? '')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import dompurify from "dompurify";
import { parse } from "marked";

const props = defineProps<{ item: any, brief?: boolean, contentOnly?: boolean }>();

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string);
}
</script>