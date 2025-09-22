<template>
  <nav class="wiki-navigation">
    <ol class="breadcrumb">
      <li class="breadcrumb-item">
        <router-link to="/wiki">
          <i class="bi bi-house"></i>
          {{ $t('wiki.home') }}
        </router-link>
      </li>
      
      <!-- Category Breadcrumb -->
      <li v-if="page.categories.length > 0" class="breadcrumb-item">
        <span>{{ page.categories[0] }}</span>
      </li>
      
      <!-- Current Page -->
      <li class="breadcrumb-item active" aria-current="page">
        {{ page.title }}
      </li>
    </ol>
  </nav>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { WikiPage } from '@/types'

interface Props {
  page: WikiPage
}

defineProps<Props>()

const { t } = useI18n()
</script>

<style scoped>
.wiki-navigation {
  margin-bottom: 2rem;
}

.breadcrumb {
  background: none;
  padding: 0;
  margin: 0;
}

.breadcrumb-item {
  font-size: 0.9rem;
}

.breadcrumb-item a {
  color: #667eea;
  text-decoration: none;
  transition: color 0.2s ease;
}

.breadcrumb-item a:hover {
  color: #5a6fd8;
}

.breadcrumb-item.active {
  color: #6c757d;
  font-weight: 500;
}

.breadcrumb-item i {
  margin-right: 0.25rem;
}
</style>
