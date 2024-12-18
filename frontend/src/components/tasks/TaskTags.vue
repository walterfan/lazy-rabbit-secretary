`<template>
  <span
    v-for="tag in tags"
    :key="tag"
    class="badge bg-secondary me-1"
  >
    {{ tag }}
  </span>
</template>

<script setup lang="ts">
defineProps<{
  tags: string[];
}>();
</script>`