export async function apiGet(path) {
  const response = await fetch(path, { headers: { Accept: "application/json" } });

  if (!response.ok) {
    throw new Error(`${path} returned ${response.status}`);
  }

  return response.json();
}

export async function apiPost(path, body) {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    throw new Error(`${path} returned ${response.status}`);
  }

  return response.status === 204 ? null : response.json();
}

export async function apiDelete(path) {
  const response = await fetch(path, { method: "DELETE" });

  if (!response.ok) {
    throw new Error(`${path} returned ${response.status}`);
  }
}
