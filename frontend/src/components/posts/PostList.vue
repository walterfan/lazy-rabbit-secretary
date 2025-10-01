<template>
  <div class="post-list">
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading posts...</span>
      </div>
      <p class="mt-3 text-muted">Loading posts...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger">
      <i class="bi bi-exclamation-triangle"></i>
      {{ error }}
      <button class="btn btn-sm btn-outline-danger ms-2" @click="$emit('retry')">
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="posts.length === 0" class="text-center py-5">
      <i class="bi bi-file-earmark-text text-muted" style="font-size: 4rem;"></i>
      <h4 class="mt-3">{{ emptyMessage }}</h4>
      <p class="text-muted">
        {{ emptyDescription }}
      </p>
    </div>

    <!-- Posts Grid/List -->
    <div v-else>
      <!-- View Toggle -->
      <div v-if="showViewToggle" class="d-flex justify-content-between align-items-center mb-3">
        <div class="text-muted">
          {{ posts.length }} {{ posts.length === 1 ? 'post' : 'posts' }}
          <span v-if="totalPosts > posts.length">
            of {{ totalPosts.toLocaleString() }}
          </span>
        </div>
        <div class="btn-group btn-group-sm">
          <button 
            class="btn"
            :class="viewMode === 'grid' ? 'btn-primary' : 'btn-outline-secondary'"
            @click="setViewMode('grid')"
          >
            <i class="bi bi-grid-3x2"></i>
            Grid
          </button>
          <button 
            class="btn"
            :class="viewMode === 'list' ? 'btn-primary' : 'btn-outline-secondary'"
            @click="setViewMode('list')"
          >
            <i class="bi bi-list"></i>
            List
          </button>
        </div>
      </div>

      <!-- Grid View -->
      <div 
        v-if="viewMode === 'grid'"
        class="row g-4"
      >
        <div 
          v-for="post in posts" 
          :key="post.id"
          :class="gridColClass"
        >
          <PostCard
            :post="post"
            :show-link="showLink"
            :show-tags="showTags"
            :show-stats="showStats"
            :show-status="showStatus"
            :show-actions="showActions"
            :show-content-preview="showContentPreview"
            :show-reading-time="showReadingTime"
            :excerpt-length="excerptLength"
            @category-click="$emit('categoryClick', $event)"
            @tag-click="$emit('tagClick', $event)"
            @edit="$emit('edit', $event)"
            @publish="$emit('publish', $event)"
            @duplicate="$emit('duplicate', $event)"
            @delete="$emit('delete', $event)"
          />
        </div>
      </div>

      <!-- List View -->
      <div v-else class="post-list-view">
        <div 
          v-for="(post, index) in posts" 
          :key="post.id"
          class="post-list-item"
          :class="{ 'border-bottom': index < posts.length - 1 }"
        >
          <div class="row g-3 align-items-center">
            <!-- Featured Image -->
            <div v-if="post.featured_image" class="col-md-2">
              <img 
                :src="post.featured_image" 
                :alt="post.title"
                class="img-fluid rounded"
                style="height: 80px; object-fit: cover; width: 100%;"
              >
            </div>

            <!-- Content -->
            <div :class="post.featured_image ? 'col-md-10' : 'col-12'">
              <div class="d-flex justify-content-between align-items-start">
                <div class="flex-grow-1">
                  <!-- Categories -->
                  <div v-if="post.categories.length > 0" class="mb-1">
                    <span 
                      v-for="category in post.categories.slice(0, 2)" 
                      :key="category"
                      class="badge bg-primary me-1"
                      @click="$emit('categoryClick', category)"
                      style="cursor: pointer; font-size: 0.7rem;"
                    >
                      {{ category }}
                    </span>
                  </div>

                  <!-- Title -->
                  <h5 class="mb-2">
                    <span class="text-muted me-2" style="font-weight: normal; font-size: 0.9em;">
                      {{ String(index + 1).padStart(2, ' ') }}.
                    </span>
                    <router-link 
                      v-if="showLink"
                      :to="`/blog/${post.slug}`"
                      class="text-decoration-none text-dark"
                    >
                      {{ post.title }}
                    </router-link>
                    <span v-else>{{ post.title }}</span>
                    <span v-if="post.is_sticky" class="badge bg-warning text-dark ms-2">
                      <i class="bi bi-pin-fill"></i>
                    </span>
                  </h5>

                  <!-- Excerpt -->
                  <p v-if="post.excerpt" class="text-muted mb-2">
                    {{ truncateText(post.excerpt, excerptLength) }}
                  </p>
                  <p v-else-if="showContentPreview" class="text-muted mb-2">
                    {{ truncateText(stripHtml(post.content), excerptLength) }}
                  </p>

                  <!-- Meta -->
                  <div class="d-flex align-items-center text-muted small">
                    <i class="bi bi-calendar3 me-1"></i>
                    <span class="me-3">{{ formatDate(post.published_at || post.created_at) }}</span>
                    
                    <i v-if="showStats" class="bi bi-eye me-1"></i>
                    <span v-if="showStats" class="me-3">{{ formatNumber(post.view_count) }}</span>
                    
                    <i v-if="showReadingTime" class="bi bi-clock me-1"></i>
                    <span v-if="showReadingTime">{{ estimateReadingTime(post.content) }} min read</span>
                  </div>

                  <!-- Tags -->
                  <div v-if="post.tags.length > 0 && showTags" class="mt-2">
                    <span 
                      v-for="tag in post.tags.slice(0, 4)" 
                      :key="tag"
                      class="badge bg-light text-dark me-1"
                      @click="$emit('tagClick', tag)"
                      style="cursor: pointer; font-size: 0.65rem;"
                    >
                      #{{ tag }}
                    </span>
                  </div>
                </div>

                <!-- Actions -->
                <div v-if="showActions" class="ms-3">
                  <div class="dropdown">
                    <button 
                      class="btn btn-sm btn-outline-secondary dropdown-toggle"
                      data-bs-toggle="dropdown"
                    >
                      Actions
                    </button>
                    <ul class="dropdown-menu">
                      <li>
                        <button class="dropdown-item" @click="$emit('edit', post)">
                          <i class="bi bi-pencil"></i> Edit
                        </button>
                      </li>
                      <li v-if="post.status !== 'published'">
                        <button class="dropdown-item" @click="$emit('publish', post)">
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
                        <button class="dropdown-item" @click="$emit('duplicate', post)">
                          <i class="bi bi-files"></i> Duplicate
                        </button>
                      </li>
                      <li><hr class="dropdown-divider"></li>
                      <li>
                        <button 
                          class="dropdown-item text-danger" 
                          @click="$emit('delete', post)"
                        >
                          <i class="bi bi-trash"></i> Delete
                        </button>
                      </li>
                    </ul>
                  </div>
                </div>

                <!-- Status (for admin) -->
                <div v-if="showStatus && !showActions" class="ms-3">
                  <span 
                    class="badge"
                    :class="getStatusBadgeClass(post.status)"
                  >
                    {{ formatStatus(post.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Load More Button -->
      <div v-if="showLoadMore && hasMore" class="text-center mt-4">
        <button 
          class="btn btn-outline-primary"
          @click="$emit('loadMore')"
          :disabled="loading"
        >
          <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
          Load More Posts
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Post } from '@/stores/postStore'
import PostCard from './PostCard.vue'

// Props
interface Props {
  posts: Post[]
  loading?: boolean
  error?: string | null
  totalPosts?: number
  hasMore?: boolean
  showLink?: boolean
  showTags?: boolean
  showStats?: boolean
  showStatus?: boolean
  showActions?: boolean
  showContentPreview?: boolean
  showReadingTime?: boolean
  showLoadMore?: boolean
  showViewToggle?: boolean
  excerptLength?: number
  emptyMessage?: string
  emptyDescription?: string
  defaultViewMode?: 'grid' | 'list'
  gridColumns?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  error: null,
  totalPosts: 0,
  hasMore: false,
  showLink: true,
  showTags: true,
  showStats: true,
  showStatus: false,
  showActions: false,
  showContentPreview: true,
  showReadingTime: true,
  showLoadMore: false,
  showViewToggle: false,
  excerptLength: 150,
  emptyMessage: 'No posts found',
  emptyDescription: 'There are no posts to display at the moment.',
  defaultViewMode: 'grid',
  gridColumns: 3,
})

// Emits
const emit = defineEmits<{
  categoryClick: [category: string]
  tagClick: [tag: string]
  edit: [post: Post]
  publish: [post: Post]
  duplicate: [post: Post]
  delete: [post: Post]
  loadMore: []
  retry: []
}>()

// State
const viewMode = ref<'grid' | 'list'>(props.defaultViewMode)

// Computed
const gridColClass = computed(() => {
  const colMap = {
    1: 'col-12',
    2: 'col-md-6',
    3: 'col-lg-4',
    4: 'col-lg-3',
  }
  return colMap[props.gridColumns as keyof typeof colMap] || 'col-lg-4'
})

// Methods
const setViewMode = (mode: 'grid' | 'list') => {
  viewMode.value = mode
}

const truncateText = (text: string, length: number) => {
  if (text.length <= length) return text
  return text.substring(0, length).trim() + '...'
}

const stripHtml = (html: string) => {
  const tmp = document.createElement('div')
  tmp.innerHTML = html
  return tmp.textContent || tmp.innerText || ''
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 1) {
    return 'Yesterday'
  } else if (diffDays < 7) {
    return `${diffDays} days ago`
  } else {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }
}

const formatNumber = (num: number) => {
  if (num < 1000) return num.toString()
  if (num < 1000000) return (num / 1000).toFixed(1) + 'K'
  return (num / 1000000).toFixed(1) + 'M'
}

const estimateReadingTime = (content: string) => {
  const wordsPerMinute = 200
  const words = stripHtml(content).split(/\s+/).length
  return Math.ceil(words / wordsPerMinute)
}

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
</script>

<style scoped>
.post-list-item {
  padding: 1.5rem 0;
}

.post-list-item:not(.border-bottom) {
  padding-bottom: 0;
}

.post-list-item.border-bottom {
  border-bottom: 1px solid #e9ecef !important;
}

.badge {
  font-size: 0.7rem;
  font-weight: 500;
}

.badge:hover {
  opacity: 0.8;
}

.btn-group-sm .btn {
  font-size: 0.8rem;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}

/* Hover effects */
.post-list-item:hover {
  background-color: #f8f9fa;
  border-radius: 0.375rem;
  margin: 0 -1rem;
  padding-left: 1rem;
  padding-right: 1rem;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .post-list-item {
    padding: 1rem 0;
  }
  
  .btn-group-sm {
    display: none;
  }
  
  .post-list-view .row .col-md-2 {
    display: none;
  }
}
</style>
