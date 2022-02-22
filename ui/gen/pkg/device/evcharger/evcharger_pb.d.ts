import * as jspb from 'google-protobuf'

import * as github_com_smart$core$os_sc$api_protobuf_traits_electric_pb
  from '@smart-core-os/sc-api-grpc-web/traits/electric_pb';
import * as github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb
  from '@smart-core-os/sc-api-grpc-web/traits/energy_storage_pb';


export class PlugInEvent extends jspb.Message {
  getFull(): github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity | undefined;

  setFull(value?: github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity): PlugInEvent;

  hasFull(): boolean;

  clearFull(): PlugInEvent;

  getLevel(): github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity | undefined;

  setLevel(value?: github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity): PlugInEvent;

  hasLevel(): boolean;

  clearLevel(): PlugInEvent;

  getModesList(): Array<PlugInEvent.ChargeMode>;

  setModesList(value: Array<PlugInEvent.ChargeMode>): PlugInEvent;

  clearModesList(): PlugInEvent;

  addModes(value?: PlugInEvent.ChargeMode, index?: number): PlugInEvent.ChargeMode;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PlugInEvent.AsObject;

  static toObject(includeInstance: boolean, msg: PlugInEvent): PlugInEvent.AsObject;

  static serializeBinaryToWriter(message: PlugInEvent, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PlugInEvent;

  static deserializeBinaryFromReader(message: PlugInEvent, reader: jspb.BinaryReader): PlugInEvent;
}

export namespace PlugInEvent {
  export type AsObject = {
    full?: github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity.AsObject,
    level?: github_com_smart$core$os_sc$api_protobuf_traits_energy_storage_pb.EnergyLevel.Quantity.AsObject,
    modesList: Array<PlugInEvent.ChargeMode.AsObject>,
  }

  export class ChargeMode extends jspb.Message {
    getId(): string;

    setId(value: string): ChargeMode;

    getTitle(): string;

    setTitle(value: string): ChargeMode;

    getDescription(): string;

    setDescription(value: string): ChargeMode;

    getSegmentsList(): Array<github_com_smart$core$os_sc$api_protobuf_traits_electric_pb.ElectricMode.Segment>;

    setSegmentsList(value: Array<github_com_smart$core$os_sc$api_protobuf_traits_electric_pb.ElectricMode.Segment>): ChargeMode;

    clearSegmentsList(): ChargeMode;

    addSegments(value?: github_com_smart$core$os_sc$api_protobuf_traits_electric_pb.ElectricMode.Segment, index?: number): github_com_smart$core$os_sc$api_protobuf_traits_electric_pb.ElectricMode.Segment;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): ChargeMode.AsObject;

    static toObject(includeInstance: boolean, msg: ChargeMode): ChargeMode.AsObject;

    static serializeBinaryToWriter(message: ChargeMode, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): ChargeMode;

    static deserializeBinaryFromReader(message: ChargeMode, reader: jspb.BinaryReader): ChargeMode;
  }

  export namespace ChargeMode {
    export type AsObject = {
      id: string,
      title: string,
      description: string,
      segmentsList: Array<github_com_smart$core$os_sc$api_protobuf_traits_electric_pb.ElectricMode.Segment.AsObject>,
    }
  }

}

export class PlugInRequest extends jspb.Message {
  getName(): string;

  setName(value: string): PlugInRequest;

  getEvent(): PlugInEvent | undefined;

  setEvent(value?: PlugInEvent): PlugInRequest;

  hasEvent(): boolean;

  clearEvent(): PlugInRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PlugInRequest.AsObject;

  static toObject(includeInstance: boolean, msg: PlugInRequest): PlugInRequest.AsObject;

  static serializeBinaryToWriter(message: PlugInRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PlugInRequest;

  static deserializeBinaryFromReader(message: PlugInRequest, reader: jspb.BinaryReader): PlugInRequest;
}

export namespace PlugInRequest {
  export type AsObject = {
    name: string,
    event?: PlugInEvent.AsObject,
  }
}

export class PlugInResponse extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PlugInResponse.AsObject;

  static toObject(includeInstance: boolean, msg: PlugInResponse): PlugInResponse.AsObject;

  static serializeBinaryToWriter(message: PlugInResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PlugInResponse;

  static deserializeBinaryFromReader(message: PlugInResponse, reader: jspb.BinaryReader): PlugInResponse;
}

export namespace PlugInResponse {
  export type AsObject = {}
}

export class ChargeStartRequest extends jspb.Message {
  getName(): string;

  setName(value: string): ChargeStartRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ChargeStartRequest.AsObject;

  static toObject(includeInstance: boolean, msg: ChargeStartRequest): ChargeStartRequest.AsObject;

  static serializeBinaryToWriter(message: ChargeStartRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ChargeStartRequest;

  static deserializeBinaryFromReader(message: ChargeStartRequest, reader: jspb.BinaryReader): ChargeStartRequest;
}

export namespace ChargeStartRequest {
  export type AsObject = {
    name: string,
  }
}

export class ChargeStartResponse extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ChargeStartResponse.AsObject;

  static toObject(includeInstance: boolean, msg: ChargeStartResponse): ChargeStartResponse.AsObject;

  static serializeBinaryToWriter(message: ChargeStartResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ChargeStartResponse;

  static deserializeBinaryFromReader(message: ChargeStartResponse, reader: jspb.BinaryReader): ChargeStartResponse;
}

export namespace ChargeStartResponse {
  export type AsObject = {}
}

export class UnplugRequest extends jspb.Message {
  getName(): string;

  setName(value: string): UnplugRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): UnplugRequest.AsObject;

  static toObject(includeInstance: boolean, msg: UnplugRequest): UnplugRequest.AsObject;

  static serializeBinaryToWriter(message: UnplugRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): UnplugRequest;

  static deserializeBinaryFromReader(message: UnplugRequest, reader: jspb.BinaryReader): UnplugRequest;
}

export namespace UnplugRequest {
  export type AsObject = {
    name: string,
  }
}

export class UnplugResponse extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): UnplugResponse.AsObject;

  static toObject(includeInstance: boolean, msg: UnplugResponse): UnplugResponse.AsObject;

  static serializeBinaryToWriter(message: UnplugResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): UnplugResponse;

  static deserializeBinaryFromReader(message: UnplugResponse, reader: jspb.BinaryReader): UnplugResponse;
}

export namespace UnplugResponse {
  export type AsObject = {}
}

