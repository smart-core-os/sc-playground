import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router'

const app = new Vue({
  render: h => h(App),
  router,
  vuetify
});
app.$mount('#app');
