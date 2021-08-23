<template>
  <v-list>
    <DeviceListItem v-for="device in devices" :device="device" :key="device.getName()"/>
  </v-list>
</template>

<script>
import {Device} from '@smart-core-os/sc-api-grpc-web/info/info_pb';
import DeviceListItem from '../../components/DeviceListItem.vue';

export default {
  name: 'DevicesView',
  components: {DeviceListItem},
  data() {
    return {
      /** @type {Device[]} */
      devices: [
        new Device()
            .setName('ABC-012'),
        new Device()
            .setName('XYZ-001')
      ]
    };
  }
};
</script>

<style scoped>

</style>
