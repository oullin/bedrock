<script setup lang="ts">
import { useForm, router } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

const props = defineProps<{
  team: { id: string; name: string; owner_id: string; personal_team: boolean };
  members: Array<{ user_id: string; role: string; name?: string; email?: string }>;
  availableRoles: Array<{ key: string; name: string; description: string }>;
  invitations: Array<{ id: string; email: string; role: string }>;
}>();

const updateForm = useForm({ name: props.team.name });
const addMemberForm = useForm({ email: "", role: "" });

function updateTeam() {
  updateForm.put(`/teams/${props.team.id}`);
}

function addMember() {
  addMemberForm.post(`/teams/${props.team.id}/members`, {
    onSuccess: () => addMemberForm.reset(),
  });
}

function updateRole(userId: string, role: string) {
  router.put(`/teams/${props.team.id}/members/${userId}`, { role });
}

function removeMember(userId: string) {
  if (confirm("Are you sure?")) {
    router.delete(`/teams/${props.team.id}/members/${userId}`);
  }
}

function cancelInvitation(invitationId: string) {
  router.delete(`/team-invitations/${invitationId}`);
}

function deleteTeam() {
  if (confirm("Are you sure you want to delete this team?")) {
    router.delete(`/teams/${props.team.id}`);
  }
}
</script>

<template>
  <AppLayout title="Team Settings">
    <template #header>
      <h2 class="text-xl font-semibold leading-tight text-gray-800">Team Settings</h2>
    </template>

    <div class="space-y-6">
      <!-- Update Team Name -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Team Name</h3>

        <form @submit.prevent="updateTeam" class="mt-4 space-y-4">
          <div>
            <label for="team-name" class="block text-sm font-medium text-gray-700">Team Name</label>
            <input id="team-name" v-model="updateForm.name" type="text"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
            <p v-if="updateForm.errors.name" class="mt-1 text-sm text-red-600">{{ updateForm.errors.name }}</p>
          </div>

          <button type="submit" :disabled="updateForm.processing"
            class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
            Save
          </button>
        </form>
      </div>

      <!-- Add Team Member -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Add Team Member</h3>
        <p class="mt-1 text-sm text-gray-600">Add a new team member by their email address.</p>

        <form @submit.prevent="addMember" class="mt-4 space-y-4">
          <div>
            <label for="member-email" class="block text-sm font-medium text-gray-700">Email</label>
            <input id="member-email" v-model="addMemberForm.email" type="email"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">Role</label>
            <div class="mt-2 space-y-2">
              <label v-for="role in availableRoles" :key="role.key" class="flex items-start gap-3 rounded-md border p-3 cursor-pointer" :class="addMemberForm.role === role.key ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200'">
                <input v-model="addMemberForm.role" type="radio" :value="role.key" class="mt-0.5" />
                <div>
                  <div class="text-sm font-medium text-gray-900">{{ role.name }}</div>
                  <div class="text-sm text-gray-500">{{ role.description }}</div>
                </div>
              </label>
            </div>
          </div>

          <button type="submit" :disabled="addMemberForm.processing"
            class="rounded-md bg-gray-800 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-gray-700 disabled:opacity-50">
            Add
          </button>
        </form>
      </div>

      <!-- Pending Invitations -->
      <div v-if="invitations.length > 0" class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Pending Invitations</h3>

        <div class="mt-4 space-y-3">
          <div v-for="invitation in invitations" :key="invitation.id" class="flex items-center justify-between">
            <div class="text-sm text-gray-600">{{ invitation.email }} &mdash; {{ invitation.role }}</div>
            <button @click="cancelInvitation(invitation.id)" class="text-sm text-red-600 hover:text-red-800">Cancel</button>
          </div>
        </div>
      </div>

      <!-- Team Members -->
      <div v-if="members.length > 0" class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Team Members</h3>

        <div class="mt-4 space-y-3">
          <div v-for="member in members" :key="member.user_id" class="flex items-center justify-between">
            <div class="text-sm text-gray-700">{{ member.name || member.email || member.user_id }}</div>

            <div class="flex items-center gap-3">
              <select @change="updateRole(member.user_id, ($event.target as HTMLSelectElement).value)" :value="member.role"
                class="rounded-md border-gray-300 text-sm shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
                <option v-for="role in availableRoles" :key="role.key" :value="role.key">{{ role.name }}</option>
              </select>

              <button @click="removeMember(member.user_id)" class="text-sm text-red-600 hover:text-red-800">Remove</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Delete Team -->
      <div v-if="!team.personal_team" class="rounded-lg bg-white p-6 shadow">
        <h3 class="text-lg font-medium text-gray-900">Delete Team</h3>
        <p class="mt-1 text-sm text-gray-600">Permanently delete this team.</p>

        <button @click="deleteTeam" class="mt-4 rounded-md bg-red-600 px-4 py-2 text-xs font-semibold uppercase text-white hover:bg-red-500">
          Delete Team
        </button>
      </div>
    </div>
  </AppLayout>
</template>
