<template>
  <div class="sparkline" ref="chartRef"></div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import * as Plot from '@observablehq/plot'

export default {
  name: 'SparklineChart',
  props: {
    data: {
      type: Array,
      required: true
    },
    value: {
      type: String,
      required: true
    },
    color: {
      type: String,
      default: '#409EFF'
    },
    height: {
      type: Number,
      default: 30
    }
  },
  setup(props) {
    const chartRef = ref(null)

    const updateChart = () => {
      if (!chartRef.value || !props.data.length) return

      const chart = Plot.plot({
        width: chartRef.value.offsetWidth,
        height: props.height,
        margin: 0,
        x: { axis: null },
        y: { axis: null },
        marks: [
          Plot.line(props.data, {
            x: "timestamp",
            y: props.value,
            stroke: props.color
          })
        ]
      })

      chartRef.value.innerHTML = ''
      chartRef.value.appendChild(chart)
    }

    onMounted(updateChart)
    watch(() => props.data, updateChart, { deep: true })

    return {
      chartRef
    }
  }
}
</script>

<style scoped>
.sparkline {
  width: 100px;
  display: inline-block;
  vertical-align: middle;
  margin-left: 8px;
}
</style>
