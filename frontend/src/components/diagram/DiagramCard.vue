<template>
  <div class="diagram-card card h-100">
    <div class="card-body">
      <!-- Header -->
      <div class="d-flex justify-content-between align-items-start mb-2">
        <h6 class="card-title mb-0">{{ diagram.name }}</h6>
        <div class="dropdown">
          <button 
            class="btn btn-sm btn-outline-secondary"
            type="button"
            data-bs-toggle="dropdown"
          >
            <i class="bi bi-three-dots"></i>
          </button>
          <ul class="dropdown-menu">
            <li>
              <button 
                class="dropdown-item"
                @click="$emit('view', diagram)"
              >
                <i class="bi bi-eye"></i>
                {{ $t('diagram.management.view') }}
              </button>
            </li>
            <li>
              <button 
                class="dropdown-item"
                @click="$emit('edit', diagram)"
              >
                <i class="bi bi-pencil"></i>
                {{ $t('diagram.management.edit') }}
              </button>
            </li>
            <li v-if="diagram.status === 'draft'">
              <button 
                class="dropdown-item"
                @click="$emit('publish', diagram)"
              >
                <i class="bi bi-check-circle"></i>
                {{ $t('diagram.management.publish') }}
              </button>
            </li>
            <li v-if="diagram.status === 'published'">
              <button 
                class="dropdown-item"
                @click="$emit('archive', diagram)"
              >
                <i class="bi bi-archive"></i>
                {{ $t('diagram.management.archive') }}
              </button>
            </li>
            <li>
              <button 
                class="dropdown-item"
                @click="diagram.shared ? $emit('unshare', diagram) : $emit('share', diagram)"
              >
                <i :class="diagram.shared ? 'bi bi-share-fill' : 'bi bi-share'"></i>
                {{ diagram.shared ? $t('diagram.management.unshare') : $t('diagram.management.share') }}
              </button>
            </li>
            <li><hr class="dropdown-divider"></li>
            <li>
              <button 
                class="dropdown-item text-danger"
                @click="$emit('delete', diagram)"
              >
                <i class="bi bi-trash"></i>
                {{ $t('diagram.management.delete') }}
              </button>
            </li>
          </ul>
        </div>
      </div>

      <!-- Description -->
      <p class="card-text text-muted small mb-3">
        {{ diagram.description || $t('diagram.management.noDescription') }}
      </p>

      <!-- Tags -->
      <div v-if="diagram.tags && diagram.tags.length > 0" class="mb-3">
        <span 
          v-for="tag in diagram.tags.slice(0, 3)" 
          :key="tag"
          class="badge bg-light text-dark me-1"
        >
          {{ tag }}
        </span>
        <span 
          v-if="diagram.tags.length > 3"
          class="badge bg-secondary"
        >
          +{{ diagram.tags.length - 3 }}
        </span>
      </div>

      <!-- Metadata -->
      <div class="row g-2 mb-3">
        <div class="col-6">
          <small class="text-muted d-block">{{ $t('diagram.management.type') }}</small>
          <span class="badge bg-secondary">{{ diagram.type }}</span>
        </div>
        <div class="col-6">
          <small class="text-muted d-block">{{ $t('diagram.management.scriptType') }}</small>
          <span class="badge bg-info">{{ diagram.script_type }}</span>
        </div>
      </div>

      <!-- Status and Stats -->
      <div class="d-flex justify-content-between align-items-center mb-3">
        <span class="badge" :class="getStatusBadgeClass(diagram.status)">
          {{ $t(`diagram.management.${diagram.status}`) }}
        </span>
        <div class="text-muted small">
          <i class="bi bi-eye"></i>
          {{ diagram.view_count }}
        </div>
      </div>

      <!-- Actions -->
      <div class="d-grid gap-2">
        <button 
          class="btn btn-primary btn-sm"
          @click="$emit('view', diagram)"
        >
          <i class="bi bi-eye"></i>
          {{ $t('diagram.management.view') }}
        </button>
        <div class="btn-group btn-group-sm">
          <button 
            class="btn btn-outline-secondary"
            @click="$emit('edit', diagram)"
            :title="$t('diagram.management.edit')"
          >
            <i class="bi bi-pencil"></i>
          </button>
          <button 
            v-if="diagram.status === 'draft'"
            class="btn btn-outline-success"
            @click="$emit('publish', diagram)"
            :title="$t('diagram.management.publish')"
          >
            <i class="bi bi-check-circle"></i>
          </button>
          <button 
            v-if="diagram.status === 'published'"
            class="btn btn-outline-warning"
            @click="$emit('archive', diagram)"
            :title="$t('diagram.management.archive')"
          >
            <i class="bi bi-archive"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="card-footer bg-transparent">
      <div class="d-flex justify-content-between align-items-center">
        <small class="text-muted">
          {{ formatDate(diagram.created_at) }}
        </small>
        <div class="d-flex gap-1">
          <i 
            v-if="diagram.public"
            class="bi bi-globe text-primary"
            :title="$t('diagram.management.public')"
          ></i>
          <i 
            v-if="diagram.shared"
            class="bi bi-share-fill text-success"
            :title="$t('diagram.management.shared')"
          ></i>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/dateUtils'
import type { Diagram } from '@/stores/diagramStore'

interface Props {
  diagram: Diagram
}

defineProps<Props>()

const emit = defineEmits<{
  view: [diagram: Diagram]
  edit: [diagram: Diagram]
  delete: [diagram: Diagram]
  publish: [diagram: Diagram]
  archive: [diagram: Diagram]
  share: [diagram: Diagram]
  unshare: [diagram: Diagram]
}>()

const { t } = useI18n()

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-secondary',
    published: 'bg-success',
    private: 'bg-warning',
    archived: 'bg-danger'
  }
  return classes[status] || 'bg-secondary'
}
</script>

<style scoped>
.diagram-card {
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.diagram-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.card-title {
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.2;
}

.card-text {
  font-size: 0.875rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.badge {
  font-size: 0.75rem;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
}

.dropdown-menu {
  min-width: 150px;
}

.card-footer {
  border-top: 1px solid rgba(0,0,0,0.05);
  padding: 0.75rem 1rem;
}

.text-muted {
  font-size: 0.75rem;
}
</style>
