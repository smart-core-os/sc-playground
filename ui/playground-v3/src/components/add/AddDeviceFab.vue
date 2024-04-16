<template>
  <v-container class="align-centerfill-height mx-auto">
      <v-row align="center" justify="center" class="align-content-center">
        <v-col cols="auto">
          <v-btn @click="addTrait" :elevation="labelElevation">
            Add Trait
            <v-icon>mdi-view-grid-plus</v-icon>
          </v-btn>
        </v-col>
        <v-col cols="auto">
        <v-btn @click="addRemote" :elevation="labelElevation">
          Add Network Device
          <v-icon>mdi-plus-network-outline</v-icon>
        </v-btn>
        </v-col>
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
