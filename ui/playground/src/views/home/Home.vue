<template>
  <v-container class="card-grid">
    <power-supply-card v-for="deviceId in devices" :device-id="deviceId" :key="deviceId"/>

    <form @submit.prevent="addDevice">
      <v-card>
        <v-card-title>Add or expose another device</v-card-title>
        <v-card-text>
          <v-text-field label="Device ID" v-model="nextDeviceId"/>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn elevation="0" color="success" type="submit">Add</v-btn>
        </v-card-actions>
      </v-card>
    </form>
  </v-container>
</template>

<script>

import PowerSupplyCard from '../../traits/powersupply/PowerSupplyCard.vue';
import {grpcWebEndpoint} from "../../util/api.js";
import {ParentApiPromiseClient} from "@smart-core-os/sc-api-grpc-web/traits/parent_grpc_web_pb.js";
import {ListChildrenRequest, PullChildrenRequest} from "@smart-core-os/sc-api-grpc-web/traits/parent_pb.js";
import Vue from "vue";

export default {
  name: 'Home',
  components: {PowerSupplyCard},
  data() {
    return {
      serverName: '',
      childrenStream: null,
      children: [],
      manualNextDeviceId: null
    };
  },
  computed: {
    devices() {
      return this.children.map(c => c.name);
    },
    nextDeviceId: {
      get() {
        if (this.manualNextDeviceId != null) {
          return this.manualNextDeviceId;
        }
        const nextIdNum = (this.devices.length + 1).toString().padStart(3, '0');
        return `POW-${nextIdNum}`;
      },
      set(v) {
        this.manualNextDeviceId = v;
      }
    }
  },
  mounted() {
    this.pull()
        .catch(err => console.error('during initial pull', err));
  },
  beforeDestroy() {
    if (this.childrenStream) this.childrenStream.cancel();
  },
  methods: {
    async addDevice() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ParentApiPromiseClient(serverEndpoint, null, null);
      await api.listChildren(new ListChildrenRequest().setName(this.nextDeviceId)); // used to trigger a new device
      this.manualNextDeviceId = null;
    },
    async pull() {
      await Promise.all([
        this.pullChildren()
      ]);
    },
    async pullChildren() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ParentApiPromiseClient(serverEndpoint, null, null);
      // children
      if (this.childrenStream) this.childrenStream.cancel();
      // get
      let childrenRes = await api.listChildren(new ListChildrenRequest().setName(this.serverName));
      this.children = childrenRes.getChildrenList().map(c => c.toObject());
      while (childrenRes.getNextPageToken()) {
        childrenRes = await api.listChildren(new ListChildrenRequest()
            .setName(this.serverName)
            .setPageToken(childrenRes.getNextPageToken()));
        this.children.push(...childrenRes.getChildrenList().map(c => c.toObject()))
      }
      // pull
      const childrenStream = api.pullChildren(new PullChildrenRequest().setName(this.serverName));
      this.childrenStream = childrenStream;
      childrenStream.on('data', res => {
        /** @type {PullChildrenResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          if (!change.getNewValue()) {
            // value was removed
            const old = change.getOldValue();
            if (!old) continue; // sanity check
            const oldIndex = this.children.findIndex(c => c.name === old.getName());
            if (oldIndex >= 0) {
              this.children.splice(oldIndex, 1);
            }
          } else {
            const newValue = change.getNewValue().toObject()
            const index = this.children.findIndex(c => c.name === newValue.name);
            if (index < 0) {
              this.children.push(newValue)
            } else {
              Vue.set(this.children, index, newValue)
            }
          }
        }
      });
    }
  }
};
</script>

<style scoped>
.card-grid {
  display: grid;
  grid-gap: 12px;
  grid-template-columns: repeat(auto-fill, 400px);
  grid-auto-columns: 400px;
  justify-content: center;
}
</style>
