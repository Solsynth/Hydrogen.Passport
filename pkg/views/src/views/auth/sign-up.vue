<template>
  <v-container class="h-screen flex flex-col gap-3 items-center justify-center">
    <callback-notify />

    <v-card class="w-full max-w-[720px] overflow-auto" :loading="loading">
      <v-card-text class="card-grid pa-9">
        <div>
          <v-avatar color="accent" icon="mdi-login-variant" size="large" class="card-rounded mb-2" />
          <h1 class="text-2xl">Create an account</h1>
          <p>Create an account on Solar Network. Then enjoy all our services.</p>
        </div>

        <div class="flex items-center">
          <v-form class="flex-grow-1" @submit.prevent="submit">
            <v-row dense class="mb-3">
              <v-col :cols="6">
                <v-text-field
                  hide-details
                  label="Name"
                  autocomplete="username"
                  variant="solo"
                  density="comfortable"
                  v-model="data.name"
                />
              </v-col>
              <v-col :cols="6">
                <v-text-field
                  hide-details
                  label="Nick"
                  autocomplete="nickname"
                  variant="solo"
                  density="comfortable"
                  v-model="data.nick"
                />
              </v-col>
              <v-col :cols="12">
                <v-text-field
                  hide-details
                  label="Email Address"
                  type="email"
                  variant="solo"
                  density="comfortable"
                  v-model="data.email"
                />
              </v-col>
              <v-col :cols="12">
                <v-text-field
                  hide-details
                  label="Password"
                  type="password"
                  autocomplete="new-password"
                  variant="solo"
                  density="comfortable"
                  v-model="data.password"
                />
              </v-col>
            </v-row>

            <v-expand-transition>
              <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
                Something went wrong... {{ error }}
              </v-alert>
            </v-expand-transition>

            <div class="flex justify-between">
              <v-btn type="button" variant="plain" color="grey-darken-3" :to="{ name: 'auth.sign-in' }">
                Sign in
              </v-btn>

              <v-btn type="submit" variant="text" color="primary" append-icon="mdi-arrow-right" :disabled="loading">
                Next
              </v-btn>
            </div>
          </v-form>
        </div>
      </v-card-text>
    </v-card>

    <v-dialog v-model="done" class="max-w-[560px]">
      <v-card title="Congratulations">
        <template #text>
          You successfully created an account on Solar Network. Now sign in to your account and start exploring!
        </template>
        <template #actions>
          <div class="flex flex-grow-1 justify-end">
            <v-btn @click="callback">Let's go</v-btn>
          </div>
        </template>
      </v-card>
    </v-dialog>

    <copyright />
  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { request } from "@/scripts/request"
import { useRoute, useRouter } from "vue-router"
import Copyright from "@/components/Copyright.vue"
import CallbackNotify from "@/components/auth/CallbackNotify.vue"

const error = ref<string | null>(null)

const route = useRoute()
const router = useRouter()

const done = ref(false)
const loading = ref(false)

const data = ref({
  name: "",
  nick: "",
  email: "",
  password: "",
})

async function submit() {
  const payload = data.value
  if (!payload.name || !payload.nick || !payload.email || !payload.password) return

  loading.value = true
  const res = await request("/api/users", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    done.value = true
    error.value = null
  }
  loading.value = false
}

function callback() {
  if (route.params["closable"]) {
    window.close()
  } else {
    router.push({ name: "auth.sign-in" })
  }
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
