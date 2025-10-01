<template>
  <div class="wiki-view">
    <!-- Header Section -->
    <div class="wiki-header">
      <div class="container">
        <div class="row align-items-center">
          <div class="col-md-8">
            <h1 class="wiki-title">
              <i class="bi bi-book"></i>
              {{ $t('wiki.title') }}
            </h1>
            <p class="wiki-subtitle">{{ $t('wiki.subtitle') }}</p>
          </div>
          <div class="col-md-4 text-end">
            <WikiSearchBar 
              @search="handleSearch"
              :loading="loading"
              class="wiki-search-bar"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="container">
      <div class="row">
        <!-- Sidebar -->
        <div class="col-lg-2">
          <WikiSidebar
            :categories="categories"
            :tags="tags"
            :recentPages="recentPages"
            :specialPages="specialPages"
            @category-click="handleCategoryClick"
            @tag-click="handleTagClick"
            @special-page-click="handleSpecialPageClick"
          />
        </div>

        <!-- Main Content Area -->
        <div class="col-lg-10">
          <!-- Quick Actions -->
          <div class="wiki-actions mb-4" v-if="isAuthenticated">
            <button 
              class="btn btn-primary"
              @click="createNewPage"
            >
              <i class="bi bi-plus-lg"></i>
              {{ $t('wiki.createPage') }}
            </button>
            <button 
              class="btn btn-outline-secondary"
              @click="goToRandomPage"
              :disabled="loading"
            >
              <i class="bi bi-shuffle"></i>
              {{ $t('wiki.randomPage') }}
            </button>
          </div>

          <!-- Content Tabs -->
          <div class="wiki-content">
            <ul class="nav nav-tabs" role="tablist">
              <li class="nav-item" role="presentation">
                <button 
                  class="nav-link"
                  :class="{ active: activeTab === 'pages' }"
                  @click="activeTab = 'pages'"
                >
                  <i class="bi bi-file-text"></i>
                  {{ $t('wiki.allPages') }}
                </button>
              </li>
              <li class="nav-item" role="presentation">
                <button 
                  class="nav-link"
                  :class="{ active: activeTab === 'recent' }"
                  @click="activeTab = 'recent'"
                >
                  <i class="bi bi-clock-history"></i>
                  {{ $t('wiki.recentChanges') }}
                </button>
              </li>
              <li class="nav-item" role="presentation">
                <button 
                  class="nav-link"
                  :class="{ active: activeTab === 'search' }"
                  @click="activeTab = 'search'"
                  v-if="hasSearchResults"
                >
                  <i class="bi bi-search"></i>
                  {{ $t('wiki.searchResults') }}
                  <span class="badge bg-primary ms-1">{{ searchResults.length }}</span>
                </button>
              </li>
            </ul>

            <div class="tab-content">
              <!-- All Pages Tab -->
              <div v-if="activeTab === 'pages'" class="tab-pane active">
                <WikiPageList 
                  :pages="pages"
                  :loading="loading"
                  :showActions="isAuthenticated"
                  @page-click="goToPage"
                  @edit-page="editPage"
                  @delete-page="deletePage"
                  @create-page="createNewPage"
                />
                
                <!-- Pagination -->
                <nav v-if="totalPagesCount > 1" class="mt-4">
                  <ul class="pagination justify-content-center">
                    <li class="page-item" :class="{ disabled: currentPageNum === 1 }">
                      <button class="page-link" @click="changePage(currentPageNum - 1)">
                        {{ $t('common.previous') }}
                      </button>
                    </li>
                    <li 
                      v-for="page in visiblePages" 
                      :key="page"
                      class="page-item"
                      :class="{ active: page === currentPageNum }"
                    >
                      <button class="page-link" @click="changePage(page)">
                        {{ page }}
                      </button>
                    </li>
                    <li class="page-item" :class="{ disabled: currentPageNum === totalPagesCount }">
                      <button class="page-link" @click="changePage(currentPageNum + 1)">
                        {{ $t('common.next') }}
                      </button>
                    </li>
                  </ul>
                </nav>
              </div>

              <!-- Recent Changes Tab -->
              <div v-if="activeTab === 'recent'" class="tab-pane">
                <WikiRevisionList 
                  :revisions="revisions"
                  :loading="loading"
                  :showCompare="true"
                  @revision-click="goToRevision"
                  @compare-revisions="compareRevisions"
                />
              </div>

              <!-- Search Results Tab -->
              <div v-if="activeTab === 'search'" class="tab-pane">
                <div class="search-results-header">
                  <h5>{{ $t('wiki.searchResultsFor', { query: searchQuery }) }}</h5>
                  <p class="text-muted">{{ $t('wiki.foundResults', { count: searchResults.length }) }}</p>
                </div>
                <WikiPageList 
                  :pages="searchResults"
                  :loading="loading"
                  :showActions="isAuthenticated"
                  @page-click="goToPage"
                  @edit-page="editPage"
                  @delete-page="deletePage"
                  @create-page="createNewPage"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">{{ $t('common.loading') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useWikiStore } from '@/stores/wikiStore'
import { useAuthStore } from '@/stores/authStore'
import WikiSearchBar from '@/components/wiki/WikiSearchBar.vue'
import WikiSidebar from '@/components/wiki/WikiSidebar.vue'
import WikiPageList from '@/components/wiki/WikiPageList.vue'
import WikiRevisionList from '@/components/wiki/WikiRevisionList.vue'

const { t } = useI18n()
const router = useRouter()
const wikiStore = useWikiStore()
const authStore = useAuthStore()

// State
const activeTab = ref('pages')
const categories = ref<Array<{ name: string; count: number }>>([])
const tags = ref<Array<{ name: string; count: number }>>([])
const recentPages = ref([])
const specialPages = ref<Array<{ type: string; title: string }>>([])

// Computed
const isAuthenticated = computed(() => authStore.isAuthenticated)
const pages = computed(() => wikiStore.pages)
const searchResults = computed(() => wikiStore.searchResults)
const revisions = computed(() => wikiStore.revisions)
const loading = computed(() => wikiStore.loading)
const searchQuery = computed(() => wikiStore.searchQuery)
const currentPageNum = computed(() => wikiStore.currentPageNum)
const totalPagesCount = computed(() => wikiStore.totalPagesCount)
const hasSearchResults = computed(() => wikiStore.hasSearchResults)

// Pagination
const visiblePages = computed(() => {
  const total = totalPagesCount.value
  const current = currentPageNum.value
  const pages = []
  
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// Methods
const handleSearch = async (query: string) => {
  if (query.trim()) {
    await wikiStore.searchPages(query)
    activeTab.value = 'search'
  }
}

const handleCategoryClick = async (category: string) => {
  await wikiStore.fetchPagesByCategory(category)
  activeTab.value = 'pages'
}

const handleTagClick = async (tag: string) => {
  // Implement tag filtering
  activeTab.value = 'pages'
}

const handleSpecialPageClick = (type: string) => {
  router.push(`/wiki/special/${type}`)
}

const goToPage = (page: any) => {
  router.push(`/wiki/page/${page.slug}`)
}

const editPage = (page: any) => {
  router.push(`/wiki/edit/${page.slug}`)
}

const deletePage = async (page: any) => {
  if (confirm(t('wiki.confirmDelete', { title: page.title }))) {
    await wikiStore.deletePage(page.id)
    await loadPages()
  }
}

const createNewPage = () => {
  router.push('/wiki/edit')
}

const goToRandomPage = async () => {
  try {
    const randomPage = await wikiStore.fetchRandomPage()
    if (randomPage) {
      router.push(`/wiki/page/${randomPage.slug}`)
    }
  } catch (error) {
    console.error('Failed to get random page:', error)
  }
}

const changePage = async (page: number) => {
  await wikiStore.fetchPages({ page })
}

const goToRevision = (revision: any) => {
  router.push(`/wiki/page/${revision.page_id}?version=${revision.version}`)
}

const compareRevisions = (from: any, to: any) => {
  router.push(`/wiki/compare/${from.page_id}/${from.version}/${to.version}`)
}

const loadPages = async () => {
  await wikiStore.fetchPages()
}

const loadRecentChanges = async () => {
  await wikiStore.fetchRecentChanges()
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    loadPages(),
    loadRecentChanges(),
    loadSidebarData()
  ])
})

const loadSidebarData = async () => {
  try {
    const sidebarData = await wikiStore.fetchSidebarData()
    recentPages.value = sidebarData.recentPages
    categories.value = sidebarData.categories
    tags.value = sidebarData.tags
    specialPages.value = sidebarData.specialPages
  } catch (error) {
    console.error('Failed to load sidebar data:', error)
  }
}

// Watch for tab changes
watch(activeTab, (newTab) => {
  if (newTab === 'recent' && revisions.value.length === 0) {
    loadRecentChanges()
  }
})
</script>

<style scoped>
.wiki-view {
  min-height: 100vh;
  background-color: #f8f9fa;
}

.wiki-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 3rem 0;
  margin-bottom: 2rem;
}

.wiki-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.wiki-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  margin-bottom: 0;
}

.wiki-search-bar {
  max-width: 300px;
}

.wiki-actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.wiki-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.nav-tabs {
  border-bottom: 1px solid #dee2e6;
  padding: 0 1.5rem;
  margin-bottom: 0;
}

.nav-tabs .nav-link {
  border: none;
  border-bottom: 3px solid transparent;
  color: #6c757d;
  font-weight: 500;
  padding: 1rem 1.5rem;
  transition: all 0.3s ease;
}

.nav-tabs .nav-link:hover {
  border-bottom-color: #dee2e6;
  color: #495057;
}

.nav-tabs .nav-link.active {
  border-bottom-color: #667eea;
  color: #667eea;
  background: none;
}

.tab-content {
  padding: 2rem;
}

.search-results-header {
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #dee2e6;
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

@media (max-width: 768px) {
  .wiki-header {
    padding: 2rem 0;
  }
  
  .wiki-title {
    font-size: 2rem;
  }
  
  .wiki-actions {
    flex-direction: column;
  }
  
  .tab-content {
    padding: 1rem;
  }
}
</style>
