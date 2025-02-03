<script setup>
defineProps({
  stats: {
    type: Object,
    required: true
  }
})
</script>

<template>
  <el-row :gutter="20" class="stats-cards">
    <el-col :span="8">
      <el-card class="stats-card">
        <template #header>
          <div class="card-header">
            <span>Total Requests</span>
          </div>
        </template>
        <div class="stats-value">{{ stats.totalRequests.toLocaleString() }}</div>
      </el-card>
    </el-col>
    <el-col :span="8">
      <el-card class="stats-card">
        <template #header>
          <div class="card-header">
            <span>Total Errors</span>
          </div>
        </template>
        <div class="stats-value">{{ stats.totalErrors.toLocaleString() }}</div>
        <div v-if="stats.totalRequests" class="error-rate">Error Rate:
          {{ ((stats.totalErrors / stats.totalRequests) * 100 || 0).toFixed(2) }}%
        </div>
      </el-card>
    </el-col>
    <el-col :span="8">
      <el-card class="stats-card">
        <template #header>
          <div class="card-header">
            <span>Unique Users</span>
          </div>
        </template>
        <div class="stats-value">{{ stats.uniqueUsers.toLocaleString() }}</div>
      </el-card>
    </el-col>
  </el-row>
</template>

<style scoped>
.stats-cards {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
  height: 100%;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  margin: 10px 0;
}

.error-rate {
  font-size: 14px;
  color: #666;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>