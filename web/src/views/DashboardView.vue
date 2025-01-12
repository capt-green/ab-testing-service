<template>
  <div>
    <div class="bg-white shadow">
      <div class="px-4 sm:px-6 lg:mx-auto lg:max-w-6xl lg:px-8">
        <div class="py-6 md:flex md:items-center md:justify-between lg:border-t lg:border-gray-200">
          <div class="min-w-0 flex-1">
            <div class="flex items-center">
              <div>
                <div class="flex items-center">
                  <h1 class="ml-3 text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:leading-9">
                    Welcome back!
                  </h1>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-8">
      <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <h2 class="text-lg font-medium leading-6 text-gray-900">Overview</h2>
        <div class="mt-2 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          <!-- Active Proxies Card -->
          <div class="overflow-hidden rounded-lg bg-white shadow">
            <div class="p-5">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="truncate text-sm font-medium text-gray-500">Active Proxies</dt>
                    <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.activeProxies }}</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <!-- Total Requests Card -->
          <div class="overflow-hidden rounded-lg bg-white shadow">
            <div class="p-5">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
                  </svg>
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="truncate text-sm font-medium text-gray-500">Total Requests (24h)</dt>
                    <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.totalRequests }}</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <!-- Success Rate Card -->
          <div class="overflow-hidden rounded-lg bg-white shadow">
            <div class="p-5">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="truncate text-sm font-medium text-gray-500">Success Rate</dt>
                    <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ stats.successRate }}%</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8 mt-8">
        <h2 class="text-lg font-medium leading-6 text-gray-900 mb-4">Recent Activity</h2>
        <div class="overflow-hidden bg-white shadow sm:rounded-md">
          <ul role="list" class="divide-y divide-gray-200">
            <li v-for="activity in recentActivity" :key="activity.id">
              <div class="px-4 py-4 sm:px-6">
                <div class="flex items-center justify-between">
                  <div class="truncate text-sm font-medium text-indigo-600">{{ activity.proxyId }}</div>
                  <div class="ml-2 flex flex-shrink-0">
                    <span :class="[
                      activity.status === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800',
                      'inline-flex rounded-full px-2 text-xs font-semibold leading-5'
                    ]">
                      {{ activity.status }}
                    </span>
                  </div>
                </div>
                <div class="mt-2 flex justify-between">
                  <div class="sm:flex">
                    <div class="mr-6 flex items-center text-sm text-gray-500">
                      <svg class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                      {{ activity.timestamp }}
                    </div>
                    <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                      <svg class="mr-1.5 h-5 w-5 flex-shrink-0 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                      </svg>
                      {{ activity.targetUrl }}
                    </div>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const stats = ref({
  activeProxies: 0,
  totalRequests: 0,
  successRate: 0
})

const recentActivity = ref([])

async function fetchDashboardData() {
  try {
    const [statsResponse, activityResponse] = await Promise.all([
      axios.get('/api/stats'),
      axios.get('/api/activity/recent')
    ])
    
    stats.value = statsResponse.data
    recentActivity.value = activityResponse.data
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  }
}

onMounted(fetchDashboardData)
</script>
