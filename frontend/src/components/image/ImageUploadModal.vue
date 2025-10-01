<template>
  <div class="modal fade show d-block" tabindex="-1" @click.self="$emit('close')">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            <i class="bi bi-cloud-upload"></i>
            {{ $t('image.upload') }}
          </h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        
        <div class="modal-body">
          <form @submit.prevent="uploadImage">
            <!-- File Upload -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.selectFile') }}</label>
              <input 
                ref="fileInput"
                type="file" 
                class="form-control" 
                accept="image/*"
                @change="onFileSelect"
                required
              >
              <div class="form-text">{{ $t('image.supportedFormats') }}</div>
            </div>

            <!-- File Preview -->
            <div v-if="selectedFile" class="mb-3">
              <label class="form-label">{{ $t('image.preview') }}</label>
              <div class="file-preview">
                <img :src="filePreview" :alt="selectedFile.name" class="preview-image">
                <div class="file-info">
                  <p><strong>{{ selectedFile.name }}</strong></p>
                  <p class="text-muted">{{ formatFileSize(selectedFile.size) }}</p>
                </div>
              </div>
            </div>

            <!-- Image Type -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.type') }} <span class="text-danger">*</span></label>
              <select v-model="formData.type" class="form-select" required>
                <option value="">{{ $t('image.selectType') }}</option>
                <option value="avatar">{{ $t('image.avatar') }}</option>
                <option value="profile">{{ $t('image.profile') }}</option>
                <option value="cover">{{ $t('image.cover') }}</option>
                <option value="gallery">{{ $t('image.gallery') }}</option>
                <option value="attachment">{{ $t('image.attachment') }}</option>
                <option value="custom">{{ $t('image.custom') }}</option>
              </select>
            </div>

            <!-- Category -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.category') }}</label>
              <input 
                v-model="formData.category" 
                type="text" 
                class="form-control" 
                :placeholder="$t('image.categoryPlaceholder')"
              >
            </div>

            <!-- Description -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.description') }}</label>
              <textarea 
                v-model="formData.description" 
                class="form-control" 
                rows="3"
                :placeholder="$t('image.descriptionPlaceholder')"
              ></textarea>
            </div>

            <!-- Tags -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.tags') }}</label>
              <input 
                v-model="tagsInput" 
                type="text" 
                class="form-control" 
                :placeholder="$t('image.tagsPlaceholder')"
                @keyup.enter="addTag"
              >
              <div class="form-text">{{ $t('image.tagsHelp') }}</div>
              <div v-if="formData.tags && formData.tags.length > 0" class="mt-2">
                <span 
                  v-for="(tag, index) in formData.tags" 
                  :key="index"
                  class="badge bg-primary me-1 mb-1"
                >
                  {{ tag }}
                  <button 
                    type="button" 
                    class="btn-close btn-close-white ms-1" 
                    @click="removeTag(index)"
                  ></button>
                </span>
              </div>
            </div>

            <!-- Privacy Settings -->
            <div class="mb-3">
              <label class="form-label">{{ $t('image.privacy') }}</label>
              <div class="form-check">
                <input 
                  v-model="formData.is_public" 
                  type="checkbox" 
                  class="form-check-input" 
                  id="isPublic"
                >
                <label class="form-check-label" for="isPublic">
                  {{ $t('image.makePublic') }}
                </label>
              </div>
              <div class="form-check">
                <input 
                  v-model="formData.is_shared" 
                  type="checkbox" 
                  class="form-check-input" 
                  id="isShared"
                >
                <label class="form-check-label" for="isShared">
                  {{ $t('image.makeShared') }}
                </label>
              </div>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="alert alert-danger" role="alert">
              <i class="bi bi-exclamation-triangle"></i>
              {{ error }}
            </div>
          </form>
        </div>
        
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            {{ $t('common.cancel') }}
          </button>
          <button 
            type="button" 
            class="btn btn-primary" 
            @click="uploadImage"
            :disabled="!selectedFile || uploading"
          >
            <i class="bi bi-cloud-upload" v-if="!uploading"></i>
            <span class="spinner-border spinner-border-sm me-2" v-if="uploading"></span>
            {{ uploading ? $t('image.uploading') : $t('image.upload') }}
          </button>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Backdrop -->
  <div class="modal-backdrop fade show"></div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useImageStore, type UploadImageRequest } from '@/stores/imageStore'

const { t } = useI18n()
const imageStore = useImageStore()

// Emits
const emit = defineEmits<{
  close: []
  uploaded: []
}>()

// Reactive data
const fileInput = ref<HTMLInputElement>()
const selectedFile = ref<File | null>(null)
const filePreview = ref<string>('')
const tagsInput = ref('')
const uploading = ref(false)
const error = ref('')


const formData = ref<UploadImageRequest>({
  type: '',
  category: '',
  description: '',
  tags: [],
  is_public: false,
  is_shared: false
})

// Methods
const onFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (file) {
    selectedFile.value = file
    
    // Create preview
    const reader = new FileReader()
    reader.onload = (e) => {
      filePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const addTag = () => {
  const tag = tagsInput.value.trim()
  if (tag && formData.value.tags && !formData.value.tags.includes(tag)) {
    formData.value.tags.push(tag)
    tagsInput.value = ''
  }
}

const removeTag = (index: number) => {
  if (formData.value.tags) {
    formData.value.tags.splice(index, 1)
  }
}

const uploadImage = async () => {
  if (!selectedFile.value) {
    error.value = t('image.errors.noFileSelected')
    return
  }

  if (!formData.value.type) {
    error.value = t('image.errors.noTypeSelected')
    return
  }

  uploading.value = true
  error.value = ''

  try {
    await imageStore.uploadImage(selectedFile.value, formData.value)
    emit('uploaded')
  } catch (err: any) {
    error.value = err.response?.data?.error || t('image.errors.uploadFailed')
  } finally {
    uploading.value = false
  }
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

.file-preview {
  display: flex;
  gap: 1rem;
  align-items: center;
  padding: 1rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  background-color: #f8f9fa;
}

.preview-image {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 4px;
}

.file-info {
  flex: 1;
}

.file-info p {
  margin: 0;
}

.badge {
  font-size: 0.8rem;
}

.btn-close {
  font-size: 0.7rem;
}
</style>
