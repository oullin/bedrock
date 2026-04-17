import { usePage } from "@inertiajs/vue3";
import type { DemoRoute, SharedPageProps } from "@/js/types";

function fillPattern(pattern: string, params: Record<string, string | number> = {}): string {
  let url = pattern ?? "#!wayfinder:unknown-route";

  Object.entries(params).forEach(([key, value]) => {
    url = url.replaceAll(`{${key}}`, encodeURIComponent(String(value)));
  });

  return url;
}

/**
 * Resolve a named route using the shared Inertia page props.
 * Must be called within Vue component setup or template context.
 */
export function useDemoRoute(
  name: string,
  params: Record<string, string | number> = {},
): DemoRoute {
  const page = usePage<SharedPageProps>();
  const pattern = page.props.routes?.[name];

  if (!pattern) {
    console.warn(`[wayfinder] unknown route "${name}", returning fallback`);
  }

  return {
    url: fillPattern(pattern, params),
  };
}

/**
 * Resolve a feature route pattern by name.
 * Must be called within Vue component setup or template context.
 */
export function featureRoute(name: string): string | null {
  const page = usePage<SharedPageProps>();
  const resolved = page.props.routes?.[name];

  if (!resolved) {
    console.warn(`[wayfinder] unknown feature route "${name}"`);

    return null;
  }

  return resolved;
}

/**
 * Route builders for app-level routes.
 * Must be called within Vue component setup or template context.
 */
export const appRoutes = {
  login: (): DemoRoute => useDemoRoute("login"),
  logout: (): DemoRoute => useDemoRoute("logout"),
  dashboard: (): DemoRoute => useDemoRoute("dashboard"),
};

/**
 * Route builders for contact routes.
 * Must be called within Vue component setup or template context.
 */
export const contactRoutes = {
  index: (): DemoRoute => useDemoRoute("contacts.index"),
  create: (): DemoRoute => useDemoRoute("contacts.create"),
  store: (): DemoRoute => useDemoRoute("contacts.store"),
  show: (contact: string | number): DemoRoute => useDemoRoute("contacts.show", { contact }),
  edit: (contact: string | number): DemoRoute => useDemoRoute("contacts.edit", { contact }),
  update: (contact: string | number): DemoRoute => useDemoRoute("contacts.update", { contact }),
  destroy: (contact: string | number): DemoRoute => useDemoRoute("contacts.destroy", { contact }),
  favorite: (contact: string | number): DemoRoute => useDemoRoute("contacts.favorite", { contact }),
  storeNote: (contact: string | number): DemoRoute =>
    useDemoRoute("contacts.notes.store", { contact }),
};

/**
 * Route builders for organization routes.
 * Must be called within Vue component setup or template context.
 */
export const organizationRoutes = {
  index: (): DemoRoute => useDemoRoute("organizations.index"),
  show: (organization: string | number): DemoRoute =>
    useDemoRoute("organizations.show", { organization }),
  update: (organization: string | number): DemoRoute =>
    useDemoRoute("organizations.update", { organization }),
};
