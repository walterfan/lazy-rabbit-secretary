<template>
  <div class="container mt-4">
    <!-- Hero Section -->
    <div class="hero-section text-center mb-5">
      <h1 class="display-4 mb-3">🐰 Lazy Rabbit Secretary</h1>
      <p class="lead text-muted mb-4">
        Your AI-powered productivity companion to boost efficiency and never miss important tasks
      </p>
      <div class="d-flex justify-content-center gap-3">
        <router-link to="/help" class="btn btn-primary btn-lg">
          <i class="fas fa-rocket"></i> Learn More
        </router-link>
        <button class="btn btn-outline-primary btn-lg" @click="scrollToOverview">
          <i class="fas fa-play"></i> Quick Overview
        </button>
      </div>
    </div>

    <!-- Quick Overview Section -->
    <div id="overview" class="overview-section mb-5">
      <h2 class="text-center mb-4">✨ What Makes Us Special</h2>
      
      <div class="row g-4">
        <!-- Core Features Summary -->
        <div class="col-lg-3 col-md-6">
          <div class="overview-card text-center">
            <div class="overview-icon mb-3">
              <i class="fas fa-clock text-danger"></i>
            </div>
            <h6>🍅 Pomodoro Timer</h6>
            <p class="text-muted small">Focus with time-boxed work sessions</p>
          </div>
        </div>

        <div class="col-lg-3 col-md-6">
          <div class="overview-card text-center">
            <div class="overview-icon mb-3">
              <i class="fas fa-list-check text-success"></i>
            </div>
            <h6>✅ Smart Tasks</h6>
            <p class="text-muted small">AI-powered task management</p>
          </div>
        </div>

        <div class="col-lg-3 col-md-6">
          <div class="overview-card text-center">
            <div class="overview-icon mb-3">
              <i class="fas fa-bell text-warning"></i>
            </div>
            <h6>🔔 Smart Reminders</h6>
            <p class="text-muted small">Never miss important tasks</p>
          </div>
        </div>

        <div class="col-lg-3 col-md-6">
          <div class="overview-card text-center">
            <div class="overview-icon mb-3">
              <i class="fas fa-robot text-info"></i>
            </div>
            <h6>🤖 AI Assistant</h6>
            <p class="text-muted small">Your personal productivity coach</p>
          </div>
        </div>
      </div>
    </div>

    <!-- News/Updates Section -->
    <div class="alert alert-info" v-if="newsPrompt">
      <i class="fas fa-info-circle"></i> <strong>Latest Update:</strong> {{ newsPrompt }}
    </div>

    <!-- Call to Action -->
    <div class="cta-section text-center py-5 mb-4">
      <h3 class="mb-3">Ready to Transform Your Productivity?</h3>
      <p class="text-muted mb-4">Join thousands of users who have already boosted their efficiency</p>
      <div class="d-flex justify-content-center gap-3">
        <router-link to="/register" class="btn btn-success btn-lg">
          <i class="fas fa-user-plus"></i> Get Started Free
        </router-link>
        <router-link to="/help" class="btn btn-outline-secondary btn-lg">
          <i class="fas fa-question-circle"></i> Learn More
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

// News prompt from API
const newsPrompt = ref<string | null>(null);

// Scroll function for smooth navigation
const scrollToOverview = () => {
  const element = document.getElementById('overview');
  element?.scrollIntoView({ behavior: 'smooth' });
};

// API call to fetch news prompt
const fetchNewsPrompt = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/news`);
    const data = await response.json();
    newsPrompt.value = data.message || 'Welcome to Lazy Rabbit Secretary! Your AI-powered productivity companion is ready to help you achieve more with less effort.';
  } catch (error) {
    console.error('Failed to fetch news prompt:', error);
    newsPrompt.value = 'Welcome to Lazy Rabbit Secretary! Boost your productivity with AI-powered task management.';
  }
};

// Initialize component
onMounted(() => {
  fetchNewsPrompt();
});
</script>

<style scoped>
/* Hero Section */
.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem 1.5rem;
  border-radius: 15px;
  margin-bottom: 2rem;
}

.hero-section .display-4 {
  font-weight: 600;
  font-size: 2.5rem;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
}

.hero-section .lead {
  font-size: 1rem;
  opacity: 0.95;
}

/* Overview Cards */
.overview-card {
  padding: 2rem 1rem;
  border-radius: 15px;
  background: white;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  height: 100%;
}

.overview-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.overview-icon i {
  font-size: 2.5rem;
  transition: transform 0.3s ease;
}

.overview-card:hover .overview-icon i {
  transform: scale(1.1);
}

/* Call to Action */
.cta-section {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border-radius: 12px;
  margin: 1rem 0;
  padding: 1.5rem 1rem;
}

.cta-section h3 {
  font-weight: 500;
  font-size: 1.25rem;
  margin-bottom: 0.75rem;
}

.cta-section p {
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

/* Button enhancements */
.btn-lg {
  padding: 0.4rem 1.2rem;
  font-weight: 500;
  border-radius: 50px;
  transition: all 0.3s ease;
  font-size: 0.85rem;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.btn-success {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.3);
}

.btn-success:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(40, 167, 69, 0.4);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero-section {
    padding: 1.5rem 1rem;
  }
  
  .hero-section .display-4 {
    font-size: 2rem;
  }
  
  .cta-section {
    padding: 1rem 0.75rem;
  }
  
  .cta-section h3 {
    font-size: 1.1rem;
  }
  
  .cta-section p {
    font-size: 0.85rem;
  }
  
  .overview-card {
    padding: 1.5rem 1rem;
  }
  
  .btn-lg {
    padding: 0.35rem 1rem;
    font-size: 0.8rem;
  }
  
  .d-flex.gap-3 {
    flex-direction: column;
    gap: 0.75rem !important;
  }
}

/* Animation for overview cards */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.overview-card {
  animation: fadeInUp 0.6s ease-out;
}

.overview-card:nth-child(2) {
  animation-delay: 0.1s;
}

.overview-card:nth-child(3) {
  animation-delay: 0.2s;
}

.overview-card:nth-child(4) {
  animation-delay: 0.3s;
}
</style>