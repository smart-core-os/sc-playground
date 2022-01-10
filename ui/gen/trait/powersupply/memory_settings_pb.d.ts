import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';
import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class MemorySettings extends jspb.Message {
  getRating(): number;

  setRating(value: number): MemorySettings;

  getVoltage(): number;

  setVoltage(value: number): MemorySettings;

  getLoad(): number;

  setLoad(value: number): MemorySettings;

  getReserved(): number;

  setReserved(value: number): MemorySettings;

  getMaxRampDuration(): google_protobuf_duration_pb.Duration | undefined;

  setMaxRampDuration(value?: google_protobuf_duration_pb.Duration): MemorySettings;

  hasMaxRampDuration(): boolean;

  clearMaxRampDuration(): MemorySettings;

  getDefaultRampDuration(): google_protobuf_duration_pb.Duration | undefined;

  setDefaultRampDuration(value?: google_protobuf_duration_pb.Duration): MemorySettings;

  hasDefaultRampDuration(): boolean;

  clearDefaultRampDuration(): MemorySettings;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): MemorySettings.AsObject;

  static toObject(includeInstance: boolean, msg: MemorySettings): MemorySettings.AsObject;

  static serializeBinaryToWriter(message: MemorySettings, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): MemorySettings;

  static deserializeBinaryFromReader(message: MemorySettings, reader: jspb.BinaryReader): MemorySettings;
}

export namespace MemorySettings {
  export type AsObject = {
    rating: number,
    voltage: number,
    load: number,
    reserved: number,
    maxRampDuration?: google_protobuf_duration_pb.Duration.AsObject,
    defaultRampDuration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class UpdateMemorySettingsReq extends jspb.Message {
  getName(): string;

  setName(value: string): UpdateMemorySettingsReq;

  getSettings(): MemorySettings | undefined;

  setSettings(value?: MemorySettings): UpdateMemorySettingsReq;

  hasSettings(): boolean;

  clearSettings(): UpdateMemorySettingsReq;

  getUpdateMask(): google_protobuf_field_mask_pb.FieldMask | undefined;

  setUpdateMask(value?: google_protobuf_field_mask_pb.FieldMask): UpdateMemorySettingsReq;

  hasUpdateMask(): boolean;

  clearUpdateMask(): UpdateMemorySettingsReq;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): UpdateMemorySettingsReq.AsObject;

  static toObject(includeInstance: boolean, msg: UpdateMemorySettingsReq): UpdateMemorySettingsReq.AsObject;

  static serializeBinaryToWriter(message: UpdateMemorySettingsReq, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): UpdateMemorySettingsReq;

  static deserializeBinaryFromReader(message: UpdateMemorySettingsReq, reader: jspb.BinaryReader): UpdateMemorySettingsReq;
}

export namespace UpdateMemorySettingsReq {
  export type AsObject = {
    name: string,
    settings?: MemorySettings.AsObject,
    updateMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class GetMemorySettingsReq extends jspb.Message {
  getName(): string;

  setName(value: string): GetMemorySettingsReq;

  getFields(): google_protobuf_field_mask_pb.FieldMask | undefined;

  setFields(value?: google_protobuf_field_mask_pb.FieldMask): GetMemorySettingsReq;

  hasFields(): boolean;

  clearFields(): GetMemorySettingsReq;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): GetMemorySettingsReq.AsObject;

  static toObject(includeInstance: boolean, msg: GetMemorySettingsReq): GetMemorySettingsReq.AsObject;

  static serializeBinaryToWriter(message: GetMemorySettingsReq, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): GetMemorySettingsReq;

  static deserializeBinaryFromReader(message: GetMemorySettingsReq, reader: jspb.BinaryReader): GetMemorySettingsReq;
}

export namespace GetMemorySettingsReq {
  export type AsObject = {
    name: string,
    fields?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullMemorySettingsReq extends jspb.Message {
  getName(): string;

  setName(value: string): PullMemorySettingsReq;

  getFields(): google_protobuf_field_mask_pb.FieldMask | undefined;

  setFields(value?: google_protobuf_field_mask_pb.FieldMask): PullMemorySettingsReq;

  hasFields(): boolean;

  clearFields(): PullMemorySettingsReq;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PullMemorySettingsReq.AsObject;

  static toObject(includeInstance: boolean, msg: PullMemorySettingsReq): PullMemorySettingsReq.AsObject;

  static serializeBinaryToWriter(message: PullMemorySettingsReq, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PullMemorySettingsReq;

  static deserializeBinaryFromReader(message: PullMemorySettingsReq, reader: jspb.BinaryReader): PullMemorySettingsReq;
}

export namespace PullMemorySettingsReq {
  export type AsObject = {
    name: string,
    fields?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullMemorySettingsRes extends jspb.Message {
  getChangesList(): Array<PullMemorySettingsRes.Change>;

  setChangesList(value: Array<PullMemorySettingsRes.Change>): PullMemorySettingsRes;

  clearChangesList(): PullMemorySettingsRes;

  addChanges(value?: PullMemorySettingsRes.Change, index?: number): PullMemorySettingsRes.Change;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): PullMemorySettingsRes.AsObject;

  static toObject(includeInstance: boolean, msg: PullMemorySettingsRes): PullMemorySettingsRes.AsObject;

  static serializeBinaryToWriter(message: PullMemorySettingsRes, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): PullMemorySettingsRes;

  static deserializeBinaryFromReader(message: PullMemorySettingsRes, reader: jspb.BinaryReader): PullMemorySettingsRes;
}

export namespace PullMemorySettingsRes {
  export type AsObject = {
    changesList: Array<PullMemorySettingsRes.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;

    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;

    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;

    hasChangeTime(): boolean;

    clearChangeTime(): Change;

    getSettings(): MemorySettings | undefined;

    setSettings(value?: MemorySettings): Change;

    hasSettings(): boolean;

    clearSettings(): Change;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): Change.AsObject;

    static toObject(includeInstance: boolean, msg: Change): Change.AsObject;

    static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): Change;

    static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
  }

  export namespace Change {
    export type AsObject = {
      name: string,
      changeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      settings?: MemorySettings.AsObject,
    }
  }

}

