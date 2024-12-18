`<template>
  <div class="container mt-4">
    <h1>Book Management</h1>
    
    <div class="mb-4">
      <button class="btn btn-primary" @click="showForm = true" v-if="!showForm">
        Add New Book
      </button>
    </div>

    <div v-if="showForm" class="card mb-4">
      <div class="card-body">
        <h2>{{ editingBook ? 'Edit Book' : 'Add New Book' }}</h2>
        <BookForm
          :book="editingBook"
          :submit-button-text="editingBook ? 'Update Book' : 'Add Book'"
          @submit="handleSubmit"
          @cancel="cancelForm"
        />
      </div>
    </div>

    <BookList
      :books="books"
      :search-query="searchQuery"
      @update:search-query="updateSearchQuery"
      @edit="editBook"
      @delete="deleteBook"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useBookStore } from '@/stores/bookStore';
import BookForm from '@/components/books/BookForm.vue';
import BookList from '@/components/books/BookList.vue';
import type { Book } from '@/types';

const bookStore = useBookStore();
const showForm = ref(false);
const editingBook = ref<Book | undefined>();

const { books, searchQuery } = bookStore;

const updateSearchQuery = (query: string) => {
  bookStore.searchQuery = query;
};

const handleSubmit = (book: Book) => {
  if (editingBook.value) {
    bookStore.updateBook(editingBook.value.isbn, book);
  } else {
    bookStore.addBook(book);
  }
  cancelForm();
};

const editBook = (book: Book) => {
  editingBook.value = book;
  showForm.value = true;
};

const deleteBook = (isbn: string) => {
  if (confirm('Are you sure you want to delete this book?')) {
    bookStore.deleteBook(isbn);
  }
};

const cancelForm = () => {
  showForm.value = false;
  editingBook.value = undefined;
};
</script>`