import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { 
  Bookmark, 
  BookmarkCategory,
  CreateBookmarkRequest,
  UpdateBookmarkRequest, 
  BookmarkListRequest,
  BookmarkListResponse,
  CreateBookmarkCategoryRequest,
  UpdateBookmarkCategoryRequest,
  BookmarkStats
} from '@/types';
import { handleHttpError, showErrorAlert, logError } from '@/utils/errorHandler';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

export const useBookmarkStore = defineStore('bookmark', () => {
  const bookmarks = ref<Bookmark[]>([]);
  const categories = ref<BookmarkCategory[]>([]);
  const allTags = ref<string[]>([]);
  const popularTags = ref<Array<{name: string, count: number}>>([]);
  const stats = ref<BookmarkStats | null>(null);
  
  const totalCount = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const searchQuery = ref('');
  const currentPage = ref(1);
  const pageSize = ref(20);

  // Fetch bookmarks with search and filters
  const fetchBookmarks = async (params: BookmarkListRequest = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      if (params.search || searchQuery.value) {
        queryParams.append('search', params.search || searchQuery.value);
      }
      if (params.category_id) queryParams.append('category_id', params.category_id);
      if (params.tags && params.tags.length > 0) {
        queryParams.append('tags', params.tags.join(','));
      }
      if (params.sort_by) queryParams.append('sort_by', params.sort_by);
      if (params.sort_order) queryParams.append('sort_order', params.sort_order);
      queryParams.append('page', (params.page || currentPage.value).toString());
      queryParams.append('page_size', (params.page_size || pageSize.value).toString());
      
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks?${queryParams}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchBookmarks');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const data: BookmarkListResponse = await response.json();
      
      bookmarks.value = data.bookmarks.map(bookmark => ({
        ...bookmark,
        created_at: new Date(bookmark.created_at),
        updated_at: new Date(bookmark.updated_at)
      }));
      
      totalCount.value = data.total;
      currentPage.value = data.page;
      pageSize.value = data.page_size;
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch bookmarks';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'fetchBookmarks');
    } finally {
      loading.value = false;
    }
  };

  // Search bookmarks (public API)
  const searchBookmarks = async (query: string, page = 1, pageSize = 20) => {
    loading.value = true;
    error.value = null;
    
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('q', query);
      queryParams.append('page', page.toString());
      queryParams.append('page_size', pageSize.toString());
      
      const response = await fetch(getApiUrl(`/api/v1/bookmarks/search?${queryParams}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'searchBookmarks');
        throw new Error(apiError.message);
      }
      
      const data: BookmarkListResponse = await response.json();
      
      return {
        ...data,
        bookmarks: data.bookmarks.map(bookmark => ({
          ...bookmark,
          created_at: new Date(bookmark.created_at),
          updated_at: new Date(bookmark.updated_at)
        }))
      };
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to search bookmarks';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'searchBookmarks');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Get recent bookmarks (public API)
  const fetchRecentBookmarks = async (limit = 10) => {
    try {
      const response = await fetch(getApiUrl(`/api/v1/bookmarks/recent?limit=${limit}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      
      return data.bookmarks.map((bookmark: any) => ({
        ...bookmark,
        created_at: new Date(bookmark.created_at),
        updated_at: new Date(bookmark.updated_at)
      }));
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchRecentBookmarks');
      throw err;
    }
  };

  // Get bookmark by ID
  const fetchBookmark = async (id: string): Promise<Bookmark> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks/${id}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'fetchBookmark');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const bookmark = await response.json();
      
      return {
        ...bookmark,
        created_at: new Date(bookmark.created_at),
        updated_at: new Date(bookmark.updated_at)
      };
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch bookmark';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'fetchBookmark');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Create bookmark
  const createBookmark = async (bookmarkData: CreateBookmarkRequest): Promise<Bookmark> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl('/api/v1/admin/bookmarks'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(bookmarkData),
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'createBookmark');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newBookmark = await response.json();
      
      const bookmark = {
        ...newBookmark,
        created_at: new Date(newBookmark.created_at),
        updated_at: new Date(newBookmark.updated_at)
      };
      
      // Add to local state
      bookmarks.value.unshift(bookmark);
      totalCount.value += 1;
      
      return bookmark;
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create bookmark';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'createBookmark');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Update bookmark
  const updateBookmark = async (id: string, bookmarkData: UpdateBookmarkRequest): Promise<Bookmark> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks/${id}`), {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(bookmarkData),
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'updateBookmark');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedBookmark = await response.json();
      
      const bookmark = {
        ...updatedBookmark,
        created_at: new Date(updatedBookmark.created_at),
        updated_at: new Date(updatedBookmark.updated_at)
      };
      
      // Update local state
      const index = bookmarks.value.findIndex(b => b.id === id);
      if (index !== -1) {
        bookmarks.value[index] = bookmark;
      }
      
      return bookmark;
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update bookmark';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'updateBookmark');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Delete bookmark
  const deleteBookmark = async (id: string): Promise<void> => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks/${id}`), {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        logError(apiError, 'deleteBookmark');
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      // Remove from local state
      const index = bookmarks.value.findIndex(b => b.id === id);
      if (index !== -1) {
        bookmarks.value.splice(index, 1);
        totalCount.value -= 1;
      }
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete bookmark';
      error.value = errorMessage;
      logError({ message: errorMessage, status: 0 }, 'deleteBookmark');
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Fetch categories
  const fetchCategories = async () => {
    try {
      const response = await makeAuthenticatedRequest(getApiUrl('/api/v1/admin/bookmarks/categories'));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      categories.value = data.categories.map((category: any) => ({
        ...category,
        created_at: new Date(category.created_at),
        updated_at: new Date(category.updated_at)
      }));
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchCategories');
      throw err;
    }
  };

  // Create category
  const createCategory = async (categoryData: CreateBookmarkCategoryRequest): Promise<BookmarkCategory> => {
    try {
      const response = await makeAuthenticatedRequest(getApiUrl('/api/v1/admin/bookmarks/categories'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(categoryData),
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const newCategory = await response.json();
      
      const category = {
        ...newCategory,
        created_at: new Date(newCategory.created_at),
        updated_at: new Date(newCategory.updated_at)
      };
      
      categories.value.push(category);
      
      return category;
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'createCategory');
      throw err;
    }
  };

  // Update category
  const updateCategory = async (id: number, categoryData: UpdateBookmarkCategoryRequest): Promise<BookmarkCategory> => {
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks/categories/${id}`), {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(categoryData),
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const updatedCategory = await response.json();
      
      const category = {
        ...updatedCategory,
        created_at: new Date(updatedCategory.created_at),
        updated_at: new Date(updatedCategory.updated_at)
      };
      
      const index = categories.value.findIndex(c => c.id === id);
      if (index !== -1) {
        categories.value[index] = category;
      }
      
      return category;
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'updateCategory');
      throw err;
    }
  };

  // Delete category
  const deleteCategory = async (id: number): Promise<void> => {
    try {
      const response = await makeAuthenticatedRequest(getApiUrl(`/api/v1/admin/bookmarks/categories/${id}`), {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        showErrorAlert(apiError);
        throw new Error(apiError.message);
      }
      
      const index = categories.value.findIndex(c => c.id === id);
      if (index !== -1) {
        categories.value.splice(index, 1);
      }
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'deleteCategory');
      throw err;
    }
  };

  // Fetch all tags
  const fetchAllTags = async () => {
    try {
      const response = await fetch(getApiUrl('/api/v1/bookmarks/tags'));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      allTags.value = data.tags || [];
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchAllTags');
      throw err;
    }
  };

  // Fetch popular tags
  const fetchPopularTags = async (limit = 10) => {
    try {
      const response = await fetch(getApiUrl(`/api/v1/bookmarks/tags/popular?limit=${limit}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data = await response.json();
      popularTags.value = data.tags || [];
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchPopularTags');
      throw err;
    }
  };

  // Fetch bookmarks by category
  const fetchBookmarksByCategory = async (categoryId: string, page = 1, pageSize = 20) => {
    try {
      const response = await fetch(getApiUrl(`/api/v1/bookmarks/category/${categoryId}?page=${page}&page_size=${pageSize}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data: BookmarkListResponse = await response.json();
      
      return {
        ...data,
        bookmarks: data.bookmarks.map(bookmark => ({
          ...bookmark,
          created_at: new Date(bookmark.created_at),
          updated_at: new Date(bookmark.updated_at)
        }))
      };
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchBookmarksByCategory');
      throw err;
    }
  };

  // Fetch bookmarks by tag
  const fetchBookmarksByTag = async (tagName: string, page = 1, pageSize = 20) => {
    try {
      const response = await fetch(getApiUrl(`/api/v1/bookmarks/tag/${tagName}?page=${page}&page_size=${pageSize}`));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      const data: BookmarkListResponse = await response.json();
      
      return {
        ...data,
        bookmarks: data.bookmarks.map(bookmark => ({
          ...bookmark,
          created_at: new Date(bookmark.created_at),
          updated_at: new Date(bookmark.updated_at)
        }))
      };
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchBookmarksByTag');
      throw err;
    }
  };

  // Fetch bookmark stats
  const fetchBookmarkStats = async () => {
    try {
      const response = await fetch(getApiUrl('/api/v1/bookmarks/stats'));
      
      if (!response.ok) {
        const apiError = await handleHttpError(response);
        throw new Error(apiError.message);
      }
      
      stats.value = await response.json();
      
    } catch (err) {
      const apiError = err instanceof Error ? { message: err.message, status: 0 } : { message: 'Unknown error', status: 0 };
      logError(apiError, 'fetchBookmarkStats');
      throw err;
    }
  };

  // Utility functions
  const setSearchQuery = (query: string) => {
    searchQuery.value = query;
  };

  const setCurrentPage = (page: number) => {
    currentPage.value = page;
  };

  const setPageSize = (size: number) => {
    pageSize.value = size;
  };

  const clearError = () => {
    error.value = null;
  };

  const resetState = () => {
    bookmarks.value = [];
    categories.value = [];
    allTags.value = [];
    popularTags.value = [];
    stats.value = null;
    totalCount.value = 0;
    loading.value = false;
    error.value = null;
    searchQuery.value = '';
    currentPage.value = 1;
    pageSize.value = 20;
  };

  return {
    // State
    bookmarks,
    categories,
    allTags,
    popularTags,
    stats,
    totalCount,
    loading,
    error,
    searchQuery,
    currentPage,
    pageSize,
    
    // Actions
    fetchBookmarks,
    searchBookmarks,
    fetchRecentBookmarks,
    fetchBookmark,
    createBookmark,
    updateBookmark,
    deleteBookmark,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    fetchAllTags,
    fetchPopularTags,
    fetchBookmarksByCategory,
    fetchBookmarksByTag,
    fetchBookmarkStats,
    
    // Utilities
    setSearchQuery,
    setCurrentPage,
    setPageSize,
    clearError,
    resetState
  };
});
