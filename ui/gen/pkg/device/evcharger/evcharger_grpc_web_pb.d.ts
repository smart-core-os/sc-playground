import * as grpcWeb from 'grpc-web';

import * as pkg_device_evcharger_evcharger_pb from '../../../pkg/device/evcharger/evcharger_pb';


export class EVChargerApiClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  plugIn(
      request: pkg_device_evcharger_evcharger_pb.PlugInRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_device_evcharger_evcharger_pb.PlugInResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_device_evcharger_evcharger_pb.PlugInResponse>;

  chargeStart(
      request: pkg_device_evcharger_evcharger_pb.ChargeStartRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_device_evcharger_evcharger_pb.ChargeStartResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_device_evcharger_evcharger_pb.ChargeStartResponse>;

  unplug(
      request: pkg_device_evcharger_evcharger_pb.UnplugRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_device_evcharger_evcharger_pb.UnplugResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_device_evcharger_evcharger_pb.UnplugResponse>;

}

export class EVChargerApiPromiseClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  plugIn(
      request: pkg_device_evcharger_evcharger_pb.PlugInRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_device_evcharger_evcharger_pb.PlugInResponse>;

  chargeStart(
      request: pkg_device_evcharger_evcharger_pb.ChargeStartRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_device_evcharger_evcharger_pb.ChargeStartResponse>;

  unplug(
      request: pkg_device_evcharger_evcharger_pb.UnplugRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_device_evcharger_evcharger_pb.UnplugResponse>;

}

