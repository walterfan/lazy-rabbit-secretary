import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './authStore'
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Types
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
  diagram_id?: string // Optional for uploaded images
  url: string
  thumbnail_url: string
  download_url: string
  aspect_ratio: number
  formatted_size: string
}

export interface ImageStatsResponse {
  total_count: number
  total_size: number
  average_size: number
  type_counts: { [key: string]: number }
  status_counts: { [key: string]: number }
}

export interface ImageListResponse {
  images: Image[]
  total: number
  page: number
  limit: number
}

export interface UploadImageRequest {
  type: string
  category?: string
  description?: string
  tags?: string[]
  is_public?: boolean
  is_shared?: boolean
}

export interface UpdateImageRequest {
  type?: string
  category?: string
  description?: string
  tags?: string[]
  is_public?: boolean
  is_shared?: boolean
}

export const useImageStore = defineStore('image', () => {
  // State
  const images = ref<Image[]>([])
  const currentImage = ref<Image | null>(null)
  const searchResults = ref<Image[]>([])
  const stats = ref<ImageStatsResponse | null>(null)
  const categories = ref<string[]>([])
  const extensions = ref<string[]>([])
  const mimeTypes = ref<string[]>([])
  const loading = ref(false)
  const uploading = ref(false)
  const updating = ref(false)
  const error = ref<string | null>(null)
  const successMessage = ref<string | null>(null)
  
  // Pagination
  const currentPage = ref(1)
  const pageSize = ref(24)
  const totalPages = ref(0)
  const total = ref(0)
  
  // Filters
  const filters = ref({
    type: undefined as string | undefined,
    status: undefined as string | undefined,
    category: undefined as string | undefined,
    user_id: undefined as string | undefined,
    source: undefined as string | undefined
  })

  // Search
  const searchQuery = ref('')
  const sortBy = ref('created_at')

  // Getters
  const authStore = useAuthStore()
  const hasImages = computed(() => images.value.length > 0)
  const hasSearchResults = computed(() => searchResults.value.length > 0)
  const totalPagesCount = computed(() => Math.ceil(total.value / pageSize.value))
  const isUploading = computed(() => uploading.value)
  const isUpdating = computed(() => updating.value)
  const hasStats = computed(() => stats.value !== null)

  // Helper function to get headers
  const getHeaders = () => {
    const headers: Record<string, string> = {}
    
    if (authStore.token) {
      headers.Authorization = `Bearer ${authStore.token}`
    }
    
    const realmId = localStorage.getItem('realm_id') || 'default'
    headers['X-Realm-ID'] = realmId
    
    return headers
  }

  // Helper function to get upload headers (without Content-Type for FormData)
  const getUploadHeaders = () => {
    const headers: Record<string, string> = {}
    
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

  // Image CRUD Operations
  const uploadImage = async (file: File, request: UploadImageRequest): Promise<Image | null> => {
    uploading.value = true
    clearMessages()
    
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('type', request.type)
      if (request.category) formData.append('category', request.category)
      if (request.description) formData.append('description', request.description)
      if (request.tags) formData.append('tags', request.tags.join(','))
      if (request.is_public !== undefined) formData.append('is_public', request.is_public.toString())
      if (request.is_shared !== undefined) formData.append('is_shared', request.is_shared.toString())

      const response = await axios.post<Image>(`${API_BASE_URL}/api/v1/admin/images/upload`, formData, {
        headers: getUploadHeaders()
      })
      
      const image = response.data
      images.value.unshift(image)
      successMessage.value = 'Image uploaded successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to upload image'
      return null
    } finally {
      uploading.value = false
    }
  }

  const updateImage = async (id: string, request: UpdateImageRequest): Promise<Image | null> => {
    updating.value = true
    clearMessages()
    
    try {
      const response = await axios.put<Image>(`${API_BASE_URL}/api/v1/admin/images/${id}`, request, {
        headers: getHeaders()
      })
      const image = response.data
      
      // Update in images array
      const index = images.value.findIndex(i => i.id === id)
      if (index !== -1) {
        images.value[index] = image
      }
      
      // Update current image if it's the same
      if (currentImage.value?.id === id) {
        currentImage.value = image
      }
      
      successMessage.value = 'Image updated successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update image'
      return null
    } finally {
      updating.value = false
    }
  }

  const deleteImage = async (id: string): Promise<boolean> => {
    loading.value = true
    clearMessages()
    
    try {
      await axios.delete(`${API_BASE_URL}/api/v1/admin/images/${id}`, {
        headers: getHeaders()
      })
      
      // Remove from images array
      images.value = images.value.filter(i => i.id !== id)
      
      // Clear current image if it's the same
      if (currentImage.value?.id === id) {
        currentImage.value = null
      }
      
      successMessage.value = 'Image deleted successfully'
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete image'
      return false
    } finally {
      loading.value = false
    }
  }

  const getImage = async (id: string): Promise<Image | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<Image>(`${API_BASE_URL}/api/v1/images/${id}`, {
        headers: getHeaders()
      })
      const image = response.data
      currentImage.value = image
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch image'
      return null
    } finally {
      loading.value = false
    }
  }

  const listImages = async (params: {
    page?: number
    limit?: number
    type?: string
    status?: string
    category?: string
    user_id?: string
  } = {}): Promise<ImageListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<ImageListResponse>(`${API_BASE_URL}/api/v1/images`, {
        headers: getHeaders(),
        params: {
          page: currentPage.value,
          limit: pageSize.value,
          ...params
        }
      })
      
      if (response.data) {
        // Handle nested data structure from backend
        const data = (response.data as any).data || response.data
        images.value = data.images || []
        total.value = data.total || 0
        totalPages.value = Math.ceil((data.total || 0) / pageSize.value)
      }
      
      return (response.data as any).data || response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch images'
      return null
    } finally {
      loading.value = false
    }
  }

  const searchImages = async (query: string, page = 1, limit = 20): Promise<ImageListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<ImageListResponse>(`${API_BASE_URL}/api/v1/images/search`, {
        headers: getHeaders(),
        params: { query, page, limit }
      })
      if (response.data) {
        searchResults.value = response.data.images
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to search images'
      return null
    } finally {
      loading.value = false
    }
  }

  const getPublicImages = async (page = 1, limit = 20): Promise<ImageListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<ImageListResponse>(`${API_BASE_URL}/api/v1/images/public`, {
        headers: getHeaders(),
        params: { page, limit }
      })
      if (response.data) {
        images.value = response.data.images || []
        total.value = response.data.total
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch public images'
      return null
    } finally {
      loading.value = false
    }
  }

  const getSharedImages = async (page = 1, limit = 20): Promise<ImageListResponse | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.get<ImageListResponse>(`${API_BASE_URL}/api/v1/images/shared`, {
        headers: getHeaders(),
        params: { page, limit }
      })
      if (response.data) {
        images.value = response.data.images || []
        total.value = response.data.total
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch shared images'
      return null
    } finally {
      loading.value = false
    }
  }

  // Image Download
  const downloadImage = async (id: string): Promise<boolean> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/v1/images/${id}/download`, {
        headers: getHeaders(),
        responseType: 'blob'
      })
      
      const blob = response.data as Blob
      const image = images.value.find(i => i.id === id)
      const fileName = image?.original_name || `image-${id}`
      
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = fileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to download image'
      return false
    }
  }

  const getImageThumbnail = async (id: string): Promise<string | null> => {
    try {
      // For now, just return the thumbnail URL from the image's metadata
      const image = images.value.find(i => i.id === id)
      return image?.thumbnail_url || ''
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch thumbnail'
      return null
    }
  }

  // Image Management
  const makeImagePublic = async (id: string): Promise<Image | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Image>(`${API_BASE_URL}/api/v1/admin/images/${id}/public`, {}, {
        headers: getHeaders()
      })
      const image = response.data
      
      // Update in images array
      const index = images.value.findIndex(i => i.id === id)
      if (index !== -1) {
        images.value[index] = image
      }
      
      successMessage.value = 'Image made public successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to make image public'
      return null
    } finally {
      loading.value = false
    }
  }

  const makeImagePrivate = async (id: string): Promise<Image | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Image>(`${API_BASE_URL}/api/v1/admin/images/${id}/private`, {}, {
        headers: getHeaders()
      })
      const image = response.data
      
      // Update in images array
      const index = images.value.findIndex(i => i.id === id)
      if (index !== -1) {
        images.value[index] = image
      }
      
      successMessage.value = 'Image made private successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to make image private'
      return null
    } finally {
      loading.value = false
    }
  }

  const shareImage = async (id: string): Promise<Image | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Image>(`${API_BASE_URL}/api/v1/admin/images/${id}/share`, {}, {
        headers: getHeaders()
      })
      const image = response.data
      
      // Update in images array
      const index = images.value.findIndex(i => i.id === id)
      if (index !== -1) {
        images.value[index] = image
      }
      
      successMessage.value = 'Image shared successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to share image'
      return null
    } finally {
      loading.value = false
    }
  }

  const unshareImage = async (id: string): Promise<Image | null> => {
    loading.value = true
    clearMessages()
    
    try {
      const response = await axios.post<Image>(`${API_BASE_URL}/api/v1/admin/images/${id}/unshare`, {}, {
        headers: getHeaders()
      })
      const image = response.data
      
      // Update in images array
      const index = images.value.findIndex(i => i.id === id)
      if (index !== -1) {
        images.value[index] = image
      }
      
      successMessage.value = 'Image unshared successfully'
      return image
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to unshare image'
      return null
    } finally {
      loading.value = false
    }
  }

  // Statistics and Metadata
  const getImageStats = async (): Promise<ImageStatsResponse | null> => {
    try {
      const response = await axios.get<ImageStatsResponse>(`${API_BASE_URL}/api/v1/images/stats`, {
        headers: getHeaders()
      })
      // Handle nested data structure from backend
      const data = (response.data as any).data || response.data
      stats.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch image statistics'
      return null
    }
  }

  const getCategories = async (): Promise<string[]> => {
    try {
      const response = await axios.get<string[]>(`${API_BASE_URL}/api/v1/images/categories`, {
        headers: getHeaders()
      })
      categories.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch categories'
      return []
    }
  }

  const getExtensions = async (): Promise<string[]> => {
    try {
      const response = await axios.get<string[]>(`${API_BASE_URL}/api/v1/images/extensions`, {
        headers: getHeaders()
      })
      extensions.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch extensions'
      return []
    }
  }

  const getMimeTypes = async (): Promise<string[]> => {
    try {
      const response = await axios.get<string[]>(`${API_BASE_URL}/api/v1/images/mime-types`, {
        headers: getHeaders()
      })
      mimeTypes.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch MIME types'
      return []
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

  // Filters and Search
  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    filters.value = { ...filters.value, ...newFilters }
    currentPage.value = 1 // Reset to first page when filters change
  }

  const clearFilters = () => {
    filters.value = {
      type: undefined,
      status: undefined,
      category: undefined,
      user_id: undefined,
      source: undefined
    }
    currentPage.value = 1
  }

  const setSearchQuery = (query: string) => {
    searchQuery.value = query
    currentPage.value = 1
  }

  const setSortBy = (sort: string) => {
    sortBy.value = sort
    currentPage.value = 1
  }

  // Utility functions
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const getSourceBadgeClass = (source: string) => {
    return source === 'uploaded' ? 'bg-primary' : 'bg-success'
  }

  const getStatusBadgeClass = (status: string) => {
    const classes: Record<string, string> = {
      uploaded: 'bg-success',
      processing: 'bg-warning',
      processed: 'bg-info',
      failed: 'bg-danger'
    }
    return classes[status] || 'bg-secondary'
  }

  // Reset store
  const reset = () => {
    images.value = []
    currentImage.value = null
    searchResults.value = []
    stats.value = null
    categories.value = []
    extensions.value = []
    mimeTypes.value = []
    loading.value = false
    uploading.value = false
    updating.value = false
    error.value = null
    successMessage.value = null
    currentPage.value = 1
    totalPages.value = 0
    total.value = 0
    searchQuery.value = ''
    sortBy.value = 'created_at'
    clearFilters()
  }

  return {
    // State
    images,
    currentImage,
    searchResults,
    stats,
    categories,
    extensions,
    mimeTypes,
    loading,
    uploading,
    updating,
    error,
    successMessage,
    currentPage,
    pageSize,
    totalPages,
    total,
    filters,
    searchQuery,
    sortBy,
    
    // Getters
    hasImages,
    hasSearchResults,
    totalPagesCount,
    isUploading,
    isUpdating,
    hasStats,
    
    // Actions
    clearError,
    clearSuccess,
    clearMessages,
    uploadImage,
    updateImage,
    deleteImage,
    getImage,
    listImages,
    searchImages,
    getPublicImages,
    getSharedImages,
    downloadImage,
    getImageThumbnail,
    makeImagePublic,
    makeImagePrivate,
    shareImage,
    unshareImage,
    getImageStats,
    getCategories,
    getExtensions,
    getMimeTypes,
    goToPage,
    nextPage,
    previousPage,
    setFilters,
    clearFilters,
    setSearchQuery,
    setSortBy,
    formatFileSize,
    getSourceBadgeClass,
    getStatusBadgeClass,
    reset
  }
})
