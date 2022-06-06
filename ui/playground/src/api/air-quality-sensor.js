import {AirQualitySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_grpc_web_pb.js';
import {PullAirQualityRequest} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb.js';
import {pullResource, setValue} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceValue<AirQuality.AsObject>} [resource]
 * @returns {ResourceValue<AirQuality.AsObject>}
 */
export function pullAirQualitySensor(name, resource = undefined) {
  return pullResource('AirQualitySensor', resource, endpoint => {
    const api = new AirQualitySensorApiPromiseClient(endpoint);
    const stream = api.pullAirQuality(new PullAirQualityRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAirQuality().toObject());
      }
    });
    return stream;
  });
}
