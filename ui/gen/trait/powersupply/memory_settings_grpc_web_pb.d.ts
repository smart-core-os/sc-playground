import * as grpcWeb from 'grpc-web';

import * as trait_powersupply_memory_settings_pb from '../../trait/powersupply/memory_settings_pb';


export class MemorySettingsApiClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  getSettings(
      request: trait_powersupply_memory_settings_pb.GetMemorySettingsReq,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.Error,
                 response: trait_powersupply_memory_settings_pb.MemorySettings) => void
  ): grpcWeb.ClientReadableStream<trait_powersupply_memory_settings_pb.MemorySettings>;

  updateSettings(
      request: trait_powersupply_memory_settings_pb.UpdateMemorySettingsReq,
      metadata: grpcWeb.Metadata | undefined,
      callback: (err: grpcWeb.Error,
                 response: trait_powersupply_memory_settings_pb.MemorySettings) => void
  ): grpcWeb.ClientReadableStream<trait_powersupply_memory_settings_pb.MemorySettings>;

  pullSettings(
      request: trait_powersupply_memory_settings_pb.PullMemorySettingsReq,
      metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<trait_powersupply_memory_settings_pb.PullMemorySettingsRes>;

}

export class MemorySettingsApiPromiseClient {
  constructor(hostname: string,
              credentials?: null | { [index: string]: string; },
              options?: null | { [index: string]: any; });

  getSettings(
      request: trait_powersupply_memory_settings_pb.GetMemorySettingsReq,
      metadata?: grpcWeb.Metadata
  ): Promise<trait_powersupply_memory_settings_pb.MemorySettings>;

  updateSettings(
      request: trait_powersupply_memory_settings_pb.UpdateMemorySettingsReq,
      metadata?: grpcWeb.Metadata
  ): Promise<trait_powersupply_memory_settings_pb.MemorySettings>;

  pullSettings(
      request: trait_powersupply_memory_settings_pb.PullMemorySettingsReq,
      metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<trait_powersupply_memory_settings_pb.PullMemorySettingsRes>;

}

