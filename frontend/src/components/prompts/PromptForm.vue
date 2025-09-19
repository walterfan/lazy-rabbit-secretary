<template>
  <form @submit.prevent="handleSubmit" class="prompt-form needs-validation" novalidate>
    <!-- Prompt Details Section -->
    <div class="form-section">
      <h6 class="section-title">
        <i class="bi bi-chat-dots"></i> Prompt Details
      </h6>
      
      <div class="form-grid">
        <div class="form-field full-width">
          <label for="name" class="form-label">
            <i class="bi bi-card-heading text-muted"></i> Prompt Name
            <span class="required">*</span>
          </label>
          <input
            type="text"
            class="form-control"
            id="name"
            v-model="promptData.name"
            :class="{ 'is-invalid': v$.name.$error }"
            placeholder="Enter prompt name"
            required
          />
          <div class="invalid-feedback" v-if="v$.name.$error">
            {{ v$.name.$errors[0].$message }}
          </div>
        </div>
        
        <div class="form-field full-width">
          <label for="description" class="form-label">
            <i class="bi bi-text-paragraph text-muted"></i> Description
          </label>
          <textarea
            class="form-control"
            id="description"
            v-model="promptData.description"
            rows="2"
            placeholder="Enter prompt description..."
          ></textarea>
        </div>
      </div>
    </div>

    <!-- System Prompt Section -->
    <div class="form-section">
      <h6 class="section-title">
        <i class="bi bi-gear"></i> System Prompt
      </h6>
      
      <div class="form-field full-width">
        <label for="system_prompt" class="form-label">
          <i class="bi bi-cpu text-muted"></i> System Prompt
          <span class="required">*</span>
        </label>
        <textarea
          class="form-control"
          id="system_prompt"
          v-model="promptData.system_prompt"
          :class="{ 'is-invalid': v$.system_prompt.$error }"
          rows="6"
          placeholder="Enter the system prompt that defines the AI's role and behavior..."
          required
        ></textarea>
        <div class="invalid-feedback" v-if="v$.system_prompt.$error">
          {{ v$.system_prompt.$errors[0].$message }}
        </div>
        <small class="text-muted">
          <i class="bi bi-info-circle"></i> This prompt sets the AI's role, personality, and instructions
        </small>
      </div>
    </div>

    <!-- User Prompt Section -->
    <div class="form-section">
      <h6 class="section-title">
        <i class="bi bi-person"></i> User Prompt Template
      </h6>
      
      <div class="form-field full-width">
        <label for="user_prompt" class="form-label">
          <i class="bi bi-chat-square-text text-muted"></i> User Prompt Template
          <span class="required">*</span>
        </label>
        <textarea
          class="form-control"
          id="user_prompt"
          v-model="promptData.user_prompt"
          :class="{ 'is-invalid': v$.user_prompt.$error }"
          rows="6"
          placeholder="Enter the user prompt template with placeholders like {input} or {context}..."
          required
        ></textarea>
        <div class="invalid-feedback" v-if="v$.user_prompt.$error">
          {{ v$.user_prompt.$errors[0].$message }}
        </div>
        <small class="text-muted">
          <i class="bi bi-info-circle"></i> Use placeholders like {input}, {context}, {data} for dynamic content
        </small>
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
        </label>
        <div class="tags-input-wrapper">
          <input
            type="text"
            class="form-control"
            id="tags"
            v-model="tagsInput"
            placeholder="Add tags... (comma-separated)"
          />
          <small class="text-muted mt-1">
            <i class="bi bi-info-circle"></i> Use tags to organize and find prompts easily
          </small>
        </div>
      </div>
    </div>

    <!-- Preview Section -->
    <div class="form-section" v-if="promptData.system_prompt || promptData.user_prompt">
      <h6 class="section-title">
        <i class="bi bi-eye"></i> Preview
      </h6>
      
      <div class="preview-container">
        <div class="preview-section">
          <h6 class="preview-title">
            <i class="bi bi-gear"></i> System Prompt
          </h6>
          <div class="preview-content">
            {{ promptData.system_prompt || 'No system prompt defined' }}
          </div>
        </div>
        
        <div class="preview-section">
          <h6 class="preview-title">
            <i class="bi bi-person"></i> User Prompt Template
          </h6>
          <div class="preview-content">
            {{ promptData.user_prompt || 'No user prompt template defined' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <button type="submit" class="btn btn-primary btn-lg">
        <i class="bi bi-check-lg me-2"></i>
        {{ submitButtonText || 'Create Prompt' }}
      </button>
      <button type="button" class="btn btn-light btn-lg" @click="$emit('cancel')">
        Cancel
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength } from '@vuelidate/validators';
import type { Prompt } from '@/types';

const props = defineProps<{
  prompt?: Prompt;
  submitButtonText?: string;
}>();

const emit = defineEmits<{
  (e: 'submit', prompt: Prompt): void;
  (e: 'cancel'): void;
}>();

const promptData = ref<Prompt>({
  id: crypto.randomUUID(),
  name: '',
  description: '',
  system_prompt: '',
  user_prompt: '',
  tags: '',
  created_by: '',
  created_at: new Date(),
  updated_by: '',
  updated_at: new Date(),
});

const tagsInput = ref('');

// Validation rules
const rules = {
  name: { required, minLength: minLength(1) },
  system_prompt: { required, minLength: minLength(1) },
  user_prompt: { required, minLength: minLength(1) }
};

const v$ = useVuelidate(rules, promptData);

onMounted(() => {
  if (props.prompt) {
    promptData.value = { ...props.prompt };
    tagsInput.value = props.prompt.tags || '';
  }
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (isValid) {
    promptData.value.tags = tagsInput.value;
    emit('submit', promptData.value);
  }
};
</script>

<style scoped>
/* Main Form Container */
.prompt-form {
  max-width: 1000px;
  margin: 0 auto;
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

textarea.form-control {
  resize: vertical;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  line-height: 1.4;
}

/* Tags Input */
.tags-input-wrapper small {
  display: block;
  color: #6c757d;
}

/* Preview Section */
.preview-container {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #e9ecef;
}

.preview-section {
  margin-bottom: 1.5rem;
}

.preview-section:last-child {
  margin-bottom: 0;
}

.preview-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #495057;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.preview-content {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.85rem;
  line-height: 1.5;
  color: #495057;
  white-space: pre-wrap;
  word-break: break-word;
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
  .prompt-form {
    padding: 0 1rem;
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
  
  .preview-container {
    padding: 1rem;
  }
  
  .preview-content {
    padding: 0.75rem;
    font-size: 0.8rem;
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
</style>
