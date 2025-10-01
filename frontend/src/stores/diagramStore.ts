import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './authStore'
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Types
export interface Diagram {
  id: string
  realm_id: string
  name: string
  type: string
  script_type: string
  status: string
  description: string
  script: string
  theme: string
  tags: string[]
  version: number
  view_count: number
  edit_count: number
  public: boolean
  shared: boolean
  share_token?: string
  language: string
  created_by: string
  created_at: string
  updated_by: string
  updated_at: string
  images?: Image[]
  tag_objects?: DiagramTag[]
}

export interface Image {
  id: string
  realm_id: string
  user_id: string
  source: string
  type: string
  status: string
  original_name: string
  file_name: string
  file_path: string
  file_size: number
  mime_type: string
  extension: string
  width: number
  height: number
  format: string
  quality: number
  color_space: string
  has_alpha: boolean
  processed_at?: string
  error_msg: string
  is_public: boolean
  is_shared: boolean
  share_token?: string
  tags: string[]
  category: string
  description: string
  view_count: number
  download_count: number
  created_by: string
  created_at: string
  updated_by: string
  updated_at: string
  diagram_id?: string
  url: string
  thumbnail_url: string
  download_url: string
  aspect_ratio: number
  formatted_size: string
}

export interface DiagramTag {
  id: string
  name: string
  slug: string
  description: string
}

export interface CreateDiagramRequest {
  name: string
  type: string
  script_type: string
  script: string
  description?: string
  tags?: string[]
  public?: boolean
  shared?: boolean
}

export interface UpdateDiagramRequest {
  name?: string
  type?: string
  script_type?: string
  script?: string
  description?: string
  tags?: string[]
  public?: boolean
  shared?: boolean
  status?: string
  theme?: string
}

export interface DrawDiagramRequest {
  script: string
  script_type: string
  width?: number
  height?: number
  format?: string
}

export interface DrawDiagramResponse {
  image_data: string
  format: string
  width: number
  height: number
  size: number
  url: string
}

export interface DiagramListResponse {
  diagrams: Diagram[]
  total: number
  page: number
  limit: number
}

export interface DiagramSearchResponse {
  diagrams: Diagram[]
  total: number
  page: number
  limit: number
}

export const useDiagramStore = defineStore('diagram', () => {
  // State
  const diagrams = ref<Diagram[]>([])
  const currentDiagram = ref<Diagram | null>(null)
  const recentDiagrams = ref<Diagram[]>([])
  const searchResults = ref<Diagram[]>([])
  const tags = ref<DiagramTag[]>([])
  const generatedImage = ref<DrawDiagramResponse | null>(null)
  const stats = ref<any>(null)
  const loading = ref(false)
  const generating = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)
  const successMessage = ref<string | null>(null)
  
  // Pagination
  const currentPage = ref(1)
  const pageSize = ref(20)
  const totalPages = ref(0)
  const total = ref(0)
  
  // Filters
  const filters = ref({
    type: undefined as string | undefined,
    script_type: undefined as string | undefined,
    status: undefined as string | undefined,
    user_id: undefined as string | undefined
  })

  // Getters
  const authStore = useAuthStore()
  const hasDiagrams = computed(() => diagrams.value.length > 0)
  const hasSearchResults = computed(() => searchResults.value.length > 0)
  const hasRecentDiagrams = computed(() => recentDiagrams.value.length > 0)
  const totalPagesCount = computed(() => Math.ceil(total.value / pageSize.value))
  const isGenerating = computed(() => generating.value)
  const isSaving = computed(() => saving.value)

  // Helper function to get headers
  const getHeaders = () => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    if (authStore.token) {
      headers.Authorization = `Bearer ${authStore.token}`
    }
    
    const realmId = localStorage.getItem('realm_id') || 'default'
    headers['X-Realm-ID'] = realmId
    
    return headers
  }

  // Actions
  const clearError = () => {
    error.value = null
  }

  const clearSuccess = () => {
    successMessage.value = null
  }

  const clearMessages = () => {
    error.value = null
    successMessage.value = null
  }

  // Diagram CRUD Operations
  const createDiagram = async (request: CreateDiagramRequest): Promise<Diagram | null> => {
    saving.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams`, request, {
        headers: getHeaders()
      })
      const diagram = response.data
      diagrams.value.unshift(diagram)
      successMessage.value = 'Diagram created successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create diagram'
      return null
    } finally {
      saving.value = false
    }
  }

  const updateDiagram = async (id: string, request: UpdateDiagramRequest): Promise<Diagram | null> => {
    saving.value = true
    clearMessages()
    
    try {
      const response = await axios.put<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}`, request, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      // Update current diagram if it's the same
      if (currentDiagram.value?.id === id) {
        currentDiagram.value = diagram
      }
      
      successMessage.value = 'Diagram updated successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update diagram'
      return null
    } finally {
      saving.value = false
    }
  }

  const deleteDiagram = async (id: string): Promise<boolean> => {
    loading.value = true
    clearMessages()
    
    try {
      await axios.delete(`${API_BASE_URL}/api/v1/admin/diagrams/${id}`, {
        headers: getHeaders()
      })
      
      // Remove from diagrams array
      diagrams.value = diagrams.value.filter(d => d.id !== id)
      
      // Clear current diagram if it's the same
      if (currentDiagram.value?.id === id) {
        currentDiagram.value = null
      }
      
      successMessage.value = 'Diagram deleted successfully'
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete diagram'
      return false
    } finally {
      loading.value = false
    }
  }

  const getDiagram = async (id: string, includeImages = false): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<Diagram>(`${API_BASE_URL}/api/v1/diagrams/${id}`, {
        headers: getHeaders(),
        params: { include_images: includeImages }
      })
      const diagram = response.data
      currentDiagram.value = diagram
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch diagram'
      return null
    } finally {
      loading.value = false
    }
  }

  const listDiagrams = async (params: {
    page?: number
    limit?: number
    type?: string
    script_type?: string
    status?: string
    user_id?: string
  } = {}): Promise<DiagramListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<DiagramListResponse>(`${API_BASE_URL}/api/v1/diagrams`, {
        headers: getHeaders(),
        params: {
          page: currentPage.value,
          limit: pageSize.value,
          ...params
        }
      })
      
      if (response.data) {
        diagrams.value = response.data.diagrams
        total.value = response.data.total
        totalPages.value = Math.ceil(response.data.total / pageSize.value)
      }
      
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch diagrams'
      return null
    } finally {
      loading.value = false
    }
  }

  const searchDiagrams = async (query: string, page = 1, limit = 20): Promise<DiagramListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<DiagramSearchResponse>(`${API_BASE_URL}/api/v1/diagrams/search`, {
        headers: getHeaders(),
        params: { query, page, limit }
      })
      if (response.data) {
        searchResults.value = response.data.diagrams
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to search diagrams'
      return null
    } finally {
      loading.value = false
    }
  }

  const getPublicDiagrams = async (page = 1, limit = 20): Promise<DiagramListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<DiagramListResponse>(`${API_BASE_URL}/api/v1/diagrams/public`, {
        headers: getHeaders(),
        params: { page, limit }
      })
      if (response.data) {
        diagrams.value = response.data.diagrams
        total.value = response.data.total
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch public diagrams'
      return null
    } finally {
      loading.value = false
    }
  }

  const getSharedDiagrams = async (page = 1, limit = 20): Promise<DiagramListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<DiagramListResponse>(`${API_BASE_URL}/api/v1/diagrams/shared`, {
        headers: getHeaders(),
        params: { page, limit }
      })
      if (response.data) {
        diagrams.value = response.data.diagrams
        total.value = response.data.total
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch shared diagrams'
      return null
    } finally {
      loading.value = false
    }
  }

  // Diagram Generation
  const drawDiagram = async (request: DrawDiagramRequest): Promise<DrawDiagramResponse | null> => {
    generating.value = true
    clearMessages()
    
    try {
      const response = await axios.post<DrawDiagramResponse>(`${API_BASE_URL}/api/v1/diagrams/draw`, request, {
        headers: getHeaders()
      })
      const result = response.data
      generatedImage.value = result
      successMessage.value = 'Diagram generated successfully'
      return result
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to generate diagram'
      return null
    } finally {
      generating.value = false
    }
  }

  const clearGeneratedImage = () => {
    generatedImage.value = null
  }

  // Diagram Management
  const publishDiagram = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/publish`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram published successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to publish diagram'
      return null
    } finally {
      loading.value = false
    }
  }

  const archiveDiagram = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/archive`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram archived successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to archive diagram'
      return null
    } finally {
      loading.value = false
    }
  }

  const makeDiagramPublic = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/public`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram made public successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to make diagram public'
      return null
    } finally {
      loading.value = false
    }
  }

  const makeDiagramPrivate = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/private`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram made private successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to make diagram private'
      return null
    } finally {
      loading.value = false
    }
  }

  const shareDiagram = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/share`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram shared successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to share diagram'
      return null
    } finally {
      loading.value = false
    }
  }

  const unshareDiagram = async (id: string): Promise<Diagram | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Diagram>(`${API_BASE_URL}/api/v1/admin/diagrams/${id}/unshare`, {}, {
        headers: getHeaders()
      })
      const diagram = response.data
      
      // Update in diagrams array
      const index = diagrams.value.findIndex(d => d.id === id)
      if (index !== -1) {
        diagrams.value[index] = diagram
      }
      
      successMessage.value = 'Diagram unshared successfully'
      return diagram
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to unshare diagram'
      return null
    } finally {
      loading.value = false
    }
  }

  // Tags
  const getAllTags = async (): Promise<DiagramTag[]> => {
    try {
      const response = await axios.get<DiagramTag[]>(`${API_BASE_URL}/api/v1/diagrams/tags`, {
        headers: getHeaders()
      })
      const tagsList = response.data
      tags.value = tagsList
      return tagsList
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch tags'
      return []
    }
  }

  const getTag = async (id: string): Promise<DiagramTag | null> => {
    try {
      const response = await axios.get<DiagramTag>(`${API_BASE_URL}/api/v1/diagrams/tags/${id}`, {
        headers: getHeaders()
      })
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch tag'
      return null
    }
  }

  // Recent Diagrams
  const loadRecentDiagrams = async (limit = 5): Promise<void> => {
    try {
      const response = await axios.get<DiagramListResponse>(`${API_BASE_URL}/api/v1/diagrams`, {
        headers: getHeaders(),
        params: { limit }
      })
      if (response.data) {
        recentDiagrams.value = response.data.diagrams
      }
    } catch (err: any) {
      console.error('Failed to load recent diagrams:', err)
    }
  }

  // Statistics
  const getDiagramStats = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/v1/diagrams/stats`, {
        headers: getHeaders()
      })
      stats.value = response.data
    } catch (err: any) {
      console.error('Failed to load diagram stats:', err)
    }
  }

  // Pagination
  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPagesCount.value) {
      currentPage.value = page
    }
  }

  const nextPage = () => {
    if (currentPage.value < totalPagesCount.value) {
      currentPage.value++
    }
  }

  const previousPage = () => {
    if (currentPage.value > 1) {
      currentPage.value--
    }
  }

  // Filters
  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    filters.value = { ...filters.value, ...newFilters }
    currentPage.value = 1 // Reset to first page when filters change
  }

  const clearFilters = () => {
    filters.value = {
      type: undefined,
      script_type: undefined,
      status: undefined,
      user_id: undefined
    }
    currentPage.value = 1
  }

  // Reset store
  const reset = () => {
    diagrams.value = []
    currentDiagram.value = null
    recentDiagrams.value = []
    searchResults.value = []
    tags.value = []
    generatedImage.value = null
    loading.value = false
    generating.value = false
    saving.value = false
    error.value = null
    successMessage.value = null
    currentPage.value = 1
    totalPages.value = 0
    total.value = 0
    clearFilters()
  }

  return {
    // State
    diagrams,
    currentDiagram,
    recentDiagrams,
    searchResults,
    tags,
    generatedImage,
    stats,
    loading,
    generating,
    saving,
    error,
    successMessage,
    currentPage,
    pageSize,
    totalPages,
    total,
    filters,
    
    // Getters
    hasDiagrams,
    hasSearchResults,
    hasRecentDiagrams,
    totalPagesCount,
    isGenerating,
    isSaving,
    
    // Actions
    clearError,
    clearSuccess,
    clearMessages,
    createDiagram,
    updateDiagram,
    deleteDiagram,
    getDiagram,
    listDiagrams,
    searchDiagrams,
    getPublicDiagrams,
    getSharedDiagrams,
    drawDiagram,
    clearGeneratedImage,
    publishDiagram,
    archiveDiagram,
    makeDiagramPublic,
    makeDiagramPrivate,
    shareDiagram,
    unshareDiagram,
    getAllTags,
    getTag,
    loadRecentDiagrams,
    getDiagramStats,
    goToPage,
    nextPage,
    previousPage,
    setFilters,
    clearFilters,
    reset
  }
})
