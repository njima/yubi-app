import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import enTranslation from "./locales/en/translation.json";
import jaTranslation from "./locales/ja/translation.json";

// Translation resources loaded from external files
const resources = {
  en: {
    translation: enTranslation,
  },
  ja: {
    translation: jaTranslation,
  },
};

i18n.use(initReactI18next).init({
  resources,
  lng: "en",
  supportedLngs: ["en", "ja"],
  fallbackLng: "en",
  interpolation: {
    escapeValue: false, // React already escapes values
  },
});

export default i18n;
