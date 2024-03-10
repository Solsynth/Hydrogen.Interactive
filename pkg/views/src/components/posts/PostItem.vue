<template>
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

      <div v-if="props.item?.model_type === 'article'" class="text-xs text-grey-darken-4 mb-2">
        Published an article
      </div>

      <component :is="renderer[props.item?.model_type]" v-bind="props" />

      <post-attachment v-if="props.item?.attachments" :attachments="props.item?.attachments" />

      <post-reaction
        size="small"
        :item="props.item"
        :model="props.item?.model_type ? props.item?.model_type + 's' : 'articles'"
        :reactions="props.item?.reaction_list ?? {}"
        @update="updateReactions"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from "vue"
import ArticleContent from "@/components/posts/ArticleContent.vue"
import MomentContent from "@/components/posts/MomentContent.vue"
import CommentContent from "@/components/posts/CommentContent.vue"
import PostAttachment from "./PostAttachment.vue"
import PostReaction from "@/components/posts/PostReaction.vue"

const props = defineProps<{ item: any; brief?: boolean }>()
const emits = defineEmits(["update:item"])

const renderer: { [id: string]: Component } = {
  article: ArticleContent,
  moment: MomentContent,
  comment: CommentContent
}

function updateReactions(symbol: string, num: number) {
  const item = JSON.parse(JSON.stringify(props.item))
  if (item.reaction_list == null) {
    item.reaction_list = {}
  }
  if (item.reaction_list.hasOwnProperty(symbol)) {
    item.reaction_list[symbol] += num
  } else {
    item.reaction_list[symbol] = num
  }
  emits("update:item", item)
}
</script>

<style scoped>
.rounded-card {
  border-radius: 8px;
}
</style>
