<template>
  <trait-card :device-id="deviceId" :trait="trait"/>
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
