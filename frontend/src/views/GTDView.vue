<template>
  <div class="gtd-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-check2-square"></i>
          {{ $t('gtd.title') }}
        </h1>
        <p class="page-description">
          {{ $t('gtd.description') }}
        </p>
      </div>
      <div class="header-stats" v-if="stats">
        <div class="stat-card">
          <div class="stat-value">{{ stats.total_items }}</div>
          <div class="stat-label">{{ $t('gtd.totalItems') }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.completion_rate.toFixed(1) }}%</div>
          <div class="stat-label">{{ $t('gtd.completionRate') }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ formatTime(stats.total_actual_time) }}</div>
          <div class="stat-label">{{ $t('gtd.totalTimeSpent') }}</div>
        </div>
      </div>
    </div>

    <!-- Tab Navigation -->
    <div class="tab-navigation">
      <div class="tab-container">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          :class="['tab-button', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          <i :class="tab.icon"></i>
          <span>{{ tab.label }}</span>
          <span v-if="tab.badge" class="tab-badge">{{ tab.badge }}</span>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Inbox Tab -->
        <div v-if="activeTab === 'inbox'" class="tab-panel">
          <InboxPanel @refresh-stats="loadStats" />
        </div>

        <!-- Daily Checklist Tab -->
        <div v-if="activeTab === 'daily'" class="tab-panel">
          <DailyPanel @refresh-stats="loadStats" />
        </div>

        <!-- Analytics Tab -->
        <div v-if="activeTab === 'analytics'" class="tab-panel">
          <AnalyticsPanel />
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="loading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useDailyStore } from '@/stores/dailyStore';
import { useInboxStore } from '@/stores/inboxStore';
import type { DailyStatsResponse } from '@/types/gtd';
import InboxPanel from '@/components/gtd/InboxPanel.vue';
import DailyPanel from '@/components/gtd/DailyPanel.vue';
import AnalyticsPanel from '@/components/gtd/AnalyticsPanel.vue';

const { t } = useI18n();
const dailyStore = useDailyStore();
const inboxStore = useInboxStore();

const activeTab = ref('inbox');
const loading = ref(false);
const stats = ref<DailyStatsResponse | null>(null);

const tabs = computed(() => [
  {
    id: 'inbox',
    label: t('gtd.inbox'),
    icon: 'bi bi-inbox',
    badge: inboxStore.totalCount > 0 ? inboxStore.totalCount : null
  },
  {
    id: 'daily',
    label: t('gtd.checklist'),
    icon: 'bi bi-calendar-check',
    badge: dailyStore.totalCount > 0 ? dailyStore.totalCount : null
  },
  {
    id: 'analytics',
    label: t('gtd.analytics'),
    icon: 'bi bi-graph-up',
    badge: null
  }
]);

const loadStats = async () => {
  try {
    loading.value = true;
    stats.value = await dailyStore.getStats(dailyStore.selectedDate);
  } catch (error) {
    console.error('Failed to load stats:', error);
  } finally {
    loading.value = false;
  }
};

const formatTime = (minutes: number): string => {
  if (minutes < 60) {
    return `${minutes}m`;
  }
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
};

onMounted(() => {
  loadStats();
});
</script>

<style scoped>
/* Page Layout */
.gtd-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #e9ecef;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.page-title i {
  color: #667eea;
}

.page-description {
  color: #6c757d;
  margin: 0;
}

.header-stats {
  display: flex;
  gap: 1rem;
}

.stat-card {
  text-align: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  min-width: 100px;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 0.25rem;
}

.stat-label {
  font-size: 0.875rem;
  color: #6c757d;
}

/* Tab Navigation */
.tab-navigation {
  margin-bottom: 2rem;
}

.tab-container {
  display: flex;
  background: white;
  border-radius: 8px;
  padding: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
  gap: 0.5rem;
}

.tab-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  color: #6c757d;
  transition: all 0.2s ease;
  cursor: pointer;
}

.tab-button:hover {
  background: #f8f9fa;
  color: #495057;
}

.tab-button.active {
  background: #667eea;
  color: white;
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.2);
}

.tab-button i {
  font-size: 1rem;
}

.tab-badge {
  background: #dc3545;
  color: white;
  font-size: 0.7rem;
  padding: 0.2rem 0.4rem;
  border-radius: 10px;
  font-weight: 600;
  min-width: 18px;
  text-align: center;
}

.tab-button.active .tab-badge {
  background: rgba(255, 255, 255, 0.2);
}

/* Content Wrapper */
.content-wrapper {
  position: relative;
}

/* Tab Content */
.tab-content {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
  min-height: 600px;
}

.tab-panel {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Loading Overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.loading-overlay .spinner-border {
  width: 3rem;
  height: 3rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .gtd-view {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .page-title {
    font-size: 1.75rem;
  }
  
  .header-stats {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .tab-container {
    flex-direction: column;
  }
  
  .tab-button {
    justify-content: flex-start;
  }
  
  .tab-content {
    padding: 1rem;
  }
}

@media (max-width: 576px) {
  .page-title {
    font-size: 1.5rem;
  }
  
  .stat-card {
    min-width: 80px;
    padding: 0.75rem;
  }
  
  .stat-value {
    font-size: 1.25rem;
  }
}
</style>
