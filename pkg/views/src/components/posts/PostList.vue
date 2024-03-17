<template>
  <div class="post-list">
    <v-infinite-scroll :items="props.posts" :onLoad="props.loader">
      <template v-for="(item, idx) in props.posts" :key="item">
        <div class="mb-3 px-1">
          <v-card>
            <template #text>
              <post-item brief :item="item" @update:item="val => updateItem(idx, val)" />
            </template>
          </v-card>
        </div>
      </template>
    </v-infinite-scroll>
  </div>
</template>

<script setup lang="ts">
import PostItem from "@/components/posts/PostItem.vue";

const props = defineProps<{ posts: any[], loader: (opts: any) => Promise<any> }>();
const emits = defineEmits(["update:posts"]);

function updateItem(idx: number, data: any) {
  const posts = JSON.parse(JSON.stringify(props.posts));
  posts[idx] = data;
  emits("update:posts", posts);
}
</script>
