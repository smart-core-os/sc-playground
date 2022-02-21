/** @type {import('vue-router').RouteConfig} */
const route = {
  path: '/devices',
  name: 'devices',
  component: () => import('./DeviceCards.vue')
}
export default route;
