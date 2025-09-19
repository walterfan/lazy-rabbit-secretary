import { defineStore } from 'pinia';
import { ref } from 'vue';
import { i18n } from '@/i18n';

export const useLanguageStore = defineStore('language', () => {
  const currentLanguage = ref(i18n.global.locale.value);

  const setLanguage = (lang: 'en' | 'zh') => {
    currentLanguage.value = lang;
    i18n.global.locale.value = lang;
    localStorage.setItem('app-language', lang);
    
    // Update document language attribute
    document.documentElement.lang = lang;
    
    // Update document title based on language
    if (lang === 'zh') {
      document.title = '懒兔秘书 - 高效工作助手';
    } else {
      document.title = 'Lazy Rabbit Secretary - Productivity Assistant';
    }
  };

  const toggleLanguage = () => {
    const newLang = currentLanguage.value === 'en' ? 'zh' : 'en';
    setLanguage(newLang);
  };

  const getCurrentLanguage = () => currentLanguage.value;

  return {
    currentLanguage,
    setLanguage,
    toggleLanguage,
    getCurrentLanguage
  };
});
