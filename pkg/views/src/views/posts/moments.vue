<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="content flex-grow-1">
      <v-card :loading="loading">
        <article>
          <v-card-text>
            <div class="flex justify-between px-3">
              <div class="flex gap-1">
                <v-avatar
                  color="grey-lighten-2"
                  icon="mdi-account-circle"
                  class="rounded-card me-2"
                  :image="post?.author.avatar"
                />

                <div>
                  <p class="font-bold">{{ post?.author.nick }}</p>
                  <p class="opacity-80">
                    {{ post?.author.description ? post?.author.description : "No description yet." }}
                  </p>
                </div>
              </div>

              <div>
                <post-action :item="post" />
              </div>
            </div>

            <v-divider class="mb-5 mt-3.5 mx-[-16px] border-opacity-50" />

            <div class="px-3">
              <moment-content :item="post" content-only />
            </div>

            <div class="mt-3 px-2">
              <post-attachment v-if="post?.attachments" :attachments="post?.attachments" />
            </div>

            <v-divider class="my-5 mx-[-16px] border-opacity-50" />

            <div class="px-3">
              <post-reaction
                model="moments"
                :item="post"
                :reactions="post?.reaction_list ?? {}"
                @update="updateReactions"
              />
            </div>
          </v-card-text>
        </article>
      </v-card>
    </div>

    <div class="aside md:sticky top-0 w-full h-fit w-full md:max-w-[380px] md:min-w-[360px]">
      <v-card title="Comments">
        <div class="px-[1rem] pb-[0.825rem] mt-[-12px]">
          <comment-list
            model="moment"
            dataset="moments"
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
import { request } from "@/scripts/request"
import { useRoute } from "vue-router"
import MomentContent from "@/components/posts/MomentContent.vue"
import PostReaction from "@/components/posts/PostReaction.vue"
import CommentList from "@/components/comments/CommentList.vue"
import PostAttachment from "@/components/posts/PostAttachment.vue"
import PostAction from "@/components/posts/PostAction.vue"

const loading = ref(false)
const error = ref<string | null>(null)

const post = ref<any>(null)
const comments = ref<any[]>([])

const route = useRoute()

async function readPost() {
  loading.value = true
  const res = await request(`/api/p/moments/${route.params.alias}`)
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
  // eslint-disable-next-line
  if (post.value.reaction_list.hasOwnProperty(symbol)) {
    post.value.reaction_list[symbol] += num
  } else {
    post.value.reaction_list[symbol] = num
  }
}
</script>

<style scoped>
.rounded-card {
  border-radius: 8px;
}
</style>
