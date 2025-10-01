<template>
  <div class="posts-view">
    <!-- Header -->
    <div class="d-flex justify-content-between align-items-center mb-4">
      <div>
        <h1 class="h3 mb-0">Posts</h1>
        <p class="text-muted mb-0">Manage your blog posts and pages</p>
      </div>
      <div class="d-flex gap-2">
        <button 
          class="btn btn-outline-primary"
          @click="showCreateModal = true"
        >
          <i class="bi bi-plus-lg"></i>
          New Post
        </button>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-4">
            <label class="form-label">Search</label>
            <div class="input-group">
              <input
                v-model="searchQuery"
                type="text"
                class="form-control"
                placeholder="Search posts..."
                @keyup.enter="handleSearch"
              >
              <button 
                class="btn btn-outline-secondary" 
                type="button"
                @click="handleSearch"
              >
                <i class="bi bi-search"></i>
              </button>
            </div>
          </div>
          <div class="col-md-3">
            <label class="form-label">Status</label>
            <select v-model="selectedStatus" class="form-select" @change="handleFilterChange">
              <option value="">All Status</option>
              <option value="draft">Draft</option>
              <option value="pending">Pending</option>
              <option value="published">Published</option>
              <option value="private">Private</option>
              <option value="scheduled">Scheduled</option>
              <option value="trash">Trash</option>
            </select>
          </div>
          <div class="col-md-3">
            <label class="form-label">Type</label>
            <select v-model="selectedType" class="form-select" @change="handleFilterChange">
              <option value="">All Types</option>
              <option value="post">Post</option>
              <option value="page">Page</option>
            </select>
          </div>
          <div class="col-md-2">
            <label class="form-label">&nbsp;</label>
            <div class="d-grid">
              <button class="btn btn-outline-secondary" @click="clearFilters">
                Clear
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Posts List -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <span>Posts ({{ postStore.totalPosts }})</span>
        <div class="d-flex gap-2">
          <button 
            class="btn btn-sm btn-outline-secondary"
            @click="refreshPosts"
            :disabled="postStore.loading"
          >
            <i class="bi bi-arrow-clockwise"></i>
            Refresh
          </button>
        </div>
      </div>
      
      <div class="card-body p-0">
        <!-- Loading State -->
        <div v-if="postStore.loading" class="text-center p-4">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
          <p class="mt-2 mb-0 text-muted">Loading posts...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="postStore.error" class="alert alert-danger m-3">
          <i class="bi bi-exclamation-triangle"></i>
          {{ postStore.error }}
          <button class="btn btn-sm btn-outline-danger ms-2" @click="refreshPosts">
            Retry
          </button>
        </div>

        <!-- Empty State -->
        <div v-else-if="!postStore.hasPosts" class="text-center p-5">
          <i class="bi bi-file-earmark-text text-muted" style="font-size: 3rem;"></i>
          <h5 class="mt-3">No posts found</h5>
          <p class="text-muted">
            {{ searchQuery ? 'No posts match your search criteria.' : 'Get started by creating your first post.' }}
          </p>
          <button 
            v-if="!searchQuery"
            class="btn btn-primary"
            @click="showCreateModal = true"
          >
            Create Your First Post
          </button>
        </div>

        <!-- Posts Table -->
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead class="table-light">
              <tr>
                <th>Title</th>
                <th>Status</th>
                <th>Type</th>
                <th>Categories</th>
                <th>Views</th>
                <th>Date</th>
                <th width="120">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="post in postStore.posts" :key="post.id">
                <td>
                  <div>
                    <strong>{{ post.title }}</strong>
                    <div class="small text-muted">
                      <i class="bi bi-link-45deg"></i>
                      /{{ post.slug }}
                    </div>
                  </div>
                </td>
                <td>
                  <span 
                    class="badge"
                    :class="getStatusBadgeClass(post.status)"
                  >
                    {{ formatStatus(post.status) }}
                  </span>
                  <div v-if="post.is_sticky" class="small text-warning">
                    <i class="bi bi-pin-fill"></i> Sticky
                  </div>
                </td>
                <td>
                  <span class="badge bg-light text-dark">
                    {{ post.type }}
                  </span>
                </td>
                <td>
                  <div v-if="post.categories.length > 0">
                    <span 
                      v-for="category in post.categories.slice(0, 2)" 
                      :key="category"
                      class="badge bg-secondary me-1"
                    >
                      {{ category }}
                    </span>
                    <span v-if="post.categories.length > 2" class="small text-muted">
                      +{{ post.categories.length - 2 }} more
                    </span>
                  </div>
                  <span v-else class="text-muted small">Uncategorized</span>
                </td>
                <td>
                  <span class="text-muted">{{ post.view_count.toLocaleString() }}</span>
                </td>
                <td>
                  <div class="small">
                    <div v-if="post.published_at">
                      <strong>Published</strong><br>
                      {{ formatDate(post.published_at) }}
                    </div>
                    <div v-else-if="post.scheduled_for">
                      <strong>Scheduled</strong><br>
                      {{ formatDate(post.scheduled_for) }}
                    </div>
                    <div v-else>
                      <strong>Modified</strong><br>
                      {{ formatDate(post.updated_at) }}
                    </div>
                  </div>
                </td>
                <td>
                  <div class="dropdown">
                    <button 
                      class="btn btn-sm btn-outline-secondary dropdown-toggle"
                      data-bs-toggle="dropdown"
                    >
                      Actions
                    </button>
                    <ul class="dropdown-menu">
                      <li>
                        <button class="dropdown-item" @click="editPost(post)">
                          <i class="bi bi-pencil"></i> Edit
                        </button>
                      </li>
                      <li v-if="post.status !== 'published'">
                        <button class="dropdown-item" @click="publishPost(post)">
                          <i class="bi bi-globe"></i> Publish
                        </button>
                      </li>
                      <li v-if="post.status === 'published'">
                        <a 
                          :href="`/blog/${post.slug}`" 
                          class="dropdown-item"
                          target="_blank"
                        >
                          <i class="bi bi-eye"></i> View
                        </a>
                      </li>
                      <li>
                        <button class="dropdown-item" @click="duplicatePost(post)">
                          <i class="bi bi-files"></i> Duplicate
                        </button>
                      </li>
                      <li><hr class="dropdown-divider"></li>
                      <li>
                        <button 
                          class="dropdown-item text-danger" 
                          @click="confirmDeletePost(post)"
                        >
                          <i class="bi bi-trash"></i> Delete
                        </button>
                      </li>
                    </ul>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="postStore.totalPages > 1" class="card-footer">
        <nav>
          <ul class="pagination pagination-sm justify-content-center mb-0">
            <li class="page-item" :class="{ disabled: postStore.currentPage <= 1 }">
              <button 
                class="page-link" 
                @click="changePage(postStore.currentPage - 1)"
                :disabled="postStore.currentPage <= 1"
              >
                Previous
              </button>
            </li>
            <li 
              v-for="page in getVisiblePages()" 
              :key="page"
              class="page-item" 
              :class="{ active: page === postStore.currentPage }"
            >
              <button class="page-link" @click="changePage(Number(page))">
                {{ page }}
              </button>
            </li>
            <li class="page-item" :class="{ disabled: postStore.currentPage >= postStore.totalPages }">
              <button 
                class="page-link" 
                @click="changePage(postStore.currentPage + 1)"
                :disabled="postStore.currentPage >= postStore.totalPages"
              >
                Next
              </button>
            </li>
          </ul>
        </nav>
      </div>
    </div>

    <!-- Create/Edit Post Modal -->
    <PostFormModal
      v-if="showCreateModal || showEditModal"
      :show="showCreateModal || showEditModal"
      :post="editingPost"
      :is-page="showCreatePageModal"
      @close="closeModals"
      @saved="handlePostSaved"
    />

    <!-- Delete Confirmation Modal -->
    <div 
      v-if="showDeleteModal"
      class="modal fade show d-block"
      tabindex="-1"
      style="background-color: rgba(0,0,0,0.5)"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Confirm Delete</h5>
            <button type="button" class="btn-close" @click="showDeleteModal = false"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to delete "<strong>{{ deletingPost?.title }}</strong>"?</p>
            <p class="text-muted small">This action cannot be undone.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showDeleteModal = false">
              Cancel
            </button>
            <button 
              type="button" 
              class="btn btn-danger"
              @click="handleDeletePost"
              :disabled="postStore.loading"
            >
              <span v-if="postStore.loading" class="spinner-border spinner-border-sm me-2"></span>
              Delete Post
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePostStore, type Post } from '@/stores/postStore'
import PostFormModal from '@/components/posts/PostFormModal.vue'

// Store
const postStore = usePostStore()
const route = useRoute()
const router = useRouter()

// Reactive state
const searchQuery = ref('')
const selectedStatus = ref('')
const selectedType = ref('')
const showCreateModal = ref(false)
const showCreatePageModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingPost = ref<Post | null>(null)
const deletingPost = ref<Post | null>(null)

// Methods
const loadPostForEdit = async (postId: string) => {
  try {
    const post = await postStore.fetchPostById(postId)
    if (post) {
      editingPost.value = post
      showEditModal.value = true
    } else {
      console.error('Post not found for editing:', postId)
    }
  } catch (error) {
    console.error('Failed to load post for editing:', error)
  }
}

// Watchers
watch(
  () => route.query.edit,
  async (editId) => {
    if (editId && typeof editId === 'string') {
      await loadPostForEdit(editId)
    }
  },
  { immediate: true }
)

// Lifecycle
onMounted(async () => {
  await refreshPosts()
  
  // Check for edit parameter on mount
  const editId = route.query.edit as string
  if (editId) {
    await loadPostForEdit(editId)
  }
})

// Methods
const refreshPosts = async () => {
  if (searchQuery.value) {
    await postStore.searchPosts(searchQuery.value, 1, selectedStatus.value, selectedType.value)
  } else {
    await postStore.fetchPosts(1, selectedStatus.value, selectedType.value)
  }
}

const handleSearch = async () => {
  if (searchQuery.value.trim()) {
    await postStore.searchPosts(searchQuery.value, 1, selectedStatus.value, selectedType.value)
  } else {
    await postStore.fetchPosts(1, selectedStatus.value, selectedType.value)
  }
}

const handleFilterChange = async () => {
  await refreshPosts()
}

const clearFilters = async () => {
  searchQuery.value = ''
  selectedStatus.value = ''
  selectedType.value = ''
  await postStore.fetchPosts(1)
}

const changePage = async (page: number) => {
  if (page >= 1 && page <= postStore.totalPages) {
    if (searchQuery.value) {
      await postStore.searchPosts(searchQuery.value, page, selectedStatus.value, selectedType.value)
    } else {
      await postStore.fetchPosts(page, selectedStatus.value, selectedType.value)
    }
  }
}

const editPost = (post: Post) => {
  editingPost.value = post
  showEditModal.value = true
}

const publishPost = async (post: Post) => {
  try {
    await postStore.publishPost(post.id)
    // Refresh the posts list to show updated status
    await refreshPosts()
  } catch (error) {
    console.error('Failed to publish post:', error)
  }
}

const duplicatePost = (post: Post) => {
  editingPost.value = {
    ...post,
    id: '',
    title: `${post.title} (Copy)`,
    slug: `${post.slug}-copy`,
    status: 'draft' as const,
    published_at: undefined,
    scheduled_for: undefined,
  }
  showEditModal.value = true
}

const confirmDeletePost = (post: Post) => {
  deletingPost.value = post
  showDeleteModal.value = true
}

const handleDeletePost = async () => {
  if (deletingPost.value) {
    try {
      await postStore.deletePost(deletingPost.value.id)
      showDeleteModal.value = false
      deletingPost.value = null
    } catch (error) {
      console.error('Failed to delete post:', error)
    }
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showCreatePageModal.value = false
  showEditModal.value = false
  editingPost.value = null
  
  // Clear edit query parameter when closing modal
  if (route.query.edit) {
    router.replace({ query: { ...route.query, edit: undefined } })
  }
}

const handlePostSaved = async () => {
  closeModals()
  await refreshPosts()
}

// Computed
const getVisiblePages = () => {
  const current = postStore.currentPage
  const total = postStore.totalPages
  const delta = 2
  const range = []
  const rangeWithDots = []

  for (let i = Math.max(2, current - delta); i <= Math.min(total - 1, current + delta); i++) {
    range.push(i)
  }

  if (current - delta > 2) {
    rangeWithDots.push(1, '...')
  } else {
    rangeWithDots.push(1)
  }

  rangeWithDots.push(...range)

  if (current + delta < total - 1) {
    rangeWithDots.push('...', total)
  } else {
    rangeWithDots.push(total)
  }

  return rangeWithDots.filter((item, index, array) => array.indexOf(item) === index)
}

// Utility functions
const getStatusBadgeClass = (status: string) => {
  const classes = {
    draft: 'bg-secondary',
    pending: 'bg-warning',
    published: 'bg-success',
    private: 'bg-info',
    scheduled: 'bg-primary',
    trash: 'bg-danger',
  }
  return classes[status as keyof typeof classes] || 'bg-secondary'
}

const formatStatus = (status: string) => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.posts-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
}

.table th {
  border-top: none;
  font-weight: 600;
  color: #495057;
}

.badge {
  font-size: 0.75rem;
}

.dropdown-toggle::after {
  font-size: 0.8rem;
}

.modal.show {
  animation: fadeIn 0.15s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.pagination .page-link {
  color: #6c757d;
  border-color: #dee2e6;
}

.pagination .page-item.active .page-link {
  background-color: #0d6efd;
  border-color: #0d6efd;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}
</style>
