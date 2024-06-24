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
import { computed, ref } from "vue"

const password = ref("")

const error = ref<string | null>(null)

const props = defineProps<{ loading?: boolean; currentFactor?: any; ticket?: any }>()
const emits = defineEmits(["swap", "update:ticket", "update:loading"])

async function submit() {
  emits("update:loading", true)
  const res = await request(`/api/auth/mfa`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      ticket_id: props.ticket?.id,
      factor_id: props.currentFactor?.id,
      code: password.value,
    }),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    error.value = null
    password.value = ""
    emits("update:ticket", data["ticket"])
    if (data["is_finished"]) emits("swap", "completed")
    else emits("swap", "mfa")
  }
  emits("update:loading", false)
}

const inputType = computed(() => {
  switch (props.currentFactor?.type) {
    case 0:
      return "text"
    case 1:
      return "one-time-password"
    default:
      return "unknown"
  }
})
</script>
