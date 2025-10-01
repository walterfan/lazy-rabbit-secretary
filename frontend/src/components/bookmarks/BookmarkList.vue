<template>
  <div class="bookmark-list">
    <!-- Filters and Search -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row align-items-end">
          <div class="col-md-4">
            <label class="form-label">Search</label>
            <div class="input-group">
              <input
                v-model="searchQuery"
                type="text"
                class="form-control"
                placeholder="Search bookmarks..."
                @keydown.enter="handleSearch"
              />
              <button class="btn btn-outline-primary" @click="handleSearch">
                <i class="bi bi-search"></i>
              </button>
            </div>
          </div>

          <div class="col-md-2">
            <label class="form-label">Category</label>
            <select
              v-model="selectedCategory"
              class="form-select"
              @change="handleFilter"
            >
              <option value="">All Categories</option>
              <option
                v-for="category in categories"
                :key="category.id"
                :value="category.id.toString()"
              >
                {{ category.name }}
              </option>
            </select>
          </div>

          <div class="col-md-2">
            <label class="form-label">Tags</label>
            <div class="dropdown">
              <button
                class="btn btn-outline-secondary dropdown-toggle w-100"
                type="button"
                data-bs-toggle="dropdown"
              >
                Tags ({{ selectedTags.length }})
              </button>
              <ul class="dropdown-menu p-2" style="min-width: 200px;">
                <li v-for="tag in allTags" :key="tag" class="form-check">
                  <input
                    :id="`tag-${tag}`"
                    v-model="selectedTags"
                    :value="tag"
                    type="checkbox"
                    class="form-check-input"
                    @change="handleFilter"
                  />
                  <label :for="`tag-${tag}`" class="form-check-label">
                    {{ tag }}
                  </label>
                </li>
              </ul>
            </div>
          </div>

          <div class="col-md-2">
            <label class="form-label">Sort By</label>
            <select
              v-model="sortBy"
              class="form-select"
              @change="handleFilter"
            >
              <option value="created_at">Created Date</option>
              <option value="updated_at">Updated Date</option>
              <option value="title">Title</option>
            </select>
          </div>

          <div class="col-md-2">
            <label class="form-label">Order</label>
            <select
              v-model="sortOrder"
              class="form-select"
              @change="handleFilter"
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
        </div>

        <!-- Active Filters -->
        <div v-if="hasActiveFilters" class="mt-3">
          <span class="text-muted me-2">Active filters:</span>
          <span
            v-if="searchQuery"
            class="badge bg-primary me-1"
          >
            Search: "{{ searchQuery }}"
            <button
              type="button"
              class="btn-close btn-close-white ms-1"
              @click="clearSearch"
              style="font-size: 0.7em"
            ></button>
          </span>
          <span
            v-if="selectedCategory"
            class="badge bg-info me-1"
          >
            Category: {{ getCategoryName(selectedCategory) }}
            <button
              type="button"
              class="btn-close btn-close-white ms-1"
              @click="clearCategory"
              style="font-size: 0.7em"
            ></button>
          </span>
          <span
            v-for="tag in selectedTags"
            :key="tag"
            class="badge bg-secondary me-1"
          >
            {{ tag }}
            <button
              type="button"
              class="btn-close btn-close-white ms-1"
              @click="removeTag(tag)"
              style="font-size: 0.7em"
            ></button>
          </span>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary"
            @click="clearAllFilters"
          >
            Clear All
          </button>
        </div>
      </div>
    </div>

    <!-- Results Info -->
    <div class="d-flex justify-content-between align-items-center mb-3">
      <span class="text-muted">
        Showing {{ bookmarks.length }} of {{ totalCount }} bookmarks
      </span>
      <div class="d-flex align-items-center">
        <label class="form-label me-2 mb-0">Page Size:</label>
        <select
          v-model="pageSize"
          class="form-select form-select-sm"
          style="width: auto;"
          @change="handlePageSizeChange"
        >
          <option value="10">10</option>
          <option value="20">20</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-4">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading bookmarks...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger">
      <i class="bi bi-exclamation-triangle me-2"></i>
      {{ error }}
      <button class="btn btn-sm btn-outline-danger ms-2" @click="retryLoad">
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="bookmarks.length === 0" class="text-center py-5">
      <i class="bi bi-bookmark display-1 text-muted"></i>
      <h4 class="mt-3">No bookmarks found</h4>
      <p class="text-muted">
        {{ hasActiveFilters ? 'Try adjusting your filters' : 'Start by creating your first bookmark' }}
      </p>
      <button
        v-if="!hasActiveFilters"
        class="btn btn-primary"
        @click="$emit('create')"
      >
        <i class="bi bi-plus me-2"></i>Create Bookmark
      </button>
    </div>

    <!-- Bookmark Cards -->
    <div v-else class="row">
      <div
        v-for="bookmark in bookmarks"
        :key="bookmark.id"
        class="col-md-4 col-lg-3 col-xl-2 mb-4"
      >
        <BookmarkCard
          :bookmark="bookmark"
          @edit="$emit('edit', $event)"
          @delete="handleDelete"
          @tag-click="handleTagFilter"
          @category-click="handleCategoryFilter"
        />
      </div>
    </div>

    <!-- Pagination -->
    <nav v-if="totalPages > 1" class="d-flex justify-content-center mt-4">
      <ul class="pagination">
        <li class="page-item" :class="{ disabled: currentPage === 1 }">
          <button
            class="page-link"
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
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
          <button class="page-link" @click="goToPage(page)">
            {{ page }}
          </button>
        </li>
        
        <li class="page-item" :class="{ disabled: currentPage === totalPages }">
          <button
            class="page-link"
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
          >
            <i class="bi bi-chevron-right"></i>
          </button>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import type { Bookmark, BookmarkCategory, BookmarkListRequest } from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';
import BookmarkCard from './BookmarkCard.vue';

interface Props {
  initialFilters?: Partial<BookmarkListRequest>;
}

interface Emits {
  (e: 'edit', bookmark: Bookmark): void;
  (e: 'delete', bookmark: Bookmark): void;
  (e: 'create'): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const bookmarkStore = useBookmarkStore();

// Local state
const searchQuery = ref('');
const selectedCategory = ref('');
const selectedTags = ref<string[]>([]);
const sortBy = ref('created_at');
const sortOrder = ref('desc');
const pageSize = ref(20);

// Computed properties
const bookmarks = computed(() => bookmarkStore.bookmarks);
const categories = computed(() => bookmarkStore.categories);
const allTags = computed(() => bookmarkStore.allTags);
const totalCount = computed(() => bookmarkStore.totalCount);
const currentPage = computed(() => bookmarkStore.currentPage);
const loading = computed(() => bookmarkStore.loading);
const error = computed(() => bookmarkStore.error);

const totalPages = computed(() => {
  return Math.ceil(totalCount.value / pageSize.value);
});

const hasActiveFilters = computed(() => {
  return searchQuery.value || selectedCategory.value || selectedTags.value.length > 0;
});

const visiblePages = computed(() => {
  const pages: number[] = [];
  const start = Math.max(1, currentPage.value - 2);
  const end = Math.min(totalPages.value, currentPage.value + 2);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  
  return pages;
});

// Methods
const getCategoryName = (categoryId: string): string => {
  const category = categories.value.find(c => c.id.toString() === categoryId);
  return category ? category.name : 'Unknown Category';
};

const handleSearch = () => {
  loadBookmarks();
};

const handleFilter = () => {
  bookmarkStore.setCurrentPage(1);
  loadBookmarks();
};

const handlePageSizeChange = () => {
  bookmarkStore.setPageSize(pageSize.value);
  bookmarkStore.setCurrentPage(1);
  loadBookmarks();
};

const handleTagFilter = (tagName: string) => {
  if (!selectedTags.value.includes(tagName)) {
    selectedTags.value.push(tagName);
    handleFilter();
  }
};

const handleCategoryFilter = (categoryId: string) => {
  selectedCategory.value = categoryId;
  handleFilter();
};

const handleDelete = async (bookmark: Bookmark) => {
  try {
    await bookmarkStore.deleteBookmark(bookmark.id);
    // Reload current page or adjust if it becomes empty
    if (bookmarks.value.length === 1 && currentPage.value > 1) {
      bookmarkStore.setCurrentPage(currentPage.value - 1);
    }
    await loadBookmarks();
  } catch (error) {
    console.error('Failed to delete bookmark:', error);
  }
};

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    bookmarkStore.setCurrentPage(page);
    loadBookmarks();
  }
};

const clearSearch = () => {
  searchQuery.value = '';
  handleFilter();
};

const clearCategory = () => {
  selectedCategory.value = '';
  handleFilter();
};

const removeTag = (tag: string) => {
  const index = selectedTags.value.indexOf(tag);
  if (index > -1) {
    selectedTags.value.splice(index, 1);
    handleFilter();
  }
};

const clearAllFilters = () => {
  searchQuery.value = '';
  selectedCategory.value = '';
  selectedTags.value = [];
  sortBy.value = 'created_at';
  sortOrder.value = 'desc';
  handleFilter();
};

const retryLoad = () => {
  bookmarkStore.clearError();
  loadBookmarks();
};

const loadBookmarks = async () => {
  const filters: BookmarkListRequest = {
    page: currentPage.value,
    page_size: pageSize.value,
    search: searchQuery.value || undefined,
    category_id: selectedCategory.value || undefined,
    tags: selectedTags.value.length > 0 ? selectedTags.value : undefined,
    sort_by: sortBy.value,
    sort_order: sortOrder.value
  };

  await bookmarkStore.fetchBookmarks(filters);
};

// Initialize filters from props
const initializeFilters = () => {
  if (props.initialFilters) {
    searchQuery.value = props.initialFilters.search || '';
    selectedCategory.value = props.initialFilters.category_id || '';
    selectedTags.value = props.initialFilters.tags || [];
    sortBy.value = props.initialFilters.sort_by || 'created_at';
    sortOrder.value = props.initialFilters.sort_order || 'desc';
    pageSize.value = props.initialFilters.page_size || 20;
  }
};

// Load initial data
onMounted(async () => {
  initializeFilters();
  
  try {
    await Promise.all([
      bookmarkStore.fetchCategories(),
      bookmarkStore.fetchAllTags()
    ]);
    
    await loadBookmarks();
  } catch (error) {
    console.error('Failed to load initial data:', error);
  }
});

// Watch for store page changes (e.g., from external navigation)
watch(
  () => bookmarkStore.currentPage,
  (newPage) => {
    if (newPage !== currentPage.value) {
      loadBookmarks();
    }
  }
);
</script>

<style scoped>
.dropdown-menu {
  max-height: 200px;
  overflow-y: auto;
}

.form-check {
  margin-bottom: 0.5rem;
}

.badge {
  display: inline-flex;
  align-items: center;
}

.btn-close {
  --bs-btn-close-bg: none;
}

.pagination {
  --bs-pagination-color: #0d6efd;
  --bs-pagination-hover-color: #0a58ca;
  --bs-pagination-active-bg: #0d6efd;
  --bs-pagination-active-border-color: #0d6efd;
}
</style>
