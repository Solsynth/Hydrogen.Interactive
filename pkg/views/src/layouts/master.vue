<template>
  <v-navigation-drawer v-model="drawerOpen" color="grey-lighten-5" floating>
    <div class="d-flex text-center justify-center items-center h-[64px]">
      <h1>Goatplaza</h1>
    </div>
  </v-navigation-drawer>

  <v-app-bar height="64" color="primary" scroll-behavior="elevate" flat>
    <div class="container mx-auto px-5">
      <v-app-bar-nav-icon variant="text" @click.stop="toggleDrawer"></v-app-bar-nav-icon>
    </div>
  </v-app-bar>

  <v-main>
    <router-view />
  </v-main>
</template>

<script setup>
import { ref } from "vue"

const drawerOpen = ref(true)

function toggleDrawer() {
  drawerOpen.value = !drawerOpen.value
}
</script>
