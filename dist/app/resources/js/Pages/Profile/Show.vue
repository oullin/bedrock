<script setup lang="ts">
import { useForm, router } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

defineProps<{
  sessions?: Array<{ id: string; ip_address: string; user_agent: string; last_active: string; is_current: boolean }>;
  twoFactorEnabled?: boolean;
  qrCode?: string;
  recoveryCodes?: string[];
}>();

const profileForm = useForm({ name: "", email: "" });
const passwordForm = useForm({ current_password: "", password: "", password_confirmation: "" });

function updateProfile() {
  profileForm.put("/user/profile-information");
}

function updatePassword() {
  passwordForm.put("/user/password", {
    onSuccess: () => passwordForm.reset(),
  });
}

function enableTwoFactor() {
  router.post("/user/two-factor-authentication");
}

function disableTwoFactor() {
  router.delete("/user/two-factor-authentication");
}

function deleteAccount() {
  if (confirm("Are you sure you want to delete your account?")) {
    router.delete("/user");
  }
}

function logoutOtherSessions() {
  router.delete("/user/other-sessions");
}
</script>

<template>
  <AppLayout title="Profile">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">Profile</h2>
    </template>

    <div class="space-y-6">
      <!-- Update Profile Information -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Profile Information</h3>
        <p class="mt-1 text-sm text-gray-600">Update your account's profile information and email address.</p>

        <form @submit.prevent="updateProfile" class="mt-6 space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
            <input id="name" v-model="profileForm.name" type="text"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
            <p v-if="profileForm.errors.name" class="mt-1 text-sm text-red-600">{{ profileForm.errors.name }}</p>
          </div>

          <div>
            <label for="profile-email" class="block text-sm font-medium text-gray-700">Email</label>
            <input id="profile-email" v-model="profileForm.email" type="email"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
            <p v-if="profileForm.errors.email" class="mt-1 text-sm text-red-600">{{ profileForm.errors.email }}</p>
          </div>

          <button type="submit" :disabled="profileForm.processing"
            class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
            Save
          </button>
        </form>
      </div>

      <!-- Update Password -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Update Password</h3>
        <p class="mt-1 text-sm text-gray-600">Ensure your account is using a long, random password to stay secure.</p>

        <form @submit.prevent="updatePassword" class="mt-6 space-y-4">
          <div>
            <label for="current_password" class="block text-sm font-medium text-gray-700">Current Password</label>
            <input id="current_password" v-model="passwordForm.current_password" type="password"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
          </div>

          <div>
            <label for="new_password" class="block text-sm font-medium text-gray-700">New Password</label>
            <input id="new_password" v-model="passwordForm.password" type="password"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
          </div>

          <div>
            <label for="password_confirmation" class="block text-sm font-medium text-gray-700">Confirm Password</label>
            <input id="password_confirmation" v-model="passwordForm.password_confirmation" type="password"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
          </div>

          <button type="submit" :disabled="passwordForm.processing"
            class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
            Save
          </button>
        </form>
      </div>

      <!-- Two Factor Authentication -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Two Factor Authentication</h3>
        <p class="mt-1 text-sm text-gray-600">Add additional security to your account using two factor authentication.</p>

        <div class="mt-4">
          <template v-if="twoFactorEnabled">
            <p class="text-sm font-medium text-green-600">You have enabled two factor authentication.</p>

            <div v-if="qrCode" class="mt-4" v-html="qrCode" />

            <div v-if="recoveryCodes?.length" class="mt-4 rounded-md bg-gray-100 p-4">
              <p class="mb-2 text-sm font-medium text-gray-700">Recovery Codes</p>
              <div class="grid gap-1 font-mono text-sm">
                <div v-for="code in recoveryCodes" :key="code">{{ code }}</div>
              </div>
            </div>

            <button @click="disableTwoFactor" class="mt-4 rounded-md bg-red-600 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-red-500">
              Disable
            </button>
          </template>

          <template v-else>
            <p class="text-sm text-gray-600">You have not enabled two factor authentication.</p>
            <button @click="enableTwoFactor" class="mt-4 rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700">
              Enable
            </button>
          </template>
        </div>
      </div>

      <!-- Browser Sessions -->
      <div v-if="sessions" class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Browser Sessions</h3>
        <p class="mt-1 text-sm text-gray-600">Manage and log out your active sessions on other browsers and devices.</p>

        <div class="mt-4 space-y-3">
          <div v-for="session in sessions" :key="session.id" class="flex items-center justify-between text-sm">
            <div>
              <span class="text-gray-700">{{ session.ip_address }}</span>
              <span class="text-gray-500"> &mdash; {{ session.user_agent }}</span>
              <span v-if="session.is_current" class="ml-2 text-xs font-semibold text-green-500">This device</span>
            </div>
          </div>
        </div>

        <button @click="logoutOtherSessions" class="mt-4 rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700">
          Log Out Other Browser Sessions
        </button>
      </div>

      <!-- Delete Account -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Delete Account</h3>
        <p class="mt-1 text-sm text-gray-600">Permanently delete your account.</p>

        <button @click="deleteAccount" class="mt-4 rounded-md bg-red-600 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-red-500">
          Delete Account
        </button>
      </div>
    </div>
  </AppLayout>
</template>
