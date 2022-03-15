<template>
  <div class="chart">
    <div v-for="(item, i) of chartBars" :key="i" :style="item.style" :class="item.class"/>
  </div>
</template>
<script>

import {maxMagnitude} from "./util.js";
import {interval} from "../../mixin/time.js";

export default {
  name: 'ElectricSegmentChart',
  mixins: [interval(200)],
  props: {
    mode: [Object]
  },
  computed: {
    /**
     * @return {ElectricMode.Segment.AsObject[]}
     */
    allSegments() {
      return this.mode?.segmentsList || [];
    },
    segments() {
      const segments = this.allSegments;
      if (segments.length === 0 || !this.mode?.startTime) {
        return segments;
      }
      const st = toDate(this.mode.startTime);
      const now = this.currentDate;

      // cut off segments starting at st until we reach now.
      const offsetMillis = now.getTime() - st.getTime();
      let curMillis = 0;
      for (let i = 0; i < segments.length; i++) {
        const segment = segments[i];
        if (!segment.length) {
          return [segment]; // only the last, infinite segment left
        }

        const lengthMillis = segment.length.seconds * 1000 + segment.length.nanos / 1_000_000;
        if (curMillis < offsetMillis - lengthMillis) {
          curMillis += lengthMillis;
          continue; // the segment is before now
        }

        // now is somewhere within segment
        const remainingMillis = curMillis + lengthMillis - offsetMillis;
        const remainingSeconds = Math.floor(remainingMillis / 1000)
        const remainingSegment = {
          length: {
            seconds: remainingSeconds,
            nanos: (remainingMillis - (remainingSeconds * 1000)) * 1_000_000
          },
          magnitude: segment.magnitude
        }

        const remainingSegments = segments.slice(i);
        remainingSegments[0] = remainingSegment;
        return remainingSegments
      }

      return []
    },
    maxMagnitude() {
      return maxMagnitude(this.allSegments);
    },
    totalLength() {
      return this.allSegments.reduce((len, s) => {
        if (!s.length) {
          return len;
        }
        return len + durationMillis(s.length);
      }, 0);
    },
    chartBars() {
      const max = this.maxMagnitude;
      const totalLength = this.totalLength;
      const scale = 300 / totalLength;
      return this.segments.map((s, i, arr) => {
        const notLast = i < arr.length - 1;

        const style = {};
        const classes = [];

        style.height = ((s.magnitude / max) * 100).toFixed(3) + '%';
        if (notLast || s.length) {
          style.width = `${(durationMillis(s.length) * scale).toFixed(3)}px`;
        } else {
          style.width = `50px`;
          style.background = 'linear-gradient(90deg, currentColor 25%, transparent)';
          style.flexGrow = '1';
        }

        return {style, class: classes};
      })
    }
  }
}

/**
 * Convert a timestamp object to a Date. Timestamp.AsObject doesn't have this method :(
 *
 * @param {Timestamp.AsObject} ts
 * @return {Date}
 */
function toDate(ts) {
  return new Date((ts.seconds * 1000) + (ts.nanos / 1000000));
}

/**
 * Convert a duration object into a millisecond value.
 *
 * @param {Duration.AsObject} d
 * @return {number}
 */
function durationMillis(d) {
  return d.seconds * 1000 + d.nanos / 1_000_000;
}
</script>
<style scoped>
.chart {
  display: flex;
  align-items: flex-end;
}

.chart > * {
  background-color: currentColor;
  flex-grow: 0;
}
</style>
