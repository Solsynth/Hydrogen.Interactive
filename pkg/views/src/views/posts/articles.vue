<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="timeline flex-grow-1">
      <v-card :loading="loading">
        <article>
          <v-card-text>
            <h1 class="text-lg font-medium">{{ post?.title }}</h1>
            <p class="text-sm">{{ post?.description }}</p>

            <v-divider class="mt-5 mx-[-16px] border-opacity-50" />

            <article-content :item="post" content-only />

            <v-divider class="my-5 mx-[-16px] border-opacity-50" />

            <div class="px-3">
              <post-reaction
                :item="post"
                :model="route.params.postType"
                :reactions="post?.reaction_list ?? {}"
                @update="updateReactions"
              />
            </div>
          </v-card-text>
        </article>
      </v-card>
    </div>

    <div class="aside sticky top-0 w-full h-fit md:max-w-[360px] md:min-w-[280px]">
      <v-card title="Comments">
        <div class="px-[1rem] pb-[0.825rem] mt-[-12px]">
          <comment-list
            model="article"
            dataset="articles"
            v-model:comments="comments"
            :item="post"
            :alias="route.params.alias"
          />
        </div>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { request } from "@/scripts/request"
import ArticleContent from "@/components/posts/ArticleContent.vue"
import PostReaction from "@/components/posts/PostReaction.vue"
import CommentList from "@/components/comments/CommentList.vue"
import { useRoute } from "vue-router"

const loading = ref(false)
const error = ref<string | null>(null)

const post = ref<any>(null)
const comments = ref<any[]>([])

const route = useRoute()

async function readPost() {
  loading.value = true
  const res = await request(`/api/p/articles/${route.params.alias}`)
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    error.value = null
    post.value = await res.json()
  }
  loading.value = false
}

readPost()

function updateReactions(symbol: string, num: number) {
  if (post.value.reaction_list.hasOwnProperty(symbol)) {
    post.value.reaction_list[symbol] += num
  } else {
    post.value.reaction_list[symbol] = num
  }
}
</script>
