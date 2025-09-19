import { createI18n } from 'vue-i18n';

// Import locale files
import en from '../locales/en.json';
import zh from '../locales/zh.json';

const messages = {
  en,
  zh
};

// Get saved language from localStorage or default to 'en'
const savedLanguage = localStorage.getItem('app-language') || 'en';

export const i18n = createI18n({
  legacy: false, // Use Composition API mode
  locale: savedLanguage,
  fallbackLocale: 'en',
  messages
});

export default i18n;
