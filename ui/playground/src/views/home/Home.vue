<template>
  <v-container class="card-grid">
    <power-supply-card v-for="deviceId in devices" :device-id="deviceId" :key="deviceId"/>
    <v-btn @click="addDevice" key="btn">Add {{ nextDeviceId }}</v-btn>
  </v-container>
</template>

<script>

import PowerSupplyCard from '../../traits/powersupply/PowerSupplyCard.vue';

export default {
  name: 'Home',
  components: {PowerSupplyCard},
  data() {
    return {
      devices: ['POW-001']
    };
  },
  computed: {
    nextDeviceId() {
      const nextIdNum = (this.devices.length + 1).toString().padStart(3, '0');
      return `POW-${nextIdNum}`;
    }
  },
  methods: {
    addDevice() {
      this.devices.push(this.nextDeviceId);
    }
  }
};
</script>

<style scoped>
.card-grid {
  display: grid;
  grid-gap: 12px;
  grid-template-columns: repeat(auto-fill, 400px);
  grid-auto-columns: 400px;
  justify-content: center;
}
</style>
