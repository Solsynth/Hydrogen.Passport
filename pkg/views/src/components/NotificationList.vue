<template>
  <v-menu eager :close-on-content-click="false">
    <template #activator="{ props }">
      <v-btn v-bind="props" stacked rounded="circle" size="small" variant="text" :loading="loading">
        <v-badge v-if="pagination.total > 0" color="error" :content="pagination.total">
          <v-icon icon="mdi-bell" />
        </v-badge>

        <v-icon v-else icon="mdi-bell" />
      </v-btn>
    </template>

    <v-list v-if="notifications.length <= 0" class="w-[380px]" density="compact">
      <v-list-item>
        <v-alert class="text-sm" variant="tonal" type="info">You are done! There is no unread notifications for you.</v-alert>
      </v-list-item>
    </v-list>

    <v-list v-else class="w-[380px]" density="compact" lines="three">
      <v-list-item v-for="item in notifications">
        <template #title>{{ item.subject }}</template>
        <template #subtitle>{{ item.content }}</template>

        <template #append>
          <v-btn icon="mdi-check" size="x-small" variant="text" :disabled="loading" @click="markAsRead(item)" />
        </template>

        <div class="flex text-xs gap-1">
          <a v-for="link in item.links" class="mt-1 underline" target="_blank" :href="link.url">{{ link.label }}</a>
        </div>
      </v-list-item>
    </v-list>
  </v-menu>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import { reactive, ref } from "vue"

const loading = ref(false)

const error = ref<string | null>(null)

const notifications = ref<any[]>([])
const pagination = reactive({ page: 1, pageSize: 25, total: 0 })

async function readNotifications() {
  loading.value = true
  const res = await request(
    "/api/notifications?" +
      new URLSearchParams({
        take: pagination.pageSize.toString(),
        offset: ((pagination.page - 1) * pagination.pageSize).toString(),
      }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status === 200) {
    const data = await res.json()
    notifications.value = data["data"]
    pagination.total = data["count"]
  }
  loading.value = false
}

readNotifications()

async function markAsRead(item: any) {
  loading.value = true
  const res = await request(`/api/notifications/${item.id}/read`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${getAtk()}` },
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await readNotifications()
    error.value = null
  }
  loading.value = false
}
</script>
