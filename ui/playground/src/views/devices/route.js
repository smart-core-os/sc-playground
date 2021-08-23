/** @type {import('vue-router').RouteConfig} */
const route = {
  path: '/devices',
  component: () => import('./DevicesView.vue')
}
export default route;
