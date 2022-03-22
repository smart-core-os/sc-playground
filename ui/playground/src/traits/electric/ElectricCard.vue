<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <template #title-append>
      <span>{{ demand.current.toFixed(2) }}A</span>
      <v-icon right>mdi-lightning-bolt</v-icon>
    </template>
    <v-card-text class="mt-2">
      <electric-mode-chooser :mode="activeMode" @update:mode="setMode" :modes="modes"
                             :label="`Active mode (of ${modes.length})`"/>
    </v-card-text>
    <v-card-text class="px-0">
      <v-expansion-panels flat accordion>
        <v-expansion-panel>
          <v-expansion-panel-header>More Details</v-expansion-panel-header>
          <v-expansion-panel-content>
            <pre>{{ allData }}</pre>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </trait-card>
</template>

<script>

import {grpcWebEndpoint} from '../../util/api.js';
import {ElectricApiPromiseClient} from "@smart-core-os/sc-api-grpc-web/traits/electric_grpc_web_pb.js";
import {
  ElectricMode,
  GetActiveModeRequest,
  GetDemandRequest,
  ListModesRequest,
  PullActiveModeRequest,
  PullDemandRequest,
  PullModesRequest,
  UpdateActiveModeRequest
} from "@smart-core-os/sc-api-grpc-web/traits/electric_pb.js";
import Vue from "vue";
import ElectricModeChooser from "./ElectricModeChooser.vue";
import TraitCard from "../../components/TraitCard.vue";
import {durationString, toDate} from "./util.js";

export default {
  name: 'ElectricCard',
  components: {TraitCard, ElectricModeChooser},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        demand: {
          /** @type {ElectricDemand.AsObject} */
          value: null,
          stream: null,
          err: null
        },
        activeMode: {
          /** @type {ElectricMode.AsObject} */
          value: null,
          stream: null,
          err: null
        },
        modes: {
          /** @type {Object<string,ElectricMode.AsObject>} */
          value: {},
          stream: null,
          err: null
        }
      }
    };
  },
  mounted() {
    this.pull()
        .catch(err => console.error('during pull', err))
  },
  beforeDestroy() {
    for (const resource of Object.values(this.resources)) {
      if (resource.stream) resource.stream.cancel();
    }
  },
  computed: {
    demand() {
      return this.resources.demand.value || {current: 0};
    },
    activeMode() {
      return this.resources.activeMode.value;
    },
    modes() {
      return Object.values(this.resources.modes.value);
    },
    allData() {
      let res = {
        demand: this.resources.demand.value,
        activeMode: this.resources.activeMode.value,
        modes: Object.values(this.resources.modes.value)
      }
      // clone
      res = JSON.parse(JSON.stringify(res))
      // convert dates, etc
      if (res.activeMode?.startTime) {
        res.activeMode.startTime = toDate(res.activeMode.startTime);
      }
      for (const segmentsListElement of res?.activeMode?.segmentsList || []) {
        if (segmentsListElement.length) {
          segmentsListElement.length = durationString(segmentsListElement.length);
        }
      }
      for (const mode of res.modes || []) {
        if (mode.startTime) {
          mode.startTime = toDate(mode.startTime);
        }
        for (const segmentsListElement of mode.segmentsList) {
          if (segmentsListElement.length) {
            segmentsListElement.length = durationString(segmentsListElement.length);
          }
        }
      }
      return res;
    }
  },
  methods: {
    async pull() {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ElectricApiPromiseClient(serverEndpoint, null, null);

      // demand resource
      const demandResource = this.resources.demand;
      if (demandResource.stream) demandResource.stream.cancel();
      try {
        const demandPb = await api.getDemand(new GetDemandRequest().setName(this.deviceId));
        demandResource.value = demandPb.toObject();
        const demandStream = api.pullDemand(new PullDemandRequest().setName(this.deviceId));
        demandResource.stream = demandStream;
        demandStream.on('data', res => {
          /** @type {PullDemandResponse.Change[]} */
          const changes = res.getChangesList();
          for (const change of changes) {
            const value = change.getDemand();
            demandResource.value = value.toObject();
          }
        });
        demandStream.on('error', err => demandResource.err = err)
      } catch (e) {
        demandResource.err = e;
      }

      // activeMode resource
      const activeModeResource = this.resources.activeMode;
      if (activeModeResource.stream) activeModeResource.stream.cancel();
      try {
        const activeModePb = await api.getActiveMode(new GetActiveModeRequest().setName(this.deviceId));
        activeModeResource.value = activeModePb.toObject();
        const activeModeStream = api.pullActiveMode(new PullActiveModeRequest().setName(this.deviceId));
        activeModeResource.stream = activeModeStream;
        activeModeStream.on('data', res => {
          /** @type {PullActiveModeResponse.Change[]} */
          const changes = res.getChangesList();
          for (const change of changes) {
            const value = change.getActiveMode();
            activeModeResource.value = value.toObject();
          }
        });
        activeModeStream.on('error', err => activeModeResource.err = err)
      } catch (e) {
        activeModeResource.err = e;
      }

      // modes resource
      const modesResource = this.resources.modes;
      if (modesResource.stream) modesResource.stream.cancel();
      try {
        let modesPb = await api.listModes(new ListModesRequest().setName(this.deviceId));
        while (true) {
          for (const mode of modesPb.getModesList()) {
            Vue.set(modesResource.value, mode.getId(), mode.toObject());
          }
          if (!modesPb.getNextPageToken()) {
            break;
          }
          modesPb = await api.listModes(new ListModesRequest().setName(this.deviceId)
              .setPageToken(modesPb.getNextPageToken()));
        }
        const modesStream = api.pullModes(new PullModesRequest().setName(this.deviceId));
        modesResource.stream = modesStream;
        modesStream.on('data', res => {
          /** @type {PullModesResponse.Change[]} */
          const changes = res.getChangesList();
          for (const change of changes) {
            const value = change.getNewValue();
            if (!value) {
              // delete
              const oldId = change.getOldValue()?.getId();
              if (oldId) {
                Vue.delete(modesResource.value, oldId)
              }
            } else {
              Vue.set(modesResource.value, value.getId(), value.toObject());
            }
          }
        });
        modesStream.on('error', err => modesResource.err = err)
      } catch (e) {
        modesResource.err = e;
      }
    },
    log(...args) {
      console.debug(...args);
    },
    async setMode(mode) {
      const serverEndpoint = await grpcWebEndpoint();
      const api = new ElectricApiPromiseClient(serverEndpoint, null, null);
      await api.updateActiveMode(new UpdateActiveModeRequest().setName(this.deviceId)
          .setActiveMode(new ElectricMode().setId(mode.id)))
      // todo: error handling
    }
  }
};
</script>

<style scoped>
</style>
