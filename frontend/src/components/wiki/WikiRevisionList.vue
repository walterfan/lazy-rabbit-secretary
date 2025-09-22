<template>
  <div class="wiki-revision-list">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">{{ $t('common.loading') }}</span>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="revisions.length === 0" class="empty-state">
      <div class="empty-icon">
        <i class="bi bi-clock-history"></i>
      </div>
      <h4>{{ $t('wiki.noRevisionsFound') }}</h4>
      <p>{{ $t('wiki.noRevisionsDescription') }}</p>
    </div>

    <!-- Revisions List -->
    <div v-else class="revisions-list">
      <div 
        v-for="revision in revisions" 
        :key="revision.id"
        class="revision-item"
        :class="{ 'current': revision.version === currentVersion }"
      >
        <!-- Revision Header -->
        <div class="revision-header">
          <div class="revision-info">
            <div class="revision-title" @click="$emit('revision-click', revision)">
              {{ revision.title }}
            </div>
            <div class="revision-meta">
              <span class="version-badge">
                v{{ revision.version }}
              </span>
              <span class="revision-date">
                {{ formatDate(revision.created_at) }}
              </span>
              <span class="revision-author">
                {{ $t('wiki.by') }} {{ revision.created_by }}
              </span>
            </div>
          </div>
          
          <div class="revision-actions" v-if="showCompare">
            <button 
              class="btn btn-sm btn-outline-primary"
              @click="$emit('compare-revisions', revision)"
            >
              <i class="bi bi-arrow-left-right"></i>
              {{ $t('wiki.compare') }}
            </button>
          </div>
        </div>

        <!-- Revision Content -->
        <div class="revision-content">
          <div class="revision-summary">{{ revision.summary || $t('wiki.noChangeNote') }}</div>
          
          <div class="revision-stats">
            <span class="stat-item">
              <i class="bi bi-file-text"></i>
              {{ revision.content_size }} {{ $t('wiki.bytes') }}
            </span>
            <span class="stat-item">
              <i class="bi bi-list"></i>
              {{ revision.line_count }} {{ $t('wiki.lines') }}
            </span>
          </div>
        </div>

        <!-- Change Note -->
        <div v-if="revision.change_note" class="change-note">
          <div class="change-note-label">{{ $t('wiki.changeNote') }}:</div>
          <div class="change-note-content">{{ revision.change_note }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/dateUtils'
import type { WikiRevision } from '@/types'

interface Props {
  revisions: WikiRevision[]
  loading?: boolean
  showCompare?: boolean
  currentVersion?: number
}

defineProps<Props>()

const emit = defineEmits<{
  'revision-click': [revision: WikiRevision]
  'compare-revisions': [revision: WikiRevision]
}>()

const { t } = useI18n()
</script>

<style scoped>
.wiki-revision-list {
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

.revisions-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.revision-item {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s ease;
}

.revision-item:hover {
  border-color: #667eea;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.1);
}

.revision-item.current {
  border-color: #28a745;
  background-color: #f8fff9;
}

.revision-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.revision-info {
  flex: 1;
}

.revision-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  cursor: pointer;
  transition: color 0.2s ease;
  margin-bottom: 0.5rem;
}

.revision-title:hover {
  color: #667eea;
}

.revision-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.875rem;
  color: #6c757d;
}

.version-badge {
  background-color: #667eea;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-weight: 500;
  font-size: 0.75rem;
}

.revision-actions {
  flex-shrink: 0;
}

.revision-content {
  margin-bottom: 1rem;
}

.revision-summary {
  color: #6c757d;
  line-height: 1.5;
  margin-bottom: 1rem;
}

.revision-stats {
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

.change-note {
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.change-note-label {
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.5rem;
}

.change-note-content {
  color: #6c757d;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .revision-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .revision-meta {
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  
  .revision-stats {
    flex-wrap: wrap;
  }
}
</style>
