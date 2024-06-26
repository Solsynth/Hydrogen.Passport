<template>
  <v-navigation-drawer :model-value="props.open" @update:model-value="(val: any) => emits('update:open', val)" location="right"
                       temporary order="0" width="400">
    <v-list-item prepend-icon="mdi-bell" title="Notifications" class="py-3"></v-list-item>

    <v-divider color="black" class="mb-1" />

    <v-list v-if="notify.notifications.length <= 0" density="compact">
      <v-list-item color="secondary" prepend-icon="mdi-check" title="All notifications read"
                   subtitle="There is no more new things for you..." />
    </v-list>

    <v-list v-else density="compact" lines="three">
      <v-list-item v-for="(item, idx) in notify.notifications" :key="idx">
        <template #title>{{ item.subject }}</template>
        <template #subtitle>{{ item.content }}</template>

        <template #append>
          <v-btn icon="mdi-check" size="x-small" variant="text" :disabled="loading" @click="markAsRead(item, idx)" />
        </template>

        <div class="flex text-xs gap-1">
          <a v-for="(link, idx) in item.links" :key="idx" class="mt-1 underline" target="_blank"
             :href="link.url">{{ link.label }}</a>
        </div>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import { computed, onMounted, onUnmounted, ref } from "vue"
import { useNotifications } from "@/stores/notifications"

const props = defineProps<{ open: boolean }>()
const emits = defineEmits(["update:open"])

const notify = useNotifications()

const error = ref<string | null>(null)
const submitting = ref(false)
const loading = computed(() => notify.loading || submitting.value)

async function markAsRead(item: any, idx: number) {
  submitting.value = true
  const res = await request(`/api/notifications/${item.id}/read`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${getAtk()}` },
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    notify.remove(idx)
    error.value = null
  }
  submitting.value = false
}

notify.list()

onMounted(() => notify.connect())
onUnmounted(() => notify.disconnect())
</script>
