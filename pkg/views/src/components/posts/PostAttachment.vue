<template>
  <v-chip size="small" variant="tonal" prepend-icon="mdi-paperclip" v-if="props.overview">
    Attached {{ props.attachments.length }} attachment(s)
  </v-chip>

  <v-card v-else variant="outlined" class="max-w-[540px]" :ripple="canLightbox" @click="openLightbox">
    <div class="content">
      <img v-if="current.type === 1" :src="getUrl(current)" />
      <video v-if="current.type === 2" controls class="w-full">
        <source :src="getUrl(current)" />
      </video>
    </div>

    <div v-if="props.attachments.length > 1" class="switcher flex justify-between items-center px-5 py-2">
      <div>{{ focus + 1 }} of {{ props.attachments.length }}</div>
      <div>
        <v-btn icon="mdi-arrow-left" variant="text" size="small" @click.stop="focusPrev" />
        <v-btn icon="mdi-arrow-right" variant="text" size="small" @click.stop="focusNext" />
      </div>
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

function focusNext() {
  if (focus.value + 1 < props.attachments.length) {
    focus.value++
  }
}

function focusPrev() {
  if (focus.value - 1 >= 0) {
    focus.value--
  }
}
</script>

<style>
.vel-model {
  z-index: 10;
}
</style>
