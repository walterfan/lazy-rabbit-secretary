<template>
  <div 
    class="modal fade show d-block"
    tabindex="-1"
    style="background-color: rgba(0,0,0,0.5)"
  >
    <div class="modal-dialog modal-xl">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            {{ isEditing ? 'Edit' : 'Create' }} {{ isPage ? 'Page' : 'Post' }}
          </h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        
        <form @submit.prevent="handleSubmit">
          <div class="modal-body">
            <!-- Error Alert -->
            <div v-if="formErrors.length > 0" class="alert alert-danger">
              <ul class="mb-0">
                <li v-for="error in formErrors" :key="error">{{ error }}</li>
              </ul>
            </div>

            <div class="row">
              <!-- Main Content Column -->
              <div class="col-lg-8">
                <!-- Title -->
                <div class="mb-3">
                  <label class="form-label">Title *</label>
                  <input
                    v-model="postData.title"
                    type="text"
                    class="form-control"
                    :class="{ 'is-invalid': v$.title.$error }"
                    placeholder="Enter post title..."
                    @input="generateSlugFromTitle"
                  >
                  <div v-if="v$.title.$error" class="invalid-feedback">
                    Title is required
                  </div>
                </div>

                <!-- Slug -->
                <div class="mb-3">
                  <label class="form-label">URL Slug</label>
                  <div class="input-group">
                    <span class="input-group-text">/blog/</span>
                    <input
                      v-model="postData.slug"
                      type="text"
                      class="form-control"
                      placeholder="url-slug"
                    >
                  </div>
                  <div class="form-text">
                    Leave empty to auto-generate from title
                  </div>
                </div>

                <!-- Content -->
                <div class="mb-3">
                  <div class="d-flex justify-content-between align-items-center mb-2">
                    <label class="form-label mb-0">Content * <small class="text-muted">(Markdown supported)</small></label>
                    
                    <!-- AI Assistant Button -->
                    <div class="dropdown">
                      <button
                        class="btn btn-outline-primary btn-sm dropdown-toggle"
                        type="button"
                        data-bs-toggle="dropdown"
                        :disabled="!postData.content?.trim() || aiProcessing"
                      >
                        <i class="bi bi-robot me-1"></i>
                        <span v-if="aiProcessing">
                          <span class="spinner-border spinner-border-sm me-1"></span>
                          Processing...
                        </span>
                        <span v-else>AI Assist</span>
                      </button>
                      <ul class="dropdown-menu dropdown-menu-end">
                        <li>
                          <h6 class="dropdown-header">
                            <i class="bi bi-pencil-square me-1"></i>
                            Writing Improvements
                          </h6>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('improve-writing')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-arrow-up-circle me-2 text-success"></i>
                            Improve Writing
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('make-shorter')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-arrow-down-circle me-2 text-info"></i>
                            Make Shorter
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('make-longer')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-arrow-up-circle me-2 text-warning"></i>
                            Make Longer
                          </button>
                        </li>
                        <li><hr class="dropdown-divider"></li>
                        <li>
                          <h6 class="dropdown-header">
                            <i class="bi bi-palette me-1"></i>
                            Change Tone
                          </h6>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('tone-formal')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-briefcase me-2 text-secondary"></i>
                            Formal
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('tone-informal')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-chat-heart me-2 text-primary"></i>
                            Informal
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('tone-friendly')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-emoji-smile me-2 text-success"></i>
                            Friendly
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('tone-persuasive')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-megaphone me-2 text-warning"></i>
                            Persuasive
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('tone-serious')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-shield-check me-2 text-danger"></i>
                            Serious
                          </button>
                        </li>
                        <li><hr class="dropdown-divider"></li>
                        <li>
                          <h6 class="dropdown-header">
                            <i class="bi bi-translate me-1"></i>
                            Translation
                          </h6>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('translate-english')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-globe me-2 text-info"></i>
                            Translate to English
                          </button>
                        </li>
                        <li>
                          <button 
                            class="dropdown-item" 
                            @click="handleAiAction('translate-chinese')"
                            :disabled="aiProcessing"
                          >
                            <i class="bi bi-globe2 me-2 text-info"></i>
                            Translate to Chinese
                          </button>
                        </li>
                      </ul>
                    </div>
                  </div>

                  <MarkdownEditor
                    v-model="postData.content"
                    :textarea-class="{ 'form-control': true, 'is-invalid': v$.content.$error }"
                    :rows="12"
                    placeholder="Write your post content in Markdown..."
                  />
                  <div v-if="v$.content.$error" class="invalid-feedback d-block">
                    Content is required
                  </div>
                </div>

                <!-- Excerpt -->
                <div class="mb-3">
                  <label class="form-label">Excerpt</label>
                  <textarea
                    v-model="postData.excerpt"
                    class="form-control"
                    rows="3"
                    placeholder="Optional short summary..."
                  ></textarea>
                  <div class="form-text">
                    Brief description for post previews and SEO
                  </div>
                </div>
              </div>

              <!-- Sidebar Column -->
              <div class="col-lg-4">
                <!-- Publish Box -->
                <div class="card mb-3">
                  <div class="card-header">
                    <h6 class="mb-0">Publish</h6>
                  </div>
                  <div class="card-body">
                    <!-- Status -->
                    <div class="mb-3">
                      <label class="form-label">Status</label>
                      <select v-model="postData.status" class="form-select">
                        <option value="draft">Draft</option>
                        <option value="pending">Pending Review</option>
                        <option value="published">Published</option>
                        <option value="private">Private</option>
                        <option value="scheduled">Scheduled</option>
                      </select>
                    </div>

                    <!-- Scheduled Date (only if status is scheduled) -->
                    <div v-if="postData.status === 'scheduled'" class="mb-3">
                      <label class="form-label">Publish Date</label>
                      <input
                        v-model="scheduledDate"
                        type="datetime-local"
                        class="form-control"
                      >
                    </div>

                    <!-- Visibility Options -->
                    <div class="mb-3">
                      <label class="form-label">Visibility</label>
                      <div class="form-check">
                        <input
                          id="sticky"
                          v-model="postData.is_sticky"
                          type="checkbox"
                          class="form-check-input"
                        >
                        <label for="sticky" class="form-check-label">
                          Stick to front page
                        </label>
                      </div>
                    </div>

                    <!-- Password Protection -->
                    <div class="mb-3">
                      <label class="form-label">Password</label>
                      <input
                        v-model="postData.password"
                        type="password"
                        class="form-control"
                        placeholder="Leave empty for public"
                      >
                      <div class="form-text">
                        Set password to protect this post
                      </div>
                    </div>

                    <!-- Comments -->
                    <div class="mb-3">
                      <label class="form-label">Comments</label>
                      <select v-model="postData.comment_status" class="form-select">
                        <option value="open">Open</option>
                        <option value="closed">Closed</option>
                        <option value="registration_required">Registration Required</option>
                      </select>
                    </div>
                  </div>
                </div>

                <!-- Categories & Tags -->
                <div class="card mb-3">
                  <div class="card-header">
                    <h6 class="mb-0">Categories & Tags</h6>
                  </div>
                  <div class="card-body">
                    <!-- Categories -->
                    <div class="mb-3">
                      <label class="form-label">Categories</label>
                      <input
                        v-model="categoriesInput"
                        type="text"
                        class="form-control"
                        placeholder="technology, web-development"
                      >
                      <div class="form-text">
                        Separate categories with commas
                      </div>
                      <div v-if="postData.categories && postData.categories.length > 0" class="mt-2">
                        <span 
                          v-for="category in postData.categories" 
                          :key="category"
                          class="badge bg-secondary me-1"
                        >
                          {{ category }}
                        </span>
                      </div>
                    </div>

                    <!-- Tags -->
                    <div class="mb-3">
                      <label class="form-label">Tags</label>
                      <input
                        v-model="tagsInput"
                        type="text"
                        class="form-control"
                        placeholder="vue, javascript, tutorial"
                      >
                      <div class="form-text">
                        Separate tags with commas
                      </div>
                      <div v-if="postData.tags && postData.tags.length > 0" class="mt-2">
                        <span 
                          v-for="tag in postData.tags" 
                          :key="tag"
                          class="badge bg-primary me-1"
                        >
                          {{ tag }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- SEO Settings -->
                <div class="card mb-3">
                  <div class="card-header">
                    <h6 class="mb-0">SEO Settings</h6>
                  </div>
                  <div class="card-body">
                    <!-- Meta Title -->
                    <div class="mb-3">
                      <label class="form-label">Meta Title</label>
                      <input
                        v-model="postData.meta_title"
                        type="text"
                        class="form-control"
                        placeholder="SEO title (leave empty to use post title)"
                      >
                      <div class="form-text">
                        {{ (postData.meta_title || postData.title).length }}/60 characters
                      </div>
                    </div>

                    <!-- Meta Description -->
                    <div class="mb-3">
                      <label class="form-label">Meta Description</label>
                      <textarea
                        v-model="postData.meta_description"
                        class="form-control"
                        rows="3"
                        placeholder="Brief description for search engines"
                      ></textarea>
                      <div class="form-text">
                        {{ (postData.meta_description || '').length }}/160 characters
                      </div>
                    </div>

                    <!-- Keywords -->
                    <div class="mb-3">
                      <label class="form-label">Keywords</label>
                      <input
                        v-model="postData.meta_keywords"
                        type="text"
                        class="form-control"
                        placeholder="keyword1, keyword2, keyword3"
                      >
                    </div>
                  </div>
                </div>

                <!-- Featured Image -->
                <div class="card mb-3">
                  <div class="card-header">
                    <h6 class="mb-0">Featured Image</h6>
                  </div>
                  <div class="card-body">
                    <div class="mb-3">
                      <label class="form-label">Image URL</label>
                      <input
                        v-model="postData.featured_image"
                        type="url"
                        class="form-control"
                        placeholder="https://example.com/image.jpg"
                      >
                    </div>
                    <div v-if="postData.featured_image" class="text-center">
                      <img 
                        :src="postData.featured_image" 
                        alt="Featured image preview"
                        class="img-fluid rounded"
                        style="max-height: 150px;"
                        @error="onImageError"
                      >
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="$emit('close')">
              Cancel
            </button>
            <button type="button" class="btn btn-outline-primary" @click="saveDraft">
              Save as Draft
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="postStore.loading || v$.$invalid"
            >
              <span v-if="postStore.loading" class="spinner-border spinner-border-sm me-2"></span>
              {{ isEditing ? 'Update' : 'Create' }} {{ isPage ? 'Page' : 'Post' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength } from '@vuelidate/validators'
import { usePostStore, type Post, type CreatePostRequest } from '@/stores/postStore'
import MarkdownEditor from '@/components/MarkdownEditor.vue'

// Props
interface Props {
  show: boolean
  post?: Post | null
  isPage?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  post: null,
  isPage: false,
})

// Emits
const emit = defineEmits<{
  close: []
  saved: [post: Post]
}>()

// Store
const postStore = usePostStore()

// Reactive state
const postData = reactive<CreatePostRequest>({
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  status: 'draft',
  type: props.isPage ? 'page' : 'post',
  format: 'standard',
  password: '',
  meta_title: '',
  meta_description: '',
  meta_keywords: '',
  featured_image: '',
  categories: [],
  tags: [],
  parent_id: '',
  menu_order: 0,
  is_sticky: false,
  allow_pings: true,
  comment_status: 'open',
  custom_fields: {},
})

const categoriesInput = ref('')
const tagsInput = ref('')
const scheduledDate = ref('')
const formErrors = ref<string[]>([])
const aiProcessing = ref(false)

// Validation rules
const rules = {
  title: { required, minLength: minLength(1) },
  content: { required, minLength: minLength(1) },
}

const v$ = useVuelidate(rules, postData)

// Computed
const isEditing = computed(() => !!props.post?.id)

// Watchers
watch(
  () => categoriesInput.value,
  (newValue) => {
    postData.categories = newValue
      .split(',')
      .map(cat => cat.trim())
      .filter(cat => cat.length > 0)
  }
)

watch(
  () => tagsInput.value,
  (newValue) => {
    postData.tags = newValue
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0)
  }
)

watch(
  () => scheduledDate.value,
  (newValue) => {
    postData.scheduled_for = newValue ? new Date(newValue).toISOString() : ''
  }
)

// Lifecycle
onMounted(() => {
  if (props.post) {
    // Populate form with existing post data
    Object.assign(postData, {
      title: props.post.title,
      slug: props.post.slug,
      content: props.post.content,
      excerpt: props.post.excerpt,
      status: props.post.status,
      type: props.post.type,
      format: props.post.format,
      password: props.post.password || '',
      meta_title: props.post.meta_title,
      meta_description: props.post.meta_description,
      meta_keywords: props.post.meta_keywords,
      featured_image: props.post.featured_image,
      categories: [...props.post.categories],
      tags: [...props.post.tags],
      parent_id: props.post.parent_id || '',
      menu_order: props.post.menu_order,
      is_sticky: props.post.is_sticky,
      allow_pings: props.post.allow_pings,
      comment_status: props.post.comment_status,
      custom_fields: props.post.custom_fields || {},
    })

    // Set input fields
    categoriesInput.value = props.post.categories.join(', ')
    tagsInput.value = props.post.tags.join(', ')
    
    if (props.post.scheduled_for) {
      scheduledDate.value = new Date(props.post.scheduled_for).toISOString().slice(0, 16)
    }
  }
})

// Methods
const generateSlugFromTitle = () => {
  if (!postData.slug && postData.title) {
    postData.slug = postData.title
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim()
  }
}

const saveDraft = async () => {
  const originalStatus = postData.status
  postData.status = 'draft'
  await handleSubmit()
  postData.status = originalStatus
}

const handleSubmit = async () => {
  formErrors.value = []
  
  // Validate form
  const isValid = await v$.value.$validate()
  if (!isValid) {
    formErrors.value = ['Please fix the validation errors above.']
    return
  }

  try {
    let savedPost: Post

    // Prepare the data for submission
    const submitData = { ...postData }
    
    // Handle scheduled_for field - only include it if we have a valid date
    if (postData.status === 'scheduled' && scheduledDate.value) {
      submitData.scheduled_for = new Date(scheduledDate.value).toISOString()
    } else {
      // Remove scheduled_for field if not needed to avoid empty string issues
      delete (submitData as any).scheduled_for
    }

    // Remove empty optional fields to avoid backend parsing issues
    if (!submitData.parent_id) {
      delete (submitData as any).parent_id
    }
    if (!submitData.password) {
      delete (submitData as any).password
    }

    if (isEditing.value && props.post) {
      // Update existing post
      savedPost = await postStore.updatePost(props.post.id, submitData)
    } else {
      // Create new post
      savedPost = await postStore.createPost(submitData)
    }

    emit('saved', savedPost)
  } catch (error) {
    if (error instanceof Error) {
      formErrors.value = [error.message]
    } else {
      formErrors.value = ['An unexpected error occurred. Please try again.']
    }
  }
}

const onImageError = () => {
  formErrors.value = ['Failed to load featured image. Please check the URL.']
}

// AI Assistant Methods
const handleAiAction = async (action: string) => {
  if (!postData.content?.trim()) {
    alert('Please add some content first before using AI assistance.')
    return
  }

  // Check if we're editing an existing post (need post ID for API call)
  if (!isEditing.value || !props.post?.id) {
    alert('Please save the post first before using AI assistance. AI refinement only works on existing posts.')
    return
  }

  aiProcessing.value = true
  
  try {
    // Map frontend action names to backend action names
    const actionMapping: Record<string, { action: string; requirement?: string }> = {
      'improve-writing': { action: 'improve_writing' },
      'make-shorter': { action: 'make_shorter' },
      'make-longer': { action: 'make_longer' },
      'tone-formal': { action: 'change_tone', requirement: 'formal' },
      'tone-informal': { action: 'change_tone', requirement: 'informal' },
      'tone-friendly': { action: 'change_tone', requirement: 'friendly' },
      'tone-persuasive': { action: 'change_tone', requirement: 'persuasive' },
      'tone-serious': { action: 'change_tone', requirement: 'serious' },
      'translate-english': { action: 'translate', requirement: 'english' },
      'translate-chinese': { action: 'translate', requirement: 'chinese' }
    }

    const mappedAction = actionMapping[action]
    if (!mappedAction) {
      throw new Error(`Unknown action: ${action}`)
    }

    // Get action description for user feedback
    const actionDescriptions: Record<string, string> = {
      'improve-writing': 'improving the writing quality',
      'make-shorter': 'making the content more concise',
      'make-longer': 'expanding the content with more details',
      'tone-formal': 'changing the tone to formal',
      'tone-informal': 'changing the tone to informal',
      'tone-friendly': 'changing the tone to friendly',
      'tone-persuasive': 'changing the tone to persuasive',
      'tone-serious': 'changing the tone to serious',
      'translate-english': 'translating to English',
      'translate-chinese': 'translating to Chinese'
    }

    const actionDescription = actionDescriptions[action] || 'processing'
    console.log(`AI is ${actionDescription}...`)

    // Prepare the refine request with current post data
    const refineRequest = {
      title: postData.title,
      slug: postData.slug,
      content: postData.content,
      excerpt: postData.excerpt,
      status: postData.status,
      type: postData.type,
      format: postData.format,
      password: postData.password,
      meta_title: postData.meta_title,
      meta_description: postData.meta_description,
      meta_keywords: postData.meta_keywords,
      featured_image: postData.featured_image,
      categories: postData.categories,
      tags: postData.tags,
      parent_id: postData.parent_id,
      menu_order: postData.menu_order,
      is_sticky: postData.is_sticky,
      allow_pings: postData.allow_pings,
      comment_status: postData.comment_status,
      scheduled_for: postData.scheduled_for,
      custom_fields: postData.custom_fields,
      action: mappedAction.action,
      requirement: mappedAction.requirement
    }

    // Call the backend API
    const refinedPost = await postStore.refinePost(props.post.id, refineRequest)
    
    if (refinedPost && refinedPost.content !== postData.content) {
      // Show confirmation dialog before applying changes
      const confirmed = confirm(
        `AI has finished ${actionDescription}. Would you like to replace the current content with the AI-processed version?\n\n` +
        `Preview (first 200 characters):\n${refinedPost.content.substring(0, 200)}${refinedPost.content.length > 200 ? '...' : ''}`
      )
      
      if (confirmed) {
        // Update the form data with the refined content
        postData.content = refinedPost.content
        // Also update other fields that might have been refined
        if (refinedPost.title !== postData.title) postData.title = refinedPost.title
        if (refinedPost.excerpt !== postData.excerpt) postData.excerpt = refinedPost.excerpt
        console.log(`Content updated with AI ${actionDescription}`)
      }
    } else {
      alert('AI processing completed, but no changes were suggested.')
    }
    
  } catch (error) {
    console.error('AI processing error:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred'
    alert(`Sorry, AI processing failed: ${errorMessage}. Please try again later.`)
  } finally {
    aiProcessing.value = false
  }
}

</script>

<style scoped>
.modal.show {
  animation: fadeIn 0.15s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.card-header h6 {
  color: #495057;
  font-weight: 600;
}

.form-text {
  font-size: 0.875rem;
}

.badge {
  font-size: 0.75rem;
}

.img-fluid {
  border: 1px solid #dee2e6;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}

.modal-xl {
  max-width: 1200px;
}

.modal-body {
  max-height: 70vh;
  overflow-y: auto;
}
</style>
