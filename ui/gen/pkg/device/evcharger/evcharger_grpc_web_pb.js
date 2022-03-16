/**
 * @fileoverview gRPC-Web generated client stub for smartcore.playground.device.evcharger
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var github_com_smart$core$os_sc$api_protobuf_traits_electric_pb = require('@smart-core-os/sc-api-grpc-web/traits/electric_pb.js')

var github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb = require('@smart-core-os/sc-api-grpc-web/traits/energy_storage_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.playground = {};
proto.smartcore.playground.device = {};
proto.smartcore.playground.device.evcharger = require('./evcharger_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.playground.device.evcharger.EVChargerApiClient =
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
proto.smartcore.playground.device.evcharger.EVChargerApiPromiseClient =
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
 *   !proto.smartcore.playground.device.evcharger.PlugInRequest,
 *   !proto.smartcore.playground.device.evcharger.PlugInResponse>}
 */
const methodDescriptor_EVChargerApi_PlugIn = new grpc.web.MethodDescriptor(
    '/smartcore.playground.device.evcharger.EVChargerApi/PlugIn',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.device.evcharger.PlugInRequest,
    proto.smartcore.playground.device.evcharger.PlugInResponse,
    /**
     * @param {!proto.smartcore.playground.device.evcharger.PlugInRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.device.evcharger.PlugInResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.device.evcharger.PlugInRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.device.evcharger.PlugInResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.device.evcharger.PlugInResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.device.evcharger.EVChargerApiClient.prototype.plugIn =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/PlugIn',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_PlugIn,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.device.evcharger.PlugInRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.device.evcharger.PlugInResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.device.evcharger.EVChargerApiPromiseClient.prototype.plugIn =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/PlugIn',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_PlugIn);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.device.evcharger.ChargeStartRequest,
 *   !proto.smartcore.playground.device.evcharger.ChargeStartResponse>}
 */
const methodDescriptor_EVChargerApi_ChargeStart = new grpc.web.MethodDescriptor(
    '/smartcore.playground.device.evcharger.EVChargerApi/ChargeStart',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.device.evcharger.ChargeStartRequest,
    proto.smartcore.playground.device.evcharger.ChargeStartResponse,
    /**
     * @param {!proto.smartcore.playground.device.evcharger.ChargeStartRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.device.evcharger.ChargeStartResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.device.evcharger.ChargeStartRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.device.evcharger.ChargeStartResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.device.evcharger.ChargeStartResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.device.evcharger.EVChargerApiClient.prototype.chargeStart =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/ChargeStart',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_ChargeStart,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.device.evcharger.ChargeStartRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.device.evcharger.ChargeStartResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.device.evcharger.EVChargerApiPromiseClient.prototype.chargeStart =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/ChargeStart',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_ChargeStart);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.playground.device.evcharger.UnplugRequest,
 *   !proto.smartcore.playground.device.evcharger.UnplugResponse>}
 */
const methodDescriptor_EVChargerApi_Unplug = new grpc.web.MethodDescriptor(
    '/smartcore.playground.device.evcharger.EVChargerApi/Unplug',
    grpc.web.MethodType.UNARY,
    proto.smartcore.playground.device.evcharger.UnplugRequest,
    proto.smartcore.playground.device.evcharger.UnplugResponse,
    /**
     * @param {!proto.smartcore.playground.device.evcharger.UnplugRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.playground.device.evcharger.UnplugResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.playground.device.evcharger.UnplugRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.playground.device.evcharger.UnplugResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.playground.device.evcharger.UnplugResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.playground.device.evcharger.EVChargerApiClient.prototype.unplug =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/Unplug',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_Unplug,
          callback);
    };


/**
 * @param {!proto.smartcore.playground.device.evcharger.UnplugRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.playground.device.evcharger.UnplugResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.playground.device.evcharger.EVChargerApiPromiseClient.prototype.unplug =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.playground.device.evcharger.EVChargerApi/Unplug',
          request,
          metadata || {},
          methodDescriptor_EVChargerApi_Unplug);
    };


module.exports = proto.smartcore.playground.device.evcharger;

