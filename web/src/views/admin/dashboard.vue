<template>
  <div class="w-full h-full flex justify-center items-center">
    <v-empty-state
      headline="Administration"
      icon="mdi-cog"
      title="What would you like to do today?"
    >
      <v-container>
        <v-row>
          <v-col cols="12" md="6">
            <v-card
              :to="{ name: 'admin.users' }"
              prepend-icon="mdi-account-group"
              text="Manage to help users do something they can't"
              title="Users"
            ></v-card>
          </v-col>

          <v-col cols="12" md="6">
            <v-card
              disabled
              prepend-icon="mdi-comment-quote"
              text="Manage the content on the platform"
              title="Posts & Articles"
            ></v-card>
          </v-col>

          <v-col cols="12" md="6">
            <v-card
              disabled
              prepend-icon="mdi-file-cabinet"
              text="Manage attachments on the platform"
              title="Attachments"
            ></v-card>
          </v-col>

          <v-col cols="12" md="6">
            <v-card
              disabled
              prepend-icon="mdi-ticket"
              text="Solve the tickets issued by users"
              title="Tickets"
            ></v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-empty-state>
  </div>
</template>
<script setup lang="ts">
</script>