<template>
  <div class="fill-height">
    <v-container class="card-grid" fluid>
      <template v-for="({child, trait}) in filteredChildTraits">
        <component :is="traitToComponent(trait)" :deviceId="child.name" :trait="trait" v-bind="overrides[child.name]"/>
      </template>
    </v-container>
    <add-device-fab/>
  </div>
</template>

<script>

import {grpcWebEndpoint} from '../../util/api.js';
import {ParentApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/parent_grpc_web_pb.js';
import {ListChildrenRequest, PullChildrenRequest} from '@smart-core-os/sc-api-grpc-web/traits/parent_pb.js';
import Vue from 'vue';
import EnergyStorageCard from '../../traits/energystorage/EnergyStorageCard.vue';
import UnknownTraitCard from '../../traits/unknown/UnknownTraitCard.vue';
import ElectricCard from '../../traits/electric/ElectricCard.vue';
import MetadataCard from '../../traits/metadata/MetadataCard.vue';
import AddDeviceFab from '../../components/add/AddDeviceFab.vue';
import OnOffCard from '../../traits/onoff/OnOffCard.vue';
import PublicationCard from '../../traits/publication/PublicationCard.vue';

export default {
  name: 'TraitCards',
  components: {
    AddDeviceFab,
    ElectricCard,
    EnergyStorageCard,
    OnOffCard,
    MetadataCard,
    PublicationCard,
    UnknownTraitCard
  },
  data() {
    return {
      serverName: '',
      childrenStream: null,
      /** @type {Child.AsObject[]} */
      children: [],
      manualNextDeviceId: null,

      first: ["scos/play/g/electric/sum", "birmingham6", "birmingham14", "AV", "Background Audio", "BMS", "Lighting", "Plant", "Signage"],
      filterName: {},
      filterTrait: {"smartcore.traits.Metadata": true},
      overrides: {
        "birmingham6": {color: "#72e0aa", elevation: 8},
        "birmingham14": {color: "#72e0aa", elevation: 8},
        "scos/play/g/electric/sum": {color: "#6bd9d9", elevation: 8, style: "grid-column: 1 / -1", width: "auto"},
      }
    };
  },
  computed: {
    filteredChildTraits() {
      const firstIndex = this.first.reduce((all, name, i) => {
        all[name] = i;
        return all;
      }, {})
      const firstItems = [];
      const res = [];
      for (const child of this.children) {
        for (const trait of child.traitsList) {
          // filter out metadata (for now)
          if (this.filterTrait[trait.name]) {
            continue
          }
          if (this.filterName[child.name]) {
            continue
          }
          if (firstIndex.hasOwnProperty(child.name)) {
            firstItems.push({child, trait});
          } else {
            res.push({child, trait});
          }
        }
      }

      firstItems.sort((a, b) => {
        return firstIndex[a.child.name] - firstIndex[b.child.name];
      })
      res.sort((a, b) => {
        return a.child.name.localeCompare(b.child.name);
      })

      return [...firstItems, ...res];
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
    },
    /**
     * @param {Trait.AsObject} trait
     * @return {string}
     */
    traitToComponent(trait) {
      switch (trait.name) {
        case 'smartcore.traits.Electric':
          return 'ElectricCard';
        case 'smartcore.traits.EnergyStorage':
          return 'EnergyStorageCard';
        case 'smartcore.traits.Metadata':
          return 'MetadataCard';
        case 'smartcore.traits.OnOff':
          return 'OnOffCard';
        case 'smartcore.traits.Publication':
          return 'PublicationCard';
        default:
          return 'UnknownTraitCard';
      }
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
