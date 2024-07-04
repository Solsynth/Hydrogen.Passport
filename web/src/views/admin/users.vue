<template>
  <div>
    <v-data-table-server
      fixed-header
      class="h-full"
      density="compact"
      :headers="dataDefinitions.users"
      :items="users"
      :items-length="pagination.total"
      :loading="reverting"
      v-model:items-per-page="pagination.pageSize"
      @update:options="readUsers"
      item-value="id"
    >
      <template v-slot:top>
        <v-toolbar color="secondary">
          <div class="max-md:px-5 md:px-12 flex flex-grow-1 items-center">
            <v-btn class="me-2" icon="mdi-account-group" density="compact" :to="{ name: 'admin.dashboard' }" exact />
            <h3 class="ml-2 text-lg font-500">Users</h3>
          </div>
        </v-toolbar>
      </template>

      <template v-slot:item="{ item }: { item: any }">
        <tr>
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.nick }}</td>
          <td>{{ new Date(item.created_at).toLocaleString() }}</td>
          <td>
            <v-tooltip text="Details">
              <template #activator="{ props }">
                <v-btn
                  v-bind="props"
                  variant="text"
                  size="x-small"
                  color="info"
                  icon="mdi-dots-horizontal"
                  @click="viewingUser = item"
                />
              </template>
            </v-tooltip>
            <v-tooltip text="Assign Permissions">
              <template #activator="{ props }">
                <v-btn
                  v-bind="props"
                  variant="text"
                  size="x-small"
                  color="teal"
                  icon="mdi-code-block-braces"
                  @click="assigningPermUser = item"
                />
              </template>
            </v-tooltip>
          </td>
        </tr>
      </template>
    </v-data-table-server>

    <user-detail-panel :data="viewingUser" @close="viewingUser = null" />
    <user-assign-perms-panel :data="assigningPermUser" @close="assigningPermUser = null"
                             @success="readUsers(pagination)"
                             @error="val => error = val" />

    <v-snackbar :timeout="3000" :model-value="error != null" @update:model-value="_ => error = null">
      {{ error }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue"
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import UserDetailPanel from "@/components/admin/UserDetailPanel.vue"
import UserAssignPermsPanel from "@/components/admin/UserAssignPermsPanel.vue"

const error = ref<string | null>(null)

const users = ref<any[]>([])

const viewingUser = ref<any>(null)
const assigningPermUser = ref<any>(null)

const dataDefinitions: { [id: string]: any[] } = {
  users: [
    { align: "start", key: "id", title: "ID" },
    { align: "start", key: "name", title: "Name" },
    { align: "start", key: "nick", title: "Nick" },
    { align: "start", key: "created_at", title: "Created At" },
    { align: "start", key: "actions", title: "Actions", sortable: false },
  ],
}

const reverting = ref(true)
const pagination = reactive({
  page: 1, pageSize: 5, total: 0,
})

async function readUsers({ page, itemsPerPage }: { page?: number; itemsPerPage?: number }) {
  reverting.value = true
  const res = await request(
    "/api/admin/users?" +
    new URLSearchParams({
      take: pagination.pageSize.toString(),
      offset: ((pagination.page - 1) * pagination.pageSize).toString(),
    }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    users.value = data["data"]
    pagination.total = data["count"]
  }
  reverting.value = false
}

onMounted(() => readUsers({}))
</script>
