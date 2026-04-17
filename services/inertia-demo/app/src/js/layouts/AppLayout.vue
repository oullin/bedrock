<script setup lang="ts">
import { computed, onMounted, watch, nextTick } from "vue";
import { Head, Link, usePage } from "@inertiajs/vue3";
import {
  IconDashboard,
  IconUsers,
  IconCode,
  IconFileText,
  IconDots,
  IconLogout,
} from "@tabler/icons-vue";
import { Avatar, AvatarFallback } from "@/js/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/js/components/ui/dropdown-menu";
import LoadingBar from "@/js/components/ui/loading-bar/LoadingBar.vue";
import { Separator } from "@/js/components/ui/separator";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/js/components/ui/sidebar";
import {
  appRoutes,
  contactRoutes,
  organizationRoutes,
  featureRoute,
  useDemoRoute,
} from "@/js/lib/routes";
import type { Breadcrumb, NavGroup, SharedPageProps } from "@/js/types";

const props = withDefaults(defineProps<{ title?: string; breadcrumbs?: Breadcrumb[] }>(), {
  title: "",
  breadcrumbs: () => [],
});

const page = usePage<SharedPageProps>();
const currentPath = computed(() => String(page.url ?? "/").split("?")[0]);
const title = computed(
  () => props.title || ((page.props as Record<string, unknown>).title as string) || "",
);
const breadcrumbs = computed(() =>
  props.breadcrumbs.length
    ? props.breadcrumbs
    : (((page.props as Record<string, unknown>).breadcrumbs as Breadcrumb[]) ?? []),
);
const user = computed(() => page.props.auth?.user ?? null);

function scrollToActiveNavItem() {
  nextTick(() => {
    const el = document.querySelector('[data-sidebar="content"] [data-active="true"]');
    if (el) {
      el.scrollIntoView({ block: "nearest" });
      return;
    }
    // If element isn't rendered yet, wait for next frame
    requestAnimationFrame(() => {
      document
        .querySelector('[data-sidebar="content"] [data-active="true"]')
        ?.scrollIntoView({ block: "nearest" });
    });
  });
}

onMounted(scrollToActiveNavItem);
watch(currentPath, scrollToActiveNavItem);

const groups = computed(() => [
  {
    label: "CRM",
    items: [
      { title: "Dashboard", href: appRoutes.dashboard().url, icon: IconDashboard },
      { title: "Contacts", href: contactRoutes.index().url, icon: IconUsers },
      { title: "Organizations", href: organizationRoutes.index().url, icon: IconUsers },
    ],
  },
]);

const featureGroups = computed(() => [
  {
    label: "Forms",
    items: [
      { title: "useForm", href: featureRoute("features.forms.use-form") },
      { title: "Form Component", href: featureRoute("features.forms.form-component") },
      { title: "File Uploads", href: featureRoute("features.forms.file-uploads") },
      { title: "Validation", href: featureRoute("features.forms.validation") },
      { title: "Precognition", href: featureRoute("features.forms.precognition") },
      { title: "Optimistic Updates", href: featureRoute("features.forms.optimistic-updates") },
      { title: "Form Context", href: featureRoute("features.forms.use-form-context") },
      { title: "Dotted Keys", href: featureRoute("features.forms.dotted-keys") },
      { title: "Wayfinder", href: featureRoute("features.forms.wayfinder") },
    ],
  },
  {
    label: "Navigation",
    items: [
      { title: "Links", href: featureRoute("features.navigation.links") },
      { title: "Preserve State", href: featureRoute("features.navigation.preserve-state") },
      { title: "Preserve Scroll", href: featureRoute("features.navigation.preserve-scroll") },
      { title: "View Transitions", href: featureRoute("features.navigation.view-transitions") },
      { title: "History", href: featureRoute("features.navigation.history-management") },
      { title: "Async Requests", href: featureRoute("features.navigation.async-requests") },
      { title: "Manual Visits", href: featureRoute("features.navigation.manual-visits") },
      { title: "Redirects", href: featureRoute("features.navigation.redirects") },
      { title: "Scroll Management", href: featureRoute("features.navigation.scroll-management") },
      { title: "Instant Visits", href: featureRoute("features.navigation.instant-visits") },
      { title: "URL Fragments", href: featureRoute("features.navigation.url-fragments") },
    ],
  },
  {
    label: "Data Loading",
    items: [
      { title: "Deferred Props", href: featureRoute("features.data-loading.deferred-props") },
      { title: "Partial Reloads", href: featureRoute("features.data-loading.partial-reloads") },
      { title: "Infinite Scroll", href: featureRoute("features.data-loading.infinite-scroll") },
      { title: "When Visible", href: featureRoute("features.data-loading.when-visible") },
      { title: "Polling", href: featureRoute("features.data-loading.polling") },
      { title: "Prop Merging", href: featureRoute("features.data-loading.prop-merging") },
      { title: "Optional Props", href: featureRoute("features.data-loading.optional-props") },
      {
        title: "Once Props",
        href: useDemoRoute("features.data-loading.once-props", { page: 1 }).url,
      },
    ],
  },
  {
    label: "Prefetching",
    items: [
      { title: "Link Prefetch", href: featureRoute("features.prefetching.link-prefetch") },
      {
        title: "Stale While Revalidate",
        href: featureRoute("features.prefetching.stale-while-revalidate"),
      },
      { title: "Manual Prefetch", href: featureRoute("features.prefetching.manual-prefetch") },
      { title: "Cache Management", href: featureRoute("features.prefetching.cache-management") },
    ],
  },
  {
    label: "State",
    items: [
      { title: "Remember", href: featureRoute("features.state.remember") },
      { title: "Flash Data", href: featureRoute("features.state.flash-data") },
      { title: "Shared Props", href: featureRoute("features.state.shared-props") },
    ],
  },
  {
    label: "Layouts",
    items: [
      { title: "Persistent Layouts", href: featureRoute("features.layouts.persistent-layouts") },
      { title: "Nested Layouts", href: featureRoute("features.layouts.nested-layouts") },
      { title: "Head", href: featureRoute("features.layouts.head") },
      { title: "Layout Props", href: featureRoute("features.layouts.layout-props") },
    ],
  },
  {
    label: "Events",
    items: [
      { title: "Global Events", href: featureRoute("features.events.global-events") },
      { title: "Visit Callbacks", href: featureRoute("features.events.visit-callbacks") },
      { title: "Progress", href: featureRoute("features.events.progress") },
    ],
  },
  {
    label: "Errors",
    items: [
      { title: "HTTP Errors", href: featureRoute("features.errors.http-error") },
      { title: "Network Errors", href: featureRoute("features.errors.network-errors") },
    ],
  },
  {
    label: "HTTP",
    items: [{ title: "useHttp", href: featureRoute("features.http.use-http") }],
  },
]);
</script>

<template>
  <Head :title="title" />
  <LoadingBar />
  <SidebarProvider
    :default-open="true"
    :style="{
      '--sidebar-width': 'calc(var(--spacing) * 68)',
      '--header-height': 'calc(var(--spacing) * 12 + 1px)',
    }"
  >
    <Sidebar collapsible="icon" class="h-auto border-r">
      <SidebarHeader class="border-b">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton as-child class="data-[slot=sidebar-menu-button]:p-1.5!">
              <Link :href="appRoutes.dashboard().url">
                <IconDashboard class="size-5!" />
                <span class="text-base font-semibold">Inertia Kitchen Sink</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup v-for="group in groups" :key="group.label">
          <SidebarGroupContent>
            <SidebarGroupLabel>{{ group.label }}</SidebarGroupLabel>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in group.items" :key="item.title">
                <SidebarMenuButton
                  as-child
                  :is-active="currentPath === item.href || currentPath.startsWith(item.href + '/')"
                  :tooltip="item.title"
                >
                  <Link :href="item.href" view-transition>
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup v-for="fGroup in featureGroups" :key="fGroup.label">
          <SidebarGroupContent>
            <SidebarGroupLabel>{{ fGroup.label }}</SidebarGroupLabel>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in fGroup.items" :key="item.title">
                <SidebarMenuButton
                  v-if="item.href"
                  as-child
                  :is-active="currentPath === item.href || currentPath.startsWith(item.href + '/')"
                  :tooltip="item.title"
                >
                  <Link :href="item.href" view-transition>
                    <IconCode class="size-4! opacity-60" />
                    <span>{{ item.title }}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter v-if="user" class="border-t">
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton size="lg" class="w-full cursor-pointer">
                  <Avatar class="size-8 rounded-lg">
                    <AvatarFallback class="rounded-lg">
                      {{
                        user.name
                          .split(" ")
                          .map((n: string) => n[0])
                          .join("")
                          .toUpperCase()
                          .slice(0, 2)
                      }}
                    </AvatarFallback>
                  </Avatar>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-medium">{{ user.name }}</span>
                    <span class="text-muted-foreground truncate text-xs">{{ user.email }}</span>
                  </div>
                  <IconDots class="ml-auto size-4!" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent side="top" align="start" class="w-(--reka-popper-anchor-width)">
                <DropdownMenuItem as-child>
                  <Link
                    :href="appRoutes.logout().url"
                    method="post"
                    as="button"
                    class="w-full cursor-pointer"
                  >
                    <IconLogout class="size-4!" />
                    Log out
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>

    <SidebarInset>
      <header
        class="bg-background/90 sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b"
      >
        <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
          <div class="flex flex-col">
            <h1 class="text-base font-medium">{{ title }}</h1>
            <p v-if="breadcrumbs.length" class="text-muted-foreground text-xs">
              {{ breadcrumbs.map((item) => item.title).join(" / ") }}
            </p>
          </div>
        </div>
      </header>

      <div
        v-if="page.props.flash?.message"
        class="mx-4 mt-4 rounded-lg border p-4 lg:mx-6"
        :class="{
          'border-green-200 bg-green-50 text-green-800': page.props.flash.kind === 'success',
          'border-red-200 bg-red-50 text-red-800': page.props.flash.kind === 'error',
          'border-yellow-200 bg-yellow-50 text-yellow-800': page.props.flash.kind === 'warning',
          'border-blue-200 bg-blue-50 text-blue-800':
            page.props.flash.kind === 'info' || !page.props.flash.kind,
        }"
      >
        <p class="text-sm font-medium">{{ page.props.flash.title || "Notice" }}</p>
        <p class="text-sm">{{ page.props.flash.message }}</p>
      </div>

      <div class="flex flex-1 flex-col">
        <slot />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
