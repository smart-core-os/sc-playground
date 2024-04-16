<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <template #title-append>
      <charge-icon :energy-level="resources.energyLevel.value"/>
    </template>
    <v-card-text>
      <template v-if="resources.energyLevel.value">
        <!-- todo: display energyLevels better       -->
        {{ resources.energyLevel.value }}
      </template>
    </v-card-text>
  </trait-card>
</template>

<script>
import TraitCard from '../../components/TraitCard.vue';
import {pullEnergyLevel} from './energy-storage.js';
import ChargeIcon from './ChargeIcon.vue';

export default {
  name: 'EnergyStorageCard',
  components: {ChargeIcon, TraitCard},
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
        .catch(err => console.error('during pull', err));
  },
  beforeDestroy() {
    for (const resource of Object.values(this.resources)) {
      if (resource.stream) resource.stream.cancel();
    }
  },
  methods: {
    async pull() {
      // EnergyLevel resource
      this.resources.energyLevel = await pullEnergyLevel(this.deviceId, this.resources.energyLevel);
    },
    log(...args) {
      console.debug(...args);
    }
  }
};
</script>

<style scoped>
</style>
