<template>
  <div class="flex items-center">
    <v-form class="flex-grow-1" @submit.prevent="submit">
      <v-text-field label="Username" variant="solo" density="comfortable" class="mb-3" :hide-details="true"
                    :disabled="props.loading" v-model="probe" />
      <v-text-field label="Password" variant="solo" density="comfortable" type="password" :disabled="props.loading"
                    v-model="password" />

      <v-expand-transition>
        <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
          Something went wrong... {{ error }}
        </v-alert>
      </v-expand-transition>

      <div class="flex justify-between">
        <v-btn type="button" variant="plain" color="grey-darken-3" :to="{ name: 'auth.sign-up' }">Sign up</v-btn>

        <v-btn
          type="submit"
          variant="text"
          color="primary"
          class="justify-self-end"
          append-icon="mdi-arrow-right"
          :disabled="props.loading"
        >
          Next
        </v-btn>
      </div>
    </v-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { request } from "@/scripts/request"

const probe = ref("")
const password = ref("")

const error = ref<string | null>(null)

const props = defineProps<{ loading?: boolean }>()
const emits = defineEmits(["swap", "update:loading", "update:ticket"])

async function submit() {
  if (!probe.value || !password.value) return

  emits("update:loading", true)
  const res = await request("/api/auth", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: probe.value, password: password.value }),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    emits("update:ticket", data["ticket"])
    if (data.is_finished) emits("swap", "completed")
    else emits("swap", "mfa")
    error.value = null
  }
  emits("update:loading", false)
}
</script>
