`<template>
  <form @submit.prevent="handleSubmit" class="secret-form needs-validation" novalidate>
    <!-- Form Content -->
    <div class="form-content">
      <!-- Basic Information Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-shield-lock"></i> Secret Information
        </h6>
        
        <div class="form-grid">
          <div class="form-field">
            <label for="name" class="form-label">
              <i class="bi bi-key text-muted"></i> Secret Name
              <span class="required">*</span>
            </label>
            <input
              type="text"
              class="form-control"
              id="name"
              v-model="secretData.name"
              :class="{ 'is-invalid': v$.name.$error }"
              placeholder="e.g., database-password"
              required
            />
            <div class="invalid-feedback" v-if="v$.name.$error">
              {{ v$.name.$errors[0].$message }}
            </div>
            <small class="text-muted">Unique identifier for this secret</small>
          </div>
        </div>

        <div class="form-grid">
          <div class="form-field">
            <label for="group" class="form-label">
              <i class="bi bi-folder text-muted"></i> Group
              <span class="required">*</span>
            </label>
            <input
              type="text"
              class="form-control"
              id="group"
              v-model="secretData.group"
              :class="{ 'is-invalid': v$.group.$error }"
              placeholder="e.g., database, api-keys, certificates"
              required
            />
            <div class="invalid-feedback" v-if="v$.group.$error">
              {{ v$.group.$errors[0].$message }}
            </div>
            <small class="text-muted">Category or group for organization</small>
          </div>

          <div class="form-field">
            <label for="path" class="form-label">
              <i class="bi bi-signpost text-muted"></i> Path
              <span class="required">*</span>
            </label>
            <input
              type="text"
              class="form-control"
              id="path"
              v-model="secretData.path"
              :class="{ 'is-invalid': v$.path.$error }"
              placeholder="e.g., /prod/database/mysql"
              required
            />
            <div class="invalid-feedback" v-if="v$.path.$error">
              {{ v$.path.$errors[0].$message }}
            </div>
            <small class="text-muted">Hierarchical path for the secret</small>
          </div>
        </div>

        <div class="form-field full-width">
          <label for="desc" class="form-label">
            <i class="bi bi-text-paragraph text-muted"></i> Description
          </label>
          <textarea
            class="form-control"
            id="desc"
            v-model="secretData.desc"
            rows="2"
            placeholder="Describe what this secret is used for..."
          ></textarea>
          <small class="text-muted">Optional description for documentation</small>
        </div>
      </div>

      <!-- Secret Value Section -->
      <div class="form-section">
        <h6 class="section-title">
          <i class="bi bi-lock"></i> Secret Value
        </h6>
        
        <div class="form-field full-width">
          <label for="value" class="form-label">
            <i class="bi bi-asterisk text-muted"></i> Value
            <span class="required">*</span>
          </label>
          <div class="input-group">
            <input
              :type="showSecret ? 'text' : 'password'"
              class="form-control font-monospace"
              id="value"
              v-model="secretData.value"
              :class="{ 'is-invalid': v$.value.$error }"
              placeholder="Enter the secret value"
              required
            />
            <button
              class="btn btn-outline-secondary"
              type="button"
              @click="showSecret = !showSecret"
            >
              <i :class="showSecret ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
            </button>
            <button
              v-if="!isEditMode"
              class="btn btn-outline-primary"
              type="button"
              @click="generateRandomSecret"
              title="Generate random secret"
            >
              <i class="bi bi-shuffle"></i>
            </button>
          </div>
          <div class="invalid-feedback" v-if="v$.value.$error">
            {{ v$.value.$errors[0].$message }}
          </div>
          <small class="text-muted">
            <i class="bi bi-shield-check"></i> This value will be encrypted before storage
          </small>
        </div>

        <!-- Secret strength indicator -->
        <div v-if="secretData.value && !isEditMode" class="secret-strength">
          <div class="strength-label">
            <span>Secret Strength:</span>
            <span :class="`strength-${secretStrength.level}`">
              {{ secretStrength.label }}
            </span>
          </div>
          <div class="strength-bar">
            <div 
              class="strength-fill"
              :class="`strength-${secretStrength.level}`"
              :style="{ width: `${secretStrength.score}%` }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <button type="submit" class="btn btn-primary btn-lg">
        <i :class="isEditMode ? 'bi bi-save' : 'bi bi-plus-lg'" class="me-2"></i>
        {{ isEditMode ? 'Update Secret' : 'Create Secret' }}
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
import { required } from '@vuelidate/validators';
import type { Secret, CreateSecretRequest } from '@/types';

const props = defineProps<{
  secret?: Secret;
  isEditMode?: boolean;
}>();

const emit = defineEmits<{
  (e: 'submit', secret: CreateSecretRequest): void;
  (e: 'cancel'): void;
}>();

const secretData = ref<CreateSecretRequest>({
  name: '',
  group: '',
  desc: '',
  path: '',
  value: ''
});

const showSecret = ref(false);

const rules = {
  name: { required },
  group: { required },
  path: { required },
  value: { required }
};

const v$ = useVuelidate(rules, secretData);

// Calculate secret strength
const secretStrength = computed(() => {
  const value = secretData.value.value;
  if (!value) return { level: 'none', label: 'None', score: 0 };
  
  let score = 0;
  
  // Length score
  if (value.length >= 8) score += 20;
  if (value.length >= 12) score += 20;
  if (value.length >= 16) score += 10;
  
  // Character variety score
  if (/[a-z]/.test(value)) score += 10;
  if (/[A-Z]/.test(value)) score += 10;
  if (/[0-9]/.test(value)) score += 10;
  if (/[^a-zA-Z0-9]/.test(value)) score += 20;
  
  // Determine level
  if (score >= 80) return { level: 'strong', label: 'Strong', score };
  if (score >= 60) return { level: 'good', label: 'Good', score };
  if (score >= 40) return { level: 'fair', label: 'Fair', score };
  return { level: 'weak', label: 'Weak', score };
});

// Generate random secret
const generateRandomSecret = () => {
  const length = 24;
  const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?';
  let result = '';
  
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);
  
  for (let i = 0; i < length; i++) {
    result += charset[array[i] % charset.length];
  }
  
  secretData.value.value = result;
};

onMounted(() => {
  if (props.secret && props.isEditMode) {
    // In edit mode, populate form with existing data (except the value)
    secretData.value = {
      name: props.secret.name,
      group: props.secret.group,
      desc: props.secret.desc,
      path: props.secret.path,
      value: '' // User must re-enter the value for updates
    };
  }
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (isValid) {
    emit('submit', secretData.value);
  }
};
</script>

<style scoped>
/* Form Container */
.secret-form {
  max-width: 800px;
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
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
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
  min-height: 80px;
}

.font-monospace {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

/* Input Groups */
.input-group .btn {
  border: 2px solid #e9ecef;
  border-left: none;
}

.input-group .btn:hover {
  background-color: #f8f9fa;
}

/* Secret Strength Indicator */
.secret-strength {
  margin-top: 1rem;
}

.strength-label {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.strength-bar {
  height: 8px;
  background-color: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
}

.strength-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.strength-weak { color: #dc3545; background-color: #dc3545; }
.strength-fair { color: #fd7e14; background-color: #fd7e14; }
.strength-good { color: #ffc107; background-color: #ffc107; }
.strength-strong { color: #28a745; background-color: #28a745; }

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

/* Help Text */
small.text-muted {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.875rem;
}

/* Validation */
.invalid-feedback {
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.is-invalid {
  border-color: #dc3545 !important;
}

/* Responsive Design */
@media (max-width: 768px) {
  .secret-form {
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
}
</style>`
