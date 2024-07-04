<template>
  <app-bar>
    <template #extension>
      <v-tabs align-tabs="title" color="white">
        <v-tab text="Dashboard" prepend-icon="mdi-view-dashboard" :to="{ name: 'dashboard' }" exact />
        <v-tab text="Personalize" prepend-icon="mdi-card-bulleted-outline" :to="{ name: 'personalize' }" exact />
        <v-tab text="Security" prepend-icon="mdi-security" :to="{ name: 'security' }" exact />
      </v-tabs>
    </template>
  </app-bar>

  <v-main>
    <v-container class="pt-6 px-6 p-container">
      <router-view />
    </v-container>

    <Copyright />
  </v-main>
</template>

<script setup lang="ts">
import AppBar from "@/components/navigation/AppBar.vue"
import Copyright from "@/components/Copyright.vue"
</script>

<style scoped>
.p-container {
  max-width: 64rem;
}
</style>