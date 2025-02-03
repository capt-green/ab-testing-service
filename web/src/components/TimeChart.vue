<script setup>
import * as Plot from "@observablehq/plot";
import {onMounted, onUnmounted, ref, watch} from "vue";

const props = defineProps({
  chartType: String,
  heatmapMetric: String,
  targetStatsArray: Array,
  availableMetrics: Array,
  stackCharts: Boolean,
  normalizeData: Boolean,
})

watch([
  () => props.targetStatsArray,
  () => props.chartType,
  () => props.stackCharts,
  () => props.normalizeData,
  () => props.heatmapMetric
], () => {
  updateChart()
})

const chartRef = ref(null)

// Chart colors for different metrics
const metricColors = {
  requests: '#409EFF',
  errors: '#F56C6C',
  uniqueUsers: '#E6A23C',
  responseTime: '#67C23A'
}

const getMetricColor = (metric) => metricColors[metric] || '#409EFF'

const updateChart = () => {
  if (!chartRef.value || !props.targetStatsArray.length) return

  const data = [...props.targetStatsArray].sort((a, b) => a.timestamp - b.timestamp)
  const marks = [
    Plot.ruleY([0])
  ]

  // Prepare data for normalization if needed
  const normalizedData = props.normalizeData ?
      data.map(d => ({
        ...d,
        normalizedValue: d[props.heatmapMetric] / Math.max(...data.map(x => x[props.heatmapMetric]))
      })) : data

  switch (props.chartType) {
    case 'line':
      marks.push(Plot.line(normalizedData, {
        x: "timestamp",
        // curve: "basis",
        marker: true,
        y: props.normalizeData ? "normalizedValue" : props.heatmapMetric,
        stroke: props.stackCharts ? "targetId" : getMetricColor(props.heatmapMetric),
        strokeWidth: 2
      }))
      break
    case 'area':
      marks.push(Plot.areaY(normalizedData, {
        x: "timestamp",
        y: props.normalizeData ? "normalizedValue" : props.heatmapMetric,
        fill: props.stackCharts ? "targetId" : getMetricColor(props.heatmapMetric),
        opacity: 0.6
      }))
      break
    case 'bar':
      marks.push(Plot.rectY(normalizedData, {
        x: "timestamp",
        interval: "hour",
        y: props.normalizeData ? "normalizedValue" : props.heatmapMetric,
        fill: props.stackCharts ? "targetId" : getMetricColor(props.heatmapMetric)
      }))
      break
  }

  const chart = Plot.plot({
    width: chartRef.value.offsetWidth,
    height: 400,
    y: {
      grid: true,
      label: props.normalizeData ? "Normalized Value" : getMetricLabel(props.heatmapMetric)
    },
    marks: marks,
    color: props.stackCharts ? {legend: true} : null
  })

  chartRef.value.innerHTML = ''
  chartRef.value.appendChild(chart)
}

const getMetricLabel = (metric) => {
  const found = props.availableMetrics.find(m => m.value === metric)
  return found ? found.label : metric
}

const handleResize = () => {
  updateChart()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div ref="chartRef" class="chart"></div>
</template>

<style scoped>
.chart {
  width: 100%;
  /*height: 400px;*/
  margin: 10px 0;
}
</style>