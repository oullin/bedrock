export function describeAuthWorkspace() {
  return {
    label: "Auth package linked",
    message:
      "Shared auth utilities resolve through pnpm workspaces and the Turborepo dependency graph.",
    workspace: "@gollin/auth",
  };
}
