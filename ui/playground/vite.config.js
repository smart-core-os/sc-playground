import {defineConfig} from 'vite';
import {createVuePlugin} from 'vite-plugin-vue2';
import VitePluginComponents, {VuetifyResolver} from 'vite-plugin-components';
import rollupPluginVuetify from 'rollup-plugin-vuetify';

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    include: [
      '@smart-core-os/sc-api-grpc-web/traits/parent_pb.js',
      '@smart-core-os/sc-api-grpc-web/traits/parent_grpc_web_pb.js',
      '@smart-core-os/sc-api-grpc-web/traits/power_supply_pb.js',
      '@smart-core-os/sc-api-grpc-web/traits/power_supply_grpc_web_pb.js'
    ]
  },
  plugins: [
    createVuePlugin(),
    VitePluginComponents({
      dirs: '',
      customComponentResolvers: [
        VuetifyResolver()
      ]
    }),
    rollupPluginVuetify()
  ],
  server: {
    fs: {
      allow: ['..']
    }
  }
})
