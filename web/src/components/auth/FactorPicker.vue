<template>
  <div class="flex items-center">
    <div class="flex-grow-1">
      <v-card class="mb-3">
        <v-list density="compact" color="primary">
          <v-list-item
            v-for="item in props.factors ?? []"
            :prepend-icon="getFactorType(item)?.icon"
            :title="getFactorType(item)?.label"
            :active="focus === item.id"
            :disabled="getFactorAvailable(item)"
            @click="focus = item.id"
          />
        </v-list>
      </v-card>

      <v-expand-transition>
        <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
          Something went wrong... {{ error }}
        </v-alert>
      </v-expand-transition>

      <div class="flex justify-end">
        <v-btn variant="text" color="primary" class="justify-self-end" append-icon="mdi-arrow-right" @click="submit">
          Next
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { request } from "@/scripts/request"

const focus = ref<number | null>(null)

const error = ref<string | null>(null)

const props = defineProps<{ factors?: any[]; challenge?: any }>()
const emits = defineEmits(["swap", "update:loading", "update:currentFactor"])

async function submit() {
  if (!focus) return

  emits("update:loading", true)
  const res = await request(`/api/auth/factors/${focus.value}`, {
    method: "POST",
  })
  if (res.status !== 200 && res.status !== 204) {
    error.value = await res.text()
  } else {
    const item = props.factors?.find((item: any) => item.id === focus.value)
    emits("update:currentFactor", item)
    emits("swap", "applicator")
    error.value = null
    focus.value = null
  }
  emits("update:loading", false)
}

function getFactorType(item: any) {
  switch (item.type) {
    case 0:
      return { icon: "mdi-form-textbox-password", label: "Password Validation" }
    case 1:
      return { icon: "mdi-email-fast", label: "Email One Time Password" }
  }
}

function getFactorAvailable(factor: any) {
  const blacklist: number[] = props.challenge?.blacklist_factors ?? []
  return blacklist.includes(factor.id)
}
</script>
