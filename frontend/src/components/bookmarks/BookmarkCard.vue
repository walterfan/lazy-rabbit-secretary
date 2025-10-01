<template>
  <div class="bookmark-card card h-100">
    <div class="card-body">
      <div class="d-flex justify-content-between align-items-start mb-2">
        <h5 class="card-title mb-1">
          <a
            :href="bookmark.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-decoration-none"
            :title="bookmark.title"
          >
            {{ truncateTitle(bookmark.title) }}
            <i class="bi bi-box-arrow-up-right ms-1 small"></i>
          </a>
        </h5>
        <div class="dropdown" v-if="showActions">
          <button
            class="btn btn-sm btn-outline-secondary dropdown-toggle"
            type="button"
            data-bs-toggle="dropdown"
          >
            <i class="bi bi-three-dots"></i>
          </button>
          <ul class="dropdown-menu">
            <li>
              <button class="dropdown-item" @click="$emit('edit', bookmark)">
                <i class="bi bi-pencil me-2"></i>Edit
              </button>
            </li>
            <li>
              <button class="dropdown-item" @click="copyUrl">
                <i class="bi bi-clipboard me-2"></i>Copy URL
              </button>
            </li>
            <li><hr class="dropdown-divider"></li>
            <li>
              <button class="dropdown-item text-danger" @click="confirmDelete">
                <i class="bi bi-trash me-2"></i>Delete
              </button>
            </li>
          </ul>
        </div>
      </div>

      <p class="card-text text-muted small mb-2">
        <i class="bi bi-link-45deg me-1"></i>
        {{ getDomainFromUrl(bookmark.url) }}
      </p>

      <p v-if="bookmark.description" class="card-text">
        {{ truncateDescription(bookmark.description) }}
      </p>

      <div v-if="bookmark.tags.length > 0" class="mb-2">
        <span
          v-for="tag in bookmark.tags"
          :key="tag.id"
          class="badge bg-light text-dark me-1 mb-1"
          @click="$emit('tag-click', tag.name)"
          style="cursor: pointer;"
          :title="`Filter by ${tag.name}`"
        >
          <i class="bi bi-tag me-1"></i>{{ tag.name }}
        </span>
      </div>


    </div>

    <div v-if="bookmark.category_id" class="card-footer bg-transparent">
      <small class="text-muted">
        <i class="bi bi-folder me-1"></i>
        <span
          @click="$emit('category-click', bookmark.category_id)"
          style="cursor: pointer;"
          class="text-primary"
          :title="`Filter by category`"
        >
          {{ getCategoryName(bookmark.category_id) }}
        </span>
      </small>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Bookmark } from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';

interface Props {
  bookmark: Bookmark;
  showActions?: boolean;
  maxTitleLength?: number;
  maxDescriptionLength?: number;
}

interface Emits {
  (e: 'edit', bookmark: Bookmark): void;
  (e: 'delete', bookmark: Bookmark): void;
  (e: 'tag-click', tagName: string): void;
  (e: 'category-click', categoryId: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true,
  maxTitleLength: 50,
  maxDescriptionLength: 100
});

const emit = defineEmits<Emits>();

const bookmarkStore = useBookmarkStore();

const categories = computed(() => bookmarkStore.categories);

const truncateTitle = (title: string): string => {
  if (title.length <= props.maxTitleLength) {
    return title;
  }
  return title.substring(0, props.maxTitleLength) + '...';
};

const truncateDescription = (description: string): string => {
  if (description.length <= props.maxDescriptionLength) {
    return description;
  }
  return description.substring(0, props.maxDescriptionLength) + '...';
};

const getDomainFromUrl = (url: string): string => {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname;
  } catch {
    return url;
  }
};

const getCategoryName = (categoryId: string): string => {
  const category = categories.value.find(c => c.id.toString() === categoryId);
  return category ? category.name : 'Unknown Category';
};

const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(new Date(date));
};

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(props.bookmark.url);
    // You could show a toast notification here
    console.log('URL copied to clipboard');
  } catch (error) {
    console.error('Failed to copy URL:', error);
    // Fallback for older browsers
    const textArea = document.createElement('textarea');
    textArea.value = props.bookmark.url;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
  }
};

const confirmDelete = () => {
  if (confirm(`Are you sure you want to delete "${props.bookmark.title}"?`)) {
    emit('delete', props.bookmark);
  }
};
</script>

<style scoped>
.bookmark-card {
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
  max-width: 320px;
  margin: 0 auto;
}

.bookmark-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.card-title a {
  color: #0d6efd;
  font-weight: 500;
}

.card-title a:hover {
  color: #0a58ca;
  text-decoration: underline !important;
}

.badge {
  transition: background-color 0.2s ease-in-out;
}

.badge:hover {
  background-color: #e9ecef !important;
  color: #495057 !important;
}

.dropdown-toggle::after {
  display: none;
}

.card-footer {
  border-top: 1px solid rgba(0, 0, 0, 0.125);
  padding: 0.5rem 1rem;
}

.text-primary {
  transition: color 0.2s ease-in-out;
}

.text-primary:hover {
  color: #0a58ca !important;
}
</style>
