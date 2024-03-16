<template>
  <div>
    <v-card class="px-1 mb-3">
      <v-card-text class="flex gap-3.5">
        <v-avatar
          color="grey-lighten-2"
          icon="mdi-account-circle"
          class="rounded-card"
          :size="54"
          :image="'/api/avatar/' + id.userinfo.data.avatar"
        />
        <div>
          <h1 class="text-2xl">{{ id.userinfo.displayName }}</h1>
          <p>What can I help you today?</p>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"

const id = useUserinfo()
</script>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
