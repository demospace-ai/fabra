import {
  ConnectionType,
  BigQueryConfigState,
  SnowflakeConfig,
  RedshiftConfig,
  SynapseConfig,
  MongoDbConfig,
  WebhookConfig,
  PostgresConfig,
  DynamoDbConfig,
} from "src/rpc/api";

export type NewConnectionConfigurationProps = {
  connectionType: ConnectionType;
  setConnectionType: (connectionType: ConnectionType | null) => void;
};

export type NewDestinationState = {
  displayName: string;
  staging_bucket: string;
  bigqueryConfig: BigQueryConfigState;
  snowflakeConfig: SnowflakeConfig;
  redshiftConfig: RedshiftConfig;
  synapseConfig: SynapseConfig;
  mongodbConfig: MongoDbConfig;
  webhookConfig: WebhookConfig;
  postgresConfig: PostgresConfig;
  dynamoDbConfig: Partial<DynamoDbConfig>;
  error: string | undefined;
};
