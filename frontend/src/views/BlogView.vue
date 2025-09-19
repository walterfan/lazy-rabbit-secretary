<template>
  <div class="blog-view">
    <!-- Hero Section -->
    <div class="hero-section bg-primary text-white py-5 mb-5">
      <div class="container">
        <div class="row align-items-center">
          <div class="col-lg-8">
            <h1 class="display-4 fw-bold mb-3">Welcome to My Blog</h1>
            <p class="lead mb-4">
              Thoughts, insights, and stories about technology, development, and life.
            </p>
            <div class="d-flex gap-3">
              <button class="btn btn-light btn-lg" @click="scrollToContent">
                <i class="bi bi-arrow-down"></i>
                Start Reading
              </button>
              <router-link to="/blog/about" class="btn btn-outline-light btn-lg">
                About Me
              </router-link>
            </div>
          </div>
          <div class="col-lg-4 text-center">
            <i class="bi bi-journal-text" style="font-size: 8rem; opacity: 0.3;"></i>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="row">
        <!-- Main Content -->
        <div class="col-lg-8">
          <!-- Search and Filters -->
          <div class="card mb-4">
            <div class="card-body">
              <div class="row g-3">
                <div class="col-md-6">
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
                  <select v-model="selectedCategory" class="form-select" @change="handleCategoryFilter">
                    <option value="">All Categories</option>
                    <option v-for="category in categories" :key="category" :value="category">
                      {{ category }}
                    </option>
                  </select>
                </div>
                <div class="col-md-3">
                  <button class="btn btn-outline-secondary w-100" @click="clearFilters">
                    Clear Filters
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Sticky Posts -->
          <div v-if="postStore.stickyPosts.length > 0 && !searchQuery && !selectedCategory" class="mb-5">
            <h3 class="h4 mb-3 d-flex align-items-center">
              <i class="bi bi-pin-fill text-warning me-2"></i>
              Featured Posts
            </h3>
            <PostList
              :posts="postStore.stickyPosts"
              :loading="false"
              :show-actions="false"
              :show-status="false"
              :grid-columns="2"
              :excerpt-length="120"
            />
          </div>

          <!-- Recent Posts -->
          <div id="content">
            <div class="d-flex justify-content-between align-items-center mb-4">
              <h2 class="h3 mb-0">
                {{ getPostsTitle() }}
              </h2>
              <div class="text-muted small">
                {{ postStore.totalPosts }} {{ postStore.totalPosts === 1 ? 'post' : 'posts' }}
              </div>
            </div>

            <PostList
              :posts="postStore.publishedPosts"
              :loading="postStore.loading"
              :error="postStore.error"
              :total-posts="postStore.totalPosts"
              :has-more="postStore.currentPage < postStore.totalPages"
              :show-actions="false"
              :show-status="false"
              :show-load-more="true"
              :show-view-toggle="true"
              :empty-message="getEmptyMessage()"
              :empty-description="getEmptyDescription()"
              @category-click="handleCategoryClick"
              @tag-click="handleTagClick"
              @load-more="loadMorePosts"
              @retry="refreshPosts"
            />
          </div>

          <!-- Pagination -->
          <nav v-if="postStore.totalPages > 1" class="mt-5">
            <ul class="pagination justify-content-center">
              <li class="page-item" :class="{ disabled: postStore.currentPage <= 1 }">
                <button 
                  class="page-link" 
                  @click="changePage(postStore.currentPage - 1)"
                  :disabled="postStore.currentPage <= 1"
                >
                  <i class="bi bi-chevron-left"></i>
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
                  <i class="bi bi-chevron-right"></i>
                </button>
              </li>
            </ul>
          </nav>
        </div>

        <!-- Sidebar -->
        <div class="col-lg-4">
          <!-- About Widget -->
          <div class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">About</h5>
            </div>
            <div class="card-body">
              <p class="card-text">
                Welcome to my personal blog where I share thoughts about technology, 
                software development, and life experiences.
              </p>
              <router-link to="/blog/about" class="btn btn-primary btn-sm">
                Learn More
              </router-link>
            </div>
          </div>

          <!-- Popular Posts -->
          <div v-if="postStore.popularPosts.length > 0" class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="bi bi-fire me-2"></i>
                Popular Posts
              </h5>
            </div>
            <div class="card-body p-0">
              <div 
                v-for="(post, index) in postStore.popularPosts" 
                :key="post.id"
                class="d-flex p-3"
                :class="{ 'border-bottom': index < postStore.popularPosts.length - 1 }"
              >
                <div v-if="post.featured_image" class="me-3 flex-shrink-0">
                  <img 
                    :src="post.featured_image" 
                    :alt="post.title"
                    class="rounded"
                    style="width: 60px; height: 60px; object-fit: cover;"
                  >
                </div>
                <div class="flex-grow-1">
                  <h6 class="mb-1">
                    <router-link 
                      :to="`/blog/${post.slug}`"
                      class="text-decoration-none text-dark"
                    >
                      {{ post.title }}
                    </router-link>
                  </h6>
                  <small class="text-muted">
                    {{ formatNumber(post.view_count) }} views
                  </small>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Posts -->
          <div v-if="postStore.recentPosts.length > 0" class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="bi bi-clock me-2"></i>
                Recent Posts
              </h5>
            </div>
            <div class="card-body p-0">
              <div 
                v-for="(post, index) in postStore.recentPosts" 
                :key="post.id"
                class="d-flex p-3"
                :class="{ 'border-bottom': index < postStore.recentPosts.length - 1 }"
              >
                <div v-if="post.featured_image" class="me-3 flex-shrink-0">
                  <img 
                    :src="post.featured_image" 
                    :alt="post.title"
                    class="rounded"
                    style="width: 60px; height: 60px; object-fit: cover;"
                  >
                </div>
                <div class="flex-grow-1">
                  <h6 class="mb-1">
                    <router-link 
                      :to="`/blog/${post.slug}`"
                      class="text-decoration-none text-dark"
                    >
                      {{ post.title }}
                    </router-link>
                  </h6>
                  <small class="text-muted">
                    {{ formatDate(post.published_at || post.created_at) }}
                  </small>
                </div>
              </div>
            </div>
          </div>

          <!-- Categories -->
          <div v-if="categories.length > 0" class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="bi bi-folder me-2"></i>
                Categories
              </h5>
            </div>
            <div class="card-body">
              <div class="d-flex flex-wrap gap-2">
                <button
                  v-for="category in categories"
                  :key="category"
                  class="btn btn-outline-primary btn-sm"
                  @click="handleCategoryClick(category)"
                >
                  {{ category }}
                </button>
              </div>
            </div>
          </div>

          <!-- Tags -->
          <div v-if="tags.length > 0" class="card mb-4">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="bi bi-tags me-2"></i>
                Tags
              </h5>
            </div>
            <div class="card-body">
              <div class="d-flex flex-wrap gap-1">
                <button
                  v-for="tag in tags"
                  :key="tag"
                  class="btn btn-light btn-sm"
                  @click="handleTagClick(tag)"
                  style="font-size: 0.8rem;"
                >
                  #{{ tag }}
                </button>
              </div>
            </div>
          </div>

          <!-- Archive -->
          <div class="card">
            <div class="card-header">
              <h5 class="mb-0">
                <i class="bi bi-calendar3 me-2"></i>
                Archive
              </h5>
            </div>
            <div class="card-body">
              <div class="list-group list-group-flush">
                <router-link
                  v-for="month in archiveMonths"
                  :key="month.link"
                  :to="month.link"
                  class="list-group-item list-group-item-action border-0 px-0"
                >
                  {{ month.label }}
                  <span class="badge bg-secondary float-end">{{ month.count }}</span>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePostStore } from '@/stores/postStore'
import PostList from '@/components/posts/PostList.vue'

// Stores
const postStore = usePostStore()
const route = useRoute()
const router = useRouter()

// State
const searchQuery = ref('')
const selectedCategory = ref('')

// Mock data for sidebar (in real app, these would come from API)
const categories = ref(['Technology', 'Web Development', 'Vue.js', 'JavaScript', 'Tutorial'])
const tags = ref(['vue', 'javascript', 'typescript', 'nodejs', 'api', 'frontend', 'backend'])
const archiveMonths = ref([
  { label: 'December 2024', link: '/blog/archive/2024/12', count: 5 },
  { label: 'November 2024', link: '/blog/archive/2024/11', count: 8 },
  { label: 'October 2024', link: '/blog/archive/2024/10', count: 12 },
  { label: 'September 2024', link: '/blog/archive/2024/9', count: 6 },
])

// Lifecycle
onMounted(async () => {
  // Load initial data
  await Promise.all([
    postStore.fetchPublishedPosts(1),
    postStore.fetchPopularPosts(5),
    postStore.fetchRecentPosts(5),
    postStore.fetchStickyPosts(),
  ])

  // Handle URL parameters
  const category = route.query.category as string
  const tag = route.query.tag as string
  const search = route.query.q as string

  if (category) {
    selectedCategory.value = category
    await postStore.fetchPostsByCategory(category)
  } else if (tag) {
    await postStore.fetchPostsByTag(tag)
  } else if (search) {
    searchQuery.value = search
    await postStore.searchPublishedPosts(search)
  }
})

// Methods
const scrollToContent = () => {
  const element = document.getElementById('content')
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

const handleSearch = async () => {
  if (searchQuery.value.trim()) {
    await postStore.searchPublishedPosts(searchQuery.value)
    updateUrl({ q: searchQuery.value })
  } else {
    await postStore.fetchPublishedPosts()
    updateUrl({})
  }
}

const handleCategoryFilter = async () => {
  if (selectedCategory.value) {
    await postStore.fetchPostsByCategory(selectedCategory.value)
    updateUrl({ category: selectedCategory.value })
  } else {
    await postStore.fetchPublishedPosts()
    updateUrl({})
  }
}

const handleCategoryClick = async (category: string) => {
  selectedCategory.value = category
  searchQuery.value = ''
  await postStore.fetchPostsByCategory(category)
  updateUrl({ category })
}

const handleTagClick = async (tag: string) => {
  selectedCategory.value = ''
  searchQuery.value = ''
  await postStore.fetchPostsByTag(tag)
  updateUrl({ tag })
}

const clearFilters = async () => {
  searchQuery.value = ''
  selectedCategory.value = ''
  await postStore.fetchPublishedPosts()
  updateUrl({})
}

const refreshPosts = async () => {
  if (searchQuery.value) {
    await postStore.searchPublishedPosts(searchQuery.value)
  } else if (selectedCategory.value) {
    await postStore.fetchPostsByCategory(selectedCategory.value)
  } else {
    await postStore.fetchPublishedPosts()
  }
}

const loadMorePosts = async () => {
  const nextPage = postStore.currentPage + 1
  if (searchQuery.value) {
    await postStore.searchPublishedPosts(searchQuery.value, nextPage)
  } else if (selectedCategory.value) {
    await postStore.fetchPostsByCategory(selectedCategory.value, nextPage)
  } else {
    await postStore.fetchPublishedPosts(nextPage)
  }
}

const changePage = async (page: number) => {
  if (page >= 1 && page <= postStore.totalPages) {
    if (searchQuery.value) {
      await postStore.searchPublishedPosts(searchQuery.value, page)
    } else if (selectedCategory.value) {
      await postStore.fetchPostsByCategory(selectedCategory.value, page)
    } else {
      await postStore.fetchPublishedPosts(page)
    }
    scrollToContent()
  }
}

const updateUrl = (params: Record<string, string>) => {
  router.push({ query: params })
}

// Computed
const getPostsTitle = () => {
  if (searchQuery.value) {
    return `Search Results for "${searchQuery.value}"`
  }
  if (selectedCategory.value) {
    return `Posts in "${selectedCategory.value}"`
  }
  return 'Latest Posts'
}

const getEmptyMessage = () => {
  if (searchQuery.value) {
    return 'No posts found'
  }
  if (selectedCategory.value) {
    return 'No posts in this category'
  }
  return 'No posts available'
}

const getEmptyDescription = () => {
  if (searchQuery.value) {
    return `No posts match your search for "${searchQuery.value}". Try different keywords.`
  }
  if (selectedCategory.value) {
    return `There are no posts in the "${selectedCategory.value}" category yet.`
  }
  return 'Check back later for new posts!'
}

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
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const formatNumber = (num: number) => {
  if (num < 1000) return num.toString()
  if (num < 1000000) return (num / 1000).toFixed(1) + 'K'
  return (num / 1000000).toFixed(1) + 'M'
}
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.blog-view {
  min-height: 100vh;
}

.card {
  border: 1px solid #e9ecef;
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
}

.card-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.list-group-item-action:hover {
  background-color: #f8f9fa;
}

.btn-outline-primary:hover {
  transform: translateY(-1px);
}

.pagination .page-link {
  color: #6c757d;
  border-color: #dee2e6;
}

.pagination .page-item.active .page-link {
  background-color: #0d6efd;
  border-color: #0d6efd;
}

/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero-section {
    padding: 3rem 0 !important;
  }
  
  .display-4 {
    font-size: 2rem !important;
  }
  
  .lead {
    font-size: 1rem !important;
  }
  
  .btn-lg {
    padding: 0.5rem 1rem !important;
    font-size: 1rem !important;
  }
}
</style>
