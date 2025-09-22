<template>
  <div class="wiki-sidebar">
    <!-- Categories Section -->
    <div v-if="categories.length > 0" class="sidebar-section">
      <h5>{{ $t('wiki.categories') }}</h5>
      <div class="category-list">
        <div 
          v-for="category in categories" 
          :key="category.name"
          class="category-item"
          @click="$emit('category-click', category.name)"
        >
          <span class="category-name">{{ category.name }}</span>
          <span class="category-count">{{ category.count }}</span>
        </div>
      </div>
    </div>

    <!-- Tags Section -->
    <div v-if="tags.length > 0" class="sidebar-section">
      <h5>{{ $t('wiki.tags') }}</h5>
      <div class="tag-list">
        <span 
          v-for="tag in tags" 
          :key="tag.name"
          class="badge bg-secondary me-1 mb-1 tag-item"
          @click="$emit('tag-click', tag.name)"
        >
          {{ tag.name }} ({{ tag.count }})
        </span>
      </div>
    </div>

    <!-- Recent Pages Section -->
    <div v-if="recentPages.length > 0" class="sidebar-section">
      <h5>{{ $t('wiki.recentPages') }}</h5>
      <div class="recent-pages">
        <div 
          v-for="page in recentPages" 
          :key="page.id"
          class="recent-page-item"
          @click="$emit('page-click', page)"
        >
          <div class="page-title">{{ page.title }}</div>
          <div class="page-meta">{{ formatDate(page.updated_at) }}</div>
        </div>
      </div>
    </div>

    <!-- Special Pages Section -->
    <div v-if="specialPages.length > 0" class="sidebar-section">
      <h5>{{ $t('wiki.specialPages') }}</h5>
      <div class="special-pages">
        <div 
          v-for="specialPage in specialPages" 
          :key="specialPage.type"
          class="special-page-item"
          @click="$emit('special-page-click', specialPage.type)"
        >
          <i :class="getSpecialPageIcon(specialPage.type)"></i>
          <span>{{ specialPage.title }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/dateUtils'

interface Props {
  categories: Array<{ name: string; count: number }>
  tags: Array<{ name: string; count: number }>
  recentPages: any[]
  specialPages: any[]
}

defineProps<Props>()

const emit = defineEmits<{
  'category-click': [category: string]
  'tag-click': [tag: string]
  'page-click': [page: any]
  'special-page-click': [type: string]
}>()

const { t } = useI18n()

const getSpecialPageIcon = (type: string) => {
  const icons: Record<string, string> = {
    'orphaned': 'bi bi-file-x',
    'wanted': 'bi bi-file-plus',
    'dead-end': 'bi bi-file-minus',
    'recent-changes': 'bi bi-clock-history',
    'random': 'bi bi-shuffle'
  }
  return icons[type] || 'bi bi-file'
}
</script>

<style scoped>
.wiki-sidebar {
  position: sticky;
  top: 2rem;
}

.sidebar-section {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 1.5rem;
}

.sidebar-section h5 {
  margin-bottom: 1rem;
  color: #333;
  font-weight: 600;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.category-item:hover {
  border-color: #667eea;
  background-color: #f8f9fa;
}

.category-name {
  font-weight: 500;
  color: #333;
}

.category-count {
  background-color: #667eea;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag-item {
  cursor: pointer;
  transition: all 0.2s ease;
}

.tag-item:hover {
  background-color: #667eea !important;
  color: white !important;
}

.recent-pages {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.recent-page-item {
  padding: 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recent-page-item:hover {
  border-color: #667eea;
  background-color: #f8f9fa;
}

.page-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}

.page-meta {
  font-size: 0.75rem;
  color: #6c757d;
}

.special-pages {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.special-page-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.special-page-item:hover {
  border-color: #667eea;
  background-color: #f8f9fa;
}

.special-page-item i {
  color: #667eea;
  font-size: 1.1rem;
}

.special-page-item span {
  font-weight: 500;
  color: #333;
}
</style>
