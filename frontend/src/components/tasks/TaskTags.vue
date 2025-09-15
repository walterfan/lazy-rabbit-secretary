`<template>
  <div class="task-tags-wrapper">
    <span
      v-for="tag in tags"
      :key="tag"
      class="task-tag"
    >
      <i class="bi bi-tag-fill"></i>
      {{ tag }}
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  tags: string[];
}>();
</script>

<style scoped>
.task-tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.task-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background-color: #e9ecef;
  color: #495057;
  padding: 0.25rem 0.5rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.task-tag:hover {
  background-color: #dee2e6;
  transform: translateY(-1px);
}

.task-tag i {
  font-size: 0.65rem;
  opacity: 0.6;
}
</style>`