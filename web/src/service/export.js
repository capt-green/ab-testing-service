import {unref} from "vue";
import {unparse} from "papaparse";

export const exportChart = (chartRef, chartName) => {
    const svg = chartRef.querySelector('svg')
    const svgData = new XMLSerializer().serializeToString(svg)
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    const img = new Image()

    img.onload = () => {
        canvas.width = svg.width.baseVal.value
        canvas.height = svg.height.baseVal.value
        ctx.drawImage(img, 0, 0)
        const link = document.createElement('a')
        link.download = `${chartName}.png`
        link.href = canvas.toDataURL('image/png')
        link.click()
    }

    img.src = 'data:image/svg+xml;base64,' + btoa(svgData)
}

export const exportToCSV = (targetStatsArray, proxyName) => {
    const targetStats = unref(targetStatsArray)
    const data = targetStats.map(stat => ({
        Target: stat.targetId,
        Requests: stat.requests,
        Errors: stat.errors,
        'Error Rate': `${(stat.errorRate * 100).toFixed(2)}%`,
        'Unique Users': stat.uniqueUsers,
        Timestamp: new Date(stat.timestamp).toLocaleString()
    }))

    const csv = unparse(data)
    const blob = new Blob([csv], {type: 'text/csv;charset=utf-8;'})
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `proxy-${proxyName}-stats.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}