<template>
  <div class="diagram-management-view">
    <div class="container-fluid">
      <!-- Header -->
      <div class="row mb-4">
        <div class="col-12">
          <div class="d-flex justify-content-between align-items-center">
            <h2>
              <i class="bi bi-diagram-3"></i>
              {{ $t('diagram.management.title') }}
            </h2>
            <button 
              class="btn btn-primary"
              @click="showCreateModal = true"
            >
              <i class="bi bi-plus-circle"></i>
              {{ $t('diagram.management.create') }}
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
                <!-- Search -->
                <div class="col-md-4">
                  <label class="form-label">{{ $t('diagram.management.search') }}</label>
                  <div class="input-group">
                    <input
                      v-model="searchQuery"
                      type="text"
                      class="form-control"
                      :placeholder="$t('diagram.management.searchPlaceholder')"
                      @keyup.enter="onSearch"
                    />
                    <button 
                      class="btn btn-outline-secondary"
                      @click="onSearch"
                    >
                      <i class="bi bi-search"></i>
                    </button>
                  </div>
                </div>

                <!-- Type Filter -->
                <div class="col-md-2">
                  <label class="form-label">{{ $t('diagram.management.type') }}</label>
                  <select v-model="selectedType" class="form-select" @change="onFilterChange">
                    <option value="">{{ $t('diagram.management.allTypes') }}</option>
                    <option value="flowchart">{{ $t('diagram.management.flowchart') }}</option>
                    <option value="sequence">{{ $t('diagram.management.sequence') }}</option>
                    <option value="class">{{ $t('diagram.management.class') }}</option>
                    <option value="mindmap">{{ $t('diagram.management.mindmap') }}</option>
                    <option value="architecture">{{ $t('diagram.management.architecture') }}</option>
                    <option value="custom">{{ $t('diagram.management.custom') }}</option>
                  </select>
                </div>

                <!-- Script Type Filter -->
                <div class="col-md-2">
                  <label class="form-label">{{ $t('diagram.management.scriptType') }}</label>
                  <select v-model="selectedScriptType" class="form-select" @change="onFilterChange">
                    <option value="">{{ $t('diagram.management.allScriptTypes') }}</option>
                    <option value="plantuml">PlantUML</option>
                    <option value="mermaid">Mermaid</option>
                    <option value="graphviz">Graphviz</option>
                  </select>
                </div>

                <!-- Status Filter -->
                <div class="col-md-2">
                  <label class="form-label">{{ $t('diagram.management.status') }}</label>
                  <select v-model="selectedStatus" class="form-select" @change="onFilterChange">
                    <option value="">{{ $t('diagram.management.allStatuses') }}</option>
                    <option value="draft">{{ $t('diagram.management.draft') }}</option>
                    <option value="published">{{ $t('diagram.management.published') }}</option>
                    <option value="private">{{ $t('diagram.management.private') }}</option>
                    <option value="archived">{{ $t('diagram.management.archived') }}</option>
                  </select>
                </div>

                <!-- Sort -->
                <div class="col-md-2">
                  <label class="form-label">{{ $t('diagram.management.sortBy') }}</label>
                  <select v-model="sortBy" class="form-select" @change="onFilterChange">
                    <option value="created_at">{{ $t('diagram.management.dateCreated') }}</option>
                    <option value="updated_at">{{ $t('diagram.management.dateUpdated') }}</option>
                    <option value="name">{{ $t('diagram.management.name') }}</option>
                    <option value="view_count">{{ $t('diagram.management.viewCount') }}</option>
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
          <div class="card bg-primary text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between">
                <div>
                  <h4>{{ stats.total_diagrams }}</h4>
                  <p class="mb-0">{{ $t('diagram.management.totalDiagrams') }}</p>
                </div>
                <div class="align-self-center">
                  <i class="bi bi-diagram-3 fs-1"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card bg-success text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between">
                <div>
                  <h4>{{ stats.published_diagrams }}</h4>
                  <p class="mb-0">{{ $t('diagram.management.publishedDiagrams') }}</p>
                </div>
                <div class="align-self-center">
                  <i class="bi bi-check-circle fs-1"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card bg-info text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between">
                <div>
                  <h4>{{ stats.public_diagrams }}</h4>
                  <p class="mb-0">{{ $t('diagram.management.publicDiagrams') }}</p>
                </div>
                <div class="align-self-center">
                  <i class="bi bi-globe fs-1"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card bg-warning text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between">
                <div>
                  <h4>{{ stats.total_views }}</h4>
                  <p class="mb-0">{{ $t('diagram.management.totalViews') }}</p>
                </div>
                <div class="align-self-center">
                  <i class="bi bi-eye fs-1"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Diagrams List -->
      <div class="row">
        <div class="col-12">
          <div class="card">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h5 class="mb-0">{{ $t('diagram.management.diagrams') }}</h5>
              <div class="d-flex gap-2">
                <button 
                  class="btn btn-outline-secondary btn-sm"
                  @click="loadDiagrams"
                  :disabled="loading"
                >
                  <i class="bi bi-arrow-clockwise" :class="{ 'spinning': loading }"></i>
                  {{ $t('diagram.management.refresh') }}
                </button>
                <div class="btn-group" role="group">
                  <button 
                    class="btn btn-outline-secondary btn-sm"
                    :class="{ active: viewMode === 'grid' }"
                    @click="viewMode = 'grid'"
                  >
                    <i class="bi bi-grid-3x3-gap"></i>
                  </button>
                  <button 
                    class="btn btn-outline-secondary btn-sm"
                    :class="{ active: viewMode === 'list' }"
                    @click="viewMode = 'list'"
                  >
                    <i class="bi bi-list"></i>
                  </button>
                </div>
              </div>
            </div>
            <div class="card-body">
              <!-- Loading -->
              <div v-if="loading" class="text-center py-4">
                <div class="spinner-border" role="status">
                  <span class="visually-hidden">{{ $t('diagram.management.loading') }}</span>
                </div>
                <p class="mt-2">{{ $t('diagram.management.loading') }}</p>
              </div>

              <!-- Error -->
              <div v-else-if="error" class="alert alert-danger">
                <i class="bi bi-exclamation-triangle"></i>
                {{ error }}
              </div>

              <!-- Empty State -->
              <div v-else-if="diagrams.length === 0" class="text-center py-5">
                <i class="bi bi-diagram-3 fs-1 text-muted"></i>
                <h5 class="mt-3">{{ $t('diagram.management.noDiagrams') }}</h5>
                <p class="text-muted">{{ $t('diagram.management.noDiagramsDescription') }}</p>
                <button 
                  class="btn btn-primary"
                  @click="showCreateModal = true"
                >
                  <i class="bi bi-plus-circle"></i>
                  {{ $t('diagram.management.createFirst') }}
                </button>
              </div>

              <!-- Grid View -->
              <div v-else-if="viewMode === 'grid'" class="row g-3">
                <div 
                  v-for="diagram in diagrams" 
                  :key="diagram.id"
                  class="col-md-6 col-lg-4"
                >
                  <DiagramCard 
                    :diagram="diagram"
                    @edit="editDiagram"
                    @delete="deleteDiagram"
                    @view="viewDiagram"
                    @publish="publishDiagram"
                    @archive="archiveDiagram"
                    @share="shareDiagram"
                    @unshare="unshareDiagram"
                  />
                </div>
              </div>

              <!-- List View -->
              <div v-else class="table-responsive">
                <table class="table table-hover">
                  <thead>
                    <tr>
                      <th>{{ $t('diagram.management.name') }}</th>
                      <th>{{ $t('diagram.management.type') }}</th>
                      <th>{{ $t('diagram.management.scriptType') }}</th>
                      <th>{{ $t('diagram.management.status') }}</th>
                      <th>{{ $t('diagram.management.views') }}</th>
                      <th>{{ $t('diagram.management.created') }}</th>
                      <th>{{ $t('diagram.management.actions') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="diagram in diagrams" :key="diagram.id">
                      <td>
                        <div class="d-flex align-items-center">
                          <i class="bi bi-diagram-3 me-2"></i>
                          <div>
                            <strong>{{ diagram.name }}</strong>
                            <br>
                            <small class="text-muted">{{ diagram.description || $t('diagram.management.noDescription') }}</small>
                          </div>
                        </div>
                      </td>
                      <td>
                        <span class="badge bg-secondary">{{ diagram.type }}</span>
                      </td>
                      <td>
                        <span class="badge bg-info">{{ diagram.script_type }}</span>
                      </td>
                      <td>
                        <span class="badge" :class="getStatusBadgeClass(diagram.status)">
                          {{ $t(`diagram.management.${diagram.status}`) }}
                        </span>
                      </td>
                      <td>{{ diagram.view_count }}</td>
                      <td>{{ formatDate(diagram.created_at) }}</td>
                      <td>
                        <div class="btn-group btn-group-sm">
                          <button 
                            class="btn btn-outline-primary"
                            @click="viewDiagram(diagram)"
                            :title="$t('diagram.management.view')"
                          >
                            <i class="bi bi-eye"></i>
                          </button>
                          <button 
                            class="btn btn-outline-secondary"
                            @click="editDiagram(diagram)"
                            :title="$t('diagram.management.edit')"
                          >
                            <i class="bi bi-pencil"></i>
                          </button>
                          <div class="dropdown">
                            <button 
                              class="btn btn-outline-secondary dropdown-toggle"
                              type="button"
                              data-bs-toggle="dropdown"
                            >
                              <i class="bi bi-three-dots"></i>
                            </button>
                            <ul class="dropdown-menu">
                              <li v-if="diagram.status === 'draft'">
                                <button 
                                  class="dropdown-item"
                                  @click="publishDiagram(diagram)"
                                >
                                  <i class="bi bi-check-circle"></i>
                                  {{ $t('diagram.management.publish') }}
                                </button>
                              </li>
                              <li v-if="diagram.status === 'published'">
                                <button 
                                  class="dropdown-item"
                                  @click="archiveDiagram(diagram)"
                                >
                                  <i class="bi bi-archive"></i>
                                  {{ $t('diagram.management.archive') }}
                                </button>
                              </li>
                              <li>
                                <button 
                                  class="dropdown-item"
                                  @click="shareDiagram(diagram)"
                                  v-if="!diagram.shared"
                                >
                                  <i class="bi bi-share"></i>
                                  {{ $t('diagram.management.share') }}
                                </button>
                                <button 
                                  class="dropdown-item"
                                  @click="unshareDiagram(diagram)"
                                  v-else
                                >
                                  <i class="bi bi-share-fill"></i>
                                  {{ $t('diagram.management.unshare') }}
                                </button>
                              </li>
                              <li><hr class="dropdown-divider"></li>
                              <li>
                                <button 
                                  class="dropdown-item text-danger"
                                  @click="deleteDiagram(diagram)"
                                >
                                  <i class="bi bi-trash"></i>
                                  {{ $t('diagram.management.delete') }}
                                </button>
                              </li>
                            </ul>
                          </div>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Pagination -->
              <div v-if="totalPages > 1" class="d-flex justify-content-between align-items-center mt-4">
                <div>
                  <small class="text-muted">
                    {{ $t('diagram.management.showing', { 
                      start: (currentPage - 1) * pageSize + 1, 
                      end: Math.min(currentPage * pageSize, total),
                      total: total 
                    }) }}
                  </small>
                </div>
                <nav>
                  <ul class="pagination pagination-sm mb-0">
                    <li class="page-item" :class="{ disabled: currentPage === 1 }">
                      <button 
                        class="page-link"
                        @click="goToPage(currentPage - 1)"
                        :disabled="currentPage === 1"
                      >
                        {{ $t('diagram.management.previous') }}
                      </button>
                    </li>
                    <li 
                      v-for="page in visiblePages" 
                      :key="page"
                      class="page-item"
                      :class="{ active: page === currentPage }"
                    >
                      <button 
                        class="page-link"
                        @click="goToPage(page)"
                      >
                        {{ page }}
                      </button>
                    </li>
                    <li class="page-item" :class="{ disabled: currentPage === totalPages }">
                      <button 
                        class="page-link"
                        @click="goToPage(currentPage + 1)"
                        :disabled="currentPage === totalPages"
                      >
                        {{ $t('diagram.management.next') }}
                      </button>
                    </li>
                  </ul>
                </nav>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <DiagramCreateModal 
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="onDiagramCreated"
    />
    
    <DiagramEditModal 
      v-if="showEditModal && editingDiagram"
      :diagram="editingDiagram"
      @close="showEditModal = false"
      @updated="onDiagramUpdated"
    />
    
    <DiagramViewModal 
      v-if="showViewModal && viewingDiagram"
      :diagram="viewingDiagram"
      @close="showViewModal = false"
      @edit="editDiagram"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDiagramStore, type Diagram } from '@/stores/diagramStore'
import { formatDate } from '@/utils/dateUtils'
import DiagramCard from '@/components/diagram/DiagramCard.vue'
import DiagramCreateModal from '@/components/diagram/DiagramCreateModal.vue'
import DiagramEditModal from '@/components/diagram/DiagramEditModal.vue'
import DiagramViewModal from '@/components/diagram/DiagramViewModal.vue'

const { t } = useI18n()
const diagramStore = useDiagramStore()

// State
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const editingDiagram = ref<Diagram | null>(null)
const viewingDiagram = ref<Diagram | null>(null)
const viewMode = ref<'grid' | 'list'>('grid')

// Filters
const searchQuery = ref('')
const selectedType = ref('')
const selectedScriptType = ref('')
const selectedStatus = ref('')
const sortBy = ref('created_at')
const pageSize = ref(12)

// Computed from store
const diagrams = computed(() => diagramStore.diagrams)
const loading = computed(() => diagramStore.loading)
const error = computed(() => diagramStore.error)
const currentPage = computed(() => diagramStore.currentPage)
const totalPages = computed(() => diagramStore.totalPages)
const total = computed(() => diagramStore.total)
const stats = computed(() => diagramStore.stats)

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
const loadDiagrams = async () => {
  if (searchQuery.value.trim()) {
    await diagramStore.searchDiagrams(searchQuery.value, currentPage.value, pageSize.value)
  } else {
    await diagramStore.listDiagrams({
      page: currentPage.value,
      limit: pageSize.value,
      type: selectedType.value || undefined,
      script_type: selectedScriptType.value || undefined,
      status: selectedStatus.value || undefined
    })
  }
}

const loadStats = async () => {
  await diagramStore.getDiagramStats()
}

const onSearch = () => {
  diagramStore.goToPage(1)
  loadDiagrams()
}

const onFilterChange = () => {
  diagramStore.goToPage(1)
  loadDiagrams()
}

const goToPage = (page: number) => {
  diagramStore.goToPage(page)
  loadDiagrams()
}

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-secondary',
    published: 'bg-success',
    private: 'bg-warning',
    archived: 'bg-danger'
  }
  return classes[status] || 'bg-secondary'
}

// Diagram actions
const viewDiagram = (diagram: Diagram) => {
  viewingDiagram.value = diagram
  showViewModal.value = true
}

const editDiagram = (diagram: Diagram) => {
  editingDiagram.value = diagram
  showEditModal.value = true
}

const deleteDiagram = async (diagram: Diagram) => {
  if (!confirm(t('diagram.management.confirmDelete', { name: diagram.name }))) {
    return
  }
  
  await diagramStore.deleteDiagram(diagram.id)
  await loadDiagrams()
  await loadStats()
}

const publishDiagram = async (diagram: Diagram) => {
  await diagramStore.publishDiagram(diagram.id)
  await loadDiagrams()
  await loadStats()
}

const archiveDiagram = async (diagram: Diagram) => {
  await diagramStore.archiveDiagram(diagram.id)
  await loadDiagrams()
  await loadStats()
}

const shareDiagram = async (diagram: Diagram) => {
  await diagramStore.shareDiagram(diagram.id)
  await loadDiagrams()
}

const unshareDiagram = async (diagram: Diagram) => {
  await diagramStore.unshareDiagram(diagram.id)
  await loadDiagrams()
}

// Modal handlers
const onDiagramCreated = () => {
  showCreateModal.value = false
  loadDiagrams()
  loadStats()
}

const onDiagramUpdated = () => {
  showEditModal.value = false
  editingDiagram.value = null
  loadDiagrams()
  loadStats()
}

// Watchers
watch(() => diagramStore.currentPage, () => {
  loadDiagrams()
})

// Lifecycle
onMounted(() => {
  loadDiagrams()
  loadStats()
})
</script>

<style scoped>
.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.card {
  transition: box-shadow 0.2s ease-in-out;
}

.card:hover {
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.table th {
  border-top: none;
  font-weight: 600;
  color: #6c757d;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
}

.dropdown-menu {
  min-width: 150px;
}

.pagination-sm .page-link {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}
</style>
