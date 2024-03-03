<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="timeline flex-grow-1 max-w-[75ch]">
      <v-card :loading="loading">
        <article>
          <v-card-title>{{ post?.title }}</v-card-title>

          <v-card-text>
            <div class="text-sm">{{ post?.description }}</div>

            <v-divider class="mt-5 mx-[-16px] border-opacity-50" />

            <article-content :item="post" content-only />
          </v-card-text>
        </article>
      </v-card>
    </div>

    <div class="aside sticky top-0 w-full h-fit md:min-w-[280px]">
      <v-card title="Comments">
        <v-list density="compact">
        </v-list>
      </v-card>

      <v-card title="Reactions" class="mt-3">
        <v-list density="compact">
        </v-list>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { request } from "@/scripts/request";
import { useRoute } from "vue-router";
import ArticleContent from "@/components/posts/ArticleContent.vue";

const loading = ref(false);
const error = ref<string | null>(null);
const post = ref<any>(null);

const route = useRoute();

async function readPost() {
  loading.value = true;
  const res = await request(`/api/${route.params.postType}/${route.params.alias}?`);
  if (res.status !== 200) {
    error.value = await res.text();
  } else {
    error.value = null;
    post.value = await res.json();
  }
  loading.value = false;
}

readPost();
</script>