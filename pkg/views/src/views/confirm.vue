<template>
  <v-container class="h-screen flex flex-col gap-3 items-center justify-center">
    <v-card class="w-full max-w-[720px]" :loading="loading">
      <v-card-text class="card-grid pa-9">
        <div>
          <v-avatar color="accent" icon="mdi-check-decagram" size="large" class="card-rounded mb-2" />
          <h1 class="text-2xl">Confirm registration</h1>
          <p>Confirm your account to keep your account longer than 48 hours.</p>
        </div>

        <v-window :touch="false" :model-value="panel" class="pa-2 mx-[-0.5rem]">
          <v-window-item value="confirm">
            <div>
              <v-expand-transition>
                <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
                  Something went wrong... {{ error }}
                </v-alert>
              </v-expand-transition>

              <v-progress-circular v-if="!error" indeterminate size="32" color="grey-darken-3" class="mb-3" />

              <h1 class="font-bold text-xl">Confirming</h1>
              <p>We are confirming your account. Please stand by, this won't took a long time...</p>
            </div>
          </v-window-item>
          <v-window-item value="callback">
            <div>
              <v-icon icon="mdi-fire" size="32" color="grey-darken-3" class="mb-3" />

              <h1 class="font-bold text-xl">Confirmed</h1>
              <p>You're done! We sucessfully confirmed your account.</p>

              <p class="mt-3">Now you can continue use Solarpass, we will redirect to dashboard you soon.</p>
            </div>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <copyright />
  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import { request } from "@/scripts/request"
import { useUserinfo } from "@/stores/userinfo"
import Copyright from "@/components/Copyright.vue"

const route = useRoute()
const router = useRouter()
const { readProfiles } = useUserinfo()

const error = ref<string | null>(null)

const loading = ref(false)

const panel = ref("confirm")

async function confirm() {
  if (!route.query["tk"]) {
    error.value = "code was not exists"
    return
  }

  const res = await request("/api/users/me/confirm", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      code: route.query["tk"],
    }),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    loading.value = true
    panel.value = "callback"
    await readProfiles()
    router.push({ name: "dashboard" })
  }
  loading.value = false
}

confirm()
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
