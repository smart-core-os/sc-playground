<template>
  <v-list dense>
    <v-list-item v-for="item in displayItems" :key="item.id">
      <v-list-item-icon>
        <v-icon>mdi-chart-line-variant</v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>{{ item.maxDraw }}&nbsp;A for {{ item.remainingSeconds }}&nbsp;s</v-list-item-title>
      </v-list-item-content>
      <v-btn icon @click="$emit('click:remove', item.id)">
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </v-list-item>
  </v-list>
</template>

<script>
export default {
  name: "PowerSupplyDrawNotificationList",
  props: {
    items: Array
  },
  data() {
    return {
      nowSecs: Math.round(Date.now() / 1000),
      nowHandle: 0
    }
  },
  mounted() {
    this.nowHandle = setInterval(() => {
      this.nowSecs = Math.round(Date.now() / 1000);
    }, 1000);
  },
  beforeDestroy() {
    clearInterval(this.nowHandle);
  },
  computed: {
    displayItems() {
      if (!this.items) return [];
      return this.items.map(item => {
        const endTime = item.notificationTime.seconds + item.rampDuration.seconds;
        const remainingSeconds = endTime - this.nowSecs;
        return {
          id: item.id,
          maxDraw: item.maxDraw,
          rampDurationSeconds: item.rampDuration.seconds,
          remainingSeconds
        };
      });
    }
  }
}
</script>

<style scoped>

</style>
