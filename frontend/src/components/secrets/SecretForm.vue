`<template>
  <form @submit.prevent="handleSubmit" class="secret-form needs-validation" novalidate>
    <!-- Edit Mode Warning -->
    <div v-if="isEditMode" class="alert alert-info mb-4" role="alert">
      <div class="d-flex align-items-center">
        <i class="bi bi-info-circle me-2"></i>
        <div>
          <strong>Security Notice:</strong> For security reasons, you must enter the current secret value for verification and provide the new value you want to store.
        </div>
      </div>
    </div>

    <!-- Form Validation Errors -->
    <div v-if="formErrors.length > 0" class="alert alert-danger mb-4" role="alert">
      <div class="d-flex align-items-start">
        <i class="bi bi-exclamation-triangle me-2 mt-1"></i>
        <div>
          <strong>Please fix the following errors:</strong>
          <ul class="mb-0 mt-2">
            <li v-for="error in formErrors" :key="error">{{ error }}</li>
          </ul>
        </div>
      </div>
    </div>

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
              placeholder="e.g., password, credentials, api-keys, certificates"
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
              placeholder="e.g., dev/database/mysql"
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
            <i class="bi bi-asterisk text-muted"></i> 
            {{ isEditMode ? 'New Secret Value' : 'Value' }}
            <span class="required">*</span>
            <span v-if="isEditMode" class="badge bg-success ms-2">
              <i class="bi bi-plus-circle"></i> New Value
            </span>
          </label>
          <div class="input-group">
            <input
              :type="showSecret ? 'text' : 'password'"
              class="form-control font-monospace"
              id="value"
              v-model="secretData.value"
              :class="{ 'is-invalid': v$.value.$error }"
              :placeholder="isEditMode ? 'Enter the new secret value' : 'Enter the secret value'"
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
            <i class="bi bi-shield-check"></i> 
            {{ isEditMode ? 'Enter the new value you want to store (required)' : 'This value will be encrypted before storage' }}
          </small>
        </div>

        <!-- Current Value Verification (Edit Mode Only) -->
        <div v-if="isEditMode" class="form-field full-width">
          <label for="currentValue" class="form-label">
            <i class="bi bi-shield-check text-muted"></i> Current Secret Value (for verification)
            <span class="required">*</span>
            <span class="badge bg-warning ms-2">
              <i class="bi bi-exclamation-triangle"></i> Verification Required
            </span>
          </label>
          <div class="input-group">
            <input
              :type="showCurrentSecret ? 'text' : 'password'"
              class="form-control font-monospace"
              id="currentValue"
              v-model="currentValue"
              :class="{ 'is-invalid': v$.currentValue?.$error }"
              placeholder="Enter the current secret value to verify your identity"
              required
            />
            <button
              class="btn btn-outline-secondary"
              type="button"
              @click="showCurrentSecret = !showCurrentSecret"
            >
              <i :class="showCurrentSecret ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
            </button>
          </div>
          <div class="invalid-feedback" v-if="v$.currentValue?.$error">
            {{ v$.currentValue?.$errors[0].$message }}
          </div>
          <small class="text-muted">
            <i class="bi bi-info-circle"></i> 
            For security, enter the current secret value to verify you have access before updating
          </small>
        </div>

        <!-- Custom KEK Section -->
        <div class="form-section">
          <h6 class="section-title">
            <i class="bi bi-key-fill"></i> Encryption Key (Optional)
          </h6>
          
          <div class="form-field">
            <label for="kek" class="form-label">
              <i class="bi bi-shield-lock text-muted"></i> Custom KEK Password
            </label>
            <div class="input-group">
              <input
                :type="showKEK ? 'text' : 'password'"
                class="form-control font-monospace"
                id="kek"
                v-model="secretData.kek"
                :class="{ 'is-invalid': v$.kek.$error }"
                placeholder="Enter your custom password/phrase (optional)"
                maxlength="256"
              />
              <button
                class="btn btn-outline-secondary"
                type="button"
                @click="showKEK = !showKEK"
              >
                <i :class="showKEK ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
              </button>
              <button
                class="btn btn-outline-primary"
                type="button"
                @click="generateRandomKEK"
                title="Generate random KEK"
              >
                <i class="bi bi-shuffle"></i>
              </button>
            </div>
            <div class="invalid-feedback" v-if="v$.kek.$error">
              {{ v$.kek.$errors[0].$message }}
            </div>
            <small class="text-muted">
              <i class="bi bi-info-circle"></i> 
              Leave empty to use system default KEK. You can enter any password or phrase - it will be securely hashed.
            </small>
          </div>
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
import { ref, computed, onMounted, watch } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { required } from '@vuelidate/validators';
import type { Secret, CreateSecretRequest, UpdateSecretRequest } from '@/types';

const props = defineProps<{
  secret?: Secret;
  isEditMode?: boolean;
}>();

const emit = defineEmits<{
  (e: 'submit', secret: CreateSecretRequest | UpdateSecretRequest): void;
  (e: 'cancel'): void;
}>();

const secretData = ref<CreateSecretRequest>({
  name: '',
  group: '',
  desc: '',
  path: '',
  value: '', // The secret value (create mode) or NEW secret value (edit mode)
  kek: ''
});

// Edit mode verification data
const currentValue = ref(''); // Current secret value for verification (edit mode only)

const showSecret = ref(false);
const showCurrentSecret = ref(false);
const showKEK = ref(false);
const formErrors = ref<string[]>([]);

const rules = computed(() => {
  const baseRules = {
    name: { required },
    group: { required },
    path: { required },
    value: { required },
    kek: {
      kekValidator: {
        $validator: (value: string | undefined) => {
          console.log('🔍 KEK validator called with:', value, 'type:', typeof value);
          // KEK is optional - empty string is allowed
          if (!value || value === '') {
            console.log('🔍 KEK validator: empty/undefined, returning true');
            return true;
          }
          // If provided, it must not be just whitespace
          const result = value.trim().length > 0;
          console.log('🔍 KEK validator: non-empty, trim length:', value.trim().length, 'result:', result);
          return result;
        },
        $message: 'KEK cannot be empty if provided'
      }
    }
  };

  // Add current value validation in edit mode
  if (props.isEditMode) {
    return {
      ...baseRules,
      currentValue: { required }
    };
  }

  return baseRules;
});

// Create the data object for Vuelidate
const vuelidateData = computed(() => ({
  ...secretData.value,
  currentValue: currentValue.value
}));

const v$ = useVuelidate(rules, vuelidateData);

// Calculate secret strength
const secretStrength = computed(() => {
  const value = secretData.value.value; // secretData.value is the reactive object, .value is the secret value property
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
  
  secretData.value.value = result; // secretData.value is the reactive object, .value is the secret value property
};

const generateRandomKEK = () => {
  // Generate a user-friendly passphrase using words and numbers
  const words = ['secure', 'encrypt', 'protect', 'shield', 'guard', 'safe', 'vault', 'lock', 'key', 'cipher'];
  const adjectives = ['strong', 'powerful', 'mighty', 'robust', 'solid', 'tough', 'firm', 'stable'];
  
  const randomWord = words[Math.floor(Math.random() * words.length)];
  const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)];
  const randomNumber = Math.floor(Math.random() * 9999) + 1000;
  
  secretData.value.kek = `${randomAdjective}-${randomWord}-${randomNumber}`; // secretData.value is the reactive object, .kek is the KEK property
};

onMounted(() => {
  if (props.secret && props.isEditMode) {
    // In edit mode, populate form with existing data (except the value and KEK)
    secretData.value = {
      name: props.secret.name,
      group: props.secret.group,
      desc: props.secret.desc,
      path: props.secret.path,
      value: '', // User must enter the NEW value for updates
      kek: '' // User must re-enter the KEK for updates
    };
    // Clear current value field
    currentValue.value = '';
  }
  
  // Clear errors when user starts typing
  watch([secretData, currentValue], () => {
    if (formErrors.value.length > 0) {
      formErrors.value = [];
    }
  }, { deep: true });
});

const handleSubmit = async () => {
  console.log('=== FORM SUBMISSION DEBUG ===');
  console.log('handleSubmit called');
  console.log('isEditMode:', props.isEditMode);
  console.log('secretData:', secretData.value, currentValue.value);
  console.log('currentValue length:', currentValue.value.length);
  console.log('KEK value:', `"${secretData.value.kek}"`);
  console.log('KEK length:', secretData.value.kek.length);
  console.log('KEK is empty?', !secretData.value.kek);
  console.log('KEK === ""?', secretData.value.kek === '');
  console.log('KEK trim length:', secretData.value.kek.trim().length);
  console.log('KEK validation result:', !secretData.value.kek || secretData.value.kek === '' ? true : secretData.value.kek.trim().length > 0);
  console.log('v$.kek.$error:', v$.value.kek.$error);
  console.log('v$.kek.$error:', v$.value.kek.$error);
  
  // Test the validator function directly
  const kekValidator = (value: string | undefined) => {
    console.log('Direct validator test - input:', value);
    if (!value || value === '') {
      console.log('Direct validator: empty, returning true');
      return true;
    }
    const result = value.trim().length > 0;
    console.log('Direct validator: non-empty, trim length:', value.trim().length, 'result:', result);
    return result;
  };
  console.log('Direct KEK validator test:', kekValidator(secretData.value.kek));

  // Clear previous errors
  formErrors.value = [];
  
  // Debug the data being passed to Vuelidate
  const vuelidateData = {
    ...secretData.value,
    currentValue: currentValue.value
  };
  console.log('Data passed to Vuelidate:', vuelidateData);
  console.log('KEK in Vuelidate data:', vuelidateData.kek);
  console.log('Vuelidate rules:', rules.value);

  const isValid = await v$.value.$validate();
  console.log('Form validation result:', isValid);
  console.log('All validation errors:', v$.value.$errors);
  
  // Debug each field's validation state
  console.log('Field validation states:');
  console.log('  name.$error:', v$.value.name.$error);
  console.log('  group.$error:', v$.value.group.$error);
  console.log('  path.$error:', v$.value.path.$error);
  console.log('  value.$error:', v$.value.value.$error);
  console.log('  kek.$error:', v$.value.kek?.$error);
  if (props.isEditMode) {
    console.log('  currentValue.$error:', v$.value.currentValue?.$error);
  }
  
  // Show actual error details
  v$.value.$errors.forEach((error, index) => {
    console.log(`Error ${index}:`, error);
  });
  
  if (isValid) {
    console.log('✅ FORM IS VALID - Proceeding with submission');
    console.log('Emitting submit event with data:', secretData.value);
    console.log('Current value for verification:', currentValue.value);
    formErrors.value = []; // Clear any previous errors
    
    if (props.isEditMode) {
      // For edit mode, create UpdateSecretRequest
      const updateRequest: UpdateSecretRequest = {
        name: secretData.value.name,
        group: secretData.value.group,
        desc: secretData.value.desc,
        path: secretData.value.path,
        value: secretData.value.value,        // NEW secret value
        current_value: currentValue.value,    // Current secret value for verification
        kek: secretData.value.kek
      };
      console.log('Sending UpdateSecretRequest:', updateRequest);
      emit('submit', updateRequest);
    } else {
      // For create mode, use CreateSecretRequest as-is
      console.log('Sending CreateSecretRequest:', secretData.value);
      emit('submit', secretData.value);
    }
  } else {
    console.log('❌ FORM VALIDATION FAILED');
    console.log('Form validation failed, not submitting');
    
    // Collect user-friendly error messages
    if (v$.value.name.$error) {
      formErrors.value.push('Name is required');
      console.log('Name validation errors:', v$.value.name.$errors);
    }
    if (v$.value.group.$error) {
      formErrors.value.push('Group is required');
      console.log('Group validation errors:', v$.value.group.$errors);
    }
    if (v$.value.path.$error) {
      formErrors.value.push('Path is required');
      console.log('Path validation errors:', v$.value.path.$errors);
    }
    if (v$.value.value.$error) {
      formErrors.value.push(props.isEditMode ? 'New secret value is required' : 'Secret value is required');
      console.log('Value validation errors:', v$.value.value.$errors);
    }
    if (props.isEditMode && v$.value.currentValue?.$error) {
      formErrors.value.push('Current secret value is required for verification');
      console.log('Current value validation errors:', v$.value.currentValue?.$errors);
    }
    if (v$.value.kek?.$error) {
      formErrors.value.push('KEK cannot be empty if provided');
      console.log('KEK validation errors:', v$.value.kek.$errors);
    }
    
    console.log('Final formErrors array:', formErrors.value);
    
    // Scroll to top to show errors
    document.querySelector('.secret-form')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
  console.log('=== END FORM SUBMISSION DEBUG ===');
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
