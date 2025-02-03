<script setup>
import SparklineChart from "@/components/SparklineChart.vue";
import {computed, ref, watch} from "vue";

const props = defineProps({
  availableMetrics: {
    type: Array,
    required: true
  },
  targetStatsArray: {
    type: Array, // flat history with target id, timestamp, requests, errors, unique users
    required: true
  },
})
const selectedMetrics = ref(['requests', 'errors', 'uniqueUsers'])
const checkAll = ref(true)
const indeterminate = ref(false)
const showSparklines = ref(true)
const availableGroups = [
  {label: 'Timestamp', value: 'timestamp'},
  {label: 'Target', value: 'targetId'},
  {label: 'Date', value: 'date'},
  {label: 'Version', value: 'version'},
  {label: 'Rule', value: 'rule'},
  {label: 'Country', value: 'country'},
  {label: 'Status', value: 'status'},
]
const groupBy = ref(['timestamp'])

const preparedStats = computed(() => {
  return props.targetStatsArray.map(stat => ({
    ...stat,
    uniqueUsers: stat.users_count,
    date: stat.timestamp.toLocaleDateString(),
    time: stat.timestamp.toLocaleTimeString(),
    timestamp: stat.timestamp.toUTCString(),
    history: []
  }))
})

const aggregatedStats = computed(() => {
  if (groupBy.value.length === 0) {
    return preparedStats.value
  }

  // Создаем результат группировки по всем полям сразу
  const groupedData = preparedStats.value.reduce((acc, stat) => {
    // Создаем составной ключ из всех полей группировки
    const groupKey = groupBy.value.map(field => stat[field]).join('|')

    if (!acc[groupKey]) {
      acc[groupKey] = {
        history: [],
        requests: 0,
        errors: 0,
        uniqueUsers: 0
      }

      // Добавляем поля группировки в результат
      groupBy.value.forEach(field => {
        acc[groupKey][field] = stat[field]
      })
    }

    // Накапливаем значения
    acc[groupKey].history.push(stat)
    acc[groupKey].requests += stat.requests
    acc[groupKey].errors += stat.errors
    acc[groupKey].uniqueUsers += stat.uniqueUsers

    return acc
  }, {})

  // Преобразуем объект группировки в массив
  return Object.values(groupedData)
})

const formatMetricValue = (metric, value) => {
  switch (metric) {
    case 'errorRate':
      return `${(value * 100).toFixed(2)}%`
    case 'requests':
    case 'errors':
    case 'uniqueUsers':
      return value.toLocaleString()
    default:
      return value
  }
}

const metricColors = {
  requests: '#409EFF',
  errors: '#F56C6C',
  uniqueUsers: '#E6A23C',
  responseTime: '#67C23A'
}

const getMetricColor = (metric) => metricColors[metric] || '#409EFF'

const getMetricLabel = (metric) => {
  const found = props.availableMetrics.find(m => m.value === metric)
  return found ? found.label : metric
}

const handleCheckAll = (val) => {
  indeterminate.value = false
  selectedMetrics.value = val ? props.availableMetrics.map(({value}) => value) : []
}

watch(selectedMetrics, (val) => {
  if (val.length === 0) {
    checkAll.value = false
    indeterminate.value = false
  } else if (val.length === props.availableMetrics.length) {
    checkAll.value = true
    indeterminate.value = false
  } else {
    indeterminate.value = true
  }
})

</script>

<template>
  <el-card class="proxy-stats">
    <template #header>
      <div class="card-header">
        <span>Target Statistics</span>
        <div class="header-controls">
          <el-select v-model="groupBy" multiple :multiple-limit="3" clearable placeholder="Select groups">
            <el-option
                v-for="field in availableGroups"
                :key="field.value"
                :label="field.label"
                :value="field.value"
            />
          </el-select>
        </div>
        <div class="header-controls">
          <el-select v-model="selectedMetrics" multiple placeholder="Select metrics">
            <template #header>
              <el-checkbox
                  v-model="checkAll"
                  :indeterminate="indeterminate"
                  @change="handleCheckAll"
              >
                All
              </el-checkbox>
            </template>
            <el-option
                v-for="metric in props.availableMetrics"
                :key="metric.value"
                :label="metric.label"
                :value="metric.value"
            />
          </el-select>
        </div>
      </div>
    </template>

    <el-table :data="aggregatedStats" table-layout="auto" style="width: 100%">
      <el-table-column v-for="group in groupBy" :key="group" :label="getMetricLabel(group)" width="180">
        <template #default="scope">
          <span>{{ scope.row[group] }}</span>
        </template>
      </el-table-column>
      <el-table-column v-for="metric in selectedMetrics" sortable :key="metric" :label="getMetricLabel(metric)">
        <template #default="scope">
          <div class="metric-cell">
            <span>{{ formatMetricValue(metric, scope.row[metric]) }}</span>
            <sparkline-chart
                v-if="showSparklines && scope.row.history.length > 1"
                :data="scope.row.history"
                :value="metric"
                :color="getMetricColor(metric)"
            />
          </div>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<style scoped>
.proxy-stats {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.header-controls {
  display: flex;
  gap: 16px;
  align-items: center;
  min-width: 200px;
}

.metric-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

</style>