<template>
  <v-card>
    <v-card-title>{{ title }}</v-card-title>
    <v-card-subtitle>{{ child.name }}</v-card-subtitle>
    <v-divider/>
    <v-card-actions>
      <v-btn text @click="asyncAction(plugIn)" :disabled="pluggedIn">Plug in</v-btn>
      <v-btn text @click="asyncAction(charge)" :disabled="!pluggedIn || charging">Begin charge</v-btn>
      <v-btn text @click="asyncAction(unplug)" :disabled="!pluggedIn">Unplug</v-btn>
    </v-card-actions>
    <v-card-subtitle>Override electric mode</v-card-subtitle>
    <v-card-text>
      <get-electric-mode-chooser :device-id="deviceId"/>
    </v-card-text>
  </v-card>
</template>

<script>
import ElectricModeChooser from '../../traits/electric/ElectricModeChooser.vue';
import GetElectricModeChooser from '../../traits/electric/GetElectricModeChooser.vue';
import {localName} from '../../util/names.js';
import {grpcWebEndpoint} from '../../util/api.js';
import {EVChargerApiPromiseClient} from '@sc-playground/gen/pkg/device/evcharger/evcharger_grpc_web_pb.js';
import {
  ChargeStartRequest,
  PlugInEvent,
  PlugInRequest,
  UnplugRequest
} from '@sc-playground/gen/pkg/device/evcharger/evcharger_pb.js';
import {EnergyLevel} from '@smart-core-os/sc-api-grpc-web/traits/energy_storage_pb.js';
import {ElectricMode} from '@smart-core-os/sc-api-grpc-web/traits/electric_pb.js';
import durationpb from 'google-protobuf/google/protobuf/duration_pb.js';
import {pullEnergyLevel} from '../../traits/energystorage/energy-storage.js';

export default {
  name: 'EVChargerCard',
  components: {GetElectricModeChooser, ElectricModeChooser},
  props: {
    child: [Object],
    metadata: [Object]
  },
  data() {
    return {
      resources: {
        energyLevel: {
          /** @type {EnergyLevel.AsObject} */
          value: null,
          stream: null
        }
      },
      action: {
        loading: false,
        error: null
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
    deviceId() {
      return this.child.name;
    },
    title() {
      const mdTitle = this.mdTitle;
      return mdTitle || localName(this.child?.name);
    },
    mdTitle() {
      return this.metadata?.appearance?.title;
    },
    pluggedIn() {
      return Boolean(this.resources.energyLevel.value?.pluggedIn);
    },
    charging() {
      return Boolean(this.resources.energyLevel.value?.charge);
    }
  },
  methods: {
    async pull() {
      this.resources.energyLevel = await pullEnergyLevel(this.deviceId, this.resources.energyLevel);
    },
    async asyncAction(action) {
      if (typeof action === 'string') {
        action = () => this[action]();
      }

      try {
        this.action.loading = true;
        await action.call(this);
        this.action.error = null;
      } catch (e) {
        console.error(e);
        this.action.error = e;
      } finally {
        this.action.loading = false;
      }
    },
    async plugIn() {
      // helps with duration calculations
      const second = 1;
      const minute = 60 * second;
      const hour = 60 * minute;
      const day = 24 * hour;

      const serverEndpoint = await grpcWebEndpoint();
      const api = new EVChargerApiPromiseClient(serverEndpoint, null, null);
      const full = new EnergyLevel.Quantity()
          .setEnergyKwh(120)
          .setDistanceKm(600);

      // Tesla (US) charge levels from https://www.quora.com/How-many-amps-does-a-Tesla-charger-pull
      const level1 = new PlugInEvent.ChargeMode()
          .setId('level1')
          .setTitle('Home Charge')
          .setDescription('Level 1 charging mode')
          .setSegmentsList([
            new ElectricMode.Segment().setMagnitude(13).setLength(new durationpb.Duration().setSeconds(2 * day))
          ]);
      const level2 = new PlugInEvent.ChargeMode()
          .setId('level2')
          .setTitle('Destination Charge')
          .setDescription('Level 2 charging mode')
          .setSegmentsList([
            new ElectricMode.Segment().setMagnitude(48).setLength(new durationpb.Duration().setSeconds(2 * hour))
          ]);
      const level3 = new PlugInEvent.ChargeMode()
          .setId('level3')
          .setTitle('Fast Charge')
          .setDescription('Level 3 charging mode')
          .setSegmentsList([
            new ElectricMode.Segment().setMagnitude(250).setLength(new durationpb.Duration().setSeconds(30 * minute))
          ]);
      await api.plugIn(new PlugInRequest().setName(this.deviceId)
          .setEvent(new PlugInEvent()
              .setFull(full)
              .setModesList([level2, level1, level3])));
    },
    async charge() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new EVChargerApiPromiseClient(serverEndpoint, null, null);
      await api.chargeStart(new ChargeStartRequest().setName(this.deviceId));
    },
    async unplug() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new EVChargerApiPromiseClient(serverEndpoint, null, null);
      await api.unplug(new UnplugRequest().setName(this.deviceId));
    }
  }
};
</script>

<style scoped>

</style>
