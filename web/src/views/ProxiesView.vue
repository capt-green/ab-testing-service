<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-base font-semibold leading-6 text-gray-900">Proxies</h1>
        <p class="mt-2 text-sm text-gray-700">A list of all proxies in your A/B testing system.</p>
      </div>
      <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
        <button
          @click="openCreateModal"
          class="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        >
          Add proxy
        </button>
      </div>
    </div>

    <!-- Filter by tags -->
    <div class="mt-4">
      <TagsInput
        v-model="selectedTags"
        label="Filter by tags"
        :available-tags="availableTags"
        @update:modelValue="filterProxies"
      />
    </div>

    <!-- Proxies List -->
    <div class="mt-8 flow-root">
      <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
          <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-300">
              <thead class="bg-gray-50">
                <tr>
                  <th 
                    scope="col" 
                    class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 cursor-pointer"
                    @click="handleSort('id')"
                  >
                    ID
                    <span v-if="sortBy === 'id'" class="ml-1">
                      {{ sortDesc ? '↓' : '↑' }}
                    </span>
                  </th>
                  <th 
                    scope="col" 
                    class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 cursor-pointer"
                    @click="handleSort('listen_url')"
                  >
                    Listen URL
                    <span v-if="sortBy === 'listen_url'" class="ml-1">
                      {{ sortDesc ? '↓' : '↑' }}
                    </span>
                  </th>
                  <th 
                    scope="col" 
                    class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 cursor-pointer"
                    @click="handleSort('mode')"
                  >
                    Mode
                    <span v-if="sortBy === 'mode'" class="ml-1">
                      {{ sortDesc ? '↓' : '↑' }}
                    </span>
                  </th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Tags</th>
                  <th 
                    scope="col" 
                    class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 cursor-pointer"
                    @click="handleSort('targets')"
                  >
                    Targets
                    <span v-if="sortBy === 'targets'" class="ml-1">
                      {{ sortDesc ? '↓' : '↑' }}
                    </span>
                  </th>
                  <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                    <span class="sr-only">Actions</span>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 bg-white">
                <tr v-for="proxy in filteredProxies" :key="proxy.id">
                  <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
                    {{ `${proxy.id.slice(0, 4)}...${proxy.id.slice(-4)}` }}
                  </td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ proxy.listen_url }}</td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ proxy.mode }}</td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    <div class="flex flex-wrap gap-1">
                      <span
                        v-for="tag in proxy.tags"
                        :key="tag"
                        class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800"
                      >
                        {{ tag }}
                      </span>
                    </div>
                  </td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    {{ proxy.targets?.length || 0 }} targets
                  </td>
                  <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                    <button
                      @click="editProxy(proxy)"
                      class="text-indigo-600 hover:text-indigo-900 mr-4"
                    >
                      Edit
                    </button>
                    <button
                      @click="viewHistory(proxy)"
                      class="text-indigo-600 hover:text-indigo-900 mr-4"
                    >
                      History
                    </button>
                    <button
                      @click="deleteProxy(proxy.id)"
                      class="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <Pagination :currentPage="currentPage" :itemsPerPage="itemsPerPage" @changePage="changePage" :totalProxies="totalProxies"/>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="relative z-10">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
            <form @submit.prevent="handleSubmit">
              <div>
                <h3 class="text-base font-semibold leading-6 text-gray-900">
                  {{ editingProxy ? 'Edit Proxy' : 'Create New Proxy' }}
                </h3>
                <div class="mt-4 space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Listen URL</label>
                    <input
                      v-model="form.listen_url"
                      type="text"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      required
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Mode</label>
                    <select
                      v-model="form.mode"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    >
                      <option value="reverse">Reverse Proxy</option>
                      <option value="redirect">Redirect</option>
                    </select>
                  </div>

                  <!-- Route Condition -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Route Condition</label>
                    <div class="mt-2 space-y-4">
                      <div>
                        <select
                          v-model="form.condition.type"
                          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        >
                          <option value="">No condition</option>
                          <option value="header">Header</option>
                          <option value="query">Query Parameter</option>
                          <option value="cookie">Cookie</option>
                          <option value="user_agent">User Agent</option>
                          <option value="language">Language</option>
                        </select>
                      </div>

                      <!-- Parameter Name -->
                      <div v-if="form.condition.type && form.condition.type !== 'language'">
                        <label class="block text-sm font-medium text-gray-700">
                          {{ form.condition.type === 'user_agent' ? 'User Agent Parameter' : 'Parameter Name' }}
                        </label>
                        <select
                          v-if="form.condition.type === 'user_agent'"
                          v-model="form.condition.param_name"
                          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        >
                          <option value="platform">Platform (Mobile/Desktop)</option>
                          <option value="browser">Browser</option>
                        </select>
                        <input
                          v-else
                          v-model="form.condition.param_name"
                          type="text"
                          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                          :placeholder="getParamPlaceholder(form.condition.type)"
                        />
                      </div>

                      <!-- Values -->
                      <div v-if="form.condition.type">
                        <label class="block text-sm font-medium text-gray-700">Values</label>
                        <div class="mt-2 space-y-2">
                          <div
                            v-for="(target, index) in form.targets"
                            :key="target.id"
                            class="flex items-center gap-2"
                          >
                            <input
                              v-model="form.condition.values[index]"
                              type="text"
                              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                              :placeholder="getValuePlaceholder(form.condition.type, form.condition.param_name)"
                            />
                            <span class="text-sm text-gray-500">→</span>
                            <span class="text-sm text-gray-700">{{ target.url }}</span>
                          </div>
                        </div>
                        <p class="mt-1 text-sm text-gray-500">
                          {{ getConditionHelp(form.condition.type, form.condition.param_name) }}
                        </p>
                      </div>

                      <!-- Default Target -->
                      <div v-if="form.condition.type">
                        <label class="block text-sm font-medium text-gray-700">Default Target</label>
                        <select
                          v-model="form.condition.default"
                          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        >
                          <option
                            v-for="target in form.targets"
                            :key="target.id"
                            :value="target.id"
                          >
                            {{ target.url }}
                          </option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <!-- Tags and Targets -->
                  <div>
                    <TagsInput
                      v-model="form.tags"
                      :available-tags="availableTags"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Targets</label>
                    <div v-for="(target, index) in form.targets" :key="index" class="mt-2 flex items-center gap-2">
                      <input
                        v-model="target.url"
                        type="text"
                        required
                        placeholder="Target URL with optional query"
                        pattern="^(https?:\/\/)?([\da-z.\-]+\.)*[\da-z.\-]+\.([a-z.]{2,6})(\/[\w.\-]*)*(\?[\w.%\-]+(=[\w.%\-]*)?(&[\w.%\-]+(=[\w.%\-]*)?)*)?(#[\w%\-]*)?$"
                        class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                      <input
                        v-model.number="target.weight"
                        type="number"
                        min="0"
                        max="100"
                        placeholder="Weight"
                        class="block w-20 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                      <button
                        type="button"
                        @click="removeTarget(index)"
                        class="text-red-600 hover:text-red-900"
                      >
                        Remove
                      </button>
                    </div>
                    <button
                      type="button"
                      @click="addTarget"
                      class="mt-2 text-sm text-indigo-600 hover:text-indigo-900"
                    >
                      Add Target
                    </button>
                  </div>
                </div>
              </div>
              <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                <button
                  type="submit"
                  class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2"
                >
                  {{ editingProxy ? 'Save Changes' : 'Create' }}
                </button>
                <button
                  type="button"
                  @click="closeModal"
                  class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- History Modal -->
    <div v-if="showHistoryModal" class="relative z-10">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-2xl sm:p-6">
            <div class="absolute right-0 top-0 pr-4 pt-4">
              <button
                type="button"
                @click="closeHistoryModal"
                class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                <span class="sr-only">Close</span>
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left w-full">
                <h3 class="text-base font-semibold leading-6 text-gray-900">
                  Proxy History
                </h3>
                <div class="mt-4">
                  <ProxyHistory :proxy-id="selectedProxyId" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import axios from 'axios'
import TagsInput from '@/components/TagsInput.vue'
import ProxyHistory from '@/components/ProxyHistory.vue'
import Pagination from "@/components/Pagination.vue";

const proxies = ref([])
const totalProxies = ref(0)
const currentPage = ref(1)
const itemsPerPage = ref(10)
const sortBy = ref('id')
const sortDesc = ref(false)
const selectedTags = ref([])
const availableTags = ref([])
const showModal = ref(false)
const showHistoryModal = ref(false)
const editingProxy = ref(null)
const selectedProxyId = ref(null)

const form = ref({
  listen_url: '',
  mode: 'reverse',
  tags: [],
  targets: [],
  condition: {
    type: '',
    param_name: '',
    values: [],
    default: ''
  },
})

function getParamPlaceholder(type) {
  switch (type) {
    case 'header':
      return 'X-Custom-Header'
    case 'query':
      return 'version'
    case 'cookie':
      return 'user_preference'
    default:
      return ''
  }
}

function getValuePlaceholder(type, param) {
  switch (type) {
    case 'user_agent':
      return param === 'platform' ? 'mobile/desktop' : 'chrome/firefox/safari/edge/ie'
    case 'language':
      return 'en/es/fr/de'
    default:
      return 'Value for target'
  }
}

function getConditionHelp(type, param) {
  switch (type) {
    case 'user_agent':
      if (param === 'platform') {
        return 'Route traffic based on device type: mobile or desktop'
      }
      return 'Route traffic based on browser type: chrome, firefox, safari, edge, ie, or other'
    case 'language':
      return 'Route traffic based on preferred language from Accept-Language header (e.g., en, es, fr)'
    case 'header':
      return 'Route traffic based on a custom HTTP header value'
    case 'query':
      return 'Route traffic based on a URL query parameter value'
    case 'cookie':
      return 'Route traffic based on a cookie value'
    default:
      return ''
  }
}

const filteredProxies = computed(() => {
  if (selectedTags.value.length === 0) return proxies.value
  return proxies.value.filter(proxy => 
    selectedTags.value.every(tag => proxy.tags?.includes(tag))
  )
})

async function loadProxies() {
  try {
    const response = await axios.get('/api/proxies', {
      params: {
        limit: itemsPerPage.value,
        offset: (currentPage.value - 1) * itemsPerPage.value,
        sortBy: sortBy.value,
        sortDesc: sortDesc.value
      }
    })
    proxies.value = response.data.items.map((proxy) => {
      return {
        ...proxy,
        condition: {
          ...proxy.condition,
          values: proxy.condition?.values ? proxy.targets.map(({ id }) => proxy.condition.values[id]) : []
        }
      }
    })
    totalProxies.value = response.data.total
  } catch (error) {
    console.error('Failed to load proxies:', error)
  }
}

function handleSort(column) {
  if (sortBy.value === column) {
    sortDesc.value = !sortDesc.value
  } else {
    sortBy.value = column
    sortDesc.value = false
  }
  loadProxies()
}

function changePage(page) {
  currentPage.value = page
  loadProxies()
}

async function filterProxies() {
  if (selectedTags.value.length === 0) {
    await loadProxies()
  } else {
    try {
      const response = await axios.get('/api/proxies/by-tags', {
        params: { tags: selectedTags.value.join(',') }
      })
      proxies.value = response.data.proxies
    } catch (error) {
      console.error('Failed to filter proxies:', error)
    }
  }
}

function openCreateModal() {
  editingProxy.value = null
  form.value = {
    listen_url: '',
    mode: 'reverse',
    tags: [],
    targets: [],
    condition: {
      type: '',
      param_name: '',
      values: [],
      default: ''
    }
  }
  showModal.value = true
}

function editProxy(proxy) {
  editingProxy.value = proxy
  form.value = {
    listen_url: proxy.listen_url,
    mode: proxy.mode,
    tags: proxy.tags || [],
    targets: proxy.targets.map(t => ({
      url: t.url,
      weight: t.weight * 100,
      is_active: t.is_active
    })),
    condition: proxy.condition || {
      type: '',
      param_name: '',
      values: [],
      default: ''
    }
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingProxy.value = null
  form.value = {
    listen_url: '',
    mode: 'reverse',
    tags: [],
    targets: [],
    condition: {
      type: '',
      param_name: '',
      values: [],
      default: ''
    },
  }
}

function viewHistory(proxy) {
  selectedProxyId.value = proxy.id
  showHistoryModal.value = true
}

function closeHistoryModal() {
  showHistoryModal.value = false
  selectedProxyId.value = null
}

async function handleSubmit() {
  try {
    const formData = {
      listen_url: form.value.listen_url,
      mode: form.value.mode,
      tags: form.value.tags,
      condition: form.value.condition,
      targets: form.value.targets.map(t => ({
        url: t.url,
        weight: t.weight / 100,
        is_active: t.is_active ?? true
      }))
    }

    if (editingProxy.value) {
      await Promise.all([
        axios.put(`/api/proxies/${editingProxy.value.id}/targets`, formData),
        axios.put(`/api/proxies/${editingProxy.value.id}/tags`, { tags: formData.tags })
      ])
    } else {
      await axios.post('/api/proxies', formData)
    }

    await loadProxies()
    closeModal()
  } catch (error) {
    console.error('Failed to save proxy:', error)
    alert(error.response?.data?.error || 'Failed to save proxy')
  }
}

async function deleteProxy(id) {
  if (!confirm('Are you sure you want to delete this proxy?')) return

  try {
    await axios.delete(`/api/proxies/${id}`)
    await loadProxies()
  } catch (error) {
    console.error('Failed to delete proxy:', error)
  }
}

function addTarget() {
  form.value.targets.push({
    url: '',
    weight: 50,
    is_active: true
  })
}

function removeTarget(index) {
  form.value.targets.splice(index, 1)
}

onMounted(async () => {
  await loadProxies()
  const tagsResponse = await axios.get('/api/tags')
  availableTags.value = tagsResponse.data.tags || []
})
</script>
