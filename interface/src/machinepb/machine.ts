/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "machine";

export enum Node {
  UNDEFINED = 0,
  HOME = 4,
  HOME_TOP = 8,
  /**
   * VIAL_1_ABOVE - Above and inside test tube positions
   * Note; INSIDE positions are valid for a range of z values, determined outside Navigation.
   */
  VIAL_1_ABOVE = 10,
  MIN_VIAL_ABOVE = 10,
  VIAL_1_INSIDE = 15,
  MIN_VIAL_INSIDE = 15,
  VIAL_2_ABOVE = 20,
  VIAL_2_INSIDE = 25,
  VIAL_3_ABOVE = 30,
  VIAL_3_INSIDE = 35,
  VIAL_4_ABOVE = 40,
  VIAL_4_INSIDE = 45,
  VIAL_5_ABOVE = 50,
  VIAL_5_INSIDE = 55,
  VIAL_6_ABOVE = 60,
  VIAL_6_INSIDE = 65,
  VIAL_7_ABOVE = 70,
  MAX_VIAL_ABOVE = 70,
  VIAL_7_INSIDE = 75,
  MAX_VIAL_INSIDE = 75,
  /** LOW_ENTRY_POINT - The node to enter the lower (vial) regions at */
  LOW_ENTRY_POINT = 30,
  /** RINSE_CONTAINER_ENTRY - High z but otherwise aligned for rinse container */
  RINSE_CONTAINER_ENTRY = 80,
  /** RINSE_CONTAINER_LOW - Low z and aligned for rinse container (in water) */
  RINSE_CONTAINER_LOW = 85,
  OUTER_HANDOVER = 90,
  INNER_HANDOVER = 110,
  INVERSE_KINEMATICS_POSITION = 150,
  IDLE_LOCATION = 80,
  UNRECOGNIZED = -1,
}

export function nodeFromJSON(object: any): Node {
  switch (object) {
    case 0:
    case "UNDEFINED":
      return Node.UNDEFINED;
    case 4:
    case "HOME":
      return Node.HOME;
    case 8:
    case "HOME_TOP":
      return Node.HOME_TOP;
    case 10:
    case "VIAL_1_ABOVE":
      return Node.VIAL_1_ABOVE;
    case 10:
    case "MIN_VIAL_ABOVE":
      return Node.MIN_VIAL_ABOVE;
    case 15:
    case "VIAL_1_INSIDE":
      return Node.VIAL_1_INSIDE;
    case 15:
    case "MIN_VIAL_INSIDE":
      return Node.MIN_VIAL_INSIDE;
    case 20:
    case "VIAL_2_ABOVE":
      return Node.VIAL_2_ABOVE;
    case 25:
    case "VIAL_2_INSIDE":
      return Node.VIAL_2_INSIDE;
    case 30:
    case "VIAL_3_ABOVE":
      return Node.VIAL_3_ABOVE;
    case 35:
    case "VIAL_3_INSIDE":
      return Node.VIAL_3_INSIDE;
    case 40:
    case "VIAL_4_ABOVE":
      return Node.VIAL_4_ABOVE;
    case 45:
    case "VIAL_4_INSIDE":
      return Node.VIAL_4_INSIDE;
    case 50:
    case "VIAL_5_ABOVE":
      return Node.VIAL_5_ABOVE;
    case 55:
    case "VIAL_5_INSIDE":
      return Node.VIAL_5_INSIDE;
    case 60:
    case "VIAL_6_ABOVE":
      return Node.VIAL_6_ABOVE;
    case 65:
    case "VIAL_6_INSIDE":
      return Node.VIAL_6_INSIDE;
    case 70:
    case "VIAL_7_ABOVE":
      return Node.VIAL_7_ABOVE;
    case 70:
    case "MAX_VIAL_ABOVE":
      return Node.MAX_VIAL_ABOVE;
    case 75:
    case "VIAL_7_INSIDE":
      return Node.VIAL_7_INSIDE;
    case 75:
    case "MAX_VIAL_INSIDE":
      return Node.MAX_VIAL_INSIDE;
    case 30:
    case "LOW_ENTRY_POINT":
      return Node.LOW_ENTRY_POINT;
    case 80:
    case "RINSE_CONTAINER_ENTRY":
      return Node.RINSE_CONTAINER_ENTRY;
    case 85:
    case "RINSE_CONTAINER_LOW":
      return Node.RINSE_CONTAINER_LOW;
    case 90:
    case "OUTER_HANDOVER":
      return Node.OUTER_HANDOVER;
    case 110:
    case "INNER_HANDOVER":
      return Node.INNER_HANDOVER;
    case 150:
    case "INVERSE_KINEMATICS_POSITION":
      return Node.INVERSE_KINEMATICS_POSITION;
    case 80:
    case "IDLE_LOCATION":
      return Node.IDLE_LOCATION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Node.UNRECOGNIZED;
  }
}

export function nodeToJSON(object: Node): string {
  switch (object) {
    case Node.UNDEFINED:
      return "UNDEFINED";
    case Node.HOME:
      return "HOME";
    case Node.HOME_TOP:
      return "HOME_TOP";
    case Node.VIAL_1_ABOVE:
      return "VIAL_1_ABOVE";
    case Node.MIN_VIAL_ABOVE:
      return "MIN_VIAL_ABOVE";
    case Node.VIAL_1_INSIDE:
      return "VIAL_1_INSIDE";
    case Node.MIN_VIAL_INSIDE:
      return "MIN_VIAL_INSIDE";
    case Node.VIAL_2_ABOVE:
      return "VIAL_2_ABOVE";
    case Node.VIAL_2_INSIDE:
      return "VIAL_2_INSIDE";
    case Node.VIAL_3_ABOVE:
      return "VIAL_3_ABOVE";
    case Node.VIAL_3_INSIDE:
      return "VIAL_3_INSIDE";
    case Node.VIAL_4_ABOVE:
      return "VIAL_4_ABOVE";
    case Node.VIAL_4_INSIDE:
      return "VIAL_4_INSIDE";
    case Node.VIAL_5_ABOVE:
      return "VIAL_5_ABOVE";
    case Node.VIAL_5_INSIDE:
      return "VIAL_5_INSIDE";
    case Node.VIAL_6_ABOVE:
      return "VIAL_6_ABOVE";
    case Node.VIAL_6_INSIDE:
      return "VIAL_6_INSIDE";
    case Node.VIAL_7_ABOVE:
      return "VIAL_7_ABOVE";
    case Node.MAX_VIAL_ABOVE:
      return "MAX_VIAL_ABOVE";
    case Node.VIAL_7_INSIDE:
      return "VIAL_7_INSIDE";
    case Node.MAX_VIAL_INSIDE:
      return "MAX_VIAL_INSIDE";
    case Node.LOW_ENTRY_POINT:
      return "LOW_ENTRY_POINT";
    case Node.RINSE_CONTAINER_ENTRY:
      return "RINSE_CONTAINER_ENTRY";
    case Node.RINSE_CONTAINER_LOW:
      return "RINSE_CONTAINER_LOW";
    case Node.OUTER_HANDOVER:
      return "OUTER_HANDOVER";
    case Node.INNER_HANDOVER:
      return "INNER_HANDOVER";
    case Node.INVERSE_KINEMATICS_POSITION:
      return "INVERSE_KINEMATICS_POSITION";
    case Node.IDLE_LOCATION:
      return "IDLE_LOCATION";
    case Node.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** used in requests */
export enum SolenoidValve {
  VALVE_UNDEFINED = 0,
  VALVE_DRAIN = 1,
  VALVE_WATER = 2,
  VALVE_MILK = 3,
  VALVE_AIR = 4,
  UNRECOGNIZED = -1,
}

export function solenoidValveFromJSON(object: any): SolenoidValve {
  switch (object) {
    case 0:
    case "VALVE_UNDEFINED":
      return SolenoidValve.VALVE_UNDEFINED;
    case 1:
    case "VALVE_DRAIN":
      return SolenoidValve.VALVE_DRAIN;
    case 2:
    case "VALVE_WATER":
      return SolenoidValve.VALVE_WATER;
    case 3:
    case "VALVE_MILK":
      return SolenoidValve.VALVE_MILK;
    case 4:
    case "VALVE_AIR":
      return SolenoidValve.VALVE_AIR;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SolenoidValve.UNRECOGNIZED;
  }
}

export function solenoidValveToJSON(object: SolenoidValve): string {
  switch (object) {
    case SolenoidValve.VALVE_UNDEFINED:
      return "VALVE_UNDEFINED";
    case SolenoidValve.VALVE_DRAIN:
      return "VALVE_DRAIN";
    case SolenoidValve.VALVE_WATER:
      return "VALVE_WATER";
    case SolenoidValve.VALVE_MILK:
      return "VALVE_MILK";
    case SolenoidValve.VALVE_AIR:
      return "VALVE_AIR";
    case SolenoidValve.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum Mode {
  UNDEFINED_MODE = 0,
  MANUAL = 1,
  AUTONOMOUS = 2,
  UNRECOGNIZED = -1,
}

export function modeFromJSON(object: any): Mode {
  switch (object) {
    case 0:
    case "UNDEFINED_MODE":
      return Mode.UNDEFINED_MODE;
    case 1:
    case "MANUAL":
      return Mode.MANUAL;
    case 2:
    case "AUTONOMOUS":
      return Mode.AUTONOMOUS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Mode.UNRECOGNIZED;
  }
}

export function modeToJSON(object: Mode): string {
  switch (object) {
    case Mode.UNDEFINED_MODE:
      return "UNDEFINED_MODE";
    case Mode.MANUAL:
      return "MANUAL";
    case Mode.AUTONOMOUS:
      return "AUTONOMOUS";
    case Mode.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum Status {
  UNDEFINED_STATUS = 0,
  ERROR = 1,
  E_STOP_ACTIVE = 5,
  SLEEPING = 6,
  SHUTTING_DOWN = 9,
  WAKING_UP = 10,
  CALIBRATING = 20,
  IDLE_STATIONARY = 30,
  IDLE_MOVING = 31,
  RINSING_PIPETTE = 40,
  DISPENSING = 50,
  WAITING_FOR_DISPENSE = 55,
  COLLECTING = 60,
  NAVIGATING_IK = 70,
  NAVIGATING_OUTER = 75,
  UNRECOGNIZED = -1,
}

export function statusFromJSON(object: any): Status {
  switch (object) {
    case 0:
    case "UNDEFINED_STATUS":
      return Status.UNDEFINED_STATUS;
    case 1:
    case "ERROR":
      return Status.ERROR;
    case 5:
    case "E_STOP_ACTIVE":
      return Status.E_STOP_ACTIVE;
    case 6:
    case "SLEEPING":
      return Status.SLEEPING;
    case 9:
    case "SHUTTING_DOWN":
      return Status.SHUTTING_DOWN;
    case 10:
    case "WAKING_UP":
      return Status.WAKING_UP;
    case 20:
    case "CALIBRATING":
      return Status.CALIBRATING;
    case 30:
    case "IDLE_STATIONARY":
      return Status.IDLE_STATIONARY;
    case 31:
    case "IDLE_MOVING":
      return Status.IDLE_MOVING;
    case 40:
    case "RINSING_PIPETTE":
      return Status.RINSING_PIPETTE;
    case 50:
    case "DISPENSING":
      return Status.DISPENSING;
    case 55:
    case "WAITING_FOR_DISPENSE":
      return Status.WAITING_FOR_DISPENSE;
    case 60:
    case "COLLECTING":
      return Status.COLLECTING;
    case 70:
    case "NAVIGATING_IK":
      return Status.NAVIGATING_IK;
    case 75:
    case "NAVIGATING_OUTER":
      return Status.NAVIGATING_OUTER;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Status.UNRECOGNIZED;
  }
}

export function statusToJSON(object: Status): string {
  switch (object) {
    case Status.UNDEFINED_STATUS:
      return "UNDEFINED_STATUS";
    case Status.ERROR:
      return "ERROR";
    case Status.E_STOP_ACTIVE:
      return "E_STOP_ACTIVE";
    case Status.SLEEPING:
      return "SLEEPING";
    case Status.SHUTTING_DOWN:
      return "SHUTTING_DOWN";
    case Status.WAKING_UP:
      return "WAKING_UP";
    case Status.CALIBRATING:
      return "CALIBRATING";
    case Status.IDLE_STATIONARY:
      return "IDLE_STATIONARY";
    case Status.IDLE_MOVING:
      return "IDLE_MOVING";
    case Status.RINSING_PIPETTE:
      return "RINSING_PIPETTE";
    case Status.DISPENSING:
      return "DISPENSING";
    case Status.WAITING_FOR_DISPENSE:
      return "WAITING_FOR_DISPENSE";
    case Status.COLLECTING:
      return "COLLECTING";
    case Status.NAVIGATING_IK:
      return "NAVIGATING_IK";
    case Status.NAVIGATING_OUTER:
      return "NAVIGATING_OUTER";
    case Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum RinseStatus {
  RINSE_UNDEFINED = 0,
  RINSE_COMPLETE = 1,
  RINSE_REQUESTED = 2,
  RINSE_EXPELLING = 3,
  UNRECOGNIZED = -1,
}

export function rinseStatusFromJSON(object: any): RinseStatus {
  switch (object) {
    case 0:
    case "RINSE_UNDEFINED":
      return RinseStatus.RINSE_UNDEFINED;
    case 1:
    case "RINSE_COMPLETE":
      return RinseStatus.RINSE_COMPLETE;
    case 2:
    case "RINSE_REQUESTED":
      return RinseStatus.RINSE_REQUESTED;
    case 3:
    case "RINSE_EXPELLING":
      return RinseStatus.RINSE_EXPELLING;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RinseStatus.UNRECOGNIZED;
  }
}

export function rinseStatusToJSON(object: RinseStatus): string {
  switch (object) {
    case RinseStatus.RINSE_UNDEFINED:
      return "RINSE_UNDEFINED";
    case RinseStatus.RINSE_COMPLETE:
      return "RINSE_COMPLETE";
    case RinseStatus.RINSE_REQUESTED:
      return "RINSE_REQUESTED";
    case RinseStatus.RINSE_EXPELLING:
      return "RINSE_EXPELLING";
    case RinseStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum FluidType {
  FLUID_UNDEFINED = 0,
  FLUID_DRAIN = 1,
  FLUID_WATER = 2,
  FLUID_MILK = 3,
  UNRECOGNIZED = -1,
}

export function fluidTypeFromJSON(object: any): FluidType {
  switch (object) {
    case 0:
    case "FLUID_UNDEFINED":
      return FluidType.FLUID_UNDEFINED;
    case 1:
    case "FLUID_DRAIN":
      return FluidType.FLUID_DRAIN;
    case 2:
    case "FLUID_WATER":
      return FluidType.FLUID_WATER;
    case 3:
    case "FLUID_MILK":
      return FluidType.FLUID_MILK;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FluidType.UNRECOGNIZED;
  }
}

export function fluidTypeToJSON(object: FluidType): string {
  switch (object) {
    case FluidType.FLUID_UNDEFINED:
      return "FLUID_UNDEFINED";
    case FluidType.FLUID_DRAIN:
      return "FLUID_DRAIN";
    case FluidType.FLUID_WATER:
      return "FLUID_WATER";
    case FluidType.FLUID_MILK:
      return "FLUID_MILK";
    case FluidType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ContentType {
  CONTENT_TYPE_UNDEFINED = 0,
  CONTENT_TYPE_LONGFORM = 1,
  CONTENT_TYPE_SHORTFORM = 2,
  CONTENT_TYPE_CLEANING = 3,
  CONTENT_TYPE_DSLR = 4,
  CONTENT_TYPE_STILL = 5,
  UNRECOGNIZED = -1,
}

export function contentTypeFromJSON(object: any): ContentType {
  switch (object) {
    case 0:
    case "CONTENT_TYPE_UNDEFINED":
      return ContentType.CONTENT_TYPE_UNDEFINED;
    case 1:
    case "CONTENT_TYPE_LONGFORM":
      return ContentType.CONTENT_TYPE_LONGFORM;
    case 2:
    case "CONTENT_TYPE_SHORTFORM":
      return ContentType.CONTENT_TYPE_SHORTFORM;
    case 3:
    case "CONTENT_TYPE_CLEANING":
      return ContentType.CONTENT_TYPE_CLEANING;
    case 4:
    case "CONTENT_TYPE_DSLR":
      return ContentType.CONTENT_TYPE_DSLR;
    case 5:
    case "CONTENT_TYPE_STILL":
      return ContentType.CONTENT_TYPE_STILL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ContentType.UNRECOGNIZED;
  }
}

export function contentTypeToJSON(object: ContentType): string {
  switch (object) {
    case ContentType.CONTENT_TYPE_UNDEFINED:
      return "CONTENT_TYPE_UNDEFINED";
    case ContentType.CONTENT_TYPE_LONGFORM:
      return "CONTENT_TYPE_LONGFORM";
    case ContentType.CONTENT_TYPE_SHORTFORM:
      return "CONTENT_TYPE_SHORTFORM";
    case ContentType.CONTENT_TYPE_CLEANING:
      return "CONTENT_TYPE_CLEANING";
    case ContentType.CONTENT_TYPE_DSLR:
      return "CONTENT_TYPE_DSLR";
    case ContentType.CONTENT_TYPE_STILL:
      return "CONTENT_TYPE_STILL";
    case ContentType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SocialPlatform {
  SOCIAL_PLATFORM_UNDEFINED = 0,
  SOCIAL_PLATFORM_YOUTUBE = 1,
  SOCIAL_PLATFORM_TIKTOK = 2,
  SOCIAL_PLATFORM_INSTAGRAM = 3,
  SOCIAL_PLATFORM_FACEBOOK = 4,
  SOCIAL_PLATFORM_TWITTER = 5,
  SOCIAL_PLATFORM_REDDIT = 6,
  UNRECOGNIZED = -1,
}

export function socialPlatformFromJSON(object: any): SocialPlatform {
  switch (object) {
    case 0:
    case "SOCIAL_PLATFORM_UNDEFINED":
      return SocialPlatform.SOCIAL_PLATFORM_UNDEFINED;
    case 1:
    case "SOCIAL_PLATFORM_YOUTUBE":
      return SocialPlatform.SOCIAL_PLATFORM_YOUTUBE;
    case 2:
    case "SOCIAL_PLATFORM_TIKTOK":
      return SocialPlatform.SOCIAL_PLATFORM_TIKTOK;
    case 3:
    case "SOCIAL_PLATFORM_INSTAGRAM":
      return SocialPlatform.SOCIAL_PLATFORM_INSTAGRAM;
    case 4:
    case "SOCIAL_PLATFORM_FACEBOOK":
      return SocialPlatform.SOCIAL_PLATFORM_FACEBOOK;
    case 5:
    case "SOCIAL_PLATFORM_TWITTER":
      return SocialPlatform.SOCIAL_PLATFORM_TWITTER;
    case 6:
    case "SOCIAL_PLATFORM_REDDIT":
      return SocialPlatform.SOCIAL_PLATFORM_REDDIT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SocialPlatform.UNRECOGNIZED;
  }
}

export function socialPlatformToJSON(object: SocialPlatform): string {
  switch (object) {
    case SocialPlatform.SOCIAL_PLATFORM_UNDEFINED:
      return "SOCIAL_PLATFORM_UNDEFINED";
    case SocialPlatform.SOCIAL_PLATFORM_YOUTUBE:
      return "SOCIAL_PLATFORM_YOUTUBE";
    case SocialPlatform.SOCIAL_PLATFORM_TIKTOK:
      return "SOCIAL_PLATFORM_TIKTOK";
    case SocialPlatform.SOCIAL_PLATFORM_INSTAGRAM:
      return "SOCIAL_PLATFORM_INSTAGRAM";
    case SocialPlatform.SOCIAL_PLATFORM_FACEBOOK:
      return "SOCIAL_PLATFORM_FACEBOOK";
    case SocialPlatform.SOCIAL_PLATFORM_TWITTER:
      return "SOCIAL_PLATFORM_TWITTER";
    case SocialPlatform.SOCIAL_PLATFORM_REDDIT:
      return "SOCIAL_PLATFORM_REDDIT";
    case SocialPlatform.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum EmailRecipient {
  EMAIL_RECIPIENT_UNDEFINED = 0,
  EMAIL_RECIPIENT_MAINTENANCE = 1,
  EMAIL_RECIPIENT_ROUTINE_OPERATIONS = 2,
  EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS = 3,
  UNRECOGNIZED = -1,
}

export function emailRecipientFromJSON(object: any): EmailRecipient {
  switch (object) {
    case 0:
    case "EMAIL_RECIPIENT_UNDEFINED":
      return EmailRecipient.EMAIL_RECIPIENT_UNDEFINED;
    case 1:
    case "EMAIL_RECIPIENT_MAINTENANCE":
      return EmailRecipient.EMAIL_RECIPIENT_MAINTENANCE;
    case 2:
    case "EMAIL_RECIPIENT_ROUTINE_OPERATIONS":
      return EmailRecipient.EMAIL_RECIPIENT_ROUTINE_OPERATIONS;
    case 3:
    case "EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS":
      return EmailRecipient.EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return EmailRecipient.UNRECOGNIZED;
  }
}

export function emailRecipientToJSON(object: EmailRecipient): string {
  switch (object) {
    case EmailRecipient.EMAIL_RECIPIENT_UNDEFINED:
      return "EMAIL_RECIPIENT_UNDEFINED";
    case EmailRecipient.EMAIL_RECIPIENT_MAINTENANCE:
      return "EMAIL_RECIPIENT_MAINTENANCE";
    case EmailRecipient.EMAIL_RECIPIENT_ROUTINE_OPERATIONS:
      return "EMAIL_RECIPIENT_ROUTINE_OPERATIONS";
    case EmailRecipient.EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS:
      return "EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS";
    case EmailRecipient.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface PipetteState {
  spent: boolean;
  vialHeld: number;
  volumeTargetUl: number;
  /** incremented every time a dispense is requested */
  dispenseRequestNumber: number;
}

export interface CollectionRequest {
  completed: boolean;
  requestNumber: number;
  vialNumber: number;
  volumeUl: number;
}

export interface MovementDetails {
  /** ik target from -1 to 1 */
  targetXUnit: number;
  /** ik target from -1 to 1 */
  targetYUnit: number;
  /** ik z target in mm */
  targetZIk: number;
  /** fk target in degrees */
  targetRingDeg: number;
  /** fk target in degrees */
  targetYawDeg: number;
}

export interface FluidRequest {
  fluidType: FluidType;
  volumeMl: number;
  complete: boolean;
  /**
   * if true, open drain while request is taking place
   * (e.g. for rinsing with water)
   */
  openDrain: boolean;
}

export interface FluidDetails {
  bowlFluidLevelMl: number;
}

export interface StateReport {
  /**
   * timestamp in microseconds since unix epoch, UTC. Added
   * by gateway since firmware doesn't know real time.
   */
  timestampUnixMicros: number;
  /** incremented on startup, currently 1 byte */
  startupCounter: number;
  mode: Mode;
  status: Status;
  /** Useful for synchronisation with footage */
  lightsOn: boolean;
  pipetteState: PipetteState | undefined;
  collectionRequest: CollectionRequest | undefined;
  movementDetails: MovementDetails | undefined;
  fluidRequest: FluidRequest | undefined;
  fluidDetails: FluidDetails | undefined;
  rinseStatus: RinseStatus;
  /** the following are populated by the backend, useful in post-processing */
  paused: boolean;
  timestampReadable: string;
  /** e.g. 1 for 0001.jpg */
  latestDslrFileNumber: number;
}

export interface StateReportList {
  StateReports: StateReport[];
}

export interface SessionStatus {
  id: number;
  paused: boolean;
  complete: boolean;
  production: boolean;
  productionId: number;
}

export interface StreamStatus {
  live: boolean;
}

export interface DispenseMetadataMap {
  /** [startupCounter]_[dispenseRequestNumber] */
  dispenseMetadata: { [key: string]: DispenseMetadata };
}

export interface DispenseMetadataMap_DispenseMetadataEntry {
  key: string;
  value: DispenseMetadata | undefined;
}

export interface DispenseMetadata {
  failedDispense: boolean;
  /** how many ms later than expected the dispense happened */
  dispenseDelayMs: number;
  /** if non-zero, override the vial profile's duration with this value. */
  minDurationOverrideMs: number;
  /** if non-zero, override the vial profile's speed with this value. */
  speedMultOverride: number;
}

/** statuses for all the content types for a specific session */
export interface ContentTypeStatuses {
  /** str(ContentType) -> ContentTypeStatus */
  contentStatuses: { [key: string]: ContentTypeStatus };
  /** splashtext for this session */
  splashtext: string;
  splashtextHue: number;
}

export interface ContentTypeStatuses_ContentStatusesEntry {
  key: string;
  value: ContentTypeStatus | undefined;
}

export interface ContentTypeStatus {
  rawTitle: string;
  rawDescription: string;
  caption: string;
  posts: Post[];
  musicFile: string;
  musicName: string;
}

export interface Post {
  platform: SocialPlatform;
  /** e.g. subreddit */
  subPlatform: string;
  title: string;
  description: string;
  uploaded: boolean;
  url: string;
  /** if true and relevant, crosspost rather than reuploading, e.g. for reddit */
  crosspost: boolean;
  /** seconds ts of when to publish. If 0, publish immediately, because 0 is in the past. */
  scheduledUnixTimetamp: number;
  /** if true, video will be posted unlisted, accessible by link only. Or not posted if the platform doesn't support it. */
  unlisted: boolean;
}

/** emails used for administration, not intended for audience distribution */
export interface Email {
  subject: string;
  body: string;
  recipient: EmailRecipient;
}

/**
 * This contains information about each vial/test tube.
 *
 * These should be maintained over time by the frontend interface and the backend
 * in response to dispenses.
 *
 * The current value is copied into session files when a session starts if it's in
 * the system.
 */
export interface VialProfile {
  /** incremental unique id for each vial in and out the system */
  id: number;
  /**
   * this should have a complete description of the mixture, including base
   * fluids and the percentage makeup of each. This may be augmented by
   * quantised makeup data in future.
   */
  description: string;
  /** the pipette slop, how much extra volume to move on the first dispense */
  slopUl: number;
  /** how much volume to dispense each time */
  dispenseVolumeUl: number;
  /** how long after dispense to slow down the footage in the videos */
  footageDelayMs: number;
  /** how long to keep the footage slowed down in the videos */
  footageMinDurationMs: number;
  /** what speed to give the footage in the videos */
  footageSpeedMult: number;
  /**
   * if true, footage of this profile will not be treated differently
   * to other footage (no slowdown etc.)
   */
  footageIgnore: boolean;
  /** Volume when this was first put in vial */
  initialVolumeUl: number;
  /**
   * Current volume. Note this will be just volume at start of session in
   * session files.
   */
  currentVolumeUl: number;
  /** friendly name for use in interfaces */
  name: string;
  vialFluid: VialProfile_VialFluid;
  /** colour to represent this in interfaces, of the form '#aa22ff' */
  colour: string;
  /** alternate names that can be used in voting */
  aliases: string[];
}

export enum VialProfile_VialFluid {
  VIAL_FLUID_UNDEFINED = 0,
  VIAL_FLUID_DYE_WATER_BASED = 1,
  VIAL_FLUID_EMULSIFIER = 2,
  VIAL_FLUID_AIR = 3,
  VIAL_FLUID_SOLVENT = 4,
  UNRECOGNIZED = -1,
}

export function vialProfile_VialFluidFromJSON(object: any): VialProfile_VialFluid {
  switch (object) {
    case 0:
    case "VIAL_FLUID_UNDEFINED":
      return VialProfile_VialFluid.VIAL_FLUID_UNDEFINED;
    case 1:
    case "VIAL_FLUID_DYE_WATER_BASED":
      return VialProfile_VialFluid.VIAL_FLUID_DYE_WATER_BASED;
    case 2:
    case "VIAL_FLUID_EMULSIFIER":
      return VialProfile_VialFluid.VIAL_FLUID_EMULSIFIER;
    case 3:
    case "VIAL_FLUID_AIR":
      return VialProfile_VialFluid.VIAL_FLUID_AIR;
    case 4:
    case "VIAL_FLUID_SOLVENT":
      return VialProfile_VialFluid.VIAL_FLUID_SOLVENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return VialProfile_VialFluid.UNRECOGNIZED;
  }
}

export function vialProfile_VialFluidToJSON(object: VialProfile_VialFluid): string {
  switch (object) {
    case VialProfile_VialFluid.VIAL_FLUID_UNDEFINED:
      return "VIAL_FLUID_UNDEFINED";
    case VialProfile_VialFluid.VIAL_FLUID_DYE_WATER_BASED:
      return "VIAL_FLUID_DYE_WATER_BASED";
    case VialProfile_VialFluid.VIAL_FLUID_EMULSIFIER:
      return "VIAL_FLUID_EMULSIFIER";
    case VialProfile_VialFluid.VIAL_FLUID_AIR:
      return "VIAL_FLUID_AIR";
    case VialProfile_VialFluid.VIAL_FLUID_SOLVENT:
      return "VIAL_FLUID_SOLVENT";
    case VialProfile_VialFluid.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * contains a map of the current vial positions to vial profile ids
 * vial position -> VialProfile id.
 */
export interface SystemVialConfiguration {
  vials: { [key: number]: number };
}

export interface SystemVialConfiguration_VialsEntry {
  key: number;
  value: number;
}

/** this is for all the VialProfiles, mapped by id. */
export interface VialProfileCollection {
  /** VialProfile ID -> VialProfile */
  profiles: { [key: number]: VialProfile };
}

export interface VialProfileCollection_ProfilesEntry {
  key: number;
  value: VialProfile | undefined;
}

/** contains a static snapshot of the VialProfiles for each system position */
export interface SystemVialConfigurationSnapshot {
  profiles: { [key: number]: VialProfile };
}

export interface SystemVialConfigurationSnapshot_ProfilesEntry {
  key: number;
  value: VialProfile | undefined;
}

function createBasePipetteState(): PipetteState {
  return { spent: false, vialHeld: 0, volumeTargetUl: 0, dispenseRequestNumber: 0 };
}

export const PipetteState = {
  encode(message: PipetteState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.spent === true) {
      writer.uint32(8).bool(message.spent);
    }
    if (message.vialHeld !== 0) {
      writer.uint32(16).uint32(message.vialHeld);
    }
    if (message.volumeTargetUl !== 0) {
      writer.uint32(29).float(message.volumeTargetUl);
    }
    if (message.dispenseRequestNumber !== 0) {
      writer.uint32(32).uint32(message.dispenseRequestNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PipetteState {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePipetteState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.spent = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.vialHeld = reader.uint32();
          continue;
        case 3:
          if (tag !== 29) {
            break;
          }

          message.volumeTargetUl = reader.float();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.dispenseRequestNumber = reader.uint32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PipetteState {
    return {
      spent: isSet(object.spent) ? globalThis.Boolean(object.spent) : false,
      vialHeld: isSet(object.vial_held) ? globalThis.Number(object.vial_held) : 0,
      volumeTargetUl: isSet(object.volume_target_ul) ? globalThis.Number(object.volume_target_ul) : 0,
      dispenseRequestNumber: isSet(object.dispense_request_number)
        ? globalThis.Number(object.dispense_request_number)
        : 0,
    };
  },

  toJSON(message: PipetteState): unknown {
    const obj: any = {};
    if (message.spent === true) {
      obj.spent = message.spent;
    }
    if (message.vialHeld !== 0) {
      obj.vial_held = Math.round(message.vialHeld);
    }
    if (message.volumeTargetUl !== 0) {
      obj.volume_target_ul = message.volumeTargetUl;
    }
    if (message.dispenseRequestNumber !== 0) {
      obj.dispense_request_number = Math.round(message.dispenseRequestNumber);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PipetteState>, I>>(base?: I): PipetteState {
    return PipetteState.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PipetteState>, I>>(object: I): PipetteState {
    const message = createBasePipetteState();
    message.spent = object.spent ?? false;
    message.vialHeld = object.vialHeld ?? 0;
    message.volumeTargetUl = object.volumeTargetUl ?? 0;
    message.dispenseRequestNumber = object.dispenseRequestNumber ?? 0;
    return message;
  },
};

function createBaseCollectionRequest(): CollectionRequest {
  return { completed: false, requestNumber: 0, vialNumber: 0, volumeUl: 0 };
}

export const CollectionRequest = {
  encode(message: CollectionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.completed === true) {
      writer.uint32(8).bool(message.completed);
    }
    if (message.requestNumber !== 0) {
      writer.uint32(16).uint64(message.requestNumber);
    }
    if (message.vialNumber !== 0) {
      writer.uint32(24).uint64(message.vialNumber);
    }
    if (message.volumeUl !== 0) {
      writer.uint32(37).float(message.volumeUl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.completed = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.requestNumber = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.vialNumber = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 37) {
            break;
          }

          message.volumeUl = reader.float();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CollectionRequest {
    return {
      completed: isSet(object.completed) ? globalThis.Boolean(object.completed) : false,
      requestNumber: isSet(object.request_number) ? globalThis.Number(object.request_number) : 0,
      vialNumber: isSet(object.vial_number) ? globalThis.Number(object.vial_number) : 0,
      volumeUl: isSet(object.volume_ul) ? globalThis.Number(object.volume_ul) : 0,
    };
  },

  toJSON(message: CollectionRequest): unknown {
    const obj: any = {};
    if (message.completed === true) {
      obj.completed = message.completed;
    }
    if (message.requestNumber !== 0) {
      obj.request_number = Math.round(message.requestNumber);
    }
    if (message.vialNumber !== 0) {
      obj.vial_number = Math.round(message.vialNumber);
    }
    if (message.volumeUl !== 0) {
      obj.volume_ul = message.volumeUl;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CollectionRequest>, I>>(base?: I): CollectionRequest {
    return CollectionRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CollectionRequest>, I>>(object: I): CollectionRequest {
    const message = createBaseCollectionRequest();
    message.completed = object.completed ?? false;
    message.requestNumber = object.requestNumber ?? 0;
    message.vialNumber = object.vialNumber ?? 0;
    message.volumeUl = object.volumeUl ?? 0;
    return message;
  },
};

function createBaseMovementDetails(): MovementDetails {
  return { targetXUnit: 0, targetYUnit: 0, targetZIk: 0, targetRingDeg: 0, targetYawDeg: 0 };
}

export const MovementDetails = {
  encode(message: MovementDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetXUnit !== 0) {
      writer.uint32(13).float(message.targetXUnit);
    }
    if (message.targetYUnit !== 0) {
      writer.uint32(21).float(message.targetYUnit);
    }
    if (message.targetZIk !== 0) {
      writer.uint32(45).float(message.targetZIk);
    }
    if (message.targetRingDeg !== 0) {
      writer.uint32(85).float(message.targetRingDeg);
    }
    if (message.targetYawDeg !== 0) {
      writer.uint32(93).float(message.targetYawDeg);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MovementDetails {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMovementDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 13) {
            break;
          }

          message.targetXUnit = reader.float();
          continue;
        case 2:
          if (tag !== 21) {
            break;
          }

          message.targetYUnit = reader.float();
          continue;
        case 5:
          if (tag !== 45) {
            break;
          }

          message.targetZIk = reader.float();
          continue;
        case 10:
          if (tag !== 85) {
            break;
          }

          message.targetRingDeg = reader.float();
          continue;
        case 11:
          if (tag !== 93) {
            break;
          }

          message.targetYawDeg = reader.float();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MovementDetails {
    return {
      targetXUnit: isSet(object.target_x_unit) ? globalThis.Number(object.target_x_unit) : 0,
      targetYUnit: isSet(object.target_y_unit) ? globalThis.Number(object.target_y_unit) : 0,
      targetZIk: isSet(object.target_z_ik) ? globalThis.Number(object.target_z_ik) : 0,
      targetRingDeg: isSet(object.target_ring_deg) ? globalThis.Number(object.target_ring_deg) : 0,
      targetYawDeg: isSet(object.target_yaw_deg) ? globalThis.Number(object.target_yaw_deg) : 0,
    };
  },

  toJSON(message: MovementDetails): unknown {
    const obj: any = {};
    if (message.targetXUnit !== 0) {
      obj.target_x_unit = message.targetXUnit;
    }
    if (message.targetYUnit !== 0) {
      obj.target_y_unit = message.targetYUnit;
    }
    if (message.targetZIk !== 0) {
      obj.target_z_ik = message.targetZIk;
    }
    if (message.targetRingDeg !== 0) {
      obj.target_ring_deg = message.targetRingDeg;
    }
    if (message.targetYawDeg !== 0) {
      obj.target_yaw_deg = message.targetYawDeg;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MovementDetails>, I>>(base?: I): MovementDetails {
    return MovementDetails.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MovementDetails>, I>>(object: I): MovementDetails {
    const message = createBaseMovementDetails();
    message.targetXUnit = object.targetXUnit ?? 0;
    message.targetYUnit = object.targetYUnit ?? 0;
    message.targetZIk = object.targetZIk ?? 0;
    message.targetRingDeg = object.targetRingDeg ?? 0;
    message.targetYawDeg = object.targetYawDeg ?? 0;
    return message;
  },
};

function createBaseFluidRequest(): FluidRequest {
  return { fluidType: 0, volumeMl: 0, complete: false, openDrain: false };
}

export const FluidRequest = {
  encode(message: FluidRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fluidType !== 0) {
      writer.uint32(8).int32(message.fluidType);
    }
    if (message.volumeMl !== 0) {
      writer.uint32(21).float(message.volumeMl);
    }
    if (message.complete === true) {
      writer.uint32(24).bool(message.complete);
    }
    if (message.openDrain === true) {
      writer.uint32(32).bool(message.openDrain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FluidRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFluidRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.fluidType = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 21) {
            break;
          }

          message.volumeMl = reader.float();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.complete = reader.bool();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.openDrain = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FluidRequest {
    return {
      fluidType: isSet(object.fluidType) ? fluidTypeFromJSON(object.fluidType) : 0,
      volumeMl: isSet(object.volume_ml) ? globalThis.Number(object.volume_ml) : 0,
      complete: isSet(object.complete) ? globalThis.Boolean(object.complete) : false,
      openDrain: isSet(object.open_drain) ? globalThis.Boolean(object.open_drain) : false,
    };
  },

  toJSON(message: FluidRequest): unknown {
    const obj: any = {};
    if (message.fluidType !== 0) {
      obj.fluidType = fluidTypeToJSON(message.fluidType);
    }
    if (message.volumeMl !== 0) {
      obj.volume_ml = message.volumeMl;
    }
    if (message.complete === true) {
      obj.complete = message.complete;
    }
    if (message.openDrain === true) {
      obj.open_drain = message.openDrain;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FluidRequest>, I>>(base?: I): FluidRequest {
    return FluidRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FluidRequest>, I>>(object: I): FluidRequest {
    const message = createBaseFluidRequest();
    message.fluidType = object.fluidType ?? 0;
    message.volumeMl = object.volumeMl ?? 0;
    message.complete = object.complete ?? false;
    message.openDrain = object.openDrain ?? false;
    return message;
  },
};

function createBaseFluidDetails(): FluidDetails {
  return { bowlFluidLevelMl: 0 };
}

export const FluidDetails = {
  encode(message: FluidDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bowlFluidLevelMl !== 0) {
      writer.uint32(13).float(message.bowlFluidLevelMl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FluidDetails {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFluidDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 13) {
            break;
          }

          message.bowlFluidLevelMl = reader.float();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FluidDetails {
    return { bowlFluidLevelMl: isSet(object.bowl_fluid_level_ml) ? globalThis.Number(object.bowl_fluid_level_ml) : 0 };
  },

  toJSON(message: FluidDetails): unknown {
    const obj: any = {};
    if (message.bowlFluidLevelMl !== 0) {
      obj.bowl_fluid_level_ml = message.bowlFluidLevelMl;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FluidDetails>, I>>(base?: I): FluidDetails {
    return FluidDetails.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FluidDetails>, I>>(object: I): FluidDetails {
    const message = createBaseFluidDetails();
    message.bowlFluidLevelMl = object.bowlFluidLevelMl ?? 0;
    return message;
  },
};

function createBaseStateReport(): StateReport {
  return {
    timestampUnixMicros: 0,
    startupCounter: 0,
    mode: 0,
    status: 0,
    lightsOn: false,
    pipetteState: undefined,
    collectionRequest: undefined,
    movementDetails: undefined,
    fluidRequest: undefined,
    fluidDetails: undefined,
    rinseStatus: 0,
    paused: false,
    timestampReadable: "",
    latestDslrFileNumber: 0,
  };
}

export const StateReport = {
  encode(message: StateReport, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timestampUnixMicros !== 0) {
      writer.uint32(16).uint64(message.timestampUnixMicros);
    }
    if (message.startupCounter !== 0) {
      writer.uint32(24).uint64(message.startupCounter);
    }
    if (message.mode !== 0) {
      writer.uint32(32).int32(message.mode);
    }
    if (message.status !== 0) {
      writer.uint32(40).int32(message.status);
    }
    if (message.lightsOn === true) {
      writer.uint32(48).bool(message.lightsOn);
    }
    if (message.pipetteState !== undefined) {
      PipetteState.encode(message.pipetteState, writer.uint32(82).fork()).ldelim();
    }
    if (message.collectionRequest !== undefined) {
      CollectionRequest.encode(message.collectionRequest, writer.uint32(90).fork()).ldelim();
    }
    if (message.movementDetails !== undefined) {
      MovementDetails.encode(message.movementDetails, writer.uint32(98).fork()).ldelim();
    }
    if (message.fluidRequest !== undefined) {
      FluidRequest.encode(message.fluidRequest, writer.uint32(106).fork()).ldelim();
    }
    if (message.fluidDetails !== undefined) {
      FluidDetails.encode(message.fluidDetails, writer.uint32(114).fork()).ldelim();
    }
    if (message.rinseStatus !== 0) {
      writer.uint32(120).int32(message.rinseStatus);
    }
    if (message.paused === true) {
      writer.uint32(400).bool(message.paused);
    }
    if (message.timestampReadable !== "") {
      writer.uint32(410).string(message.timestampReadable);
    }
    if (message.latestDslrFileNumber !== 0) {
      writer.uint32(416).uint64(message.latestDslrFileNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateReport {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateReport();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag !== 16) {
            break;
          }

          message.timestampUnixMicros = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.startupCounter = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.mode = reader.int32() as any;
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.status = reader.int32() as any;
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.lightsOn = reader.bool();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.pipetteState = PipetteState.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.collectionRequest = CollectionRequest.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.movementDetails = MovementDetails.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.fluidRequest = FluidRequest.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.fluidDetails = FluidDetails.decode(reader, reader.uint32());
          continue;
        case 15:
          if (tag !== 120) {
            break;
          }

          message.rinseStatus = reader.int32() as any;
          continue;
        case 50:
          if (tag !== 400) {
            break;
          }

          message.paused = reader.bool();
          continue;
        case 51:
          if (tag !== 410) {
            break;
          }

          message.timestampReadable = reader.string();
          continue;
        case 52:
          if (tag !== 416) {
            break;
          }

          message.latestDslrFileNumber = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StateReport {
    return {
      timestampUnixMicros: isSet(object.timestamp_unix_micros) ? globalThis.Number(object.timestamp_unix_micros) : 0,
      startupCounter: isSet(object.startup_counter) ? globalThis.Number(object.startup_counter) : 0,
      mode: isSet(object.mode) ? modeFromJSON(object.mode) : 0,
      status: isSet(object.status) ? statusFromJSON(object.status) : 0,
      lightsOn: isSet(object.lights_on) ? globalThis.Boolean(object.lights_on) : false,
      pipetteState: isSet(object.pipette_state) ? PipetteState.fromJSON(object.pipette_state) : undefined,
      collectionRequest: isSet(object.collection_request)
        ? CollectionRequest.fromJSON(object.collection_request)
        : undefined,
      movementDetails: isSet(object.movement_details) ? MovementDetails.fromJSON(object.movement_details) : undefined,
      fluidRequest: isSet(object.fluid_request) ? FluidRequest.fromJSON(object.fluid_request) : undefined,
      fluidDetails: isSet(object.fluid_details) ? FluidDetails.fromJSON(object.fluid_details) : undefined,
      rinseStatus: isSet(object.rinse_status) ? rinseStatusFromJSON(object.rinse_status) : 0,
      paused: isSet(object.paused) ? globalThis.Boolean(object.paused) : false,
      timestampReadable: isSet(object.timestamp_readable) ? globalThis.String(object.timestamp_readable) : "",
      latestDslrFileNumber: isSet(object.latest_dslr_file_number)
        ? globalThis.Number(object.latest_dslr_file_number)
        : 0,
    };
  },

  toJSON(message: StateReport): unknown {
    const obj: any = {};
    if (message.timestampUnixMicros !== 0) {
      obj.timestamp_unix_micros = Math.round(message.timestampUnixMicros);
    }
    if (message.startupCounter !== 0) {
      obj.startup_counter = Math.round(message.startupCounter);
    }
    if (message.mode !== 0) {
      obj.mode = modeToJSON(message.mode);
    }
    if (message.status !== 0) {
      obj.status = statusToJSON(message.status);
    }
    if (message.lightsOn === true) {
      obj.lights_on = message.lightsOn;
    }
    if (message.pipetteState !== undefined) {
      obj.pipette_state = PipetteState.toJSON(message.pipetteState);
    }
    if (message.collectionRequest !== undefined) {
      obj.collection_request = CollectionRequest.toJSON(message.collectionRequest);
    }
    if (message.movementDetails !== undefined) {
      obj.movement_details = MovementDetails.toJSON(message.movementDetails);
    }
    if (message.fluidRequest !== undefined) {
      obj.fluid_request = FluidRequest.toJSON(message.fluidRequest);
    }
    if (message.fluidDetails !== undefined) {
      obj.fluid_details = FluidDetails.toJSON(message.fluidDetails);
    }
    if (message.rinseStatus !== 0) {
      obj.rinse_status = rinseStatusToJSON(message.rinseStatus);
    }
    if (message.paused === true) {
      obj.paused = message.paused;
    }
    if (message.timestampReadable !== "") {
      obj.timestamp_readable = message.timestampReadable;
    }
    if (message.latestDslrFileNumber !== 0) {
      obj.latest_dslr_file_number = Math.round(message.latestDslrFileNumber);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StateReport>, I>>(base?: I): StateReport {
    return StateReport.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StateReport>, I>>(object: I): StateReport {
    const message = createBaseStateReport();
    message.timestampUnixMicros = object.timestampUnixMicros ?? 0;
    message.startupCounter = object.startupCounter ?? 0;
    message.mode = object.mode ?? 0;
    message.status = object.status ?? 0;
    message.lightsOn = object.lightsOn ?? false;
    message.pipetteState = (object.pipetteState !== undefined && object.pipetteState !== null)
      ? PipetteState.fromPartial(object.pipetteState)
      : undefined;
    message.collectionRequest = (object.collectionRequest !== undefined && object.collectionRequest !== null)
      ? CollectionRequest.fromPartial(object.collectionRequest)
      : undefined;
    message.movementDetails = (object.movementDetails !== undefined && object.movementDetails !== null)
      ? MovementDetails.fromPartial(object.movementDetails)
      : undefined;
    message.fluidRequest = (object.fluidRequest !== undefined && object.fluidRequest !== null)
      ? FluidRequest.fromPartial(object.fluidRequest)
      : undefined;
    message.fluidDetails = (object.fluidDetails !== undefined && object.fluidDetails !== null)
      ? FluidDetails.fromPartial(object.fluidDetails)
      : undefined;
    message.rinseStatus = object.rinseStatus ?? 0;
    message.paused = object.paused ?? false;
    message.timestampReadable = object.timestampReadable ?? "";
    message.latestDslrFileNumber = object.latestDslrFileNumber ?? 0;
    return message;
  },
};

function createBaseStateReportList(): StateReportList {
  return { StateReports: [] };
}

export const StateReportList = {
  encode(message: StateReportList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.StateReports) {
      StateReport.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateReportList {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateReportList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.StateReports.push(StateReport.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StateReportList {
    return {
      StateReports: globalThis.Array.isArray(object?.StateReports)
        ? object.StateReports.map((e: any) => StateReport.fromJSON(e))
        : [],
    };
  },

  toJSON(message: StateReportList): unknown {
    const obj: any = {};
    if (message.StateReports?.length) {
      obj.StateReports = message.StateReports.map((e) => StateReport.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StateReportList>, I>>(base?: I): StateReportList {
    return StateReportList.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StateReportList>, I>>(object: I): StateReportList {
    const message = createBaseStateReportList();
    message.StateReports = object.StateReports?.map((e) => StateReport.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSessionStatus(): SessionStatus {
  return { id: 0, paused: false, complete: false, production: false, productionId: 0 };
}

export const SessionStatus = {
  encode(message: SessionStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.paused === true) {
      writer.uint32(16).bool(message.paused);
    }
    if (message.complete === true) {
      writer.uint32(24).bool(message.complete);
    }
    if (message.production === true) {
      writer.uint32(32).bool(message.production);
    }
    if (message.productionId !== 0) {
      writer.uint32(40).uint64(message.productionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SessionStatus {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSessionStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.paused = reader.bool();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.complete = reader.bool();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.production = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.productionId = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SessionStatus {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      paused: isSet(object.paused) ? globalThis.Boolean(object.paused) : false,
      complete: isSet(object.complete) ? globalThis.Boolean(object.complete) : false,
      production: isSet(object.production) ? globalThis.Boolean(object.production) : false,
      productionId: isSet(object.production_id) ? globalThis.Number(object.production_id) : 0,
    };
  },

  toJSON(message: SessionStatus): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.paused === true) {
      obj.paused = message.paused;
    }
    if (message.complete === true) {
      obj.complete = message.complete;
    }
    if (message.production === true) {
      obj.production = message.production;
    }
    if (message.productionId !== 0) {
      obj.production_id = Math.round(message.productionId);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SessionStatus>, I>>(base?: I): SessionStatus {
    return SessionStatus.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SessionStatus>, I>>(object: I): SessionStatus {
    const message = createBaseSessionStatus();
    message.id = object.id ?? 0;
    message.paused = object.paused ?? false;
    message.complete = object.complete ?? false;
    message.production = object.production ?? false;
    message.productionId = object.productionId ?? 0;
    return message;
  },
};

function createBaseStreamStatus(): StreamStatus {
  return { live: false };
}

export const StreamStatus = {
  encode(message: StreamStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.live === true) {
      writer.uint32(8).bool(message.live);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamStatus {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.live = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StreamStatus {
    return { live: isSet(object.live) ? globalThis.Boolean(object.live) : false };
  },

  toJSON(message: StreamStatus): unknown {
    const obj: any = {};
    if (message.live === true) {
      obj.live = message.live;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StreamStatus>, I>>(base?: I): StreamStatus {
    return StreamStatus.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StreamStatus>, I>>(object: I): StreamStatus {
    const message = createBaseStreamStatus();
    message.live = object.live ?? false;
    return message;
  },
};

function createBaseDispenseMetadataMap(): DispenseMetadataMap {
  return { dispenseMetadata: {} };
}

export const DispenseMetadataMap = {
  encode(message: DispenseMetadataMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.dispenseMetadata).forEach(([key, value]) => {
      DispenseMetadataMap_DispenseMetadataEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DispenseMetadataMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDispenseMetadataMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = DispenseMetadataMap_DispenseMetadataEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.dispenseMetadata[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DispenseMetadataMap {
    return {
      dispenseMetadata: isObject(object.dispense_metadata)
        ? Object.entries(object.dispense_metadata).reduce<{ [key: string]: DispenseMetadata }>((acc, [key, value]) => {
          acc[key] = DispenseMetadata.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: DispenseMetadataMap): unknown {
    const obj: any = {};
    if (message.dispenseMetadata) {
      const entries = Object.entries(message.dispenseMetadata);
      if (entries.length > 0) {
        obj.dispense_metadata = {};
        entries.forEach(([k, v]) => {
          obj.dispense_metadata[k] = DispenseMetadata.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DispenseMetadataMap>, I>>(base?: I): DispenseMetadataMap {
    return DispenseMetadataMap.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DispenseMetadataMap>, I>>(object: I): DispenseMetadataMap {
    const message = createBaseDispenseMetadataMap();
    message.dispenseMetadata = Object.entries(object.dispenseMetadata ?? {}).reduce<
      { [key: string]: DispenseMetadata }
    >((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = DispenseMetadata.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseDispenseMetadataMap_DispenseMetadataEntry(): DispenseMetadataMap_DispenseMetadataEntry {
  return { key: "", value: undefined };
}

export const DispenseMetadataMap_DispenseMetadataEntry = {
  encode(message: DispenseMetadataMap_DispenseMetadataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DispenseMetadata.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DispenseMetadataMap_DispenseMetadataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDispenseMetadataMap_DispenseMetadataEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = DispenseMetadata.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DispenseMetadataMap_DispenseMetadataEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DispenseMetadata.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: DispenseMetadataMap_DispenseMetadataEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DispenseMetadata.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DispenseMetadataMap_DispenseMetadataEntry>, I>>(
    base?: I,
  ): DispenseMetadataMap_DispenseMetadataEntry {
    return DispenseMetadataMap_DispenseMetadataEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DispenseMetadataMap_DispenseMetadataEntry>, I>>(
    object: I,
  ): DispenseMetadataMap_DispenseMetadataEntry {
    const message = createBaseDispenseMetadataMap_DispenseMetadataEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DispenseMetadata.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseDispenseMetadata(): DispenseMetadata {
  return { failedDispense: false, dispenseDelayMs: 0, minDurationOverrideMs: 0, speedMultOverride: 0 };
}

export const DispenseMetadata = {
  encode(message: DispenseMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.failedDispense === true) {
      writer.uint32(8).bool(message.failedDispense);
    }
    if (message.dispenseDelayMs !== 0) {
      writer.uint32(16).uint64(message.dispenseDelayMs);
    }
    if (message.minDurationOverrideMs !== 0) {
      writer.uint32(24).uint64(message.minDurationOverrideMs);
    }
    if (message.speedMultOverride !== 0) {
      writer.uint32(32).uint64(message.speedMultOverride);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DispenseMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDispenseMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.failedDispense = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.dispenseDelayMs = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.minDurationOverrideMs = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.speedMultOverride = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DispenseMetadata {
    return {
      failedDispense: isSet(object.failed_dispense) ? globalThis.Boolean(object.failed_dispense) : false,
      dispenseDelayMs: isSet(object.dispense_delay_ms) ? globalThis.Number(object.dispense_delay_ms) : 0,
      minDurationOverrideMs: isSet(object.min_duration_override_ms)
        ? globalThis.Number(object.min_duration_override_ms)
        : 0,
      speedMultOverride: isSet(object.speed_mult_override) ? globalThis.Number(object.speed_mult_override) : 0,
    };
  },

  toJSON(message: DispenseMetadata): unknown {
    const obj: any = {};
    if (message.failedDispense === true) {
      obj.failed_dispense = message.failedDispense;
    }
    if (message.dispenseDelayMs !== 0) {
      obj.dispense_delay_ms = Math.round(message.dispenseDelayMs);
    }
    if (message.minDurationOverrideMs !== 0) {
      obj.min_duration_override_ms = Math.round(message.minDurationOverrideMs);
    }
    if (message.speedMultOverride !== 0) {
      obj.speed_mult_override = Math.round(message.speedMultOverride);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DispenseMetadata>, I>>(base?: I): DispenseMetadata {
    return DispenseMetadata.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DispenseMetadata>, I>>(object: I): DispenseMetadata {
    const message = createBaseDispenseMetadata();
    message.failedDispense = object.failedDispense ?? false;
    message.dispenseDelayMs = object.dispenseDelayMs ?? 0;
    message.minDurationOverrideMs = object.minDurationOverrideMs ?? 0;
    message.speedMultOverride = object.speedMultOverride ?? 0;
    return message;
  },
};

function createBaseContentTypeStatuses(): ContentTypeStatuses {
  return { contentStatuses: {}, splashtext: "", splashtextHue: 0 };
}

export const ContentTypeStatuses = {
  encode(message: ContentTypeStatuses, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.contentStatuses).forEach(([key, value]) => {
      ContentTypeStatuses_ContentStatusesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    if (message.splashtext !== "") {
      writer.uint32(18).string(message.splashtext);
    }
    if (message.splashtextHue !== 0) {
      writer.uint32(24).uint64(message.splashtextHue);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContentTypeStatuses {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContentTypeStatuses();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = ContentTypeStatuses_ContentStatusesEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.contentStatuses[entry1.key] = entry1.value;
          }
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.splashtext = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.splashtextHue = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ContentTypeStatuses {
    return {
      contentStatuses: isObject(object.content_statuses)
        ? Object.entries(object.content_statuses).reduce<{ [key: string]: ContentTypeStatus }>((acc, [key, value]) => {
          acc[key] = ContentTypeStatus.fromJSON(value);
          return acc;
        }, {})
        : {},
      splashtext: isSet(object.splashtext) ? globalThis.String(object.splashtext) : "",
      splashtextHue: isSet(object.splashtext_hue) ? globalThis.Number(object.splashtext_hue) : 0,
    };
  },

  toJSON(message: ContentTypeStatuses): unknown {
    const obj: any = {};
    if (message.contentStatuses) {
      const entries = Object.entries(message.contentStatuses);
      if (entries.length > 0) {
        obj.content_statuses = {};
        entries.forEach(([k, v]) => {
          obj.content_statuses[k] = ContentTypeStatus.toJSON(v);
        });
      }
    }
    if (message.splashtext !== "") {
      obj.splashtext = message.splashtext;
    }
    if (message.splashtextHue !== 0) {
      obj.splashtext_hue = Math.round(message.splashtextHue);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ContentTypeStatuses>, I>>(base?: I): ContentTypeStatuses {
    return ContentTypeStatuses.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ContentTypeStatuses>, I>>(object: I): ContentTypeStatuses {
    const message = createBaseContentTypeStatuses();
    message.contentStatuses = Object.entries(object.contentStatuses ?? {}).reduce<{ [key: string]: ContentTypeStatus }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = ContentTypeStatus.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.splashtext = object.splashtext ?? "";
    message.splashtextHue = object.splashtextHue ?? 0;
    return message;
  },
};

function createBaseContentTypeStatuses_ContentStatusesEntry(): ContentTypeStatuses_ContentStatusesEntry {
  return { key: "", value: undefined };
}

export const ContentTypeStatuses_ContentStatusesEntry = {
  encode(message: ContentTypeStatuses_ContentStatusesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      ContentTypeStatus.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContentTypeStatuses_ContentStatusesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContentTypeStatuses_ContentStatusesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = ContentTypeStatus.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ContentTypeStatuses_ContentStatusesEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? ContentTypeStatus.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ContentTypeStatuses_ContentStatusesEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = ContentTypeStatus.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ContentTypeStatuses_ContentStatusesEntry>, I>>(
    base?: I,
  ): ContentTypeStatuses_ContentStatusesEntry {
    return ContentTypeStatuses_ContentStatusesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ContentTypeStatuses_ContentStatusesEntry>, I>>(
    object: I,
  ): ContentTypeStatuses_ContentStatusesEntry {
    const message = createBaseContentTypeStatuses_ContentStatusesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? ContentTypeStatus.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseContentTypeStatus(): ContentTypeStatus {
  return { rawTitle: "", rawDescription: "", caption: "", posts: [], musicFile: "", musicName: "" };
}

export const ContentTypeStatus = {
  encode(message: ContentTypeStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rawTitle !== "") {
      writer.uint32(10).string(message.rawTitle);
    }
    if (message.rawDescription !== "") {
      writer.uint32(18).string(message.rawDescription);
    }
    if (message.caption !== "") {
      writer.uint32(26).string(message.caption);
    }
    for (const v of message.posts) {
      Post.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.musicFile !== "") {
      writer.uint32(58).string(message.musicFile);
    }
    if (message.musicName !== "") {
      writer.uint32(66).string(message.musicName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContentTypeStatus {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContentTypeStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.rawTitle = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.rawDescription = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.caption = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.posts.push(Post.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.musicFile = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.musicName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ContentTypeStatus {
    return {
      rawTitle: isSet(object.raw_title) ? globalThis.String(object.raw_title) : "",
      rawDescription: isSet(object.raw_description) ? globalThis.String(object.raw_description) : "",
      caption: isSet(object.caption) ? globalThis.String(object.caption) : "",
      posts: globalThis.Array.isArray(object?.posts) ? object.posts.map((e: any) => Post.fromJSON(e)) : [],
      musicFile: isSet(object.music_file) ? globalThis.String(object.music_file) : "",
      musicName: isSet(object.music_name) ? globalThis.String(object.music_name) : "",
    };
  },

  toJSON(message: ContentTypeStatus): unknown {
    const obj: any = {};
    if (message.rawTitle !== "") {
      obj.raw_title = message.rawTitle;
    }
    if (message.rawDescription !== "") {
      obj.raw_description = message.rawDescription;
    }
    if (message.caption !== "") {
      obj.caption = message.caption;
    }
    if (message.posts?.length) {
      obj.posts = message.posts.map((e) => Post.toJSON(e));
    }
    if (message.musicFile !== "") {
      obj.music_file = message.musicFile;
    }
    if (message.musicName !== "") {
      obj.music_name = message.musicName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ContentTypeStatus>, I>>(base?: I): ContentTypeStatus {
    return ContentTypeStatus.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ContentTypeStatus>, I>>(object: I): ContentTypeStatus {
    const message = createBaseContentTypeStatus();
    message.rawTitle = object.rawTitle ?? "";
    message.rawDescription = object.rawDescription ?? "";
    message.caption = object.caption ?? "";
    message.posts = object.posts?.map((e) => Post.fromPartial(e)) || [];
    message.musicFile = object.musicFile ?? "";
    message.musicName = object.musicName ?? "";
    return message;
  },
};

function createBasePost(): Post {
  return {
    platform: 0,
    subPlatform: "",
    title: "",
    description: "",
    uploaded: false,
    url: "",
    crosspost: false,
    scheduledUnixTimetamp: 0,
    unlisted: false,
  };
}

export const Post = {
  encode(message: Post, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.platform !== 0) {
      writer.uint32(8).int32(message.platform);
    }
    if (message.subPlatform !== "") {
      writer.uint32(18).string(message.subPlatform);
    }
    if (message.title !== "") {
      writer.uint32(26).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    if (message.uploaded === true) {
      writer.uint32(40).bool(message.uploaded);
    }
    if (message.url !== "") {
      writer.uint32(50).string(message.url);
    }
    if (message.crosspost === true) {
      writer.uint32(56).bool(message.crosspost);
    }
    if (message.scheduledUnixTimetamp !== 0) {
      writer.uint32(64).uint64(message.scheduledUnixTimetamp);
    }
    if (message.unlisted === true) {
      writer.uint32(72).bool(message.unlisted);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Post {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePost();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.platform = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.subPlatform = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.title = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.description = reader.string();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.uploaded = reader.bool();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.url = reader.string();
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.crosspost = reader.bool();
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.scheduledUnixTimetamp = longToNumber(reader.uint64() as Long);
          continue;
        case 9:
          if (tag !== 72) {
            break;
          }

          message.unlisted = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Post {
    return {
      platform: isSet(object.platform) ? socialPlatformFromJSON(object.platform) : 0,
      subPlatform: isSet(object.sub_platform) ? globalThis.String(object.sub_platform) : "",
      title: isSet(object.title) ? globalThis.String(object.title) : "",
      description: isSet(object.description) ? globalThis.String(object.description) : "",
      uploaded: isSet(object.uploaded) ? globalThis.Boolean(object.uploaded) : false,
      url: isSet(object.url) ? globalThis.String(object.url) : "",
      crosspost: isSet(object.crosspost) ? globalThis.Boolean(object.crosspost) : false,
      scheduledUnixTimetamp: isSet(object.scheduled_unix_timetamp)
        ? globalThis.Number(object.scheduled_unix_timetamp)
        : 0,
      unlisted: isSet(object.unlisted) ? globalThis.Boolean(object.unlisted) : false,
    };
  },

  toJSON(message: Post): unknown {
    const obj: any = {};
    if (message.platform !== 0) {
      obj.platform = socialPlatformToJSON(message.platform);
    }
    if (message.subPlatform !== "") {
      obj.sub_platform = message.subPlatform;
    }
    if (message.title !== "") {
      obj.title = message.title;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.uploaded === true) {
      obj.uploaded = message.uploaded;
    }
    if (message.url !== "") {
      obj.url = message.url;
    }
    if (message.crosspost === true) {
      obj.crosspost = message.crosspost;
    }
    if (message.scheduledUnixTimetamp !== 0) {
      obj.scheduled_unix_timetamp = Math.round(message.scheduledUnixTimetamp);
    }
    if (message.unlisted === true) {
      obj.unlisted = message.unlisted;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Post>, I>>(base?: I): Post {
    return Post.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Post>, I>>(object: I): Post {
    const message = createBasePost();
    message.platform = object.platform ?? 0;
    message.subPlatform = object.subPlatform ?? "";
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.uploaded = object.uploaded ?? false;
    message.url = object.url ?? "";
    message.crosspost = object.crosspost ?? false;
    message.scheduledUnixTimetamp = object.scheduledUnixTimetamp ?? 0;
    message.unlisted = object.unlisted ?? false;
    return message;
  },
};

function createBaseEmail(): Email {
  return { subject: "", body: "", recipient: 0 };
}

export const Email = {
  encode(message: Email, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subject !== "") {
      writer.uint32(10).string(message.subject);
    }
    if (message.body !== "") {
      writer.uint32(18).string(message.body);
    }
    if (message.recipient !== 0) {
      writer.uint32(24).int32(message.recipient);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Email {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.subject = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.body = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.recipient = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Email {
    return {
      subject: isSet(object.subject) ? globalThis.String(object.subject) : "",
      body: isSet(object.body) ? globalThis.String(object.body) : "",
      recipient: isSet(object.recipient) ? emailRecipientFromJSON(object.recipient) : 0,
    };
  },

  toJSON(message: Email): unknown {
    const obj: any = {};
    if (message.subject !== "") {
      obj.subject = message.subject;
    }
    if (message.body !== "") {
      obj.body = message.body;
    }
    if (message.recipient !== 0) {
      obj.recipient = emailRecipientToJSON(message.recipient);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Email>, I>>(base?: I): Email {
    return Email.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Email>, I>>(object: I): Email {
    const message = createBaseEmail();
    message.subject = object.subject ?? "";
    message.body = object.body ?? "";
    message.recipient = object.recipient ?? 0;
    return message;
  },
};

function createBaseVialProfile(): VialProfile {
  return {
    id: 0,
    description: "",
    slopUl: 0,
    dispenseVolumeUl: 0,
    footageDelayMs: 0,
    footageMinDurationMs: 0,
    footageSpeedMult: 0,
    footageIgnore: false,
    initialVolumeUl: 0,
    currentVolumeUl: 0,
    name: "",
    vialFluid: 0,
    colour: "",
    aliases: [],
  };
}

export const VialProfile = {
  encode(message: VialProfile, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.slopUl !== 0) {
      writer.uint32(29).float(message.slopUl);
    }
    if (message.dispenseVolumeUl !== 0) {
      writer.uint32(37).float(message.dispenseVolumeUl);
    }
    if (message.footageDelayMs !== 0) {
      writer.uint32(40).uint64(message.footageDelayMs);
    }
    if (message.footageMinDurationMs !== 0) {
      writer.uint32(48).uint64(message.footageMinDurationMs);
    }
    if (message.footageSpeedMult !== 0) {
      writer.uint32(61).float(message.footageSpeedMult);
    }
    if (message.footageIgnore === true) {
      writer.uint32(64).bool(message.footageIgnore);
    }
    if (message.initialVolumeUl !== 0) {
      writer.uint32(77).float(message.initialVolumeUl);
    }
    if (message.currentVolumeUl !== 0) {
      writer.uint32(85).float(message.currentVolumeUl);
    }
    if (message.name !== "") {
      writer.uint32(90).string(message.name);
    }
    if (message.vialFluid !== 0) {
      writer.uint32(96).int32(message.vialFluid);
    }
    if (message.colour !== "") {
      writer.uint32(106).string(message.colour);
    }
    for (const v of message.aliases) {
      writer.uint32(114).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VialProfile {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVialProfile();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.description = reader.string();
          continue;
        case 3:
          if (tag !== 29) {
            break;
          }

          message.slopUl = reader.float();
          continue;
        case 4:
          if (tag !== 37) {
            break;
          }

          message.dispenseVolumeUl = reader.float();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.footageDelayMs = longToNumber(reader.uint64() as Long);
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.footageMinDurationMs = longToNumber(reader.uint64() as Long);
          continue;
        case 7:
          if (tag !== 61) {
            break;
          }

          message.footageSpeedMult = reader.float();
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.footageIgnore = reader.bool();
          continue;
        case 9:
          if (tag !== 77) {
            break;
          }

          message.initialVolumeUl = reader.float();
          continue;
        case 10:
          if (tag !== 85) {
            break;
          }

          message.currentVolumeUl = reader.float();
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.name = reader.string();
          continue;
        case 12:
          if (tag !== 96) {
            break;
          }

          message.vialFluid = reader.int32() as any;
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.colour = reader.string();
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.aliases.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VialProfile {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      description: isSet(object.description) ? globalThis.String(object.description) : "",
      slopUl: isSet(object.slop_ul) ? globalThis.Number(object.slop_ul) : 0,
      dispenseVolumeUl: isSet(object.dispense_volume_ul) ? globalThis.Number(object.dispense_volume_ul) : 0,
      footageDelayMs: isSet(object.footage_delay_ms) ? globalThis.Number(object.footage_delay_ms) : 0,
      footageMinDurationMs: isSet(object.footage_min_duration_ms)
        ? globalThis.Number(object.footage_min_duration_ms)
        : 0,
      footageSpeedMult: isSet(object.footage_speed_mult) ? globalThis.Number(object.footage_speed_mult) : 0,
      footageIgnore: isSet(object.footage_ignore) ? globalThis.Boolean(object.footage_ignore) : false,
      initialVolumeUl: isSet(object.initial_volume_ul) ? globalThis.Number(object.initial_volume_ul) : 0,
      currentVolumeUl: isSet(object.current_volume_ul) ? globalThis.Number(object.current_volume_ul) : 0,
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      vialFluid: isSet(object.vial_fluid) ? vialProfile_VialFluidFromJSON(object.vial_fluid) : 0,
      colour: isSet(object.colour) ? globalThis.String(object.colour) : "",
      aliases: globalThis.Array.isArray(object?.aliases) ? object.aliases.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: VialProfile): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.slopUl !== 0) {
      obj.slop_ul = message.slopUl;
    }
    if (message.dispenseVolumeUl !== 0) {
      obj.dispense_volume_ul = message.dispenseVolumeUl;
    }
    if (message.footageDelayMs !== 0) {
      obj.footage_delay_ms = Math.round(message.footageDelayMs);
    }
    if (message.footageMinDurationMs !== 0) {
      obj.footage_min_duration_ms = Math.round(message.footageMinDurationMs);
    }
    if (message.footageSpeedMult !== 0) {
      obj.footage_speed_mult = message.footageSpeedMult;
    }
    if (message.footageIgnore === true) {
      obj.footage_ignore = message.footageIgnore;
    }
    if (message.initialVolumeUl !== 0) {
      obj.initial_volume_ul = message.initialVolumeUl;
    }
    if (message.currentVolumeUl !== 0) {
      obj.current_volume_ul = message.currentVolumeUl;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.vialFluid !== 0) {
      obj.vial_fluid = vialProfile_VialFluidToJSON(message.vialFluid);
    }
    if (message.colour !== "") {
      obj.colour = message.colour;
    }
    if (message.aliases?.length) {
      obj.aliases = message.aliases;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<VialProfile>, I>>(base?: I): VialProfile {
    return VialProfile.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<VialProfile>, I>>(object: I): VialProfile {
    const message = createBaseVialProfile();
    message.id = object.id ?? 0;
    message.description = object.description ?? "";
    message.slopUl = object.slopUl ?? 0;
    message.dispenseVolumeUl = object.dispenseVolumeUl ?? 0;
    message.footageDelayMs = object.footageDelayMs ?? 0;
    message.footageMinDurationMs = object.footageMinDurationMs ?? 0;
    message.footageSpeedMult = object.footageSpeedMult ?? 0;
    message.footageIgnore = object.footageIgnore ?? false;
    message.initialVolumeUl = object.initialVolumeUl ?? 0;
    message.currentVolumeUl = object.currentVolumeUl ?? 0;
    message.name = object.name ?? "";
    message.vialFluid = object.vialFluid ?? 0;
    message.colour = object.colour ?? "";
    message.aliases = object.aliases?.map((e) => e) || [];
    return message;
  },
};

function createBaseSystemVialConfiguration(): SystemVialConfiguration {
  return { vials: {} };
}

export const SystemVialConfiguration = {
  encode(message: SystemVialConfiguration, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.vials).forEach(([key, value]) => {
      SystemVialConfiguration_VialsEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SystemVialConfiguration {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSystemVialConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = SystemVialConfiguration_VialsEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.vials[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SystemVialConfiguration {
    return {
      vials: isObject(object.vials)
        ? Object.entries(object.vials).reduce<{ [key: number]: number }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = Number(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SystemVialConfiguration): unknown {
    const obj: any = {};
    if (message.vials) {
      const entries = Object.entries(message.vials);
      if (entries.length > 0) {
        obj.vials = {};
        entries.forEach(([k, v]) => {
          obj.vials[k] = Math.round(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SystemVialConfiguration>, I>>(base?: I): SystemVialConfiguration {
    return SystemVialConfiguration.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SystemVialConfiguration>, I>>(object: I): SystemVialConfiguration {
    const message = createBaseSystemVialConfiguration();
    message.vials = Object.entries(object.vials ?? {}).reduce<{ [key: number]: number }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[globalThis.Number(key)] = globalThis.Number(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseSystemVialConfiguration_VialsEntry(): SystemVialConfiguration_VialsEntry {
  return { key: 0, value: 0 };
}

export const SystemVialConfiguration_VialsEntry = {
  encode(message: SystemVialConfiguration_VialsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint64(message.key);
    }
    if (message.value !== 0) {
      writer.uint32(16).uint64(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SystemVialConfiguration_VialsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSystemVialConfiguration_VialsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.value = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SystemVialConfiguration_VialsEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? globalThis.Number(object.value) : 0,
    };
  },

  toJSON(message: SystemVialConfiguration_VialsEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== 0) {
      obj.value = Math.round(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SystemVialConfiguration_VialsEntry>, I>>(
    base?: I,
  ): SystemVialConfiguration_VialsEntry {
    return SystemVialConfiguration_VialsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SystemVialConfiguration_VialsEntry>, I>>(
    object: I,
  ): SystemVialConfiguration_VialsEntry {
    const message = createBaseSystemVialConfiguration_VialsEntry();
    message.key = object.key ?? 0;
    message.value = object.value ?? 0;
    return message;
  },
};

function createBaseVialProfileCollection(): VialProfileCollection {
  return { profiles: {} };
}

export const VialProfileCollection = {
  encode(message: VialProfileCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.profiles).forEach(([key, value]) => {
      VialProfileCollection_ProfilesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VialProfileCollection {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVialProfileCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = VialProfileCollection_ProfilesEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.profiles[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VialProfileCollection {
    return {
      profiles: isObject(object.profiles)
        ? Object.entries(object.profiles).reduce<{ [key: number]: VialProfile }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = VialProfile.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: VialProfileCollection): unknown {
    const obj: any = {};
    if (message.profiles) {
      const entries = Object.entries(message.profiles);
      if (entries.length > 0) {
        obj.profiles = {};
        entries.forEach(([k, v]) => {
          obj.profiles[k] = VialProfile.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<VialProfileCollection>, I>>(base?: I): VialProfileCollection {
    return VialProfileCollection.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<VialProfileCollection>, I>>(object: I): VialProfileCollection {
    const message = createBaseVialProfileCollection();
    message.profiles = Object.entries(object.profiles ?? {}).reduce<{ [key: number]: VialProfile }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[globalThis.Number(key)] = VialProfile.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseVialProfileCollection_ProfilesEntry(): VialProfileCollection_ProfilesEntry {
  return { key: 0, value: undefined };
}

export const VialProfileCollection_ProfilesEntry = {
  encode(message: VialProfileCollection_ProfilesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint64(message.key);
    }
    if (message.value !== undefined) {
      VialProfile.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VialProfileCollection_ProfilesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVialProfileCollection_ProfilesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = VialProfile.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VialProfileCollection_ProfilesEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? VialProfile.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: VialProfileCollection_ProfilesEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== undefined) {
      obj.value = VialProfile.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<VialProfileCollection_ProfilesEntry>, I>>(
    base?: I,
  ): VialProfileCollection_ProfilesEntry {
    return VialProfileCollection_ProfilesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<VialProfileCollection_ProfilesEntry>, I>>(
    object: I,
  ): VialProfileCollection_ProfilesEntry {
    const message = createBaseVialProfileCollection_ProfilesEntry();
    message.key = object.key ?? 0;
    message.value = (object.value !== undefined && object.value !== null)
      ? VialProfile.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSystemVialConfigurationSnapshot(): SystemVialConfigurationSnapshot {
  return { profiles: {} };
}

export const SystemVialConfigurationSnapshot = {
  encode(message: SystemVialConfigurationSnapshot, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.profiles).forEach(([key, value]) => {
      SystemVialConfigurationSnapshot_ProfilesEntry.encode({ key: key as any, value }, writer.uint32(10).fork())
        .ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SystemVialConfigurationSnapshot {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSystemVialConfigurationSnapshot();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = SystemVialConfigurationSnapshot_ProfilesEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.profiles[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SystemVialConfigurationSnapshot {
    return {
      profiles: isObject(object.profiles)
        ? Object.entries(object.profiles).reduce<{ [key: number]: VialProfile }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = VialProfile.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SystemVialConfigurationSnapshot): unknown {
    const obj: any = {};
    if (message.profiles) {
      const entries = Object.entries(message.profiles);
      if (entries.length > 0) {
        obj.profiles = {};
        entries.forEach(([k, v]) => {
          obj.profiles[k] = VialProfile.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SystemVialConfigurationSnapshot>, I>>(base?: I): SystemVialConfigurationSnapshot {
    return SystemVialConfigurationSnapshot.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SystemVialConfigurationSnapshot>, I>>(
    object: I,
  ): SystemVialConfigurationSnapshot {
    const message = createBaseSystemVialConfigurationSnapshot();
    message.profiles = Object.entries(object.profiles ?? {}).reduce<{ [key: number]: VialProfile }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[globalThis.Number(key)] = VialProfile.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseSystemVialConfigurationSnapshot_ProfilesEntry(): SystemVialConfigurationSnapshot_ProfilesEntry {
  return { key: 0, value: undefined };
}

export const SystemVialConfigurationSnapshot_ProfilesEntry = {
  encode(message: SystemVialConfigurationSnapshot_ProfilesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint64(message.key);
    }
    if (message.value !== undefined) {
      VialProfile.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SystemVialConfigurationSnapshot_ProfilesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSystemVialConfigurationSnapshot_ProfilesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = VialProfile.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SystemVialConfigurationSnapshot_ProfilesEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? VialProfile.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SystemVialConfigurationSnapshot_ProfilesEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== undefined) {
      obj.value = VialProfile.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SystemVialConfigurationSnapshot_ProfilesEntry>, I>>(
    base?: I,
  ): SystemVialConfigurationSnapshot_ProfilesEntry {
    return SystemVialConfigurationSnapshot_ProfilesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SystemVialConfigurationSnapshot_ProfilesEntry>, I>>(
    object: I,
  ): SystemVialConfigurationSnapshot_ProfilesEntry {
    const message = createBaseSystemVialConfigurationSnapshot_ProfilesEntry();
    message.key = object.key ?? 0;
    message.value = (object.value !== undefined && object.value !== null)
      ? VialProfile.fromPartial(object.value)
      : undefined;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
