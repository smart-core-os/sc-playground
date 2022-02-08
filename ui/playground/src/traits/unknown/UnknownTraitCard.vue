<template>
  <v-card width="400">
    <v-card-title>
      <span>{{ deviceId }}</span>
      <v-spacer/>
    </v-card-title>
    <v-card-text>
      This device implements an unknown trait {{ trait }}
    </v-card-text>
  </v-card>
</template>

<script>

export default {
  name: 'UnknownTraitCard',
  props: {
    deviceId: null,
    trait: null
  }
};
</script>

<style scoped>
</style>
