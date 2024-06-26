<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn flat exact v-bind="props" icon>
        <v-avatar color="transparent" icon="mdi-account-circle" :image="'/api/avatar/' + id.userinfo.data?.avatar" />
      </v-btn>
    </template>

    <v-list density="compact" v-if="!id.userinfo.isLoggedIn">
      <v-list-item title="Sign in" prepend-icon="mdi-login-variant" :to="{ name: 'auth.sign-in' }" />
      <v-list-item title="Create account" prepend-icon="mdi-account-plus" :to="{ name: 'auth.sign-up' }" />
    </v-list>
    <v-list density="compact" v-else>
      <v-list-item :title="nickname" :subtitle="username" />

      <v-divider class="border-opacity-50 my-2" />

      <v-list-item title="Dashboard" prepend-icon="mdi-account-supervisor" exact :to="{ name: 'dashboard' }" />
      <v-list-item title="Sign out" prepend-icon="mdi-logout" @click="signout"></v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { defaultUserinfo, useUserinfo } from "@/stores/userinfo"
import { computed } from "vue"
import Cookie from "universal-cookie"

const id = useUserinfo()

const username = computed(() => {
  if (id.userinfo.isLoggedIn) {
    return "@" + id.userinfo.data?.name
  } else {
    return "@visitor"
  }
})
const nickname = computed(() => {
  if (id.userinfo.isLoggedIn) {
    return id.userinfo.data?.nick
  } else {
    return "Anonymous"
  }
})

function signout() {
  const ck = new Cookie();
  ck.remove("__hydrogen_atk");
  ck.remove("__hydrogen_rtk")
  id.userinfo = defaultUserinfo
  window.location.reload()
}
</script>
