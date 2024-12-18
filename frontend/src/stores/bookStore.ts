import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Book } from '@/types';

export const useBookStore = defineStore('book', () => {
  const books = ref<Book[]>([]);
  const searchQuery = ref('');

  const addBook = (book: Book) => {
    books.value.push(book);
  };

  const updateBook = (isbn: string, updatedBook: Book) => {
    const index = books.value.findIndex(book => book.isbn === isbn);
    if (index !== -1) {
      books.value[index] = updatedBook;
    }
  };

  const deleteBook = (isbn: string) => {
    books.value = books.value.filter(book => book.isbn !== isbn);
  };

  const searchBooks = (query: string) => {
    return books.value.filter(book => 
      book.title.toLowerCase().includes(query.toLowerCase()) ||
      book.author.toLowerCase().includes(query.toLowerCase()) ||
      book.isbn.includes(query)
    );
  };

  return {
    books,
    searchQuery,
    addBook,
    updateBook,
    deleteBook,
    searchBooks
  };
});