<template>
  <div class="wiki-search-bar">
    <div class="input-group">
      <input
        type="text"
        class="form-control"
        :placeholder="placeholder"
        v-model="query"
        @keyup.enter="handleSearch"
        @input="handleInput"
        :disabled="loading"
        ref="searchInput"
      />
      <button 
        class="btn btn-outline-secondary"
        type="button"
        @click="handleSearch"
        :disabled="loading || !query.trim()"
      >
        <i v-if="loading" class="spinner-border spinner-border-sm"></i>
        <i v-else class="bi bi-search"></i>
      </button>
    </div>
    
    <!-- Search Suggestions -->
    <div v-if="showSuggestions && suggestions.length > 0" class="search-suggestions">
      <div 
        v-for="suggestion in suggestions" 
        :key="suggestion.slug"
        class="suggestion-item"
        @click="selectSuggestion(suggestion)"
      >
        <div class="suggestion-title">{{ suggestion.title }}</div>
        <div class="suggestion-summary">{{ suggestion.summary }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useWikiStore } from '@/stores/wikiStore'

interface Props {
  placeholder?: string
  showFilters?: boolean
  autoFocus?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Search wiki pages...',
  showFilters: false,
  autoFocus: false,
  loading: false
})

const emit = defineEmits<{
  search: [query: string]
  filter: [filters: any]
}>()

const wikiStore = useWikiStore()
const query = ref('')
const suggestions = ref<any[]>([])
const showSuggestions = ref(false)
const searchInput = ref<HTMLInputElement>()

// Debounced search
let searchTimeout: NodeJS.Timeout

const handleInput = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  
  if (query.value.trim()) {
    searchTimeout = setTimeout(() => {
      // Implement suggestions logic here
      showSuggestions.value = true
    }, 300)
  } else {
    showSuggestions.value = false
  }
}

const handleSearch = () => {
  if (query.value.trim()) {
    emit('search', query.value.trim())
    showSuggestions.value = false
  }
}

const selectSuggestion = (suggestion: any) => {
  query.value = suggestion.title
  showSuggestions.value = false
  handleSearch()
}

// Auto focus
onMounted(() => {
  if (props.autoFocus && searchInput.value) {
    nextTick(() => {
      searchInput.value?.focus()
    })
  }
})

// Hide suggestions when clicking outside
const handleClickOutside = (event: Event) => {
  const target = event.target as HTMLElement
  if (!target.closest('.wiki-search-bar')) {
    showSuggestions.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

// Cleanup
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
})
</script>

<style scoped>
.wiki-search-bar {
  position: relative;
}

.search-suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #dee2e6;
  border-top: none;
  border-radius: 0 0 8px 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-height: 300px;
  overflow-y: auto;
}

.suggestion-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-bottom: 1px solid #f8f9fa;
  transition: background-color 0.2s ease;
}

.suggestion-item:hover {
  background-color: #f8f9fa;
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-title {
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.25rem;
}

.suggestion-summary {
  font-size: 0.875rem;
  color: #6c757d;
  line-height: 1.4;
}
</style>
