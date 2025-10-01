import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './authStore'
import { getApiUrl } from '@/utils/apiConfig'

// Types
export interface Post {
  id: string
  title: string
  slug: string
  content: string
  excerpt: string
  status: 'draft' | 'pending' | 'published' | 'private' | 'trash' | 'scheduled'
  type: 'post' | 'page' | 'attachment' | 'revision' | 'custom'
  format: 'standard' | 'aside' | 'gallery' | 'link' | 'image' | 'quote' | 'status' | 'video' | 'audio' | 'chat'
  password?: string
  meta_title: string
  meta_description: string
  meta_keywords: string
  featured_image: string
  categories: string[]
  tags: string[]
  published_at?: string
  scheduled_for?: string
  view_count: number
  comment_count: number
  parent_id?: string
  menu_order: number
  is_sticky: boolean
  allow_pings: boolean
  comment_status: 'open' | 'closed' | 'registration_required'
  language: string
  custom_fields?: Record<string, any>
  created_by: string
  created_at: string
  updated_by: string
  updated_at: string
}

export interface CreatePostRequest {
  title: string
  slug?: string
  content: string
  excerpt?: string
  status?: 'draft' | 'pending' | 'published' | 'private' | 'scheduled'
  type?: 'post' | 'page'
  format?: string
  password?: string
  meta_title?: string
  meta_description?: string
  meta_keywords?: string
  featured_image?: string
  categories?: string[]
  tags?: string[]
  parent_id?: string
  menu_order?: number
  is_sticky?: boolean
  allow_pings?: boolean
  comment_status?: string
  scheduled_for?: string
  custom_fields?: Record<string, any>
}

export interface UpdatePostRequest extends Partial<CreatePostRequest> {}

export interface PostListResponse {
  posts: Post[]
  total: number
  page: number
  limit: number
}

export const usePostStore = defineStore('post', () => {
  // State
  const posts = ref<Post[]>([])
  const currentPost = ref<Post | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const totalPosts = ref(0)
  const currentPage = ref(1)
  const postsPerPage = ref(20)
  
  // Published posts for public view
  const publishedPosts = ref<Post[]>([])
  const popularPosts = ref<Post[]>([])
  const recentPosts = ref<Post[]>([])
  const stickyPosts = ref<Post[]>([])

  // Getters
  const authStore = useAuthStore()
  const totalPages = computed(() => Math.ceil(totalPosts.value / postsPerPage.value))
  const hasPosts = computed(() => posts.value.length > 0)
  const hasPublishedPosts = computed(() => publishedPosts.value.length > 0)

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

  // Actions - Admin API (authenticated)
  const fetchPosts = async (page = 1, status?: string, type?: string) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: postsPerPage.value.toString(),
      })
      
      if (status) params.append('status', status)
      if (type) params.append('type', type)
      
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts?${params.toString()}`),
        {
          headers: getHeaders(),
        }
      )
      
      if (!response.ok) {
        throw new Error(`Failed to fetch posts: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      posts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch posts'
    } finally {
      loading.value = false
    }
  }

  const searchPosts = async (query: string, page = 1, status?: string, type?: string) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        q: query,
        page: page.toString(),
        limit: postsPerPage.value.toString(),
      })
      
      if (status) params.append('status', status)
      if (type) params.append('type', type)
      
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/search?${params.toString()}`),
        {
          headers: getHeaders(),
        }
      )
      
      if (!response.ok) {
        throw new Error(`Failed to search posts: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      posts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to search posts'
    } finally {
      loading.value = false
    }
  }

  const fetchPost = async (id: string) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}`),
        {
          headers: getHeaders(),
        }
      )
      
      if (!response.ok) {
        throw new Error(`Failed to fetch post: ${response.statusText}`)
      }
      
      currentPost.value = await response.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch post'
    } finally {
      loading.value = false
    }
  }

  const createPost = async (postData: CreatePostRequest): Promise<Post> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl('/api/v1/admin/posts'),
        {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(postData),
        }
      )
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Failed to create post: ${response.statusText}`)
      }
      
      const newPost: Post = await response.json()
      posts.value.unshift(newPost)
      totalPosts.value++
      return newPost
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create post'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updatePost = async (id: string, postData: UpdatePostRequest): Promise<Post> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}`),
        {
          method: 'PUT',
          headers: getHeaders(),
          body: JSON.stringify(postData),
        }
      )
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Failed to update post: ${response.statusText}`)
      }
      
      const updatedPost: Post = await response.json()
      const index = posts.value.findIndex(p => p.id === id)
      if (index !== -1) {
        posts.value[index] = updatedPost
      }
      if (currentPost.value?.id === id) {
        currentPost.value = updatedPost
      }
      return updatedPost
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update post'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deletePost = async (id: string) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}`),
        {
          method: 'DELETE',
          headers: getHeaders(),
        }
      )
      
      if (!response.ok) {
        throw new Error(`Failed to delete post: ${response.statusText}`)
      }
      
      posts.value = posts.value.filter(p => p.id !== id)
      totalPosts.value--
      if (currentPost.value?.id === id) {
        currentPost.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete post'
      throw err
    } finally {
      loading.value = false
    }
  }

  const publishPost = async (id: string): Promise<Post> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}/publish`),
        {
          method: 'POST',
          headers: getHeaders(),
        }
      )
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Failed to publish post: ${response.statusText}`)
      }
      
      const publishedPost: Post = await response.json()
      const index = posts.value.findIndex(p => p.id === id)
      if (index !== -1) {
        posts.value[index] = publishedPost
      }
      if (currentPost.value?.id === id) {
        currentPost.value = publishedPost
      }
      return publishedPost
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to publish post'
      throw err
    } finally {
      loading.value = false
    }
  }

  const schedulePost = async (id: string, scheduledFor: string): Promise<Post> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}/schedule`),
        {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify({ scheduled_for: scheduledFor }),
        }
      )
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Failed to schedule post: ${response.statusText}`)
      }
      
      const scheduledPost: Post = await response.json()
      const index = posts.value.findIndex(p => p.id === id)
      if (index !== -1) {
        posts.value[index] = scheduledPost
      }
      if (currentPost.value?.id === id) {
        currentPost.value = scheduledPost
      }
      return scheduledPost
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to schedule post'
      throw err
    } finally {
      loading.value = false
    }
  }

  const refinePost = async (id: string, refineData: any): Promise<Post> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(
        getApiUrl(`/api/v1/admin/posts/${id}/refine`),
        {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(refineData),
        }
      )
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `Failed to refine post: ${response.statusText}`)
      }
      
      const refinedPost: Post = await response.json()
      
      // Update the post in the posts array if it exists
      const index = posts.value.findIndex(p => p.id === id)
      if (index !== -1) {
        posts.value[index] = refinedPost
      }
      
      // Update current post if it's the same one
      if (currentPost.value?.id === id) {
        currentPost.value = refinedPost
      }
      
      return refinedPost
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to refine post'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Actions - Public API (no authentication)
  const fetchPublishedPosts = async (page = 1, type = 'post') => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '10',
        type,
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/published?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error(`Failed to fetch published posts: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      publishedPosts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch published posts'
    } finally {
      loading.value = false
    }
  }

  const fetchPostById = async (id: string): Promise<Post | null> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl(`/api/v1/admin/posts/${id}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        if (response.status === 404) {
          return null
        }
        throw new Error(`Failed to fetch post: ${response.statusText}`)
      }
      
      const responseData = await response.json()
      const post: Post = responseData.post || responseData
      currentPost.value = post
      return post
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch post'
      return null
    } finally {
      loading.value = false
    }
  }

  const fetchPublishedPostBySlug = async (slug: string): Promise<Post | null> => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(getApiUrl(`/api/v1/posts/published/${slug}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        if (response.status === 404) {
          return null
        }
        throw new Error(`Failed to fetch post: ${response.statusText}`)
      }
      
      const responseData = await response.json()
      const post: Post = responseData.post || responseData
      currentPost.value = post
      return post
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch post'
      return null
    } finally {
      loading.value = false
    }
  }

  const fetchPostsByCategory = async (category: string, page = 1) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '10',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/category/${category}?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error(`Failed to fetch posts by category: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      publishedPosts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch posts by category'
    } finally {
      loading.value = false
    }
  }

  const fetchPostsByTag = async (tag: string, page = 1) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '10',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/tag/${tag}?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error(`Failed to fetch posts by tag: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      publishedPosts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch posts by tag'
    } finally {
      loading.value = false
    }
  }

  const searchPublishedPosts = async (query: string, page = 1) => {
    loading.value = true
    error.value = null
    
    try {
      const params = new URLSearchParams({
        q: query,
        page: page.toString(),
        limit: '10',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/search?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (!response.ok) {
        throw new Error(`Failed to search posts: ${response.statusText}`)
      }
      
      const data: PostListResponse = await response.json()
      publishedPosts.value = data.posts
      totalPosts.value = data.total
      currentPage.value = data.page
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to search posts'
    } finally {
      loading.value = false
    }
  }

  const fetchPopularPosts = async (limit = 5, days = 30) => {
    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        days: days.toString(),
        type: 'post',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/popular?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (response.ok) {
        const data = await response.json()
        popularPosts.value = data.posts || []
      }
    } catch (err) {
      console.warn('Failed to fetch popular posts:', err)
    }
  }

  const fetchRecentPosts = async (limit = 5) => {
    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        type: 'post',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/recent?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (response.ok) {
        const data = await response.json()
        recentPosts.value = data.posts || []
      }
    } catch (err) {
      console.warn('Failed to fetch recent posts:', err)
    }
  }

  const fetchStickyPosts = async () => {
    try {
      const params = new URLSearchParams({
        type: 'post',
      })
      
      const response = await fetch(getApiUrl(`/api/v1/posts/sticky?${params.toString()}`), {
        headers: getHeaders()
      })
      
      if (response.ok) {
        const data = await response.json()
        stickyPosts.value = data.posts || []
      }
    } catch (err) {
      console.warn('Failed to fetch sticky posts:', err)
    }
  }

  // Utility actions
  const clearError = () => {
    error.value = null
  }

  const clearCurrentPost = () => {
    currentPost.value = null
  }

  const setCurrentPage = (page: number) => {
    currentPage.value = page
  }

  return {
    // State
    posts,
    currentPost,
    loading,
    error,
    totalPosts,
    currentPage,
    postsPerPage,
    publishedPosts,
    popularPosts,
    recentPosts,
    stickyPosts,
    
    // Getters
    totalPages,
    hasPosts,
    hasPublishedPosts,
    
    // Admin actions
    fetchPosts,
    searchPosts,
    fetchPost,
    fetchPostById,
    createPost,
    updatePost,
    deletePost,
    publishPost,
    schedulePost,
    refinePost,
    
    // Public actions
    fetchPublishedPosts,
    fetchPublishedPostBySlug,
    fetchPostsByCategory,
    fetchPostsByTag,
    searchPublishedPosts,
    fetchPopularPosts,
    fetchRecentPosts,
    fetchStickyPosts,
    
    // Utility actions
    clearError,
    clearCurrentPost,
    setCurrentPage,
  }
})
