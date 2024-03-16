<template>
  <v-container class="h-screen flex flex-col gap-3 items-center justify-center">
    <v-card class="w-full max-w-[720px]" :loading="loading">
      <v-card-text class="card-grid pa-9">
        <div>
          <v-avatar color="accent" icon="mdi-connection" size="large" class="card-rounded mb-2" />
          <h1 class="text-2xl">Connect to third-party</h1>
          <p>One Solarpass, entire internet.</p>
        </div>

        <v-window :model-value="panel" class="pa-2 mx-[-0.5rem]">
          <v-window-item value="confirm">
            <div class="flex flex-col gap-2">
              <v-expand-transition>
                <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
                  <p>Something went wrong... {{ error }}</p>
                  <br />

                  <p class="font-bold">
                    It's usually not our fault. Try bringing this link to give feedback to the developer of the app you
                    came from.
                  </p>
                </v-alert>
              </v-expand-transition>

              <div v-if="!error">
                <h1 class="font-bold text-xl">{{ metadata?.name ?? "Loading" }}</h1>
                <p>{{ metadata?.description ?? "Hold on a second please!" }}</p>

                <div class="mt-3">
                  <p class="opacity-80 text-xs">Permissions they requested</p>
                  <v-card variant="tonal" class="mt-1 mx-[-4px]">
                    <v-list density="compact">
                      <v-list-item v-for="claim in requestedClaims" lines="two">
                        <template #title>
                          <span class="capitalize">{{ getClaimDescription(claim)?.name }}</span>
                        </template>
                        <template #subtitle>
                          <span>{{ getClaimDescription(claim)?.description }}</span>
                        </template>
                        <template #prepend>
                          <v-icon :icon="getClaimDescription(claim)?.icon" size="x-large" />
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card>

                  <div class="mt-5 flex justify-between">
                    <v-btn prepend-icon="mdi-close" variant="text" color="error" :disabled="loading" @click="decline">
                      Decline
                    </v-btn>
                    <v-btn append-icon="mdi-check" variant="tonal" color="success" :disabled="loading" @click="approve">
                      Approve
                    </v-btn>
                  </div>

                  <div class="mt-5 text-xs text-center opacity-75">
                    <p>After approve their request, you will be redirect to</p>
                    <p class="text-mono">{{ route.query["redirect_uri"] }}</p>
                  </div>
                </div>
              </div>
            </div>
          </v-window-item>

          <v-window-item value="callback">
            <div>
              <v-icon icon="mdi-fire" size="32" color="grey-darken-3" class="mb-3" />

              <h1 class="font-bold text-xl">Authoirzed</h1>
              <p>You're done! We sucessfully established connection between you and {{ metadata?.name }}.</p>

              <p class="mt-3">Now you can continue your their app, we will redirect you soon.</p>

              <p class="mt-3">Teleporting you to...</p>
              <p class="text-xs text-mono">{{ route.query["redirect_uri"] }}</p>
            </div>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <copyright />
  </v-container>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { useRoute } from "vue-router"
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import { claims, type ClaimType } from "@/views/auth/claims"
import Copyright from "@/components/Copyright.vue"

const route = useRoute()

const error = ref<string | null>(null)

const loading = ref(false)

const metadata = ref<any>(null)
const requestedClaims = computed(() => {
  const scope: string = (route.query["scope"] as string) ?? ""
  return scope.split(" ")
})

const panel = ref("confirm")

async function preconnect() {
  const res = await request(`/api/auth/o/connect${location.search}`, {
    headers: { Authorization: `Bearer ${getAtk()}` },
  })

  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()

    if (data["session"]) {
      panel.value = "callback"
      callback(data["session"])
    } else {
      document.title = `Solarpass | Connect to ${data["client"]?.name}`
      metadata.value = data["client"]
      loading.value = false
    }
  }
}

preconnect()

function decline() {
  if (window.history.length > 0) {
    window.history.back()
  } else {
    window.close()
  }
}

async function approve() {
  loading.value = true
  const res = await request(
    "/api/auth/o/connect?" +
      new URLSearchParams({
        client_id: route.query["client_id"] as string,
        redirect_uri: encodeURIComponent(route.query["redirect_uri"] as string),
        response_type: "code",
        scope: route.query["scope"] as string,
      }),
    {
      method: "POST",
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )

  if (res.status !== 200) {
    error.value = await res.text()
    loading.value = false
  } else {
    const data = await res.json()
    panel.value = "callback"
    setTimeout(() => callback(data["session"]), 1850)
  }
}

function callback(session: any) {
  const url = `${route.query["redirect_uri"]}?code=${session["grant_token"]}&state=${route.query["state"]}`
  window.open(url, "_self")
}

function getClaimDescription(key: string): ClaimType {
  return claims.hasOwnProperty(key) ? claims[key] : { icon: "mdi-asterisk", name: key, description: "Unknown claim..." }
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
