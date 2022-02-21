<template>
  <v-container fluid class="card-grid">
    <component v-for="child of resources.children.value" :key="child.name"
               :is="childToComponent(child)"
               :child="child" :metadata="metadataForChild(child)"/>
  </v-container>
</template>

<script>
import {grpcWebEndpoint} from '../../util/api.js';
import {MetadataApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/metadata_grpc_web_pb.js';
import {GetMetadataRequest, PullMetadataRequest} from '@smart-core-os/sc-api-grpc-web/traits/metadata_pb.js';
import {ParentApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/parent_grpc_web_pb.js';
import {ListChildrenRequest, PullChildrenRequest} from '@smart-core-os/sc-api-grpc-web/traits/parent_pb.js';
import Vue from 'vue';
import MetadataPanel from '../../traits/metadata/MetadataPanel.vue';
import EVChargerCard from '../../devices/evcharger/EVChargerCard.vue';
import UnknownDeviceCard from '../../devices/unknown/UnknownDeviceCard.vue';

export default {
  name: 'DevicesView',
  components: {MetadataPanel, EVChargerCard, UnknownDeviceCard},
  data() {
    return {
      serverName: '',
      resources: {
        children: {
          /** @type {Child.AsObject[]} */
          value: [],
          stream: null
        },
        /** @type {Object<string, {value: Metadata.AsObject, stream}>} */
        metadataByDeviceId: {}
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
    for (const resource of Object.values(this.resources.metadataByDeviceId)) {
      if (resource.stream) resource.stream.cancel();
    }
  },
  computed: {},
  watch: {
    'resources.children.value': {
      handler(newVal, oldVal) {
        this.syncMetadataPull(newVal, oldVal);
      }
    }
  },
  methods: {
    async pull() {
      return this.pullChildren();
    },
    async pullChildren() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ParentApiPromiseClient(serverEndpoint, null, null);
      // children
      const resource = this.resources.children;
      if (resource.stream) resource.stream.cancel();
      // get
      let listVal = await api.listChildren(new ListChildrenRequest().setName(this.serverName));
      resource.value = listVal.getChildrenList().map(c => c.toObject());
      while (listVal.getNextPageToken()) {
        listVal = await api.listChildren(new ListChildrenRequest()
            .setName(this.serverName)
            .setPageToken(listVal.getNextPageToken()));
        resource.value.push(...listVal.getChildrenList().map(c => c.toObject()));
      }
      // pull
      const stream = api.pullChildren(new PullChildrenRequest().setName(this.serverName));
      resource.stream = stream;
      stream.on('data', res => {
        /** @type {PullChildrenResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          if (!change.getNewValue()) {
            // value was removed
            const old = change.getOldValue();
            if (!old) continue; // sanity check
            const oldIndex = resource.value.findIndex(c => c.name === old.getName());
            if (oldIndex >= 0) {
              resource.value.splice(oldIndex, 1);
            }
          } else {
            const newValue = change.getNewValue().toObject();
            const index = resource.value.findIndex(c => c.name === newValue.name);
            if (index < 0) {
              resource.value.push(newValue);
            } else {
              Vue.set(resource.value, index, newValue);
            }
          }
        }
      });
    },
    async pullMetadata(deviceId) {
      const serverEndpoint = await grpcWebEndpoint();
      // Metadata resource
      const api = new MetadataApiPromiseClient(serverEndpoint, null, null);
      const resource = this.resources.metadataByDeviceId[deviceId] ||
          Vue.set(this.resources.metadataByDeviceId, deviceId, {value: null, stream: null});
      if (resource.stream) resource.stream.cancel();
      const getVal = await api.getMetadata(new GetMetadataRequest().setName(deviceId));
      resource.value = getVal.toObject();
      const stream = api.pullMetadata(new PullMetadataRequest().setName(deviceId));
      resource.stream = stream;
      stream.on('data', res => {
        /** @type {PullMetadataResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const value = change.getMetadata();
          resource.value = value.toObject();
        }
      });
    },
    /**
     * @param {Child.AsObject[]} newChildren
     * @param {Child.AsObject[]} oldChildren
     */
    syncMetadataPull(newChildren, oldChildren) {
      // fetch metadata (or cancel streams) for each child we know about
      const oldIndex = (oldChildren || []).reduce((arr, item) => {
        arr[item.name] = item;
        return arr;
      }, {});
      const toCancel = new Set(Object.keys(oldIndex));
      /** @type {Child.AsObject[]} */
      const toAdd = [];
      for (const child of newChildren) {
        toCancel.delete(child.name); // don't cancel watchers we still using
        if (!oldIndex[child.name]) {
          toAdd.push(child); // add watchers for children we haven't seen before
        }
      }

      for (const deviceId of toCancel) {
        const resource = this.resources.metadataByDeviceId[deviceId];
        if (resource) {
          if (resource.stream) {
            resource.stream.cancel();
          }
          Vue.delete(this.resources.metadataByDeviceId, deviceId);
        }
      }
      for (const child of toAdd) {
        this.pullMetadata(child.name);
      }
    },

    metadataForChild(child) {
      return this.resources.metadataByDeviceId[child.name]?.value;
    },
    log(...args) {
      console.debug(...args);
    },

    childToComponent(child) {
      const md = this.metadataForChild(child);
      const entry = md?.moreMap?.find(([k]) => k === 'scos.playground.device-type');
      if (!entry) {
        return 'UnknownDeviceCard';
      }
      const deviceType = entry[1];
      switch (deviceType) {
        case 'evcharger':
          return 'EVChargerCard';
      }
      return 'UnknownDeviceCard';
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
