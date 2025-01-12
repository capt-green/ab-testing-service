<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-medium leading-6 text-gray-900">History</h3>
      <button
        @click="loadHistory"
        class="text-sm text-indigo-600 hover:text-indigo-900"
      >
        Refresh
      </button>
    </div>

    <div class="flow-root">
      <ul role="list" class="-mb-8">
        <li v-for="(event, eventIdx) in history" :key="event.id">
          <div class="relative pb-8">
            <span
              v-if="eventIdx !== history.length - 1"
              class="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200"
              aria-hidden="true"
            />
            <div class="relative flex space-x-3">
              <!-- Event icon -->
              <div>
                <span
                  :class="[
                    getEventColor(event.type),
                    'h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white'
                  ]"
                >
                  <component
                    :is="getEventIcon(event.type)"
                    class="h-5 w-5 text-white"
                    aria-hidden="true"
                  />
                </span>
              </div>

              <!-- Event content -->
              <div class="flex min-w-0 flex-1 justify-between space-x-4">
                <div>
                  <p class="text-sm text-gray-500">
                    {{ getEventDescription(event) }}
                  </p>
                  <div
                    v-if="event.changes"
                    class="mt-2 text-sm text-gray-700"
                  >
                    <div
                      v-for="(change, field) in event.changes"
                      :key="field"
                      class="mt-1"
                    >
                      <span class="font-medium">{{ formatFieldName(field) }}:</span>
                      <template v-if="Array.isArray(change)">
                        <span class="text-red-600 line-through">{{ formatValue(change[0]) }}</span>
                        <span class="mx-1">→</span>
                        <span class="text-green-600">{{ formatValue(change[1]) }}</span>
                      </template>
                      <template v-else>
                        <span class="text-green-600">{{ formatValue(change) }}</span>
                      </template>
                    </div>
                  </div>
                </div>
                <div class="whitespace-nowrap text-right text-sm text-gray-500">
                  <time :datetime="event.timestamp">{{ formatDate(event.timestamp) }}</time>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <div v-if="history.length === 0" class="text-center py-6">
      <p class="text-sm text-gray-500">No history available</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import {
  DocumentPlusIcon,
  TrashIcon,
  PencilSquareIcon,
  TagIcon,
  ArrowPathIcon,
} from '@heroicons/vue/24/outline'

const props = defineProps({
  proxyId: {
    type: String,
    required: true
  }
})

const history = ref([])

const eventIcons = {
  'create': DocumentPlusIcon,
  'delete': TrashIcon,
  'update': PencilSquareIcon,
  'update_tags': TagIcon,
  'update_targets': ArrowPathIcon
}

const eventColors = {
  'create': 'bg-green-500',
  'delete': 'bg-red-500',
  'update': 'bg-blue-500',
  'update_tags': 'bg-indigo-500',
  'update_targets': 'bg-yellow-500'
}

function getEventIcon(type) {
  return eventIcons[type] || PencilSquareIcon
}

function getEventColor(type) {
  return eventColors[type] || 'bg-gray-500'
}

function formatFieldName(field) {
  return field
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

function formatValue(value) {
  if (value === null || value === undefined) return 'None'
  if (Array.isArray(value)) return value.join(', ')
  if (typeof value === 'boolean') return value ? 'Yes' : 'No'
  return value.toString()
}

function formatDate(timestamp) {
  return new Date(timestamp).toLocaleString()
}

function getEventDescription(event) {
  const descriptions = {
    'create': 'Created proxy',
    'delete': 'Deleted proxy',
    'update': 'Updated proxy configuration',
    'update_tags': 'Updated proxy tags',
    'update_targets': 'Updated proxy targets'
  }
  return descriptions[event.type] || 'Modified proxy'
}

async function loadHistory() {
  try {
    const response = await axios.get(`/api/proxies/${props.proxyId}/history`)
    history.value = response.data.history
  } catch (error) {
    console.error('Failed to load proxy history:', error)
  }
}

onMounted(loadHistory)
</script>
