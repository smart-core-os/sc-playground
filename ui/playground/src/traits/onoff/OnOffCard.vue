<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <v-card-text>
      State: {{ state }}
    </v-card-text>
    <v-card-actions>
      <v-btn @click="doOn" :loading="Boolean(resources.onOff.loading)" depressed>On</v-btn>
      <v-btn @click="doOff" :loading="Boolean(resources.onOff.loading)" depressed>Off</v-btn>
      <v-btn @click="doReboot" :loading="Boolean(resources.onOff.loading)" depressed>Reboot</v-btn>
    </v-card-actions>
  </trait-card>
</template>

<script>
import TraitCard from '../../components/TraitCard.vue';
import {OnOff} from "@smart-core-os/sc-api-grpc-web/traits/on_off_pb.js";
import {pullOnOff, updateOnOff} from "./on-off.js";

export default {
  name: 'OnOffCard',
  components: {TraitCard},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        onOff: {
          /** @type {OnOff.AsObject} */
          value: null,
          stream: null,
          loading: 0
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
    state() {
      if (!this.resources.onOff.value) {
        return "LOADING";
      }
      switch (this.resources.onOff.value.state) {
        case OnOff.State.UNKNOWN:
          return "UNKNOWN";
        case OnOff.State.ON:
          return "ON";
        case OnOff.State.OFF:
          return "OFF";
        default:
          return "< unsupported >";
      }
    }
  },
  methods: {
    async pull() {
      // EnergyLevel resource
      this.resources.onOff = await pullOnOff(this.deviceId, this.resources.onOff);
    },
    log(...args) {
      console.debug(...args);
    },
    doOn() {
      return updateOnOff(this.deviceId, {state: OnOff.State.ON}, this.resources.onOff);
    },
    doOff() {
      return updateOnOff(this.deviceId, {state: OnOff.State.OFF}, this.resources.onOff);
    },
    async doReboot() {
      try {
        this.resources.onOff.loading++;
        await this.doOff();
        await new Promise(resolve => setTimeout(resolve, 5000));
        await this.doOn()
      } finally {
        this.resources.onOff.loading--;
      }
    }
  }
};
</script>

<style scoped>
</style>
