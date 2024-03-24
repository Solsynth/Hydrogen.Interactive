<template>
  <v-chip size="small" variant="tonal" prepend-icon="mdi-paperclip" v-if="props.overview">
    Attached {{ props.attachments.length }} attachment(s)
  </v-chip>

  <v-card v-else variant="outlined" class="max-w-[540px]">
    <v-carousel hide-delimiters progress="primary" show-arrows="hover" height="100%">
      <v-carousel-item v-for="item in attachments">
        <img v-if="item.type === 1" :src="getUrl(item)" class="cursor-zoom-in" @click="openLightbox" />
        <video v-if="item.type === 2" controls class="w-full">
          <source :src="getUrl(item)" />
        </video>
        <div v-if="item.type === 3" class="w-full px-7 py-12">
          <audio controls :src="getUrl(item)" class="mx-auto"></audio>
        </div>
      </v-carousel-item>
    </v-carousel>

    <vue-easy-lightbox teleport="#app" :visible="lightbox" :imgs="[getUrl(current)]" @hide="lightbox = false">
      <template v-slot:close-btn="{ close }">
        <v-btn class="fixed left-2 top-2" icon="mdi-close" variant="text" color="white" @click="close" />
      </template>
    </vue-easy-lightbox>
  </v-card>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import VueEasyLightbox from "vue-easy-lightbox"

const props = defineProps<{ attachments: any[]; overview?: boolean }>()

const lightbox = ref(false)
const focus = ref(0)

const current = computed(() => props.attachments[focus.value])
const canLightbox = computed(() => current.value.type === 1)

function getUrl(item: any) {
  return item.external_url ? item.external_url : `/api/attachments/o/${item.file_id}`
}

function openLightbox() {
  if (canLightbox.value) {
    lightbox.value = true
  }
}
</script>

<style>
.vel-model {
  z-index: 10;
}
</style>
