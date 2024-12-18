`<template>
  <form @submit.prevent="handleSubmit" class="needs-validation" novalidate>
    <div class="mb-3">
      <label for="isbn" class="form-label">ISBN</label>
      <input
        type="text"
        class="form-control"
        id="isbn"
        v-model="bookData.isbn"
        :class="{ 'is-invalid': v$.isbn.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.isbn.$error">
        {{ v$.isbn.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="title" class="form-label">Title</label>
      <input
        type="text"
        class="form-control"
        id="title"
        v-model="bookData.title"
        :class="{ 'is-invalid': v$.title.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.title.$error">
        {{ v$.title.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="author" class="form-label">Author</label>
      <input
        type="text"
        class="form-control"
        id="author"
        v-model="bookData.author"
        :class="{ 'is-invalid': v$.author.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.author.$error">
        {{ v$.author.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="price" class="form-label">Price</label>
      <input
        type="number"
        class="form-control"
        id="price"
        v-model.number="bookData.price"
        :class="{ 'is-invalid': v$.price.$error }"
        required
      />
      <div class="invalid-feedback" v-if="v$.price.$error">
        {{ v$.price.$errors[0].$message }}
      </div>
    </div>

    <div class="mb-3">
      <label for="borrowTime" class="form-label">Borrow Time</label>
      <input
        type="datetime-local"
        class="form-control"
        id="borrowTime"
        v-model="bookData.borrowTime"
      />
    </div>

    <div class="mb-3">
      <label for="returnTime" class="form-label">Return Time</label>
      <input
        type="datetime-local"
        class="form-control"
        id="returnTime"
        v-model="bookData.returnTime"
      />
    </div>

    <button type="submit" class="btn btn-primary">{{ submitButtonText }}</button>
    <button type="button" class="btn btn-secondary ms-2" @click="$emit('cancel')">
      Cancel
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import type { Book } from '@/types';
import { useBookValidation } from './useBookValidation';

const props = defineProps<{
  book?: Book;
  submitButtonText?: string;
}>();

const emit = defineEmits<{
  (e: 'submit', book: Book): void;
  (e: 'cancel'): void;
}>();

const bookData = ref<Book>({
  isbn: '',
  title: '',
  author: '',
  price: 0,
  borrowTime: undefined,
  returnTime: undefined,
});

const { rules } = useBookValidation();
const v$ = useVuelidate(rules, bookData);

onMounted(() => {
  if (props.book) {
    bookData.value = { ...props.book };
  }
});

const handleSubmit = async () => {
  const isValid = await v$.value.$validate();
  if (isValid) {
    emit('submit', bookData.value);
  }
};
</script>`