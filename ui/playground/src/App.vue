<template>
  <v-app id="root">
    <v-system-bar>
      <template v-if="serverConfigError">
        {{ serverConfigError }}
      </template>
      <template v-else-if="!serverConfig">
        Loading...
      </template>
      <template v-else>
        Smart Core API address:
        <code>{{ serverConfig.grpcAddress }}</code>
        {{ serverConfig.insecure ? ' (insecure)' : '' }}
      </template>
    </v-system-bar>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>
import {serverConfig} from './util/api.js';

export default {
  name: 'App',
  data() {
    return {
      /** @type {ServerConfig|null} */
      serverConfig: null,
      serverConfigError: null
    };
  },
  mounted() {
    serverConfig()
        .then(config => this.serverConfig = config)
        .catch(err => this.serverConfigError = err);
  }
};
</script>

<style scoped>
#root {
  background: #f1f1f1;
}
</style>
