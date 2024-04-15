<template>
  <v-container>

    <v-row>
      <v-btn fab small @click="addTrait" :elevation="labelElevation">
        Add Trait
        <v-icon>mdi-view-grid-plus</v-icon>
      </v-btn>
    </v-row>

    <v-row>
      <v-btn fab small @click="addRemote" :elevation="labelElevation">
        Add Network Device
        <v-icon>mdi-plus-network-outline</v-icon>
      </v-btn>
    </v-row>
    <v-dialog v-model="dialogOpen" width="400">
      <add-trait-card v-if="dialogFlavor === 'trait'" @done="dialogOpen = false"/>
      <add-remote-device-card v-if="dialogFlavor === 'remote'" @done="dialogOpen = false"/>
    </v-dialog>
  </v-container>
</template>

<script>
import AddTraitCard from "./AddTraitCard.vue";
import AddRemoteDeviceCard from "./AddRemoteDeviceCard.vue";

export default {
  name: "AddDeviceFab",
  components: {AddRemoteDeviceCard, AddTraitCard},
  data() {
    return {
      expanded: false,
      dialogOpen: false,
      labelElevation: 3,
      dialogFlavor: 'trait'
    }
  },
  methods: {
    addTrait() {
      this.dialogFlavor = 'trait';
      this.dialogOpen = true;
    },
    addRemote() {
      this.dialogFlavor = 'remote';
      this.dialogOpen = true;
    }
  }
}
</script>

<style scoped>
.label {
  position: absolute;
  right: 55px;
  text-transform: initial;
}
</style>
