<template>
  <v-card>
    <v-card-title>Add a remote device</v-card-title>
    <v-card-text>
      <v-text-field label="Device name" v-model="name"/>
      <v-text-field type="url" label="Endpoint" v-model="endpoint"/>
      <v-item-group multiple v-model="selectedTraits">
        <v-item v-for="trait in supportedTraits" :key="trait" v-slot="{ active, toggle }">
          <v-checkbox :value="active" @change="toggle" :label="trait" hide-details/>
        </v-item>
      </v-item-group>
    </v-card-text>
    <v-card-actions>
      <v-spacer/>
      <v-btn @click="cancel" text>Cancel</v-btn>
      <v-btn @click="add" depressed>Add</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import {grpcWebEndpoint} from "../../util/api.js";
import {PlaygroundApiPromiseClient} from "@sc-playground/gen/pkg/playpb/playground_grpc_web_pb.js";
import {AddRemoteDeviceRequest, ListSupportedTraitsRequest} from "@sc-playground/gen/pkg/playpb/playground_pb.js";

export default {
  name: "AddRemoteDeviceCard",
  data() {
    return {
      loading: false,

      supportedTraitsLoading: true,
      supportedTraits: [],

      name: '',
      endpoint: '',
      selectedTraits: []
    };
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
      } finally {
        this.supportedTraitsLoading = false;
      }
    },
    clear() {
      this.loading = false;
      this.name = '';
      this.endpoint = '';
      this.selectedTraits = [];
    },
    async add() {
      this.loading = true;
      const serverEndpoint = await grpcWebEndpoint();
      const api = new PlaygroundApiPromiseClient(serverEndpoint, null, null);
      try {
        const traitNames = this.selectedTraits.map(i => this.supportedTraits[i])
        await api.addRemoteDevice(new AddRemoteDeviceRequest()
            .setName(this.name)
            .setEndpoint(this.endpoint)
            .setTraitNameList(traitNames)
        );
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
