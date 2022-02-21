<template>
  <trait-card :device-id="deviceId" :trait="trait" icon="mdi-information">
    <metadata-panel :metadata="resources.metadata.value" multiple flat class="mb-4"/>
  </trait-card>
</template>

<script>
import TraitCard from '../../components/TraitCard.vue';
import {grpcWebEndpoint} from '../../util/api.js';
import {GetMetadataRequest, PullMetadataRequest} from '@smart-core-os/sc-api-grpc-web/traits/metadata_pb.js';
import {MetadataApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/metadata_grpc_web_pb.js';
import MetadataPanel from './MetadataPanel.vue';

export default {
  name: 'MetadataCard',
  components: {MetadataPanel, TraitCard},
  props: {
    deviceId: null,
    trait: null
  },
  data() {
    return {
      resources: {
        metadata: {
          /** @type {Metadata.AsObject} */
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
      const serverEndpoint = await grpcWebEndpoint();

      // Metadata resource
      const metadataApi = new MetadataApiPromiseClient(serverEndpoint, null, null);
      const metadataResource = this.resources.metadata;
      if (metadataResource.stream) metadataResource.stream.cancel();
      const metadataPb = await metadataApi.getMetadata(new GetMetadataRequest().setName(this.deviceId));
      metadataResource.value = metadataPb.toObject();
      const metadataStream = metadataApi.pullMetadata(new PullMetadataRequest().setName(this.deviceId));
      metadataResource.stream = metadataStream;
      metadataStream.on('data', res => {
        /** @type {PullMetadataResponse.Change[]} */
        const changes = res.getChangesList();
        for (const change of changes) {
          const value = change.getMetadata();
          metadataResource.value = value.toObject();
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
