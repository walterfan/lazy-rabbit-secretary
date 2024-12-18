`<template>
  <div class="mb-3">
    <input
      type="text"
      class="form-control"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
  placeholder: string;
}>();

defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();
</script>`