<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="12">
        <h1>Dashboard</h1>
      </el-col>
      <el-col :span="12">
        <el-select v-model="selectedProxy" placeholder="Selected Proxy">
          <el-option v-for="proxy in proxies" :key="proxy.id" :label="proxy.listen_url" :value="proxy.id"/>
        </el-select>
      </el-col>
    </el-row>
    <SelectPeriod
        v-model:selected-period="selectedPeriod"
        v-model:time-range="timeRange"
        v-model:refresh-interval="refreshInterval"
        @periodChange="handlePeriodChange"
        @setupAutoRefresh="setupAutoRefresh"
        @force-refresh="fetchStats"
    />

    <!-- Factoids -->
    <StatsCards :stats="stats"/>

    <!-- Time Chart -->
    <el-row :gutter="20" class="charts-section">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <div class="chart-controls">
                <span>Requests Over Time</span>
                <el-switch
                    v-model="compareWithPrevious"
                    active-text="Compare with previous period"
                    @change="fetchComparisonData"
                />
                <div class="control-group">
                  <el-select v-model="chartType" placeholder="Chart Type">
                    <el-option label="Line" value="line"/>
                    <el-option label="Area" value="area"/>
                    <el-option label="Bar" value="bar"/>
                  </el-select>
                  <el-checkbox v-model="stackCharts">Stack Charts</el-checkbox>
                  <el-checkbox v-model="normalizeData">Normalize</el-checkbox>
                </div>
              </div>
            </div>
          </template>
          <TimeChart
              :chart-type="chartType"
              :heatmap-metric="heatmapMetric"
              :stack-charts="stackCharts"
              :normalize-data="normalizeData"
              :available-metrics="availableMetrics"
              :target-stats-array="targetStatsArray"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- Time Heatmap -->
    <el-row :gutter="20" class="charts-section">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>Requests by Time</span>
              <el-select v-model="heatmapMetric" style="width: 150px">
                <el-option label="Request Count" value="requests"/>
                <el-option label="Errors Count" value="errors"/>
                <el-option label="Users Count" value="users_count"/>
              </el-select>
            </div>
          </template>
          <TimeHeatmap
              :heatmap-metric="heatmapMetric"
              :available-metrics="availableMetrics"
              :target-stats-array="targetStatsArray"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- Stats Table -->
    <StatsTable
        :available-metrics="availableMetrics"
        :target-stats-array="targetStatsArray"
    />

  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import axios from 'axios'
import SelectPeriod from '../components/SelectPeriod.vue'
import StatsCards from '../components/StatsCards.vue'
import TimeChart from "@/components/TimeChart.vue";
import TimeHeatmap from "@/components/TimeHeatmap.vue";
import StatsTable from "@/components/StatsTable.vue";

const timeRange = ref([new Date().setDate(new Date().getDate() - 7), new Date()])
const stats = ref({
  totalRequests: 0,
  totalErrors: 0,
  uniqueUsers: 0
})
const proxies = ref([])
const selectedProxy = ref('')
const proxyStats = ref({})
const selectedPeriod = ref('week')
const refreshInterval = ref(0)
const compareWithPrevious = ref(false)

let refreshTimer = null

// Metrics configuration
const availableMetrics = [
  {label: 'Requests', value: 'requests'},
  {label: 'Errors', value: 'errors'},
  {label: 'Unique Users', value: 'uniqueUsers'}
]

const targetStatsArray = computed(() => {
  if (!proxyStats.value.targetStats) return []
  return Object.entries(proxyStats.value.targetStats).flatMap(([targetId, stats]) => stats.map(stat => ({
    targetId,
    ...stat,
    timestamp: new Date(stat.timestamp)
  })))
})

const formatDate = (date) => {
  return new Date(date).toISOString()
}

const chartType = ref('line')
const stackCharts = ref(false)
const normalizeData = ref(false)
const heatmapMetric = ref('requests')

const handlePeriodChange = (period) => {
  const now = new Date()
  switch (period) {
    case 'today':
      timeRange.value = [
        new Date(now.setHours(0, 0, 0, 0)),
        new Date()
      ]
      break
    case 'week':
      timeRange.value = [
        new Date(now.setDate(now.getDate() - 7)),
        new Date()
      ]
      break
    case 'month':
      timeRange.value = [
        new Date(now.setMonth(now.getMonth() - 1)),
        new Date()
      ]
      break
  }
  fetchStats()
}

const setupAutoRefresh = (interval) => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
  if (interval > 0) {
    refreshTimer = setInterval(() => {
      fetchStats()
      if (selectedProxy.value) {
        fetchProxyStats()
      }
    }, interval * 1000)
  }
}

const fetchStats = async () => {
  try {
    const [start, end] = timeRange.value
    const response = await axios.get('/api/stats', {
      params: {
        start_time: formatDate(start),
        end_time: formatDate(end)
      }
    })

    stats.value = {
      totalRequests: response.data.total_requests,
      totalErrors: response.data.total_errors,
      uniqueUsers: response.data.unique_users
    }
  } catch (error) {
    ElMessage.error('Failed to fetch statistics')
    console.error('Error fetching stats:', error)
  }
}

const fetchProxies = async () => {
  try {
    const response = await axios.get('/api/proxies')
    proxies.value = response.data?.items
  } catch (error) {
    ElMessage.error('Failed to fetch proxies')
    console.error('Error fetching proxies:', error)
  }
}

const fetchProxyStats = async () => {
  if (!selectedProxy.value) return

  try {
    const [start, end] = timeRange.value
    const response = await axios.get(`/api/stats/${selectedProxy.value}`, {
      params: {
        start_time: formatDate(start),
        end_time: formatDate(end)
      }
    })
    proxyStats.value = {
      ...response.data,
      targetStats: response.data.target_stats,
    }
  } catch (error) {
    ElMessage.error('Failed to fetch proxy statistics')
    console.error('Error fetching proxy stats:', error)
  }
}

const fetchComparisonData = async () => {
  if (!compareWithPrevious.value) return

  const [start, end] = timeRange.value
  const duration = end - start
  const previousStart = new Date(start - duration)
  const previousEnd = new Date(end - duration)

  try {
    const response = await axios.get('/api/stats', {
      params: {
        start_time: formatDate(previousStart),
        end_time: formatDate(previousEnd)
      }
    })
    // Update charts with comparison data...
  } catch (error) {
    ElMessage.error('Failed to fetch comparison data')
  }
}

onMounted(() => {
  fetchStats()
  fetchProxies()
  handlePeriodChange('week')
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})

watch(timeRange, () => {
  fetchStats()
  if (selectedProxy.value) {
    fetchProxyStats()
  }
})

</script>

<style scoped>
@import "dashboard.css";
</style>