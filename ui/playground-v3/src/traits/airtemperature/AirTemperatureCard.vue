<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <template #title-append>
      <v-icon right>mdi-thermometer</v-icon>
    </template>
    <v-card-text>
      Mode: {{ mode }}<br/>
      Temperature: {{ temperature }}°C<br/>
      Humidity: {{ humidity }}%
    </v-card-text>
  </trait-card>
</template>

<script>
import {AirTemperature} from "@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb";
import {pullAirTemperature} from "./airtemp";
import TraitCard from "../../components/TraitCard.vue";

export default {
  name: "AirTemperatureCard",
  components:{TraitCard},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        airTemperature: {
          /** @type {AirTemperature.AsObject} */
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
    mode() {
      if (!this.resources.airTemperature.value) {
        return "LOADING";
      }
      switch (this.resources.airTemperature.value.mode) {
        case AirTemperature.Mode.MODE_UNSPECIFIED:
          return "UNKNOWN";
        case AirTemperature.Mode.AUTO:
          return "AUTO";
        case AirTemperature.Mode.ON:
          return "ON";
        case AirTemperature.Mode.OFF:
          return "OFF";
        default:
          return `< unsupported: ${this.resources.airTemperature.value.mode}>`;
      }
    },
    temperature() {
      if (!this.resources.airTemperature.value) {
        return "...";
      }
      return this.resources.airTemperature.value.ambientTemperature.valueCelsius.toFixed(1);
    },
    humidity() {
      if (!this.resources.airTemperature.value) {
        return "...";
      }
      return (this.resources.airTemperature.value.ambientHumidity*100).toFixed(1);
    }
  },
  methods: {
    async pull() {
      // Occupancy resource
      this.resources.airTemperature = await pullAirTemperature(this.deviceId, this.resources.airTemperature);
    },
    log(...args) {
      console.debug(...args);
    }
  }
}
</script>

<style scoped>

</style>
