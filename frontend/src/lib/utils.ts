import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isAlphanum(str: string): boolean {
  return /^[a-z0-9]+$/i.test(str);
}

export function currentUrl(): string {
  let currentUrl = window.location.href;
  if (currentUrl.endsWith("/")) {
    currentUrl = currentUrl.slice(0, -1);
  }
  return currentUrl;
}
