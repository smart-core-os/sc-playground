/**
 * @fileoverview gRPC-Web generated client stub for smartcore.go.trait.powersupply
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js')

var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.go = {};
proto.smartcore.go.trait = {};
proto.smartcore.go.trait.powersupply = require('./memory_settings_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options['format'] = 'text';

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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiPromiseClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options['format'] = 'text';

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
 *   !proto.smartcore.go.trait.powersupply.GetMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.MemorySettings>}
 */
const methodDescriptor_MemorySettingsApi_GetSettings = new grpc.web.MethodDescriptor(
    '/smartcore.go.trait.powersupply.MemorySettingsApi/GetSettings',
    grpc.web.MethodType.UNARY,
    proto.smartcore.go.trait.powersupply.GetMemorySettingsReq,
    proto.smartcore.go.trait.powersupply.MemorySettings,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.GetMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.MemorySettings.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.go.trait.powersupply.GetMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.MemorySettings>}
 */
const methodInfo_MemorySettingsApi_GetSettings = new grpc.web.AbstractClientBase.MethodInfo(
    proto.smartcore.go.trait.powersupply.MemorySettings,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.GetMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.MemorySettings.deserializeBinary
);


/**
 * @param {!proto.smartcore.go.trait.powersupply.GetMemorySettingsReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.go.trait.powersupply.MemorySettings)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.go.trait.powersupply.MemorySettings>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiClient.prototype.getSettings =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/GetSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_GetSettings,
          callback);
    };


/**
 * @param {!proto.smartcore.go.trait.powersupply.GetMemorySettingsReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.go.trait.powersupply.MemorySettings>}
 *     Promise that resolves to the response
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiPromiseClient.prototype.getSettings =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/GetSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_GetSettings);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.MemorySettings>}
 */
const methodDescriptor_MemorySettingsApi_UpdateSettings = new grpc.web.MethodDescriptor(
    '/smartcore.go.trait.powersupply.MemorySettingsApi/UpdateSettings',
    grpc.web.MethodType.UNARY,
    proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq,
    proto.smartcore.go.trait.powersupply.MemorySettings,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.MemorySettings.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.MemorySettings>}
 */
const methodInfo_MemorySettingsApi_UpdateSettings = new grpc.web.AbstractClientBase.MethodInfo(
    proto.smartcore.go.trait.powersupply.MemorySettings,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.MemorySettings.deserializeBinary
);


/**
 * @param {!proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.go.trait.powersupply.MemorySettings)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.go.trait.powersupply.MemorySettings>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiClient.prototype.updateSettings =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/UpdateSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_UpdateSettings,
          callback);
    };


/**
 * @param {!proto.smartcore.go.trait.powersupply.UpdateMemorySettingsReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.go.trait.powersupply.MemorySettings>}
 *     Promise that resolves to the response
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiPromiseClient.prototype.updateSettings =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/UpdateSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_UpdateSettings);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.go.trait.powersupply.PullMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.PullMemorySettingsRes>}
 */
const methodDescriptor_MemorySettingsApi_PullSettings = new grpc.web.MethodDescriptor(
    '/smartcore.go.trait.powersupply.MemorySettingsApi/PullSettings',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.smartcore.go.trait.powersupply.PullMemorySettingsReq,
    proto.smartcore.go.trait.powersupply.PullMemorySettingsRes,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.PullMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.PullMemorySettingsRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.go.trait.powersupply.PullMemorySettingsReq,
 *   !proto.smartcore.go.trait.powersupply.PullMemorySettingsRes>}
 */
const methodInfo_MemorySettingsApi_PullSettings = new grpc.web.AbstractClientBase.MethodInfo(
    proto.smartcore.go.trait.powersupply.PullMemorySettingsRes,
    /**
     * @param {!proto.smartcore.go.trait.powersupply.PullMemorySettingsReq} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.smartcore.go.trait.powersupply.PullMemorySettingsRes.deserializeBinary
);


/**
 * @param {!proto.smartcore.go.trait.powersupply.PullMemorySettingsReq} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.go.trait.powersupply.PullMemorySettingsRes>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiClient.prototype.pullSettings =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/PullSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_PullSettings);
    };


/**
 * @param {!proto.smartcore.go.trait.powersupply.PullMemorySettingsReq} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.go.trait.powersupply.PullMemorySettingsRes>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.go.trait.powersupply.MemorySettingsApiPromiseClient.prototype.pullSettings =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/smartcore.go.trait.powersupply.MemorySettingsApi/PullSettings',
          request,
          metadata || {},
          methodDescriptor_MemorySettingsApi_PullSettings);
    };


module.exports = proto.smartcore.go.trait.powersupply;

