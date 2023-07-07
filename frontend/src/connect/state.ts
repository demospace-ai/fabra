import { useNavigate } from "react-router-dom";
import { ShowToastFunction, useConnectShowToast } from "src/components/notifications/Notifications";
import { sendLinkTokenRequest } from "src/rpc/ajax";
import {
  BigQueryConfigState,
  ConnectionType,
  FabraObject,
  Field,
  FieldMappingInput,
  FrequencyUnits,
  GetSources,
  LinkCreateSource,
  LinkCreateSourceRequest,
  LinkCreateSync,
  LinkCreateSyncRequest,
  LinkGetSources,
  LinkGetSyncs,
  MongoDbConfig,
  MySqlConfig,
  ObjectField,
  PostgresConfig,
  RedshiftConfig,
  SnowflakeConfig,
  Source,
  SynapseConfig,
} from "src/rpc/api";
import { consumeError, HttpError } from "src/utils/errors";
import { mutate } from "swr";

export type SetupSyncProps = {
  linkToken: string;
  state: SetupSyncState;
  setState: React.Dispatch<React.SetStateAction<SetupSyncState>>;
};

export enum SyncSetupStep {
  ExistingSources = 1,
  ChooseSourceType,
  ConnectionDetails,
  ChooseData,
  Finalize,
}

export type NewSourceState = {
  sourceCreated: boolean;
  error: string | undefined;
  displayName: string;
  bigqueryConfig: BigQueryConfigState;
  snowflakeConfig: SnowflakeConfig;
  redshiftConfig: RedshiftConfig;
  synapseConfig: SynapseConfig;
  postgresConfig: PostgresConfig;
  mysqlConfig: MySqlConfig;
  mongodbConfig: MongoDbConfig;
};

// Values must be empty strings otherwise the input will be uncontrolled
export const INITIAL_SOURCE_STATE: NewSourceState = {
  sourceCreated: false,
  error: undefined,
  displayName: "",
  bigqueryConfig: {
    credentials: "",
    location: undefined,
  },
  snowflakeConfig: {
    username: "",
    password: "",
    database_name: "",
    warehouse_name: "",
    role: "",
    host: "",
  },
  redshiftConfig: {
    username: "",
    password: "",
    database_name: "",
    endpoint: "",
  },
  synapseConfig: {
    username: "",
    password: "",
    database_name: "",
    endpoint: "",
  },
  mongodbConfig: {
    username: "",
    password: "",
    host: "",
    connection_options: "",
  },
  postgresConfig: {
    username: "",
    password: "",
    database_name: "",
    endpoint: "",
  },
  mysqlConfig: {
    username: "",
    password: "",
    database_name: "",
    endpoint: "",
  },
};

export const resetState = (setState: React.Dispatch<React.SetStateAction<SetupSyncState>>) => {
  setState((_) => {
    return INITIAL_SETUP_STATE;
  });
};

export interface FieldMappingState {
  sourceField: Field | undefined;
  destinationField: ObjectField;
  expandedJson: boolean;
  jsonFields: (Field | undefined)[];
}

export type SetupSyncState = {
  step: SyncSetupStep;
  error: string | undefined;
  skippedSourceSetup: boolean;
  skippedSourceSelection: boolean;
  object: FabraObject | undefined;
  namespace: string | undefined;
  tableName: string | undefined;
  customJoin: string | undefined;
  connectionType: ConnectionType | undefined;
  source: Source | undefined;
  newSourceState: NewSourceState;
  displayName: string | undefined;
  frequency: number | undefined;
  frequencyUnits: FrequencyUnits | undefined;
  fieldMappings: FieldMappingState[] | undefined;
};

export const INITIAL_SETUP_STATE: SetupSyncState = {
  step: SyncSetupStep.ExistingSources,
  error: undefined,
  skippedSourceSetup: false,
  skippedSourceSelection: false,
  object: undefined,
  namespace: undefined,
  tableName: undefined,
  customJoin: undefined,
  connectionType: undefined,
  source: undefined,
  newSourceState: INITIAL_SOURCE_STATE,
  displayName: undefined,
  frequency: undefined,
  frequencyUnits: undefined,
  fieldMappings: undefined,
};

export const validateConnectionSetup = (connectionType: ConnectionType | undefined, state: NewSourceState): boolean => {
  if (!connectionType) {
    consumeError(new Error("Connection type not set for source setup"));
    return false;
  }

  switch (connectionType) {
    case ConnectionType.Snowflake:
      return (
        state.displayName.length > 0 &&
        state.snowflakeConfig.username.length > 0 &&
        state.snowflakeConfig.password.length > 0 &&
        state.snowflakeConfig.database_name.length > 0 &&
        state.snowflakeConfig.warehouse_name.length > 0 &&
        state.snowflakeConfig.role.length > 0 &&
        state.snowflakeConfig.host.length > 0
      );
    case ConnectionType.BigQuery:
      return state.displayName.length > 0 && state.bigqueryConfig.credentials.length > 0;
    case ConnectionType.Redshift:
      return (
        state.displayName.length > 0 &&
        state.redshiftConfig.username.length > 0 &&
        state.redshiftConfig.password.length > 0 &&
        state.redshiftConfig.database_name.length > 0 &&
        state.redshiftConfig.endpoint.length > 0
      );
    case ConnectionType.Synapse:
      return (
        state.displayName.length > 0 &&
        state.synapseConfig.username.length > 0 &&
        state.synapseConfig.password.length > 0 &&
        state.synapseConfig.database_name.length > 0 &&
        state.synapseConfig.endpoint.length > 0
      );
    case ConnectionType.MongoDb:
      return (
        state.displayName.length > 0 &&
        state.mongodbConfig.username.length > 0 &&
        state.mongodbConfig.password.length > 0 &&
        state.mongodbConfig.host.length > 0
      ); // connection options is optional
    case ConnectionType.Postgres:
      return (
        state.displayName.length > 0 &&
        state.postgresConfig.username.length > 0 &&
        state.postgresConfig.password.length > 0 &&
        state.postgresConfig.database_name.length > 0 &&
        state.postgresConfig.endpoint.length > 0
      );
    case ConnectionType.MySQL:
      return (
        state.displayName.length > 0 &&
        state.mysqlConfig.username.length > 0 &&
        state.mysqlConfig.password.length > 0 &&
        state.mysqlConfig.database_name.length > 0 &&
        state.mysqlConfig.endpoint.length > 0
      );
    case ConnectionType.Webhook:
      return false; // cannot create a sync with a webhook source
    case ConnectionType.DynamoDb:
      return false; // TODO: DynamoDB not supported as a source yet
  }
};

export const createNewSource = async (
  linkToken: string,
  state: SetupSyncState,
  setState: React.Dispatch<React.SetStateAction<SetupSyncState>>,
) => {
  if (!validateConnectionSetup(state.connectionType, state.newSourceState)) {
    // TODO: make each required input field red if it's not filled out
    setState((state) => {
      return { ...state, newSourceState: { ...state.newSourceState, error: "Must fill out all required fields" } };
    });
    return;
  }

  if (state.newSourceState.sourceCreated) {
    // TODO: clear success if one of the inputs change and just update the already created source
    // Already created the source, just continue again
    setState((state) => ({ ...state, step: SyncSetupStep.ChooseData }));
    return;
  }

  const payload: LinkCreateSourceRequest = {
    display_name: state.newSourceState.displayName,
    connection_type: state.connectionType!,
  };

  switch (state.connectionType!) {
    case ConnectionType.BigQuery:
      payload.bigquery_config = {
        location: state.newSourceState.bigqueryConfig.location!.code,
        credentials: state.newSourceState.bigqueryConfig.credentials,
      };
      break;
    case ConnectionType.Snowflake:
      payload.snowflake_config = state.newSourceState.snowflakeConfig;
      break;
    case ConnectionType.Redshift:
      payload.redshift_config = state.newSourceState.redshiftConfig;
      break;
    case ConnectionType.Synapse:
      payload.synapse_config = state.newSourceState.synapseConfig;
      break;
    case ConnectionType.MongoDb:
      payload.mongodb_config = state.newSourceState.mongodbConfig;
      break;
    case ConnectionType.Postgres:
      payload.postgres_config = state.newSourceState.postgresConfig;
      break;
    case ConnectionType.MySQL:
      payload.mysql_config = state.newSourceState.mysqlConfig;
      break;
    case ConnectionType.DynamoDb:
      // TODO: throw an error
      return;
    case ConnectionType.Webhook:
      // TODO: throw an error
      return;
  }

  try {
    const response = await sendLinkTokenRequest(LinkCreateSource, linkToken, payload);
    // Tell SWRs to refetch sources
    mutate({ GetSources });
    mutate({ LinkGetSources }); // Tell SWRs to refetch sources
    setState((state) => ({
      ...state,
      source: response.source,
      step: SyncSetupStep.ChooseData,
      newSourceState: { ...state.newSourceState, sourceCreated: true },
      namespace: undefined, // set namespace and table name to undefined since we"re using a new source
      tableName: undefined,
    }));
  } catch (e) {
    if (e instanceof HttpError) {
      const errorMessage = e.message;
      setState((state) => ({ ...state, newSourceState: { ...state.newSourceState, error: errorMessage } }));
    }
    consumeError(e);
  }
};

export const validateObjectSetup = (state: SetupSyncState, showToast: ShowToastFunction): boolean => {
  if (state.object === undefined) {
    showToast("error", "Must choose an object to sync.", 5000);
    return false;
  }

  if (state.namespace === undefined) {
    showToast("error", "Must choose a source namespace.", 5000);
    return false;
  }

  if (state.tableName === undefined) {
    showToast("error", "Must choose a source table.", 5000);
    return false;
  }

  return true;
};

export const validateSyncSetup = (state: SetupSyncState, showToast: ShowToastFunction): boolean => {
  if (state.displayName === undefined || state.displayName.length <= 0) {
    showToast("error", "Must set a display name.", 5000);
    return false;
  }
  if (state.source === undefined) {
    showToast("error", "Must choose a source.", 5000);
    return false;
  }
  if (state.object === undefined) {
    showToast("error", "Must choose a destination object.", 5000);
    return false;
  }
  if (state.namespace === undefined && state.tableName === undefined && state.customJoin === undefined) {
    showToast("error", "Must configure the source namespace and table.", 5000);
    return false;
  }
  // TODO: validate frequency once we allow end customers to customize this
  //if (state.frequency === undefined) return false;
  //if (state.frequencyUnits === undefined) return false;
  if (!validateFieldMappings(state.fieldMappings)) {
    showToast("error", "Field mappings are invalid.", 5000);
    return false;
  }

  return true;
};

const validateFieldMappings = (fieldMappings: FieldMappingState[] | undefined): boolean => {
  if (fieldMappings === undefined) {
    return false;
  }

  return fieldMappings.every((fieldMapping) => {
    if (fieldMapping.expandedJson) {
      return fieldMapping.jsonFields.every((jsonField) => jsonField !== undefined);
    } else {
      return fieldMapping.destinationField.optional || fieldMapping.sourceField !== undefined;
    }
  });
};

export const useCreateNewSync = () => {
  const navigate = useNavigate();
  const showToast = useConnectShowToast();
  return async (
    linkToken: string,
    state: SetupSyncState,
    setState: React.Dispatch<React.SetStateAction<SetupSyncState>>,
  ) => {
    if (!validateSyncSetup(state, showToast)) {
      setState((state) => ({ ...state, error: "Please fill out all required fields." }));
      return;
    }

    const convertFieldMappings = (fieldMapping: FieldMappingState): FieldMappingInput[] => {
      if (fieldMapping.expandedJson) {
        return fieldMapping.jsonFields.flatMap((jsonMapping) => {
          if (jsonMapping === undefined) {
            consumeError(new Error("JSON mapping is undefined"));
            return [];
          }
          return [
            {
              source_field_name: jsonMapping.name,
              source_field_type: jsonMapping.type,
              destination_field_id: fieldMapping.destinationField.id,
              is_json_field: true,
            },
          ];
        });
      } else {
        if (fieldMapping.sourceField === undefined) {
          consumeError(new Error("Field mapping source field is undefined"));
          return [];
        }

        return [
          {
            source_field_name: fieldMapping.sourceField.name,
            source_field_type: fieldMapping.sourceField.type,
            destination_field_id: fieldMapping.destinationField.id,
            is_json_field: false,
          },
        ];
      }
    };

    if (state.fieldMappings === undefined) {
      consumeError(new Error("Field mappings are undefined"));
      return;
    }

    const fieldMappings: FieldMappingInput[] = state.fieldMappings.flatMap(convertFieldMappings);
    const payload: LinkCreateSyncRequest = {
      display_name: state.displayName!,
      source_id: state.source!.id,
      object_id: state.object!.id,
      namespace: state.namespace,
      table_name: state.tableName,
      frequency: state.frequency!,
      frequency_units: state.frequencyUnits!,
      field_mappings: fieldMappings,
    };

    try {
      await sendLinkTokenRequest(LinkCreateSync, linkToken, payload);
      // Tell SWRs to refetch syncs
      mutate({ LinkGetSyncs });
      showToast("success", "Success! Your sync has been created.", 4000);

      // Reset state so a new sync can be created
      resetState(setState);
      navigate("/");
    } catch (e) {
      if (e instanceof HttpError) {
        const errorMessage = e.message;
        setState((state) => ({ ...state, error: errorMessage }));
      }
      consumeError(e);
    }
  };
};
