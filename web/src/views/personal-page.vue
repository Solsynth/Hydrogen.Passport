<template>
  <div>
    <v-card class="mb-3" title="Design" prepend-icon="mdi-pencil-ruler" :loading="loading">
      <template #text>
        <v-form class="mt-1" @submit.prevent="submit">
          <v-row dense>
            <v-col :cols="12">
              <v-textarea hide-details label="Content" density="comfortable" variant="outlined"
                          v-model="data.content" />
            </v-col>
          </v-row>

          <v-btn type="submit" class="mt-2" variant="text" prepend-icon="mdi-content-save" :disabled="loading">
            Apply Changes
          </v-btn>
        </v-form>
      </template>
    </v-card>

    <v-snackbar v-model="done" :timeout="3000"> Your personal page has been updated.</v-snackbar>

    <!-- @vue-ignore -->
    <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { getAtk } from "@/stores/userinfo";
import { request } from "@/scripts/request";

const error = ref<string | null>(null);
const done = ref(false);
const loading = ref(false);

const data = ref<any>({});

async function read() {
  loading.value = true;
  const res = await request("/api/users/me/page", {
    headers: { Authorization: `Bearer ${(getAtk())}` }
  });
  if (res.status !== 200) {
    error.value = await res.text();
  } else {
    data.value = await res.json();
  }
  loading.value = false;
}

async function submit() {
  const payload = data.value;

  loading.value = true;
  const res = await request("/api/users/me/page", {
    method: "PUT",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  });
  if (res.status !== 200) {
    error.value = await res.text();
  } else {
    await read();
    done.value = true;
    error.value = null;
  }
  loading.value = false;
}

read();
</script>
