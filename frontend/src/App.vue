<template>
  <q-toolbar class="bg-primary text-white shadow-2">
    <q-toolbar-title shrink>ClickHouse Executor</q-toolbar-title>
    <router-link custom to="/" v-slot:default="{ href, isExactActive }">
      <q-btn
        :outline="isExactActive"
        :flat="!isExactActive"
        rounded
        label="Query Executor"
        :to="href"
      />
    </router-link>
  </q-toolbar>
  <router-view />
</template>
