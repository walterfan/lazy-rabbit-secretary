<template>
  <div class="wiki-page-list">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">{{ $t('common.loading') }}</span>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="pages.length === 0" class="empty-state">
      <div class="empty-icon">
        <i class="bi bi-file-text"></i>
      </div>
      <h4>{{ $t('wiki.noPagesFound') }}</h4>
      <p>{{ $t('wiki.noPagesDescription') }}</p>
      <button v-if="showActions" class="btn btn-primary" @click="$emit('create-page')">
        <i class="bi bi-plus-lg"></i>
        {{ $t('wiki.createFirstPage') }}
      </button>
    </div>

    <!-- Pages List -->
    <div v-else class="pages-grid">
      <div 
        v-for="page in pages" 
        :key="page.id"
        class="page-card"
        :class="{ 'locked': page.is_locked, 'protected': page.is_protected }"
      >
        <!-- Page Header -->
        <div class="page-header">
          <div class="page-title" @click="$emit('page-click', page)">
            {{ page.title }}
          </div>
          <div class="page-actions" v-if="showActions">
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
                    @click="$emit('edit-page', page)"
                    :disabled="page.is_locked"
                  >
                    <i class="bi bi-pencil"></i>
                    {{ $t('wiki.edit') }}
                  </button>
                </li>
                <li>
                  <button 
                    class="dropdown-item"
                    @click="$emit('view-history', page)"
                  >
                    <i class="bi bi-clock-history"></i>
                    {{ $t('wiki.history') }}
                  </button>
                </li>
                <li><hr class="dropdown-divider"></li>
                <li>
                  <button 
                    class="dropdown-item text-danger"
                    @click="$emit('delete-page', page)"
                    :disabled="page.is_locked"
                  >
                    <i class="bi bi-trash"></i>
                    {{ $t('wiki.delete') }}
                  </button>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Page Content -->
        <div class="page-content">
          <div class="page-summary">{{ page.summary || $t('wiki.noSummary') }}</div>
          
          <!-- Categories and Tags -->
          <div v-if="page.categories.length > 0 || page.tags.length > 0" class="page-meta">
            <div v-if="page.categories.length > 0" class="categories">
              <span 
                v-for="category in page.categories.slice(0, 3)" 
                :key="category"
                class="badge bg-primary me-1"
              >
                {{ category }}
              </span>
              <span v-if="page.categories.length > 3" class="text-muted">
                +{{ page.categories.length - 3 }} more
              </span>
            </div>
            
            <div v-if="page.tags.length > 0" class="tags">
              <span 
                v-for="tag in page.tags.slice(0, 3)" 
                :key="tag"
                class="badge bg-secondary me-1"
              >
                {{ tag }}
              </span>
              <span v-if="page.tags.length > 3" class="text-muted">
                +{{ page.tags.length - 3 }} more
              </span>
            </div>
          </div>
        </div>

        <!-- Page Footer -->
        <div class="page-footer">
          <div class="page-stats">
            <span class="stat-item">
              <i class="bi bi-eye"></i>
              {{ page.view_count }}
            </span>
            <span class="stat-item">
              <i class="bi bi-pencil"></i>
              {{ page.edit_count }}
            </span>
            <span class="stat-item">
              <i class="bi bi-clock"></i>
              {{ formatDate(page.updated_at) }}
            </span>
          </div>
          
          <div class="page-status">
            <span class="badge" :class="getStatusBadgeClass(page.status)">
              {{ $t(`wiki.status.${page.status}`) }}
            </span>
            <span v-if="page.is_locked" class="badge bg-danger">
              <i class="bi bi-lock"></i>
              {{ $t('wiki.locked') }}
            </span>
            <span v-if="page.is_protected" class="badge bg-warning">
              <i class="bi bi-shield-lock"></i>
              {{ $t('wiki.protected') }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/dateUtils'
import type { WikiPage } from '@/types'

interface Props {
  pages: WikiPage[]
  loading?: boolean
  showActions?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  'page-click': [page: WikiPage]
  'edit-page': [page: WikiPage]
  'delete-page': [page: WikiPage]
  'view-history': [page: WikiPage]
  'create-page': []
}>()

const { t } = useI18n()

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
.wiki-page-list {
  min-height: 200px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.empty-state {
  text-align: center;
  padding: 3rem 2rem;
  color: #6c757d;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h4 {
  margin-bottom: 1rem;
  color: #495057;
}

.empty-state p {
  margin-bottom: 2rem;
}

.pages-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.page-card {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s ease;
  cursor: pointer;
}

.page-card:hover {
  border-color: #667eea;
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
}

.page-card.locked {
  border-color: #dc3545;
  background-color: #fff5f5;
}

.page-card.protected {
  border-color: #ffc107;
  background-color: #fffbf0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
  line-height: 1.3;
  cursor: pointer;
  transition: color 0.2s ease;
}

.page-title:hover {
  color: #667eea;
}

.page-actions {
  flex-shrink: 0;
}

.page-content {
  margin-bottom: 1rem;
}

.page-summary {
  color: #6c757d;
  line-height: 1.5;
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.page-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.categories,
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.page-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #f8f9fa;
}

.page-stats {
  display: flex;
  gap: 1rem;
  font-size: 0.875rem;
  color: #6c757d;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.page-status {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .pages-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .page-card {
    padding: 1rem;
  }
  
  .page-footer {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .page-stats {
    flex-wrap: wrap;
  }
}
</style>
