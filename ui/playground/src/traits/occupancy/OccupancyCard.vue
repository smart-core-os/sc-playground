<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <template #title-append>
      <span>{{ count }}</span>
      <v-icon right>mdi-account</v-icon>
    </template>
    <v-card-text>
      State: {{ state }}
    </v-card-text>
  </trait-card>
</template>

<script>
import {Occupancy} from "@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb";
import {pullOccupancy} from "./occupancy";
import TraitCard from "../../components/TraitCard.vue";

export default {
  name: "OccupancyCard",
  components:{TraitCard},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        occupancy: {
          /** @type {Occupancy.AsObject} */
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
      if (!this.resources.occupancy.value) {
        return "LOADING";
      }
      switch (this.resources.occupancy.value.state) {
        case Occupancy.State.STATE_UNSPECIFIED:
          return "UNKNOWN";
        case Occupancy.State.OCCUPIED:
          return "OCCUPIED";
        case Occupancy.State.UNOCCUPIED:
          return "UNOCCUPIED";
        default:
          return "< unsupported >";
      }
    },
    count() {
      if (!this.resources.occupancy.value) {
        return "...";
      }
      return this.resources.occupancy.value.peopleCount;
    }
  },
  methods: {
    async pull() {
      // Occupancy resource
      this.resources.occupancy = await pullOccupancy(this.deviceId, this.resources.occupancy);
    },
    log(...args) {
      console.debug(...args);
    }
  }
}
</script>

<style scoped>

</style>
