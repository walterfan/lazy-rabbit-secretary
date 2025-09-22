<template>
  <div class="wiki-edit-view">
    <!-- Header -->
    <div class="edit-header">
      <div class="container">
        <div class="row align-items-center">
          <div class="col-md-8">
            <h1 class="edit-title">
              <i class="bi bi-pencil"></i>
              {{ isEditing ? $t('wiki.editPage') : $t('wiki.createPage') }}
            </h1>
            <p class="edit-subtitle">
              {{ isEditing ? $t('wiki.editPageDesc') : $t('wiki.createPageDesc') }}
            </p>
          </div>
          <div class="col-md-4 text-end">
            <div class="edit-actions">
              <button 
                class="btn btn-outline-secondary me-2"
                @click="goBack"
              >
                <i class="bi bi-arrow-left"></i>
                {{ $t('common.cancel') }}
              </button>
              <button 
                class="btn btn-primary"
                @click="savePage"
                :disabled="loading || !canSave"
              >
                <i v-if="loading" class="spinner-border spinner-border-sm me-1"></i>
                <i v-else class="bi bi-check-lg me-1"></i>
                {{ $t('common.save') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="container">
      <div class="row">
        <div class="col-lg-8">
          <div class="edit-form">
            <!-- Basic Information -->
            <div class="form-section">
              <h5>{{ $t('wiki.basicInfo') }}</h5>
              
              <div class="row">
                <div class="col-md-8">
                  <div class="mb-3">
                    <label for="title" class="form-label">
                      {{ $t('wiki.title') }} <span class="text-danger">*</span>
                    </label>
                    <input
                      type="text"
                      class="form-control"
                      id="title"
                      v-model="formData.title"
                      :placeholder="$t('wiki.titlePlaceholder')"
                      :class="{ 'is-invalid': errors.title }"
                    />
                    <div v-if="errors.title" class="invalid-feedback">
                      {{ errors.title }}
                    </div>
                  </div>
                </div>
                <div class="col-md-4">
                  <div class="mb-3">
                    <label for="slug" class="form-label">
                      {{ $t('wiki.slug') }}
                    </label>
                    <input
                      type="text"
                      class="form-control"
                      id="slug"
                      v-model="formData.slug"
                      :placeholder="$t('wiki.slugPlaceholder')"
                      :class="{ 'is-invalid': errors.slug }"
                    />
                    <div v-if="errors.slug" class="invalid-feedback">
                      {{ errors.slug }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="mb-3">
                <label for="summary" class="form-label">
                  {{ $t('wiki.summary') }}
                </label>
                <textarea
                  class="form-control"
                  id="summary"
                  v-model="formData.summary"
                  :placeholder="$t('wiki.summaryPlaceholder')"
                  rows="2"
                ></textarea>
              </div>
            </div>

            <!-- Content -->
            <div class="form-section">
              <h5>{{ $t('wiki.content') }}</h5>
              
              <div class="mb-3">
                <label for="content" class="form-label">
                  {{ $t('wiki.content') }} <span class="text-danger">*</span>
                </label>
                <textarea
                  class="form-control"
                  id="content"
                  v-model="formData.content"
                  :placeholder="$t('wiki.contentPlaceholder')"
                  rows="15"
                  :class="{ 'is-invalid': errors.content }"
                ></textarea>
                <div v-if="errors.content" class="invalid-feedback">
                  {{ errors.content }}
                </div>
              </div>
            </div>

            <!-- Change Note -->
            <div class="form-section">
              <h5>{{ $t('wiki.changeNote') }}</h5>
              
              <div class="mb-3">
                <label for="changeNote" class="form-label">
                  {{ $t('wiki.changeNoteDesc') }}
                </label>
                <input
                  type="text"
                  class="form-control"
                  id="changeNote"
                  v-model="formData.change_note"
                  :placeholder="$t('wiki.changeNotePlaceholder')"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="col-lg-4">
          <div class="edit-sidebar">
            <!-- Page Settings -->
            <div class="sidebar-section">
              <h5>{{ $t('wiki.pageSettings') }}</h5>
              
              <div class="mb-3">
                <label for="status" class="form-label">
                  {{ $t('wiki.status') }}
                </label>
                <select class="form-select" id="status" v-model="formData.status">
                  <option value="draft">{{ $t('wiki.status.draft') }}</option>
                  <option value="published">{{ $t('wiki.status.published') }}</option>
                  <option value="protected">{{ $t('wiki.status.protected') }}</option>
                </select>
              </div>

              <div class="mb-3">
                <label for="type" class="form-label">
                  {{ $t('wiki.type') }}
                </label>
                <select class="form-select" id="type" v-model="formData.type">
                  <option value="article">{{ $t('wiki.type.article') }}</option>
                  <option value="template">{{ $t('wiki.type.template') }}</option>
                  <option value="category">{{ $t('wiki.type.category') }}</option>
                  <option value="stub">{{ $t('wiki.type.stub') }}</option>
                </select>
              </div>

              <div class="mb-3">
                <label for="language" class="form-label">
                  {{ $t('wiki.language') }}
                </label>
                <select class="form-select" id="language" v-model="formData.language">
                  <option value="en">English</option>
                  <option value="zh">中文</option>
                </select>
              </div>

              <div class="form-check mb-3">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="isProtected"
                  v-model="formData.is_protected"
                />
                <label class="form-check-label" for="isProtected">
                  {{ $t('wiki.isProtected') }}
                </label>
              </div>
            </div>

            <!-- Categories and Tags -->
            <div class="sidebar-section">
              <h5>{{ $t('wiki.organization') }}</h5>
              
              <div class="mb-3">
                <label for="categories" class="form-label">
                  {{ $t('wiki.categories') }}
                </label>
                <input
                  type="text"
                  class="form-control"
                  id="categories"
                  v-model="categoriesInput"
                  :placeholder="$t('wiki.categoriesPlaceholder')"
                />
                <div class="form-text">{{ $t('wiki.categoriesHelp') }}</div>
              </div>

              <div class="mb-3">
                <label for="tags" class="form-label">
                  {{ $t('wiki.tags') }}
                </label>
                <input
                  type="text"
                  class="form-control"
                  id="tags"
                  v-model="tagsInput"
                  :placeholder="$t('wiki.tagsPlaceholder')"
                />
                <div class="form-text">{{ $t('wiki.tagsHelp') }}</div>
              </div>
            </div>

            <!-- Preview -->
            <div class="sidebar-section">
              <h5>{{ $t('wiki.preview') }}</h5>
              <div class="preview-content">
                <div v-if="formData.title" class="preview-title">
                  {{ formData.title }}
                </div>
                <div v-if="formData.summary" class="preview-summary">
                  {{ formData.summary }}
                </div>
                <div v-if="renderedPreview" class="preview-markdown" v-html="renderedPreview"></div>
                <div v-if="!formData.title && !formData.summary && !renderedPreview" class="text-muted">
                  {{ $t('wiki.noPreview') }}
                </div>
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
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useWikiStore } from '@/stores/wikiStore'
import { useAuthStore } from '@/stores/authStore'
import { marked } from 'marked'
import type { CreateWikiPageRequest } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const wikiStore = useWikiStore()
const authStore = useAuthStore()

// State
const loading = ref(false)
const errors = ref<Record<string, string>>({})

// Form data
const formData = ref<CreateWikiPageRequest>({
  title: '',
  content: '',
  summary: '',
  status: 'draft',
  type: 'article',
  is_protected: false,
  language: 'en',
  change_note: ''
})

const categoriesInput = ref('')
const tagsInput = ref('')

// Computed
const isEditing = computed(() => !!route.params.slug)
const isAuthenticated = computed(() => authStore.isAuthenticated)

const canSave = computed(() => {
  return formData.value.title.trim() && formData.value.content.trim()
})

const renderedPreview = computed(() => {
  if (!formData.value.content) return ''
  
  // Configure marked options
  marked.setOptions({
    breaks: true,
    gfm: true,
  })
  
  // Render markdown to HTML
  return marked(formData.value.content)
})

// Watch for changes to update categories and tags
watch(categoriesInput, (newValue) => {
  formData.value.categories = newValue.split(',').map(c => c.trim()).filter(c => c)
})

watch(tagsInput, (newValue) => {
  formData.value.tags = newValue.split(',').map(t => t.trim()).filter(t => t)
})

// Methods
const loadPage = async () => {
  if (!isEditing.value) return

  const slug = route.params.slug as string
  loading.value = true

  try {
    const page = await wikiStore.fetchPageBySlug(slug)
    if (page) {
      formData.value = {
        title: page.title,
        content: page.content,
        summary: page.summary,
        status: page.status,
        type: page.type,
        is_protected: page.is_protected,
        language: page.language,
        change_note: ''
      }
      
      categoriesInput.value = page.categories.join(', ')
      tagsInput.value = page.tags.join(', ')
    }
  } catch (error) {
    console.error('Failed to load page:', error)
    // If page doesn't exist, we'll create a new one with the slug
    if (route.params.slug) {
      formData.value.slug = route.params.slug as string
      // Generate a title from the slug
      formData.value.title = (route.params.slug as string)
        .split('-')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ')
    }
  } finally {
    loading.value = false
  }
}

const savePage = async () => {
  if (!canSave.value) return

  loading.value = true
  errors.value = {}

  try {
    if (isEditing.value) {
      const slug = route.params.slug as string
      try {
        const page = await wikiStore.fetchPageBySlug(slug)
        if (page) {
          await wikiStore.updatePage(page.id, formData.value)
        }
      } catch (error) {
        // If page doesn't exist, create a new one
        await wikiStore.createPage(formData.value)
      }
    } else {
      await wikiStore.createPage(formData.value)
    }
    
    // Redirect to the created/updated page
    const redirectSlug = formData.value.slug || route.params.slug
    if (redirectSlug) {
      router.push(`/wiki/page/${redirectSlug}`)
    } else {
      router.push('/wiki')
    }
  } catch (error) {
    console.error('Failed to save page:', error)
    errors.value = { general: 'Failed to save page' }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

// Lifecycle
onMounted(() => {
  if (!isAuthenticated.value) {
    router.push('/signin')
    return
  }
  
  loadPage()
})
</script>

<style scoped>
.wiki-edit-view {
  min-height: 100vh;
  background-color: #f8f9fa;
}

.edit-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem 0;
  margin-bottom: 2rem;
}

.edit-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.edit-subtitle {
  font-size: 1rem;
  opacity: 0.9;
  margin-bottom: 0;
}

.edit-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-form {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  padding: 2rem;
  margin-bottom: 2rem;
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #dee2e6;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.form-section h5 {
  margin-bottom: 1.5rem;
  color: #333;
  font-weight: 600;
}

.edit-sidebar {
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

.preview-content {
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #dee2e6;
}

.preview-title {
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
}

.preview-summary {
  color: #6c757d;
  font-size: 0.9rem;
  line-height: 1.4;
}

.preview-markdown {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #dee2e6;
  line-height: 1.6;
  color: #333;
}

.preview-markdown h1,
.preview-markdown h2,
.preview-markdown h3,
.preview-markdown h4,
.preview-markdown h5,
.preview-markdown h6 {
  margin-top: 1rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.preview-markdown h1:first-child,
.preview-markdown h2:first-child,
.preview-markdown h3:first-child {
  margin-top: 0;
}

.preview-markdown p {
  margin-bottom: 0.75rem;
}

.preview-markdown ul,
.preview-markdown ol {
  margin-bottom: 0.75rem;
  padding-left: 1.5rem;
}

.preview-markdown blockquote {
  border-left: 4px solid #667eea;
  padding-left: 0.75rem;
  margin: 0.75rem 0;
  color: #6c757d;
}

.preview-markdown code {
  background-color: #f8f9fa;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
}

.preview-markdown pre {
  background-color: #f8f9fa;
  padding: 0.75rem;
  border-radius: 6px;
  overflow-x: auto;
  font-size: 0.875rem;
}

.preview-markdown a {
  color: #667eea;
  text-decoration: none;
}

.preview-markdown a:hover {
  color: #5a67d8;
  text-decoration: underline;
}

.preview-markdown table {
  width: 100%;
  border-collapse: collapse;
  margin: 0.75rem 0;
  font-size: 0.875rem;
}

.preview-markdown th,
.preview-markdown td {
  border: 1px solid #dee2e6;
  padding: 0.5rem;
  text-align: left;
}

.preview-markdown th {
  background-color: #f8f9fa;
  font-weight: 600;
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
  .edit-header {
    padding: 1.5rem 0;
  }
  
  .edit-title {
    font-size: 1.5rem;
  }
  
  .edit-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .edit-form {
    padding: 1rem;
  }
  
  .sidebar-section {
    padding: 1rem;
  }
}
</style>
