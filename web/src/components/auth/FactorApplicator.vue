<template>
  <div class="flex items-center">
    <v-form class="flex-grow-1" @submit.prevent="submit">
      <div v-if="inputType === 'one-time-password'" class="text-center">
        <p class="text-xs opacity-90">Check your inbox!</p>
        <v-otp-input
          class="pt-0"
          variant="solo"
          density="compact"
          type="text"
          :length="6"
          v-model="password"
          :loading="loading"
        />
      </div>
      <v-text-field
        v-else
        label="Password"
        type="password"
        variant="solo"
        density="comfortable"
        :disabled="loading"
        v-model="password"
      />

      <v-expand-transition>
        <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
          Something went wrong... {{ error }}
        </v-alert>
      </v-expand-transition>

      <div class="flex justify-end">
        <v-btn
          type="submit"
          variant="text"
          color="primary"
          class="justify-self-end"
          append-icon="mdi-arrow-right"
          :disabled="loading"
        >
          Next
        </v-btn>
      </div>
    </v-form>
  </div>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useUserinfo } from "@/stores/userinfo"
import { computed, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const password = ref("")

const error = ref<string | null>(null)

const props = defineProps<{ loading?: boolean; currentFactor?: any; challenge?: any }>()
const emits = defineEmits(["swap", "update:challenge"])

const route = useRoute()
const router = useRouter()

const { readProfiles } = useUserinfo()

async function submit() {
  const res = await request(`/api/auth`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      challenge_id: props.challenge?.id,
      factor_id: props.currentFactor?.id,
      secret: password.value,
    }),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    if (data["is_finished"]) {
      await getToken(data["session"]["grant_token"])
      await readProfiles()
      callback()
    } else {
      emits("swap", "pick")
      emits("update:challenge", data["challenge"])
      error.value = null
      password.value = ""
    }
  }
}

async function getToken(tk: string) {
  const res = await request("/api/auth/token", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      code: tk,
      grant_type: "grant_token",
    }),
  })
  if (res.status !== 200) {
    const err = await res.text()
    error.value = err
    throw new Error(err)
  } else {
    error.value = null
  }
}

function callback() {
  if (route.query["closable"]) {
    window.close()
  } else if (route.query["redirect_uri"]) {
    window.open((route.query["redirect_uri"] as string) ?? "/", "_self")
  } else {
    router.push({ name: "dashboard" })
  }
}

const inputType = computed(() => {
  switch (props.currentFactor?.type) {
    case 0:
      return "text"
    case 1:
      return "one-time-password"
  }
})
</script>
