import * as grpcWeb from 'grpc-web';

import * as pkg_playpb_playground_pb from '../../pkg/playpb/playground_pb';


export class PlaygroundApiClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  addDeviceTrait(
      request: pkg_playpb_playground_pb.AddDeviceTraitRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_playpb_playground_pb.AddDeviceTraitResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_playpb_playground_pb.AddDeviceTraitResponse>;

  listSupportedTraits(
      request: pkg_playpb_playground_pb.ListSupportedTraitsRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_playpb_playground_pb.ListSupportedTraitsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_playpb_playground_pb.ListSupportedTraitsResponse>;

  addRemoteDevice(
      request: pkg_playpb_playground_pb.AddRemoteDeviceRequest,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.RpcError,
                 response: pkg_playpb_playground_pb.AddRemoteDeviceResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_playpb_playground_pb.AddRemoteDeviceResponse>;

  pullPerformance(
      request: pkg_playpb_playground_pb.PullPerformanceRequest,
      metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<pkg_playpb_playground_pb.PullPerformanceResponse>;

}

export class PlaygroundApiPromiseClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  addDeviceTrait(
      request: pkg_playpb_playground_pb.AddDeviceTraitRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_playpb_playground_pb.AddDeviceTraitResponse>;

  listSupportedTraits(
      request: pkg_playpb_playground_pb.ListSupportedTraitsRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_playpb_playground_pb.ListSupportedTraitsResponse>;

  addRemoteDevice(
      request: pkg_playpb_playground_pb.AddRemoteDeviceRequest,
      metadata?: grpcWeb.Metadata
  ): Promise<pkg_playpb_playground_pb.AddRemoteDeviceResponse>;

  pullPerformance(
      request: pkg_playpb_playground_pb.PullPerformanceRequest,
      metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<pkg_playpb_playground_pb.PullPerformanceResponse>;

}

