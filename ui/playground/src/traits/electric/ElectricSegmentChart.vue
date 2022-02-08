<template>
  <div class="chart">
    <div v-for="(item, i) of chartBars" :key="i" :style="item.style" :class="item.class"/>
  </div>
</template>
<script>

import {maxMagnitude} from "./util.js";

export default {
  name: 'ElectricSegmentChart',
  props: {
    mode: [Object]
  },
  computed: {
    /**
     * @return {ElectricMode.Segment.AsObject[]}
     */
    segments() {
      return this.mode?.segmentsList || [];
    },
    maxMagnitude() {
      return maxMagnitude(this.segments);
    },
    maxOrderOfMagnitude() {
      return this.segments.reduce((max, s) => Math.max(max, `${s.length?.seconds || 0}`.length), 0);
    },
    chartBars() {
      const max = this.maxMagnitude;
      return this.segments.map((s, i, arr) => {
        const notLast = i < arr.length - 1;

        const style = {};
        const classes = [];

        style.height = ((s.magnitude / max) * 100).toFixed(3) + '%';
        if (notLast || s.length) {
          style.width = `${s.length.seconds}px`;
        } else {
          style.width = `${Math.pow(10, this.maxOrderOfMagnitude) * 0.3}px`;
          style.background = 'linear-gradient(90deg, currentColor 25%, transparent)'
        }

        return {style, class: classes};
      })
    }
  }
}
</script>
<style scoped>
.chart {
  display: flex;
  align-items: flex-end;
}

.chart > * {
  background-color: currentColor;
  flex-grow: 1;
}
</style>
