import * as grpcWeb from "grpc-web"
import {Timestamp} from "google-protobuf/google/protobuf/timestamp_pb";
import {ChangeType} from "@smart-core-os/sc-api-grpc-web/types/change_pb";

export interface RemoteResource {
  loading?: boolean;
  stream?: grpcWeb.ClientReadableStream<any>;
  streamError?: Error;
  updateTime?: Date;
}

export interface ResourceValue<T> extends RemoteResource {
  value?: T;
}

export interface ResourceCollection<T> extends RemoteResource {
  value?: { [id: string]: T };
}

export interface ResourceCallback<T> {
  data(val: T);

  error(e: Error);
}

export interface StreamFactory<T> {
  (endpoint: string): grpcWeb.ClientReadableStream<T>;
}

export interface CollectionChange<T extends { toObject(): O }, O> {
  getName(): string;

  getChangeTime(): Timestamp;

  getChangeType(): ChangeType;

  getOldValue(): T | undefined;

  getNewValue(): T | undefined;
}

export interface ActionTracker<T> {
  loading?: boolean;
  error?: Error;
  response?: T;
  duration?: number;
}
