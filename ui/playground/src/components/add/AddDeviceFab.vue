<template>
  <div>
    <v-speed-dial bottom right fixed v-model="expanded" transition="slide-y-reverse-transition">
      <template #activator>
        <v-btn fab v-model="expanded">
          <v-icon v-if="expanded">mdi-close</v-icon>
          <v-icon v-else>mdi-plus</v-icon>
        </v-btn>
      </template>
      <v-btn fab small @click="addTrait" :elevation="labelElevation">
        <v-icon>mdi-view-grid-plus</v-icon>
        <v-sheet class="label px-4 py-2" rounded :elevation="labelElevation">Trait</v-sheet>
      </v-btn>
      <v-btn fab small :elevation="labelElevation" disabled>
        <v-icon>mdi-plus-box-outline</v-icon>
        <v-sheet class="label px-4 py-2" rounded :elevation="labelElevation">Virtual Device</v-sheet>
      </v-btn>
      <v-btn fab small @click="addRemote" :elevation="labelElevation">
        <v-icon>mdi-plus-network-outline</v-icon>
        <v-sheet class="label px-4 py-2" rounded :elevation="labelElevation">Network Device</v-sheet>
      </v-btn>
    </v-speed-dial>
    <v-dialog v-model="dialogOpen" width="400">
      <add-trait-card v-if="dialogFlavor === 'trait'" @done="dialogOpen = false"/>
      <add-remote-device-card v-if="dialogFlavor === 'remote'" @done="dialogOpen = false"/>
    </v-dialog>
  </div>
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
