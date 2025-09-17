<template>
  <form @submit.prevent="handleSubmit" class="reminder-form needs-validation" novalidate>
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
            placeholder="Describe your reminder naturally... e.g., 'Remind me to call John tomorrow at 3pm'"
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
      <!-- Reminder Details Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-bell"></i> Reminder Details
        </h6>
        
        <div class="form-grid">
          <div class="form-field full-width">
            <label for="name" class="form-label">
              <i class="bi bi-card-heading text-muted"></i> Reminder Name
              <span v-if="fieldSources.name" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <input
              type="text"
              class="form-control form-control-lg"
              id="name"
              v-model="reminderData.name"
              :class="{ 
                'is-invalid': v$.name.$error,
                'auto-filled': fieldSources.name
              }"
              placeholder="What should I remind you about?"
              required
            />
            <div class="invalid-feedback" v-if="v$.name.$error">
              {{ v$.name.$errors[0].$message }}
            </div>
          </div>
          
          <div class="form-field full-width">
            <label for="content" class="form-label">
              <i class="bi bi-text-paragraph text-muted"></i> Content
              <span v-if="fieldSources.content" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <textarea
              class="form-control"
              id="content"
              v-model="reminderData.content"
              :class="{ 
                'is-invalid': v$.content.$error,
                'auto-filled': fieldSources.content
              }"
              rows="3"
              placeholder="Detailed reminder content..."
              required
            ></textarea>
            <div class="invalid-feedback" v-if="v$.content.$error">
              {{ v$.content.$errors[0].$message }}
            </div>
          </div>
        </div>
      </div>

      <!-- Scheduling Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-calendar-event"></i> Scheduling
        </h6>
        
        <div class="form-grid">
          <div class="form-field">
            <label for="remindTime" class="form-label">
              <i class="bi bi-clock text-muted"></i> Remind Time
              <span v-if="fieldSources.remind_time" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
              <span class="required">*</span>
            </label>
            <input
              type="datetime-local"
              class="form-control"
              id="remindTime"
              v-model="remindTimeString"
              :class="{ 
                'is-invalid': v$.remind_time.$error,
                'auto-filled': fieldSources.remind_time
              }"
              required
            />
            <div class="invalid-feedback" v-if="v$.remind_time.$error">
              {{ v$.remind_time.$errors[0].$message }}
            </div>
          </div>

          <div class="form-field">
            <label for="tags" class="form-label">
              <i class="bi bi-tags text-muted"></i> Tags
              <span v-if="fieldSources.tags" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
            </label>
            <input
              type="text"
              class="form-control"
              id="tags"
              v-model="reminderData.tags"
              :class="{ 
                'auto-filled': fieldSources.tags
              }"
              placeholder="work, personal, meeting..."
            />
            <small class="form-text text-muted">
              Comma-separated tags for organization
            </small>
          </div>
        </div>
      </div>

      <!-- Notification Settings Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-bell-fill"></i> Notification Settings
        </h6>
        
        <div class="form-grid">
          <div class="form-field">
            <label for="remind_methods" class="form-label">
              <i class="bi bi-send text-muted"></i> Notification Methods
              <span v-if="fieldSources.remind_methods" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
            </label>
            <select
              class="form-control form-select"
              id="remind_methods"
              v-model="reminderData.remind_methods"
              :class="{ 
                'auto-filled': fieldSources.remind_methods
              }"
            >
              <option value="email">Email only</option>
              <option value="webhook">Webhook only</option>
              <option value="email,webhook">Email + Webhook</option>
            </select>
            <small class="form-text text-muted">
              Choose how you want to be notified
            </small>
          </div>

          <div class="form-field">
            <label for="remind_targets" class="form-label">
              <i class="bi bi-bullseye text-muted"></i> Notification Targets
              <span v-if="fieldSources.remind_targets" class="auto-fill-badge">
                <i class="bi bi-magic"></i> AI
              </span>
            </label>
            <textarea
              class="form-control"
              id="remind_targets"
              v-model="reminderData.remind_targets"
              :class="{ 
                'auto-filled': fieldSources.remind_targets
              }"
              placeholder="Enter targets based on selected methods:&#10;• Email: user@example.com&#10;• Webhook: https://hooks.slack.com/...&#10;&#10;For multiple targets, separate with commas"
              rows="4"
            />
            <small class="form-text text-muted">
              Specify where notifications should be sent
            </small>
          </div>
        </div>
      </div>

      <!-- Status Section (only for edit mode) -->
      <div v-if="isEditMode" class="form-section">
        <h6 class="section-title">
          <i class="bi bi-flag"></i> Status
        </h6>
        
        <div class="form-field">
          <label for="status" class="form-label">
            <i class="bi bi-circle-fill text-muted"></i> Current Status
          </label>
          <select
            class="form-select"
            id="status"
            v-model="reminderData.status"
          >
            <option value="pending">Pending</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <button type="button" class="btn btn-outline-secondary" @click="$emit('cancel')">
        <i class="bi bi-x-lg me-2"></i>
        Cancel
      </button>
      <button type="submit" class="btn btn-primary" :disabled="loading">
        <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
        <i v-else :class="isEditMode ? 'bi-pencil-square' : 'bi-plus-lg'" class="bi me-2"></i>
        {{ isEditMode ? 'Update Reminder' : 'Create Reminder' }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength } from '@vuelidate/validators';
import type { Reminder, CreateReminderRequest, UpdateReminderRequest } from '@/types';

// Props
interface Props {
  reminder?: Reminder;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
});

// Emits
const emit = defineEmits<{
  (e: 'submit', reminder: CreateReminderRequest | UpdateReminderRequest): void;
  (e: 'cancel'): void;
}>();

// Form data
const reminderData = ref<CreateReminderRequest & { status?: string }>({
  name: '',
  content: '',
  remind_time: new Date(),
  tags: '',
  remind_methods: 'email', // Default to email
  remind_targets: ''
});

// Natural language input
const naturalLanguageInput = ref('');
const parsingFeedback = ref<Array<{ type: 'success' | 'warning' | 'error'; message: string }>>([]);

// Field sources tracking (for AI auto-fill indicators)
const fieldSources = ref<Record<string, boolean>>({
  name: false,
  content: false,
  remind_time: false,
  tags: false,
  remind_methods: false,
  remind_targets: false
});

// Computed properties
const isEditMode = computed(() => !!props.reminder);

const remindTimeString = computed({
  get: () => {
    const date = reminderData.value.remind_time;
    if (!date) return '';
    
    // Convert to local datetime-local format
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  },
  set: (value: string) => {
    if (value) {
      reminderData.value.remind_time = new Date(value);
    }
  }
});

// Validation rules
const rules = {
  name: { required, minLength: minLength(1) },
  content: { required, minLength: minLength(1) },
  remind_time: { 
    required,
    futureDate: (value: Date) => {
      if (!value) return false;
      return value > new Date();
    }
  }
};

const v$ = useVuelidate(rules, reminderData);

// Natural language parsing
const parseNaturalLanguage = () => {
  const input = naturalLanguageInput.value.trim();
  if (!input) return;

  parsingFeedback.value = [];
  
  // Reset field sources
  Object.keys(fieldSources.value).forEach(key => {
    fieldSources.value[key] = false;
  });

  // Simple parsing logic (can be enhanced with more sophisticated NLP)
  try {
    // Extract reminder name (usually the main action)
    const reminderMatch = input.match(/remind me to (.+?)(?:\s+(?:at|on|tomorrow|today|next))/i);
    if (reminderMatch) {
      reminderData.value.name = reminderMatch[1].trim();
      fieldSources.value.name = true;
      parsingFeedback.value.push({
        type: 'success',
        message: `Extracted reminder: "${reminderMatch[1].trim()}"`
      });
    }

    // Extract time information
    const timePatterns = [
      { pattern: /at (\d{1,2}):?(\d{0,2})\s*(am|pm)?/i, type: 'time' },
      { pattern: /(\d{1,2}):(\d{2})\s*(am|pm)?/i, type: 'time' },
      { pattern: /(tomorrow)/i, type: 'relative' },
      { pattern: /(today)/i, type: 'relative' },
      { pattern: /(next week)/i, type: 'relative' }
    ];

    let timeExtracted = false;
    for (const { pattern, type } of timePatterns) {
      const match = input.match(pattern);
      if (match && !timeExtracted) {
        const now = new Date();
        let targetDate = new Date(now);

        if (type === 'time') {
          let hours = parseInt(match[1]);
          const minutes = match[2] ? parseInt(match[2]) : 0;
          const period = match[3]?.toLowerCase();

          if (period === 'pm' && hours !== 12) hours += 12;
          if (period === 'am' && hours === 12) hours = 0;

          targetDate.setHours(hours, minutes, 0, 0);
          
          // If time is in the past today, set for tomorrow
          if (targetDate <= now) {
            targetDate.setDate(targetDate.getDate() + 1);
          }
        } else if (type === 'relative') {
          if (match[1].toLowerCase() === 'tomorrow') {
            targetDate.setDate(targetDate.getDate() + 1);
            targetDate.setHours(9, 0, 0, 0); // Default to 9 AM
          } else if (match[1].toLowerCase() === 'today') {
            targetDate.setHours(targetDate.getHours() + 1, 0, 0, 0); // 1 hour from now
          } else if (match[1].toLowerCase() === 'next week') {
            targetDate.setDate(targetDate.getDate() + 7);
            targetDate.setHours(9, 0, 0, 0); // Default to 9 AM
          }
        }

        reminderData.value.remind_time = targetDate;
        fieldSources.value.remind_time = true;
        timeExtracted = true;
        
        parsingFeedback.value.push({
          type: 'success',
          message: `Scheduled for: ${targetDate.toLocaleString()}`
        });
        break;
      }
    }

    // Use the full input as content if no specific name was extracted
    if (!reminderMatch) {
      reminderData.value.name = input.length > 50 ? input.substring(0, 47) + '...' : input;
      fieldSources.value.name = true;
    }
    
    reminderData.value.content = input;
    fieldSources.value.content = true;

    // Extract tags from common keywords
    const tagKeywords = ['work', 'meeting', 'personal', 'call', 'appointment', 'deadline'];
    const extractedTags = tagKeywords.filter(tag => 
      input.toLowerCase().includes(tag)
    );
    
    if (extractedTags.length > 0) {
      reminderData.value.tags = extractedTags.join(', ');
      fieldSources.value.tags = true;
      parsingFeedback.value.push({
        type: 'success',
        message: `Added tags: ${extractedTags.join(', ')}`
      });
    }

    if (!timeExtracted) {
      parsingFeedback.value.push({
        type: 'warning',
        message: 'Could not extract time information. Please set manually.'
      });
    }

  } catch (error) {
    parsingFeedback.value.push({
      type: 'error',
      message: 'Failed to parse input. Please fill the form manually.'
    });
  }
};

// Form submission
const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (!isValid) return;

  if (isEditMode.value) {
    const updateData: UpdateReminderRequest = {
      name: reminderData.value.name,
      content: reminderData.value.content,
      remind_time: reminderData.value.remind_time,
      tags: reminderData.value.tags,
      remind_methods: reminderData.value.remind_methods,
      remind_targets: reminderData.value.remind_targets,
      status: reminderData.value.status as any
    };
    emit('submit', updateData);
  } else {
    const createData: CreateReminderRequest = {
      name: reminderData.value.name,
      content: reminderData.value.content,
      remind_time: reminderData.value.remind_time,
      tags: reminderData.value.tags,
      remind_methods: reminderData.value.remind_methods,
      remind_targets: reminderData.value.remind_targets
    };
    emit('submit', createData);
  }
};

// Initialize form data when editing
onMounted(() => {
  if (props.reminder) {
    reminderData.value = {
      name: props.reminder.name,
      content: props.reminder.content,
      remind_time: new Date(props.reminder.remind_time),
      tags: props.reminder.tags,
      remind_methods: props.reminder.remind_methods || 'email',
      remind_targets: props.reminder.remind_targets || '',
      status: props.reminder.status
    };
  }
});

// Watch for prop changes
watch(() => props.reminder, (newReminder) => {
  if (newReminder) {
    reminderData.value = {
      name: newReminder.name,
      content: newReminder.content,
      remind_time: new Date(newReminder.remind_time),
      tags: newReminder.tags,
      remind_methods: newReminder.remind_methods || 'email',
      remind_targets: newReminder.remind_targets || '',
      status: newReminder.status
    };
  }
}, { deep: true });
</script>

<style scoped>
/* AI Input Section */
.ai-input-section {
  background: linear-gradient(135deg, #f8f9ff 0%, #e3f2fd 100%);
  border: 2px solid #e3f2fd;
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.ai-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.ai-header h5 {
  color: #1976d2;
  font-weight: 600;
}

.ai-input-wrapper {
  position: relative;
}

.ai-input {
  border: 2px solid #e3f2fd;
  border-radius: 12px;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.ai-input:focus {
  border-color: #1976d2;
  box-shadow: 0 0 0 0.25rem rgba(25, 118, 210, 0.1);
}

.feedback-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.feedback-item {
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  display: flex;
  align-items: center;
}

.feedback-success {
  background-color: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
}

.feedback-warning {
  background-color: #fff8e1;
  color: #f57c00;
  border: 1px solid #ffecb3;
}

.feedback-error {
  background-color: #ffebee;
  color: #c62828;
  border: 1px solid #ffcdd2;
}

/* Form Sections */
.form-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-section {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
}

.section-title {
  color: #495057;
  font-weight: 600;
  margin-bottom: 1.25rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #f8f9fa;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-field.full-width {
  grid-column: 1 / -1;
}

.form-label {
  font-weight: 600;
  color: #495057;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.required {
  color: #dc3545;
}

.auto-fill-badge {
  background: linear-gradient(45deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  padding: 0.125rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  margin-left: auto;
}

.auto-filled {
  background-color: #f0f8ff !important;
  border-color: #4facfe !important;
  animation: glow 0.5s ease-in-out;
}

@keyframes glow {
  0% { box-shadow: 0 0 5px rgba(79, 172, 254, 0.5); }
  50% { box-shadow: 0 0 20px rgba(79, 172, 254, 0.8); }
  100% { box-shadow: 0 0 5px rgba(79, 172, 254, 0.5); }
}

.form-control, .form-select {
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 0.75rem;
  transition: all 0.3s ease;
}

.form-control:focus, .form-select:focus {
  border-color: #4facfe;
  box-shadow: 0 0 0 0.25rem rgba(79, 172, 254, 0.1);
}

.form-control-lg {
  padding: 1rem;
  font-size: 1.1rem;
  font-weight: 500;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e9ecef;
  margin-top: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(45deg, #4facfe 0%, #00f2fe 100%);
  border: none;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
}

/* Transitions */
.feedback-enter-active, .feedback-leave-active {
  transition: all 0.3s ease;
}

.feedback-enter-from, .feedback-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Responsive */
@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
  
  .ai-input-section {
    padding: 1rem;
  }
}
</style>
