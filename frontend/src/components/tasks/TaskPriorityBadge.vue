`<template>
  <span :class="badgeClass">
    {{ priority }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  priority: 'low' | 'medium' | 'high';
}>();

const badgeClass = computed(() => {
  const classes = {
    low: 'badge bg-success',
    medium: 'badge bg-warning',
    high: 'badge bg-danger'
  };
  return classes[props.priority];
});
</script>`