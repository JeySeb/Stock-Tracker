import * as echarts from 'echarts'

export interface ChartTheme {
  backgroundColor: string
  textColor: string
  primaryColor: string
  successColor: string
  dangerColor: string
  warningColor: string
}

export const defaultTheme: ChartTheme = {
  backgroundColor: '#ffffff',
  textColor: '#374151',
  primaryColor: '#3b82f6',
  successColor: '#10b981',
  dangerColor: '#ef4444',
  warningColor: '#f59e0b'
}

export const darkTheme: ChartTheme = {
  backgroundColor: '#1f2937',
  textColor: '#f9fafb',
  primaryColor: '#60a5fa',
  successColor: '#34d399',
  dangerColor: '#f87171',
  warningColor: '#fbbf24'
}

interface ChartDataPoint {
  value: number
  name: string
}

interface PriceTargetDataPoint {
  date: string
  target: number
  current: number
}

interface HeatmapDataPoint {
  x: string
  y: string
  value: number
}

export function createRecommendationChart(data: ChartDataPoint[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10,
      textStyle: {
        color: theme.textColor
      }
    },
    series: [
      {
        name: 'Recommendations',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold',
            color: theme.textColor
          }
        },
        labelLine: {
          show: false
        },
        data: data.map(item => ({
          value: item.value,
          name: item.name,
          itemStyle: {
            color: getColorForRecommendationType(item.name, theme)
          }
        }))
      }
    ]
  }
}

export function createPriceTargetChart(data: PriceTargetDataPoint[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['Price Target', 'Current Price'],
      textStyle: {
        color: theme.textColor
      }
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      axisLabel: {
        color: theme.textColor
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: theme.textColor,
        formatter: '${value}'
      }
    },
    series: [
      {
        name: 'Price Target',
        type: 'line',
        data: data.map(item => item.target),
        lineStyle: {
          color: theme.primaryColor
        },
        itemStyle: {
          color: theme.primaryColor
        }
      },
      {
        name: 'Current Price',
        type: 'line',
        data: data.map(item => item.current),
        lineStyle: {
          color: theme.successColor
        },
        itemStyle: {
          color: theme.successColor
        }
      }
    ]
  }
}

interface SentimentChartDataPoint {
  date: string
  sentiment: number
  volume: number
}

export function createSentimentChart(data: SentimentChartDataPoint[], theme: ChartTheme = defaultTheme) {
  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['Sentiment Score', 'Volume'],
      textStyle: {
        color: theme.textColor
      }
    },
    xAxis: [
      {
        type: 'category',
        data: data.map(item => item.date),
        axisLabel: {
          color: theme.textColor
        }
      }
    ],
    yAxis: [
      {
        type: 'value',
        name: 'Sentiment',
        min: -1,
        max: 1,
        axisLabel: {
          color: theme.textColor
        }
      },
      {
        type: 'value',
        name: 'Volume',
        axisLabel: {
          color: theme.textColor
        }
      }
    ],
    series: [
      {
        name: 'Sentiment Score',
        type: 'line',
        yAxisIndex: 0,
        data: data.map(item => item.sentiment),
        lineStyle: {
          color: theme.primaryColor
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: theme.primaryColor + '40' },
            { offset: 1, color: theme.primaryColor + '10' }
          ])
        }
      },
      {
        name: 'Volume',
        type: 'bar',
        yAxisIndex: 1,
        data: data.map(item => item.volume),
        itemStyle: {
          color: theme.warningColor + '60'
        }
      }
    ]
  }
}

function getColorForRecommendationType(type: string, theme: ChartTheme): string {
  switch (type.toLowerCase()) {
    case 'strong buy':
      return theme.successColor
    case 'buy':
      return '#22c55e'
    case 'hold':
      return theme.warningColor
    case 'sell':
      return '#f97316'
    case 'strong sell':
      return theme.dangerColor
    default:
      return theme.textColor
  }
}

export function createHeatmapChart(data: HeatmapDataPoint[], theme: ChartTheme = defaultTheme) {
  const hours = ['12a', '1a', '2a', '3a', '4a', '5a', '6a', '7a', '8a', '9a', '10a', '11a',
                 '12p', '1p', '2p', '3p', '4p', '5p', '6p', '7p', '8p', '9p', '10p', '11p']
  const days = ['Saturday', 'Friday', 'Thursday', 'Wednesday', 'Tuesday', 'Monday', 'Sunday']

  return {
    backgroundColor: theme.backgroundColor,
    tooltip: {
      position: 'top',
      formatter: function(params: { data: [number, number, number] }) {
        return `${days[params.data[1]]} ${hours[params.data[0]]}<br/>Activity: ${params.data[2]}`
      }
    },
    grid: {
      height: '50%',
      top: '10%'
    },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: {
        show: true
      },
      axisLabel: {
        color: theme.textColor
      }
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: {
        show: true
      },
      axisLabel: {
        color: theme.textColor
      }
    },
    visualMap: {
      min: 0,
      max: 10,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '15%',
      textStyle: {
        color: theme.textColor
      }
    },
    series: [{
      name: 'Trading Activity',
      type: 'heatmap',
      data: data,
      label: {
        show: true,
        color: theme.textColor
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }]
  }
}