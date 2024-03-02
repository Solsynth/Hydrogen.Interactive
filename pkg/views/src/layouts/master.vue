<template>
  <v-navigation-drawer v-model="drawerOpen" color="grey-lighten-5" floating>
    <v-list density="compact" nav>
    </v-list>
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
          <v-btn flat v-bind="props" :to="{ name: item.to }" size="small" :icon="item.icon" />
        </template>
      </v-tooltip>
    </div>
  </v-app-bar>

  <v-main>
    <router-view />
  </v-main>

  <v-fab
    class="editor-fab"
    icon="mdi-pencil"
    color="primary"
    size="64"
    appear
    @click="editor.show = true"
  />

  <post-editor />
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useEditor } from "@/stores/editor";
import PostEditor from "@/components/publish/PostEditor.vue";

const editor = useEditor()
const navigationMenu = [
  { name: "Explore", icon: "mdi-compass", to: "explore" }
];

const drawerOpen = ref(true);

function toggleDrawer() {
  drawerOpen.value = !drawerOpen.value;
}
</script>

<style scoped>
.editor-fab {
  position: fixed !important;
  bottom: 16px;
  right: 20px;
}
</style>
