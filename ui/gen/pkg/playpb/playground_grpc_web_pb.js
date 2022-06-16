/**
 * @fileoverview gRPC-Web generated client stub for smartcore.playground.api
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.playground = {};
proto.smartcore.playground.api = require('./playground_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.playground.api.PlaygroundApiClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.playground.api.PlaygroundApiPromiseClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.api.AddDeviceTraitRequest,
 *   !proto.smartcore.playground.api.AddDeviceTraitResponse>}
 */
const methodDescriptor_PlaygroundApi_AddDeviceTrait = new grpc.web.MethodDescriptor(
    '/smartcore.playground.api.PlaygroundApi/AddDeviceTrait',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.api.AddDeviceTraitRequest,
    proto.smartcore.playground.api.AddDeviceTraitResponse,
    /**
     * @param {!proto.smartcore.playground.api.AddDeviceTraitRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.api.AddDeviceTraitResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.api.AddDeviceTraitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.api.AddDeviceTraitResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.api.AddDeviceTraitResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.api.PlaygroundApiClient.prototype.addDeviceTrait =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/AddDeviceTrait',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_AddDeviceTrait,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.api.AddDeviceTraitRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.api.AddDeviceTraitResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.api.PlaygroundApiPromiseClient.prototype.addDeviceTrait =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/AddDeviceTrait',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_AddDeviceTrait);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.api.ListSupportedTraitsRequest,
 *   !proto.smartcore.playground.api.ListSupportedTraitsResponse>}
 */
const methodDescriptor_PlaygroundApi_ListSupportedTraits = new grpc.web.MethodDescriptor(
    '/smartcore.playground.api.PlaygroundApi/ListSupportedTraits',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.api.ListSupportedTraitsRequest,
    proto.smartcore.playground.api.ListSupportedTraitsResponse,
    /**
     * @param {!proto.smartcore.playground.api.ListSupportedTraitsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.api.ListSupportedTraitsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.api.ListSupportedTraitsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.api.ListSupportedTraitsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.api.ListSupportedTraitsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.api.PlaygroundApiClient.prototype.listSupportedTraits =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/ListSupportedTraits',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_ListSupportedTraits,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.api.ListSupportedTraitsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.api.ListSupportedTraitsResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.api.PlaygroundApiPromiseClient.prototype.listSupportedTraits =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/ListSupportedTraits',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_ListSupportedTraits);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.api.AddRemoteDeviceRequest,
 *   !proto.smartcore.playground.api.AddRemoteDeviceResponse>}
 */
const methodDescriptor_PlaygroundApi_AddRemoteDevice = new grpc.web.MethodDescriptor(
    '/smartcore.playground.api.PlaygroundApi/AddRemoteDevice',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.api.AddRemoteDeviceRequest,
    proto.smartcore.playground.api.AddRemoteDeviceResponse,
    /**
     * @param {!proto.smartcore.playground.api.AddRemoteDeviceRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.api.AddRemoteDeviceResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.api.AddRemoteDeviceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.api.AddRemoteDeviceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.api.AddRemoteDeviceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.api.PlaygroundApiClient.prototype.addRemoteDevice =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/AddRemoteDevice',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_AddRemoteDevice,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.api.AddRemoteDeviceRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.api.AddRemoteDeviceResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.api.PlaygroundApiPromiseClient.prototype.addRemoteDevice =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/AddRemoteDevice',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_AddRemoteDevice);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.api.PullPerformanceRequest,
 *   !proto.smartcore.playground.api.PullPerformanceResponse>}
 */
const methodDescriptor_PlaygroundApi_PullPerformance = new grpc.web.MethodDescriptor(
    '/smartcore.playground.api.PlaygroundApi/PullPerformance',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.smartcore.playground.api.PullPerformanceRequest,
    proto.smartcore.playground.api.PullPerformanceResponse,
    /**
     * @param {!proto.smartcore.playground.api.PullPerformanceRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.api.PullPerformanceResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.api.PullPerformanceRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.api.PullPerformanceResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.api.PlaygroundApiClient.prototype.pullPerformance =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/PullPerformance',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_PullPerformance);
    };


/**
 * @param {!proto.smartcore.playground.api.PullPerformanceRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.api.PullPerformanceResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.api.PlaygroundApiPromiseClient.prototype.pullPerformance =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/smartcore.playground.api.PlaygroundApi/PullPerformance',
          request,
          metadata || {},
          methodDescriptor_PlaygroundApi_PullPerformance);
    };


module.exports = proto.smartcore.playground.api;

