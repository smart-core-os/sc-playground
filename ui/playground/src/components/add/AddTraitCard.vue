<template>
  <v-card :loading="supportedTraitsLoading">
    <v-card-title>Add a virtual device</v-card-title>
    <v-card-text>
      <v-text-field label="Device Name" v-model="name" hide-details autofocus/>
      <v-item-group multiple v-model="selectedTraits">
        <v-item v-for="trait in supportedTraits" :key="trait" v-slot="{ active, toggle }">
          <v-checkbox :value="active" @change="toggle" :label="trait" hide-details/>
        </v-item>
      </v-item-group>
    </v-card-text>
    <v-card-actions>
      <v-spacer/>
      <v-btn :loading="loading" text @click="cancel">Cancel</v-btn>
      <v-btn :loading="loading" depressed @click="add" :disabled="selectedTraits.length === 0 || name.length === 0">
        Add
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import {grpcWebEndpoint} from "../../util/api.js";
import {
  AddDeviceTraitRequest,
  PlaygroundApiPromiseClient
} from "@sc-playground/gen/pkg/playpb/playground_grpc_web_pb.js";
import {ListSupportedTraitsRequest} from "@sc-playground/gen/pkg/playpb/playground_pb.js";

export default {
  name: "AddTraitCard",
  data() {
    return {
      loading: false,

      supportedTraitsLoading: true,
      supportedTraits: [],

      name: '',
      selectedTraits: []
    }
  },
  mounted() {
    this.pull();
  },
  methods: {
    async pull() {
      try {
        const serverEndpoint = await grpcWebEndpoint();
        const api = new PlaygroundApiPromiseClient(serverEndpoint, null, null);
        const response = await api.listSupportedTraits(new ListSupportedTraitsRequest());
        this.supportedTraits = response.getTraitNameList();
        this.supportedTraits.sort();
      } finally {
        this.supportedTraitsLoading = false;
      }
    },
    clear() {
      this.loading = false;
      this.name = '';
      this.selectedTraits = [];
    },
    async add() {
      this.loading = true;
      const serverEndpoint = await grpcWebEndpoint();
      const api = new PlaygroundApiPromiseClient(serverEndpoint, null, null);
      try {
        for (const selectedTrait of this.selectedTraits) {
          await api.addDeviceTrait(new AddDeviceTraitRequest().setName(this.name).setTraitName(this.supportedTraits[selectedTrait]));
        }
      } finally {
        this.$emit('done');
        this.clear();
      }
    },
    cancel() {
      this.$emit('done');
      this.clear();
    }
  }
}
</script>

<style scoped>

</style>
