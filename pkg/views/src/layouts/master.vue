<template>
  <v-navigation-drawer v-model="drawerOpen" color="grey-lighten-5" floating>
    <v-list density="compact" nav> </v-list>
  </v-navigation-drawer>

  <v-app-bar height="64" color="primary" scroll-behavior="elevate" flat>
    <div class="max-md:px-5 md:px-12 flex flex-grow-1 items-center">
      <v-app-bar-nav-icon variant="text" @click.stop="toggleDrawer" />

      <router-link :to="{ name: 'explore' }">
        <h2 class="ml-2 text-lg font-500">Goatplaza</h2>
      </router-link>

      <v-spacer />

      <v-tooltip v-for="item in navigationMenu" :text="item.name" location="bottom">
        <template #activator="{ props }">
          <v-btn flat exact v-bind="props" :to="{ name: item.to }" size="small" :icon="item.icon" />
        </template>
      </v-tooltip>
    </div>
  </v-app-bar>

  <v-main>
    <router-view />
  </v-main>

  <v-menu
    open-on-hover
    open-on-click
    :open-delay="0"
    :close-delay="0"
    location="top"
    transition="scroll-y-reverse-transition"
  >
    <template v-slot:activator="{ props }">
      <v-fab v-bind="props" class="editor-fab" icon="mdi-pencil" color="primary" size="64" appear />
    </template>

    <div class="flex flex-col items-center gap-4 mb-4">
      <v-btn variant="elevated" color="secondary" icon="mdi-newspaper-variant" @click="editor.show.article = true" />
      <v-btn variant="elevated" color="accent" icon="mdi-camera-iris" @click="editor.show.moment = true" />
    </div>
  </v-menu>

  <post-action />
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useEditor } from "@/stores/editor"
import PostAction from "@/components/publish/PostAction.vue"

const editor = useEditor()
const navigationMenu = [{ name: "Explore", icon: "mdi-compass", to: "explore" }]

const drawerOpen = ref(true)

function toggleDrawer() {
  drawerOpen.value = !drawerOpen.value
}
</script>

<style scoped>
.editor-fab {
  position: fixed !important;
  bottom: 16px;
  right: 20px;
}
</style>
