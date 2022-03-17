<template>
  <v-app id="root">
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click="navDrawer = !navDrawer"/>
      <v-app-bar-title>
        <!--
        This span works around a rounding issue in the width calculations of v-app-bar-title.
        The component duplicates the title content (default slot) and calculates the width of that duplicate
        dom tree. The width is then set explicitly on the visible content.
        If the width happens to be a fractional pixel (i.e. 123.55) then some rounding happens that causes the explicit
        width to be truncated to an integer, causing the displayed elements to be narrower than the need to be, resulting
        in ellipses (...) instead of the full text.

        Adding this 1px padding forces the bounding box to be wider so the hard coded width doesn't ellipses the text.
        -->
        <span style="padding-right: 1px">Smart Core: Playground</span>
      </v-app-bar-title>
      <v-spacer/>
      <template v-if="serverConfigError">
        {{ serverConfigError }}
      </template>
      <template v-else-if="!serverConfig">
        Loading...
      </template>
      <template v-else>
        <span>
          gRPC address:
        <code>{{ serverConfig.grpcAddress }}</code>
        <template v-if="serverConfig.insecure">(insecure)</template>
        <template v-else-if="showCaDownload">
          <v-tooltip bottom>
            <template #activator="{ on, attrs }">
              <v-btn v-on="on" v-bind="attrs" outlined download :href="caCertPath" class="ml-2">
                CA Cert <v-icon right>mdi-download</v-icon></v-btn>
            </template>
            Download the self-signed CA certificate in PEM format.<br/>
            Add this cert into your root-ca pool to verify the identity of the server.
          </v-tooltip>
        </template></span>
        <template v-if="showClientCertDownload">
          <v-tooltip bottom>
            <template #activator="{ on, attrs }">
              <v-btn v-on="on" v-bind="attrs" outlined download :href="clientCertPath" class="ml-2">
                Client Creds
                <v-icon right>mdi-download-lock</v-icon>
              </v-btn>
            </template>
            Download a set of new client credentials in PEM format.<br/>
            Will contain a certificate and private key block.
          </v-tooltip>
        </template>
      </template>
    </v-app-bar>
    <v-navigation-drawer app clipped v-model="navDrawer">
      <v-list>
        <v-list-item :to="{name: 'devices'}">
          <v-list-item-title>Devices</v-list-item-title>
        </v-list-item>
        <v-list-item :to="{name: 'traits'}">
          <v-list-item-title>Traits</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>
import {caCertPath, clientCertPath, serverConfig} from './util/api.js';

export default {
  name: 'App',
  data() {
    return {
      /** @type {ServerConfig|null} */
      serverConfig: null,
      serverConfigError: null,

      navDrawer: true
    };
  },
  computed: {
    caCertPath() {
      return caCertPath();
    },
    clientCertPath() {
      return clientCertPath();
    },
    showCaDownload() {
      return this.serverConfig && this.serverConfig.selfSigned && !this.serverConfig.insecure;
    },
    showClientCertDownload() {
      return this.serverConfig && this.serverConfig.selfSigned && this.serverConfig.mutualTls;
    }
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
