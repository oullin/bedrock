import type { Component } from "vue";

export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Organization {
  id: number;
  name: string;
  contacts_count?: number;
}

export interface Contact {
  id: number;
  first_name: string;
  last_name: string;
  name?: string;
  email: string;
  phone: string | null;
  is_favorite: boolean;
  organization_id: string | null;
  organization?: { id: number; name: string } | null;
}

export interface Note {
  id: number;
  body: string;
  created_at: string;
  user?: { name: string };
  contact: { id: number; name: string };
}

export interface SelectOption {
  value: string;
  label: string;
}

export interface CursorPaginated<T> {
  data: T[];
  next_cursor: string | null;
  prev_cursor?: string | null;
}

export interface OffsetPaginated<T> {
  data: T[];
  total: number;
  per_page: number;
  current_page: number;
  last_page: number;
}

export interface FlashMessage {
  kind?: "success" | "error" | "warning" | "info";
  title?: string;
  message?: string;
}

export interface Breadcrumb {
  title: string;
  href?: string;
}

export interface SharedPageProps {
  auth: { user: User | null };
  routes: Record<string, string>;
  flash: FlashMessage;
  [key: string]: unknown;
}

export interface ContactFormData {
  organization_id: string | null;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
}

export interface NavItem {
  title: string;
  href: string | null;
  icon?: Component;
}

export interface NavGroup {
  label: string;
  items: NavItem[];
}

export interface DemoRoute {
  url: string;
}
