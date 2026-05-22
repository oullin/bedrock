import { ref, type Ref } from "vue";

export interface Route {
  path: string;
  params: Record<string, string>;
}

function parseHash(): Route {
  const raw = window.location.hash.replace(/^#/, "") || "/";
  const [path, query = ""] = raw.split("?");
  const params = Object.fromEntries(new URLSearchParams(query));

  return { path, params };
}

export const currentRoute: Ref<Route> = ref(parseHash());

window.addEventListener("hashchange", () => {
  currentRoute.value = parseHash();
});

export function navigate(path: string): void {
  window.location.hash = path;
}
