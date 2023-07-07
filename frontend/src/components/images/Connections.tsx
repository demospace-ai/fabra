import bigquery from "src/components/images/bigquery.svg";
import mongodb from "src/components/images/mongodb.svg";
import mysql from "src/components/images/mysql.svg";
import postgres from "src/components/images/postgres.svg";
import redshift from "src/components/images/redshift.svg";
import snowflake from "src/components/images/snowflake.svg";
import synapse from "src/components/images/synapse.svg";
import webhook from "src/components/images/webhook.svg";
import dynamodb from "src/components/images/dynamodb.svg";
import { ConnectionType } from "src/rpc/api";

function getConnectionTypeImg(connectionType: ConnectionType): string {
  switch (connectionType) {
    case ConnectionType.BigQuery:
      return bigquery;
    case ConnectionType.DynamoDb:
      return dynamodb;
    case ConnectionType.Snowflake:
      return snowflake;
    case ConnectionType.Redshift:
      return redshift;
    case ConnectionType.MongoDb:
      return mongodb;
    case ConnectionType.Synapse:
      return synapse;
    case ConnectionType.Postgres:
      return postgres;
    case ConnectionType.MySQL:
      return mysql;
    case ConnectionType.Webhook:
      return webhook;
  }
}

export const ConnectionImage: React.FC<{ connectionType: ConnectionType; className?: string }> = ({
  connectionType,
  className,
}) => {
  return <img src={getConnectionTypeImg(connectionType)} alt="data source logo" className={className} />;
};
