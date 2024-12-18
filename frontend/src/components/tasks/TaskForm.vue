`<template>
  <form @submit.prevent="handleSubmit" class="needs-validation" novalidate>
    <div class="mb-3">
      <label for="name" class="form-label">Name</label>
      <input
        type="text"
        class="form-control"
        id="name"
        v-model="taskData.name"
        :class="{ 'is-invalid': v$.name.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.name.$error">
        {{ v$.name.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="description" class="form-label">Description</label>
      <textarea
        class="form-control"
        id="description"
        v-model="taskData.description"
        :class="{ 'is-invalid': v$.description.$error }"
        required
      ></textarea>
      <div class="invalid-feedback" v-if="v$.description.$error">
        {{ v$.description.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="priority" class="form-label">Priority</label>
      <select
        class="form-select"
        id="priority"
        v-model="taskData.priority"
        :class="{ 'is-invalid': v$.priority.$error }"
        required
      >
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
      </select>
      <div class="invalid-feedback" v-if="v$.priority.$error">
        {{ v$.priority.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="minutes" class="form-label">Duration (minutes)</label>
      <input
        type="number"
        class="form-control"
        id="minutes"
        v-model.number="taskData.minutes"
        :class="{ 'is-invalid': v$.minutes.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.minutes.$error">
        {{ v$.minutes.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="deadline" class="form-label">Deadline</label>
      <input
        type="datetime-local"
        class="form-control"
        id="deadline"
        v-model="taskData.deadline"
        :class="{ 'is-invalid': v$.deadline.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.deadline.$error">
        {{ v$.deadline.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="start_time" class="form-label">Start Time</label>
      <input
        type="datetime-local"
        class="form-control"
        id="start_time"
        v-model="taskData.start_time"
      />
    </div>

    <div class="mb-3">
      <label for="end_time" class="form-label">End Time</label>
      <input
        type="datetime-local"
        class="form-control"
        id="end_time"
        v-model="taskData.end_time"
      />
    </div>

    <div class="mb-3">
      <label for="tags" class="form-label">Tags (comma-separated)</label>
      <input
        type="text"
        class="form-control"
        id="tags"
        v-model="tagsInput"
        placeholder="Enter tags separated by commas"
      />
    </div>

    <button type="submit" class="btn btn-primary">{{ submitButtonText }}</button>
    <button type="button" class="btn btn-secondary ms-2" @click="$emit('cancel')">
      Cancel
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import type { Task } from '@/types';
import { useTaskValidation } from './useTaskValidation';

const props = defineProps<{
  task?: Task;
  submitButtonText?: string;
}>();

const emit = defineEmits<{
  (e: 'submit', task: Task): void;
  (e: 'cancel'): void;
}>();

const taskData = ref<Task>({
  id: crypto.randomUUID(),
  name: '',
  description: '',
  priority: 'medium',
  minutes: 0,
  deadline: new Date(),
  tags: [],
});

const tagsInput = ref('');

const { rules } = useTaskValidation();
const v$ = useVuelidate(rules, taskData);

onMounted(() => {
  if (props.task) {
    taskData.value = { ...props.task };
    tagsInput.value = props.task.tags.join(', ');
  }
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (isValid) {
    taskData.value.tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);
    emit('submit', taskData.value);
  }
};
</script>`