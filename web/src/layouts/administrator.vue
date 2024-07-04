<template>
  <app-bar title="Solarpass Administration" />

  <v-main>
    <router-view />
  </v-main>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"
import { useRouter } from "vue-router"
import { onMounted } from "vue"
import AppBar from "@/components/navigation/AppBar.vue"

const id = useUserinfo()
const router = useRouter()

onMounted(async () => {
  await id.readProfiles()
  if (!id.userinfo.data.perm_nodes["AdminView"]) {
    await router.push({ name: "dashboard" })
  }
})
</script>

<style scoped>
.icon-filter {
  filter: invert(100%) sepia(100%) saturate(14%) hue-rotate(212deg) brightness(104%) contrast(104%);
}
</style>
