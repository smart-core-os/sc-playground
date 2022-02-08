<template>
  <v-card width="400">
    <v-card-title>
      <span>{{ deviceId }}</span>
      <v-spacer/>
      <slot name="title-append">
        <v-icon v-if="icon">{{ icon }}</v-icon>
      </slot>
    </v-card-title>
    <v-card-subtitle>
      Implements the
      <a-new :href="traitProtoHREF" target="_blank">{{ traitLocalName }}</a-new>
      trait
    </v-card-subtitle>
    <slot/>
  </v-card>
</template>

<script>
import {snakeCase} from 'change-case';
import ANew from "./ANew.vue";

export default {
  name: "TraitCard",
  components: {ANew},
  props: {
    deviceId: [String],
    trait: [Object],
    icon: [String]
  },
  computed: {
    traitLocalName() {
      if (!this.trait) return '';
      const traitParts = this.trait.name.split('.');
      return traitParts[traitParts.length - 1];
    },
    traitProtoHREF() {
      const base = 'https://github.com/smart-core-os/sc-api/blob/main/protobuf/traits/';
      const traitProtoFile = snakeCase(this.traitLocalName) + '.proto';
      return base + traitProtoFile;
    }
  }
}
</script>

<style scoped>

</style>
