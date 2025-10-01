<template>
  <div class="bookmark-form">
    <form @submit.prevent="handleSubmit">
      <div class="row">
        <div class="col-md-8">
          <div class="mb-3">
            <label for="url" class="form-label">URL <span class="text-danger">*</span></label>
            <input
              id="url"
              v-model="formData.url"
              type="url"
              class="form-control"
              :class="{ 'is-invalid': errors.url }"
              placeholder="https://example.com"
              required
            />
            <div v-if="errors.url" class="invalid-feedback">
              {{ errors.url }}
            </div>
          </div>

          <div class="mb-3">
            <label for="title" class="form-label">Title <span class="text-danger">*</span></label>
            <input
              id="title"
              v-model="formData.title"
              type="text"
              class="form-control"
              :class="{ 'is-invalid': errors.title }"
              placeholder="Bookmark title"
              required
            />
            <div v-if="errors.title" class="invalid-feedback">
              {{ errors.title }}
            </div>
          </div>

          <div class="mb-3">
            <label for="description" class="form-label">Description</label>
            <textarea
              id="description"
              v-model="formData.description"
              class="form-control"
              :class="{ 'is-invalid': errors.description }"
              rows="3"
              placeholder="Optional description"
            ></textarea>
            <div v-if="errors.description" class="invalid-feedback">
              {{ errors.description }}
            </div>
          </div>
        </div>

        <div class="col-md-4">
          <div class="mb-3">
            <label for="category" class="form-label">Category</label>
            <select
              id="category"
              v-model="formData.category_id"
              class="form-select"
              :class="{ 'is-invalid': errors.category_id }"
            >
              <option value="">Select a category</option>
              <option
                v-for="category in categories"
                :key="category.id"
                :value="category.id.toString()"
              >
                {{ category.name }}
              </option>
            </select>
            <div v-if="errors.category_id" class="invalid-feedback">
              {{ errors.category_id }}
            </div>
          </div>

          <div class="mb-3">
            <label for="tags" class="form-label">Tags</label>
            <div class="input-group">
              <input
                id="newTag"
                v-model="newTag"
                type="text"
                class="form-control"
                placeholder="Add a tag"
                @keydown.enter.prevent="addTag"
              />
              <button
                type="button"
                class="btn btn-outline-secondary"
                @click="addTag"
                :disabled="!newTag.trim()"
              >
                <i class="bi bi-plus"></i>
              </button>
            </div>
            <div class="mt-2">
              <span
                v-for="(tag, index) in formData.tags"
                :key="index"
                class="badge bg-secondary me-1 mb-1"
              >
                {{ tag }}
                <button
                  type="button"
                  class="btn-close btn-close-white ms-1"
                  @click="removeTag(index)"
                  style="font-size: 0.7em"
                ></button>
              </span>
            </div>
            <!-- Popular tags suggestions -->
            <div v-if="popularTags.length > 0" class="mt-2">
              <small class="text-muted">Popular tags:</small>
              <div class="mt-1">
                <button
                  v-for="tag in popularTags"
                  :key="tag.name"
                  type="button"
                  class="btn btn-outline-info btn-sm me-1 mb-1"
                  @click="addPopularTag(tag.name)"
                  :disabled="(formData.tags || []).includes(tag.name)"
                >
                  {{ tag.name }} ({{ tag.count }})
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="d-flex justify-content-between">
        <button
          type="button"
          class="btn btn-secondary"
          @click="$emit('cancel')"
        >
          Cancel
        </button>
        <div>
          <button
            v-if="mode === 'edit'"
            type="button"
            class="btn btn-danger me-2"
            @click="handleDelete"
            :disabled="loading"
          >
            <i class="bi bi-trash"></i> Delete
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="loading || !isFormValid"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
            <i v-else class="bi" :class="mode === 'create' ? 'bi-plus' : 'bi-check'"></i>
            {{ mode === 'create' ? 'Create Bookmark' : 'Update Bookmark' }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import type { 
  Bookmark, 
  BookmarkCategory, 
  CreateBookmarkRequest, 
  UpdateBookmarkRequest 
} from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';

interface Props {
  bookmark?: Bookmark;
  mode: 'create' | 'edit';
  loading?: boolean;
}

interface Emits {
  (e: 'submit', data: CreateBookmarkRequest | UpdateBookmarkRequest): void;
  (e: 'delete'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
});

const emit = defineEmits<Emits>();

const bookmarkStore = useBookmarkStore();

const formData = ref<CreateBookmarkRequest>({
  url: '',
  title: '',
  description: '',
  tags: [],
  category_id: ''
});

const newTag = ref('');
const errors = ref<Record<string, string>>({});

const categories = computed(() => bookmarkStore.categories);
const popularTags = computed(() => bookmarkStore.popularTags);

const isFormValid = computed(() => {
  return formData.value.url.trim() !== '' && 
         formData.value.title.trim() !== '' &&
         isValidUrl(formData.value.url);
});

const isValidUrl = (url: string): boolean => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

const validateForm = () => {
  errors.value = {};

  if (!formData.value.url.trim()) {
    errors.value.url = 'URL is required';
  } else if (!isValidUrl(formData.value.url)) {
    errors.value.url = 'Please enter a valid URL';
  }

  if (!formData.value.title.trim()) {
    errors.value.title = 'Title is required';
  } else if (formData.value.title.length > 200) {
    errors.value.title = 'Title must be less than 200 characters';
  }

  if (formData.value.description && formData.value.description.length > 1000) {
    errors.value.description = 'Description must be less than 1000 characters';
  }

  return Object.keys(errors.value).length === 0;
};

const addTag = () => {
  const tag = newTag.value.trim().toLowerCase();
  if (tag && !(formData.value.tags || []).includes(tag)) {
    if (!formData.value.tags) {
      formData.value.tags = [];
    }
    formData.value.tags.push(tag);
    newTag.value = '';
  }
};

const addPopularTag = (tagName: string) => {
  if (!formData.value.tags?.includes(tagName)) {
    if (!formData.value.tags) {
      formData.value.tags = [];
    }
    formData.value.tags.push(tagName);
  }
};

const removeTag = (index: number) => {
  if (formData.value.tags) {
    formData.value.tags.splice(index, 1);
  }
};

const handleSubmit = () => {
  if (!validateForm()) {
    return;
  }

  emit('submit', formData.value);
};

const handleDelete = () => {
  if (confirm('Are you sure you want to delete this bookmark?')) {
    emit('delete');
  }
};

// Initialize form data when bookmark prop changes
watch(
  () => props.bookmark,
  (bookmark) => {
    if (bookmark) {
      formData.value = {
        url: bookmark.url,
        title: bookmark.title,
        description: bookmark.description,
        tags: bookmark.tags.map(tag => tag.name),
        category_id: bookmark.category_id
      };
    } else {
      formData.value = {
        url: '',
        title: '',
        description: '',
        tags: [],
        category_id: ''
      };
    }
  },
  { immediate: true }
);

// Load data on mount
onMounted(async () => {
  try {
    await Promise.all([
      bookmarkStore.fetchCategories(),
      bookmarkStore.fetchPopularTags(10)
    ]);
  } catch (error) {
    console.error('Failed to load form data:', error);
  }
});
</script>

<style scoped>
.bookmark-form {
  max-width: 100%;
}

.badge {
  display: inline-flex;
  align-items: center;
}

.btn-close {
  --bs-btn-close-bg: none;
}

.form-control:focus,
.form-select:focus {
  border-color: #86b7fe;
  box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
}

.invalid-feedback {
  display: block;
}
</style>
