<template>
  <div>
    <v-expansion-panels>
      <v-expansion-panel eager title="Challenges">
        <template #text>
          <v-card :loading="reverting.challenges" variant="outlined">
            <v-data-table-server
              density="compact"
              :headers="dataDefinitions.challenges"
              :items="challenges"
              :items-length="pagination.challenges.total"
              :loading="reverting.challenges"
              v-model:items-per-page="pagination.challenges.pageSize"
              @update:options="readChallenges"
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
                </tr>
              </template>
            </v-data-table-server>
          </v-card>
        </template>
      </v-expansion-panel>

      <v-expansion-panel eager title="Sessions">
        <template #text>
          <v-card :loading="reverting.sessions" variant="outlined">
            <v-data-table-server
              density="compact"
              :headers="dataDefinitions.sessions"
              :items="sessions"
              :items-length="pagination.sessions.total"
              :loading="reverting.sessions"
              v-model:items-per-page="pagination.sessions.pageSize"
              @update:options="readSessions"
              item-value="id"
            >
              <template v-slot:item="{ item }: { item: any }">
                <tr>
                  <td>{{ item.id }}</td>
                  <td>
                    <v-chip v-for="value in item.audiences" size="x-small" color="warning" class="capitalize">
                      {{ value }}
                    </v-chip>
                  </td>
                  <td>
                    <v-chip v-for="value in item.claims" size="x-small" color="info" class="font-mono">
                      {{ value }}
                    </v-chip>
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
                          @click="killSession(item)"
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
import { getAtk, useUserinfo } from "@/stores/userinfo"
import { reactive, ref } from "vue"

const id = useUserinfo()

const error = ref<string | null>(null)

const dataDefinitions: { [id: string]: any[] } = {
  challenges: [
    { align: "start", key: "id", title: "ID" },
    { align: "start", key: "ip_address", title: "IP Address" },
    { align: "start", key: "user_agent", title: "User Agent" },
    { align: "start", key: "created_at", title: "Issued At" },
  ],
  sessions: [
    { align: "start", key: "id", title: "ID" },
    { align: "start", key: "audiences", title: "Audiences" },
    { align: "start", key: "claims", title: "Claims" },
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

const challenges = ref<any>([])
const sessions = ref<any>([])
const events = ref<any>([])

const reverting = reactive({ challenges: false, sessions: false, events: false })
const pagination = reactive({
  challenges: { page: 1, pageSize: 5, total: 0 },
  sessions: { page: 1, pageSize: 5, total: 0 },
  events: { page: 1, pageSize: 5, total: 0 },
})

async function readChallenges({ page, itemsPerPage }: { page?: number; itemsPerPage?: number }) {
  if (itemsPerPage) pagination.challenges.pageSize = itemsPerPage
  if (page) pagination.challenges.page = page

  reverting.challenges = true
  const res = await request(
    "/api/users/me/challenges?" +
      new URLSearchParams({
        take: pagination.challenges.pageSize.toString(),
        offset: ((pagination.challenges.page - 1) * pagination.challenges.pageSize).toString(),
      }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    challenges.value = data["data"]
    pagination.challenges.total = data["count"]
  }
  reverting.challenges = false
}

async function readSessions({ page, itemsPerPage }: { page?: number; itemsPerPage?: number }) {
  if (itemsPerPage) pagination.sessions.pageSize = itemsPerPage
  if (page) pagination.sessions.page = page

  reverting.sessions = true
  const res = await request(
    "/api/users/me/sessions?" +
      new URLSearchParams({
        take: pagination.sessions.pageSize.toString(),
        offset: ((pagination.sessions.page - 1) * pagination.sessions.pageSize).toString(),
      }),
    {
      headers: { Authorization: `Bearer ${getAtk()}` },
    },
  )
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    const data = await res.json()
    sessions.value = data["data"]
    pagination.sessions.total = data["count"]
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

Promise.all([readChallenges({}), readSessions({}), readEvents({})])

async function killSession(item: any) {
  reverting.sessions = true
  const res = await request(`/api/users/me/sessions/${item.id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${getAtk()}` },
  })
  if (res.status !== 200) {
    error.value = await res.text()
  } else {
    await readSessions({})
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
