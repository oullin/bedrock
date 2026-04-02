export function describeAuthWorkspace() {
  return {
    label: "Auth package linked",
    message:
      "The dashboard app resolves @gollin/auth through pnpm workspaces and the Turborepo dependency graph.",
    workspace: "@gollin/auth",
  };
}
