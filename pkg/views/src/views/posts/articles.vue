<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="content flex-grow-1">
      <v-card :loading="loading">
        <article>
          <v-card-text>
            <div class="flex justify-between px-3">
              <div>
                <h1 class="text-lg font-medium">{{ post?.title }}</h1>
                <p class="text-sm">{{ post?.description }}</p>
              </div>

              <div>
                <post-action :item="post" />
              </div>
            </div>

            <v-divider class="my-5 mx-[-16px] border-opacity-50" />

            <div class="px-3 text-xs opacity-80 flex gap-1">
              <span>Written by {{ post?.author?.nick }}</span>
              <span>·</span>
              <span>Published at {{ new Date(post?.created_at).toLocaleString() }}</span>
            </div>

            <v-divider class="mt-5 mx-[-16px] border-opacity-50" />

            <div class="px-3">
              <article-content :item="post" content-only />
            </div>

            <v-divider class="my-5 mx-[-16px] border-opacity-50" />

            <div class="px-3">
              <post-reaction
                model="articles"
                :item="post"
                :reactions="post?.reaction_list ?? {}"
                @update="updateReactions"
              />
            </div>
          </v-card-text>
        </article>
      </v-card>
    </div>

    <div class="aside sticky top-0 w-full h-fit w-full md:max-w-[380px] md:min-w-[360px]">
      <v-card title="Comments">
        <div class="px-[1rem] pb-[0.825rem] mt-[-12px]">
          <comment-list
            model="article"
            dataset="articles"
            :item="post"
            :alias="route.params.alias"
            v-model:comments="comments"
          />
        </div>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRoute } from "vue-router"
import { request } from "@/scripts/request"
import ArticleContent from "@/components/posts/ArticleContent.vue"
import PostReaction from "@/components/posts/PostReaction.vue"
import PostAction from "@/components/posts/PostAction.vue"
import CommentList from "@/components/comments/CommentList.vue"

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
