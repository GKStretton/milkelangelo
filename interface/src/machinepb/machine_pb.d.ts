// package: machine
// file: machine.proto

import * as jspb from "google-protobuf";

export class PipetteState extends jspb.Message {
  getSpent(): boolean;
  setSpent(value: boolean): void;

  getVialHeld(): number;
  setVialHeld(value: number): void;

  getVolumeTargetUl(): number;
  setVolumeTargetUl(value: number): void;

  getDispenseRequestNumber(): number;
  setDispenseRequestNumber(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PipetteState.AsObject;
  static toObject(includeInstance: boolean, msg: PipetteState): PipetteState.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PipetteState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PipetteState;
  static deserializeBinaryFromReader(message: PipetteState, reader: jspb.BinaryReader): PipetteState;
}

export namespace PipetteState {
  export type AsObject = {
    spent: boolean,
    vialHeld: number,
    volumeTargetUl: number,
    dispenseRequestNumber: number,
  }
}

export class CollectionRequest extends jspb.Message {
  getCompleted(): boolean;
  setCompleted(value: boolean): void;

  getRequestNumber(): number;
  setRequestNumber(value: number): void;

  getVialNumber(): number;
  setVialNumber(value: number): void;

  getVolumeUl(): number;
  setVolumeUl(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CollectionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CollectionRequest): CollectionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CollectionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CollectionRequest;
  static deserializeBinaryFromReader(message: CollectionRequest, reader: jspb.BinaryReader): CollectionRequest;
}

export namespace CollectionRequest {
  export type AsObject = {
    completed: boolean,
    requestNumber: number,
    vialNumber: number,
    volumeUl: number,
  }
}

export class MovementDetails extends jspb.Message {
  getTargetXUnit(): number;
  setTargetXUnit(value: number): void;

  getTargetYUnit(): number;
  setTargetYUnit(value: number): void;

  getTargetZIk(): number;
  setTargetZIk(value: number): void;

  getTargetRingDeg(): number;
  setTargetRingDeg(value: number): void;

  getTargetYawDeg(): number;
  setTargetYawDeg(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MovementDetails.AsObject;
  static toObject(includeInstance: boolean, msg: MovementDetails): MovementDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MovementDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MovementDetails;
  static deserializeBinaryFromReader(message: MovementDetails, reader: jspb.BinaryReader): MovementDetails;
}

export namespace MovementDetails {
  export type AsObject = {
    targetXUnit: number,
    targetYUnit: number,
    targetZIk: number,
    targetRingDeg: number,
    targetYawDeg: number,
  }
}

export class FluidRequest extends jspb.Message {
  getFluidtype(): FluidTypeMap[keyof FluidTypeMap];
  setFluidtype(value: FluidTypeMap[keyof FluidTypeMap]): void;

  getVolumeMl(): number;
  setVolumeMl(value: number): void;

  getComplete(): boolean;
  setComplete(value: boolean): void;

  getOpenDrain(): boolean;
  setOpenDrain(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FluidRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FluidRequest): FluidRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FluidRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FluidRequest;
  static deserializeBinaryFromReader(message: FluidRequest, reader: jspb.BinaryReader): FluidRequest;
}

export namespace FluidRequest {
  export type AsObject = {
    fluidtype: FluidTypeMap[keyof FluidTypeMap],
    volumeMl: number,
    complete: boolean,
    openDrain: boolean,
  }
}

export class FluidDetails extends jspb.Message {
  getBowlFluidLevelMl(): number;
  setBowlFluidLevelMl(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FluidDetails.AsObject;
  static toObject(includeInstance: boolean, msg: FluidDetails): FluidDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FluidDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FluidDetails;
  static deserializeBinaryFromReader(message: FluidDetails, reader: jspb.BinaryReader): FluidDetails;
}

export namespace FluidDetails {
  export type AsObject = {
    bowlFluidLevelMl: number,
  }
}

export class StateReport extends jspb.Message {
  getTimestampUnixMicros(): number;
  setTimestampUnixMicros(value: number): void;

  getStartupCounter(): number;
  setStartupCounter(value: number): void;

  getMode(): ModeMap[keyof ModeMap];
  setMode(value: ModeMap[keyof ModeMap]): void;

  getStatus(): StatusMap[keyof StatusMap];
  setStatus(value: StatusMap[keyof StatusMap]): void;

  getLightsOn(): boolean;
  setLightsOn(value: boolean): void;

  hasPipetteState(): boolean;
  clearPipetteState(): void;
  getPipetteState(): PipetteState | undefined;
  setPipetteState(value?: PipetteState): void;

  hasCollectionRequest(): boolean;
  clearCollectionRequest(): void;
  getCollectionRequest(): CollectionRequest | undefined;
  setCollectionRequest(value?: CollectionRequest): void;

  hasMovementDetails(): boolean;
  clearMovementDetails(): void;
  getMovementDetails(): MovementDetails | undefined;
  setMovementDetails(value?: MovementDetails): void;

  hasFluidRequest(): boolean;
  clearFluidRequest(): void;
  getFluidRequest(): FluidRequest | undefined;
  setFluidRequest(value?: FluidRequest): void;

  hasFluidDetails(): boolean;
  clearFluidDetails(): void;
  getFluidDetails(): FluidDetails | undefined;
  setFluidDetails(value?: FluidDetails): void;

  getPaused(): boolean;
  setPaused(value: boolean): void;

  getTimestampReadable(): string;
  setTimestampReadable(value: string): void;

  getLatestDslrFileNumber(): number;
  setLatestDslrFileNumber(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateReport.AsObject;
  static toObject(includeInstance: boolean, msg: StateReport): StateReport.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StateReport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateReport;
  static deserializeBinaryFromReader(message: StateReport, reader: jspb.BinaryReader): StateReport;
}

export namespace StateReport {
  export type AsObject = {
    timestampUnixMicros: number,
    startupCounter: number,
    mode: ModeMap[keyof ModeMap],
    status: StatusMap[keyof StatusMap],
    lightsOn: boolean,
    pipetteState?: PipetteState.AsObject,
    collectionRequest?: CollectionRequest.AsObject,
    movementDetails?: MovementDetails.AsObject,
    fluidRequest?: FluidRequest.AsObject,
    fluidDetails?: FluidDetails.AsObject,
    paused: boolean,
    timestampReadable: string,
    latestDslrFileNumber: number,
  }
}

export class StateReportList extends jspb.Message {
  clearStatereportsList(): void;
  getStatereportsList(): Array<StateReport>;
  setStatereportsList(value: Array<StateReport>): void;
  addStatereports(value?: StateReport, index?: number): StateReport;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateReportList.AsObject;
  static toObject(includeInstance: boolean, msg: StateReportList): StateReportList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StateReportList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateReportList;
  static deserializeBinaryFromReader(message: StateReportList, reader: jspb.BinaryReader): StateReportList;
}

export namespace StateReportList {
  export type AsObject = {
    statereportsList: Array<StateReport.AsObject>,
  }
}

export class SessionStatus extends jspb.Message {
  getId(): number;
  setId(value: number): void;

  getPaused(): boolean;
  setPaused(value: boolean): void;

  getComplete(): boolean;
  setComplete(value: boolean): void;

  getProduction(): boolean;
  setProduction(value: boolean): void;

  getProductionId(): number;
  setProductionId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionStatus.AsObject;
  static toObject(includeInstance: boolean, msg: SessionStatus): SessionStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SessionStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionStatus;
  static deserializeBinaryFromReader(message: SessionStatus, reader: jspb.BinaryReader): SessionStatus;
}

export namespace SessionStatus {
  export type AsObject = {
    id: number,
    paused: boolean,
    complete: boolean,
    production: boolean,
    productionId: number,
  }
}

export class StreamStatus extends jspb.Message {
  getLive(): boolean;
  setLive(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamStatus.AsObject;
  static toObject(includeInstance: boolean, msg: StreamStatus): StreamStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StreamStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamStatus;
  static deserializeBinaryFromReader(message: StreamStatus, reader: jspb.BinaryReader): StreamStatus;
}

export namespace StreamStatus {
  export type AsObject = {
    live: boolean,
  }
}

export class DispenseMetadataMap extends jspb.Message {
  getDispenseMetadataMap(): jspb.Map<string, DispenseMetadata>;
  clearDispenseMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DispenseMetadataMap.AsObject;
  static toObject(includeInstance: boolean, msg: DispenseMetadataMap): DispenseMetadataMap.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DispenseMetadataMap, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DispenseMetadataMap;
  static deserializeBinaryFromReader(message: DispenseMetadataMap, reader: jspb.BinaryReader): DispenseMetadataMap;
}

export namespace DispenseMetadataMap {
  export type AsObject = {
    dispenseMetadataMap: Array<[string, DispenseMetadata.AsObject]>,
  }
}

export class DispenseMetadata extends jspb.Message {
  getFaileddispense(): boolean;
  setFaileddispense(value: boolean): void;

  getDispensedelayms(): number;
  setDispensedelayms(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DispenseMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: DispenseMetadata): DispenseMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DispenseMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DispenseMetadata;
  static deserializeBinaryFromReader(message: DispenseMetadata, reader: jspb.BinaryReader): DispenseMetadata;
}

export namespace DispenseMetadata {
  export type AsObject = {
    faileddispense: boolean,
    dispensedelayms: number,
  }
}

export class ContentTypeStatuses extends jspb.Message {
  getContentStatusesMap(): jspb.Map<string, ContentTypeStatus>;
  clearContentStatusesMap(): void;
  getSplashtext(): string;
  setSplashtext(value: string): void;

  getSplashtextHue(): number;
  setSplashtextHue(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ContentTypeStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: ContentTypeStatuses): ContentTypeStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ContentTypeStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ContentTypeStatuses;
  static deserializeBinaryFromReader(message: ContentTypeStatuses, reader: jspb.BinaryReader): ContentTypeStatuses;
}

export namespace ContentTypeStatuses {
  export type AsObject = {
    contentStatusesMap: Array<[string, ContentTypeStatus.AsObject]>,
    splashtext: string,
    splashtextHue: number,
  }
}

export class ContentTypeStatus extends jspb.Message {
  getRawTitle(): string;
  setRawTitle(value: string): void;

  getRawDescription(): string;
  setRawDescription(value: string): void;

  getCaption(): string;
  setCaption(value: string): void;

  clearPostsList(): void;
  getPostsList(): Array<Post>;
  setPostsList(value: Array<Post>): void;
  addPosts(value?: Post, index?: number): Post;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ContentTypeStatus.AsObject;
  static toObject(includeInstance: boolean, msg: ContentTypeStatus): ContentTypeStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ContentTypeStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ContentTypeStatus;
  static deserializeBinaryFromReader(message: ContentTypeStatus, reader: jspb.BinaryReader): ContentTypeStatus;
}

export namespace ContentTypeStatus {
  export type AsObject = {
    rawTitle: string,
    rawDescription: string,
    caption: string,
    postsList: Array<Post.AsObject>,
  }
}

export class Post extends jspb.Message {
  getPlatform(): SocialPlatformMap[keyof SocialPlatformMap];
  setPlatform(value: SocialPlatformMap[keyof SocialPlatformMap]): void;

  getSubPlatform(): string;
  setSubPlatform(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getUploaded(): boolean;
  setUploaded(value: boolean): void;

  getUrl(): string;
  setUrl(value: string): void;

  getCrosspost(): boolean;
  setCrosspost(value: boolean): void;

  getScheduledUnixTimetamp(): number;
  setScheduledUnixTimetamp(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Post.AsObject;
  static toObject(includeInstance: boolean, msg: Post): Post.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Post, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Post;
  static deserializeBinaryFromReader(message: Post, reader: jspb.BinaryReader): Post;
}

export namespace Post {
  export type AsObject = {
    platform: SocialPlatformMap[keyof SocialPlatformMap],
    subPlatform: string,
    title: string,
    description: string,
    uploaded: boolean,
    url: string,
    crosspost: boolean,
    scheduledUnixTimetamp: number,
  }
}

export class Email extends jspb.Message {
  getSubject(): string;
  setSubject(value: string): void;

  getBody(): string;
  setBody(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Email.AsObject;
  static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Email;
  static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
}

export namespace Email {
  export type AsObject = {
    subject: string,
    body: string,
  }
}

export interface NodeMap {
  UNDEFINED: 0;
  HOME: 4;
  HOME_TOP: 8;
  VIAL_1_ABOVE: 10;
  MIN_VIAL_ABOVE: 10;
  VIAL_1_INSIDE: 15;
  MIN_VIAL_INSIDE: 15;
  VIAL_2_ABOVE: 20;
  VIAL_2_INSIDE: 25;
  VIAL_3_ABOVE: 30;
  VIAL_3_INSIDE: 35;
  VIAL_4_ABOVE: 40;
  VIAL_4_INSIDE: 45;
  VIAL_5_ABOVE: 50;
  VIAL_5_INSIDE: 55;
  VIAL_6_ABOVE: 60;
  VIAL_6_INSIDE: 65;
  VIAL_7_ABOVE: 70;
  MAX_VIAL_ABOVE: 70;
  VIAL_7_INSIDE: 75;
  MAX_VIAL_INSIDE: 75;
  LOW_ENTRY_POINT: 30;
  RINSE_CONTAINER_ENTRY: 80;
  RINSE_CONTAINER_LOW: 85;
  OUTER_HANDOVER: 90;
  INNER_HANDOVER: 110;
  INVERSE_KINEMATICS_POSITION: 150;
  IDLE_LOCATION: 8;
}

export const Node: NodeMap;

export interface SolenoidValveMap {
  VALVE_UNDEFINED: 0;
  VALVE_DRAIN: 1;
  VALVE_WATER: 2;
  VALVE_MILK: 3;
  VALVE_AIR: 4;
}

export const SolenoidValve: SolenoidValveMap;

export interface ModeMap {
  UNDEFINED_MODE: 0;
  MANUAL: 1;
  AUTONOMOUS: 2;
}

export const Mode: ModeMap;

export interface StatusMap {
  UNDEFINED_STATUS: 0;
  ERROR: 1;
  E_STOP_ACTIVE: 5;
  SLEEPING: 6;
  SHUTTING_DOWN: 9;
  WAKING_UP: 10;
  CALIBRATING: 20;
  IDLE_STATIONARY: 30;
  IDLE_MOVING: 31;
  RINSING_PIPETTE: 40;
  DISPENSING: 50;
  WAITING_FOR_DISPENSE: 55;
  COLLECTING: 60;
  NAVIGATING_IK: 70;
  NAVIGATING_OUTER: 75;
}

export const Status: StatusMap;

export interface FluidTypeMap {
  FLUID_UNDEFINED: 0;
  FLUID_DRAIN: 1;
  FLUID_WATER: 2;
  FLUID_MILK: 3;
}

export const FluidType: FluidTypeMap;

export interface ContentTypeMap {
  CONTENT_TYPE_UNDEFINED: 0;
  CONTENT_TYPE_LONGFORM: 1;
  CONTENT_TYPE_SHORTFORM: 2;
  CONTENT_TYPE_CLEANING: 3;
  CONTENT_TYPE_DSLR: 4;
  CONTENT_TYPE_STILL: 5;
}

export const ContentType: ContentTypeMap;

export interface SocialPlatformMap {
  SOCIAL_PLATFORM_UNDEFINED: 0;
  SOCIAL_PLATFORM_YOUTUBE: 1;
  SOCIAL_PLATFORM_TIKTOK: 2;
  SOCIAL_PLATFORM_INSTAGRAM: 3;
  SOCIAL_PLATFORM_FACEBOOK: 4;
  SOCIAL_PLATFORM_TWITTER: 5;
  SOCIAL_PLATFORM_REDDIT: 6;
}

export const SocialPlatform: SocialPlatformMap;

