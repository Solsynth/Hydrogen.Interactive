<template>
  <v-card variant="tonal" class="max-w-[540px]" :ripple="canLightbox" @click="openLightbox">
    <div class="content">
      <img v-if="current.type === 1" :src="getUrl(current)" />
      <video v-if="current.type === 2" controls class="w-full">
        <source :src="getUrl(current)"></source>
      </video>
    </div>

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

const props = defineProps<{ attachments: any[] }>()

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
