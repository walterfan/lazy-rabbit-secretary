<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container">
      <router-link class="navbar-brand" to="/">Lazy Rabbit Secretary</router-link>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNav"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <router-link class="nav-link" to="/home">
              <i class="bi bi-house"></i> {{ $t('nav.home') }}
            </router-link>
          </li>
          <li class="nav-item dropdown" v-if="isAuthenticated">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-check2-square"></i> {{ $t('nav.getThingsDone') }}
            </a>
            <ul class="dropdown-menu">

              <li>
                <router-link class="dropdown-item" to="/gtd">
                  <i class="bi bi-check2-square"></i> {{ $t('nav.gtdSystem') }}
                </router-link>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <router-link class="dropdown-item" to="/tasks">
                  <i class="bi bi-list-task"></i> {{ $t('nav.tasks') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/reminders">
                  <i class="bi bi-bell"></i> {{ $t('nav.reminders') }}
                </router-link>
              </li>

            </ul>
          </li>
          <li class="nav-item" v-if="isAuthenticated">
            <router-link class="nav-link" to="/books">
              <i class="bi bi-book"></i> {{ $t('nav.books') }}
            </router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/blog">
              <i class="bi bi-journal-text"></i> {{ $t('nav.blog') }}
            </router-link>
          </li>
          <li class="nav-item dropdown" v-if="isAuthenticated">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-tools"></i> {{ $t('nav.tools') }}
            </a>
            <ul class="dropdown-menu">
              <li>
                <router-link class="dropdown-item" to="/prompts">
                  <i class="bi bi-chat-dots"></i> {{ $t('nav.prompts') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/commands">
                  <i class="bi bi-terminal"></i> {{ $t('nav.commands') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/secrets">
                  <i class="bi bi-shield-lock"></i> {{ $t('nav.secrets') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/tools/encoding">
                  <i class="bi bi-code-slash"></i> {{ $t('nav.encodingTools') }}
                </router-link>
              </li>

            </ul>
          </li>
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-translate"></i> {{ $t('language.switch') }}
            </a>
            <ul class="dropdown-menu">
              <li>
                <a class="dropdown-item" href="#" @click.prevent="setLanguage('en')">
                  <i class="bi bi-flag-us"></i> {{ $t('language.english') }}
                </a>
              </li>
              <li>
                <a class="dropdown-item" href="#" @click.prevent="setLanguage('zh')">
                  <i class="bi bi-flag-cn"></i> {{ $t('language.chinese') }}
                </a>
              </li>
            </ul>
          </li>
          <li class="nav-item dropdown" v-if="isAuthenticated">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-gear"></i> {{ $t('nav.admin') }}
            </a>
            <ul class="dropdown-menu">
              <li>
                <router-link class="dropdown-item" to="/posts">
                  <i class="bi bi-file-earmark-text"></i> {{ $t('nav.postsManagement') }}
                </router-link>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <router-link class="dropdown-item" to="/admin/users">
                  <i class="bi bi-people"></i> {{ $t('nav.userManagement') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/admin/roles">
                  <i class="bi bi-shield"></i> {{ $t('nav.roleManagement') }}
                </router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/admin/permissions">
                  <i class="bi bi-key"></i> {{ $t('nav.permissionManagement') }}
                </router-link>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <router-link class="dropdown-item" to="/admin/registrations">
                  <i class="bi bi-person-check"></i> {{ $t('nav.registrationManagement') }}
                </router-link>
              </li>
            </ul>
          </li>
        </ul>
        <ul class="navbar-nav">
          <li class="nav-item">
            <router-link class="nav-link" to="/help">
              <i class="bi bi-question-circle"></i> {{ $t('nav.help') }}
            </router-link>
          </li>
          <li class="nav-item" v-if="!isAuthenticated">
            <router-link class="nav-link" to="/signin">
              <i class="bi bi-box-arrow-in-right"></i> {{ $t('nav.signIn') }}
            </router-link>
          </li>
          <li class="nav-item" v-if="!isAuthenticated">
            <router-link class="nav-link" to="/signup">
              <i class="bi bi-person-plus"></i> {{ $t('nav.signUp') }}
            </router-link>
          </li>
          <li class="nav-item dropdown" v-if="isAuthenticated">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-person-circle"></i> {{ currentUser?.username || 'User' }}
            </a>
            <ul class="dropdown-menu">
              <li>
                <router-link class="dropdown-item" to="/profile">
                  <i class="bi bi-person"></i> {{ $t('nav.profile') }}
                </router-link>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <a class="dropdown-item" href="#" @click.prevent="signOut">
                  <i class="bi bi-box-arrow-right"></i> {{ $t('nav.signOut') }}
                </a>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>

  <router-view></router-view>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '@/stores/authStore';
import { useLanguageStore } from '@/stores/languageStore';

const authStore = useAuthStore();
const languageStore = useLanguageStore();

const isAuthenticated = computed(() => authStore.isAuthenticated);
const currentUser = computed(() => authStore.currentUser);

const signOut = () => {
  authStore.signOut();
};

const setLanguage = (lang: 'en' | 'zh') => {
  languageStore.setLanguage(lang);
};
</script>

<style>
@import 'bootstrap/dist/css/bootstrap.min.css';
@import 'bootstrap-icons/font/bootstrap-icons.css';

.navbar-brand {
  font-weight: bold;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}
</style>