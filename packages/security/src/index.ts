export function describeSecurityWorkspace() {
  return {
    label: "Security package linked",
    message:
      "The security workspace is available for shared policy, validation, and hardening utilities across the monorepo.",
    workspace: "@gollin/security",
  };
}
