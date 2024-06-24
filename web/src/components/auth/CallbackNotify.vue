<template>
  <div class="w-full max-w-[720px]">
    <v-expand-transition>
      <v-alert v-show="route.query['redirect_uri']" variant="tonal" type="info" class="text-xs">
        You need to sign in before access that page. After you signed in, we will redirect you to: <br />
        <span class="font-mono">{{ route.query["redirect_uri"] }}</span>
      </v-alert>
    </v-expand-transition>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router"

const route = useRoute()
</script>
