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

    <!-- Bookmark Cards (Grid View) -->
    <div v-else-if="props.viewMode === 'grid'" class="row">
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

    <!-- Bookmark List (List View) -->
    <div v-else class="list-group">
      <div
        v-for="bookmark in bookmarks"
        :key="bookmark.id"
        class="list-group-item list-group-item-action"
      >
        <div class="d-flex w-100 justify-content-between align-items-start">
          <div class="flex-grow-1">
            <div class="d-flex align-items-center mb-2">
              <img
                v-if="getFaviconUrl(bookmark.url)"
                :src="getFaviconUrl(bookmark.url)"
                :alt="`${getDomainFromUrl(bookmark.url)} favicon`"
                class="favicon me-2"
                @error="handleImageError"
              />
              <i v-else class="bi bi-bookmark me-2 text-muted"></i>
              <h6 class="mb-0 me-2">
                <a
                  :href="bookmark.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-decoration-none"
                >
                  {{ bookmark.title }}
                </a>
              </h6>
              <small class="text-muted">{{ getDomainFromUrl(bookmark.url) }}</small>
            </div>
            
            <p v-if="bookmark.description" class="mb-2 text-muted">
              {{ bookmark.description }}
            </p>
            
            <div class="d-flex align-items-center gap-3">
              <div v-if="bookmark.tags && bookmark.tags.length > 0" class="d-flex flex-wrap gap-1">
                <span
                  v-for="tag in bookmark.tags"
                  :key="tag.id"
                  class="badge bg-light text-dark cursor-pointer"
                  @click="handleTagFilter(tag.name)"
                >
                  {{ tag.name }}
                </span>
              </div>
              
              <small v-if="getCategoryName(bookmark.category_id)" class="text-muted">
                <i class="bi bi-folder me-1"></i>
                <span
                  class="cursor-pointer"
                  @click="handleCategoryFilter(bookmark.category_id?.toString() || '')"
                >
                  {{ getCategoryName(bookmark.category_id) }}
                </span>
              </small>
              
              <small class="text-muted">
                {{ formatRelativeTime(bookmark.created_at) }}
              </small>
            </div>
          </div>
          
          <div class="dropdown ms-3">
            <button
              class="btn btn-sm btn-outline-secondary dropdown-toggle"
              type="button"
              :id="`bookmark-menu-${bookmark.id}`"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <i class="bi bi-three-dots"></i>
            </button>
            <ul class="dropdown-menu" :aria-labelledby="`bookmark-menu-${bookmark.id}`">
              <li>
                <button class="dropdown-item" @click="$emit('edit', bookmark)">
                  <i class="bi bi-pencil me-2"></i>
                  Edit
                </button>
              </li>
              <li>
                <a
                  class="dropdown-item"
                  :href="bookmark.url"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <i class="bi bi-box-arrow-up-right me-2"></i>
                  Open
                </a>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <button
                  class="dropdown-item text-danger"
                  @click="handleDelete(bookmark)"
                >
                  <i class="bi bi-trash me-2"></i>
                  Delete
                </button>
              </li>
            </ul>
          </div>
        </div>
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
import { ref, computed, onMounted, watch, withDefaults } from 'vue';
import type { Bookmark, BookmarkCategory, BookmarkListRequest } from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';
import BookmarkCard from './BookmarkCard.vue';

interface Props {
  initialFilters?: Partial<BookmarkListRequest>;
  viewMode?: 'grid' | 'list';
}

interface Emits {
  (e: 'edit', bookmark: Bookmark): void;
  (e: 'delete', bookmark: Bookmark): void;
  (e: 'create'): void;
}

const props = withDefaults(defineProps<Props>(), {
  viewMode: 'grid'
});
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
const getDomainFromUrl = (url: string): string => {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname;
  } catch {
    return url;
  }
};

const formatRelativeTime = (date: string | Date): string => {
  const now = new Date();
  const bookmarkDate = new Date(date);
  const diff = now.getTime() - bookmarkDate.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60));
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60));
      return `${minutes}m ago`;
    }
    return `${hours}h ago`;
  } else if (days === 1) {
    return 'Yesterday';
  } else if (days < 7) {
    return `${days}d ago`;
  } else {
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric'
    }).format(bookmarkDate);
  }
};

const getCategoryName = (categoryId: string): string => {
  const category = categories.value.find(c => c.id.toString() === categoryId);
  return category ? category.name : '';
};

const getFaviconUrl = (url: string): string => {
  try {
    const urlObj = new URL(url);
    return `https://www.google.com/s2/favicons?domain=${urlObj.hostname}&sz=16`;
  } catch {
    return '';
  }
};

const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement;
  if (target) {
    target.style.display = 'none';
  }
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

.cursor-pointer {
  cursor: pointer;
}

.cursor-pointer:hover {
  text-decoration: underline;
}

.favicon {
  width: 16px;
  height: 16px;
  object-fit: contain;
}

.list-group-item {
  border: 1px solid rgba(0, 0, 0, 0.125);
  padding: 1rem;
}

.list-group-item:hover {
  background-color: rgba(0, 0, 0, 0.025);
}

.list-group-item .badge {
  font-size: 0.75em;
}

.list-group-item h6 a {
  color: #0d6efd;
  font-weight: 600;
}

.list-group-item h6 a:hover {
  color: #0a58ca;
  text-decoration: underline;
}
</style>
