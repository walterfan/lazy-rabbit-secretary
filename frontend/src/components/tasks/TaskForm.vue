`<template>
  <form @submit.prevent="handleSubmit" class="task-form needs-validation" novalidate>
    <!-- Natural Language Input Section -->
    <div class="ai-input-section">
      <div class="ai-header">
        <h5 class="mb-0">
          <i class="bi bi-stars"></i> AI Quick Input
        </h5>
        <span class="badge bg-primary">Beta</span>
      </div>
      
      <div class="ai-input-wrapper">
        <div class="input-group input-group-lg">
          <span class="input-group-text bg-white">
            <i class="bi bi-chat-dots text-primary"></i>
          </span>
          <input
            type="text"
            class="form-control ai-input"
            id="nlInput"
            v-model="naturalLanguageInput"
            placeholder="Describe your task naturally... e.g., 'Team meeting tomorrow at 2pm for 1 hour'"
            @keyup.enter="parseNaturalLanguage"
          />
          <button 
            class="btn btn-primary px-4" 
            type="button"
            @click="parseNaturalLanguage"
            :disabled="!naturalLanguageInput.trim()"
          >
            <i class="bi bi-lightning-fill me-2"></i>
            Parse
          </button>
        </div>
        
        <!-- Parsing Feedback -->
        <transition-group name="feedback" tag="div" class="feedback-container mt-3">
          <div 
            v-for="(feedback, index) in parsingFeedback" 
            :key="`feedback-${index}`"
            class="feedback-item"
            :class="`feedback-${feedback.type}`"
          >
            <i :class="{
              'bi-check-circle-fill': feedback.type === 'success',
              'bi-exclamation-triangle-fill': feedback.type === 'warning',
              'bi-x-circle-fill': feedback.type === 'error'
            }" class="bi me-2"></i>
            {{ feedback.message }}
          </div>
        </transition-group>
      </div>
    </div>

    <!-- Main Form Content -->
    <div class="form-content">
      <!-- Task Details Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-pencil-square"></i> Task Details
        </h6>
        
        <div class="form-grid">
          <div class="form-field full-width">
            <label for="name" class="form-label">
              <i class="bi bi-card-heading text-muted"></i> Task Name
              <span v-if="fieldSources.name" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <input
              type="text"
              class="form-control form-control-lg"
              id="name"
              v-model="taskData.name"
              :class="{ 
                'is-invalid': v$.name.$error,
                'auto-filled': fieldSources.name
              }"
              placeholder="What needs to be done?"
              required
            />
            <div class="invalid-feedback" v-if="v$.name.$error">
              {{ v$.name.$errors[0].$message }}
            </div>
          </div>
          
          <div class="form-field full-width">
            <label for="description" class="form-label">
              <i class="bi bi-text-paragraph text-muted"></i> Description
              <span v-if="fieldSources.description" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <textarea
              class="form-control"
              id="description"
              v-model="taskData.description"
              :class="{ 
                'is-invalid': v$.description.$error,
                'auto-filled': fieldSources.description
              }"
              rows="3"
              placeholder="Add more details about this task..."
              required
            ></textarea>
            <div class="invalid-feedback" v-if="v$.description.$error">
              {{ v$.description.$errors[0].$message }}
            </div>
          </div>
        </div>
      </div>

      <!-- Priority, Difficulty & Duration Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-flag"></i> Priority, Difficulty & Duration
        </h6>
        
        <div class="form-grid three-columns">
          <div class="form-field">
            <label for="priority" class="form-label">
              <i class="bi bi-exclamation-circle text-muted"></i> Priority
              <span v-if="fieldSources.priority" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <div class="priority-selector">
              <input type="radio" class="btn-check" name="priority" id="priority-low" value="low" v-model="taskData.priority">
              <label class="btn btn-outline-success" for="priority-low">
                <i class="bi bi-arrow-down"></i> Low
              </label>
              
              <input type="radio" class="btn-check" name="priority" id="priority-medium" value="medium" v-model="taskData.priority">
              <label class="btn btn-outline-warning" for="priority-medium">
                <i class="bi bi-dash"></i> Medium
              </label>
              
              <input type="radio" class="btn-check" name="priority" id="priority-high" value="high" v-model="taskData.priority">
              <label class="btn btn-outline-danger" for="priority-high">
                <i class="bi bi-arrow-up"></i> High
              </label>
            </div>
            <div class="invalid-feedback d-block" v-if="v$.priority.$error">
              {{ v$.priority.$errors[0].$message }}
            </div>
          </div>

          <div class="form-field">
            <label for="difficulty" class="form-label">
              <i class="bi bi-speedometer text-muted"></i> Difficulty
              <span v-if="fieldSources.difficulty" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <div class="difficulty-selector">
              <input type="radio" class="btn-check" name="difficulty" id="difficulty-easy" value="easy" v-model="taskData.difficulty">
              <label class="btn btn-outline-info" for="difficulty-easy">
                <i class="bi bi-emoji-smile"></i> Easy
              </label>
              
              <input type="radio" class="btn-check" name="difficulty" id="difficulty-medium" value="medium" v-model="taskData.difficulty">
              <label class="btn btn-outline-primary" for="difficulty-medium">
                <i class="bi bi-emoji-neutral"></i> Medium
              </label>
              
              <input type="radio" class="btn-check" name="difficulty" id="difficulty-hard" value="hard" v-model="taskData.difficulty">
              <label class="btn btn-outline-dark" for="difficulty-hard">
                <i class="bi bi-emoji-dizzy"></i> Hard
              </label>
            </div>
            <div class="invalid-feedback d-block" v-if="v$.difficulty && v$.difficulty.$error">
              {{ v$.difficulty.$errors[0].$message }}
            </div>
          </div>

          <div class="form-field">
            <label for="minutes" class="form-label">
              <i class="bi bi-clock text-muted"></i> Duration
              <span v-if="fieldSources.minutes" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <div class="input-group">
              <input
                type="number"
                class="form-control"
                id="minutes"
                v-model.number="taskData.minutes"
                :class="{ 
                  'is-invalid': v$.minutes.$error,
                  'auto-filled': fieldSources.minutes
                }"
                placeholder="30"
                min="1"
                required
              />
              <span class="input-group-text">minutes</span>
              <div class="invalid-feedback" v-if="v$.minutes.$error">
                {{ v$.minutes.$errors[0].$message }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Schedule Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-calendar-event"></i> Schedule
        </h6>
        
        <div class="form-grid">
          <div class="form-field">
            <label for="schedule_time" class="form-label">
              <i class="bi bi-calendar-week text-muted"></i> Schedule Time
              <span v-if="fieldSources.schedule_time" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <input
              type="datetime-local"
              class="form-control"
              id="schedule_time"
              v-model="scheduleTimeFormatted"
              :class="{ 
                'is-invalid': v$.schedule_time && v$.schedule_time.$error,
                'auto-filled': fieldSources.schedule_time
              }"
              required
            />
            <div class="invalid-feedback" v-if="v$.schedule_time && v$.schedule_time.$error">
              {{ v$.schedule_time.$errors[0].$message }}
            </div>
            <small class="text-muted">When should this task be scheduled?</small>
          </div>

          <div class="form-field">
            <label for="deadline" class="form-label">
              <i class="bi bi-calendar-check text-muted"></i> Deadline
              <span v-if="fieldSources.deadline" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <input
              type="datetime-local"
              class="form-control"
              id="deadline"
              v-model="deadlineFormatted"
              :class="{ 
                'is-invalid': v$.deadline.$error,
                'auto-filled': fieldSources.deadline
              }"
              required
            />
            <div class="invalid-feedback" v-if="v$.deadline.$error">
              {{ v$.deadline.$errors[0].$message }}
            </div>
            <small class="text-muted">When must this task be completed?</small>
          </div>
        </div>

        <div class="form-grid mt-3">
          <div class="form-field">
            <label for="start_time" class="form-label">
              <i class="bi bi-play-circle text-muted"></i> Start Time
              <span v-if="fieldSources.start_time" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
            </label>
            <input
              type="datetime-local"
              class="form-control"
              id="start_time"
              v-model="startTimeFormatted"
              :class="{ 
                'auto-filled': fieldSources.start_time
              }"
            />
            <small class="text-muted">Optional: Actual start time</small>
          </div>

          <div class="form-field">
            <label for="end_time" class="form-label">
              <i class="bi bi-stop-circle text-muted"></i> End Time
            </label>
            <input
              type="datetime-local"
              class="form-control"
              id="end_time"
              v-model="endTimeFormatted"
            />
            <small class="text-muted">Optional: Actual end time</small>
          </div>
        </div>
      </div>

      <!-- Tags Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-tags"></i> Tags & Categories
        </h6>
        
        <div class="form-field">
          <label for="tags" class="form-label">
            <i class="bi bi-tag text-muted"></i> Tags
            <span v-if="fieldSources.tags" class="auto-fill-badge">
              <i class="bi bi-magic"></i> AI
            </span>
          </label>
          <div class="tags-input-wrapper">
            <input
              type="text"
              class="form-control"
              id="tags"
              v-model="tagsInput"
              placeholder="Add tags... (comma-separated)"
              :class="{ 
                'auto-filled': fieldSources.tags
              }"
            />
            <small class="text-muted mt-1">
              <i class="bi bi-info-circle"></i> Use tags to organize and find tasks easily
            </small>
          </div>
        </div>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <button type="submit" class="btn btn-primary btn-lg">
        <i class="bi bi-check-lg me-2"></i>
        {{ submitButtonText || 'Create Task' }}
      </button>
      <button type="button" class="btn btn-light btn-lg" @click="$emit('cancel')">
        Cancel
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import type { Task } from '@/types';
import { useTaskValidation } from './useTaskValidation';
import { parseISO, addDays, addHours, startOfDay, setHours, setMinutes } from 'date-fns';
import { useAuthStore } from '@/stores/authStore';

const authStore = useAuthStore();

// Helper function to get headers with optional authentication
const getHeaders = () => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  if (authStore.token) {
    headers.Authorization = `Bearer ${authStore.token}`;
  }
  
  return headers;
};

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
  difficulty: 'medium',
  status: 'pending',
  minutes: 0,
  deadline: new Date(),
  schedule_time: new Date(),
  tags: [],
});

const tagsInput = ref('');
const naturalLanguageInput = ref('');

// Track which fields were auto-filled
const fieldSources = ref<Record<string, boolean>>({
  name: false,
  description: false,
  priority: false,
  difficulty: false,
  minutes: false,
  deadline: false,
  schedule_time: false,
  start_time: false,
  tags: false
});

// Parsing feedback messages
interface FeedbackMessage {
  type: 'success' | 'warning' | 'error';
  message: string;
}
const parsingFeedback = ref<FeedbackMessage[]>([]);

const { rules } = useTaskValidation();
const v$ = useVuelidate(rules, taskData);

// Format date for datetime-local input
const formatDateForInput = (date: Date | undefined): string => {
  if (!date) return '';
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
};

// Computed properties for datetime-local inputs
const deadlineFormatted = computed({
  get: () => formatDateForInput(taskData.value.deadline),
  set: (value: string) => {
    if (value) taskData.value.deadline = new Date(value);
  }
});

const scheduleTimeFormatted = computed({
  get: () => formatDateForInput(taskData.value.schedule_time),
  set: (value: string) => {
    if (value) taskData.value.schedule_time = new Date(value);
  }
});

const startTimeFormatted = computed({
  get: () => formatDateForInput(taskData.value.start_time),
  set: (value: string) => {
    if (value) taskData.value.start_time = new Date(value);
  }
});

const endTimeFormatted = computed({
  get: () => formatDateForInput(taskData.value.end_time),
  set: (value: string) => {
    if (value) taskData.value.end_time = new Date(value);
  }
});

onMounted(() => {
  if (props.task) {
    taskData.value = { ...props.task };
    tagsInput.value = props.task.tags.join(', ');
  }
});

// Natural language parsing function
const parseNaturalLanguage = async () => {
  if (!naturalLanguageInput.value.trim()) return;
  
  // Reset feedback and field sources
  parsingFeedback.value = [];
  Object.keys(fieldSources.value).forEach(key => {
    fieldSources.value[key] = false;
  });
  
  const input = naturalLanguageInput.value.toLowerCase();
  const originalInput = naturalLanguageInput.value;
  
  try {
    // First, try to call the backend LLM API for parsing
    const response = await fetch('/api/v1/tasks/parse', {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ text: originalInput })
    });
    
    if (response.ok) {
      const parsedData = await response.json();
      applyParsedData(parsedData);
      return;
    }
  } catch (error) {
    console.log('Backend parsing failed, using client-side parsing:', error);
  }
  
  // Fallback to client-side parsing
  clientSideParsing(input, originalInput);
};

// Apply parsed data from backend
const applyParsedData = (parsedData: any) => {
  if (parsedData.name) {
    taskData.value.name = parsedData.name;
    fieldSources.value.name = true;
    parsingFeedback.value.push({ type: 'success', message: 'Task name identified' });
  }
  
  if (parsedData.description) {
    taskData.value.description = parsedData.description;
    fieldSources.value.description = true;
    parsingFeedback.value.push({ type: 'success', message: 'Description parsed' });
  }
  
  if (parsedData.priority) {
    taskData.value.priority = parsedData.priority;
    fieldSources.value.priority = true;
    parsingFeedback.value.push({ type: 'success', message: `Priority set to ${parsedData.priority}` });
  }
  
  if (parsedData.difficulty) {
    taskData.value.difficulty = parsedData.difficulty;
    fieldSources.value.difficulty = true;
    parsingFeedback.value.push({ type: 'success', message: `Difficulty set to ${parsedData.difficulty}` });
  }
  
  if (parsedData.minutes) {
    taskData.value.minutes = parsedData.minutes;
    fieldSources.value.minutes = true;
    parsingFeedback.value.push({ type: 'success', message: `Duration set to ${parsedData.minutes} minutes` });
  }
  
  if (parsedData.schedule_time) {
    taskData.value.schedule_time = new Date(parsedData.schedule_time);
    fieldSources.value.schedule_time = true;
    parsingFeedback.value.push({ type: 'success', message: 'Schedule time parsed' });
  }
  
  if (parsedData.deadline) {
    taskData.value.deadline = new Date(parsedData.deadline);
    fieldSources.value.deadline = true;
    parsingFeedback.value.push({ type: 'success', message: 'Deadline parsed' });
  }
  
  if (parsedData.start_time) {
    taskData.value.start_time = new Date(parsedData.start_time);
    fieldSources.value.start_time = true;
    parsingFeedback.value.push({ type: 'success', message: 'Start time parsed' });
  }
  
  if (parsedData.tags && parsedData.tags.length > 0) {
    tagsInput.value = parsedData.tags.join(', ');
    fieldSources.value.tags = true;
    parsingFeedback.value.push({ type: 'success', message: 'Tags identified' });
  }
};

// Client-side parsing logic
const clientSideParsing = (input: string, originalInput: string) => {
  // Parse priority
  if (input.includes('high priority') || input.includes('urgent') || input.includes('asap')) {
    taskData.value.priority = 'high';
    fieldSources.value.priority = true;
    parsingFeedback.value.push({ type: 'success', message: 'Priority set to high' });
  } else if (input.includes('low priority') || input.includes('whenever')) {
    taskData.value.priority = 'low';
    fieldSources.value.priority = true;
    parsingFeedback.value.push({ type: 'success', message: 'Priority set to low' });
  } else {
    parsingFeedback.value.push({ type: 'warning', message: 'Priority not specified, using default: medium' });
  }
  
  // Parse difficulty
  if (input.includes('easy') || input.includes('simple') || input.includes('quick')) {
    taskData.value.difficulty = 'easy';
    fieldSources.value.difficulty = true;
    parsingFeedback.value.push({ type: 'success', message: 'Difficulty set to easy' });
  } else if (input.includes('hard') || input.includes('difficult') || input.includes('complex') || input.includes('challenging')) {
    taskData.value.difficulty = 'hard';
    fieldSources.value.difficulty = true;
    parsingFeedback.value.push({ type: 'success', message: 'Difficulty set to hard' });
  } else {
    parsingFeedback.value.push({ type: 'warning', message: 'Difficulty not specified, using default: medium' });
  }
  
  // Parse duration
  const durationMatch = input.match(/(\d+)\s*(hour|hr|minute|min)/);
  if (durationMatch) {
    const value = parseInt(durationMatch[1]);
    const unit = durationMatch[2];
    if (unit.startsWith('hour') || unit.startsWith('hr')) {
      taskData.value.minutes = value * 60;
    } else {
      taskData.value.minutes = value;
    }
    fieldSources.value.minutes = true;
    parsingFeedback.value.push({ type: 'success', message: `Duration set to ${taskData.value.minutes} minutes` });
  } else {
    taskData.value.minutes = 30; // Default 30 minutes
    parsingFeedback.value.push({ type: 'warning', message: 'Duration not specified, using default: 30 minutes' });
  }
  
  // Parse date and time
  const now = new Date();
  let deadlineSet = false;
  let scheduleTimeSet = false;
  let startTimeSet = false;
  
  // Tomorrow
  if (input.includes('tomorrow')) {
    const tomorrow = addDays(now, 1);
    taskData.value.deadline = tomorrow;
    taskData.value.schedule_time = tomorrow;
    deadlineSet = true;
    scheduleTimeSet = true;
    fieldSources.value.deadline = true;
    fieldSources.value.schedule_time = true;
    
    // Parse time for tomorrow
    const timeMatch = input.match(/(\d{1,2})\s*(am|pm|:)/);
    if (timeMatch) {
      let hours = parseInt(timeMatch[1]);
      if (timeMatch[2] === 'pm' && hours < 12) hours += 12;
      if (timeMatch[2] === 'am' && hours === 12) hours = 0;
      
      const withTime = setHours(setMinutes(tomorrow, 0), hours);
      taskData.value.deadline = withTime;
      taskData.value.schedule_time = withTime;
      taskData.value.start_time = withTime;
      startTimeSet = true;
      fieldSources.value.start_time = true;
      parsingFeedback.value.push({ type: 'success', message: `Scheduled for tomorrow at ${hours}:00` });
    } else {
      parsingFeedback.value.push({ type: 'success', message: 'Deadline and schedule time set to tomorrow' });
    }
  }
  
  // Today with time
  else if (input.match(/today|(\d{1,2})\s*(am|pm)/)) {
    const timeMatch = input.match(/(\d{1,2})\s*(am|pm)/);
    if (timeMatch) {
      let hours = parseInt(timeMatch[1]);
      if (timeMatch[2] === 'pm' && hours < 12) hours += 12;
      if (timeMatch[2] === 'am' && hours === 12) hours = 0;
      
      const withTime = setHours(setMinutes(now, 0), hours);
      taskData.value.deadline = withTime;
      taskData.value.schedule_time = withTime;
      taskData.value.start_time = withTime;
      deadlineSet = true;
      scheduleTimeSet = true;
      startTimeSet = true;
      fieldSources.value.deadline = true;
      fieldSources.value.schedule_time = true;
      fieldSources.value.start_time = true;
      parsingFeedback.value.push({ type: 'success', message: `Scheduled for today at ${hours}:00` });
    }
  }
  
  // Next week
  else if (input.includes('next week')) {
    const nextWeek = addDays(now, 7);
    taskData.value.deadline = nextWeek;
    taskData.value.schedule_time = nextWeek;
    deadlineSet = true;
    scheduleTimeSet = true;
    fieldSources.value.deadline = true;
    fieldSources.value.schedule_time = true;
    parsingFeedback.value.push({ type: 'success', message: 'Deadline and schedule time set to next week' });
  }
  
  // Default deadline and schedule time if not set
  if (!deadlineSet) {
    const tomorrow = addDays(now, 1);
    taskData.value.deadline = tomorrow;
    taskData.value.schedule_time = tomorrow;
    parsingFeedback.value.push({ type: 'warning', message: 'Deadline and schedule time not specified, using default: tomorrow' });
  }
  
  // Parse task name and description
  // Remove parsed elements to get the core task description
  let cleanedInput = originalInput
    .replace(/tomorrow|today|next week/gi, '')
    .replace(/\d+\s*(hour|hr|minute|min)/gi, '')
    .replace(/\d+\s*(am|pm)/gi, '')
    .replace(/high priority|low priority|urgent|asap|whenever/gi, '')
    .replace(/at\s+/gi, '')
    .replace(/for\s+/gi, '')
    .trim();
  
  // Extract potential tags (words starting with # or common task categories)
  const tagMatches = cleanedInput.match(/#\w+/g) || [];
  const categoryWords = ['meeting', 'call', 'email', 'review', 'task', 'project'];
  const foundCategories = categoryWords.filter(cat => input.includes(cat));
  
  if (tagMatches.length > 0 || foundCategories.length > 0) {
    const allTags = [...tagMatches.map(t => t.substring(1)), ...foundCategories];
    tagsInput.value = allTags.join(', ');
    fieldSources.value.tags = true;
    parsingFeedback.value.push({ type: 'success', message: `Tags identified: ${allTags.join(', ')}` });
    
    // Remove tags from the cleaned input
    tagMatches.forEach(tag => {
      cleanedInput = cleanedInput.replace(tag, '');
    });
  }
  
  // Set name and description
  if (cleanedInput) {
    // First 50 characters as name
    taskData.value.name = cleanedInput.substring(0, 50);
    fieldSources.value.name = true;
    
    // Full text as description if longer than 50 chars
    if (cleanedInput.length > 50) {
      taskData.value.description = cleanedInput;
      fieldSources.value.description = true;
      parsingFeedback.value.push({ type: 'success', message: 'Task name and description set' });
    } else {
      taskData.value.description = cleanedInput;
      fieldSources.value.description = true;
      parsingFeedback.value.push({ type: 'success', message: 'Task name set' });
    }
  } else {
    parsingFeedback.value.push({ type: 'error', message: 'Could not extract task name from input' });
  }
};

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
</script>

<style scoped>
/* Main Form Container */
.task-form {
  max-width: 900px;
  margin: 0 auto;
}

/* AI Input Section */
.ai-input-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.ai-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  color: white;
}

.ai-header h5 {
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.ai-header .badge {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
}

.ai-input-wrapper {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
}

.ai-input {
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1.1rem;
  padding: 0.75rem 1rem;
  transition: all 0.3s ease;
}

.ai-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 0.2rem rgba(102, 126, 234, 0.1);
}

.input-group-text {
  border: 2px solid #e9ecef;
  border-right: none;
}

/* Feedback Messages */
.feedback-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.feedback-item {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
}

.feedback-success {
  background-color: #d1fae5;
  color: #065f46;
}

.feedback-warning {
  background-color: #fef3c7;
  color: #92400e;
}

.feedback-error {
  background-color: #fee2e2;
  color: #991b1b;
}

/* Feedback animations */
.feedback-enter-active {
  animation: slideIn 0.3s ease;
}

.feedback-leave-active {
  animation: slideOut 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideOut {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-10px);
  }
}

/* Form Content */
.form-content {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

/* Form Sections */
.form-section {
  margin-bottom: 2.5rem;
  padding-bottom: 2.5rem;
  border-bottom: 1px solid #e9ecef;
}

.form-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.section-title {
  color: #495057;
  font-weight: 600;
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Form Grid */
.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.form-field {
  position: relative;
}

.form-field.full-width {
  grid-column: 1 / -1;
}

/* Form Labels */
.form-label {
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.required {
  color: #dc3545;
  font-weight: normal;
}

.auto-fill-badge {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-size: 0.7rem;
  padding: 0.15rem 0.4rem;
  border-radius: 12px;
  font-weight: normal;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

/* Form Controls */
.form-control, .form-select {
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 0.625rem 0.875rem;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-control:focus, .form-select:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 0.2rem rgba(102, 126, 234, 0.1);
}

.form-control.auto-filled {
  background-color: #f0f9ff;
  border-color: #667eea;
}

.form-control-lg {
  font-size: 1.1rem;
  padding: 0.75rem 1rem;
}

textarea.form-control {
  resize: vertical;
}

/* Priority Selector */
.priority-selector, .difficulty-selector {
  display: flex;
  gap: 0.5rem;
}

.priority-selector .btn, .difficulty-selector .btn {
  flex: 1;
  padding: 0.5rem 1rem;
  font-weight: 500;
  border-width: 2px;
  transition: all 0.3s ease;
}

.priority-selector .btn:hover, .difficulty-selector .btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Three column grid for priority, difficulty, duration */
.form-grid.three-columns {
  grid-template-columns: repeat(3, 1fr);
}

@media (max-width: 992px) {
  .form-grid.three-columns {
    grid-template-columns: 1fr;
  }
}

/* Input Groups */
.input-group {
  position: relative;
}

.input-group-text {
  background-color: #f8f9fa;
  border: 2px solid #e9ecef;
  color: #6c757d;
  font-weight: 500;
}

/* Tags Input */
.tags-input-wrapper small {
  display: block;
  color: #6c757d;
}

/* Form Actions */
.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e9ecef;
}

.form-actions .btn {
  padding: 0.75rem 2rem;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.form-actions .btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.form-actions .btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.form-actions .btn-light {
  background-color: #f8f9fa;
  border: 2px solid #e9ecef;
  color: #495057;
}

.form-actions .btn-light:hover {
  background-color: #e9ecef;
  border-color: #dee2e6;
}

/* Responsive Design */
@media (max-width: 768px) {
  .task-form {
    padding: 0 1rem;
  }
  
  .ai-input-section {
    padding: 1.5rem;
  }
  
  .form-content {
    padding: 1.5rem;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions .btn {
    width: 100%;
  }
}

/* Validation Feedback */
.invalid-feedback {
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.is-invalid {
  border-color: #dc3545 !important;
}

/* Smooth Transitions */
* {
  transition: border-color 0.3s ease, background-color 0.3s ease, box-shadow 0.3s ease;
}
</style>`