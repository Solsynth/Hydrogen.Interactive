<template>
  <div class="flex gap-3">
    <div>
      <v-avatar
        color="grey-lighten-2"
        icon="mdi-account-circle"
        class="rounded-card"
        :image="props.item?.author.avatar"
      />
    </div>

    <div class="flex-grow-1 relative">
      <div class="font-bold">{{ props.item?.author.nick }}</div>

      <div v-if="props.item?.model_type === 'article'" class="text-xs text-grey-darken-4 mb-2">
        Published an article
      </div>

      <component :is="renderer[props.item?.model_type]" v-bind="props" />

      <post-attachment
        v-if="props.item?.attachments"
        class="mt-1.5"
        :overview="props.item?.model_type !== 'moment'"
        :attachments="props.item?.attachments"
      />

      <post-reaction
        size="small"
        :item="props.item"
        :model="props.item?.model_type ? props.item?.model_type + 's' : 'articles'"
        :reactions="props.item?.reaction_list ?? {}"
        @update="updateReactions"
      />

      <div class="mt-1 text-xs opacity-80 flex gap-2 items-center">
        <span>Posted at {{ new Date(props.item?.created_at).toLocaleString() }}</span>
      </div>

      <v-menu>
        <template #activator="{ props }">
          <div class="absolute right-0 top-0">
            <v-btn v-bind="props" icon="mdi-dots-vertical" variant="text" size="x-small" />
          </div>
        </template>

        <v-list density="compact" lines="one">
          <v-list-item disabled append-icon="mdi-flag" title="Report" />
          <v-list-item v-if="isOwned" append-icon="mdi-pencil" title="Edit" @click="editPost" />
          <v-list-item v-if="isOwned" append-icon="mdi-delete" title="Delete" @click="deletePost" />
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from "vue"
import { useUserinfo } from "@/stores/userinfo"
import { useEditor } from "@/stores/editor"
import ArticleContent from "@/components/posts/ArticleContent.vue"
import MomentContent from "@/components/posts/MomentContent.vue"
import CommentContent from "@/components/posts/CommentContent.vue"
import PostAttachment from "./PostAttachment.vue"
import PostReaction from "@/components/posts/PostReaction.vue"

const id = useUserinfo()

const props = defineProps<{ item: any; brief?: boolean }>()
const emits = defineEmits(["update:item"])

const editor = useEditor()

const renderer: { [id: string]: Component } = {
  article: ArticleContent,
  moment: MomentContent,
  comment: CommentContent
}

const isOwned = computed(() => props.item?.author_id === id.userinfo.data.id)

function editPost() {
  editor.related.edit_to = props.item
  if (editor.show.hasOwnProperty(props.item.model_type)) {
    // @ts-ignore
    editor.show[props.item.model_type] = true
  }
  if (props.item.model_type === "comment") {
    editor.related.comment_to = props.item
  }
}

function deletePost() {
  editor.related.delete_to = JSON.parse(JSON.stringify(props.item))
  editor.related.delete_to.model_type = props.item.model_type + "s"
  editor.show.delete = true
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
