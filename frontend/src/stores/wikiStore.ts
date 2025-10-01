import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './authStore'
import { getApiUrl } from '@/utils/apiConfig'
import type { 
  WikiPage, 
  WikiRevision, 
  CreateWikiPageRequest, 
  UpdateWikiPageRequest,
  WikiPageListResponse,
  WikiRevisionListResponse 
} from '@/types'

export const useWikiStore = defineStore('wiki', () => {
  // State
  const pages = ref<WikiPage[]>([])
  const currentPage = ref<WikiPage | null>(null)
  const revisions = ref<WikiRevision[]>([])
  const searchResults = ref<WikiPage[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')
  const currentPageNum = ref(1)
  const pageSize = ref(20)
  const totalPages = ref(0)
  
  // Filters
  const filters = ref({
    status: undefined as string | undefined,
    type: undefined as string | undefined,
    category: undefined as string | undefined,
    tag: undefined as string | undefined
  })

  // Getters
  const authStore = useAuthStore()
  const hasPages = computed(() => pages.value.length > 0)
  const hasSearchResults = computed(() => searchResults.value.length > 0)
  const totalPagesCount = computed(() => Math.ceil(totalPages.value / pageSize.value))

  // Helper function to get headers
  const getHeaders = () => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    if (authStore.token) {
      headers.Authorization = `Bearer ${authStore.token}`
    }
    
    return headers
  }

  // Public API calls (with optional authentication)
  const fetchPages = async (params: {
    page?: number
    limit?: number
    status?: string
    type?: string
  } = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.append('page', (params.page || currentPageNum.value).toString())
      queryParams.append('limit', (params.limit || pageSize.value).toString())
      
      if (params.status) queryParams.append('status', params.status)
      if (params.type) queryParams.append('type', params.type)
      
      const response = await fetch(getApiUrl(`/api/v1/wiki/pages?${queryParams}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch pages')
      }
      
      const data: WikiPageListResponse = await response.json()
      pages.value = data.pages
      totalPages.value = data.total
      currentPageNum.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchPageBySlug = async (slug: string) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl(`/api/v1/wiki/page/${slug}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Page not found')
        }
        throw new Error('Failed to fetch page')
      }
      
      const data = await response.json()
      currentPage.value = data.page
      return data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const searchPages = async (query: string, params: {
    page?: number
    limit?: number
  } = {}) => {
    loading.value = true
    error.value = null
    searchQuery.value = query
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.append('q', query)
      queryParams.append('page', (params.page || 1).toString())
      queryParams.append('limit', (params.limit || pageSize.value).toString())
      
      const response = await fetch(getApiUrl(`/api/v1/wiki/search?${queryParams}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to search pages')
      }
      
      const data: WikiPageListResponse = await response.json()
      searchResults.value = data.pages
      totalPages.value = data.total
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchPagesByCategory = async (category: string, params: {
    page?: number
    limit?: number
  } = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.append('page', (params.page || 1).toString())
      queryParams.append('limit', (params.limit || pageSize.value).toString())
      
      const response = await fetch(getApiUrl(`/api/v1/wiki/category/${category}?${queryParams}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch pages by category')
      }
      
      const data: WikiPageListResponse = await response.json()
      pages.value = data.pages
      totalPages.value = data.total
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchRandomPage = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl('/api/v1/wiki/random'), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch random page')
      }
      
      const data = await response.json()
      return data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchRecentChanges = async (params: {
    page?: number
    limit?: number
  } = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.append('page', (params.page || 1).toString())
      queryParams.append('limit', (params.limit || 20).toString())
      
      const response = await fetch(getApiUrl(`/api/v1/wiki/recent-changes?${queryParams}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch recent changes')
      }
      
      const data = await response.json()
      revisions.value = data.revisions
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Sidebar data methods
  const fetchSidebarData = async () => {
    try {
      // Fetch recent pages (limit to 5 for sidebar)
      const recentResponse = await fetch(getApiUrl('/api/v1/wiki/pages?limit=5'), {
        headers: getHeaders()
      })
      
      if (recentResponse.ok) {
        const recentData = await recentResponse.json()
        return {
          recentPages: recentData.pages || [],
          categories: [], // TODO: Implement categories endpoint
          tags: [], // TODO: Implement tags endpoint
          specialPages: [
            { type: 'recent-changes', title: 'Recent Changes' },
            { type: 'random', title: 'Random Page' },
            { type: 'orphaned', title: 'Orphaned Pages' },
            { type: 'wanted', title: 'Wanted Pages' }
          ]
        }
      }
      
      return {
        recentPages: [],
        categories: [],
        tags: [],
        specialPages: [
          { type: 'recent-changes', title: 'Recent Changes' },
          { type: 'random', title: 'Random Page' },
          { type: 'orphaned', title: 'Orphaned Pages' },
          { type: 'wanted', title: 'Wanted Pages' }
        ]
      }
    } catch (err) {
      console.error('Failed to fetch sidebar data:', err)
      return {
        recentPages: [],
        categories: [],
        tags: [],
        specialPages: [
          { type: 'recent-changes', title: 'Recent Changes' },
          { type: 'random', title: 'Random Page' },
          { type: 'orphaned', title: 'Orphaned Pages' },
          { type: 'wanted', title: 'Wanted Pages' }
        ]
      }
    }
  }

  const fetchPageHistory = async (slug: string, params: {
    page?: number
    limit?: number
  } = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.append('page', (params.page || 1).toString())
      queryParams.append('limit', (params.limit || 20).toString())
      
      const response = await fetch(getApiUrl(`/api/v1/wiki/page/${slug}/history?${queryParams}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch page history')
      }
      
      const data = await response.json()
      revisions.value = data.revisions
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Admin API calls (require authentication)
  const createPage = async (pageData: CreateWikiPageRequest) => {
    loading.value = true
    error.value = null
    
    try {
      // Add realm_id from auth store if not provided
      const requestData = {
        ...pageData,
        realm_id: pageData.realm_id || authStore.currentUser?.realm_id
      }
      
      const response = await fetch(getApiUrl('/api/v1/admin/wiki/pages'), {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(requestData)
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create page')
      }
      
      const data = await response.json()
      pages.value.unshift(data)
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updatePage = async (id: string, pageData: Partial<CreateWikiPageRequest>) => {
    loading.value = true
    error.value = null
    
    try {
      // Add realm_id from auth store if not provided
      const requestData = {
        ...pageData,
        realm_id: pageData.realm_id || authStore.currentUser?.realm_id
      }
      
      const response = await fetch(getApiUrl(`/api/v1/admin/wiki/pages/${id}`), {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(requestData)
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update page')
      }
      
      const data = await response.json()
      const index = pages.value.findIndex(p => p.id === id)
      if (index !== -1) {
        pages.value[index] = data
      }
      if (currentPage.value?.id === id) {
        currentPage.value = data
      }
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deletePage = async (id: string) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl(`/api/v1/admin/wiki/pages/${id}`), {
        method: 'DELETE',
        headers: getHeaders()
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to delete page')
      }
      
      pages.value = pages.value.filter(p => p.id !== id)
      if (currentPage.value?.id === id) {
        currentPage.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createRevision = async (pageId: string, revisionData: {
    title: string
    content: string
    summary?: string
    change_note?: string
  }) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl(`/api/v1/admin/wiki/pages/${pageId}/revisions`), {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(revisionData)
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create revision')
      }
      
      const data = await response.json()
      revisions.value.unshift(data)
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Utility functions
  const clearError = () => {
    error.value = null
  }

  const clearSearchResults = () => {
    searchResults.value = []
    searchQuery.value = ''
  }

  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  const clearFilters = () => {
    filters.value = {
      status: undefined,
      type: undefined,
      category: undefined,
      tag: undefined
    }
  }

  return {
    // State
    pages,
    currentPage,
    revisions,
    searchResults,
    loading,
    error,
    searchQuery,
    currentPageNum,
    pageSize,
    totalPages,
    filters,
    
    // Getters
    hasPages,
    hasSearchResults,
    totalPagesCount,
    
    // Actions
    fetchPages,
    fetchPageBySlug,
    searchPages,
    fetchPagesByCategory,
    fetchRandomPage,
    fetchRecentChanges,
    fetchPageHistory,
    fetchSidebarData,
    createPage,
    updatePage,
    deletePage,
    createRevision,
    clearError,
    clearSearchResults,
    setFilters,
    clearFilters
  }
})
