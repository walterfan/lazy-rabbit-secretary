<template>
  <article class="post-card card h-100">
    <!-- Featured Image -->
    <div v-if="post.featured_image" class="card-img-top-wrapper">
      <img 
        :src="post.featured_image" 
        :alt="post.title"
        class="card-img-top"
        @error="onImageError"
      >
      <div v-if="post.is_sticky" class="sticky-badge">
        <i class="bi bi-pin-fill"></i>
        Featured
      </div>
    </div>

    <div class="card-body d-flex flex-column">
      <!-- Categories -->
      <div v-if="post.categories.length > 0" class="mb-2">
        <span 
          v-for="category in post.categories.slice(0, 2)" 
          :key="category"
          class="badge bg-primary me-1"
          @click="$emit('categoryClick', category)"
          style="cursor: pointer;"
        >
          {{ category }}
        </span>
        <span v-if="post.categories.length > 2" class="text-muted small">
          +{{ post.categories.length - 2 }} more
        </span>
      </div>

      <!-- Title -->
      <h5 class="card-title">
        <router-link 
          v-if="showLink"
          :to="getPostLink()"
          class="text-decoration-none text-dark"
        >
          {{ post.title }}
        </router-link>
        <span v-else>{{ post.title }}</span>
      </h5>

      <!-- Excerpt -->
      <p v-if="post.excerpt" class="card-text text-muted">
        {{ truncateText(post.excerpt, excerptLength) }}
      </p>

      <!-- Content preview if no excerpt -->
      <p v-else-if="showContentPreview" class="card-text text-muted">
        {{ truncateText(stripHtml(post.content), excerptLength) }}
      </p>

      <!-- Tags -->
      <div v-if="post.tags.length > 0 && showTags" class="mb-2">
        <span 
          v-for="tag in post.tags.slice(0, 3)" 
          :key="tag"
          class="badge bg-light text-dark me-1"
          @click="$emit('tagClick', tag)"
          style="cursor: pointer; font-size: 0.7rem;"
        >
          #{{ tag }}
        </span>
        <span v-if="post.tags.length > 3" class="text-muted small">
          +{{ post.tags.length - 3 }}
        </span>
      </div>

      <!-- Spacer -->
      <div class="flex-grow-1"></div>

      <!-- Meta Information -->
      <div class="post-meta mt-auto">
        <div class="d-flex justify-content-between align-items-center text-muted small">
          <div class="d-flex align-items-center">
            <i class="bi bi-calendar3 me-1"></i>
            <span>{{ formatDate(post.published_at || post.created_at) }}</span>
          </div>
          
          <div class="d-flex align-items-center gap-3">
            <!-- View count -->
            <div v-if="showStats" class="d-flex align-items-center">
              <i class="bi bi-eye me-1"></i>
              <span>{{ formatNumber(post.view_count) }}</span>
            </div>
            
            <!-- Comment count -->
            <div v-if="showStats && post.comment_count > 0" class="d-flex align-items-center">
              <i class="bi bi-chat me-1"></i>
              <span>{{ formatNumber(post.comment_count) }}</span>
            </div>

            <!-- Reading time estimate -->
            <div v-if="showReadingTime" class="d-flex align-items-center">
              <i class="bi bi-clock me-1"></i>
              <span>{{ estimateReadingTime(post.content) }} min</span>
            </div>
          </div>
        </div>

        <!-- Status badge (for admin view) -->
        <div v-if="showStatus" class="mt-2">
          <span 
            class="badge"
            :class="getStatusBadgeClass(post.status)"
          >
            {{ formatStatus(post.status) }}
          </span>
          <span v-if="post.is_sticky" class="badge bg-warning text-dark ms-1">
            <i class="bi bi-pin-fill"></i> Sticky
          </span>
        </div>
      </div>
    </div>

    <!-- Action buttons (for admin view) -->
    <div v-if="showActions" class="card-footer bg-transparent">
      <div class="d-flex gap-2">
        <button 
          class="btn btn-sm btn-outline-primary"
          @click="$emit('edit', post)"
        >
          <i class="bi bi-pencil"></i>
          Edit
        </button>
        
        <button 
          v-if="post.status !== 'published'"
          class="btn btn-sm btn-outline-success"
          @click="$emit('publish', post)"
        >
          <i class="bi bi-globe"></i>
          Publish
        </button>
        
        <a 
          v-if="post.status === 'published'"
          :href="getPostLink()"
          target="_blank"
          class="btn btn-sm btn-outline-info"
        >
          <i class="bi bi-eye"></i>
          View
        </a>
        
        <div class="dropdown ms-auto">
          <button 
            class="btn btn-sm btn-outline-secondary dropdown-toggle"
            data-bs-toggle="dropdown"
          >
            More
          </button>
          <ul class="dropdown-menu">
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
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Post } from '@/stores/postStore'

// Props
interface Props {
  post: Post
  showLink?: boolean
  showTags?: boolean
  showStats?: boolean
  showStatus?: boolean
  showActions?: boolean
  showContentPreview?: boolean
  showReadingTime?: boolean
  excerptLength?: number
}

const props = withDefaults(defineProps<Props>(), {
  showLink: true,
  showTags: true,
  showStats: true,
  showStatus: false,
  showActions: false,
  showContentPreview: true,
  showReadingTime: true,
  excerptLength: 150,
})

// Emits
const emit = defineEmits<{
  categoryClick: [category: string]
  tagClick: [tag: string]
  edit: [post: Post]
  publish: [post: Post]
  duplicate: [post: Post]
  delete: [post: Post]
}>()

// Methods
const getPostLink = () => {
  return `/blog/${props.post.slug}`
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
  } else if (diffDays < 30) {
    const weeks = Math.floor(diffDays / 7)
    return `${weeks} week${weeks > 1 ? 's' : ''} ago`
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

const onImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
}
</script>

<style scoped>
.post-card {
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
  border: 1px solid #e9ecef;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-img-top-wrapper {
  position: relative;
  overflow: hidden;
}

.card-img-top {
  height: 200px;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.post-card:hover .card-img-top {
  transform: scale(1.05);
}

.sticky-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(255, 193, 7, 0.9);
  color: #000;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 600;
  line-height: 1.3;
  margin-bottom: 0.75rem;
}

.card-title a:hover {
  color: #0d6efd !important;
}

.card-text {
  font-size: 0.9rem;
  line-height: 1.5;
}

.post-meta {
  font-size: 0.8rem;
  border-top: 1px solid #f8f9fa;
  padding-top: 0.75rem;
}

.badge {
  font-size: 0.7rem;
  font-weight: 500;
}

.badge:hover {
  opacity: 0.8;
}

.dropdown-toggle::after {
  font-size: 0.7rem;
}

.card-footer {
  border-top: 1px solid #f8f9fa;
  background: transparent !important;
}

.btn-sm {
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .card-img-top {
    height: 150px;
  }
  
  .card-title {
    font-size: 1rem;
  }
  
  .card-text {
    font-size: 0.85rem;
  }
}
</style>
