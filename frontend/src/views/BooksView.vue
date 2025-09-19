<template>
  <div class="books-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-book"></i>
          Book Management
        </h1>
        <p class="page-description">
          Manage your book collection with borrowing tracking, due dates, and organization features
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm">
          <i class="bi bi-plus-lg me-2"></i>
          Add Book
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Create/Edit Form -->
      <div v-if="showForm" class="form-section">
        <div class="section-header">
          <h2>{{ editingBook ? 'Edit Book' : 'Add New Book' }}</h2>
          <button 
            class="btn btn-sm btn-light"
            @click="closeForm"
          >
            <i class="bi bi-x-lg"></i>
          </button>
        </div>
        <BookForm
          :book="editingBook"
          :submit-button-text="editingBook ? 'Update Book' : 'Add Book'"
          @submit="handleFormSubmit"
          @cancel="closeForm"
        />
      </div>

      <!-- Book List -->
      <div v-else>
        <BookList
          :books="bookStore.books"
          :search-query="searchQuery"
          :current-page="currentPage"
          :page-size="pageSize"
          :total-count="bookStore.totalCount"
          :loading="bookStore.loading"
          @view="handleView"
          @edit="handleEdit"
          @delete="handleDelete"
          @borrow="handleBorrow"
          @return="handleReturn"
          @update:search-query="searchQuery = $event"
          @update:page="currentPage = $event"
        />
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="bookStore.loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Book Details Modal -->
    <div 
      v-if="viewingBook"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="viewingBook = null"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-book me-2"></i>
              Book Details
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="viewingBook = null"
            ></button>
          </div>
          <div class="modal-body">
            <div class="book-details">
              <div class="detail-row">
                <span class="detail-label">Title:</span>
                <span class="detail-value fw-bold">{{ viewingBook.name }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Author:</span>
                <span class="detail-value">{{ viewingBook.author || 'Unknown Author' }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">ISBN:</span>
                <span class="detail-value">{{ viewingBook.isbn }}</span>
              </div>
              <div class="detail-row" v-if="viewingBook.description">
                <span class="detail-label">Description:</span>
                <span class="detail-value">{{ viewingBook.description }}</span>
              </div>
              <div class="detail-row" v-if="viewingBook.price > 0">
                <span class="detail-label">Price:</span>
                <span class="detail-value">${{ viewingBook.price.toFixed(2) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Status:</span>
                <span class="detail-value">
                  <span :class="getStatusBadgeClass(viewingBook)">
                    {{ getBookStatus(viewingBook) }}
                  </span>
                </span>
              </div>
              <div class="detail-row" v-if="isBorrowed(viewingBook)">
                <span class="detail-label">Borrowed:</span>
                <span class="detail-value">{{ formatDate(viewingBook.borrow_time!) }}</span>
              </div>
              <div class="detail-row" v-if="isBorrowed(viewingBook) && viewingBook.deadline">
                <span class="detail-label">Due Date:</span>
                <span class="detail-value">{{ formatDate(viewingBook.deadline!) }}</span>
              </div>
              <div class="detail-row" v-if="viewingBook.tags && viewingBook.tags.length > 0">
                <span class="detail-label">Tags:</span>
                <span class="detail-value">
                  <span 
                    v-for="tag in viewingBook.tags" 
                    :key="tag" 
                    class="badge bg-light text-dark me-1"
                  >
                    {{ tag }}
                  </span>
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Added:</span>
                <span class="detail-value">{{ formatDate(viewingBook.created_at) }}</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="viewingBook = null"
            >
              Close
            </button>
            <button 
              v-if="!isBorrowed(viewingBook)"
              type="button" 
              class="btn btn-success" 
              @click="handleBorrow(viewingBook)"
            >
              <i class="bi bi-bookmark-plus me-2"></i>
              Borrow Book
            </button>
            <button 
              v-if="isBorrowed(viewingBook)"
              type="button" 
              class="btn btn-warning" 
              @click="handleReturn(viewingBook)"
            >
              <i class="bi bi-bookmark-check me-2"></i>
              Return Book
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="viewingBook" class="modal-backdrop fade show"></div>

    <!-- Borrow Modal -->
    <div 
      v-if="borrowingBook"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="borrowingBook = null"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-bookmark-plus me-2"></i>
              Borrow Book
            </h5>
            <button 
              type="button" 
              class="btn-close" 
              @click="borrowingBook = null"
            ></button>
          </div>
          <div class="modal-body">
            <p>You are about to borrow: <strong>{{ borrowingBook.name }}</strong></p>
            <div class="mb-3">
              <label for="borrowDeadline" class="form-label">Return Deadline</label>
              <input
                type="datetime-local"
                class="form-control"
                id="borrowDeadline"
                v-model="borrowDeadline"
                required
              />
            </div>
          </div>
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="borrowingBook = null"
            >
              Cancel
            </button>
            <button 
              type="button" 
              class="btn btn-success" 
              @click="confirmBorrow"
              :disabled="!borrowDeadline"
            >
              <i class="bi bi-bookmark-plus me-2"></i>
              Borrow Book
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="borrowingBook" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useBookStore } from '@/stores/bookStore';
import type { Book } from '@/types';
import { formatDate } from '@/utils/dateUtils';
import BookForm from '@/components/books/BookForm.vue';
import BookList from '@/components/books/BookList.vue';

const bookStore = useBookStore();

// UI State
const showForm = ref(false);
const editingBook = ref<Book | undefined>(undefined);
const viewingBook = ref<Book | null>(null);
const borrowingBook = ref<Book | null>(null);
const borrowDeadline = ref('');

// Search and Pagination State
const searchQuery = ref('');
const currentPage = ref(1);
const pageSize = ref(20);

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>;
watch([searchQuery], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    loadBooks();
  }, 300);
});

watch(currentPage, () => {
  loadBooks();
});

// Load books
const loadBooks = async () => {
  try {
    await bookStore.fetchBooks({
      q: searchQuery.value,
      page: currentPage.value,
      page_size: pageSize.value
    });
  } catch (error) {
    console.error('Failed to load books:', error);
  }
};

// Form handlers
const showCreateForm = () => {
  editingBook.value = undefined;
  showForm.value = true;
};

const closeForm = () => {
  showForm.value = false;
  editingBook.value = undefined;
};

const handleFormSubmit = async (book: Book) => {
  try {
    if (editingBook.value) {
      await bookStore.updateBook(editingBook.value.id, book);
    } else {
      await bookStore.addBook(book);
    }
    closeForm();
    await loadBooks();
  } catch (error) {
    console.error('Failed to save book:', error);
    alert('Failed to save book. Please try again.');
  }
};

// List action handlers
const handleView = (book: Book) => {
  viewingBook.value = book;
};

const handleEdit = (book: Book) => {
  editingBook.value = book;
  showForm.value = true;
};

const handleDelete = async (id: string) => {
  if (confirm('Are you sure you want to delete this book?')) {
    try {
      await bookStore.deleteBook(id);
      await loadBooks();
    } catch (error) {
      console.error('Failed to delete book:', error);
      alert('Failed to delete book. Please try again.');
    }
  }
};

const handleBorrow = (book: Book) => {
  borrowingBook.value = book;
  // Set default deadline to 2 weeks from now
  const defaultDeadline = new Date();
  defaultDeadline.setDate(defaultDeadline.getDate() + 14);
  borrowDeadline.value = defaultDeadline.toISOString().slice(0, 16);
};

const confirmBorrow = async () => {
  if (!borrowingBook.value || !borrowDeadline.value) return;
  
  try {
    await bookStore.borrowBook(borrowingBook.value.id, new Date(borrowDeadline.value));
    borrowingBook.value = null;
    borrowDeadline.value = '';
    await loadBooks();
    if (viewingBook.value) {
      viewingBook.value = null;
    }
  } catch (error) {
    console.error('Failed to borrow book:', error);
    alert('Failed to borrow book. Please try again.');
  }
};

const handleReturn = async (book: Book) => {
  if (confirm(`Are you sure you want to return "${book.name}"?`)) {
    try {
      await bookStore.returnBook(book.id);
      await loadBooks();
      if (viewingBook.value) {
        viewingBook.value = null;
      }
    } catch (error) {
      console.error('Failed to return book:', error);
      alert('Failed to return book. Please try again.');
    }
  }
};

// Helper functions
const isBorrowed = (book: Book): boolean => {
  return book.borrow_time != null && book.return_time == null;
};

const isOverdue = (book: Book): boolean => {
  if (!isBorrowed(book) || !book.deadline) return false;
  return new Date() > book.deadline;
};

const getBookStatus = (book: Book): string => {
  if (isOverdue(book)) return 'Overdue';
  if (isBorrowed(book)) return 'Borrowed';
  return 'Available';
};

const getStatusBadgeClass = (book: Book) => {
  if (isOverdue(book)) return 'badge bg-danger';
  if (isBorrowed(book)) return 'badge bg-warning';
  return 'badge bg-success';
};

// Initialize
onMounted(() => {
  loadBooks();
});
</script>

<style scoped>
/* Page Layout */
.books-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #e9ecef;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.page-title i {
  color: #667eea;
}

.page-description {
  color: #6c757d;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

/* Content Wrapper */
.content-wrapper {
  position: relative;
}

/* Form Section */
.form-section {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e9ecef;
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

/* Loading Overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* Modal Styles */
.modal {
  background: rgba(0, 0, 0, 0.5);
}

.modal-dialog {
  margin-top: 5rem;
}

.book-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f8f9fa;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #6c757d;
  min-width: 120px;
}

.detail-value {
  color: #212529;
}

/* Badge Styles */
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .books-view {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .header-actions {
    width: 100%;
  }
  
  .header-actions .btn {
    flex: 1;
  }
  
  .form-section {
    padding: 1rem;
  }
  
  .detail-row {
    flex-direction: column;
    align-items: start;
    gap: 0.25rem;
  }
}
</style>