import {PublicationApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/publication_grpc_web_pb.js';
import {
  AcknowledgePublicationRequest,
  CreatePublicationRequest,
  Publication,
  PullPublicationsRequest,
  UpdatePublicationRequest
} from '@smart-core-os/sc-api-grpc-web/traits/publication_pb.js';
import {pullResource, setCollection} from './resource.js';
import {grpcWebEndpoint} from './config.js';
import Vue from 'vue';

/**
 * @param {string} name
 * @param {ResourceCollection<Publication.AsObject>} [resource]
 * @returns {ResourceCollection<Publication.AsObject>}
 */
export function pullPublications(name, resource = null) {
  return pullResource('Publication', resource, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint);
    const stream = api.pullPublications(new PullPublicationsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resource, change, v => v.id);
      }
    });
    return stream;
  });
}

/**
 * @param {string} name
 * @param {Publication.AsObject} publication
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function createPublication(name, publication, tracker) {
  tracker = tracker ?? {};
  Vue.set(tracker, 'loading', true);
  const endpoint = await grpcWebEndpoint();
  const api = new PublicationApiPromiseClient(endpoint);
  try {
    const response = await api.createPublication(new CreatePublicationRequest()
        .setName(name)
        .setPublication(fromObject(publication)))
    const responseObj = response.toObject();
    Vue.set(tracker, 'response', responseObj);
    return responseObj;
  } catch (err) {
    Vue.set(tracker, 'error', err);
    throw err;
  } finally {
    Vue.set(tracker, 'loading', false);
  }
}

/**
 * @param {string} name
 * @param {Publication.AsObject} publication
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function updatePublication(name, publication, tracker) {
  tracker = tracker ?? {};
  Vue.set(tracker, 'loading', true);
  const endpoint = await grpcWebEndpoint();
  const api = new PublicationApiPromiseClient(endpoint);
  try {
    const response = await api.updatePublication(new UpdatePublicationRequest()
        .setName(name)
        .setVersion(publication.version)
        .setPublication(fromObject(publication)))
    const responseObj = response.toObject();
    Vue.set(tracker, 'response', responseObj);
    return responseObj;
  } catch (err) {
    Vue.set(tracker, 'error', err);
    throw err;
  } finally {
    Vue.set(tracker, 'loading', false);
  }
}

/**
 * @param {AcknowledgePublicationRequest.AsObject} request
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function acknowledgePublication(request, tracker) {
  tracker = tracker ?? {};
  Vue.set(tracker, 'loading', true);
  const endpoint = await grpcWebEndpoint();
  const api = new PublicationApiPromiseClient(endpoint);
  try {
    const response = await api.acknowledgePublication(new AcknowledgePublicationRequest()
        .setName(request.name)
        .setId(request.id)
        .setVersion(request.version)
        .setReceipt(request.receipt ?? Publication.Audience.Receipt.ACCEPTED)
        .setReceiptRejectedReason(request.receiptRejectedReason ?? '')
        .setAllowAcknowledged(request.allowAcknowledged ?? false));
    const responseObj = response.toObject();
    Vue.set(tracker, 'response', responseObj);
    return responseObj;
  } catch (err) {
    Vue.set(tracker, 'error', err);
    throw err;
  } finally {
    Vue.set(tracker, 'loading', false);
  }
}

/**
 * @param {Publication.AsObject} obj
 * @returns {Publication}
 */
export function fromObject(obj) {
  const publication = new Publication();
  for (const prop of ['id', 'body', 'mediaType', 'version']) {
    if (obj.hasOwnProperty(prop)) {
      publication['set' + prop[0].toUpperCase() + prop.substring(1)](obj[prop]);
    }
  }

  if (obj.hasOwnProperty('audience')) {
    const src = obj.audience;
    const audience = new Publication.Audience();
    publication.setAudience(audience);
    for (const prop of ['name', 'receipt', 'receiptRejectedReason']) {
      if (src.hasOwnProperty(prop)) {
        audience['set' + prop[0].toUpperCase() + prop.substring(1)](src[prop]);
      }
    }
  }

  return publication;
}
