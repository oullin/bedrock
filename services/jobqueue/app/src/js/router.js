import { ref } from "vue";

export const currentRoute = ref(parseHash());

function parseHash() {
  const raw = window.location.hash.replace(/^#/, "") || "/";
  const [path, query = ""] = raw.split("?");
  const params = Object.fromEntries(new URLSearchParams(query));

  return { path, params };
}

window.addEventListener("hashchange", () => {
  currentRoute.value = parseHash();
});

export function navigate(path) {
  window.location.hash = path;
}
