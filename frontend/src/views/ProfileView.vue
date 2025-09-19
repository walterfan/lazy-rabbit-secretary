<template>
  <div class="profile-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <i class="bi bi-person-circle"></i>
          {{ $t('nav.profile') }}
        </h1>
        <p class="page-description">
          {{ $t('profile.description') }}
        </p>
      </div>
    </div>

    <!-- Profile Content -->
    <div class="content-wrapper">
      <div class="profile-content">
        <div class="profile-card">
          <div class="profile-header">
            <div class="profile-avatar">
              <i class="bi bi-person-circle"></i>
            </div>
            <div class="profile-info">
              <h2>{{ currentUser?.username || 'User' }}</h2>
              <p class="profile-email">{{ currentUser?.email || 'No email provided' }}</p>
            </div>
          </div>
          
          <div class="profile-details">
            <div class="detail-section">
              <h3>{{ $t('profile.accountInfo') }}</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="detail-label">{{ $t('profile.username') }}:</span>
                  <span class="detail-value">{{ currentUser?.username || '-' }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">{{ $t('profile.email') }}:</span>
                  <span class="detail-value">{{ currentUser?.email || '-' }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">{{ $t('profile.createdAt') }}:</span>
                  <span class="detail-value">{{ currentUser?.created_at ? formatDate(currentUser.created_at) : '-' }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">{{ $t('profile.lastLogin') }}:</span>
                  <span class="detail-value">{{ currentUser?.updated_at ? formatDate(currentUser.updated_at) : '-' }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h3>{{ $t('profile.preferences') }}</h3>
              <div class="preferences-grid">
                <div class="preference-item">
                  <label>{{ $t('language.switch') }}:</label>
                  <select v-model="selectedLanguage" @change="changeLanguage" class="form-select">
                    <option value="en">{{ $t('language.english') }}</option>
                    <option value="zh">{{ $t('language.chinese') }}</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/stores/authStore';
import { useLanguageStore } from '@/stores/languageStore';
import { formatDate } from '@/utils/dateUtils';

const { t, locale } = useI18n();
const authStore = useAuthStore();
const languageStore = useLanguageStore();

const selectedLanguage = ref(locale.value);

const currentUser = computed(() => authStore.currentUser);

const changeLanguage = () => {
  languageStore.setLanguage(selectedLanguage.value as 'en' | 'zh');
};

onMounted(() => {
  selectedLanguage.value = locale.value;
});
</script>

<style scoped>
/* Page Layout */
.profile-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #e9ecef;
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

/* Content Wrapper */
.content-wrapper {
  position: relative;
}

/* Profile Content */
.profile-content {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
}

.profile-card {
  max-width: 800px;
  margin: 0 auto;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-bottom: 3rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #e9ecef;
}

.profile-avatar {
  font-size: 4rem;
  color: #667eea;
}

.profile-info h2 {
  font-size: 1.75rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 0.5rem;
}

.profile-email {
  color: #6c757d;
  margin: 0;
}

.detail-section {
  margin-bottom: 2rem;
}

.detail-section h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #212529;
  margin-bottom: 1rem;
}

.detail-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f8f9fa;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #6c757d;
  min-width: 120px;
}

.detail-value {
  color: #212529;
}

.preferences-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.preference-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.preference-item label {
  font-weight: 600;
  color: #6c757d;
  min-width: 120px;
}

.form-select {
  max-width: 200px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-view {
    padding: 1rem;
  }
  
  .profile-header {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }
  
  .profile-avatar {
    font-size: 3rem;
  }
  
  .detail-item {
    flex-direction: column;
    align-items: start;
    gap: 0.25rem;
  }
  
  .preference-item {
    flex-direction: column;
    align-items: start;
    gap: 0.5rem;
  }
  
  .form-select {
    max-width: 100%;
  }
}
</style>
