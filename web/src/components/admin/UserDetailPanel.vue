<template>
  <v-dialog class="max-w-[720px]" :model-value="props.data != null"
            @update:model-value="(val: boolean) => !val && emits('close')">
    <template v-slot:default="{ isActive }">
      <v-card :title="`User @${props.data?.name}`">
        <v-card-text>
          <v-row>
            <v-col cols="12" md="6">
              <h4 class="field-title">Name</h4>
              <p>{{ props.data?.name }}</p>
            </v-col>
            <v-col cols="12" md="6">
              <h4 class="field-title">Nick</h4>
              <p>{{ props.data?.nick }}</p>
            </v-col>
            <v-col cols="12">
              <h4 class="field-title">Entire Payload</h4>
              <v-code class="font-mono overflow-x-scroll max-h-[360px]">
                <pre>{{ JSON.stringify(props.data, null, 4) }}</pre>
              </v-code>
            </v-col>
          </v-row>
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
const props = defineProps<{ data: any }>()
const emits = defineEmits(["close"])
</script>

<style scoped>
.field-title {
  font-weight: bold;
}
</style>