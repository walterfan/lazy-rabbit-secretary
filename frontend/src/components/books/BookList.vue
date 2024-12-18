`<template>
  <div class="book-list">
    <SearchInput
      :model-value="searchQuery"
      @update:model-value="$emit('update:searchQuery', $event)"
      placeholder="Search books..."
    />

    <div class="table-responsive">
      <table class="table table-striped">
        <thead>
          <tr>
            <th>ISBN</th>
            <th>Title</th>
            <th>Author</th>
            <th>Price</th>
            <th>Borrow Time</th>
            <th>Return Time</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="book in filteredBooks" :key="book.isbn">
            <td>{{ book.isbn }}</td>
            <td>{{ book.title }}</td>
            <td>{{ book.author }}</td>
            <td>{{ formatPrice(book.price) }}</td>
            <td>{{ formatDate(book.borrowTime) }}</td>
            <td>{{ formatDate(book.returnTime) }}</td>
            <td>
              <button
                class="btn btn-sm btn-primary me-2"
                @click="$emit('edit', book)"
              >
                Edit
              </button>
              <button
                class="btn btn-sm btn-danger"
                @click="handleDelete(book.isbn)"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Book } from '@/types';
import { formatDate, formatPrice } from '@/utils/dateUtils';
import { useConfirmDialog } from '@/components/common/ConfirmDialog';
import SearchInput from '@/components/common/SearchInput.vue';

const props = defineProps<{
  books: Book[];
  searchQuery: string;
}>();

const emit = defineEmits<{
  (e: 'edit', book: Book): void;
  (e: 'delete', isbn: string): void;
  (e: 'update:searchQuery', value: string): void;
}>();

const { confirm } = useConfirmDialog();

const filteredBooks = computed(() => {
  if (!props.searchQuery) return props.books;
  
  return props.books.filter(book => 
    book.title.toLowerCase().includes(props.searchQuery.toLowerCase()) ||
    book.author.toLowerCase().includes(props.searchQuery.toLowerCase()) ||
    book.isbn.includes(props.searchQuery)
  );
});

const handleDelete = async (isbn: string) => {
  if (await confirm('Are you sure you want to delete this book?')) {
    emit('delete', isbn);
  }
};
</script>`