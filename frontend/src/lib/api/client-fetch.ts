import type { z } from "zod";

export async function fetchJson<T>(
  path: string,
  errorLabel: string,
  init?: RequestInit
): Promise<T> {
  const response = await fetch(path, init);
  if (!response.ok) {
    throw new Error(`${errorLabel}: ${response.statusText}`);
  }
  return response.json();
}

export async function fetchAndParse<T>(
  path: string,
  schema: z.ZodType<T>,
  errorLabel: string,
  init?: RequestInit
): Promise<T> {
  const data = await fetchJson<unknown>(path, errorLabel, init);
  return schema.parse(data);
}
