import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isAlphanum(str: string): boolean {
  return /^[a-z0-9]+$/i.test(str);
}
