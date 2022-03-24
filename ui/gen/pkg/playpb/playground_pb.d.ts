import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class AddDeviceTraitRequest extends jspb.Message {
  getName(): string;

  setName(value: string): AddDeviceTraitRequest;

  getTraitName(): string;

  setTraitName(value: string): AddDeviceTraitRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDeviceTraitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddDeviceTraitRequest): AddDeviceTraitRequest.AsObject;
  static serializeBinaryToWriter(message: AddDeviceTraitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDeviceTraitRequest;
  static deserializeBinaryFromReader(message: AddDeviceTraitRequest, reader: jspb.BinaryReader): AddDeviceTraitRequest;
}

export namespace AddDeviceTraitRequest {
  export type AsObject = {
    name: string,
    traitName: string,
  }
}

export class AddDeviceTraitResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDeviceTraitResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddDeviceTraitResponse): AddDeviceTraitResponse.AsObject;
  static serializeBinaryToWriter(message: AddDeviceTraitResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDeviceTraitResponse;
  static deserializeBinaryFromReader(message: AddDeviceTraitResponse, reader: jspb.BinaryReader): AddDeviceTraitResponse;
}

export namespace AddDeviceTraitResponse {
  export type AsObject = {}
}

export class ListSupportedTraitsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSupportedTraitsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSupportedTraitsRequest): ListSupportedTraitsRequest.AsObject;
  static serializeBinaryToWriter(message: ListSupportedTraitsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSupportedTraitsRequest;
  static deserializeBinaryFromReader(message: ListSupportedTraitsRequest, reader: jspb.BinaryReader): ListSupportedTraitsRequest;
}

export namespace ListSupportedTraitsRequest {
  export type AsObject = {}
}

export class ListSupportedTraitsResponse extends jspb.Message {
  getTraitNameList(): Array<string>;
  setTraitNameList(value: Array<string>): ListSupportedTraitsResponse;
  clearTraitNameList(): ListSupportedTraitsResponse;
  addTraitName(value: string, index?: number): ListSupportedTraitsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSupportedTraitsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSupportedTraitsResponse): ListSupportedTraitsResponse.AsObject;
  static serializeBinaryToWriter(message: ListSupportedTraitsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSupportedTraitsResponse;
  static deserializeBinaryFromReader(message: ListSupportedTraitsResponse, reader: jspb.BinaryReader): ListSupportedTraitsResponse;
}

export namespace ListSupportedTraitsResponse {
  export type AsObject = {
    traitNameList: Array<string>,
  }
}

export class AddRemoteDeviceRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddRemoteDeviceRequest;

  getEndpoint(): string;
  setEndpoint(value: string): AddRemoteDeviceRequest;

  getTraitNameList(): Array<string>;
  setTraitNameList(value: Array<string>): AddRemoteDeviceRequest;
  clearTraitNameList(): AddRemoteDeviceRequest;
  addTraitName(value: string, index?: number): AddRemoteDeviceRequest;

  getTls(): RemoteTLS | undefined;
  setTls(value?: RemoteTLS): AddRemoteDeviceRequest;
  hasTls(): boolean;
  clearTls(): AddRemoteDeviceRequest;

  getInsecure(): boolean;
  setInsecure(value: boolean): AddRemoteDeviceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddRemoteDeviceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddRemoteDeviceRequest): AddRemoteDeviceRequest.AsObject;
  static serializeBinaryToWriter(message: AddRemoteDeviceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddRemoteDeviceRequest;
  static deserializeBinaryFromReader(message: AddRemoteDeviceRequest, reader: jspb.BinaryReader): AddRemoteDeviceRequest;
}

export namespace AddRemoteDeviceRequest {
  export type AsObject = {
    name: string,
    endpoint: string,
    traitNameList: Array<string>,
    tls?: RemoteTLS.AsObject,
    insecure: boolean,
  }
}

export class RemoteTLS extends jspb.Message {
  getServerCaCert(): string;
  setServerCaCert(value: string): RemoteTLS;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoteTLS.AsObject;
  static toObject(includeInstance: boolean, msg: RemoteTLS): RemoteTLS.AsObject;
  static serializeBinaryToWriter(message: RemoteTLS, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoteTLS;
  static deserializeBinaryFromReader(message: RemoteTLS, reader: jspb.BinaryReader): RemoteTLS;
}

export namespace RemoteTLS {
  export type AsObject = {
    serverCaCert: string,
  }
}

export class AddRemoteDeviceResponse extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): AddRemoteDeviceResponse.AsObject;

  static toObject(includeInstance: boolean, msg: AddRemoteDeviceResponse): AddRemoteDeviceResponse.AsObject;

  static serializeBinaryToWriter(message: AddRemoteDeviceResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): AddRemoteDeviceResponse;

  static deserializeBinaryFromReader(message: AddRemoteDeviceResponse, reader: jspb.BinaryReader): AddRemoteDeviceResponse;
}

export namespace AddRemoteDeviceResponse {
  export type AsObject = {}
}

export class Performance extends jspb.Message {
  getFrame(): google_protobuf_duration_pb.Duration | undefined;

  setFrame(value?: google_protobuf_duration_pb.Duration): Performance;

  hasFrame(): boolean;

  clearFrame(): Performance;

  getCapture(): google_protobuf_duration_pb.Duration | undefined;

  setCapture(value?: google_protobuf_duration_pb.Duration): Performance;

  hasCapture(): boolean;

  clearCapture(): Performance;

  getScrub(): google_protobuf_duration_pb.Duration | undefined;

  setScrub(value?: google_protobuf_duration_pb.Duration): Performance;

  hasScrub(): boolean;

  clearScrub(): Performance;

  getRespond(): google_protobuf_duration_pb.Duration | undefined;

  setRespond(value?: google_protobuf_duration_pb.Duration): Performance;

  hasRespond(): boolean;

  clearRespond(): Performance;

  getIdle(): google_protobuf_duration_pb.Duration | undefined;

  setIdle(value?: google_protobuf_duration_pb.Duration): Performance;

  hasIdle(): boolean;

  clearIdle(): Performance;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): Performance.AsObject;

  static toObject(includeInstance: boolean, msg: Performance): Performance.AsObject;

  static serializeBinaryToWriter(message: Performance, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): Performance;

  static deserializeBinaryFromReader(message: Performance, reader: jspb.BinaryReader): Performance;
}

export namespace Performance {
  export type AsObject = {
    frame?: google_protobuf_duration_pb.Duration.AsObject,
    capture?: google_protobuf_duration_pb.Duration.AsObject,
    scrub?: google_protobuf_duration_pb.Duration.AsObject,
    respond?: google_protobuf_duration_pb.Duration.AsObject,
    idle?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class PullPerformanceRequest extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PullPerformanceRequest.AsObject;

  static toObject(includeInstance: boolean, msg: PullPerformanceRequest): PullPerformanceRequest.AsObject;

  static serializeBinaryToWriter(message: PullPerformanceRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PullPerformanceRequest;

  static deserializeBinaryFromReader(message: PullPerformanceRequest, reader: jspb.BinaryReader): PullPerformanceRequest;
}

export namespace PullPerformanceRequest {
  export type AsObject = {}
}

export class PullPerformanceResponse extends jspb.Message {
  getPerformance(): Performance | undefined;

  setPerformance(value?: Performance): PullPerformanceResponse;

  hasPerformance(): boolean;

  clearPerformance(): PullPerformanceResponse;

  getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;

  setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): PullPerformanceResponse;

  hasChangeTime(): boolean;

  clearChangeTime(): PullPerformanceResponse;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PullPerformanceResponse.AsObject;

  static toObject(includeInstance: boolean, msg: PullPerformanceResponse): PullPerformanceResponse.AsObject;

  static serializeBinaryToWriter(message: PullPerformanceResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PullPerformanceResponse;

  static deserializeBinaryFromReader(message: PullPerformanceResponse, reader: jspb.BinaryReader): PullPerformanceResponse;
}

export namespace PullPerformanceResponse {
  export type AsObject = {
    performance?: Performance.AsObject,
    changeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

