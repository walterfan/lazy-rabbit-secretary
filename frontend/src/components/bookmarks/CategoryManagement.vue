<template>
  <div class="category-management">
    <!-- Header -->
    <div class="d-flex justify-content-between align-items-center mb-4">
      <div>
        <h4 class="mb-0">
          <i class="bi bi-folder me-2"></i>
          Bookmark Categories
        </h4>
        <p class="text-muted mb-0">Organize your bookmarks with categories</p>
      </div>
      <button
        class="btn btn-primary"
        @click="showCreateModal = true"
      >
        <i class="bi bi-plus me-2"></i>
        Add Category
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading && categories.length === 0" class="text-center py-4">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2">Loading categories...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger">
      <i class="bi bi-exclamation-triangle me-2"></i>
      {{ error }}
      <button class="btn btn-sm btn-outline-danger ms-2" @click="loadCategories">
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="categories.length === 0" class="text-center py-5">
      <i class="bi bi-folder-plus display-1 text-muted"></i>
      <h5 class="mt-3">No categories yet</h5>
      <p class="text-muted">Create your first category to organize bookmarks</p>
      <button
        class="btn btn-primary"
        @click="showCreateModal = true"
      >
        <i class="bi bi-plus me-2"></i>Create Category
      </button>
    </div>

    <!-- Categories List -->
    <div v-else class="row">
      <div class="col-12">
        <div class="card">
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Parent Category</th>
                    <th>Created</th>
                    <th>Updated</th>
                    <th width="120">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="category in categories" :key="category.id">
                    <td>
                      <div class="d-flex align-items-center">
                        <i class="bi bi-folder me-2 text-warning"></i>
                        <strong>{{ category.name }}</strong>
                      </div>
                    </td>
                    <td>
                      <span v-if="category.parent_id" class="text-muted">
                        {{ getParentCategoryName(category.parent_id) }}
                      </span>
                      <span v-else class="text-muted">—</span>
                    </td>
                    <td>
                      <small class="text-muted">
                        {{ formatDate(category.created_at) }}
                      </small>
                    </td>
                    <td>
                      <small class="text-muted">
                        {{ formatDate(category.updated_at) }}
                      </small>
                    </td>
                    <td>
                      <div class="btn-group btn-group-sm">
                        <button
                          class="btn btn-outline-primary"
                          @click="editCategory(category)"
                          :title="'Edit ' + category.name"
                        >
                          <i class="bi bi-pencil"></i>
                        </button>
                        <button
                          class="btn btn-outline-danger"
                          @click="confirmDelete(category)"
                          :title="'Delete ' + category.name"
                        >
                          <i class="bi bi-trash"></i>
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <div
      class="modal fade"
      :class="{ show: showCreateModal }"
      :style="{ display: showCreateModal ? 'block' : 'none' }"
      tabindex="-1"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-plus me-2"></i>
              Create Category
            </h5>
            <button
              type="button"
              class="btn-close"
              @click="closeCreateModal"
            ></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleCreate">
              <div class="mb-3">
                <label for="categoryName" class="form-label">
                  Category Name <span class="text-danger">*</span>
                </label>
                <input
                  id="categoryName"
                  v-model="createForm.name"
                  type="text"
                  class="form-control"
                  :class="{ 'is-invalid': createErrors.name }"
                  placeholder="Enter category name"
                  required
                  maxlength="100"
                />
                <div v-if="createErrors.name" class="invalid-feedback">
                  {{ createErrors.name }}
                </div>
              </div>

              <div class="mb-3">
                <label for="parentCategory" class="form-label">
                  Parent Category
                </label>
                <select
                  id="parentCategory"
                  v-model="createForm.parent_id"
                  class="form-select"
                >
                  <option :value="undefined">None (Top Level)</option>
                  <option
                    v-for="category in categories"
                    :key="category.id"
                    :value="category.id"
                  >
                    {{ category.name }}
                  </option>
                </select>
                <div class="form-text">
                  Select a parent category to create a subcategory
                </div>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button
                  type="button"
                  class="btn btn-secondary"
                  @click="closeCreateModal"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  class="btn btn-primary"
                  :disabled="creating || !createForm.name.trim()"
                >
                  <span v-if="creating" class="spinner-border spinner-border-sm me-2"></span>
                  <i v-else class="bi bi-check me-2"></i>
                  Create Category
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div
      class="modal fade"
      :class="{ show: showEditModal }"
      :style="{ display: showEditModal ? 'block' : 'none' }"
      tabindex="-1"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              <i class="bi bi-pencil me-2"></i>
              Edit Category
            </h5>
            <button
              type="button"
              class="btn-close"
              @click="closeEditModal"
            ></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleUpdate">
              <div class="mb-3">
                <label for="editCategoryName" class="form-label">
                  Category Name <span class="text-danger">*</span>
                </label>
                <input
                  id="editCategoryName"
                  v-model="editForm.name"
                  type="text"
                  class="form-control"
                  :class="{ 'is-invalid': editErrors.name }"
                  placeholder="Enter category name"
                  required
                  maxlength="100"
                />
                <div v-if="editErrors.name" class="invalid-feedback">
                  {{ editErrors.name }}
                </div>
              </div>

              <div class="mb-3">
                <label for="editParentCategory" class="form-label">
                  Parent Category
                </label>
                <select
                  id="editParentCategory"
                  v-model="editForm.parent_id"
                  class="form-select"
                >
                  <option :value="undefined">None (Top Level)</option>
                  <option
                    v-for="category in availableParentCategories"
                    :key="category.id"
                    :value="category.id"
                  >
                    {{ category.name }}
                  </option>
                </select>
                <div class="form-text">
                  Select a parent category to create a subcategory
                </div>
              </div>

              <div class="d-flex justify-content-between">
                <button
                  type="button"
                  class="btn btn-danger"
                  @click="confirmDelete(editingCategory!)"
                >
                  <i class="bi bi-trash me-2"></i>
                  Delete Category
                </button>
                <div class="d-flex gap-2">
                  <button
                    type="button"
                    class="btn btn-secondary"
                    @click="closeEditModal"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    class="btn btn-primary"
                    :disabled="updating || !editForm.name?.trim()"
                  >
                    <span v-if="updating" class="spinner-border spinner-border-sm me-2"></span>
                    <i v-else class="bi bi-check me-2"></i>
                    Update Category
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Backdrop -->
    <div
      v-if="showCreateModal || showEditModal"
      class="modal-backdrop fade show"
      @click="closeModals"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { BookmarkCategory, CreateBookmarkCategoryRequest, UpdateBookmarkCategoryRequest } from '@/types';
import { useBookmarkStore } from '@/stores/bookmarkStore';

const bookmarkStore = useBookmarkStore();

// Local state
const loading = ref(false);
const error = ref<string | null>(null);
const creating = ref(false);
const updating = ref(false);

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingCategory = ref<BookmarkCategory | null>(null);

// Form data
const createForm = ref<CreateBookmarkCategoryRequest>({
  name: '',
  parent_id: undefined
});

const editForm = ref<UpdateBookmarkCategoryRequest>({
  name: '',
  parent_id: undefined
});

// Form errors
const createErrors = ref<Record<string, string>>({});
const editErrors = ref<Record<string, string>>({});

// Computed properties
const categories = computed(() => bookmarkStore.categories);

const availableParentCategories = computed(() => {
  if (!editingCategory.value) return categories.value;
  // Exclude the current category and its descendants to prevent circular references
  return categories.value.filter(cat => cat.id !== editingCategory.value!.id);
});

// Methods
const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(date));
};

const getParentCategoryName = (parentId: number): string => {
  const parent = categories.value.find(cat => cat.id === parentId);
  return parent ? parent.name : 'Unknown';
};

const validateCreateForm = (): boolean => {
  createErrors.value = {};

  if (!createForm.value.name.trim()) {
    createErrors.value.name = 'Category name is required';
  } else if (createForm.value.name.length > 100) {
    createErrors.value.name = 'Category name must be less than 100 characters';
  }

  return Object.keys(createErrors.value).length === 0;
};

const validateEditForm = (): boolean => {
  editErrors.value = {};

  if (!editForm.value.name?.trim()) {
    editErrors.value.name = 'Category name is required';
  } else if (editForm.value.name && editForm.value.name.length > 100) {
    editErrors.value.name = 'Category name must be less than 100 characters';
  }

  return Object.keys(editErrors.value).length === 0;
};

const loadCategories = async () => {
  loading.value = true;
  error.value = null;

  try {
    await bookmarkStore.fetchCategories();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load categories';
  } finally {
    loading.value = false;
  }
};

const handleCreate = async () => {
  if (!validateCreateForm()) return;

  creating.value = true;
  try {
    await bookmarkStore.createCategory(createForm.value);
    closeCreateModal();
    // Categories are automatically updated in the store
  } catch (err) {
    console.error('Failed to create category:', err);
  } finally {
    creating.value = false;
  }
};

const editCategory = (category: BookmarkCategory) => {
  editingCategory.value = category;
  editForm.value = {
    name: category.name,
    parent_id: category.parent_id
  };
  editErrors.value = {};
  showEditModal.value = true;
};

const handleUpdate = async () => {
  if (!editingCategory.value || !validateEditForm()) return;

  updating.value = true;
  try {
    await bookmarkStore.updateCategory(editingCategory.value.id, editForm.value);
    closeEditModal();
    // Categories are automatically updated in the store
  } catch (err) {
    console.error('Failed to update category:', err);
  } finally {
    updating.value = false;
  }
};

const confirmDelete = (category: BookmarkCategory) => {
  if (confirm(`Are you sure you want to delete the category "${category.name}"?\n\nThis action cannot be undone.`)) {
    deleteCategory(category);
  }
};

const deleteCategory = async (category: BookmarkCategory) => {
  try {
    await bookmarkStore.deleteCategory(category.id);
    closeModals();
    // Categories are automatically updated in the store
  } catch (err) {
    console.error('Failed to delete category:', err);
    alert('Failed to delete category. It may be in use by existing bookmarks.');
  }
};

const closeCreateModal = () => {
  showCreateModal.value = false;
  createForm.value = { name: '', parent_id: undefined };
  createErrors.value = {};
};

const closeEditModal = () => {
  showEditModal.value = false;
  editingCategory.value = null;
  editForm.value = { name: '', parent_id: undefined };
  editErrors.value = {};
};

const closeModals = () => {
  closeCreateModal();
  closeEditModal();
};

// Load categories on mount
onMounted(() => {
  loadCategories();
});
</script>

<style scoped>
.category-management {
  padding: 1rem 0;
}

.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.modal.show {
  display: block !important;
}

.modal-backdrop {
  background-color: rgba(0, 0, 0, 0.5);
}

.table th {
  border-top: none;
  font-weight: 600;
  color: #495057;
  background-color: #f8f9fa;
}

.table td {
  vertical-align: middle;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
}

.invalid-feedback {
  display: block;
}

.card {
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  border: 1px solid rgba(0, 0, 0, 0.125);
}
</style>
