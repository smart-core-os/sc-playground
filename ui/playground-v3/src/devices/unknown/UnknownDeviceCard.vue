<template>
  <v-card width="400">
    <v-card-title>{{ child.name }}</v-card-title>
    <v-card-subtitle>Device type {{ deviceType }}</v-card-subtitle>
    <v-card-text>
      <metadata-panel :metadata="metadata"/>
    </v-card-text>
  </v-card>
</template>

<script>
import MetadataPanel from '../../traits/metadata/MetadataPanel.vue';

export default {
  name: 'UnknownDeviceCard',
  components: {MetadataPanel},
  props: {
    child: [Object],
    metadata: [Object]
  },
  computed: {
    deviceType() {
      const entry = this.metadata?.moreMap?.find(([k]) => k === 'scos.playground.device-type');
      return entry ? entry[1] : '<unknown>';
    }
  }
};
</script>

<style scoped>

</style>
