<template>
  <v-dialog class="max-w-[720px]" :model-value="props.data != null"
            @update:model-value="(val) => !val && emits('close')"
            :loading="reverting">
    <template v-slot:default="{ isActive }">
      <v-card title="Auth Factors" :subtitle="`Of user @${props.data?.name}`">
        <v-card-text>
          <v-sheet elevation="2" rounded="lg">
            <v-table density="compact">
              <thead>
              <tr>
                <th class="text-left">
                  Name
                </th>
                <th class="text-left">
                  Secret
                </th>
              </tr>
              </thead>
              <tbody>
              <tr
                v-for="item in factors"
                :key="item.name"
              >
                <td class="w-1/2">{{ item.id }}</td>
                <td class="w-1/2"><code>{{ item.secret }}</code></td>
              </tr>
              </tbody>
            </v-table>
          </v-sheet>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            text="Close"
            @click="isActive.value = false"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { request } from "@/scripts/request"

const props = defineProps<{ data: any }>()
const emits = defineEmits(["close", "error"])

const reverting = ref(false)

const factors = ref<any[]>([])

async function load() {
  reverting.value = true
  const res = await request(`/api/admin/users/${props.data.id}/factors`)
  if (res.status !== 200) {
    emits("error", await res.text())
  } else {
    factors.value = await res.json()
  }
  reverting.value = false
}

watch(props, (v) => {
  if (v.data != null) {
    factors.value = []
    load()
  }
}, { immediate: true, deep: true })
</script>