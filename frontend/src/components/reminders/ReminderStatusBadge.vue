<template>
  <span 
    class="badge status-badge"
    :class="statusClass"
    :title="statusTooltip"
  >
    <i :class="statusIcon" class="bi me-1"></i>
    {{ statusLabel }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

// Props
interface Props {
  status: 'pending' | 'active' | 'completed' | 'cancelled';
}

const props = defineProps<Props>();

// Computed properties
const statusClass = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'bg-warning text-dark';
    case 'active':
      return 'bg-info text-dark';
    case 'completed':
      return 'bg-success';
    case 'cancelled':
      return 'bg-secondary';
    default:
      return 'bg-light text-dark';
  }
});

const statusIcon = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'bi-clock';
    case 'active':
      return 'bi-bell';
    case 'completed':
      return 'bi-check-circle';
    case 'cancelled':
      return 'bi-x-circle';
    default:
      return 'bi-question-circle';
  }
});

const statusLabel = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'Pending';
    case 'active':
      return 'Active';
    case 'completed':
      return 'Completed';
    case 'cancelled':
      return 'Cancelled';
    default:
      return 'Unknown';
  }
});

const statusTooltip = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'Waiting for remind time';
    case 'active':
      return 'Currently active';
    case 'completed':
      return 'Reminder completed';
    case 'cancelled':
      return 'Reminder cancelled';
    default:
      return 'Unknown status';
  }
});
</script>

<style scoped>
.status-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.375rem 0.75rem;
  border-radius: 20px;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.status-badge:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
