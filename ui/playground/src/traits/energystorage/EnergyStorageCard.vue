<template>
  <trait-card :device-id="deviceId" :trait="trait" :icon="iconStr">
    <v-card-text>
      <template v-if="resources.energyLevel.value">
        <!-- todo: display energyLevels better       -->
        {{ resources.energyLevel.value }}
      </template>
    </v-card-text>
  </trait-card>
</template>

<script>

import {grpcWebEndpoint} from '../../util/api.js';
import {EnergyStorageApiPromiseClient} from "@smart-core-os/sc-api-grpc-web/traits/energy_storage_grpc_web_pb.js";
import {
  GetEnergyLevelRequest,
  PullEnergyLevelRequest
} from "@smart-core-os/sc-api-grpc-web/traits/energy_storage_pb.js";
import TraitCard from "../../components/TraitCard.vue";

export default {
  name: 'EnergyStorageCard',
  components: {TraitCard},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        energyLevel: {
          /** @type {EnergyLevel.AsObject} */
          value: null,
          stream: null
        }
      }
    };
  },
  mounted() {
    this.pull()
        .catch(err => console.error('during pull', err))
  },
  beforeDestroy() {
    for (const resource of Object.values(this.resources)) {
      if (resource.stream) resource.stream.cancel();
    }
  },
  computed: {
    iconStr() {
      if (this.charging) {
        return 'mdi-battery';
      }
      if (this.pluggedIn) {
        return 'mdi-power-plug';
      }
      return 'mdi-power-plug-off';
    },
    pluggedIn() {
      return Boolean(this.resources.energyLevel.value?.pluggedIn);
    },
    charging() {
      return false;
    }
  },
  methods: {
    async pull() {
      const serverEndpoint = await grpcWebEndpoint();

      // EnergyLevel resource
      const energyLevelApi = new EnergyStorageApiPromiseClient(serverEndpoint, null, null);
      const energyLevelResource = this.resources.energyLevel;
      if (energyLevelResource.stream) energyLevelResource.stream.cancel();
      const energyLevelPb = await energyLevelApi.getEnergyLevel(new GetEnergyLevelRequest().setName(this.deviceId));
      energyLevelResource.value = energyLevelPb.toObject();
      const energyLevelStream = energyLevelApi.pullEnergyLevel(new PullEnergyLevelRequest().setName(this.deviceId));
      energyLevelResource.stream = energyLevelStream;
      energyLevelStream.on('data', res => {
        /** @type {PullEnergyLevelResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const value = change.getEnergyLevel();
          energyLevelResource.value = value.toObject();
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
</style>
