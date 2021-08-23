<template>
  <div>
    <v-slider label="Rating" :value="rating" @input="$emit('update:rating', $event)" hide-details/>
    <v-slider label="Voltage" :value="voltage" @input="$emit('update:voltage', $event)" hide-details :max="300"/>
    <v-slider label="Load" :value="load" @input="$emit('update:load', $event)" hide-details :max="rating"/>
    <v-slider label="Reserved" :value="reserved" @input="$emit('update:reserved', $event)" hide-details :max="rating"/>
  </div>
</template>

<script>
export default {
  name: 'PowerSupplySettingsEditor',
  props: ['rating', 'voltage', 'load', 'reserved']
};
</script>

<style scoped>

</style>
