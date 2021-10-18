<template>
  <div class="power-supply-bar">
    <div class="bar rating" :style="styles.rating" :title="`Rating: ${rating.toFixed(1)} A`"/>
    <div class="bar load" :style="styles.load" :title="`Load: ${load.toFixed(1)} A`"/>
    <div class="bar reserved" :style="styles.reserved" :title="`Reserved: ${reserved.toFixed(1)} A`"/>
    <div class="bar free" :style="styles.free" :title="`Free: ${free.toFixed(1)} A`"/>
    <div class="bar notified" :style="styles.notified" :title="`Notified: ${notified.toFixed(1)} A`"/>
  </div>
</template>

<script>
import {linearConversion} from '../../util/convert.js';

export default {
  name: 'PowerSupplyBar',
  props: {
    rating: Number,
    load: Number,
    capacity: Number,
    free: Number,
    notified: Number,
    warnAt: {
      type: Number,
      default: 75
    },
    warnColor: {
      default: 'warning'
    },
    errorAt: {
      type: Number,
      default: 95
    },
    errorColor: {
      default: 'error'
    }
  },
  computed: {
    reserved() {
      return this.capacity - this.free;
    },
    styles() {
      const styles = {};
      styles.rating = {
        width: 100
      };
      styles.load = {
        width: linearConversion(0, this.rating, 0, 100, this.load)
      };
      styles.reserved = {
        left: styles.load.width,
        width: linearConversion(0, this.rating, 0, 100, this.reserved)
      };
      styles.free = {
        left: styles.reserved.left + styles.reserved.width,
        width: Math.max(0, linearConversion(0, this.rating, 0, 100, this.free))
      };
      styles.notified = {
        left: styles.free.left,
        width: linearConversion(0, this.rating, 0, 100, this.notified)
      };

      // power usage can exceed the rating, translate things to fit in the actual max
      const max = linearConversion(0, this.rating, 0, 100,
          Math.max(this.rating, this.rating - this.free + this.notified));
      const props = ['left', 'right', 'width'];
      for (const style of Object.values(styles)) {
        for (const prop of props) {
          if (style[prop]) {
            style[prop] = linearConversion(0, max, 0, 100, style[prop]) + '%';
          }
        }
      }
      return styles;
    }
  }
};
</script>

<style scoped>
.power-supply-bar {
  position: relative;
  height: 24px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #0004;
}

.power-supply-bar > .bar {
  position: absolute;
  top: 25%;
  height: 50%;
  left: 0;
  width: 0;
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.power-supply-bar .load {
  background: #F44336;
}

.power-supply-bar .free {
  background: #4CAF50;
}

.power-supply-bar .reserved {
  background: #0002;
}

.power-supply-bar .rating {
  border-radius: 3px;
  background: #0001;
  top: 0;
  height: 100%;
}

.power-supply-bar .notified {
  background: #F4433670;
  top: 50%;
  height: 25%;
}
</style>
