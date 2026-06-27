/**
 * Language storage utilities for managing user language preference via cookies
 */

const LANGUAGE_COOKIE_NAME = "language";
const COOKIE_MAX_AGE = 365 * 24 * 60 * 60; // 1 year in seconds

export type Language = "en" | "ja";

/**
 * Get the stored language preference from cookies
 */
export function getStoredLanguage(): Language | null {
  if (typeof document === "undefined") return null;

  const cookies = document.cookie.split(";");
  const languageCookie = cookies.find((cookie) =>
    cookie.trim().startsWith(`${LANGUAGE_COOKIE_NAME}=`)
  );

  if (!languageCookie) return null;

  const value = languageCookie.split("=")[1]?.trim();
  if (value === "en" || value === "ja") {
    return value;
  }

  return null;
}

/**
 * Store the language preference in cookies
 */
export function setStoredLanguage(language: Language): void {
  if (typeof document === "undefined") return;

  const secureAttribute =
    window.location.protocol === "https:" ? "; Secure" : "";

  document.cookie = `${LANGUAGE_COOKIE_NAME}=${language}; path=/; max-age=${COOKIE_MAX_AGE}; SameSite=Lax${secureAttribute}`;
}

/**
 * Detect browser's preferred language
 */
export function detectBrowserLanguage(): Language {
  if (typeof window === "undefined") return "en";
  return window.navigator.language.toLowerCase().startsWith("ja") ? "ja" : "en";
}

/**
 * Get the language to use: stored preference > browser preference > default
 */
export function getPreferredLanguage(): Language {
  return getStoredLanguage() ?? detectBrowserLanguage();
}
