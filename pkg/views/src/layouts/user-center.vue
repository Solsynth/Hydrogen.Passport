<template>
  <v-container class="pt-6 px-6">
    <v-row>
      <v-col :cols="12" :xs="12" :sm="12" :md="4" :lg="3">
        <v-card title="Navigation">
          <v-list density="comfortable">
            <v-list-item title="Dashboard" prepend-icon="mdi-view-dashboard" :to="{ name: 'dashboard' }" exact />
            <v-list-item title="Personalize" prepend-icon="mdi-card-bulleted-outline" :to="{ name: 'personalize' }" />
            <v-list-item title="Security" prepend-icon="mdi-security" :to="{ name: 'security' }" />
          </v-list>
        </v-card>
      </v-col>

      <v-col :cols="12" :xs="12" :sm="12" :md="8" :lg="9">
        <router-view />
      </v-col>
    </v-row>
  </v-container>
</template>
