<template>
  <div class="post-list">
    <div v-if="props.loading" class="text-center py-8">
      <v-progress-circular indeterminate />
    </div>

    <v-infinite-scroll :items="props.posts" :onLoad="props.loader">
      <template v-for="item in props.posts" :key="item">
        <div class="mb-3 px-1">
          <post-item :item="item" brief />
        </div>
      </template>
    </v-infinite-scroll>
  </div>
</template>

<script setup lang="ts">
import PostItem from "@/components/posts/PostItem.vue";

const props = defineProps<{ loading: boolean, posts: any[], loader: (opts: any) => Promise<any> }>();
</script>
