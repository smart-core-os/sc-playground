<template>
  <v-card width="400">
    <v-card-title>
      <span>{{ deviceId }}</span>
      <v-spacer/>
      <span title="Free / Rating" v-if="capacity">
        {{ capacity.free.toFixed(1) }} free of {{ capacity.rating.toFixed(1) }} A
      </span>
    </v-card-title>
    <v-card-text>
      <power-supply-bar v-bind="capacity" v-if="capacity"/>
    </v-card-text>
    <v-card-subtitle>Adjust Device</v-card-subtitle>
    <v-card-text>
      <power-supply-settings-editor v-if="settings" v-bind.sync="settings" class="settings-editor"/>
    </v-card-text>
  </v-card>
</template>

<script>

import {PowerSupplyApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/power_supply_grpc_web_pb.js';
import {
  GetPowerCapacityRequest,
  PullPowerCapacityRequest
} from '@smart-core-os/sc-api-grpc-web/traits/power_supply_pb.js';
import {MemorySettingsApiPromiseClient} from 'gen/trait/powersupply/memory_settings_grpc_web_pb.js';
import {
  GetMemorySettingsReq,
  MemorySettings,
  PullMemorySettingsReq,
  UpdateMemorySettingsReq
} from 'gen/trait/powersupply/memory_settings_pb.js';
import {FieldMask} from 'google-protobuf/google/protobuf/field_mask_pb.js';
import PowerSupplyBar from './PowerSupplyBar.vue';
import PowerSupplySettingsEditor from './PowerSupplySettingsEditor.vue';
import {grpcWebEndpoint} from '../../util/api.js';

export default {
  name: 'PowerSupplyCard',
  components: {PowerSupplySettingsEditor, PowerSupplyBar},
  props: {
    deviceId: {
      type: String,
      default: 'POW-001'
    }
  },
  data() {
    return {
      capacityStream: null,
      capacity: null,
      settingsStream: null,
      settings: null,
      serverSettings: null
    };
  },
  mounted() {
    this.pull();
  },
  beforeDestroy() {
    if (this.capacityStream) this.capacityStream.cancel();
  },
  watch: {
    settings: {
      deep: true,
      async handler(v) {
        const old = this.serverSettings || {};
        console.debug('settings updated:', old, '->', v);
        const serverEndpoint = await grpcWebEndpoint();
        const settingsApi = new MemorySettingsApiPromiseClient(serverEndpoint, null, null);
        const changed = [];
        const newSettings = new MemorySettings();
        if (v.rating !== old.rating) {
          changed.push('rating');
          newSettings.setRating(v.rating);
        }
        if (v.voltage !== old.voltage) {
          changed.push('voltage');
          newSettings.setVoltage(v.voltage);
        }
        if (v.load !== old.load) {
          changed.push('load');
          newSettings.setLoad(v.load);
        }
        if (v.reserved !== old.reserved) {
          changed.push('reserved');
          newSettings.setReserved(v.reserved);
        }
        if (changed.length > 0) {
          const req = new UpdateMemorySettingsReq()
              .setName(this.deviceId)
              .setSettings(newSettings)
              .setUpdateMask(new FieldMask().setPathsList(changed));
          console.debug('Updating settings:', req.toObject());
          const settingsRes = await settingsApi.updateSettings(req, undefined);
          this.settings = settingsRes.toObject();
          this.serverSettings = settingsRes.toObject();
        }
      }
    }
  },
  methods: {
    async pull() {
      if (this.capacityStream) this.capacityStream.cancel();
      const serverEndpoint = await grpcWebEndpoint();
      const api = new PowerSupplyApiPromiseClient(serverEndpoint, null, null);
      const capacityRes = await api.getPowerCapacity(new GetPowerCapacityRequest().setName(this.deviceId));
      this.capacity = capacityRes.toObject();
      const capacityStream = api.pullPowerCapacity(new PullPowerCapacityRequest().setName(this.deviceId));
      this.capacityStream = capacityStream;
      capacityStream.on('data', res => {
        /** @type {PullPowerCapacityResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const capacity = change.getAvailableCapacity();
          this.capacity = capacity.toObject();
        }
      });

      if (this.settingsStream) this.settingsStream.cancel();
      const settingsApi = new MemorySettingsApiPromiseClient(serverEndpoint, null, null);
      const settingsRes = await settingsApi.getSettings(new GetMemorySettingsReq().setName(this.deviceId), undefined);
      this.settings = settingsRes.toObject();
      this.serverSettings = settingsRes.toObject();

      const settingsStream = settingsApi.pullSettings(new PullMemorySettingsReq().setName(this.deviceId), undefined);
      this.settingsStream = settingsStream;
      settingsStream.on('data', res => {
        /** @type {PullMemorySettingsRes.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const settings = change.getSettings();
          this.settings = settings.toObject();
          this.serverSettings = settings.toObject();
        }
      });
    },
    log(...args) {
      console.debug(...args);
    }
  }
};
</script>

<style scoped>
.settings-editor >>> label.v-label {
  min-width: 5em;
}
</style>
