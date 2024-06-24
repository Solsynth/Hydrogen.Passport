<template>
  <v-container class="h-screen flex flex-col gap-3 items-center justify-center">
    <callback-notify />

    <v-card class="w-full max-w-[720px] overflow-auto" :loading="loading">
      <v-card-text class="card-grid pa-9">
        <div>
          <v-avatar color="accent" icon="mdi-login-variant" size="large" class="card-rounded mb-2" />
          <h1 class="text-2xl">Sign in</h1>
          <p v-if="ticket">We need to verify that the person trying to access your account is you.</p>
          <p v-else>Sign in via your Solar ID to access the entire Solar Network.</p>
        </div>

        <v-window :touch="false" :model-value="panel" class="pa-2 mx-[-0.5rem]">
          <v-window-item v-for="(k, idx) in Object.keys(panels)" :key="idx" :value="k">
            <component :is="panels[k]" @swap="(val: string) => (panel = val)" v-model:loading="loading"
                       v-model:currentFactor="currentFactor" v-model:ticket="ticket" />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <copyright />
  </v-container>
</template>

<script setup lang="ts">
import { type Component, ref } from "vue"
import Copyright from "@/components/Copyright.vue"
import CallbackNotify from "@/components/auth/CallbackNotify.vue"
import FactorPicker from "@/components/auth/FactorPicker.vue"
import FactorApplicator from "@/components/auth/FactorApplicator.vue"
import AccountAuthenticate from "@/components/auth/Authenticate.vue"
import AuthenticateCompleted from "@/components/auth/AuthenticateCompleted.vue"

const loading = ref(false)

const currentFactor = ref<any>(null)
const ticket = ref<any>(null)

const panel = ref("authenticate")

const panels: { [id: string]: Component } = {
  authenticate: AccountAuthenticate,
  mfa: FactorPicker,
  applicator: FactorApplicator,
  completed: AuthenticateCompleted,
}
</script>

<style scoped>
.card-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

@media (max-width: 768px) {
  .card-grid {
    grid-template-columns: 1fr;
  }
}

.card-rounded {
  border-radius: 8px;
}
</style>
