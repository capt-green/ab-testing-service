<template>
  <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
    <div class="space-y-6 sm:px-6 lg:col-span-9 lg:px-0">
      <!-- Profile Section -->
      <section aria-labelledby="profile-section">
        <div class="shadow sm:overflow-hidden sm:rounded-md">
          <div class="bg-white py-6 px-4 sm:p-6">
            <div>
              <h2 id="profile-section" class="text-lg font-medium leading-6 text-gray-900">Profile</h2>
              <p class="mt-1 text-sm text-gray-500">
                Update your profile information.
              </p>
            </div>

            <div class="mt-6 grid grid-cols-4 gap-6">
              <div class="col-span-4 sm:col-span-2">
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input
                  type="email"
                  name="email"
                  id="email"
                  v-model="profile.email"
                  disabled
                  class="mt-1 block w-full rounded-md border-gray-300 bg-gray-100 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Password Section -->
      <section aria-labelledby="password-section">
        <form @submit.prevent="updatePassword" class="shadow sm:overflow-hidden sm:rounded-md">
          <div class="bg-white py-6 px-4 sm:p-6">
            <div>
              <h2 id="password-section" class="text-lg font-medium leading-6 text-gray-900">Password</h2>
              <p class="mt-1 text-sm text-gray-500">
                Update your password.
              </p>
            </div>

            <div class="mt-6 grid grid-cols-4 gap-6">
              <div class="col-span-4 sm:col-span-2">
                <label for="current-password" class="block text-sm font-medium text-gray-700">Current Password</label>
                <input
                  type="password"
                  name="current-password"
                  id="current-password"
                  v-model="passwordForm.currentPassword"
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div class="col-span-4 sm:col-span-2">
                <label for="new-password" class="block text-sm font-medium text-gray-700">New Password</label>
                <input
                  type="password"
                  name="new-password"
                  id="new-password"
                  v-model="passwordForm.newPassword"
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div class="col-span-4 sm:col-span-2">
                <label for="confirm-password" class="block text-sm font-medium text-gray-700">Confirm New Password</label>
                <input
                  type="password"
                  name="confirm-password"
                  id="confirm-password"
                  v-model="passwordForm.confirmPassword"
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
            <button
              type="submit"
              :disabled="loading"
              class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            >
              {{ loading ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </section>

      <!-- API Tokens Section -->
      <section aria-labelledby="tokens-section">
        <div class="shadow sm:overflow-hidden sm:rounded-md">
          <div class="bg-white py-6 px-4 sm:p-6">
            <div>
              <h2 id="tokens-section" class="text-lg font-medium leading-6 text-gray-900">API Tokens</h2>
              <p class="mt-1 text-sm text-gray-500">
                Manage your API tokens for programmatic access.
              </p>
            </div>

            <div class="mt-6">
              <button
                type="button"
                @click="generateToken"
                :disabled="loading"
                class="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                Generate New Token
              </button>

              <div class="mt-4">
                <div v-if="tokens.length === 0" class="text-sm text-gray-500">
                  No tokens generated yet.
                </div>
                <ul v-else role="list" class="divide-y divide-gray-200">
                  <li v-for="token in tokens" :key="token.id" class="flex items-center justify-between py-4">
                    <div>
                      <p class="text-sm font-medium text-gray-900">{{ token.name }}</p>
                      <p class="text-sm text-gray-500">Created: {{ token.created_at }}</p>
                    </div>
                    <button
                      type="button"
                      @click="revokeToken(token.id)"
                      class="inline-flex items-center rounded-md border border-transparent bg-red-100 px-3 py-2 text-sm font-medium leading-4 text-red-700 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                    >
                      Revoke
                    </button>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const authStore = useAuthStore()
const loading = ref(false)
const tokens = ref([])

const profile = ref({
  email: authStore.user?.email || ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

async function updatePassword() {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('New passwords do not match')
    return
  }

  loading.value = true
  try {
    await axios.post('/api/auth/change-password', {
      current_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword
    })
    
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    alert('Password updated successfully')
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to update password')
  } finally {
    loading.value = false
  }
}

async function loadTokens() {
  try {
    const response = await axios.get('/api/auth/tokens')
    tokens.value = response.data
  } catch (error) {
    console.error('Failed to load tokens:', error)
  }
}

async function generateToken() {
  loading.value = true
  try {
    const response = await axios.post('/api/auth/tokens')
    alert(`Your new token is: ${response.data.token}\nPlease save it now as you won't be able to see it again.`)
    await loadTokens()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to generate token')
  } finally {
    loading.value = false
  }
}

async function revokeToken(tokenId) {
  if (!confirm('Are you sure you want to revoke this token?')) return

  loading.value = true
  try {
    await axios.delete(`/api/auth/tokens/${tokenId}`)
    await loadTokens()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to revoke token')
  } finally {
    loading.value = false
  }
}

onMounted(loadTokens)
</script>
