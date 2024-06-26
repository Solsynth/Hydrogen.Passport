<template>
  <div>
    <GoUseSolian class="mb-3" />

    <v-card title="Information" prepend-icon="mdi-face-man-profile" :loading="loading">
      <template #text>
        <v-form class="mt-1" @submit.prevent="submit">
          <v-row dense>
            <v-col :xs="12" :md="6">
              <v-text-field readonly hide-details label="Username" density="comfortable" variant="outlined"
                v-model="data.name" />
            </v-col>
            <v-col :xs="12" :md="6">
              <v-text-field hide-details label="Nickname" density="comfortable" variant="outlined"
                v-model="data.nick" />
            </v-col>
            <v-col :cols="12">
              <v-textarea hide-details label="Description" density="comfortable" variant="outlined"
                v-model="data.description" />
            </v-col>
            <v-col :xs="12" :md="6" :lg="4">
              <v-text-field hide-details label="First Name" density="comfortable" variant="outlined"
                v-model="data.first_name" />
            </v-col>
            <v-col :xs="12" :md="6" :lg="4">
              <v-text-field hide-details label="Last Name" density="comfortable" variant="outlined"
                v-model="data.last_name" />
            </v-col>
            <v-col :xs="12" :lg="4">
              <v-text-field hide-details label="Birthday" density="comfortable" variant="outlined" type="datetime-local"
                v-model="data.birthday" />
            </v-col>
          </v-row>

          <v-btn type="submit" class="mt-2" variant="text" prepend-icon="mdi-content-save" :disabled="loading">
            Apply Changes
          </v-btn>
        </v-form>
      </template>
    </v-card>

    <v-snackbar v-model="done" :timeout="3000"> Your personal information has been updated. </v-snackbar>

    <!-- @vue-ignore -->
    <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { useUserinfo, getAtk } from "@/stores/userinfo"
import { request } from "@/scripts/request"
import GoUseSolian from "@/components/GoUseSolian.vue"

const id = useUserinfo()

const error = ref<string | null>(null)
const done = ref(false)
const loading = ref(false)

const data = ref<any>({})
const avatar = ref<any>(null)
const banner = ref<any>(null)

watch(
  id,
  (val) => {
    if (val.isReady) {
      data.value.name = id.userinfo.data.name
      data.value.nick = id.userinfo.data.nick
      data.value.description = id.userinfo.data.description
      data.value.first_name = id.userinfo.data.profile.first_name
      data.value.last_name = id.userinfo.data.profile.last_name
      data.value.birthday = id.userinfo.data.profile.birthday

      if (data.value.birthday) data.value.birthday = data.value.birthday.substring(0, 16)
    }
  },
  { immediate: true, deep: true },
)

async function submit() {
  const payload = data.value
  if (payload.birthday) payload.birthday = new Date(payload.birthday).toISOString()

  loading.value = true
  const res = await request("/api/users/me", {
    method: "PUT",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload),
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await id.readProfiles()
    done.value = true
    error.value = null
  }
  loading.value = false
}

async function applyAvatar() {
  if (!avatar.value) return

  if (loading.value) return

  const payload = new FormData()
  payload.set("avatar", avatar.value[0])

  loading.value = true
  const res = await request("/api/users/me/avatar", {
    method: "PUT",
    headers: { Authorization: `Bearer ${getAtk()}` },
    body: payload,
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await id.readProfiles()
    done.value = true
    error.value = null
    avatar.value = null
  }
  loading.value = false
}

async function applyBanner() {
  if (!banner.value) return

  if (loading.value) return

  const payload = new FormData()
  payload.set("banner", banner.value[0])

  loading.value = true
  const res = await request("/api/users/me/banner", {
    method: "PUT",
    headers: { Authorization: `Bearer ${getAtk()}` },
    body: payload,
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await id.readProfiles()
    done.value = true
    error.value = null
    banner.value = null
  }
  loading.value = false
}
</script>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
