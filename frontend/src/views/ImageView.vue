<template>
  <div class="image-view">
    <div class="container-fluid">
      <!-- Header -->
      <div class="row mb-4">
        <div class="col-12">
          <div class="d-flex justify-content-between align-items-center">
            <h2>
              <i class="bi bi-images"></i>
              {{ $t('image.title') }}
            </h2>
            <button @click="showUploadModal = true" class="btn btn-primary">
              <i class="bi bi-cloud-upload"></i>
              {{ $t('image.upload') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Filters and Search -->
      <div class="row mb-4">
        <div class="col-12">
          <div class="card">
            <div class="card-body">
              <div class="row g-3">
                <div class="col-md-3">
                  <label class="form-label">{{ $t('image.search') }}</label>
                  <input 
                    v-model="searchQuery" 
                    @input="onSearch" 
                    type="text" 
                    class="form-control" 
                    :placeholder="$t('image.searchPlaceholder')"
                  >
                </div>
                <div class="col-md-2">
                  <label class="form-label">{{ $t('image.type') }}</label>
                  <select v-model="selectedType" @change="onFilterChange" class="form-select">
                    <option value="">{{ $t('image.allTypes') }}</option>
                    <option value="avatar">{{ $t('image.avatar') }}</option>
                    <option value="profile">{{ $t('image.profile') }}</option>
                    <option value="cover">{{ $t('image.cover') }}</option>
                    <option value="gallery">{{ $t('image.gallery') }}</option>
                    <option value="attachment">{{ $t('image.attachment') }}</option>
                    <option value="custom">{{ $t('image.custom') }}</option>
                  </select>
                </div>
                <div class="col-md-2">
                  <label class="form-label">{{ $t('image.source') }}</label>
                  <select v-model="selectedSource" @change="onFilterChange" class="form-select">
                    <option value="">{{ $t('image.allSources') }}</option>
                    <option value="uploaded">{{ $t('image.uploaded') }}</option>
                    <option value="generated">{{ $t('image.generated') }}</option>
                  </select>
                </div>
                <div class="col-md-2">
                  <label class="form-label">{{ $t('image.status') }}</label>
                  <select v-model="selectedStatus" @change="onFilterChange" class="form-select">
                    <option value="">{{ $t('image.allStatuses') }}</option>
                    <option value="uploaded">{{ $t('image.uploaded') }}</option>
                    <option value="processing">{{ $t('image.processing') }}</option>
                    <option value="processed">{{ $t('image.processed') }}</option>
                    <option value="failed">{{ $t('image.failed') }}</option>
                  </select>
                </div>
                <div class="col-md-2">
                  <label class="form-label">{{ $t('image.sortBy') }}</label>
                  <select v-model="sortBy" @change="onFilterChange" class="form-select">
                    <option value="created_at">{{ $t('image.dateCreated') }}</option>
                    <option value="updated_at">{{ $t('image.dateUpdated') }}</option>
                    <option value="view_count">{{ $t('image.viewCount') }}</option>
                    <option value="file_size">{{ $t('image.fileSize') }}</option>
                  </select>
                </div>
                <div class="col-md-1">
                  <label class="form-label">{{ $t('image.limit') }}</label>
                  <select v-model="limit" @change="onFilterChange" class="form-select">
                    <option value="12">12</option>
                    <option value="24">24</option>
                    <option value="48">48</option>
                    <option value="96">96</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Statistics -->
      <div class="row mb-4" v-if="stats">
        <div class="col-md-3">
          <div class="card text-center">
            <div class="card-body">
              <h5 class="card-title text-primary">{{ stats.total_count }}</h5>
              <p class="card-text">{{ $t('image.totalImages') }}</p>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card text-center">
            <div class="card-body">
              <h5 class="card-title text-success">{{ formatFileSize(stats.total_size || 0) }}</h5>
              <p class="card-text">{{ $t('image.totalSize') }}</p>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card text-center">
            <div class="card-body">
              <h5 class="card-title text-info">{{ formatFileSize(stats.average_size || 0) }}</h5>
              <p class="card-text">{{ $t('image.averageSize') }}</p>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card text-center">
            <div class="card-body">
              <h5 class="card-title text-warning">{{ Object.keys(stats.type_counts || {}).length }}</h5>
              <p class="card-text">{{ $t('image.types') }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Image Grid -->
      <div class="row">
        <div class="col-12">
          <div v-if="loading" class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">{{ $t('image.loading') }}</span>
            </div>
          </div>
          
          <div v-else-if="images.length === 0" class="text-center py-5">
            <i class="bi bi-image display-1 text-muted"></i>
            <h4 class="text-muted mt-3">{{ $t('image.noImages') }}</h4>
            <p class="text-muted">{{ $t('image.noImagesDescription') }}</p>
            <button @click="showUploadModal = true" class="btn btn-primary">
              <i class="bi bi-cloud-upload"></i>
              {{ $t('image.uploadFirst') }}
            </button>
          </div>
          
          <div v-else class="image-grid">
            <div 
              v-for="image in images" 
              :key="image.id" 
              class="image-card"
              @click="selectImage(image)"
            >
              <div class="image-container">
                <img 
                  :src="image.thumbnail_url || image.url" 
                  :alt="image.original_name"
                  class="image-thumbnail"
                  @error="onImageError"
                >
                <div class="image-overlay">
                  <div class="image-actions">
                    <button @click.stop="viewImage(image)" class="btn btn-sm btn-light">
                      <i class="bi bi-eye"></i>
                    </button>
                    <button @click.stop="downloadImage(image)" class="btn btn-sm btn-light">
                      <i class="bi bi-download"></i>
                    </button>
                    <button @click.stop="editImage(image)" class="btn btn-sm btn-light">
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button @click.stop="deleteImage(image)" class="btn btn-sm btn-danger">
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </div>
              </div>
              <div class="image-info">
                <h6 class="image-title" :title="image.original_name">
                  {{ image.original_name }}
                </h6>
                <div class="image-meta">
                  <span class="badge" :class="getSourceBadgeClass(image.source)">
                    {{ image.source }}
                  </span>
                  <span class="badge bg-secondary">{{ image.type }}</span>
                  <small class="text-muted">{{ formatFileSize(image.file_size) }}</small>
                </div>
                <div class="image-stats">
                  <small class="text-muted">
                    <i class="bi bi-eye"></i> {{ image.view_count }}
                    <i class="bi bi-download ms-2"></i> {{ image.download_count }}
                  </small>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="row mt-4" v-if="totalPages > 1">
        <div class="col-12">
          <nav aria-label="Image pagination">
            <ul class="pagination justify-content-center">
              <li class="page-item" :class="{ disabled: currentPage === 1 }">
                <button @click="goToPage(currentPage - 1)" class="page-link">
                  {{ $t('image.previous') }}
                </button>
              </li>
              <li 
                v-for="page in visiblePages" 
                :key="page"
                class="page-item" 
                :class="{ active: page === currentPage }"
              >
                <button @click="goToPage(page)" class="page-link">{{ page }}</button>
              </li>
              <li class="page-item" :class="{ disabled: currentPage === totalPages }">
                <button @click="goToPage(currentPage + 1)" class="page-link">
                  {{ $t('image.next') }}
                </button>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>

    <!-- Upload Modal -->
    <ImageUploadModal
      v-if="showUploadModal" 
      @close="showUploadModal = false"
      @uploaded="onImageUploaded"
    />

    <!-- Image View Modal -->
    <ImageViewModal 
      v-if="selectedImage" 
      :image="selectedImage"
      @close="selectedImage = null"
      @updated="onImageUpdated"
      @deleted="onImageDeleted"
    />

    <!-- Image Edit Modal -->
    <ImageEditModal 
      v-if="editingImage" 
      :image="editingImage"
      @close="editingImage = null"
      @updated="onImageUpdated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useImageStore, type Image } from '@/stores/imageStore'
import ImageUploadModal from '@/components/image/ImageUploadModal.vue'
import ImageViewModal from '@/components/image/ImageViewModal.vue'
import ImageEditModal from '@/components/image/ImageEditModal.vue'

const { t } = useI18n()
const imageStore = useImageStore()

// Reactive data
const showUploadModal = ref(false)
const selectedImage = ref<any>(null)
const editingImage = ref<any>(null)

// Filters
const searchQuery = ref('')
const selectedType = ref('')
const selectedSource = ref('')
const selectedStatus = ref('')
const limit = ref(24)

// Computed from store  
const images = computed(() => imageStore.images || [])
const stats = computed(() => imageStore.stats)
const loading = computed(() => imageStore.loading)
const currentPage = computed(() => imageStore.currentPage)
const totalPages = computed(() => imageStore.totalPages)
const total = computed(() => imageStore.total)
const sortBy = computed(() => imageStore.sortBy)

// Computed
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
const loadImages = async () => {
  if (searchQuery.value.trim()) {
    await imageStore.searchImages(searchQuery.value, currentPage.value, limit.value)
  } else {
    await imageStore.listImages({
      type: selectedType.value || undefined,
      status: selectedStatus.value || undefined
    })
  }
}

const loadStats = async () => {
  await imageStore.getImageStats()
}

const onSearch = () => {
  imageStore.setSearchQuery(searchQuery.value)
  imageStore.goToPage(1)
  loadImages()
}

const onFilterChange = () => {
  imageStore.setFilters({
    type: selectedType.value || undefined,
    status: selectedStatus.value || undefined,
    source: selectedSource.value || undefined
  })
  imageStore.goToPage(1)
  loadImages()
  // NOTE: loadStats() called only on mount, not on every filter change
}

const goToPage = (page: number) => {
  imageStore.goToPage(page)
  loadImages()
}

const selectImage = (image: any) => {
  selectedImage.value = image
}

const viewImage = (image: any) => {
  selectedImage.value = image
}

const editImage = (image: any) => {
  editingImage.value = image
}

const downloadImage = async (image: any) => {
  await imageStore.downloadImage(image.id)
}

const deleteImage = async (image: any) => {
  if (!confirm(t('image.confirmDelete', { name: image.original_name }))) {
    return
  }
  
  await imageStore.deleteImage(image.id)
  await loadImages()
  await loadStats()
}

const onImageUploaded = () => {
  showUploadModal.value = false
  loadImages()
  loadStats()
}

const onImageUpdated = () => {
  selectedImage.value = null
  editingImage.value = null
  loadImages()
}

const onImageDeleted = () => {
  selectedImage.value = null
  loadImages()
  loadStats()
}

const onImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgdmlld0JveD0iMCAwIDIwMCAyMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIiBmaWxsPSIjRjVGNUY1Ii8+CjxwYXRoIGQ9Ik04MCA5MEgxMjBWNTEwSDEwMEwxMDAgOTBWOTBIODBWOTBaIiBmaWxsPSIjQ0NDQ0NDIi8+Cjwvc3ZnPg==' // Fallback SVG image
}

const getSourceBadgeClass = (source: string) => {
  return imageStore.getSourceBadgeClass(source)
}

const formatFileSize = (bytes: number): string => {
  return imageStore.formatFileSize(bytes)
}

// Watchers
watch(() => imageStore.currentPage, () => {
  loadImages()
})

// Lifecycle
onMounted(() => {
  loadImages()
  loadStats()
})
</script>

<style scoped>
.image-view {
  min-height: 100vh;
  background-color: #f8f9fa;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.image-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.image-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.image-container {
  position: relative;
  aspect-ratio: 16/9;
  overflow: hidden;
}

.image-thumbnail {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-card:hover .image-overlay {
  opacity: 1;
}

.image-actions {
  display: flex;
  gap: 0.5rem;
}

.image-info {
  padding: 1rem;
}

.image-title {
  font-size: 0.9rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  flex-wrap: wrap;
}

.image-stats {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.badge {
  font-size: 0.7rem;
}

@media (max-width: 768px) {
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .image-info {
    padding: 0.75rem;
  }
  
  .image-title {
    font-size: 0.8rem;
  }
}
</style>
