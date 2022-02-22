import {grpcWebEndpoint} from '../../util/api.js';
import {EnergyStorageApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/energy_storage_grpc_web_pb.js';
import {
  GetEnergyLevelRequest,
  PullEnergyLevelRequest
} from '@smart-core-os/sc-api-grpc-web/traits/energy_storage_pb.js';

/**
 * @param {string} deviceId
 * @param {{stream?: *, value?: EnergyLevel.AsObject}} [resource]
 * @return {Promise<{stream: *, value: EnergyLevel.AsObject}>}
 */
export async function pullEnergyLevel(deviceId, resource = {value: undefined, stream: undefined}) {
  const serverEndpoint = await grpcWebEndpoint();

  // EnergyLevel resource
  const api = new EnergyStorageApiPromiseClient(serverEndpoint, null, null);
  if (resource.stream) resource.stream.cancel();
  const getResult = await api.getEnergyLevel(new GetEnergyLevelRequest().setName(deviceId));
  resource.value = getResult.toObject();
  const stream = api.pullEnergyLevel(new PullEnergyLevelRequest().setName(deviceId));
  resource.stream = stream;
  stream.on('data', res => {
    /** @type {PullEnergyLevelResponse.Change[]} */
    const changes = res.getChangesList();
    for (const change of changes) {
      const value = change.getEnergyLevel();
      resource.value = value.toObject();
    }
  });
  return resource;
}
