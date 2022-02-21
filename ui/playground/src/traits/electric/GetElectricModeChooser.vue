<template>
  <electric-mode-chooser :mode="activeMode" @update:mode="setMode" :modes="modes"
                         :label="`Active mode (of ${modes.length})`"/>
</template>

<script>

import {grpcWebEndpoint} from '../../util/api.js';
import {ElectricApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/electric_grpc_web_pb.js';
import {
  ElectricMode,
  GetActiveModeRequest,
  ListModesRequest,
  PullActiveModeRequest,
  PullModesRequest,
  UpdateActiveModeRequest
} from '@smart-core-os/sc-api-grpc-web/traits/electric_pb.js';
import Vue from 'vue';
import ElectricModeChooser from './ElectricModeChooser.vue';

export default {
  name: 'GetElectricModeChooser',
  components: {ElectricModeChooser},
  props: {
    deviceId: [String]
  },
  data() {
    return {
      resources: {
        activeMode: {
          /** @type {ElectricMode.AsObject} */
          value: null,
          stream: null
        },
        modes: {
          /** @type {Object<string,ElectricMode.AsObject>} */
          value: {},
          stream: null
        }
      }
    };
  },
  mounted() {
    this.pull()
        .catch(err => console.error('during pull', err));
  },
  beforeDestroy() {
    for (const resource of Object.values(this.resources)) {
      if (resource.stream) resource.stream.cancel();
    }
  },
  computed: {
    activeMode() {
      return this.resources.activeMode.value;
    },
    modes() {
      return Object.values(this.resources.modes.value);
    }
  },
  methods: {
    async pull() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ElectricApiPromiseClient(serverEndpoint, null, null);

      // activeMode resource
      const activeModeResource = this.resources.activeMode;
      if (activeModeResource.stream) activeModeResource.stream.cancel();
      const activeModePb = await api.getActiveMode(new GetActiveModeRequest().setName(this.deviceId));
      activeModeResource.value = activeModePb.toObject();
      const activeModeStream = api.pullActiveMode(new PullActiveModeRequest().setName(this.deviceId));
      activeModeResource.stream = activeModeStream;
      activeModeStream.on('data', res => {
        /** @type {PullActiveModeResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const value = change.getActiveMode();
          activeModeResource.value = value.toObject();
        }
      });

      // modes resource
      const modesResource = this.resources.modes;
      if (modesResource.stream) modesResource.stream.cancel();
      let modesPb = await api.listModes(new ListModesRequest().setName(this.deviceId));
      while (true) {
        for (const mode of modesPb.getModesList()) {
          Vue.set(modesResource.value, mode.getId(), mode.toObject());
        }
        if (!modesPb.getNextPageToken()) {
          break;
        }
        modesPb = await api.listModes(new ListModesRequest().setName(this.deviceId)
            .setPageToken(modesPb.getNextPageToken()));
      }
      const modesStream = api.pullModes(new PullModesRequest().setName(this.deviceId));
      modesResource.stream = modesStream;
      modesStream.on('data', res => {
        /** @type {PullModesResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const value = change.getNewValue();
          if (!value) {
            // delete
            const oldId = change.getOldValue()?.getId();
            if (oldId) {
              Vue.delete(modesResource.value, oldId);
            }
          } else {
            Vue.set(modesResource.value, value.getId(), value.toObject());
          }
        }
      });
    },
    log(...args) {
      console.debug(...args);
    },
    async setMode(mode) {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ElectricApiPromiseClient(serverEndpoint, null, null);
      await api.updateActiveMode(new UpdateActiveModeRequest().setName(this.deviceId)
          .setActiveMode(new ElectricMode().setId(mode.id)));
      // todo: error handling
    }
  }
};
</script>

<style scoped>
</style>
