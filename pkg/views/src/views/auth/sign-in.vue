<template>
  <v-container class="h-screen flex flex-col gap-3 items-center justify-center">
    <v-card class="w-full max-w-[720px]" :loading="loading">
      <v-card-text class="card-grid pa-9">
        <div>
          <v-avatar color="accent" icon="mdi-login-variant" size="large" class="card-rounded mb-2" />
          <h1 class="text-2xl">Sign in</h1>
          <p>Through sign in to access the entire Solar Network.</p>
        </div>

        <v-window :model-value="panel" class="pa-2 mx-[-0.5rem]">
          <v-window-item v-for="k in Object.keys(panels)" :value="k">
            <component
              :is="panels[k]"
              @swap="(val: string) => (panel = val)"
              v-model:loading="loading"
              v-model:factors="factors"
              v-model:currentFactor="currentFactor"
              v-model:challenge="challenge"
            />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <copyright />
  </v-container>
</template>

<script setup lang="ts">
import { ref, type Component } from "vue"
import Copyright from "@/components/Copyright.vue"
import AccountLocator from "@/components/auth/AccountLocator.vue"
import FactorPicker from "@/components/auth/FactorPicker.vue"
import FactorApplicator from "@/components/auth/FactorApplicator.vue"

const loading = ref(false)

const factors = ref<any>(null)
const currentFactor = ref<any>(null)
const challenge = ref<any>(null)

const panel = ref("locate")

const panels: { [id: string]: Component } = {
  locate: AccountLocator,
  pick: FactorPicker,
  applicator: FactorApplicator,
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
@/components/Copyright.vue
