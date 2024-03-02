<template>
  <v-card>
    <template #text>
      <div class="flex gap-3">
        <div>
          <v-avatar
            color="grey-lighten-2"
            icon="mdi-account-circle"
            class="rounded-card"
            :src="props.item?.author.avatar"
          />
        </div>

        <div>
          <div class="font-bold">{{ props.item?.author.nick }}</div>
          <div class="prose prose-post" v-html="parseContent(props.item.content)"></div>
        </div>
      </div>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import dompurify from "dompurify";
import { parse } from "marked";

const props = defineProps<{ item: any }>();

function parseContent(src: string): string {
  return dompurify().sanitize(parse(src) as string);
}
</script>

<style scoped>
.rounded-card {
  border-radius: 8px;
}
</style>

<style>
.prose.prose-post, p {
  margin: 0 !important;
}
</style>