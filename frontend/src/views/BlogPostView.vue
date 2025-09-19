<template>
  <div class="blog-post-view">
    <!-- Loading State -->
    <div v-if="postStore.loading" class="container py-5">
      <div class="text-center">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Loading post...</span>
        </div>
        <p class="mt-3 text-muted">Loading post...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="postStore.error" class="container py-5">
      <div class="row justify-content-center">
        <div class="col-md-6">
          <div class="alert alert-danger text-center">
            <i class="bi bi-exclamation-triangle-fill fs-1 mb-3"></i>
            <h4>Post Not Found</h4>
            <p>{{ postStore.error }}</p>
            <router-link to="/blog" class="btn btn-primary">
              <i class="bi bi-arrow-left"></i>
              Back to Blog
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Post Content -->
    <div v-else-if="post" class="post-content">
      <!-- Featured Image -->
      <div v-if="post.featured_image" class="featured-image-section">
        <div class="featured-image-wrapper">
          <img 
            :src="post.featured_image" 
            :alt="post.title"
            class="featured-image"
          >
          <div class="featured-image-overlay"></div>
        </div>
      </div>

      <!-- Post Header -->
      <div class="post-header py-5" :class="{ 'with-image': post.featured_image }">
        <div class="container">
          <div class="row justify-content-center">
            <div class="col-lg-8">
              <!-- Breadcrumb -->
              <nav class="mb-4">
                <ol class="breadcrumb">
                  <li class="breadcrumb-item">
                    <router-link to="/" class="text-decoration-none">Home</router-link>
                  </li>
                  <li class="breadcrumb-item">
                    <router-link to="/blog" class="text-decoration-none">Blog</router-link>
                  </li>
                  <li v-if="post.categories.length > 0" class="breadcrumb-item">
                    <button 
                      class="btn btn-link p-0 text-decoration-none"
                      @click="goToCategory(post.categories[0])"
                    >
                      {{ post.categories[0] }}
                    </button>
                  </li>
                  <li class="breadcrumb-item active">{{ post.title }}</li>
                </ol>
              </nav>

              <!-- Categories -->
              <div v-if="post.categories.length > 0" class="mb-3">
                <button
                  v-for="category in post.categories"
                  :key="category"
                  class="badge bg-primary me-2"
                  @click="goToCategory(category)"
                  style="cursor: pointer;"
                >
                  {{ category }}
                </button>
              </div>

              <!-- Title -->
              <h1 class="display-5 fw-bold mb-4" :class="{ 'text-white': post.featured_image }">
                {{ post.title }}
              </h1>

              <!-- Excerpt -->
              <p v-if="post.excerpt" class="lead mb-4" :class="{ 'text-white-50': post.featured_image }">
                {{ post.excerpt }}
              </p>

              <!-- Meta Information -->
              <div class="post-meta d-flex flex-wrap align-items-center gap-4 mb-4">
                <div class="d-flex align-items-center" :class="{ 'text-white-50': post.featured_image }">
                  <i class="bi bi-calendar3 me-2"></i>
                  <span>{{ formatDate(post.published_at || post.created_at) }}</span>
                </div>
                <div class="d-flex align-items-center" :class="{ 'text-white-50': post.featured_image }">
                  <i class="bi bi-clock me-2"></i>
                  <span>{{ estimateReadingTime(post.content) }} min read</span>
                </div>
                <div class="d-flex align-items-center" :class="{ 'text-white-50': post.featured_image }">
                  <i class="bi bi-eye me-2"></i>
                  <span>{{ formatNumber(post.view_count) }} views</span>
                </div>
                <div v-if="post.comment_count > 0" class="d-flex align-items-center" :class="{ 'text-white-50': post.featured_image }">
                  <i class="bi bi-chat me-2"></i>
                  <span>{{ formatNumber(post.comment_count) }} comments</span>
                </div>
              </div>

              <!-- Share Buttons -->
              <div class="share-buttons d-flex gap-2">
                <button 
                  class="btn btn-outline-primary btn-sm"
                  :class="{ 'btn-outline-light': post.featured_image }"
                  @click="shareOnTwitter"
                >
                  <i class="bi bi-twitter"></i>
                  Share
                </button>
                <button 
                  class="btn btn-outline-primary btn-sm"
                  :class="{ 'btn-outline-light': post.featured_image }"
                  @click="copyLink"
                >
                  <i class="bi bi-link-45deg"></i>
                  Copy Link
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Post Body -->
      <div class="post-body py-5">
        <div class="container">
          <div class="row justify-content-center">
            <div class="col-lg-8">
              <!-- Content -->
              <div class="post-content-wrapper">
                <div v-html="formattedContent" class="post-content-html"></div>
              </div>

              <!-- Tags -->
              <div v-if="post.tags.length > 0" class="tags-section mt-5 pt-4 border-top">
                <h6 class="mb-3">Tags:</h6>
                <div class="d-flex flex-wrap gap-2">
                  <button
                    v-for="tag in post.tags"
                    :key="tag"
                    class="btn btn-outline-secondary btn-sm"
                    @click="goToTag(tag)"
                  >
                    #{{ tag }}
                  </button>
                </div>
              </div>

              <!-- Share Section -->
              <div class="share-section mt-5 pt-4 border-top">
                <h6 class="mb-3">Share this post:</h6>
                <div class="d-flex gap-2">
                  <button class="btn btn-primary btn-sm" @click="shareOnTwitter">
                    <i class="bi bi-twitter"></i>
                    Twitter
                  </button>
                  <button class="btn btn-primary btn-sm" @click="shareOnLinkedIn">
                    <i class="bi bi-linkedin"></i>
                    LinkedIn
                  </button>
                  <button class="btn btn-secondary btn-sm" @click="copyLink">
                    <i class="bi bi-link-45deg"></i>
                    Copy Link
                  </button>
                </div>
              </div>

              <!-- Navigation -->
              <div class="post-navigation mt-5 pt-4 border-top">
                <div class="row">
                  <div class="col-md-6">
                    <div v-if="previousPost" class="d-flex align-items-center">
                      <i class="bi bi-arrow-left me-3 text-muted"></i>
                      <div>
                        <div class="small text-muted">Previous Post</div>
                        <router-link 
                          :to="`/blog/${previousPost.slug}`"
                          class="text-decoration-none fw-semibold"
                        >
                          {{ previousPost.title }}
                        </router-link>
                      </div>
                    </div>
                  </div>
                  <div class="col-md-6 text-md-end mt-3 mt-md-0">
                    <div v-if="nextPost" class="d-flex align-items-center justify-content-md-end">
                      <div class="text-md-end me-md-3">
                        <div class="small text-muted">Next Post</div>
                        <router-link 
                          :to="`/blog/${nextPost.slug}`"
                          class="text-decoration-none fw-semibold"
                        >
                          {{ nextPost.title }}
                        </router-link>
                      </div>
                      <i class="bi bi-arrow-right ms-3 text-muted"></i>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Related Posts -->
              <div v-if="relatedPosts.length > 0" class="related-posts mt-5 pt-4 border-top">
                <h4 class="mb-4">Related Posts</h4>
                <div class="row g-4">
                  <div 
                    v-for="relatedPost in relatedPosts" 
                    :key="relatedPost.id"
                    class="col-md-6"
                  >
                    <div class="card h-100">
                      <div v-if="relatedPost.featured_image" class="card-img-wrapper">
                        <img 
                          :src="relatedPost.featured_image" 
                          :alt="relatedPost.title"
                          class="card-img-top"
                        >
                      </div>
                      <div class="card-body">
                        <h6 class="card-title">
                          <router-link 
                            :to="`/blog/${relatedPost.slug}`"
                            class="text-decoration-none text-dark"
                          >
                            {{ relatedPost.title }}
                          </router-link>
                        </h6>
                        <p v-if="relatedPost.excerpt" class="card-text text-muted small">
                          {{ truncateText(relatedPost.excerpt, 100) }}
                        </p>
                        <div class="small text-muted">
                          {{ formatDate(relatedPost.published_at || relatedPost.created_at) }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Link Copied Toast -->
    <div 
      v-if="showToast"
      class="toast-container position-fixed bottom-0 end-0 p-3"
    >
      <div class="toast show">
        <div class="toast-body">
          <i class="bi bi-check-circle-fill text-success me-2"></i>
          Link copied to clipboard!
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePostStore, type Post } from '@/stores/postStore'

// Stores
const postStore = usePostStore()
const route = useRoute()
const router = useRouter()

// State
const showToast = ref(false)
const previousPost = ref<Post | null>(null)
const nextPost = ref<Post | null>(null)
const relatedPosts = ref<Post[]>([])

// Computed
const post = computed(() => postStore.currentPost)

const formattedContent = computed(() => {
  if (!post.value?.content) return ''
  
  // Basic HTML formatting (in a real app, you'd use a proper markdown parser)
  return post.value.content
    .replace(/\n\n/g, '</p><p>')
    .replace(/\n/g, '<br>')
    .replace(/^/, '<p>')
    .replace(/$/, '</p>')
})

// Watchers
watch(
  () => route.params.slug,
  async (newSlug) => {
    if (newSlug && typeof newSlug === 'string') {
      await loadPost(newSlug)
    }
  }
)

// Lifecycle
onMounted(async () => {
  const slug = route.params.slug as string
  if (slug) {
    await loadPost(slug)
  }
})

// Methods
const loadPost = async (slug: string) => {
  const loadedPost = await postStore.fetchPublishedPostBySlug(slug)
  
  if (loadedPost) {
    // Update page title and meta tags
    document.title = `${loadedPost.title} | My Blog`
    
    // Update meta description
    const metaDescription = document.querySelector('meta[name="description"]')
    if (metaDescription) {
      metaDescription.setAttribute('content', loadedPost.meta_description || loadedPost.excerpt || 'Blog post')
    }

    // Load related content (mock implementation)
    await loadRelatedContent(loadedPost)
  }
}

const loadRelatedContent = async (currentPost: Post) => {
  // In a real app, you'd call APIs to get related posts
  // For now, we'll use mock data
  relatedPosts.value = []
  
  // Mock previous/next posts
  previousPost.value = null
  nextPost.value = null
}

const goToCategory = (category: string) => {
  router.push({ path: '/blog', query: { category } })
}

const goToTag = (tag: string) => {
  router.push({ path: '/blog', query: { tag } })
}

const shareOnTwitter = () => {
  const url = encodeURIComponent(window.location.href)
  const text = encodeURIComponent(`Check out this post: ${post.value?.title}`)
  const twitterUrl = `https://twitter.com/intent/tweet?url=${url}&text=${text}`
  window.open(twitterUrl, '_blank', 'width=600,height=400')
}

const shareOnLinkedIn = () => {
  const url = encodeURIComponent(window.location.href)
  const linkedInUrl = `https://www.linkedin.com/sharing/share-offsite/?url=${url}`
  window.open(linkedInUrl, '_blank', 'width=600,height=400')
}

const copyLink = async () => {
  try {
    await navigator.clipboard.writeText(window.location.href)
    showToast.value = true
    setTimeout(() => {
      showToast.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to copy link:', err)
    // Fallback for older browsers
    const textArea = document.createElement('textarea')
    textArea.value = window.location.href
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showToast.value = true
    setTimeout(() => {
      showToast.value = false
    }, 3000)
  }
}

// Utility functions
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const formatNumber = (num: number) => {
  if (num < 1000) return num.toString()
  if (num < 1000000) return (num / 1000).toFixed(1) + 'K'
  return (num / 1000000).toFixed(1) + 'M'
}

const estimateReadingTime = (content: string) => {
  const wordsPerMinute = 200
  const words = content.replace(/<[^>]*>/g, '').split(/\s+/).length
  return Math.ceil(words / wordsPerMinute)
}

const truncateText = (text: string, length: number) => {
  if (text.length <= length) return text
  return text.substring(0, length).trim() + '...'
}
</script>

<style scoped>
.featured-image-section {
  position: relative;
  height: 60vh;
  min-height: 400px;
  overflow: hidden;
}

.featured-image-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.featured-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.featured-image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.7) 100%);
}

.post-header.with-image {
  position: relative;
  margin-top: -60vh;
  padding-top: 20vh;
  z-index: 10;
}

.post-content-html {
  font-size: 1.1rem;
  line-height: 1.8;
  color: #333;
}

.post-content-html p {
  margin-bottom: 1.5rem;
}

.post-content-html h2,
.post-content-html h3,
.post-content-html h4 {
  margin-top: 2rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.post-content-html h2 {
  font-size: 1.5rem;
  border-bottom: 2px solid #e9ecef;
  padding-bottom: 0.5rem;
}

.post-content-html h3 {
  font-size: 1.3rem;
}

.post-content-html h4 {
  font-size: 1.1rem;
}

.post-content-html blockquote {
  border-left: 4px solid #0d6efd;
  margin: 2rem 0;
  padding: 1rem 1.5rem;
  background-color: #f8f9fa;
  font-style: italic;
}

.post-content-html code {
  background-color: #f8f9fa;
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.9em;
}

.post-content-html pre {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1.5rem 0;
}

.post-content-html pre code {
  background: none;
  padding: 0;
}

.post-content-html img {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 1.5rem 0;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.post-content-html ul,
.post-content-html ol {
  margin-bottom: 1.5rem;
  padding-left: 2rem;
}

.post-content-html li {
  margin-bottom: 0.5rem;
}

.breadcrumb {
  background: none;
  padding: 0;
}

.breadcrumb-item + .breadcrumb-item::before {
  content: "›";
  color: #6c757d;
}

.card-img-wrapper {
  height: 150px;
  overflow: hidden;
}

.card-img-top {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.card:hover .card-img-top {
  transform: scale(1.05);
}

.toast-container {
  z-index: 1060;
}

.toast.show {
  display: block;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .featured-image-section {
    height: 40vh;
    min-height: 300px;
  }
  
  .post-header.with-image {
    margin-top: -40vh;
    padding-top: 15vh;
  }
  
  .display-5 {
    font-size: 1.8rem !important;
  }
  
  .post-content-html {
    font-size: 1rem;
  }
  
  .post-meta {
    flex-direction: column;
    align-items: flex-start !important;
    gap: 0.5rem !important;
  }
}

@media (max-width: 576px) {
  .share-buttons {
    flex-direction: column;
  }
  
  .share-buttons .btn {
    width: 100%;
  }
}
</style>
