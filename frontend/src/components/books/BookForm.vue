<template>
  <form @submit.prevent="handleSubmit" class="book-form needs-validation" novalidate>
    <!-- Book Details Section -->
    <div class="form-section">
      <h6 class="section-title">
        <i class="bi bi-book"></i> Book Details
      </h6>
      
      <div class="form-grid">
        <div class="form-field">
          <label for="isbn" class="form-label">
            <i class="bi bi-upc-scan text-muted"></i> ISBN
            <span class="required">*</span>
          </label>
          <input
            type="text"
            class="form-control"
            id="isbn"
            v-model="bookData.isbn"
            :class="{ 'is-invalid': v$.isbn.$error }"
            placeholder="978-0-123456-78-9"
            required
          />
          <div class="invalid-feedback" v-if="v$.isbn.$error">
            {{ v$.isbn.$errors[0].$message }}
          </div>
        </div>
        
        <div class="form-field">
          <label for="name" class="form-label">
            <i class="bi bi-card-heading text-muted"></i> Book Name
            <span class="required">*</span>
          </label>
          <input
            type="text"
            class="form-control"
            id="name"
            v-model="bookData.name"
            :class="{ 'is-invalid': v$.name.$error }"
            placeholder="Enter book title"
            required
          />
          <div class="invalid-feedback" v-if="v$.name.$error">
            {{ v$.name.$errors[0].$message }}
          </div>
        </div>
        
        <div class="form-field">
          <label for="author" class="form-label">
            <i class="bi bi-person text-muted"></i> Author
          </label>
          <input
            type="text"
            class="form-control"
            id="author"
            v-model="bookData.author"
            placeholder="Enter author name"
          />
        </div>
        
        <div class="form-field">
          <label for="price" class="form-label">
            <i class="bi bi-currency-dollar text-muted"></i> Price
          </label>
          <div class="input-group">
            <span class="input-group-text">$</span>
            <input
              type="number"
              class="form-control"
              id="price"
              v-model.number="bookData.price"
              step="0.01"
              min="0"
              placeholder="0.00"
            />
          </div>
        </div>
      </div>
      
      <div class="form-field full-width">
        <label for="description" class="form-label">
          <i class="bi bi-text-paragraph text-muted"></i> Description
        </label>
        <textarea
          class="form-control"
          id="description"
          v-model="bookData.description"
          rows="3"
          placeholder="Enter book description..."
        ></textarea>
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
            <i class="bi bi-info-circle"></i> Use tags to organize and find books easily
          </small>
        </div>
      </div>
    </div>

    <!-- Deadline Section -->
    <div class="form-section">
      <h6 class="section-title">
        <i class="bi bi-calendar-event"></i> Borrowing Information
      </h6>
      
      <div class="form-field">
        <label for="deadline" class="form-label">
          <i class="bi bi-calendar-check text-muted"></i> Default Return Deadline
        </label>
        <input
          type="datetime-local"
          class="form-control"
          id="deadline"
          v-model="deadlineFormatted"
        />
        <small class="text-muted">Default deadline when someone borrows this book</small>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <button type="submit" class="btn btn-primary btn-lg">
        <i class="bi bi-check-lg me-2"></i>
        {{ submitButtonText || 'Create Book' }}
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
import { required, minLength } from '@vuelidate/validators';
import type { Book } from '@/types';

const props = defineProps<{
  book?: Book;
  submitButtonText?: string;
}>();

const emit = defineEmits<{
  (e: 'submit', book: Book): void;
  (e: 'cancel'): void;
}>();

const bookData = ref<Book>({
  id: crypto.randomUUID(),
  realm_id: '',
  isbn: '',
  name: '',
  author: '',
  description: '',
  price: 0,
  tags: [],
  created_by: '',
  created_at: new Date(),
  updated_by: '',
  updated_time: new Date(),
});

const tagsInput = ref('');

// Validation rules
const rules = {
  isbn: { required, minLength: minLength(10) },
  name: { required, minLength: minLength(1) }
};

const v$ = useVuelidate(rules, bookData);

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

// Computed property for deadline input
const deadlineFormatted = computed({
  get: () => formatDateForInput(bookData.value.deadline),
  set: (value: string) => {
    if (value) bookData.value.deadline = new Date(value);
  }
});

onMounted(() => {
  if (props.book) {
    bookData.value = { ...props.book };
    tagsInput.value = props.book.tags.join(', ');
  }
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (isValid) {
    bookData.value.tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);
    emit('submit', bookData.value);
  }
};
</script>

<style scoped>
/* Main Form Container */
.book-form {
  max-width: 900px;
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
  .book-form {
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