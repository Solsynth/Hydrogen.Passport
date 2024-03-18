<template>
  <v-app-bar height="64" color="primary" scroll-behavior="elevate" flat>
    <div class="max-md:px-5 md:px-12 flex flex-grow-1 items-center">
      <router-link :to="{ name: 'dashboard' }" class="flex gap-1">
        <img src="/favicon.png" width="24" height="24" class="icon-filter" />
        <h2 class="ml-2 text-lg font-500">Solarpass</h2>
      </router-link>

      <v-spacer />

      <div class="me-2">
        <notification-list />
      </div>

      <div>
        <user-menu />
      </div>
    </div>
  </v-app-bar>

  <v-main>
    <router-view />
  </v-main>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"
import NotificationList from "@/components/NotificationList.vue"
import UserMenu from "@/components/UserMenu.vue"

const id = useUserinfo()

id.readProfiles()
</script>

<style scoped>
.editor-fab {
  position: fixed !important;
  bottom: 16px;
  right: 20px;
}

.icon-filter {
  filter: invert(100%) sepia(100%) saturate(14%) hue-rotate(212deg) brightness(104%) contrast(104%);
}
</style>
