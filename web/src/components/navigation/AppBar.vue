<template>
  <v-app-bar height="64" color="primary" scroll-behavior="elevate" flat>
    <div class="max-md:px-5 md:px-12 flex flex-grow-1 items-center">
      <router-link :to="{ name: 'dashboard' }" class="flex gap-1">
        <img src="/favicon.png" alt="logo" width="27" height="24" class="icon-filter" />
        <h2 class="ml-2 text-lg font-500">Solarpass</h2>
      </router-link>

      <v-spacer />

      <div class="me-2">
        <v-btn icon size="small" variant="text" @click="openNotify = !openNotify">
          <v-badge v-if="notify.total > 0" color="error" :content="notify.total">
            <v-icon icon="mdi-bell" />
          </v-badge>

          <v-icon v-else icon="mdi-bell" />
        </v-btn>
      </div>

      <div>
        <user-menu />
      </div>
    </div>

    <template #extension>
      <slot name="extension" />
    </template>
  </v-app-bar>

  <NotificationList v-model:open="openNotify" />
</template>

<script setup lang="ts">
import NotificationList from "@/components/NotificationList.vue"
import UserMenu from "@/components/UserMenu.vue"
import { useNotifications } from "@/stores/notifications"
import { ref } from "vue"

const notify = useNotifications()

const openNotify = ref(false)
</script>

<style scoped>
.icon-filter {
  filter: invert(100%) sepia(100%) saturate(14%) hue-rotate(212deg) brightness(104%) contrast(104%);
}
</style>
