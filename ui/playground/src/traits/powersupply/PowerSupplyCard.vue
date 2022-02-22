<template>
  <v-card width="400">
    <v-card-title>
      <span>{{ deviceId }}</span>
      <v-spacer/>
      <span title="Free / Rating" v-if="capacity">
        {{ capacity.free.toFixed(1) }} free of {{ capacity.rating.toFixed(1) }}&nbsp;A
      </span>
    </v-card-title>
    <v-card-text>
      <power-supply-bar v-bind="capacity" v-if="capacity"/>
      <power-supply-draw-notification-list v-if="drawNotifications.length > 0"
                                           class="mx-n4" :items="drawNotifications"
                                           @click:remove="removeDrawNotification"/>
    </v-card-text>
    <v-card-actions class="align-center">
      <v-text-field label="Max draw" suffix="Amps" outlined class="mr-2" hide-details dense type="number"
                    v-model.number="draw.max"/>
      <v-text-field label="Min draw" suffix="Amps" outlined class="mr-2" hide-details dense type="number"
                    v-model.number="draw.min"/>
      <v-btn @click="addDrawNotification">Notify</v-btn>
    </v-card-actions>
    <v-card-actions class="align-center pt-0 pb-4">
      <v-checkbox v-model="draw.force" class="ms-2" label="Force" hide-details dense/>
    </v-card-actions>
    <v-expand-transition>
      <v-card-text v-if="draw.message" key="notifyMsg">
        {{ draw.message }}
      </v-card-text>
    </v-expand-transition>
    <v-divider/>
    <v-card-subtitle>Adjust Device</v-card-subtitle>
    <v-card-text>
      <power-supply-settings-editor v-if="settings" v-bind.sync="settings" class="settings-editor"/>
    </v-card-text>
  </v-card>
</template>

<script>

import {PowerSupplyApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/power_supply_grpc_web_pb.js';
import {
  CreateDrawNotificationRequest,
  DeleteDrawNotificationRequest,
  DrawNotification,
  GetPowerCapacityRequest,
  ListDrawNotificationsRequest,
  PullDrawNotificationsRequest,
  PullPowerCapacityRequest
} from '@smart-core-os/sc-api-grpc-web/traits/power_supply_pb.js';
import {MemorySettingsApiPromiseClient} from '@sc-playground/gen/trait/powersupply/memory_settings_grpc_web_pb.js';
import {
  GetMemorySettingsReq,
  MemorySettings,
  PullMemorySettingsReq,
  UpdateMemorySettingsReq
} from '@sc-playground/gen/trait/powersupply/memory_settings_pb.js';
import {FieldMask} from 'google-protobuf/google/protobuf/field_mask_pb.js';
import PowerSupplyBar from './PowerSupplyBar.vue';
import PowerSupplySettingsEditor from './PowerSupplySettingsEditor.vue';
import {grpcWebEndpoint} from '../../util/api.js';
import PowerSupplyDrawNotificationList from './PowerSupplyDrawNotificationList.vue';
import {Duration} from 'google-protobuf/google/protobuf/duration_pb.js';
import Vue from 'vue';

export default {
  name: 'PowerSupplyCard',
  components: {PowerSupplyDrawNotificationList, PowerSupplySettingsEditor, PowerSupplyBar},
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
      drawNotificationStream: null,
      drawNotifications: [],
      settingsStream: null,
      settings: null,
      serverSettings: null,

      draw: {
        working: false,
        message: null,
        clearMessageHandle: 0,
        max: 20,
        min: 0,
        force: false,
        durationSec: 30
      }
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
      const serverEndpoint = await grpcWebEndpoint();
      const api = new PowerSupplyApiPromiseClient(serverEndpoint, null, null);
      // capacity
      if (this.capacityStream) this.capacityStream.cancel();
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
      // draw notifications
      if (this.drawNotificationStream) this.drawNotificationStream.cancel();
      const notificationListReq = new ListDrawNotificationsRequest().setName(this.deviceId);
      let notificationRes = await api.listDrawNotifications(notificationListReq);
      this.drawNotifications = notificationRes.getDrawNotificationsList().map(n => n.toObject());
      // get all pages
      while (notificationRes.getNextPageToken()) {
        notificationListReq.setPageToken(notificationRes.getNextPageToken());
        notificationRes = await api.listDrawNotifications(notificationListReq);
        this.drawNotifications.push(...notificationRes.getDrawNotificationsList().map(n => n.toObject()));
      }
      this.drawNotificationStream = api.pullDrawNotifications(
          new PullDrawNotificationsRequest().setName(this.deviceId));
      this.drawNotificationStream.on('data', res => {
        /** @type {PullDrawNotificationsResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          // We handle change by looking at new and old values instead of the change type.
          // There are three cases
          // Old value only -> delete the index
          // New value only -> push to the array
          // Old and new -> replace the value at the same index
          //
          // Having an old value that can't be found it treated as if there were no old value
          let index = -1;
          if (change.hasOldValue()) {
            const oldId = change.getOldValue().getId();
            index = this.drawNotifications.findIndex(n => n.id === oldId);
          }
          if (change.hasNewValue()) {
            if (index >= 0) {
              Vue.set(this.drawNotifications, index, change.getNewValue().toObject());
            } else {
              this.drawNotifications.push(change.getNewValue().toObject());
            }
          } else {
            if (index >= 0) {
              this.drawNotifications.splice(index, 1);
            }
          }
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
    async addDrawNotification() {
      this.draw.working = true;
      this.draw.message = null;
      clearTimeout(this.draw.clearMessageHandle);
      try {
        const n = new DrawNotification();
        if (this.draw.min) n.setMinDraw(this.draw.min);
        n.setMaxDraw(this.draw.max);
        n.setRampDuration(new Duration().setSeconds(this.draw.durationSec));
        if (this.draw.force) n.setForce(this.draw.force);
        const req = new CreateDrawNotificationRequest()
            .setName(this.deviceId)
            .setDrawNotification(n);

        const serverEndpoint = await grpcWebEndpoint();
        const api = new PowerSupplyApiPromiseClient(serverEndpoint, null, null);
        const res = await api.createDrawNotification(req);
        if (res.getMaxDraw() < n.getMaxDraw()) {
          this.draw.message = `Unable to reserve all power, ${res.getMaxDraw()} A reserved`;
          this.draw.clearMessageHandle = setTimeout(() => {
            this.draw.message = null;
          }, 10 * 1000)
        }
      } finally {
        this.draw.working = false;
      }
    },
    async removeDrawNotification(id) {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new PowerSupplyApiPromiseClient(serverEndpoint, null, null);
      await api.deleteDrawNotification(new DeleteDrawNotificationRequest().setName(this.deviceId).setId(id));
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
