`<template>
  <span :class="badgeClass">
    <i :class="iconClass"></i>
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

const iconClass = computed(() => {
  const icons = {
    low: 'bi bi-arrow-down',
    medium: 'bi bi-dash',
    high: 'bi bi-arrow-up'
  };
  return icons[props.priority];
});
</script>

<style scoped>
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.badge i {
  font-size: 0.7rem;
}
</style>`