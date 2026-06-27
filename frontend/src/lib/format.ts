/**
 * Formatting utilities
 */

/**
 * Truncate UUID to first 8 characters
 * @param uuid - Full UUID string
 * @returns First 8 characters of UUID
 */
export function truncateUuid(uuid: string): string {
  return uuid.slice(0, 8);
}
