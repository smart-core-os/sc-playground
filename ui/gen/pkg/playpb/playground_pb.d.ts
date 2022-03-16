import * as jspb from 'google-protobuf'


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

