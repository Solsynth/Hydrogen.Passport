<template>
  <div>
    <v-card>
      <v-img cover class="bg-grey-lighten-2" :height="240" src="/api/avatar" />

      <v-card-text class="flex gap-3.5 px-5 pb-5">
        <v-avatar
          color="grey-lighten-2"
          icon="mdi-account-circle"
          class="rounded-card"
          image="/api/banner"
          :size="54"
        />
        <div>
          <h1 class="text-2xl cursor-pointer" @click="show.realname = !show.realname">{{ displayName }}</h1>
          <p v-html="description"></p>

          <div class="mt-5">
            <p class="opacity-80 desc-line">
              <v-icon icon="mdi-calendar-blank" size="16" />
              <span>Joined at {{ new Date(id.userinfo.data?.created_at)?.toLocaleString() }}</span>
            </p>
            <p class="opacity-80 desc-line">
              <v-icon icon="mdi-cake-variant" size="16" />
              <span>Birthday is {{ new Date(id.userinfo.data?.profile.birthday)?.toLocaleString() }}</span>
            </p>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { useUserinfo } from "@/stores/userinfo"
import { computed } from "vue"
import { reactive } from "vue"
import { parse } from "marked"
import dompurify from "dompurify"

const id = useUserinfo()

const displayName = computed(() => {
  if (show.realname) {
    return (
      (id.userinfo.data?.profile?.first_name ?? "Unknown") + " " + (id.userinfo.data?.profile?.last_name ?? "Unknown")
    )
  } else {
    return id.userinfo.displayName
  }
})
const description = computed(() => {
  if (id.userinfo.data?.description) {
    return dompurify().sanitize(parse(id.userinfo.data?.description) as string)
  } else {
    return "No description yet."
  }
})

const show = reactive({
  realname: false,
})
</script>

<style scoped>
.desc-line {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>

<style>
.rounded-card {
  border-radius: 8px;
}
</style>
