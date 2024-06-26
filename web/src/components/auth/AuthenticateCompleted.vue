<template>
  <div>
    <v-icon icon="mdi-lan-check" size="32" color="grey-darken-3" class="mb-3" />

    <h1 class="font-bold text-xl">All Done!</h1>
    <p>Welcome back! You just signed in right now! We're going to send you to jesus...</p>

    <v-expand-transition>
      <v-alert v-show="error" variant="tonal" type="error" class="text-xs mb-3">
        Something went wrong... {{ error }}
      </v-alert>
    </v-expand-transition>
  </div>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { useUserinfo } from "@/stores/userinfo"
import { onMounted, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()
const userinfo = useUserinfo()

const props = defineProps<{ loading?: boolean; currentFactor?: any; ticket?: any }>()
const emits = defineEmits(["update:loading"])

const error = ref<string | null>(null)

async function load() {
  emits("update:loading", true)
  await getToken(props.ticket.grant_token)
  await userinfo.readProfiles()
  emits("update:loading", false)
  setTimeout(() => callback(), 3000)
}

onMounted(() => load())

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
  if (route.query["close"]) {
    window.close()
  } else if (route.query["redirect_uri"]) {
    window.open((route.query["redirect_uri"] as string) ?? "/", "_self")
  } else {
    router.push({ name: "dashboard" })
  }
}
</script>