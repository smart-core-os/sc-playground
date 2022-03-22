<template>
  <v-select v-model="selectedMode" :items="modesWithActive" item-text="title" item-value="id" return-object hide-details
            outlined :menu-props="{offsetY: true}" v-bind="$attrs">
    <template #item="{ item, on, attrs }">
      <v-list-item v-on="on" v-bind="attrs" class="mode-item">
        <v-list-item-title>{{ item.title }}</v-list-item-title>
        <v-list-item-action-text>
          {{ maxModeMagnitude(item).toFixed(1) }}A
        </v-list-item-action-text>
        <electric-segment-chart v-if="showSegmentChart(item)" :mode="item" class="chart"/>
      </v-list-item>
    </template>
    <template #selection="{ item }">
      <div style="width: 100%" class="d-flex">
        <span>{{ item.title }}</span>
        <span class="grey--text ml-auto">{{ maxModeMagnitude(item).toFixed(1) }}A</span>
      </div>
      <electric-segment-chart v-if="showSegmentChart(item)" :mode="item" class="chart mx-n3"/>
    </template>
  </v-select>
</template>

<script>
import ElectricSegmentChart from "./ElectricSegmentChart.vue";
import {maxMagnitude} from "./util.js";

export default {
  name: "ElectricModeChooser",
  components: {ElectricSegmentChart},
  props: {
    mode: [Object],
    modes: [Array]
  },
  data() {
    return {
      _selectedMode: null
    };
  },
  computed: {
    modesWithActive() {
      return this.modes.map(m => {
        if (m.id === this.mode.id) {
          return this.mode;
        }
        return m;
      })
    },
    singleMode() {
      const modes = this.modes;
      return !modes || modes.length === 0 || (modes.length === 1 && modes[0].id === this.mode.id);
    },
    activeModeTitle() {
      return this.mode?.title;
    },
    selectedMode: {
      get() {
        return this._selectedMode || this.mode;
      },
      set(mode) {
        this._selectedMode = mode;
        this.$emit("update:mode", mode);
      }
    }
  },
  watch: {
    mode() {
      this._selectedMode = null;
    }
  },
  methods: {
    maxModeMagnitude(mode) {
      return maxMagnitude(mode.segmentsList);
    },
    showSegmentChart(mode) {
      // only show the chart if there are segments and the first segment isn't infinite
      return mode.segmentsList && mode.segmentsList.length > 0 && mode.segmentsList[0].length
    }
  }
}
</script>

<style scoped>
.hero {
  font-size: 85px;
}

.v-card {
  text-align: center;
}

.mode-item {
  position: relative;
}

.chart {
  color: var(--v-primary-base);
  opacity: 0.1;
  position: absolute;
  inset: 12px 0 0 0;
}
</style>
