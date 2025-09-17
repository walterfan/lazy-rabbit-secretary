<template>
  <div class="reminder-tags">
    <span
      v-for="tag in tagList"
      :key="tag"
      class="badge tag-badge"
      :class="getTagClass(tag)"
      :title="`Tag: ${tag}`"
    >
      <i class="bi bi-tag-fill me-1"></i>
      {{ tag }}
    </span>
    <span v-if="tagList.length === 0" class="text-muted no-tags">
      <i class="bi bi-dash"></i>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

// Props
interface Props {
  tags: string;
}

const props = defineProps<Props>();

// Computed properties
const tagList = computed(() => {
  if (!props.tags) return [];
  return props.tags
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag.length > 0)
    .slice(0, 3); // Limit to 3 tags for display
});

// Helper functions
const getTagClass = (tag: string): string => {
  // Generate consistent colors based on tag content
  const colors = [
    'bg-primary',
    'bg-success', 
    'bg-info',
    'bg-warning text-dark',
    'bg-danger',
    'bg-secondary',
    'bg-dark'
  ];
  
  // Simple hash function to get consistent color
  let hash = 0;
  for (let i = 0; i < tag.length; i++) {
    const char = tag.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash & hash; // Convert to 32bit integer
  }
  
  const colorIndex = Math.abs(hash) % colors.length;
  return colors[colorIndex];
};
</script>

<style scoped>
.reminder-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  align-items: center;
}

.tag-badge {
  font-size: 0.7rem;
  font-weight: 500;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  transition: all 0.2s ease;
  cursor: default;
}

.tag-badge:hover {
  transform: translateY(-1px);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.no-tags {
  font-size: 0.875rem;
  font-style: italic;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .tag-badge {
    font-size: 0.65rem;
    padding: 0.2rem 0.4rem;
  }
}
</style>
