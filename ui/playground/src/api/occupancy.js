import {OccupancySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_grpc_web_pb';
import {PullOccupancyRequest} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {pullResource, setValue} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceValue<Occupancy.AsObject>} [resource]
 * @returns {ResourceValue<Occupancy.AsObject>}
 */
export function pullOccupancy(name, resource = null) {
  return pullResource('OccupancySensor.Occupancy', resource, endpoint => {
    const api = new OccupancySensorApiPromiseClient(endpoint);
    const stream = api.pullOccupancy(new PullOccupancyRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getOccupancy().toObject());
      }
    });
    return stream;
  });
}
