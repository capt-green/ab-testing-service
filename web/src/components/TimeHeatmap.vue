<script setup>
import {onMounted, onUnmounted, ref, watch} from "vue";
import * as Plot from "@observablehq/plot";

const props = defineProps({
  targetStatsArray: Array,
  heatmapMetric: String,
  availableMetrics: Array
})

// Watch for changes that should trigger chart updates
watch([() => props.heatmapMetric, () => props.targetStatsArray], () => {
  updateHeatmap()
})

const timeHeatmapRef = ref(null)

const updateHeatmap = () => {
  if (!timeHeatmapRef.value || !props.targetStatsArray.length) return

  const heatmap = Plot.plot({
    width: timeHeatmapRef.value.offsetWidth,
    height: 300,
    x: {
      label: "Hour",
      tickFormat: d => `${d}:00`
    },
    y: {
      padding: 0,
      domain: Array.from({length: 7}, (_, i) => i),
      tickFormat: Plot.formatWeekday("en", "short"),
      tickSize: 0
    },
    marks: [
      Plot.cell(props.targetStatsArray, {
        x: (d) => d.timestamp.getUTCHours(),
        y: (d) => d.timestamp.getUTCDay(),
        fill: d => d[props.heatmapMetric],
        inset: 0.5,
      })
    ],
    color: {
      scheme: "viridis",
      legend: true,
      label: getMetricLabel(props.heatmapMetric)
    }
  })

  timeHeatmapRef.value.innerHTML = ''
  timeHeatmapRef.value.appendChild(heatmap)
}

const getMetricLabel = (metric) => {
  const found = props.availableMetrics.find(m => m.value === metric)
  return found ? found.label : metric
}

const handleResize = () => {
  updateHeatmap()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div ref="timeHeatmapRef" class="heatmap-chart"></div>
</template>

<style scoped>
.heatmap-chart {
  width: 100%;
  margin: 10px 0;
}
</style>