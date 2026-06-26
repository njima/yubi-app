"use client";

import { Fragment, useEffect } from "react";
import { useTranslation } from "react-i18next";

import {
  type Language,
  setStoredLanguage,
} from "@/shared/lib/language-storage";

type LanguageOption = {
  code: Language;
  label: string;
  ariaLabel: string;
};

function LanguageOptionButton({
  isActive,
  label,
  ariaLabel,
  onClick,
}: {
  isActive: boolean;
  label: string;
  ariaLabel: string;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={
        isActive
          ? "px-2 py-1 transition-colors font-semibold text-blue-600 dark:text-blue-400 underline"
          : "px-2 py-1 transition-colors text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100"
      }
      aria-label={ariaLabel}
    >
      {label}
    </button>
  );
}

const languageOptions: LanguageOption[] = [
  { code: "en", label: "EN", ariaLabel: "Switch to English" },
  { code: "ja", label: "JP", ariaLabel: "日本語に切り替え" },
];

export function LanguageSwitcher() {
  const { i18n } = useTranslation();
  const currentLanguage = i18n.language as Language;

  useEffect(() => {
    document.documentElement.lang = currentLanguage;
  }, [currentLanguage]);

  const handleLanguageChange = (lang: Language) => {
    if (lang !== currentLanguage) {
      setStoredLanguage(lang);
      i18n.changeLanguage(lang);
    }
  };

  return (
    <div className="flex items-center gap-1 text-sm">
      {languageOptions.map((option, index) => (
        <Fragment key={option.code}>
          <LanguageOptionButton
            isActive={currentLanguage === option.code}
            label={option.label}
            ariaLabel={option.ariaLabel}
            onClick={() => handleLanguageChange(option.code)}
          />
          {index < languageOptions.length - 1 ? (
            <span className="text-gray-400 dark:text-gray-600">|</span>
          ) : null}
        </Fragment>
      ))}
    </div>
  );
}
