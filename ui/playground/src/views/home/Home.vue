<template>
  <v-container class="card-grid">
    <power-supply-card v-for="deviceId in devices" :device-id="deviceId" :key="deviceId"/>

    <form @submit.prevent="addDevice">
      <v-card>
        <v-card-title>Add or expose another device</v-card-title>
        <v-card-text>
          <v-text-field label="Device ID" v-model="nextDeviceId"/>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn elevation="0" color="success" type="submit">Add</v-btn>
        </v-card-actions>
      </v-card>
    </form>
  </v-container>
</template>

<script>

import PowerSupplyCard from '../../traits/powersupply/PowerSupplyCard.vue';

export default {
  name: 'Home',
  components: {PowerSupplyCard},
  data() {
    return {
      devices: ['POW-001'],
      manualNextDeviceId: null
    };
  },
  computed: {
    nextDeviceId: {
      get() {
        if (this.manualNextDeviceId != null) {
          return this.manualNextDeviceId;
        }
        const nextIdNum = (this.devices.length + 1).toString().padStart(3, '0');
        return `POW-${nextIdNum}`;
      },
      set(v) {
        this.manualNextDeviceId = v;
      }
    }
  },
  methods: {
    addDevice() {
      this.devices.push(this.nextDeviceId);
      this.manualNextDeviceId = null;
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
