<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  cancelSubscription,
  createSubscription,
  fetchBillingState,
  markPendingCheckout,
  resumeSubscription,
  updatePaymentMethod,
  updateSubscription,
} from "./api";
import {
  checkoutOpenOptions,
  handlePaddleCheckoutEvent,
  initializePaddle,
  openPaddleCheckout,
  type CheckoutMode,
} from "./checkout";
import { planPrice, statusLabel as formatStatusLabel } from "./format";
import type { BillingPlan, BillingPortalState } from "./types";

const fallbackState: BillingPortalState = {
  appName: "Billing",
  sandbox: false,
  billableId: "",
  billableName: "",
  billableType: "team",
  brandColor: "bg-gray-800",
  dashboardUrl: "/",
  defaultInterval: "monthly",
  invoices: [],
  monthlyPlans: [],
  yearlyPlans: [],
  sparkPath: "billing",
  state: "none",
  subscription: {
    status: "",
    plan_code: "",
    plan_name: "",
    pending_expires_at: null,
    payment_ready_at: null,
    portal_url: "",
    pay_now: false,
  },
  cta: {
    visible: false,
    label: "",
    remaining_days: 0,
  },
};

const state = ref<BillingPortalState>(window.__BILLING_STATE__ ?? fallbackState);
const interval = ref(state.value.defaultInterval === "yearly" ? "yearly" : "monthly");
const busyAction = ref<string | null>(null);
const checkoutMode = ref<CheckoutMode>(null);
const notice = ref(state.value.message ?? "");
const error = ref("");

const plans = computed(() =>
  interval.value === "yearly" ? state.value.yearlyPlans : state.value.monthlyPlans,
);

const activePlanId = computed(() => state.value.plan?.id ?? "");

const statusLabel = computed(() => formatStatusLabel(state.value.state));

async function refresh(message = "") {
  state.value = await fetchBillingState();
  notice.value = message;
}

async function run(action: string, callback: () => Promise<void>) {
  busyAction.value = action;
  error.value = "";
  notice.value = "";

  try {
    await callback();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Something went wrong.";
  } finally {
    busyAction.value = null;
  }
}

onMounted(() => {
  initializePaddle(state.value, (event) => {
    void handlePaddleCheckoutEvent(event, {
      getMode: () => checkoutMode.value,
      markPendingCheckout,
      refresh,
      setMode: (mode) => {
        checkoutMode.value = mode;
      },
      setPendingState: () => {
        state.value.state = "pending";
      },
    }).catch((err: unknown) => {
      error.value = err instanceof Error ? err.message : "Something went wrong.";
      checkoutMode.value = null;
    });
  });
});

function subscribeLabel(plan: BillingPlan) {
  if (state.value.state === "active" || state.value.state === "past_due") {
    return activePlanId.value === plan.id ? "Current" : "Switch";
  }

  return "Subscribe";
}

function price(plan: BillingPlan) {
  return planPrice(plan);
}

async function choosePlan(plan: BillingPlan) {
  if (activePlanId.value === plan.id) {
    return;
  }

  await run(`plan:${plan.id}`, async () => {
    if (state.value.state === "active" || state.value.state === "past_due") {
      await updateSubscription(plan.id);
      await refresh("Subscription updated.");
      return;
    }

    const checkout = await createSubscription(plan.id);
    checkoutMode.value = "subscription";
    openPaddleCheckout(checkoutOpenOptions(checkout));
  });
}

async function updatePayment() {
  await run("payment-method", async () => {
    const transaction = await updatePaymentMethod();
    checkoutMode.value = "payment-method";
    openPaddleCheckout(checkoutOpenOptions(transaction));
  });
}

async function cancelCurrent() {
  await run("cancel", async () => {
    await cancelSubscription();
    await refresh("Subscription cancellation scheduled.");
  });
}

async function resumeCurrent() {
  await run("resume", async () => {
    await resumeSubscription();
    await refresh("Subscription resumed.");
  });
}
</script>

<template>
  <main class="min-h-screen">
    <header class="border-b border-gray-200 bg-white">
      <div class="mx-auto flex max-w-7xl flex-col gap-5 px-4 py-5 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500">{{ state.appName }}</p>
            <h1 class="mt-1 text-2xl font-semibold text-gray-950">Billing</h1>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <span
              class="rounded-md border border-gray-200 bg-gray-50 px-3 py-1.5 text-sm text-gray-700"
            >
              {{ state.billableName || state.billableType }}
            </span>
            <span
              class="rounded-md border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-sm font-medium text-emerald-800"
            >
              {{ statusLabel }}
            </span>
            <a
              class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
              :href="state.dashboardUrl"
            >
              Dashboard
            </a>
          </div>
        </div>

        <div
          v-if="notice"
          class="rounded-md border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
        >
          {{ notice }}
        </div>
        <div
          v-if="error"
          class="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
        >
          {{ error }}
        </div>
      </div>
    </header>

    <div
      class="mx-auto grid max-w-7xl gap-6 px-4 py-6 sm:px-6 lg:grid-cols-[minmax(0,1fr)_22rem] lg:px-8"
    >
      <section class="space-y-6">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 class="text-lg font-semibold text-gray-950">Plans</h2>
            <p class="mt-1 text-sm text-gray-600">
              Choose the billing interval and plan for this account.
            </p>
          </div>
          <div class="inline-flex w-fit rounded-md border border-gray-300 bg-white p-1">
            <button
              class="rounded px-3 py-1.5 text-sm font-medium"
              :class="
                interval === 'monthly' ? 'bg-gray-900 text-white' : 'text-gray-700 hover:bg-gray-50'
              "
              type="button"
              @click="interval = 'monthly'"
            >
              Monthly
            </button>
            <button
              class="rounded px-3 py-1.5 text-sm font-medium"
              :class="
                interval === 'yearly' ? 'bg-gray-900 text-white' : 'text-gray-700 hover:bg-gray-50'
              "
              type="button"
              @click="interval = 'yearly'"
            >
              Yearly
            </button>
          </div>
        </div>

        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <article
            v-for="plan in plans"
            :key="plan.id"
            class="flex min-h-64 flex-col justify-between rounded-lg border bg-white p-5"
            :class="
              activePlanId === plan.id
                ? 'border-emerald-400 ring-1 ring-emerald-300'
                : 'border-gray-200'
            "
          >
            <div>
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h3 class="text-base font-semibold text-gray-950">{{ plan.name }}</h3>
                  <p class="mt-1 text-sm text-gray-600">{{ plan.short_description }}</p>
                </div>
                <span
                  v-if="activePlanId === plan.id"
                  class="rounded bg-emerald-100 px-2 py-1 text-xs font-medium text-emerald-800"
                >
                  Active
                </span>
              </div>
              <p class="mt-5 text-2xl font-semibold text-gray-950">{{ price(plan) }}</p>
              <ul class="mt-4 space-y-2 text-sm text-gray-700">
                <li v-for="feature in plan.features ?? []" :key="feature" class="flex gap-2">
                  <span class="mt-1 h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
                  <span>{{ feature }}</span>
                </li>
              </ul>
            </div>
            <button
              class="mt-6 rounded-md px-3 py-2 text-sm font-semibold"
              :class="
                activePlanId === plan.id
                  ? 'bg-gray-100 text-gray-500'
                  : 'bg-gray-900 text-white hover:bg-gray-800'
              "
              type="button"
              :disabled="busyAction === `plan:${plan.id}` || activePlanId === plan.id"
              @click="choosePlan(plan)"
            >
              {{ busyAction === `plan:${plan.id}` ? "Working" : subscribeLabel(plan) }}
            </button>
          </article>
        </div>

        <section class="rounded-lg border border-gray-200 bg-white">
          <div class="border-b border-gray-200 px-5 py-4">
            <h2 class="text-base font-semibold text-gray-950">Invoices</h2>
          </div>
          <div v-if="state.invoices.length === 0" class="px-5 py-8 text-sm text-gray-600">
            No invoices yet.
          </div>
          <div v-else class="divide-y divide-gray-200">
            <a
              v-for="invoice in state.invoices"
              :key="invoice.id"
              class="flex items-center justify-between gap-4 px-5 py-4 text-sm hover:bg-gray-50"
              :href="invoice.invoice_url"
            >
              <span class="font-medium text-gray-900">{{ invoice.billed_at || invoice.id }}</span>
              <span class="text-gray-600">{{ invoice.total }}</span>
            </a>
          </div>
        </section>
      </section>

      <aside class="space-y-4">
        <section class="rounded-lg border border-gray-200 bg-white p-5">
          <h2 class="text-base font-semibold text-gray-950">Subscription</h2>
          <dl class="mt-4 space-y-3 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-gray-500">Status</dt>
              <dd class="font-medium text-gray-900">{{ statusLabel }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-gray-500">Plan</dt>
              <dd class="font-medium text-gray-900">
                {{ state.subscription.plan_name || "None" }}
              </dd>
            </div>
            <div v-if="state.genericTrialEndsAt" class="flex justify-between gap-4">
              <dt class="text-gray-500">Trial ends</dt>
              <dd class="font-medium text-gray-900">{{ state.genericTrialEndsAt }}</dd>
            </div>
            <div v-if="state.nextPayment" class="flex justify-between gap-4">
              <dt class="text-gray-500">Next payment</dt>
              <dd class="font-medium text-gray-900">{{ state.nextPayment.amount }}</dd>
            </div>
          </dl>
        </section>

        <section class="rounded-lg border border-gray-200 bg-white p-5">
          <h2 class="text-base font-semibold text-gray-950">Actions</h2>
          <div class="mt-4 grid gap-2">
            <button
              class="rounded-md border border-gray-300 px-3 py-2 text-sm font-semibold text-gray-800 hover:bg-gray-50 disabled:opacity-50"
              type="button"
              :disabled="busyAction === 'payment-method' || state.state === 'none'"
              @click="updatePayment"
            >
              {{ busyAction === "payment-method" ? "Working" : "Update payment method" }}
            </button>
            <button
              v-if="state.state === 'active' || state.state === 'past_due'"
              class="rounded-md border border-red-200 px-3 py-2 text-sm font-semibold text-red-700 hover:bg-red-50 disabled:opacity-50"
              type="button"
              :disabled="busyAction === 'cancel'"
              @click="cancelCurrent"
            >
              {{ busyAction === "cancel" ? "Working" : "Cancel subscription" }}
            </button>
            <button
              v-if="state.state === 'onGracePeriod'"
              class="rounded-md border border-emerald-200 px-3 py-2 text-sm font-semibold text-emerald-800 hover:bg-emerald-50 disabled:opacity-50"
              type="button"
              :disabled="busyAction === 'resume'"
              @click="resumeCurrent"
            >
              {{ busyAction === "resume" ? "Working" : "Resume subscription" }}
            </button>
          </div>
        </section>
      </aside>
    </div>
  </main>
</template>
