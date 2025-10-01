import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Book, CreateBookRequest, UpdateBookRequest, BookListResponse } from '@/types';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';

export const useBookStore = defineStore('book', () => {
  const books = ref<Book[]>([]);
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);

  // Fetch books with search and filters
  const fetchBooks = async (params: {
    realm_id?: string;
    q?: string;
    author?: string;
    tags?: string;
    page?: number;
    page_size?: number;
  } = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      if (params.realm_id) queryParams.append('realm_id', params.realm_id);
      if (params.q || searchQuery.value) queryParams.append('q', params.q || searchQuery.value);
      if (params.author) queryParams.append('author', params.author);
      if (params.tags) queryParams.append('tags', params.tags);
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('page_size', (params.page_size || pageSize.value).toString());
      
      const response = await makeAuthenticatedRequest(`/api/v1/books?${queryParams}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchBooks');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      // Convert backend book format to frontend format
      books.value = (data.books || []).map((book: any) => ({
        ...book,
        borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
        return_time: book.return_time ? new Date(book.return_time) : undefined,
        deadline: book.deadline ? new Date(book.deadline) : undefined,
        created_at: new Date(book.created_at),
        updated_at: new Date(book.updated_at),
        tags: book.tags || []
      }));
      
      totalCount.value = data.total || 0;
      currentPage.value = data.page || 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get a single book
  const getBook = async (id: string): Promise<Book> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/books/${id}`);
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'getBook');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const book = await response.json();
      
      // Convert backend format to frontend format
      return {
        ...book,
        borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
        return_time: book.return_time ? new Date(book.return_time) : undefined,
        deadline: book.deadline ? new Date(book.deadline) : undefined,
        created_at: new Date(book.created_at),
        updated_at: new Date(book.updated_at),
        tags: book.tags || []
      };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create a new book
  const addBook = async (book: Book): Promise<Book> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: CreateBookRequest = {
        isbn: book.isbn,
        name: book.name,
        author: book.author,
        description: book.description,
        price: book.price,
        tags: book.tags.join(','),
        deadline: book.deadline
      };
      
      const response = await makeAuthenticatedRequest('/api/v1/books', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'addBook');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newBook = await response.json();
      
      // Convert and add to local store
      const convertedBook = {
        ...newBook,
        borrow_time: newBook.borrow_time ? new Date(newBook.borrow_time) : undefined,
        return_time: newBook.return_time ? new Date(newBook.return_time) : undefined,
        deadline: newBook.deadline ? new Date(newBook.deadline) : undefined,
        created_at: new Date(newBook.created_at),
        updated_at: new Date(newBook.updated_at),
        tags: newBook.tags || []
      };
      
      books.value.push(convertedBook);
      totalCount.value += 1;
      return convertedBook;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update an existing book
  const updateBook = async (id: string, updates: Partial<Book>): Promise<Book> => {
    loading.value = true;
    error.value = null;
    
    try {
      const requestData: UpdateBookRequest = {};
      
      if (updates.isbn !== undefined) requestData.isbn = updates.isbn;
      if (updates.name !== undefined) requestData.name = updates.name;
      if (updates.author !== undefined) requestData.author = updates.author;
      if (updates.description !== undefined) requestData.description = updates.description;
      if (updates.price !== undefined) requestData.price = updates.price;
      if (updates.tags !== undefined) requestData.tags = updates.tags.join(',');
      if (updates.deadline !== undefined) requestData.deadline = updates.deadline;
      
      const response = await makeAuthenticatedRequest(`/api/v1/books/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateBook');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedBook = await response.json();
      
      // Convert and update in local store
      const convertedBook = {
        ...updatedBook,
        borrow_time: updatedBook.borrow_time ? new Date(updatedBook.borrow_time) : undefined,
        return_time: updatedBook.return_time ? new Date(updatedBook.return_time) : undefined,
        deadline: updatedBook.deadline ? new Date(updatedBook.deadline) : undefined,
        created_at: new Date(updatedBook.created_at),
        updated_at: new Date(updatedBook.updated_at),
        tags: updatedBook.tags || []
      };
      
      const index = books.value.findIndex(b => b.id === id);
      if (index !== -1) {
        books.value[index] = convertedBook;
      }
      return convertedBook;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete a book
  const deleteBook = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(`/api/v1/books/${id}`, {
        method: 'DELETE'
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteBook');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      books.value = books.value.filter(b => b.id !== id);
      totalCount.value -= 1;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Borrow a book
  const borrowBook = async (id: string, deadline: Date): Promise<Book> => {
    const response = await makeAuthenticatedRequest(`/api/v1/books/${id}/borrow`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        deadline: deadline.toISOString()
      })
    });
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'borrowBook');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const updatedBook = await response.json();
    await fetchBooks(); // Refresh the list
    return updatedBook;
  };

  // Return a book
  const returnBook = async (id: string): Promise<Book> => {
    const response = await makeAuthenticatedRequest(`/api/v1/books/${id}/return`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'returnBook');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const updatedBook = await response.json();
    await fetchBooks(); // Refresh the list
    return updatedBook;
  };

  // Get borrowed books
  const getBorrowedBooks = async (): Promise<Book[]> => {
    const response = await makeAuthenticatedRequest('/api/v1/books/borrowed');
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getBorrowedBooks');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return (data.books || []).map((book: any) => ({
      ...book,
      borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
      return_time: book.return_time ? new Date(book.return_time) : undefined,
      deadline: book.deadline ? new Date(book.deadline) : undefined,
      created_at: new Date(book.created_at),
      updated_at: new Date(book.updated_at),
      tags: book.tags || []
    }));
  };

  // Get overdue books
  const getOverdueBooks = async (): Promise<Book[]> => {
    const response = await makeAuthenticatedRequest('/api/v1/books/overdue');
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getOverdueBooks');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return (data.books || []).map((book: any) => ({
      ...book,
      borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
      return_time: book.return_time ? new Date(book.return_time) : undefined,
      deadline: book.deadline ? new Date(book.deadline) : undefined,
      created_at: new Date(book.created_at),
      updated_at: new Date(book.updated_at),
      tags: book.tags || []
    }));
  };

  // Get available books
  const getAvailableBooks = async (): Promise<Book[]> => {
    const response = await makeAuthenticatedRequest('/api/v1/books/available');
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getAvailableBooks');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return (data.books || []).map((book: any) => ({
      ...book,
      borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
      return_time: book.return_time ? new Date(book.return_time) : undefined,
      deadline: book.deadline ? new Date(book.deadline) : undefined,
      created_at: new Date(book.created_at),
      updated_at: new Date(book.updated_at),
      tags: book.tags || []
    }));
  };

  // Get books by tag
  const getBooksByTag = async (tagName: string): Promise<Book[]> => {
    const response = await makeAuthenticatedRequest(`/api/v1/books/tag/${encodeURIComponent(tagName)}`);
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getBooksByTag');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return (data.books || []).map((book: any) => ({
      ...book,
      borrow_time: book.borrow_time ? new Date(book.borrow_time) : undefined,
      return_time: book.return_time ? new Date(book.return_time) : undefined,
      deadline: book.deadline ? new Date(book.deadline) : undefined,
      created_at: new Date(book.created_at),
      updated_at: new Date(book.updated_at),
      tags: book.tags || []
    }));
  };

  // Get all tags
  const getAllTags = async (): Promise<string[]> => {
    const response = await makeAuthenticatedRequest('/api/v1/books/tags');
    
    if (!response.ok) {
      const apiError = await handleHttpError(response);
      logError(apiError, 'getAllTags');
      showErrorAlert(apiError);
      throw new Error(apiError.message);
    }
    
    const data = await response.json();
    return data.tags || [];
  };

  // Clear store
  const clearBooks = () => {
    books.value = [];
    totalCount.value = 0;
    error.value = null;
    searchQuery.value = '';
    currentPage.value = 1;
  };

  return {
    books,
    totalCount,
    loading,
    error,
    searchQuery,
    currentPage,
    pageSize,
    fetchBooks,
    getBook,
    addBook,
    updateBook,
    deleteBook,
    borrowBook,
    returnBook,
    getBorrowedBooks,
    getOverdueBooks,
    getAvailableBooks,
    getBooksByTag,
    getAllTags,
    clearBooks
  };
});