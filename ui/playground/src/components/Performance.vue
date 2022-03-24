<template>
  <div class="px-4 py-2 grey--text text--lighten-1">{{ frameMs }}ms/frame</div>
</template>

<script>
import {grpcWebEndpoint} from "../util/api.js";
import {PlaygroundApiClient} from "@sc-playground/gen/pkg/playpb/playground_grpc_web_pb.js";
import {PullPerformanceRequest} from "@sc-playground/gen/pkg/playpb/playground_pb.js";
import {durationMillis} from "../util/time.js";

export default {
  name: "Performance",
  data() {
    return {
      resources: {
        performance: {
          stream: null,
          /** @type {Performance.AsObject} */
          value: null,
          error: null
        }
      }
    }
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
    frameMs() {
      return this.toMillis(this.resources.performance.value?.frame);
    },
    captureMs() {
      return this.toMillis(this.resources.performance.value?.capture);
    },
    scrubMs() {
      return this.toMillis(this.resources.performance.value?.scrub);
    },
    respondMs() {
      return this.toMillis(this.resources.performance.value?.respond);
    },
  },
  methods: {
    toMillis(d) {
      if (!d) {
        return 0;
      }
      return durationMillis(d);
    },
    async pull() {
      const serverEndpoint = await grpcWebEndpoint();
      return this.pullPerformance(serverEndpoint);
    },
    async pullPerformance(serverEndpoint) {
      // Performance resource
      const playgroundApiClient = new PlaygroundApiClient(serverEndpoint, null, null);
      const performanceResource = this.resources.performance;
      if (performanceResource.stream) performanceResource.stream.cancel();
      const performanceStream = playgroundApiClient.pullPerformance(new PullPerformanceRequest());
      performanceResource.stream = performanceStream;
      performanceStream.on('data', (/** @type {PullPerformanceResponse} */ res) => {
        performanceResource.value = res.getPerformance().toObject();
      });
    }
  }
}
</script>

<style scoped>

</style>
