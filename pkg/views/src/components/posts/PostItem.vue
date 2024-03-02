<template>
  <v-card>
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

        <div>
          <div class="font-bold">{{ props.item?.author.nick }}</div>
          {{ props.item?.content }}
        </div>
      </div>
    </template>
  </v-card>
</template>

<script setup lang="ts">
const props = defineProps<{ item: any }>();
</script>

<style scoped>
.rounded-card {
  border-radius: 8px;
}
</style>