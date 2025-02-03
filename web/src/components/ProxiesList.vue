<template>
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
            <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ proxy.mode === 'path' ? proxy.path_key : proxy.listen_url }}</td>
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
</template>

<script setup lang="ts">
type Proxy = {
  id: string,
  listen_url?: string,
  path_key?: string,
  mode: string,
  tags: string[],
  targets: Array<{ id: string, url: string, weight: number, is_active: boolean }>
}
defineProps<{
  filteredProxies: Array<Proxy>,
  sortBy: string,
  sortDesc: boolean
}>()

const emit = defineEmits(['delete', 'edit', 'viewHistory', 'sort'])

const handleSort = (column) => {
  emit('sort', column)
}

const editProxy = (proxy) => {
  emit('edit', proxy)
}

const viewHistory = (proxy) => {
  emit('viewHistory', proxy)
}

const deleteProxy = (id) => {
  emit('delete', id)
}

</script>