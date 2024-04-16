<template>
  <v-icon v-bind="$attrs">{{ iconStr }}</v-icon>
</template>

<script>
export default {
  name: 'ChargeIcon',
  props: {
    /** @type {import('vue').PropType<EnergyLevel.AsObject>} */
    energyLevel: [Object]
  },
  computed: {
    iconStr() {
      if (this.charging) {

        let prefix = 'mdi-battery';
        if (this.charging) prefix += '-charging';
        const keys = [10, 20, 30, 40, 50, 60, 70, 80, 90];
        const val = this.progress;
        if (typeof val === 'number') {
          for (const key of keys) {
            if (val <= key) {
              return `${prefix}-${key}`;
            }
          }
        }
        return prefix;
      }
      if (this.pluggedIn) {
        return 'mdi-power-plug';
      }
      return 'mdi-power-plug-off';
    },

    pluggedIn() {
      return Boolean(this.energyLevel?.pluggedIn);
    },
    charging() {
      return Boolean(this.energyLevel?.charge);
    },
    progress() {
      const el = /** @type {EnergyLevel.AsObject} */ this.energyLevel;
      return el.quantity?.percentage;
    }
  }
};
</script>

<style scoped>

</style>
