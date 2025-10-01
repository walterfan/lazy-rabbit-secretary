<template>
  <div class="wiki-page-actions">
    <div class="action-buttons">
      <!-- History Button -->
      <button 
        class="btn btn-outline-secondary"
        @click="$emit('history')"
        :title="$t('wiki.history')"
      >
        <i class="bi bi-clock-history"></i>
      </button>

      <!-- Edit Button -->
      <button 
        v-if="canEdit"
        class="btn btn-primary"
        @click="$emit('edit')"
        :disabled="page.is_locked"
        :title="$t('wiki.edit')"
      >
        <i class="bi bi-pencil"></i>
      </button>

      <!-- Lock/Unlock Button -->
      <button
        v-if="isAuthenticated"
        class="btn"
        :class="page.is_locked ? 'btn-warning' : 'btn-outline-warning'"
        @click="page.is_locked ? $emit('unlock') : $emit('lock')"
        :title="page.is_locked ? $t('wiki.unlock') : $t('wiki.lock')"
      >
        <i :class="page.is_locked ? 'bi bi-unlock' : 'bi bi-lock'"></i>
      </button>

      <!-- More Actions Dropdown -->
      <div class="dropdown">
        <button 
          class="btn btn-outline-secondary"
          type="button"
          data-bs-toggle="dropdown"
        >
          <i class="bi bi-three-dots"></i>
        </button>
        <ul class="dropdown-menu">
          <li>
            <button 
              class="dropdown-item"
              @click="$emit('view-source')"
            >
              <i class="bi bi-code"></i>
              {{ $t('wiki.viewSource') }}
            </button>
          </li>
          <li>
            <button 
              class="dropdown-item"
              @click="$emit('print')"
            >
              <i class="bi bi-printer"></i>
              {{ $t('wiki.print') }}
            </button>
          </li>
          <li>
            <button 
              class="dropdown-item"
              @click="$emit('share')"
            >
              <i class="bi bi-share"></i>
              {{ $t('wiki.share') }}
            </button>
          </li>
          <li><hr class="dropdown-divider"></li>
          <li>
            <button 
              v-if="canDelete"
              class="dropdown-item text-danger"
              @click="$emit('delete')"
              :disabled="page.is_locked"
            >
              <i class="bi bi-trash"></i>
              {{ $t('wiki.delete') }}
            </button>
          </li>
        </ul>
      </div>
    </div>

    <!-- Page Status Indicators -->
    <div class="status-indicators">
      <span v-if="page.is_locked" class="badge bg-danger">
        <i class="bi bi-lock"></i>
        {{ $t('wiki.locked') }}
      </span>
      <span v-if="page.is_protected" class="badge bg-warning">
        <i class="bi bi-shield-lock"></i>
        {{ $t('wiki.protected') }}
      </span>
      <span class="badge" :class="getStatusBadgeClass(page.status)">
        {{ $t(`wiki.status.${page.status}`) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/authStore'
import type { WikiPage } from '@/types'

interface Props {
  page: WikiPage
  canEdit?: boolean
  canDelete?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  edit: []
  history: []
  lock: []
  unlock: []
  delete: []
  'view-source': []
  print: []
  share: []
}>()

const { t } = useI18n()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-secondary',
    published: 'bg-success',
    archived: 'bg-warning',
    protected: 'bg-info'
  }
  return classes[status] || 'bg-secondary'
}
</script>

<style scoped>
.wiki-page-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.status-indicators {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .wiki-page-actions {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .action-buttons {
    width: 100%;
    justify-content: flex-start;
  }
  
  .status-indicators {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
