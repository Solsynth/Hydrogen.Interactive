<template>
  <v-container class="flex max-md:flex-col gap-3 overflow-auto max-h-[calc(100vh-64px)] no-scrollbar">
    <div class="timeline flex-grow-1 mt-[-16px]">
      <post-list :loading="loading" :posts="posts" :loader="readMore" />
    </div>

    <div class="aside sticky top-0 w-full h-fit md:min-w-[280px] md:max-w-[320px]">
      <v-card title="Categories">
        <v-list density="compact">
        </v-list>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import PostList from "@/components/posts/PostList.vue";
import { reactive, ref } from "vue";
import { request } from "@/scripts/request";

const error = ref<string | null>(null);
const loading = ref(false);
const pagination = reactive({ page: 1, pageSize: 10, total: 0 });

const posts = ref<any[]>([]);

async function readPosts() {
  loading.value = true;
  const res = await request(`/api/posts?` + new URLSearchParams({
    take: pagination.pageSize.toString(),
    offset: ((pagination.page - 1) * pagination.pageSize).toString()
  }));
  if (res.status !== 200) {
    loading.value = false;
    error.value = await res.text();
  } else {
    error.value = null;
    loading.value = false;
    const data = await res.json();
    pagination.total = data["count"];
    posts.value.push(...data["data"]);
  }
}

async function readMore({ done }: any) {
  // Reach the end of data
  if (pagination.total <= pagination.page * pagination.pageSize) {
    done("empty");
    return;
  }

  pagination.page++;
  await readPosts();

  done("ok");
}

readPosts();
</script>
