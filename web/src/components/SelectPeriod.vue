<script setup>
const selectedPeriod = defineModel('selectedPeriod')
const timeRange = defineModel('timeRange')
const refreshInterval = defineModel('refreshInterval')
const emit = defineEmits(['periodChange', 'force-refresh', 'setupAutoRefresh'])
const defaultTime = [
  new Date().setDate(new Date().getDate() - 7),
  new Date()
]
</script>

<template>
  <div class="stats-header">
    <!-- Preset time periods -->
    <el-radio-group v-model="selectedPeriod" @change="emit('periodChange', $event)">
      <el-radio-button value="today">Today</el-radio-button>
      <el-radio-button value="week">Week</el-radio-button>
      <el-radio-button value="month">Month</el-radio-button>
      <el-radio-button value="custom">Custom</el-radio-button>
    </el-radio-group>

    <!-- Custom date range -->
    <el-date-picker
        v-if="selectedPeriod === 'custom'"
        v-model="timeRange"
        type="daterange"
        range-separator="to"
        start-placeholder="Start date"
        end-placeholder="End date"
        :default-time="defaultTime"
        @change="emit('periodChange', $event)"
    />
    <div class="header-controls">

      <!-- Auto-refresh control -->
      <el-select v-model="refreshInterval" placeholder="Auto-refresh" @change="emit('setupAutoRefresh', $event)">
        <el-option label="Off" :value="0"/>
        <el-option label="30 seconds" :value="30"/>
        <el-option label="1 minute" :value="60"/>
        <el-option label="5 minutes" :value="300"/>
      </el-select>

      <el-button type="primary" @click="$emit('force-refresh')">Refresh</el-button>
    </div>
  </div>
</template>

<style scoped>
.stats-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  margin-block: 20px;
}

.header-controls {
  display: flex;
  gap: 16px;
  align-items: center;
  min-width: 200px;
}
</style>