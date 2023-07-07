// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { z } from "zod";

export enum FrequencyUnits {
  Minutes = "minutes",
  Hours = "hours",
  Days = "days",
  Weeks = "weeks",
}

export enum SyncMode {
  FullOverwrite = "full_overwrite",
  IncrementalAppend = "incremental_append",
  IncrementalUpdate = "incremental_update",
}

export enum TargetType {
  SingleExisting = "single_existing",
  SingleNew = "single_new",
  TablePerCustomer = "table_per_customer",
  Webhook = "webhook",
}

export enum FieldType {
  String = "STRING",
  Integer = "INTEGER",
  Number = "NUMBER",
  Timestamp = "TIMESTAMP",
  TimeTz = "TIME_TZ",
  TimeNtz = "TIME_NTZ",
  DatetimeTz = "DATETIME_TZ",
  DatetimeNtz = "DATETIME_NTZ",
  Date = "DATE",
  Boolean = "BOOLEAN",
  Array = "ARRAY",
  Json = "JSON",
}

export const FieldSchema = z.object({
  name: z.string(),
  type: z.nativeEnum(FieldType),
});

export type Field = z.infer<typeof FieldSchema>;

export interface IEndpoint<RequestType, ResponseType> {
  name: string;
  method: "GET" | "POST" | "DELETE" | "PUT" | "PATCH";
  path: string;
  track?: boolean;
  queryParams?: string[]; // These will be used as query params instead of being used as path params
  noJson?: boolean; // TODO: do this better
}

export const Login: IEndpoint<LoginRequest, LoginResponse> = {
  name: "Login",
  method: "POST",
  path: "/login",
  track: true,
};

export const OAuthRedirect: IEndpoint<{ provider: OAuthProvider }, undefined> = {
  name: "OAuth Redirect",
  method: "GET",
  path: "/oauth_redirect",
  queryParams: ["provider"],
  track: true,
};

export const GetAllUsers: IEndpoint<undefined, GetAllUsersResponse> = {
  name: "All Users Fetched",
  method: "GET",
  path: "/users",
};

export const GetDestinations: IEndpoint<undefined, GetDestinationsResponse> = {
  name: "Destinations Fetched",
  method: "GET",
  path: "/destinations",
};

export const GetDestination: IEndpoint<{ destinationID: number }, GetDestinationResponse> = {
  name: "Destination Fetched",
  method: "GET",
  path: "/destination/:destinationID",
};

export const GetSources: IEndpoint<undefined, GetSourcesResponse> = {
  name: "Sources Fetched",
  method: "GET",
  path: "/sources",
};

export const GetObjects: IEndpoint<{ destinationID?: number }, GetObjectsResponse> = {
  name: "Objects Fetched",
  method: "GET",
  path: "/objects",
  queryParams: ["destinationID"],
};

export const GetObject: IEndpoint<{ objectID: number }, GetObjectResponse> = {
  name: "Object Fetched",
  method: "GET",
  path: "/object/:objectID",
};

export const UpdateObject: IEndpoint<UpdateObjectRequest, UpdateObjectResponse> = {
  name: "Object Updated",
  method: "PATCH",
  path: "/object/:objectID",
};

export const UpdateObjectFields: IEndpoint<UpdateObjectFieldsRequest, UpdateObjectFieldsResponse> = {
  name: "ObjectField Updated",
  method: "PATCH",
  path: "/object/:objectID/object_fields",
};

export const GetSyncs: IEndpoint<undefined, GetSyncsResponse> = {
  name: "Syncs Fetched",
  method: "GET",
  path: "/syncs",
};

export const GetSync: IEndpoint<{ syncID: number }, GetSyncResponse> = {
  name: "Sync Fetched",
  method: "GET",
  path: "/sync/:syncID",
  track: true,
};

export const GetNamespaces: IEndpoint<{ connectionID: number }, GetNamespacesResponse> = {
  name: "Namespaces Fetched",
  method: "GET",
  path: "/connection/namespaces",
  queryParams: ["connectionID"],
};

export const GetTables: IEndpoint<{ connectionID: number; namespace: string }, GetTablesResponse> = {
  name: "Tables Fetched",
  method: "GET",
  path: "/connection/tables",
  queryParams: ["connectionID", "namespace"],
};

export const GetSchema: IEndpoint<GetSchemaRequest, GetSchemaResponse> = {
  name: "Schema Fetched",
  method: "GET",
  path: "/connection/schema",
  queryParams: ["connectionID", "namespace", "tableName", "customJoin"],
};

export const LinkGetPreview: IEndpoint<LinkGetPreviewRequest, QueryResults> = {
  name: "Get Preview",
  method: "POST",
  path: "/link/preview",
};

export const LinkGetSources: IEndpoint<undefined, GetSourcesResponse> = {
  name: "Sources Fetched",
  method: "GET",
  path: "/link/sources",
};

export const LinkGetNamespaces: IEndpoint<{ sourceID: number }, GetNamespacesResponse> = {
  name: "Namespaces Fetched",
  method: "GET",
  path: "/link/namespaces",
  queryParams: ["sourceID"],
};

export const LinkGetTables: IEndpoint<{ sourceID: number; namespace: string }, GetTablesResponse> = {
  name: "Tables Fetched",
  method: "GET",
  path: "/link/tables",
  queryParams: ["sourceID", "namespace"],
};

export const LinkGetSchema: IEndpoint<LinkGetSchemaRequest, GetSchemaResponse> = {
  name: "Schema Fetched",
  method: "GET",
  path: "/link/schema",
  queryParams: ["sourceID", "namespace", "tableName", "customJoin"],
};

export const GetApiKey: IEndpoint<undefined, string> = {
  name: "API Key Fetched",
  method: "GET",
  path: "/api_key",
  noJson: true,
};

export const CheckSession: IEndpoint<undefined, CheckSessionResponse> = {
  name: "Session Checked",
  method: "GET",
  path: "/check_session",
};

export const Logout: IEndpoint<undefined, undefined> = {
  name: "Logout",
  method: "DELETE",
  path: "/logout",
  track: true,
};

export const SetOrganization: IEndpoint<SetOrganizationRequest, SetOrganizationResponse> = {
  name: "Organization Set",
  method: "POST",
  path: "/organization",
  track: true,
};

export const TestDataConnection: IEndpoint<TestDataConnectionRequest, undefined> = {
  name: "Test Data Connection",
  method: "POST",
  path: "/connection/test",
};

export const GetFieldValues: IEndpoint<GetFieldValuesRequest, GetFieldValuesResponse> = {
  name: "Field Values Fetched",
  method: "GET",
  path: "/connection/field_values",
  queryParams: ["connectionID", "namespace", "tableName", "fieldName"],
  track: true,
};

export const CreateDestination: IEndpoint<CreateDestinationRequest, CreateDestinationResponse> = {
  name: "Destination Created",
  method: "POST",
  path: "/destination",
  track: true,
};

export const LinkCreateSource: IEndpoint<LinkCreateSourceRequest, CreateSourceResponse> = {
  name: "Source Created",
  method: "POST",
  path: "/link/source",
  track: true,
};

export const RunSync: IEndpoint<undefined, RunSyncResponse> = {
  name: "Sync Run",
  method: "POST",
  path: "/sync/:syncID/run",
  track: true,
};

export const LinkCreateSync: IEndpoint<LinkCreateSyncRequest, CreateSyncResponse> = {
  name: "Sync Created",
  method: "POST",
  path: "/link/sync",
  track: true,
};

export const LinkGetSyncs: IEndpoint<undefined, GetSyncsResponse> = {
  name: "Syncs Fetched",
  method: "GET",
  path: "/link/syncs",
  track: true,
};

export const LinkGetSync: IEndpoint<{ syncID: number }, GetSyncResponse> = {
  name: "Sync Fetched",
  method: "GET",
  path: "/link/sync/:syncID",
  track: true,
};

export const LinkRunSync: IEndpoint<{ syncID: string }, RunSyncResponse> = {
  name: "Sync Run",
  method: "POST",
  path: "/link/sync/:syncID/run",
  track: true,
};

export const CreateObject: IEndpoint<CreateObjectRequest, CreateObjectResponse> = {
  name: "Object Created",
  method: "POST",
  path: "/object",
  track: true,
};

export const CreateLinkToken: IEndpoint<CreateLinkTokenRequest, CreateLinkTokenResponse> = {
  name: "Link Token Created",
  method: "POST",
  path: "/link_token",
  track: true,
};

export interface TestDataConnectionRequest {
  display_name: string;
  connection_type: ConnectionType;
  bigquery_config?: BigQueryConfig;
  snowflake_config?: SnowflakeConfig;
  mongodb_config?: MongoDbConfig;
  redshift_config?: RedshiftConfig;
  synapse_config?: SynapseConfig;
  postgres_config?: PostgresConfig;
  mysql_config?: MySqlConfig;
  webhook_config?: WebhookConfig;
  dynamodb_config?: CreateDynamoDbConfig;
}

export interface CreateDestinationRequest {
  display_name: string;
  connection_type: ConnectionType;
  staging_bucket?: string;
  bigquery_config?: BigQueryConfig;
  snowflake_config?: SnowflakeConfig;
  redshift_config?: RedshiftConfig;
  mongodb_config?: MongoDbConfig;
  synapse_config?: SynapseConfig;
  webhook_config?: WebhookConfig;
  postgres_config?: PostgresConfig;
  dynamodb_config?: CreateDynamoDbConfig;
}

export interface CreateDestinationResponse {
  destination: Destination;
}

export interface LinkCreateSourceRequest {
  display_name: string;
  connection_type: ConnectionType;
  bigquery_config?: BigQueryConfig;
  snowflake_config?: SnowflakeConfig;
  redshift_config?: RedshiftConfig;
  synapse_config?: SynapseConfig;
  mongodb_config?: MongoDbConfig;
  postgres_config?: PostgresConfig;
  mysql_config?: MySqlConfig;
}

export interface CreateSourceResponse {
  source: Source;
}

export interface CreateObjectRequest {
  display_name: string;
  destination_id: number;
  target_type: TargetType;
  namespace?: string;
  table_name?: string;
  end_customer_id_field?: string;
  sync_mode: SyncMode;
  recurring: boolean;
  frequency?: number;
  frequency_units?: FrequencyUnits;
  object_fields: ObjectFieldInput[];
  cursor_field?: string; // required for incremental append: need cursor field to detect new data
  primary_key?: string; // required  for incremental update: need primary key to match up rows
}

export interface CreateObjectResponse {
  object: FabraObject;
}

export type UpdateObjectRequest = Partial<Omit<CreateObjectRequest, "object_fields">> & {
  objectID: number;
};

export interface CreateLinkTokenRequest {
  end_customer_id: string;
  destination_ids?: number[];
}

export interface CreateLinkTokenResponse {
  link_token: string;
}

export const ObjectFieldSchema = z.object({
  id: z.number(),
  name: z.string().min(1, { message: "Field name must be at least 1 character long" }),
  field_type: z.nativeEnum(FieldType), // Use field type because "type" is already used in Zod objects
  omit: z.boolean(),
  optional: z.boolean(),
  display_name: z.string().optional(),
  description: z.string().optional(),
});

export type ObjectField = z.infer<typeof ObjectFieldSchema>;

export type ObjectFieldInput = Partial<ObjectField>;

export interface FieldMappingInput {
  source_field_name: string;
  source_field_type: FieldType;
  destination_field_id: number;
  is_json_field: boolean;
}

export interface FieldMapping {
  source_field_name: string;
  source_field_type: FieldType;
  destination_field_name: string;
  destination_field_type: FieldType;
  is_json_field: boolean;
}

export type GCPLocation = {
  name: string;
  code: string;
};

export interface BigQueryConfigState {
  credentials: string;
  location: GCPLocation | undefined;
}

export const AwsLocationSchema = z.object({
  name: z.string(),
  code: z.string(),
});

export type AwsLocation = z.infer<typeof AwsLocationSchema>;

export const DynamoDbConfigSchema = z.object({
  accessKey: z.string(),
  secretKey: z.string(),
  region: AwsLocationSchema,
});

export const CreateDynamoDbConfigSchema = z.object({
  access_key: z.string(),
  secret_key: z.string(),
  region: z.string(),
});

export type CreateDynamoDbConfig = z.infer<typeof CreateDynamoDbConfigSchema>;

export type DynamoDbConfig = z.infer<typeof DynamoDbConfigSchema>;

export interface BigQueryConfig {
  credentials: string;
  location: string;
}

export interface SnowflakeConfig {
  username: string;
  password: string;
  database_name: string;
  warehouse_name: string;
  role: string;
  host: string;
}

export interface RedshiftConfig {
  username: string;
  password: string;
  database_name: string;
  endpoint: string;
}

export interface PostgresConfig {
  username: string;
  password: string;
  database_name: string;
  endpoint: string;
}

export interface MySqlConfig {
  username: string;
  password: string;
  database_name: string;
  endpoint: string;
}

export interface SynapseConfig {
  username: string;
  password: string;
  database_name: string;
  endpoint: string;
}

export interface MongoDbConfig {
  username: string;
  password: string;
  host: string;
  connection_options: string;
}

export interface WebhookConfig {
  url: string;
  headers: HeaderInput[];
}

export interface HeaderInput {
  name: string;
  value: string;
}

export interface RunSyncResponse {
  sync: Sync;
  next_run_time: string;
  result: "success" | "failure";
}

export interface LinkCreateSyncRequest {
  display_name: string;
  source_id: number;
  object_id: number;
  namespace?: string;
  table_name?: string;
  custom_join?: string;
  field_mappings: FieldMappingInput[];
  // remaining fields will inherit from object if empty
  source_cursor_field?: string;
  source_primary_key?: string;
  sync_mode?: SyncMode;
  recurring?: boolean;
  frequency?: number;
  frequency_units?: FrequencyUnits;
}

export interface CreateSyncResponse {
  sync: Sync;
}

export interface GetSyncsResponse {
  syncs: Sync[];
  sources: Source[];
  objects: FabraObject[];
}

export interface GetSyncResponse {
  sync: Sync;
  field_mappings: FieldMapping[];
  next_run_time: string;
  sync_runs: SyncRun[];
}

export type JSONValue = string | number | boolean | JSONObject | JSONArray;

export interface JSONObject {
  [x: string]: JSONValue;
}

export interface JSONArray extends Array<JSONValue> {}

export interface ResultRow extends Array<string | number> {}

export interface Schema extends Array<Field> {}

export interface LinkGetPreviewRequest {
  source_id: number;
  namespace: string;
  table_name: string;
}

export interface GetFieldValuesRequest {
  connectionID: number;
  namespace: string;
  tableName: string;
  fieldName: string;
}

export interface GetFieldValuesResponse {
  field_values: string[];
}

export interface SetOrganizationRequest {
  organization_name?: string;
  organization_id?: number;
}

export interface SetOrganizationResponse {
  organization: Organization;
}

export interface ValidationCodeRequest {
  email: string;
}

export interface LoginRequest {
  code?: string;
  state?: string;
  organization_name?: string;
  organization_id?: string;
}

export enum OAuthProvider {
  Google = "google",
  Github = "github",
}

export interface User {
  id: number;
  name: string;
  email: string;
  intercom_hash: string;
}

export interface GetAllUsersResponse {
  users: User[];
}

export interface GetDestinationsResponse {
  destinations: Destination[];
}

export interface GetDestinationResponse {
  destination: Destination;
}

export interface GetSourcesResponse {
  sources: Source[];
}

export interface GetObjectsResponse {
  objects: FabraObject[];
}

export interface GetObjectResponse {
  object: FabraObject;
}

export type UpdateObjectFieldsRequest = {
  objectID: number;
  object_fields: (ObjectFieldInput & { id: number })[];
};

export interface UpdateObjectFieldsResponse {
  object_fields: FabraObject;
  failures: FabraObject[];
}

export type CreateObjectFieldsRequest = ObjectField[];

export interface CreateObjectFieldsResponse {
  object_fields: FabraObject;
  failures: FabraObject[];
}

export interface UpdateObjectResponse {
  object: FabraObject;
}

export interface FabraObject {
  id: number;
  display_name: string;
  destination_id: number;
  target_type: TargetType;
  namespace?: string;
  table_name?: string;
  custom_join?: string;
  object_fields: ObjectField[];
  end_customer_id_field: string;
  sync_mode: SyncMode;
  recurring: boolean;
  frequency: number;
  frequency_units: FrequencyUnits;
  cursor_field?: string;
  primary_key?: string;
}

export interface GetNamespacesResponse {
  namespaces: string[];
}

export interface GetTablesResponse {
  tables: string[];
}

export interface GetSchemaRequest {
  connectionID: number;
  namespace?: string;
  tableName?: string;
  customJoin?: string;
}

export interface LinkGetSchemaRequest {
  sourceID: number;
  namespace?: string;
  tableName?: string;
  customJoin?: string;
}

export interface QueryResults {
  schema: Schema;
  data: ResultRow[];
}

export interface GetSchemaResponse {
  schema: Schema;
}

export interface LoginResponse {
  user: User;
  organization?: Organization;
  suggested_organizations?: Organization[];
}

export interface Organization {
  id: number;
  name: string;
  free_trial_end?: string;
}

export interface CheckSessionResponse {
  user: User;
  organization?: Organization;
  suggested_organizations?: Organization[];
}

export enum ConnectionType {
  BigQuery = "bigquery",
  Snowflake = "snowflake",
  Redshift = "redshift",
  MongoDb = "mongodb",
  Synapse = "synapse",
  Postgres = "postgres",
  MySQL = "mysql",
  Webhook = "webhook",
  DynamoDb = "dynamodb",
}

export const ConnectionSchema = z.object({
  id: z.number(),
  connection_type: z.nativeEnum(ConnectionType),
});

export type Connection = z.infer<typeof ConnectionSchema>;

export const DestinationSchema = z.object({
  id: z.number(),
  display_name: z.string(),
  connection: ConnectionSchema,
  webhook_signing_key: z.string().optional(),
});

export type Destination = z.infer<typeof DestinationSchema>;

export const CreateWarehouseObjectSchema = z.object({
  displayName: z.string(),
  destination: DestinationSchema,
  targetType: z.enum([TargetType.SingleExisting, TargetType.SingleNew, TargetType.TablePerCustomer]),
  namespace: z.string(),
  tableName: z.string(),
  endCustomerIdField: FieldSchema,
  syncMode: z.nativeEnum(SyncMode),
  frequency: z.number(),
  frequencyUnits: z.nativeEnum(FrequencyUnits),
  objectFields: z.array(
    z.object({
      name: z.string(),
      type: z.nativeEnum(FieldType),
      omit: z.boolean(),
      optional: z.boolean(),
      displayName: z.string().optional(),
      description: z.string().optional(),
    }),
  ),
  cursorField: FieldSchema.optional(),
  primaryKey: FieldSchema.optional(),
});

export const CreateWebhookObjectSchema = z.object({
  targetType: z.enum([TargetType.Webhook]),
  namespace: z.string().optional(),
  tableName: z.string().optional(),
  cursorField: FieldSchema.optional(),
  primaryKey: FieldSchema.optional(),
  endCustomerIdField: z.object({
    name: z.string(),
    type: z.nativeEnum(FieldType),
  }),
  objectFields: z.array(z.any()),
  destination: DestinationSchema,
  displayName: z.string(),
  syncMode: z.nativeEnum(SyncMode),
  frequency: z.number(),
  frequencyUnits: z.nativeEnum(FrequencyUnits),
});

export const CreateObjectSchema = z.discriminatedUnion("targetType", [
  CreateWarehouseObjectSchema,
  CreateWebhookObjectSchema,
]);

export type CreateObjectSchemaType = z.infer<typeof CreateObjectSchema>;
export interface Source {
  id: number;
  display_name: string;
  connection: Connection;
  end_customer_id: string;
}

export interface Sync {
  id: number;
  display_name: string;
  end_customer_id: string;
  source_id: number;
  object_id: number;
  namespace: string | undefined;
  table_name: string | undefined;
  custom_join: string | undefined;
  // the following settings will override object settings if set
  source_cursor_field: string | undefined;
  source_primary_key: string | undefined;
  sync_mode: SyncMode | undefined;
  frequency: number | undefined;
  frequency_units: FrequencyUnits | undefined;
}

export interface SyncRun {
  id: number;
  sync_id: number;
  status: SyncRunStatus;
  error: string | undefined;
  started_at: string;
  completed_at: string;
  duration: string | undefined;
  rows_written: number;
}

export enum SyncRunStatus {
  Running = "running",
  Failed = "failed",
  Completed = "completed",
}

export const needsCursorField = (syncMode: SyncMode): boolean => {
  switch (syncMode) {
    case SyncMode.FullOverwrite:
      return false;
    case SyncMode.IncrementalAppend:
      return true;
    case SyncMode.IncrementalUpdate:
      return true;
  }
};

export const needsPrimaryKey = (syncMode: SyncMode): boolean => {
  switch (syncMode) {
    case SyncMode.FullOverwrite:
      return false;
    case SyncMode.IncrementalAppend:
      return false;
    case SyncMode.IncrementalUpdate:
      return true;
  }
};

export const needsEndCustomerId = (targetType: TargetType): boolean => {
  switch (targetType) {
    case TargetType.Webhook:
      return false;
    case TargetType.SingleExisting:
    case TargetType.SingleNew:
    case TargetType.TablePerCustomer:
      return true;
  }
};

export const shouldCreateFields = (connectionType: ConnectionType, targetType: TargetType): boolean => {
  if (connectionType === ConnectionType.Webhook || connectionType === ConnectionType.DynamoDb) {
    return true;
  }

  if (targetType === TargetType.SingleNew) {
    return true;
  }

  return false;
};

export function targetTypeToString(targetType: TargetType) {
  switch (targetType) {
    case TargetType.SingleExisting:
      return "Single Existing Table";
    case TargetType.SingleNew:
      return "Single New Table";
    case TargetType.TablePerCustomer:
      return "Table Per Customer";
    case TargetType.Webhook:
      return "Webhook";
  }
}

export function syncModeToString(syncMode: SyncMode) {
  switch (syncMode) {
    case SyncMode.FullOverwrite:
      return "Full Overwrite";
    case SyncMode.IncrementalAppend:
      return "Incremental Append";
    case SyncMode.IncrementalUpdate:
      return "Incremental Update";
  }
}

export function toReadableFrequency(frequency: number, frequencyUnits: FrequencyUnits): string {
  if (frequency === 1) {
    return `Every ${toSingularTime(frequencyUnits).toLowerCase()}`;
  }

  return `Every ${frequency} ${frequencyUnits.toLowerCase()}`;
}

export function toSingularTime(frequencyUnits: FrequencyUnits): string {
  switch (frequencyUnits) {
    case FrequencyUnits.Minutes:
      return "Minute";
    case FrequencyUnits.Hours:
      return "Hour";
    case FrequencyUnits.Days:
      return "Day";
    case FrequencyUnits.Weeks:
      return "Week";
  }
}

export function getConnectionType(connectionType: ConnectionType): string {
  switch (connectionType) {
    case ConnectionType.DynamoDb:
      return "DynamoDB";
    case ConnectionType.BigQuery:
      return "BigQuery";
    case ConnectionType.Snowflake:
      return "Snowflake";
    case ConnectionType.Redshift:
      return "Redshift";
    case ConnectionType.MongoDb:
      return "MongoDB";
    case ConnectionType.Synapse:
      return "Synapse";
    case ConnectionType.Postgres:
      return "Postgres";
    case ConnectionType.MySQL:
      return "MySQL";
    case ConnectionType.Webhook:
      return "Webhook";
  }
}
