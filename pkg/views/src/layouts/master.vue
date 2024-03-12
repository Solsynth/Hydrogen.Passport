<template>
  <v-app-bar height="64" color="primary" scroll-behavior="elevate" flat>
    <div class="max-md:px-5 md:px-12 flex flex-grow-1 items-center">
      <router-link :to="{ name: 'dashboard' }">
        <h2 class="ml-2 text-lg font-500">Solarpass</h2>
      </router-link>

      <v-spacer />

      <v-menu>
        <template #activator="{ props }">
          <v-btn flat exact v-bind="props" icon>
            <v-avatar color="transparent" icon="mdi-account-circle" :src="id.userinfo.data?.avatar" />
          </v-btn>
        </template>

        <v-list density="compact">
          <v-list-item title="Sign in" prepend-icon="mdi-login-variant" />
          <v-list-item title="Create account" prepend-icon="mdi-account-plus" />
        </v-list>
      </v-menu>
    </div>
  </v-app-bar>

  <v-main>
    <router-view />
  </v-main>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useUserinfo } from "@/stores/userinfo"

const id = useUserinfo()

const username = computed(() => {
  if (id.userinfo.isLoggedIn) {
    return "@" + id.userinfo.data?.name
  } else {
    return "@vistor"
  }
})
const nickname = computed(() => {
  if (id.userinfo.isLoggedIn) {
    return id.userinfo.data?.nick
  } else {
    return "Anonymous"
  }
})

id.readProfiles()
</script>

<style scoped>
.editor-fab {
  position: fixed !important;
  bottom: 16px;
  right: 20px;
}
</style>
