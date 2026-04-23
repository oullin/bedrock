export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly path: string,
  ) {
    super(`${path} returned ${status}`);
    this.name = "ApiError";
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const response = await fetch(path, { headers: { Accept: "application/json" } });

  if (!response.ok) {
    throw new ApiError(response.status, path);
  }

  return (await response.json()) as T;
}

export async function apiPost<T = unknown>(path: string, body?: unknown): Promise<T | null> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    throw new ApiError(response.status, path);
  }

  return response.status === 204 ? null : ((await response.json()) as T);
}

export async function apiDelete(path: string): Promise<void> {
  const response = await fetch(path, { method: "DELETE" });

  if (!response.ok) {
    throw new ApiError(response.status, path);
  }
}
