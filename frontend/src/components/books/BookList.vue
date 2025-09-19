<template>
  <div class="book-list">
    <!-- Search and Filters -->
    <div class="search-section">
      <div class="search-bar">
        <div class="input-group">
          <span class="input-group-text">
            <i class="bi bi-search"></i>
          </span>
          <input
            type="text"
            class="form-control"
            :value="searchQuery"
            @input="handleSearchChange"
            placeholder="Search books by name, author, or description..."
          />
        </div>
      </div>
      
      <div class="filter-controls">
        <div class="filter-group">
          <label class="form-label">Author:</label>
          <input
            type="text"
            class="form-control form-control-sm"
            :value="authorFilter"
            @input="handleAuthorFilterChange"
            placeholder="Filter by author"
          />
        </div>
        
        <div class="filter-group">
          <label class="form-label">Tags:</label>
          <input
            type="text"
            class="form-control form-control-sm"
            :value="tagsFilter"
            @input="handleTagsFilterChange"
            placeholder="Filter by tags"
          />
        </div>
        
        <div class="filter-group">
          <label class="form-label">Status:</label>
          <select class="form-select form-select-sm" :value="statusFilter" @change="handleStatusFilterChange">
            <option value="">All Books</option>
            <option value="available">Available</option>
            <option value="borrowed">Borrowed</option>
            <option value="overdue">Overdue</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Books Grid -->
    <div class="books-grid" v-if="!loading && books.length > 0">
      <div 
        v-for="book in books" 
        :key="book.id" 
        class="book-card"
        :class="{ 'borrowed': isBorrowed(book), 'overdue': isOverdue(book) }"
      >
        <div class="book-header">
          <h5 class="book-title">{{ book.name }}</h5>
          <div class="book-status">
            <span v-if="isOverdue(book)" class="badge bg-danger">
              <i class="bi bi-exclamation-triangle"></i> Overdue
            </span>
            <span v-else-if="isBorrowed(book)" class="badge bg-warning">
              <i class="bi bi-bookmark-check"></i> Borrowed
            </span>
            <span v-else class="badge bg-success">
              <i class="bi bi-check-circle"></i> Available
            </span>
          </div>
        </div>
        
        <div class="book-details">
          <div class="book-info">
            <div class="info-item">
              <i class="bi bi-person"></i>
              <span>{{ book.author || 'Unknown Author' }}</span>
            </div>
            <div class="info-item">
              <i class="bi bi-upc-scan"></i>
              <span>{{ book.isbn }}</span>
            </div>
            <div class="info-item" v-if="book.price > 0">
              <i class="bi bi-currency-dollar"></i>
              <span>${{ book.price.toFixed(2) }}</span>
            </div>
          </div>
          
          <div class="book-description" v-if="book.description">
            <p>{{ book.description }}</p>
          </div>
          
          <div class="book-tags" v-if="book.tags && book.tags.length > 0">
            <span 
              v-for="tag in book.tags" 
              :key="tag" 
              class="badge bg-light text-dark me-1"
            >
              {{ tag }}
            </span>
          </div>
          
          <div class="borrow-info" v-if="isBorrowed(book)">
            <div class="borrow-details">
              <small class="text-muted">
                <i class="bi bi-calendar-event"></i>
                Borrowed: {{ formatDate(book.borrow_time!) }}
              </small>
              <small class="text-muted" v-if="book.deadline">
                <i class="bi bi-calendar-check"></i>
                Due: {{ formatDate(book.deadline!) }}
              </small>
            </div>
          </div>
        </div>
        
        <div class="book-actions">
          <button 
            class="btn btn-sm btn-outline-primary" 
            @click="$emit('view', book)"
            title="View Details"
          >
            <i class="bi bi-eye"></i>
          </button>
          <button 
            class="btn btn-sm btn-outline-secondary" 
            @click="$emit('edit', book)"
            title="Edit Book"
          >
            <i class="bi bi-pencil"></i>
          </button>
          <button 
            v-if="!isBorrowed(book)"
            class="btn btn-sm btn-outline-success" 
            @click="$emit('borrow', book)"
            title="Borrow Book"
          >
            <i class="bi bi-bookmark-plus"></i>
          </button>
          <button 
            v-if="isBorrowed(book)"
            class="btn btn-sm btn-outline-warning" 
            @click="$emit('return', book)"
            title="Return Book"
          >
            <i class="bi bi-bookmark-check"></i>
          </button>
          <button 
            class="btn btn-sm btn-outline-danger" 
            @click="$emit('delete', book.id)"
            title="Delete Book"
          >
            <i class="bi bi-trash"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div class="empty-state" v-if="!loading && books.length === 0">
      <div class="empty-icon">
        <i class="bi bi-book"></i>
      </div>
      <h4>No books found</h4>
      <p>Try adjusting your search criteria or add some books to get started.</p>
    </div>

    <!-- Loading State -->
    <div class="loading-state" v-if="loading">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p>Loading books...</p>
    </div>

    <!-- Pagination -->
    <div class="pagination-section" v-if="totalCount > pageSize">
      <nav aria-label="Books pagination">
        <ul class="pagination justify-content-center">
          <li class="page-item" :class="{ disabled: currentPage <= 1 }">
            <button 
              class="page-link" 
              @click="changePage(currentPage - 1)"
              :disabled="currentPage <= 1"
            >
              <i class="bi bi-chevron-left"></i>
            </button>
          </li>
          
          <li 
            v-for="page in visiblePages" 
            :key="page" 
            class="page-item" 
            :class="{ active: page === currentPage }"
          >
            <button 
              class="page-link" 
              @click="changePage(Number(page))"
            >
              {{ page }}
            </button>
          </li>
          
          <li class="page-item" :class="{ disabled: currentPage >= totalPages }">
            <button 
              class="page-link" 
              @click="changePage(currentPage + 1)"
              :disabled="currentPage >= totalPages"
            >
              <i class="bi bi-chevron-right"></i>
            </button>
          </li>
        </ul>
      </nav>
      
      <div class="pagination-info">
        <small class="text-muted">
          Showing {{ startItem }} to {{ endItem }} of {{ totalCount }} books
        </small>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { Book } from '@/types';
import { formatDate } from '@/utils/dateUtils';

const props = defineProps<{
  books: Book[];
  searchQuery: string;
  currentPage: number;
  pageSize: number;
  totalCount: number;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'view', book: Book): void;
  (e: 'edit', book: Book): void;
  (e: 'delete', id: string): void;
  (e: 'borrow', book: Book): void;
  (e: 'return', book: Book): void;
  (e: 'update:search-query', value: string): void;
  (e: 'update:page', value: number): void;
}>();

// Local filter state
const authorFilter = ref('');
const tagsFilter = ref('');
const statusFilter = ref('');

// Computed properties
const totalPages = computed(() => Math.ceil(props.totalCount / props.pageSize));
const startItem = computed(() => (props.currentPage - 1) * props.pageSize + 1);
const endItem = computed(() => Math.min(props.currentPage * props.pageSize, props.totalCount));

const visiblePages = computed(() => {
  const pages = [];
  const start = Math.max(1, props.currentPage - 2);
  const end = Math.min(totalPages.value, props.currentPage + 2);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});

// Helper functions
const isBorrowed = (book: Book): boolean => {
  return book.borrow_time != null && book.return_time == null;
};

const isOverdue = (book: Book): boolean => {
  if (!isBorrowed(book) || !book.deadline) return false;
  return new Date() > book.deadline;
};

// Event handlers
const handleSearchChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:search-query', target.value);
};

const handleAuthorFilterChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  authorFilter.value = target.value;
  // For now, we'll just emit the search query change
  // In a real implementation, you'd emit separate filter events
  emit('update:search-query', props.searchQuery);
};

const handleTagsFilterChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  tagsFilter.value = target.value;
  // For now, we'll just emit the search query change
  // In a real implementation, you'd emit separate filter events
  emit('update:search-query', props.searchQuery);
};

const handleStatusFilterChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  statusFilter.value = target.value;
  // For now, we'll just emit the search query change
  // In a real implementation, you'd emit separate filter events
  emit('update:search-query', props.searchQuery);
};

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    emit('update:page', page);
  }
};
</script>

<style scoped>
/* Book List Container */
.book-list {
  padding: 1rem 0;
}

/* Search Section */
.search-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.search-bar {
  margin-bottom: 1rem;
}

.search-bar .input-group-text {
  background-color: #f8f9fa;
  border-color: #e9ecef;
}

.filter-controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.filter-group .form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6c757d;
  margin: 0;
}

/* Books Grid */
.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

/* Book Card */
.book-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 2px solid transparent;
  transition: all 0.3s ease;
  position: relative;
}

.book-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.book-card.borrowed {
  border-color: #ffc107;
  background: linear-gradient(135deg, #fff9e6 0%, #ffffff 100%);
}

.book-card.overdue {
  border-color: #dc3545;
  background: linear-gradient(135deg, #ffe6e6 0%, #ffffff 100%);
}

/* Book Header */
.book-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.book-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
  margin: 0;
  flex: 1;
  margin-right: 1rem;
}

.book-status {
  flex-shrink: 0;
}

/* Book Details */
.book-details {
  margin-bottom: 1rem;
}

.book-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: #6c757d;
}

.info-item i {
  width: 16px;
  text-align: center;
}

.book-description {
  margin-bottom: 1rem;
}

.book-description p {
  font-size: 0.9rem;
  color: #495057;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.book-tags {
  margin-bottom: 1rem;
}

.borrow-info {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 0.75rem;
  margin-bottom: 1rem;
}

.borrow-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

/* Book Actions */
.book-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.book-actions .btn {
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.book-actions .btn:hover {
  transform: translateY(-1px);
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6c757d;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h4 {
  color: #495057;
  margin-bottom: 0.5rem;
}

/* Loading State */
.loading-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6c757d;
}

.loading-state .spinner-border {
  margin-bottom: 1rem;
}

/* Pagination */
.pagination-section {
  margin-top: 2rem;
}

.pagination {
  margin-bottom: 1rem;
}

.page-link {
  border-radius: 6px;
  margin: 0 2px;
  border: 1px solid #e9ecef;
  color: #495057;
  transition: all 0.2s ease;
}

.page-link:hover {
  background-color: #e9ecef;
  border-color: #dee2e6;
}

.page-item.active .page-link {
  background-color: #667eea;
  border-color: #667eea;
}

.page-item.disabled .page-link {
  color: #6c757d;
  background-color: #f8f9fa;
  border-color: #e9ecef;
}

.pagination-info {
  text-align: center;
}

/* Badges */
.badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .books-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-controls {
    grid-template-columns: 1fr;
  }
  
  .book-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .book-title {
    margin-right: 0;
  }
  
  .book-actions {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .book-actions .btn {
    flex: 1;
    min-width: 40px;
  }
}

@media (max-width: 576px) {
  .search-section {
    padding: 1rem;
  }
  
  .book-card {
    padding: 1rem;
  }
  
  .book-actions {
    gap: 0.25rem;
  }
  
  .book-actions .btn {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }
}
</style>