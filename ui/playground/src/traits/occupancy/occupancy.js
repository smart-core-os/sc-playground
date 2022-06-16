import {grpcWebEndpoint} from '../../util/api.js';
import {OccupancySensorApiPromiseClient} from "@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_grpc_web_pb";
import {GetOccupancyRequest, Occupancy, PullOccupancyRequest} from "@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb";

/**
 * @param {string} deviceId
 * @param {{stream?: *, value?: Occupancy.AsObject}} [resource]
 * @return {Promise<{stream: *, value: Occupancy.AsObject}>}
 */
export async function pullOccupancy(deviceId, resource = {value: undefined, stream: undefined}) {
  const serverEndpoint = await grpcWebEndpoint();

  const api = new OccupancySensorApiPromiseClient(serverEndpoint, null, null);
  if (resource.stream) resource.stream.cancel();
  const getResult = await api.getOccupancy(new GetOccupancyRequest().setName(deviceId));
  resource.value = getResult.toObject();
  const stream = api.pullOccupancy(new PullOccupancyRequest().setName(deviceId));
  resource.stream = stream;
  stream.on('data', res => {
    /** @type {PullOccupancyResponse.Change[]} */
    const changes = res.getChangesList();
    for (const change of changes) {
      const value = change.getOccupancy();
      resource.value = value.toObject();
    }
  });
  return resource;
}

/**
 * @param {string} deviceId
 * @param {Occupancy.AsObject} occupancy
 * @param {{loading: number}} [resource]
 * @return {Promise<Occupancy.AsObject>}
 */
/*export async function updateOccupancy(deviceId, occupancy, resource = {loading: 0}) {
  resource.loading++;
  try {
    const serverEndpoint = await grpcWebEndpoint();

    const api = new OnOffApiPromiseClient(serverEndpoint, null, null);
    const req = new UpdateOnOffRequest().setName(deviceId);
    req.setOnOff(new OnOff().setState(onOff.state))
    const getResult = await api.updateOnOff(req);
    return getResult.toObject();
  } finally {
    resource.loading--;
  }
}*/
