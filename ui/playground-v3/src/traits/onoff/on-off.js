import {grpcWebEndpoint} from '../../util/api.js';
import {OnOffApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/on_off_grpc_web_pb.js';
import {
  GetOnOffRequest,
  OnOff,
  PullOnOffRequest,
  UpdateOnOffRequest
} from '@smart-core-os/sc-api-grpc-web/traits/on_off_pb.js';

/**
 * @param {string} deviceId
 * @param {{stream?: *, value?: OnOff.AsObject}} [resource]
 * @return {Promise<{stream: *, value: OnOff.AsObject}>}
 */
export async function pullOnOff(deviceId, resource = {value: undefined, stream: undefined}) {
  const serverEndpoint = await grpcWebEndpoint();

  const api = new OnOffApiPromiseClient(serverEndpoint, null, null);
  if (resource.stream) resource.stream.cancel();
  const getResult = await api.getOnOff(new GetOnOffRequest().setName(deviceId));
  resource.value = getResult.toObject();
  const stream = api.pullOnOff(new PullOnOffRequest().setName(deviceId));
  resource.stream = stream;
  stream.on('data', res => {
    /** @type {PullOnOffResponse.Change[]} */
    const changes = res.getChangesList();
    for (const change of changes) {
      const value = change.getOnOff();
      resource.value = value.toObject();
    }
  });
  return resource;
}

/**
 * @param {string} deviceId
 * @param {OnOff.AsObject} onOff
 * @param {{loading: number}} [resource]
 * @return {Promise<OnOff.AsObject>}
 */
export async function updateOnOff(deviceId, onOff, resource = {loading: 0}) {
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
}
