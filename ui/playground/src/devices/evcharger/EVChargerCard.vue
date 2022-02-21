<template>
  <v-card>
    <v-card-text>This is an EV charger</v-card-text>
    <v-card-subtitle>Manual mode control</v-card-subtitle>
    <v-card-text>
      <get-electric-mode-chooser :device-id="deviceId"/>
    </v-card-text>
  </v-card>
</template>

<script>
import ElectricModeChooser from '../../traits/electric/ElectricModeChooser.vue';
import GetElectricModeChooser from '../../traits/electric/GetElectricModeChooser.vue';

export default {
  name: 'EVChargerCard',
  components: {GetElectricModeChooser, ElectricModeChooser},
  props: {
    child: [Object],
    metadata: [Object]
  },
  computed: {
    deviceId() {
      return this.child.name;
    }
  }
};
</script>

<style scoped>

</style>
