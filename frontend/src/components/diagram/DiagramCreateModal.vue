<template>
  <div class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            <i class="bi bi-plus-circle"></i>
            {{ $t('diagram.management.create') }}
          </h5>
          <button 
            type="button" 
            class="btn-close" 
            @click="$emit('close')"
          ></button>
        </div>
        
        <form @submit.prevent="onSubmit">
          <div class="modal-body">
            <!-- Basic Information -->
            <div class="row g-3">
              <div class="col-md-6">
                <label class="form-label required">{{ $t('diagram.management.name') }}</label>
                <input
                  v-model="form.name"
                  type="text"
                  class="form-control"
                  :class="{ 'is-invalid': errors.name }"
                  :placeholder="$t('diagram.management.namePlaceholder')"
                  required
                />
                <div v-if="errors.name" class="invalid-feedback">
                  {{ errors.name }}
                </div>
              </div>
              
              <div class="col-md-6">
                <label class="form-label required">{{ $t('diagram.management.type') }}</label>
                <select
                  v-model="form.type"
                  class="form-select"
                  :class="{ 'is-invalid': errors.type }"
                  required
                >
                  <option value="">{{ $t('diagram.management.selectType') }}</option>
                  <option value="flowchart">{{ $t('diagram.management.flowchart') }}</option>
                  <option value="sequence">{{ $t('diagram.management.sequence') }}</option>
                  <option value="class">{{ $t('diagram.management.class') }}</option>
                  <option value="mindmap">{{ $t('diagram.management.mindmap') }}</option>
                  <option value="architecture">{{ $t('diagram.management.architecture') }}</option>
                  <option value="custom">{{ $t('diagram.management.custom') }}</option>
                </select>
                <div v-if="errors.type" class="invalid-feedback">
                  {{ errors.type }}
                </div>
              </div>
            </div>

            <div class="row g-3 mt-2">
              <div class="col-md-6">
                <label class="form-label required">{{ $t('diagram.management.scriptType') }}</label>
                <select
                  v-model="form.script_type"
                  class="form-select"
                  :class="{ 'is-invalid': errors.script_type }"
                  required
                >
                  <option value="">{{ $t('diagram.management.selectScriptType') }}</option>
                  <option value="plantuml">PlantUML</option>
                  <option value="mermaid">Mermaid</option>
                  <option value="graphviz">Graphviz</option>
                </select>
                <div v-if="errors.script_type" class="invalid-feedback">
                  {{ errors.script_type }}
                </div>
              </div>
              
              <div class="col-md-6">
                <label class="form-label">{{ $t('diagram.management.theme') }}</label>
                <select v-model="form.theme" class="form-select">
                  <option value="default">{{ $t('diagram.management.default') }}</option>
                  <option value="dark">{{ $t('diagram.management.dark') }}</option>
                  <option value="light">{{ $t('diagram.management.light') }}</option>
                  <option value="colorful">{{ $t('diagram.management.colorful') }}</option>
                </select>
              </div>
            </div>

            <!-- Description -->
            <div class="mt-3">
              <label class="form-label">{{ $t('diagram.management.description') }}</label>
              <textarea
                v-model="form.description"
                class="form-control"
                rows="3"
                :placeholder="$t('diagram.management.descriptionPlaceholder')"
              ></textarea>
            </div>

            <!-- Script -->
            <div class="mt-3">
              <label class="form-label required">{{ $t('diagram.management.script') }}</label>
              <textarea
                v-model="form.script"
                class="form-control"
                :class="{ 'is-invalid': errors.script }"
                rows="8"
                :placeholder="getScriptPlaceholder()"
                required
              ></textarea>
              <div v-if="errors.script" class="invalid-feedback">
                {{ errors.script }}
              </div>
              <div class="form-text">
                {{ $t('diagram.management.scriptHelp') }}
              </div>
            </div>

            <!-- Tags -->
            <div class="mt-3">
              <label class="form-label">{{ $t('diagram.management.tags') }}</label>
              <input
                v-model="tagInput"
                type="text"
                class="form-control"
                :placeholder="$t('diagram.management.tagsPlaceholder')"
                @keyup.enter="addTag"
              />
              <div class="form-text">
                {{ $t('diagram.management.tagsHelp') }}
              </div>
              <div v-if="form.tags && form.tags.length > 0" class="mt-2">
                <span 
                  v-for="(tag, index) in form.tags" 
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
            <div class="mt-3">
              <label class="form-label">{{ $t('diagram.management.privacy') }}</label>
              <div class="form-check">
                <input
                  v-model="form.public"
                  class="form-check-input"
                  type="checkbox"
                  id="public"
                />
                <label class="form-check-label" for="public">
                  {{ $t('diagram.management.makePublic') }}
                </label>
              </div>
              <div class="form-check">
                <input
                  v-model="form.shared"
                  class="form-check-input"
                  type="checkbox"
                  id="shared"
                />
                <label class="form-check-label" for="shared">
                  {{ $t('diagram.management.makeShared') }}
                </label>
              </div>
            </div>
          </div>

          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary" 
              @click="$emit('close')"
            >
              {{ $t('diagram.management.cancel') }}
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="saving"
            >
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              <i v-else class="bi bi-check-circle me-2"></i>
              {{ saving ? $t('diagram.management.creating') : $t('diagram.management.create') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDiagramStore } from '@/stores/diagramStore'

const { t } = useI18n()
const diagramStore = useDiagramStore()

// State
const saving = ref(false)
const tagInput = ref('')

// Form data
const form = reactive({
  name: '',
  type: '',
  script_type: '',
  script: '',
  description: '',
  theme: 'default',
  tags: [] as string[],
  public: false,
  shared: false
})

// Validation errors
const errors = reactive({
  name: '',
  type: '',
  script_type: '',
  script: ''
})

// Methods
const getScriptPlaceholder = () => {
  const placeholders = {
    plantuml: `@startuml
Alice -> Bob: Authentication Request
Bob --> Alice: Authentication Response
@enduml`,
    mermaid: `graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;`,
    graphviz: `digraph G {
    A -> B;
    B -> C;
    C -> A;
}`
  }
  return placeholders[form.script_type as keyof typeof placeholders] || t('diagram.management.scriptPlaceholder')
}

const addTag = () => {
  const tag = tagInput.value.trim()
  if (tag && !form.tags.includes(tag)) {
    form.tags.push(tag)
    tagInput.value = ''
  }
}

const removeTag = (index: number) => {
  form.tags.splice(index, 1)
}

const validateForm = () => {
  // Clear previous errors
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  if (!form.name.trim()) {
    errors.name = t('diagram.management.errors.nameRequired')
    isValid = false
  }

  if (!form.type) {
    errors.type = t('diagram.management.errors.typeRequired')
    isValid = false
  }

  if (!form.script_type) {
    errors.script_type = t('diagram.management.errors.scriptTypeRequired')
    isValid = false
  }

  if (!form.script.trim()) {
    errors.script = t('diagram.management.errors.scriptRequired')
    isValid = false
  }

  return isValid
}

const onSubmit = async () => {
  if (!validateForm()) {
    return
  }

  saving.value = true
  try {
    await diagramStore.createDiagram({
      name: form.name,
      type: form.type,
      script_type: form.script_type,
      script: form.script,
      description: form.description,
      tags: form.tags,
      public: form.public,
      shared: form.shared
    })
    
    // Reset form
    Object.assign(form, {
      name: '',
      type: '',
      script_type: '',
      script: '',
      description: '',
      theme: 'default',
      tags: [],
      public: false,
      shared: false
    })
    
    emit('created')
  } catch (error) {
    console.error('Failed to create diagram:', error)
  } finally {
    saving.value = false
  }
}

// Watchers
watch(() => form.script_type, () => {
  if (form.script_type) {
    form.script = getScriptPlaceholder()
  }
})

const emit = defineEmits<{
  close: []
  created: []
}>()
</script>

<style scoped>
.required::after {
  content: ' *';
  color: #dc3545;
}

.modal-dialog {
  max-width: 800px;
}

.form-control.is-invalid,
.form-select.is-invalid {
  border-color: #dc3545;
}

.invalid-feedback {
  display: block;
}

.badge .btn-close {
  font-size: 0.6em;
  padding: 0;
  margin: 0;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
}
</style>
