<template>
  <div class="bookmarks-view">
    <div class="container-fluid">
      <!-- Header -->
      <div class="d-flex justify-content-between align-items-center mb-4">
        <div>
          <h1 class="h2 mb-0">
            <i class="bi bi-bookmark-fill me-2 text-primary"></i>
            Bookmarks
          </h1>
          <p class="text-muted mb-0">Manage your saved bookmarks and categories</p>
        </div>
        <div class="d-flex gap-2" v-if="activeTab === 'bookmarks'">
          <button
            class="btn btn-outline-info"
            @click="showStats = !showStats"
          >
            <i class="bi bi-graph-up me-2"></i>
            {{ showStats ? 'Hide' : 'Show' }} Stats
          </button>
          <button
            class="btn btn-primary"
            @click="showCreateModal = true"
          >
            <i class="bi bi-plus me-2"></i>
            Add Bookmark
          </button>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="mb-4">
        <ul class="nav nav-tabs">
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'bookmarks' }"
              @click="activeTab = 'bookmarks'"
            >
              <i class="bi bi-bookmark me-2"></i>
              Bookmarks
            </button>
          </li>
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'categories' }"
              @click="activeTab = 'categories'"
            >
              <i class="bi bi-folder me-2"></i>
              Categories
            </button>
          </li>
        </ul>
      </div>

      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Bookmarks Tab -->
        <div v-if="activeTab === 'bookmarks'" class="tab-pane fade show active">
          <!-- Statistics Card -->
          <div v-if="showStats" class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">
            <i class="bi bi-graph-up me-2"></i>
            Statistics
          </h5>
        </div>
        <div class="card-body">
          <div class="row" v-if="stats">
            <div class="col-md-4">
              <div class="text-center">
                <div class="h3 text-primary">{{ stats.total_bookmarks }}</div>
                <small class="text-muted">Total Bookmarks</small>
              </div>
            </div>
            <div class="col-md-4">
              <div class="text-center">
                <div class="h3 text-info">{{ stats.total_categories }}</div>
                <small class="text-muted">Categories</small>
              </div>
            </div>
            <div class="col-md-4">
              <div class="text-center">
                <div class="h3 text-success">{{ stats.total_tags }}</div>
                <small class="text-muted">Unique Tags</small>
              </div>
            </div>
          </div>
          <div v-else class="text-center">
            <div class="spinner-border spinner-border-sm me-2"></div>
            Loading statistics...
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="row mb-4">
        <div class="col-md-6">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">
                <i class="bi bi-clock-history me-2"></i>
                Recent Bookmarks
              </h5>
              <div v-if="recentBookmarks.length > 0">
                <div
                  v-for="bookmark in recentBookmarks"
                  :key="bookmark.id"
                  class="d-flex justify-content-between align-items-center py-2 border-bottom"
                >
                  <div class="flex-grow-1">
                    <a
                      :href="bookmark.url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-decoration-none fw-medium"
                    >
                      {{ bookmark.title }}
                    </a>
                    <div class="text-muted small">
                      {{ getDomainFromUrl(bookmark.url) }}
                    </div>
                  </div>
                  <small class="text-muted">
                    {{ formatRelativeTime(bookmark.created_at) }}
                  </small>
                </div>
              </div>
              <div v-else class="text-muted text-center py-3">
                No recent bookmarks
              </div>
            </div>
          </div>
        </div>

        <div class="col-md-6">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">
                <i class="bi bi-tags me-2"></i>
                Popular Tags
              </h5>
              <div v-if="popularTags.length > 0">
                <div class="d-flex flex-wrap gap-2">
                  <span
                    v-for="tag in popularTags"
                    :key="tag.name"
                    class="badge bg-light text-dark"
                    style="cursor: pointer;"
                    @click="filterByTag(tag.name)"
                  >
                    {{ tag.name }} ({{ tag.count }})
                  </span>
                </div>
              </div>
              <div v-else class="text-muted text-center py-3">
                No tags available
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Bookmark List -->
      <BookmarkList
        @edit="handleEdit"
        @delete="handleDelete"
        @create="showCreateModal = true"
      />
    </div>

    <!-- Create Modal -->
    <div
      class="modal fade"
      id="createBookmarkModal"
      tabindex="-1"
      :class="{ show: showCreateModal }"
      :style="{ display: showCreateModal ? 'block' : 'none' }"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-plus me-2"></i>
              Create Bookmark
            </h5>
            <button
              type="button"
              class="btn-close"
              @click="showCreateModal = false"
            ></button>
          </div>
          <div class="modal-body">
            <BookmarkForm
              mode="create"
              :loading="creating"
              @submit="handleCreate"
              @cancel="showCreateModal = false"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div
      class="modal fade"
      id="editBookmarkModal"
      tabindex="-1"
      :class="{ show: showEditModal }"
      :style="{ display: showEditModal ? 'block' : 'none' }"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-pencil me-2"></i>
              Edit Bookmark
            </h5>
            <button
              type="button"
              class="btn-close"
              @click="showEditModal = false"
            ></button>
          </div>
          <div class="modal-body">
            <BookmarkForm
              v-if="editingBookmark"
              mode="edit"
              :bookmark="editingBookmark"
              :loading="updating"
              @submit="handleUpdate"
              @delete="handleDelete"
              @cancel="showEditModal = false"
            />
          </div>
        </div>
      </div>
    </div>
        </div>

        <!-- Categories Tab -->
        <div v-if="activeTab === 'categories'" class="tab-pane fade show active">
          <CategoryManagement />
        </div>
      </div>

    <!-- Modal Backdrop -->
    <div
      v-if="showCreateModal || showEditModal"
      class="modal-backdrop fade show"
      @click="closeModals"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Bookmark, CreateBookmarkRequest, UpdateBookmarkRequest } from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';
import BookmarkList from '@/components/bookmarks/BookmarkList.vue';
import BookmarkForm from '@/components/bookmarks/BookmarkForm.vue';
import CategoryManagement from '@/components/bookmarks/CategoryManagement.vue';

const bookmarkStore = useBookmarkStore();

// Local state
const activeTab = ref('bookmarks');
const showStats = ref(false);
const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingBookmark = ref<Bookmark | null>(null);
const creating = ref(false);
const updating = ref(false);
const recentBookmarks = ref<Bookmark[]>([]);

// Computed properties
const stats = computed(() => bookmarkStore.stats);
const popularTags = computed(() => bookmarkStore.popularTags);

// Methods
const getDomainFromUrl = (url: string): string => {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname;
  } catch {
    return url;
  }
};

const formatRelativeTime = (date: Date): string => {
  const now = new Date();
  const diff = now.getTime() - new Date(date).getTime();
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
    }).format(new Date(date));
  }
};

const filterByTag = (tagName: string) => {
  // This could emit an event to the BookmarkList component or use router
  // For now, we'll just log it
  console.log('Filter by tag:', tagName);
};

const handleCreate = async (data: CreateBookmarkRequest | UpdateBookmarkRequest) => {
  creating.value = true;
  try {
    // Ensure required fields are present for creation
    const createData: CreateBookmarkRequest = {
      url: data.url || '',
      title: data.title || '',
      description: data.description,
      tags: data.tags,
      category_id: data.category_id
    };
    await bookmarkStore.createBookmark(createData);
    showCreateModal.value = false;
    // Refresh recent bookmarks
    await loadRecentBookmarks();
  } catch (error) {
    console.error('Failed to create bookmark:', error);
  } finally {
    creating.value = false;
  }
};

const handleEdit = (bookmark: Bookmark) => {
  editingBookmark.value = bookmark;
  showEditModal.value = true;
};

const handleUpdate = async (data: UpdateBookmarkRequest) => {
  if (!editingBookmark.value) return;
  
  updating.value = true;
  try {
    await bookmarkStore.updateBookmark(editingBookmark.value.id, data);
    showEditModal.value = false;
    editingBookmark.value = null;
    // Refresh recent bookmarks
    await loadRecentBookmarks();
  } catch (error) {
    console.error('Failed to update bookmark:', error);
  } finally {
    updating.value = false;
  }
};

const handleDelete = async (bookmark?: Bookmark) => {
  const bookmarkToDelete = bookmark || editingBookmark.value;
  if (!bookmarkToDelete) return;
  
  try {
    await bookmarkStore.deleteBookmark(bookmarkToDelete.id);
    if (showEditModal.value) {
      showEditModal.value = false;
      editingBookmark.value = null;
    }
    // Refresh recent bookmarks
    await loadRecentBookmarks();
  } catch (error) {
    console.error('Failed to delete bookmark:', error);
  }
};

const closeModals = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  editingBookmark.value = null;
};

const loadRecentBookmarks = async () => {
  try {
    recentBookmarks.value = await bookmarkStore.fetchRecentBookmarks(5);
  } catch (error) {
    console.error('Failed to load recent bookmarks:', error);
  }
};

// Load initial data
onMounted(async () => {
  try {
    await Promise.all([
      bookmarkStore.fetchBookmarkStats(),
      bookmarkStore.fetchPopularTags(10),
      loadRecentBookmarks()
    ]);
  } catch (error) {
    console.error('Failed to load initial data:', error);
  }
});
</script>

<style scoped>
.bookmarks-view {
  min-height: 100vh;
  padding: 2rem 0;
}

.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.modal.show {
  display: block !important;
}

.modal-backdrop {
  background-color: rgba(0, 0, 0, 0.5);
}

.card {
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  border: 1px solid rgba(0, 0, 0, 0.125);
}

.card:hover {
  box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15);
  transition: box-shadow 0.15s ease-in-out;
}

.badge {
  transition: background-color 0.2s ease-in-out;
}

.badge:hover {
  background-color: #e9ecef !important;
  color: #495057 !important;
}

.border-bottom:last-child {
  border-bottom: none !important;
}

h1.h2 {
  font-weight: 600;
}

.text-primary {
  color: #0d6efd !important;
}

.nav-tabs .nav-link {
  border: 1px solid transparent;
  border-radius: 0.375rem 0.375rem 0 0;
  background: none;
  color: #6c757d;
  font-weight: 500;
}

.nav-tabs .nav-link:hover {
  border-color: #e9ecef #e9ecef #dee2e6;
  color: #495057;
}

.nav-tabs .nav-link.active {
  color: #495057;
  background-color: #fff;
  border-color: #dee2e6 #dee2e6 #fff;
}

.tab-content {
  border: 1px solid #dee2e6;
  border-top: none;
  border-radius: 0 0 0.375rem 0.375rem;
  padding: 1.5rem;
  background-color: #fff;
}
</style>
