<template>
  <div class="diagram-view">
    <div class="container-fluid">
      <div class="row">
        <!-- Sidebar -->
        <div class="col-md-3">
          <div class="sidebar">
            <h5 class="mb-3">
              <i class="bi bi-diagram-3"></i>
              {{ $t('diagram.title') }}
            </h5>
            
            <!-- Diagram Type Selection -->
            <div class="mb-3">
              <label class="form-label">{{ $t('diagram.scriptType') }}</label>
              <select v-model="scriptType" class="form-select" @change="onScriptTypeChange">
                <option value="plantuml">PlantUML</option>
                <option value="mermaid">Mermaid</option>
                <option value="graphviz">Graphviz</option>
              </select>
            </div>

            <!-- Diagram Settings -->
            <div class="mb-3">
              <label class="form-label">{{ $t('diagram.width') }}</label>
              <input 
                v-model.number="width" 
                type="number" 
                class="form-control" 
                min="100" 
                max="2000"
                :placeholder="$t('diagram.widthPlaceholder')"
              >
            </div>

            <div class="mb-3">
              <label class="form-label">{{ $t('diagram.height') }}</label>
              <input 
                v-model.number="height" 
                type="number" 
                class="form-control" 
                min="100" 
                max="2000"
                :placeholder="$t('diagram.heightPlaceholder')"
              >
            </div>

            <!-- Action Buttons -->
            <div class="d-grid gap-2">
              <button 
                @click="generateDiagram"
                class="btn btn-primary"
                :disabled="!script.trim() || generating"
              >
                <i class="bi bi-play-circle" v-if="!generating"></i>
                <span class="spinner-border spinner-border-sm me-2" v-if="generating"></span>
                {{ generating ? $t('diagram.generating') : $t('diagram.generate') }}
              </button>
              
              <button 
                @click="saveDiagram" 
                class="btn btn-success"
                :disabled="!generatedImage || saving"
              >
                <i class="bi bi-save" v-if="!saving"></i>
                <span class="spinner-border spinner-border-sm me-2" v-if="saving"></span>
                {{ saving ? $t('diagram.saving') : $t('diagram.save') }}
              </button>
              
              <button @click="clearAll" class="btn btn-outline-secondary">
                <i class="bi bi-trash"></i>
                {{ $t('diagram.clear') }}
              </button>
            </div>

            <!-- Recent Diagrams -->
            <div class="mt-4" v-if="recentDiagrams && recentDiagrams.length > 0">
              <h6>{{ $t('diagram.recent') }}</h6>
              <div class="list-group">
                <button 
                  v-for="diagram in (recentDiagrams || [])" 
                  :key="diagram.id"
                  @click="loadDiagram(diagram)"
                  class="list-group-item list-group-item-action"
                >
                  <div class="d-flex w-100 justify-content-between">
                    <h6 class="mb-1">{{ diagram.name }}</h6>
                    <small>{{ formatDate(diagram.updated_at) }}</small>
                  </div>
                  <p class="mb-1 text-muted">{{ diagram.description || $t('diagram.noDescription') }}</p>
                  <small>{{ diagram.script_type }} • {{ diagram.type }}</small>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Main Content -->
        <div class="col-md-9">
          <div class="main-content">
            <!-- Script Editor -->
            <div class="editor-section mb-4">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <h5>{{ $t('diagram.scriptEditor') }}</h5>
                <div class="btn-group" role="group">
                  <button 
                    @click="showPreview = !showPreview" 
                    class="btn btn-outline-primary btn-sm"
                  >
                    <i class="bi bi-eye" v-if="!showPreview"></i>
                    <i class="bi bi-eye-slash" v-if="showPreview"></i>
                    {{ showPreview ? $t('diagram.hidePreview') : $t('diagram.showPreview') }}
                  </button>
                </div>
              </div>
              
              <div class="row" v-if="showPreview">
                <div class="col-md-6">
                  <label class="form-label">{{ $t('diagram.script') }}</label>
                  <textarea 
                    v-model="script" 
                    class="form-control script-editor" 
                    rows="15"
                    :placeholder="getScriptPlaceholder()"
                  ></textarea>
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ $t('diagram.preview') }}</label>
                  <div class="preview-container">
                    <div v-if="generatedImage" class="generated-image">
                      <img :src="generatedImage.image_data" :alt="$t('diagram.generatedImage')" class="img-fluid">
                    </div>
                    <div v-else class="preview-placeholder">
                      <i class="bi bi-image"></i>
                      <p>{{ $t('diagram.noPreview') }}</p>
                    </div>
                  </div>
                </div>
              </div>
              
              <div v-else>
                <label class="form-label">{{ $t('diagram.script') }}</label>
                <textarea 
                  v-model="script" 
                  class="form-control script-editor" 
                  rows="15"
                  :placeholder="getScriptPlaceholder()"
                ></textarea>
              </div>
            </div>

            <!-- Generated Image Display -->
            <div class="result-section" v-if="generatedImage">
              <div class="d-flex justify-content-between align-items-center mb-3">
                <h5>{{ $t('diagram.result') }}</h5>
                <div class="btn-group" role="group">
                  <button @click="downloadImage" class="btn btn-outline-primary btn-sm">
                    <i class="bi bi-download"></i>
                    {{ $t('diagram.download') }}
                  </button>
                  <button @click="copyImageUrl" class="btn btn-outline-secondary btn-sm">
                    <i class="bi bi-clipboard"></i>
                    {{ $t('diagram.copyUrl') }}
                  </button>
                </div>
              </div>
              
              <div class="generated-result">
                <img :src="generatedImage.image_data" :alt="$t('diagram.generatedImage')" class="img-fluid">
                <div class="image-info mt-2">
                  <small class="text-muted">
                    {{ $t('diagram.format') }}: {{ generatedImage.format }} • 
                    {{ $t('diagram.size') }}: {{ generatedImage.width }}x{{ generatedImage.height }} • 
                    {{ $t('diagram.fileSize') }}: {{ formatFileSize(generatedImage.size) }}
                  </small>
                </div>
              </div>
            </div>

            <!-- Error Display -->
            <div v-if="error" class="alert alert-danger" role="alert">
              <i class="bi bi-exclamation-triangle"></i>
              {{ error }}
            </div>

            <!-- Success Message -->
            <div v-if="successMessage" class="alert alert-success" role="alert">
              <i class="bi bi-check-circle"></i>
              {{ successMessage }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDiagramStore, type Diagram } from '@/stores/diagramStore'
import { formatDate } from '@/utils/dateUtils'

const { t } = useI18n()
const diagramStore = useDiagramStore()

// Reactive data
const script = ref('')
const scriptType = ref('plantuml')
const width = ref(800)
const height = ref(600)
const showPreview = ref(true)

// Computed from store
const generatedImage = computed(() => diagramStore.generatedImage)
const generating = computed(() => diagramStore.generating)
const saving = computed(() => diagramStore.saving)
const error = computed(() => diagramStore.error)
const successMessage = computed(() => diagramStore.successMessage)
const recentDiagrams = computed(() => diagramStore.recentDiagrams)

// Methods
const getScriptPlaceholder = () => {
  const placeholders = {
    plantuml: `@startuml
Alice -> Bob: Hello
Bob -> Alice: Hi there!
@enduml`,
    mermaid: `graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
    C --> E[End]
    D --> E`,
    graphviz: `digraph G {
    A -> B
    B -> C
    C -> A
}`
  }
  return placeholders[scriptType.value as keyof typeof placeholders] || ''
}

const onScriptTypeChange = () => {
  script.value = getScriptPlaceholder()
  diagramStore.clearGeneratedImage()
  diagramStore.clearMessages()
}

const generateDiagram = async () => {
  if (!script.value.trim()) {
    diagramStore.error = t('diagram.errors.noScript')
    return
  }

  await diagramStore.drawDiagram({
    script: script.value,
    script_type: scriptType.value,
    width: width.value,
    height: height.value,
    format: 'png'
  })
}

const saveDiagram = async () => {
  if (!generatedImage.value) {
    diagramStore.error = t('diagram.errors.noImageToSave')
    return
  }

  const diagramName = prompt(t('diagram.prompt.diagramName'))
  if (!diagramName) {
    return
  }

  await diagramStore.createDiagram({
    name: diagramName,
    type: 'custom',
    script_type: scriptType.value,
    script: script.value,
    description: t('diagram.generatedDescription', { type: scriptType.value }),
    tags: [scriptType.value, 'generated']
  })
  
  await diagramStore.loadRecentDiagrams()
}

const downloadImage = () => {
  if (!generatedImage.value) return
  
  const link = document.createElement('a')
  link.href = generatedImage.value.image_data
  link.download = `diagram-${Date.now()}.${generatedImage.value.format}`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const copyImageUrl = async () => {
  if (!generatedImage.value) return
  
  try {
    await navigator.clipboard.writeText(generatedImage.value.url)
    diagramStore.successMessage = t('diagram.success.urlCopied')
  } catch (err) {
    diagramStore.error = t('diagram.errors.copyFailed')
  }
}

const clearAll = () => {
  script.value = ''
  diagramStore.clearGeneratedImage()
  diagramStore.clearMessages()
}

const loadDiagram = (diagram: any) => {
  script.value = diagram.script
  scriptType.value = diagram.script_type
  diagramStore.clearGeneratedImage()
  diagramStore.clearMessages()
}

const loadRecentDiagrams = async () => {
  await diagramStore.loadRecentDiagrams()
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Lifecycle
onMounted(() => {
  script.value = getScriptPlaceholder()
  loadRecentDiagrams()
})
</script>

<style scoped>
.diagram-view {
  min-height: 100vh;
  background-color: #f8f9fa;
}

.sidebar {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  height: fit-content;
  position: sticky;
  top: 1rem;
}

.main-content {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.script-editor {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.preview-container {
  border: 1px solid #dee2e6;
  border-radius: 4px;
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
}

.preview-placeholder {
  text-align: center;
  color: #6c757d;
}

.preview-placeholder i {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.generated-image img {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.generated-result {
  text-align: center;
  padding: 1rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  background-color: #f8f9fa;
}

.image-info {
  text-align: center;
}

.list-group-item {
  border: none;
  border-bottom: 1px solid #dee2e6;
}

.list-group-item:last-child {
  border-bottom: none;
}

.btn-group .btn {
  margin-right: 0.25rem;
}

.btn-group .btn:last-child {
  margin-right: 0;
}

@media (max-width: 768px) {
  .sidebar {
    position: static;
    margin-bottom: 1rem;
  }
  
  .main-content {
    margin-bottom: 1rem;
  }
}
</style>
