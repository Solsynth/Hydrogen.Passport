<template>
  <v-dialog class="max-w-[720px]" :model-value="data != null" @update:model-value="(val) => !val && emits('close')">
    <template v-slot:default="{ isActive }">
      <v-card title="Assign permissions" :subtitle="`To user @${props.data?.name}`" :loading="submitting">
        <v-card-text>
          <v-sheet elevation="2" rounded="lg">
            <v-table density="comfortable">
              <thead>
              <tr>
                <th class="text-left">
                  Key
                </th>
                <th class="text-left">
                  Value
                </th>
              </tr>
              </thead>
              <tbody>
              <tr
                v-for="[key, val] in Object.entries(perms)"
                :key="key"
              >
                <td class="w-1/2">
                  <div>
                    <p>{{ key }}</p>
                    <div class="flex mx-[-8px]">
                      <v-btn color="error" text="Delete" variant="plain" size="x-small"
                             @click="() => deleteNode(key)" />
                      <v-btn class="ms-[-8px]" color="info" text="Change" variant="plain" size="x-small"
                             @click="() => changeNodeType(key)" />
                    </div>
                  </div>
                </td>
                <td class="w-1/2">
                  <div class="w-full flex items-center">
                    <v-checkbox v-if="typeof val === 'boolean'" class="my-1" density="comfortable"
                                :hide-details="true"
                                v-model="perms[key]" />
                    <v-number-input v-else-if="typeof val === 'number'"
                                    controlVariant="default"
                                    :reverse="false"
                                    :hideInput="false"
                                    :inset="false"
                                    class="font-mono my-2"
                                    density="compact" :hide-details="true"
                                    v-model="perms[key]" />
                    <v-text-field v-else class="font-mono my-2" density="compact" :hide-details="true"
                                  v-model="perms[key]" />
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <v-text-field class="my-3.5" label="Key" density="compact" variant="solo-filled"
                                v-model="pendingNodeKey"
                                :hide-details="true" />
                </td>
                <td>
                  <div class="w-full flex justify-center">
                    <v-btn prepend-icon="mdi-plus-circle" text="Add one" block rounded="md" @click="addNode" />
                  </div>
                </td>
              </tr>
              </tbody>
            </v-table>
          </v-sheet>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            :disabled="submitting"
            text="Cancel"
            color="grey"
            @click="isActive.value = false"
          ></v-btn>
          <v-btn
            :disabled="submitting"
            text="Apply Changes"
            @click="saveNode"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"

const perms = ref<any>({})

const pendingNodeKey = ref("")

const props = defineProps<{ data: any }>()
const emits = defineEmits(["close", "success", "error"])

watch(props, (v) => {
  if (v.data != null) {
    perms.value = v.data["perm_nodes"]
  }
}, { immediate: true, deep: true })

function addNode() {
  if (pendingNodeKey.value) {
    perms.value[pendingNodeKey.value] = false
    pendingNodeKey.value = ""
  }
}

function deleteNode(key: string) {
  delete perms.value[key]
}

function changeNodeType(key: string) {
  const typelist = [
    "boolean",
    "number",
    "string",
  ]
  const idx = typelist.indexOf(typeof perms.value[key])
  if (idx == -1 || idx == typelist.length - 1) {
    perms.value[key] = false
    return
  }
  switch (typelist[idx + 1]) {
    case "boolean":
      perms.value[key] = false
      break
    case "number":
      perms.value[key] = 0
      break
    default:
      perms.value[key] = ""
      break
  }
}

const submitting = ref(false)

async function saveNode() {
  submitting.value = true
  const res = await request(`/api/admin/users/${props.data.id}/permissions`, {
    method: 'PUT',
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${getAtk()}`,
    },
    body: JSON.stringify({
      'perm_nodes': perms.value,
    }),
  })
  if (res.status !== 200) {
    emits("error", await res.text())
  } else {
    emits("success")
    emits("close")
  }
  submitting.value = false
}
</script>
