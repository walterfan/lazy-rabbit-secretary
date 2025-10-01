<template>
  <div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-xl">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            <i class="bi bi-diagram-3"></i>
            {{ diagram.name }}
          </h5>
          <div class="d-flex gap-2">
            <button 
              class="btn btn-outline-primary btn-sm"
              @click="$emit('edit', diagram)"
            >
              <i class="bi bi-pencil"></i>
              {{ $t('diagram.management.edit') }}
            </button>
            <button 
              type="button" 
              class="btn-close" 
              @click="$emit('close')"
            ></button>
          </div>
        </div>
        
        <div class="modal-body">
          <div class="row">
            <!-- Diagram Preview -->
            <div class="col-md-8">
              <div class="card">
                <div class="card-header d-flex justify-content-between align-items-center">
                  <h6 class="mb-0">{{ $t('diagram.management.preview') }}</h6>
                  <div class="btn-group btn-group-sm">
                    <button 
                      class="btn btn-outline-secondary"
                      @click="generatePreview"
                      :disabled="generating"
                    >
                      <span v-if="generating" class="spinner-border spinner-border-sm me-1"></span>
                      <i v-else class="bi bi-arrow-clockwise me-1"></i>
                      {{ generating ? $t('diagram.management.generating') : $t('diagram.management.regenerate') }}
                    </button>
                    <button 
                      v-if="previewImage"
                      class="btn btn-outline-primary"
                      @click="downloadPreview"
                    >
                      <i class="bi bi-download me-1"></i>
                      {{ $t('diagram.management.download') }}
                    </button>
                  </div>
                </div>
                <div class="card-body text-center">
                  <div v-if="generating" class="py-5">
                    <div class="spinner-border text-primary" role="status">
                      <span class="visually-hidden">{{ $t('diagram.management.generating') }}</span>
                    </div>
                    <p class="mt-3">{{ $t('diagram.management.generating') }}</p>
                  </div>
                  <div v-else-if="previewImage" class="diagram-preview">
                    <img 
                      :src="previewImage.image_data" 
                      :alt="diagram.name"
                      class="img-fluid"
                      style="max-height: 500px;"
                    />
                  </div>
                  <div v-else class="py-5 text-muted">
                    <i class="bi bi-image fs-1"></i>
                    <p class="mt-3">{{ $t('diagram.management.noPreview') }}</p>
                    <button 
                      class="btn btn-primary"
                      @click="generatePreview"
                    >
                      <i class="bi bi-play-circle me-1"></i>
                      {{ $t('diagram.management.generatePreview') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Diagram Details -->
            <div class="col-md-4">
              <div class="card">
                <div class="card-header">
                  <h6 class="mb-0">{{ $t('diagram.management.details') }}</h6>
                </div>
                <div class="card-body">
                  <!-- Basic Info -->
                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.name') }}</label>
                    <p class="mb-0">{{ diagram.name }}</p>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.type') }}</label>
                    <p class="mb-0">
                      <span class="badge bg-secondary">{{ diagram.type }}</span>
                    </p>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.scriptType') }}</label>
                    <p class="mb-0">
                      <span class="badge bg-info">{{ diagram.script_type }}</span>
                    </p>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.status') }}</label>
                    <p class="mb-0">
                      <span class="badge" :class="getStatusBadgeClass(diagram.status)">
                        {{ $t(`diagram.management.${diagram.status}`) }}
                      </span>
                    </p>
                  </div>

                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.theme') }}</label>
                    <p class="mb-0">{{ diagram.theme }}</p>
                  </div>

                  <!-- Description -->
                  <div class="mb-3" v-if="diagram.description">
                    <label class="form-label small text-muted">{{ $t('diagram.management.description') }}</label>
                    <p class="mb-0">{{ diagram.description }}</p>
                  </div>

                  <!-- Tags -->
                  <div class="mb-3" v-if="diagram.tags && diagram.tags.length > 0">
                    <label class="form-label small text-muted">{{ $t('diagram.management.tags') }}</label>
                    <div>
                      <span 
                        v-for="tag in diagram.tags" 
                        :key="tag"
                        class="badge bg-light text-dark me-1 mb-1"
                      >
                        {{ tag }}
                      </span>
                    </div>
                  </div>

                  <!-- Privacy -->
                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.privacy') }}</label>
                    <div class="d-flex gap-2">
                      <span 
                        v-if="diagram.public"
                        class="badge bg-success"
                      >
                        <i class="bi bi-globe me-1"></i>
                        {{ $t('diagram.management.public') }}
                      </span>
                      <span 
                        v-if="diagram.shared"
                        class="badge bg-info"
                      >
                        <i class="bi bi-share-fill me-1"></i>
                        {{ $t('diagram.management.shared') }}
                      </span>
                      <span 
                        v-if="!diagram.public && !diagram.shared"
                        class="badge bg-secondary"
                      >
                        <i class="bi bi-lock me-1"></i>
                        {{ $t('diagram.management.private') }}
                      </span>
                    </div>
                  </div>

                  <!-- Statistics -->
                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.statistics') }}</label>
                    <div class="row g-2">
                      <div class="col-6">
                        <div class="text-center">
                          <div class="h5 mb-0">{{ diagram.view_count }}</div>
                          <small class="text-muted">{{ $t('diagram.management.views') }}</small>
                        </div>
                      </div>
                      <div class="col-6">
                        <div class="text-center">
                          <div class="h5 mb-0">{{ diagram.edit_count }}</div>
                          <small class="text-muted">{{ $t('diagram.management.edits') }}</small>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Timestamps -->
                  <div class="mb-3">
                    <label class="form-label small text-muted">{{ $t('diagram.management.created') }}</label>
                    <p class="mb-1">{{ formatDate(diagram.created_at) }}</p>
                    <label class="form-label small text-muted">{{ $t('diagram.management.updated') }}</label>
                    <p class="mb-0">{{ formatDate(diagram.updated_at) }}</p>
                  </div>
                </div>
              </div>

              <!-- Script -->
              <div class="card mt-3">
                <div class="card-header">
                  <h6 class="mb-0">{{ $t('diagram.management.script') }}</h6>
                </div>
                <div class="card-body">
                  <pre class="bg-light p-3 rounded small" style="max-height: 200px; overflow-y: auto;"><code>{{ diagram.script }}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDiagramStore, type Diagram, type DrawDiagramResponse } from '@/stores/diagramStore'
import { formatDate } from '@/utils/dateUtils'

interface Props {
  diagram: Diagram
}

const props = defineProps<Props>()
const { t } = useI18n()
const diagramStore = useDiagramStore()

// State
const generating = ref(false)
const previewImage = ref<DrawDiagramResponse | null>(null)

// Methods
const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-secondary',
    published: 'bg-success',
    private: 'bg-warning',
    archived: 'bg-danger'
  }
  return classes[status] || 'bg-secondary'
}

const generatePreview = async () => {
  generating.value = true
  try {
    const result = await diagramStore.drawDiagram({
      script: props.diagram.script,
      script_type: props.diagram.script_type,
      format: 'png'
    })
    if (result) {
      previewImage.value = result
    }
  } catch (error) {
    console.error('Failed to generate preview:', error)
  } finally {
    generating.value = false
  }
}

const downloadPreview = () => {
  if (!previewImage.value) return
  
  const link = document.createElement('a')
  link.href = previewImage.value.image_data
  link.download = `${props.diagram.name}-preview.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// Lifecycle
onMounted(() => {
  // Auto-generate preview on mount
  generatePreview()
})

const emit = defineEmits<{
  close: []
  edit: [diagram: Diagram]
}>()
</script>

<style scoped>
.modal-xl {
  max-width: 1200px;
}

.diagram-preview {
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.form-label.small {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.badge {
  font-size: 0.75rem;
}

pre {
  font-size: 0.8rem;
  line-height: 1.4;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}
</style>
