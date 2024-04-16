import {grpcWebEndpoint} from '../../util/api.js';
import {AirTemperatureApiPromiseClient} from "@smart-core-os/sc-api-grpc-web/traits/air_temperature_grpc_web_pb.js";
import {GetAirTemperatureRequest,PullAirTemperatureRequest} from "@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb";

/**
 * @param {string} deviceId
 * @param {{stream?: *, value?: AirTemperature.AsObject}} [resource]
 * @return {Promise<{stream: *, value: AirTemperature.AsObject}>}
 */
export async function pullAirTemperature(deviceId, resource = {value: undefined, stream: undefined}) {
  const serverEndpoint = await grpcWebEndpoint();

  const api = new AirTemperatureApiPromiseClient(serverEndpoint, null, null);
  if (resource.stream) resource.stream.cancel();
  const getResult = await api.getAirTemperature(new GetAirTemperatureRequest().setName(deviceId));
  resource.value = getResult.toObject();
  const stream = api.pullAirTemperature(new PullAirTemperatureRequest().setName(deviceId));
  resource.stream = stream;
  stream.on('data', res => {
    /** @type {PullAirTemperatureRespnse.Change[]} */
    const changes = res.getChangesList();
    for (const change of changes) {
      const value = change.getAirTemperature();
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
