<template>
  <trait-card :device-id="deviceId" :trait="trait">
    <v-card-text>
      <v-alert color="info" outlined>The UI for this trait has not yet been implemented.</v-alert>
    </v-card-text>
  </trait-card>
</template>

<script>
import TraitCard from '../../components/TraitCard.vue';

export default {
  name: 'UnknownTraitCard',
  components: {TraitCard},
  props: {
    deviceId: null,
    trait: null
  }
};
</script>

<style scoped>
</style>
