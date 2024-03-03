<template>
  <v-card :loading="props.loading">
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

        <div class="flex-grow-1">
          <div class="font-bold">{{ props.item?.author.nick }}</div>

          <div v-if="props.item?.modal_type === 'article'" class="text-xs text-grey-darken-4 mb-2">Published an article</div>

          <component :is="renderer[props.item?.model_type]" v-bind="props" />
        </div>
      </div>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import type { Component } from "vue";
import ArticleContent from "@/components/posts/ArticleContent.vue";
import MomentContent from "@/components/posts/MomentContent.vue";

const props = defineProps<{ item: any, brief?: boolean, loading?: boolean }>();

const renderer: { [id: string]: Component } = {
  article: ArticleContent,
  moment: MomentContent
};
</script>

<style scoped>
.rounded-card {
  border-radius: 8px;
}
</style>
