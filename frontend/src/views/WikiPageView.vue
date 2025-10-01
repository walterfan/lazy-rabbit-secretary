<template>
  <div class="wiki-page-view">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">{{ $t('common.loading') }}</span>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <div class="alert" :class="isPageNotFound ? 'alert-info' : 'alert-danger'">
        <div class="text-center">
          <i class="bi" :class="isPageNotFound ? 'bi-file-earmark-plus display-1 text-primary mb-3' : 'bi-exclamation-triangle display-1 text-danger mb-3'"></i>
          <h4>{{ $t('wiki.pageNotFound') }}</h4>
          <p class="mb-4">{{ isPageNotFound ? $t('wiki.pageNotFoundDesc', { slug: route.params.slug }) : error }}</p>
          
          <div v-if="isPageNotFound && isAuthenticated" class="d-flex justify-content-center gap-3">
            <button class="btn btn-primary btn-lg" @click="createPage">
              <i class="bi bi-plus-lg me-2"></i>
              {{ $t('wiki.createThisPage') }}
            </button>
            <button class="btn btn-outline-secondary btn-lg" @click="goHome">
              <i class="bi bi-house me-2"></i>
              {{ $t('wiki.backToHome') }}
            </button>
          </div>
          
          <div v-else-if="isPageNotFound && !isAuthenticated" class="d-flex justify-content-center gap-3">
            <button class="btn btn-outline-primary btn-lg" @click="goToSignIn">
              <i class="bi bi-box-arrow-in-right me-2"></i>
              {{ $t('nav.signIn') }}
            </button>
            <button class="btn btn-outline-secondary btn-lg" @click="goHome">
              <i class="bi bi-house me-2"></i>
              {{ $t('wiki.backToHome') }}
            </button>
          </div>
          
          <div v-else class="d-flex justify-content-center gap-3">
            <button class="btn btn-outline-secondary btn-lg" @click="goHome">
              <i class="bi bi-house me-2"></i>
              {{ $t('wiki.backToHome') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Page Content -->
    <div v-else-if="page" class="container">
      <!-- Breadcrumb Navigation -->
      <WikiNavigation :page="page" />

      <!-- Page Header -->
      <div class="page-header">
        <div class="row">
          <div class="col-md-8">
            <h1 class="page-title">{{ page.title }}</h1>
            <div class="page-meta">
              <span class="meta-item">
                <i class="bi bi-person"></i>
                {{ $t('wiki.lastEditedBy', { user: page.updated_by }) }}
              </span>
              <span class="meta-item">
                <i class="bi bi-clock"></i>
                {{ formatDate(page.updated_at) }}
              </span>
              <span class="meta-item">
                <i class="bi bi-eye"></i>
                {{ $t('wiki.viewCount', { count: page.view_count }) }}
              </span>
              <span class="meta-item">
                <i class="bi bi-pencil"></i>
                {{ $t('wiki.editCount', { count: page.edit_count }) }}
              </span>
            </div>
          </div>
          <div class="col-md-4 text-end">
            <WikiPageActions
              :page="page"
              :canEdit="page.can_edit"
              :canDelete="page.can_delete"
              @edit="editPage"
              @delete="deletePage"
              @history="viewHistory"
              @lock="lockPage"
              @unlock="unlockPage"
            />
          </div>
        </div>
      </div>

      <!-- Page Content -->
      <div class="row">
        <div class="col-lg-8">
          <div class="page-content">
            <!-- Redirect Notice -->
            <div v-if="page.type === 'redirect' && page.redirect_to" class="alert alert-info">
              <i class="bi bi-arrow-right"></i>
              {{ $t('wiki.redirectNotice', { target: page.redirect_to }) }}
            </div>

            <!-- Stub Notice -->
            <div v-if="page.type === 'stub'" class="alert alert-warning">
              <i class="bi bi-exclamation-triangle"></i>
              {{ $t('wiki.stubNotice') }}
            </div>

            <!-- Protected Page Notice -->
            <div v-if="page.is_protected && !isAuthenticated" class="alert alert-warning">
              <i class="bi bi-shield-lock"></i>
              {{ $t('wiki.protectedPageNotice') }}
            </div>

            <!-- Locked Page Notice -->
            <div v-if="page.is_locked" class="alert alert-danger">
              <i class="bi bi-lock"></i>
              {{ $t('wiki.lockedPageNotice') }}
            </div>

            <!-- Markdown Content -->
            <div class="markdown-content" v-html="renderedContent"></div>

            <!-- Categories and Tags -->
            <div v-if="page.categories.length > 0 || page.tags.length > 0" class="page-footer">
              <div class="categories-tags">
                <div v-if="page.categories.length > 0" class="categories">
                  <h6>{{ $t('wiki.categories') }}:</h6>
                  <div class="tag-list">
                    <span 
                      v-for="category in page.categories" 
                      :key="category"
                      class="badge bg-primary me-1 mb-1"
                    >
                      {{ category }}
                    </span>
                  </div>
                </div>
                
                <div v-if="page.tags.length > 0" class="tags">
                  <h6>{{ $t('wiki.tags') }}:</h6>
                  <div class="tag-list">
                    <span 
                      v-for="tag in page.tags" 
                      :key="tag"
                      class="badge bg-secondary me-1 mb-1"
                    >
                      {{ tag }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="col-lg-4">
          <div class="page-sidebar">
            <!-- Page Info -->
            <div class="sidebar-section">
              <h5>{{ $t('wiki.pageInfo') }}</h5>
              <div class="info-list">
                <div class="info-item">
                  <span class="info-label">{{ $t('wiki.status') }}:</span>
                  <span class="badge" :class="getStatusBadgeClass(page.status)">
                    {{ $t(`wiki.status.${page.status}`) }}
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">{{ $t('wiki.type') }}:</span>
                  <span class="badge bg-info">
                    {{ $t(`wiki.type.${page.type}`) }}
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">{{ $t('wiki.language') }}:</span>
                  <span>{{ page.language }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">{{ $t('wiki.version') }}:</span>
                  <span>{{ page.current_version }}</span>
                </div>
              </div>
            </div>

            <!-- Related Pages -->
            <div v-if="relatedPages.length > 0" class="sidebar-section">
              <h5>{{ $t('wiki.relatedPages') }}</h5>
              <div class="related-pages">
                <div 
                  v-for="relatedPage in relatedPages" 
                  :key="relatedPage.id"
                  class="related-page-item"
                  @click="goToPage(relatedPage.slug)"
                >
                  <div class="related-page-title">{{ relatedPage.title }}</div>
                  <div class="related-page-summary">{{ relatedPage.summary }}</div>
                </div>
              </div>
            </div>

            <!-- Recent Changes -->
            <div class="sidebar-section">
              <h5>{{ $t('wiki.recentChanges') }}</h5>
              <div class="recent-changes">
                <div 
                  v-for="revision in recentRevisions" 
                  :key="revision.id"
                  class="revision-item"
                  @click="goToRevision(revision)"
                >
                  <div class="revision-title">{{ revision.title }}</div>
                  <div class="revision-meta">
                    {{ formatDate(revision.created_at) }} by {{ revision.created_by }}
                  </div>
                </div>
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
import { useI18n } from 'vue-i18n'
import { useWikiStore } from '@/stores/wikiStore'
import { useAuthStore } from '@/stores/authStore'
import { formatDate } from '@/utils/dateUtils'
import { marked } from 'marked'
import type { WikiPage, WikiRevision } from '@/types'
import WikiNavigation from '@/components/wiki/WikiNavigation.vue'
import WikiPageActions from '@/components/wiki/WikiPageActions.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const wikiStore = useWikiStore()
const authStore = useAuthStore()

// State
const page = ref<WikiPage | null>(null)
const relatedPages = ref<WikiPage[]>([])
const recentRevisions = ref<WikiRevision[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const isPageNotFound = ref(false)

// Computed
const isAuthenticated = computed(() => authStore.isAuthenticated)
const renderedContent = computed(() => {
  if (!page.value?.content) return ''
  
  // Configure marked options
  marked.setOptions({
    breaks: true, // Convert line breaks to <br>
    gfm: true, // GitHub Flavored Markdown
  })
  
  // Render markdown to HTML
  return marked(page.value.content)
})

// Methods
const loadPage = async () => {
  const slug = route.params.slug as string
  if (!slug) return

  loading.value = true
  error.value = null
  isPageNotFound.value = false

  try {
    const pageData = await wikiStore.fetchPageBySlug(slug)
    page.value = pageData
    
    // Load related data
    await Promise.all([
      loadRelatedPages(),
      loadRecentRevisions()
    ])
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'An error occurred'
    
    // Check if it's a 404 error (page not found)
    if (errorMessage.includes('404') || errorMessage.includes('not found') || errorMessage.includes('Page not found')) {
      isPageNotFound.value = true
      error.value = 'Page not found'
    } else {
      error.value = errorMessage
    }
  } finally {
    loading.value = false
  }
}

const loadRelatedPages = async () => {
  // Implement related pages logic
  relatedPages.value = []
}

const loadRecentRevisions = async () => {
  if (page.value) {
    const revisions = await wikiStore.fetchPageHistory(page.value.slug)
    recentRevisions.value = revisions.revisions.slice(0, 5)
  }
}

const editPage = () => {
  if (page.value) {
    router.push(`/wiki/edit/${page.value.slug}`)
  }
}

const deletePage = async () => {
  if (page.value && confirm(t('wiki.confirmDelete', { title: page.value.title }))) {
    await wikiStore.deletePage(page.value.id)
    router.push('/wiki')
  }
}

const viewHistory = () => {
  if (page.value) {
    router.push(`/wiki/history/${page.value.slug}`)
  }
}

const lockPage = async () => {
  // Implement lock page logic
}

const unlockPage = async () => {
  // Implement unlock page logic
}

const goToPage = (slug: string) => {
  router.push(`/wiki/page/${slug}`)
}

const goToRevision = (revision: WikiRevision) => {
  if (page.value) {
    router.push(`/wiki/page/${page.value.slug}?version=${revision.version}`)
  }
}

const goHome = () => {
  router.push('/wiki')
}

const createPage = () => {
  const slug = route.params.slug as string
  router.push(`/wiki/edit/${slug}`)
}

const goToSignIn = () => {
  router.push('/signin')
}

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-secondary',
    published: 'bg-success',
    archived: 'bg-warning',
    protected: 'bg-info'
  }
  return classes[status] || 'bg-secondary'
}

// Lifecycle
onMounted(() => {
  loadPage()
})
</script>

<style scoped>
.wiki-page-view {
  min-height: 100vh;
  background-color: #f8f9fa;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

.error-container {
  padding: 4rem 2rem;
  min-height: 50vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-container .alert {
  max-width: 600px;
  width: 100%;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.error-container .display-1 {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.page-header {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #333;
}

.page-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  color: #6c757d;
  font-size: 0.9rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.page-content {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 2rem;
}

.markdown-content {
  line-height: 1.6;
  color: #333;
}

.markdown-content h1,
.markdown-content h2,
.markdown-content h3,
.markdown-content h4,
.markdown-content h5,
.markdown-content h6 {
  margin-top: 2rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.markdown-content p {
  margin-bottom: 1rem;
}

.markdown-content ul,
.markdown-content ol {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.markdown-content blockquote {
  border-left: 4px solid #667eea;
  padding-left: 1rem;
  margin: 1rem 0;
  color: #6c757d;
}

.markdown-content code {
  background-color: #f8f9fa;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.markdown-content a {
  color: #667eea;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s ease;
}

.markdown-content a:hover {
  color: #5a67d8;
  border-bottom-color: #5a67d8;
}

.markdown-content a:visited {
  color: #805ad5;
}

.markdown-content table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.markdown-content th,
.markdown-content td {
  border: 1px solid #dee2e6;
  padding: 0.75rem;
  text-align: left;
}

.markdown-content th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.markdown-content pre {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-content img {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.page-footer {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid #dee2e6;
}

.categories-tags {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.categories,
.tags {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.page-sidebar {
  position: sticky;
  top: 2rem;
}

.sidebar-section {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 1.5rem;
}

.sidebar-section h5 {
  margin-bottom: 1rem;
  color: #333;
  font-weight: 600;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-weight: 500;
  color: #6c757d;
}

.related-pages {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.related-page-item {
  padding: 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.related-page-item:hover {
  border-color: #667eea;
  background-color: #f8f9fa;
}

.related-page-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 0.25rem;
}

.related-page-summary {
  font-size: 0.875rem;
  color: #6c757d;
  line-height: 1.4;
}

.recent-changes {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.revision-item {
  padding: 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.revision-item:hover {
  border-color: #667eea;
  background-color: #f8f9fa;
}

.revision-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 0.25rem;
}

.revision-meta {
  font-size: 0.875rem;
  color: #6c757d;
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .page-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .page-content,
  .page-sidebar .sidebar-section {
    padding: 1rem;
  }
}
</style>
