<template>
  <div class="analytics-panel">
    <!-- Analytics Header -->
    <div class="analytics-header">
      <h3>
        <i class="bi bi-graph-up"></i>
        {{ $t('gtd.analytics') }}
      </h3>
      <p>{{ $t('gtd.analyticsDescription') }}</p>
    </div>

    <!-- Coming Soon -->
    <div class="coming-soon">
      <div class="coming-soon-icon">
        <i class="bi bi-tools"></i>
      </div>
      <h4>{{ $t('gtd.featureInDevelopment') }}</h4>
      <p>{{ $t('gtd.featureInDevelopmentDesc') }}</p>
      <div class="features-preview">
        <div class="feature-item">
          <i class="bi bi-graph-up"></i>
          <span>{{ $t('gtd.completionRateTrend') }}</span>
        </div>
        <div class="feature-item">
          <i class="bi bi-pie-chart"></i>
          <span>{{ $t('gtd.priorityDistribution') }}</span>
        </div>
        <div class="feature-item">
          <i class="bi bi-stopwatch"></i>
          <span>{{ $t('gtd.timeEfficiencyAnalysis') }}</span>
        </div>
        <div class="feature-item">
          <i class="bi bi-geo-alt"></i>
          <span>{{ $t('gtd.contextUsageStats') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
/* Analytics Panel Styles */
.analytics-panel {
  padding: 0;
}

/* Analytics Header */
.analytics-header {
  text-align: center;
  margin-bottom: 2rem;
  padding: 2rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 16px;
  color: white;
}

.analytics-header h3 {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.analytics-header p {
  opacity: 0.9;
  margin: 0;
  font-size: 1rem;
}

/* Coming Soon */
.coming-soon {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.coming-soon-icon {
  font-size: 4rem;
  color: #667eea;
  margin-bottom: 1.5rem;
  opacity: 0.8;
}

.coming-soon h4 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.5rem;
}

.coming-soon p {
  font-size: 1rem;
  color: #718096;
  margin-bottom: 2rem;
}

.features-preview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  max-width: 600px;
  margin: 0 auto;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  color: #4a5568;
  font-size: 0.9rem;
}

.feature-item i {
  color: #667eea;
  font-size: 1.1rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .analytics-header {
    padding: 1.5rem;
  }
  
  .analytics-header h3 {
    font-size: 1.5rem;
  }
  
  .coming-soon {
    padding: 3rem 1.5rem;
  }
  
  .features-preview {
    grid-template-columns: 1fr;
  }
}
</style>
