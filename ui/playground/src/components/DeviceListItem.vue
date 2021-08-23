<template>
  <v-list-item>
    <v-list-item-title>{{ title }}</v-list-item-title>
  </v-list-item>
</template>

<script>
import {Device} from '@smart-core-os/sc-api-grpc-web/info/info_pb';

export default {
  name: 'DeviceListItem',
  props: {
    device: {
      type: Device
    }
  },
  computed: {
    title() {
      return this.device.getName();
    }
  }
};
</script>

<style scoped>

</style>
