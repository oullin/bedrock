export function summarizeBillingWorkspace() {
  return {
    tiers: [
      {
        name: "Starter",
        price: "$0",
        description: "Workspace bootstrap and internal packages ready for feature work.",
        featured: false,
      },
      {
        name: "Growth",
        price: "$49",
        description:
          "Adds shared foundations for the first customer-facing billing and policy flows.",
        featured: true,
      },
      {
        name: "Scale",
        price: "$149",
        description:
          "Reserves headroom for queueing, notifications, files, and media automation across the monorepo.",
        featured: false,
      },
    ],
  };
}
