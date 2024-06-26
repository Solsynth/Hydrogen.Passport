<template>
  <app-bar />

  <v-main>
    <router-view />
  </v-main>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"
import AppBar from "@/components/navigation/AppBar.vue"

const id = useUserinfo()

id.readProfiles()
</script>

<style scoped>
.icon-filter {
  filter: invert(100%) sepia(100%) saturate(14%) hue-rotate(212deg) brightness(104%) contrast(104%);
}
</style>
