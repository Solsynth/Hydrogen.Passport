<template>
  <div>
    <v-expansion-panels>
      <v-expansion-panel eager title="Tickets">
        <template #text>
          <v-card :loading="reverting.tickets" variant="outlined">
            <v-data-table-server
              density="compact"
              :headers="dataDefinitions.tickets"
              :items="tickets"
              :items-length="pagination.tickets.total"
              :loading="reverting.tickets"
              v-model:items-per-page="pagination.tickets.pageSize"
              @update:options="readTickets"
              item-value="id"
            >
              <template v-slot:item="{ item }: { item: any }">
                <tr>
                  <td>{{ item.id }}</td>
                  <td>{{ item.ip_address }}</td>
                  <td>
                    <v-tooltip :text="item.user_agent" location="top">
                      <template #activator="{ props }">
                        <div v-bind="props" class="text-ellipsis whitespace-nowrap overflow-hidden max-w-[280px]">
                          {{ item.user_agent }}
                        </div>
                      </template>
                    </v-tooltip>
                  </td>
                  <td>{{ new Date(item.created_at).toLocaleString() }}</td>
                  <td>
                    <v-tooltip text="Sign out">
                      <template #activator="{ props }">
                        <v-btn
                          v-bind="props"
                          variant="text"
                          size="x-small"
                          color="error"
                          icon="mdi-logout-variant"
                          @click="killTicket(item)"
                        />
                      </template>
                    </v-tooltip>
                  </td>
                </tr>
              </template>
            </v-data-table-server>
          </v-card>
        </template>
      </v-expansion-panel>

      <v-expansion-panel eager title="Events">
        <template #text>
          <v-card :loading="reverting.events" variant="outlined">
            <v-data-table-server
              density="compact"
              :headers="dataDefinitions.events"
              :items="events"
              :items-length="pagination.events.total"
              :loading="reverting.events"
              v-model:items-per-page="pagination.events.pageSize"
              @update:options="readEvents"
              item-value="id"
            >
              <template v-slot:item="{ item }: { item: any }">
                <tr>
                  <td>{{ item.id }}</td>
                  <td>{{ item.type }}</td>
                  <td>{{ item.target }}</td>
                  <td>{{ item.ip_address }}</td>
                  <td>
                    <v-tooltip :text="item.user_agent" location="top">
                      <template #activator="{ props }">
                        <div v-bind="props" class="text-ellipsis whitespace-nowrap overflow-hidden max-w-[180px]">
                          {{ item.user_agent }}
                        </div>
                      </template>
                    </v-tooltip>
                  </td>
                  <td>{{ new Date(item.created_at).toLocaleString() }}</td>
                </tr>
              </template>
            </v-data-table-server>
          </v-card>
        </template>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script setup lang="ts">
import { request } from "@/scripts/request"
import { getAtk } from "@/stores/userinfo"
import { reactive, ref } from "vue"

const error = ref<string | null>(null)

const dataDefinitions: { [id: string]: any[] } = {
  tickets: [
    { align: "start", key: "id", title: "ID" },
    { align: "start", key: "ip_address", title: "IP Address" },
    { align: "start", key: "user_agent", title: "User Agent" },
    { align: "start", key: "created_at", title: "Issued At" },
    { align: "start", key: "actions", title: "Actions", sortable: false },
  ],
  events: [
    { align: "start", key: "id", title: "ID" },
    { align: "start", key: "type", title: "Type" },
    { align: "start", key: "target", title: "Affected Object" },
    { align: "start", key: "ip_address", title: "IP Address" },
    { align: "start", key: "user_agent", title: "User Agent" },
    { align: "start", key: "created_at", title: "Performed At" },
  ],
}

const tickets = ref<any>([])
const events = ref<any>([])

const reverting = reactive({ tickets: false, sessions: false, events: false })
const pagination = reactive({
  tickets: { page: 1, pageSize: 5, total: 0 },
  events: { page: 1, pageSize: 5, total: 0 },
})

async function readTickets({ page, itemsPerPage }: { page?: number; itemsPerPage?: number }) {
  if (itemsPerPage) pagination.tickets.pageSize = itemsPerPage
  if (page) pagination.tickets.page = page

  reverting.sessions = true
  const res = await request(
    "/api/users/me/tickets?" +
      new URLSearchParams({
        take: pagination.tickets.pageSize.toString(),
        offset: ((pagination.tickets.page - 1) * pagination.tickets.pageSize).toString(),
      }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    tickets.value = data["data"]
    pagination.tickets.total = data["count"]
  }
  reverting.sessions = false
}

async function readEvents({ page, itemsPerPage }: { page?: number; itemsPerPage?: number }) {
  if (itemsPerPage) pagination.events.pageSize = itemsPerPage
  if (page) pagination.events.page = page

  reverting.events = true
  const res = await request(
    "/api/users/me/events?" +
      new URLSearchParams({
        take: pagination.events.pageSize.toString(),
        offset: ((pagination.events.page - 1) * pagination.events.pageSize).toString(),
      }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    events.value = data["data"]
    pagination.events.total = data["count"]
  }
  reverting.events = false
}

Promise.all([readTickets({}), readEvents({})])

async function killTicket(item: any) {
  reverting.sessions = true
  const res = await request(`/api/users/me/tickets/${item.id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${getAtk()}` },
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await readTickets({})
    error.value = null
  }
  reverting.sessions = false
}
</script>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
