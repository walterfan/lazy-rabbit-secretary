<template>
  <div class="modal fade show d-block" tabindex="-1" @click.self="$emit('close')">
    <div class="modal-dialog modal-xl">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            <i class="bi bi-image"></i>
            {{ image.original_name }}
          </h5>
          <div class="btn-group">
            <button @click="downloadImage" class="btn btn-outline-primary btn-sm">
              <i class="bi bi-download"></i>
              {{ $t('image.download') }}
            </button>
            <button @click="editImage" class="btn btn-outline-secondary btn-sm">
              <i class="bi bi-pencil"></i>
              {{ $t('image.edit') }}
            </button>
            <button @click="deleteImage" class="btn btn-outline-danger btn-sm">
              <i class="bi bi-trash"></i>
              {{ $t('image.delete') }}
            </button>
          </div>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        
        <div class="modal-body">
          <div class="row">
            <!-- Image Display -->
            <div class="col-md-8">
              <div class="image-display">
                <img 
                  :src="image.url" 
                  :alt="image.original_name"
                  class="img-fluid"
                  @error="onImageError"
                >
              </div>
            </div>
            
            <!-- Image Details -->
            <div class="col-md-4">
              <div class="image-details">
                <h6>{{ $t('image.details') }}</h6>
                
                <div class="detail-item">
                  <label>{{ $t('image.originalName') }}:</label>
                  <span>{{ image.original_name }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.type') }}:</label>
                  <span class="badge bg-secondary">{{ image.type }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.source') }}:</label>
                  <span class="badge" :class="getSourceBadgeClass(image.source)">
                    {{ image.source }}
                  </span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.status') }}:</label>
                  <span class="badge" :class="getStatusBadgeClass(image.status)">
                    {{ image.status }}
                  </span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.dimensions') }}:</label>
                  <span>{{ image.width }} × {{ image.height }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.fileSize') }}:</label>
                  <span>{{ formatFileSize(image.file_size) }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.format') }}:</label>
                  <span>{{ image.format.toUpperCase() }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.mimeType') }}:</label>
                  <span>{{ image.mime_type }}</span>
                </div>
                
                <div class="detail-item" v-if="image.category">
                  <label>{{ $t('image.category') }}:</label>
                  <span>{{ image.category }}</span>
                </div>
                
                <div class="detail-item" v-if="image.description">
                  <label>{{ $t('image.description') }}:</label>
                  <span>{{ image.description }}</span>
                </div>
                
                <div class="detail-item" v-if="image.tags.length > 0">
                  <label>{{ $t('image.tags') }}:</label>
                  <div class="tags">
                    <span 
                      v-for="tag in image.tags" 
                      :key="tag"
                      class="badge bg-primary me-1 mb-1"
                    >
                      {{ tag }}
                    </span>
                  </div>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.privacy') }}:</label>
                  <div class="privacy-status">
                    <span v-if="image.is_public" class="badge bg-success me-1">
                      <i class="bi bi-globe"></i> {{ $t('image.public') }}
                    </span>
                    <span v-if="image.is_shared" class="badge bg-info me-1">
                      <i class="bi bi-share"></i> {{ $t('image.shared') }}
                    </span>
                    <span v-if="!image.is_public && !image.is_shared" class="badge bg-secondary">
                      <i class="bi bi-lock"></i> {{ $t('image.private') }}
                    </span>
                  </div>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.stats') }}:</label>
                  <div class="stats">
                    <small class="text-muted">
                      <i class="bi bi-eye"></i> {{ image.view_count }} {{ $t('image.views') }}
                    </small>
                    <br>
                    <small class="text-muted">
                      <i class="bi bi-download"></i> {{ image.download_count }} {{ $t('image.downloads') }}
                    </small>
                  </div>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.created') }}:</label>
                  <span>{{ formatDate(image.created_at) }}</span>
                </div>
                
                <div class="detail-item">
                  <label>{{ $t('image.updated') }}:</label>
                  <span>{{ formatDate(image.updated_at) }}</span>
                </div>
                
                <div class="detail-item" v-if="image.diagram_id">
                  <label>{{ $t('image.diagram') }}:</label>
                  <span class="badge bg-warning">{{ image.diagram_id }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            {{ $t('common.close') }}
          </button>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Backdrop -->
  <div class="modal-backdrop fade show"></div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useImageStore, type Image } from '@/stores/imageStore'
import { formatDate } from '@/utils/dateUtils'

const { t } = useI18n()
const imageStore = useImageStore()

// Props
interface Props {
  image: Image
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  close: []
  updated: []
  deleted: []
}>()

// Methods
const downloadImage = async () => {
  try {
    await imageStore.downloadImage(props.image.id)
  } catch (error) {
    console.error('Failed to download image:', error)
  }
}

const editImage = () => {
  emit('updated')
}

const deleteImage = async () => {
  if (!confirm(t('image.confirmDelete', { name: props.image.original_name }))) {
    return
  }
  
  try {
    await imageStore.deleteImage(props.image.id)
    emit('deleted')
  } catch (error) {
    console.error('Failed to delete image:', error)
  }
}

const onImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgdmlld0JveD0iMCAwIDIwMCAyMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIiBmaWxsPSIjRjVGNUY1Ii8+CjxwYXRoIGQ9Ik04MCA5MEgxMjBWNTEwSDEwMEwxMDAgOTBWOTBIODBWOTBaIiBmaWxsPSIjQ0NDQ0NDIi8+Cjwvc3ZnPg==' // Fallback SVG image
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

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.modal {
  z-index: 1055;
}

.modal-backdrop {
  z-index: 1050;
}

.image-display {
  text-align: center;
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 1rem;
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-display img {
  max-width: 100%;
  max-height: 500px;
  object-fit: contain;
  border-radius: 4px;
}

.image-details {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
}

.detail-item {
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #dee2e6;
}

.detail-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.detail-item label {
  font-weight: 600;
  color: #495057;
  display: block;
  margin-bottom: 0.25rem;
}

.detail-item span {
  color: #6c757d;
}

.tags {
  margin-top: 0.25rem;
}

.privacy-status {
  margin-top: 0.25rem;
}

.stats {
  margin-top: 0.25rem;
}

.badge {
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .image-display {
    min-height: 300px;
  }
  
  .image-details {
    padding: 1rem;
  }
}
</style>
