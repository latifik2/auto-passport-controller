import type { CommonPassport } from "../types";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "/api";

export async function fetchPassports(): Promise<CommonPassport[]> {
  const res = await fetch(`${API_BASE_URL}/v1/passports`);

  if (!res.ok) {
    throw new Error(`HTTP ${res.status}`);
  }

  return res.json();
}
