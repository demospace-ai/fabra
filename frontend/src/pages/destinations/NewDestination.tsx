import React, { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { BackButton, Button, FormButton } from "src/components/button/Button";
import { ErrorDisplay } from "src/components/error/Error";
import { InfoIcon } from "src/components/icons/Icons";
import { ConnectionImage } from "src/components/images/Connections";
import { ValidatedInput } from "src/components/input/Input";
import { Loading } from "src/components/loading/Loading";
import { useShowToast } from "src/components/notifications/Notifications";
import { GoogleLocationSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { DynamoDbInputs } from "src/pages/destinations/DynamoDbInputs";
import { NewDestinationState } from "src/pages/destinations/helpers";
import { sendRequest } from "src/rpc/ajax";
import {
  ConnectionType,
  CreateDestination,
  CreateDestinationRequest,
  CreateDynamoDbConfigSchema,
  DynamoDbConfigSchema,
  getConnectionType,
  GetDestinations,
  TestDataConnection,
  TestDataConnectionRequest,
} from "src/rpc/api";
import { forceError } from "src/utils/errors";
import { useMutation } from "src/utils/queryHelpers";
import { mutate } from "swr";

export const NewDestination: React.FC = () => {
  const [connectionType, setConnectionType] = useState<ConnectionType | null>(null);
  const navigate = useNavigate();
  const onBack = () => {
    if (connectionType) {
      setConnectionType(null);
    } else {
      navigate("/destinations");
    }
  };

  return (
    <div className="tw-py-5 tw-px-10 tw-overflow-scroll">
      <BackButton className="tw-mt-3" onClick={onBack} />
      <div className="tw-flex tw-flex-col tw-mt-8 tw-mb-24 tw-py-12 tw-px-10 tw-bg-white tw-rounded-lg tw-shadow-md tw-items-center">
        <div className="tw-text-center tw-mb-5 tw-font-bold tw-text-lg">New Destination</div>
        {connectionType ? (
          <NewDestinationConfiguration connectionType={connectionType} setConnectionType={setConnectionType} />
        ) : (
          <ConnectionTypeSelector setConnectionType={setConnectionType} />
        )}
      </div>
    </div>
  );
};

type NewConnectionConfigurationProps = {
  connectionType: ConnectionType;
  setConnectionType: (connectionType: ConnectionType | null) => void;
};

// Values must be empty strings otherwise the input will be uncontrolled
const INITIAL_DESTINATION_STATE: NewDestinationState = {
  displayName: "",
  staging_bucket: "",
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
  webhookConfig: {
    url: "",
    headers: [],
  },
  postgresConfig: {
    username: "",
    password: "",
    database_name: "",
    endpoint: "",
  },
  dynamoDbConfig: {
    secretKey: "",
    accessKey: "",
    region: undefined,
  },
  error: undefined,
};

const validateAll = (
  connectionType: ConnectionType,
  state: NewDestinationState,
  setState: React.Dispatch<React.SetStateAction<NewDestinationState>>,
): boolean => {
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
      return (
        state.displayName.length > 0 &&
        state.staging_bucket.length > 0 &&
        state.bigqueryConfig.location !== undefined &&
        state.bigqueryConfig.credentials.length > 0
      );
    case ConnectionType.Webhook:
      if (state.displayName.length <= 0) {
        setState((state) => ({ ...state, error: "Must set a display name." }));
        return false;
      }

      if (state.webhookConfig.url.length <= 0) {
        setState((state) => ({ ...state, error: "Must specify Webhook URL" }));
        return false;
      }

      if (!state.webhookConfig.url.includes("https")) {
        setState((state) => ({ ...state, error: "Webhook URL must use HTTPS" }));
        return false;
      }

      return true;
    case ConnectionType.DynamoDb: {
      const result = DynamoDbConfigSchema.safeParse(state.dynamoDbConfig);
      return result.success;
    }
    case ConnectionType.Redshift:
    case ConnectionType.Synapse:
    case ConnectionType.MongoDb:
    case ConnectionType.Postgres:
    case ConnectionType.MySQL:
      // TODO: throw error
      return false;
  }
};

const NewDestinationConfiguration: React.FC<NewConnectionConfigurationProps> = (props) => {
  const [state, setState] = useState<NewDestinationState>(INITIAL_DESTINATION_STATE);
  const navigate = useNavigate();
  const showToast = useShowToast();

  const newDestinationMutation = useMutation(
    async () => {
      const payload: CreateDestinationRequest = {
        display_name: state.displayName,
        connection_type: props.connectionType,
      };

      if (state.staging_bucket) {
        payload.staging_bucket = state.staging_bucket;
      }

      switch (props.connectionType) {
        case ConnectionType.BigQuery:
          payload.bigquery_config = {
            location: state.bigqueryConfig.location!.code,
            credentials: state.bigqueryConfig.credentials,
          };
          break;
        case ConnectionType.Snowflake:
          payload.snowflake_config = state.snowflakeConfig;
          break;
        case ConnectionType.Webhook:
          payload.webhook_config = state.webhookConfig;
          break;
        case ConnectionType.DynamoDb: {
          const config = DynamoDbConfigSchema.parse(state.dynamoDbConfig);
          payload.dynamodb_config = CreateDynamoDbConfigSchema.parse({
            access_key: config.accessKey,
            secret_key: config.secretKey,
            region: config.region.code,
          });
          break;
        }
        case ConnectionType.Redshift:
        case ConnectionType.MongoDb:
        case ConnectionType.Synapse:
        case ConnectionType.Postgres:
        case ConnectionType.MySQL:
          // TODO: throw error
          break;
      }

      return await sendRequest(CreateDestination, payload);
    },
    {
      onSuccess: (destination) => {
        mutate({ GetDestinations });
        showToast("success", "Created destination!", 5000);
        navigate(`/destinations/${destination.destination.id}`);
      },
      onError: (err) => {
        showToast("error", err?.message, 5000);
      },
    },
  );

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!validateAll(props.connectionType, state, setState)) {
      return;
    }

    newDestinationMutation.mutate();
  };

  let inputs: React.ReactElement;
  switch (props.connectionType) {
    case ConnectionType.Snowflake:
      inputs = <SnowflakeInputs state={state} setState={setState} />;
      break;
    case ConnectionType.BigQuery:
      inputs = <BigQueryInputs state={state} setState={setState} />;
      break;
    case ConnectionType.Webhook:
      inputs = <WebhookInputs state={state} setState={setState} />;
      break;
    case ConnectionType.DynamoDb:
      inputs = <DynamoDbInputs state={state} setState={setState} />;
      break;
    case ConnectionType.Redshift:
    case ConnectionType.MongoDb:
    case ConnectionType.Synapse:
    case ConnectionType.Postgres:
    case ConnectionType.MySQL:
      inputs = <></>;
      break; // TODO: throw error
  }

  return (
    <div>
      <div className="tw-flex tw-items-center tw-mb-8">
        <ConnectionImage connectionType={props.connectionType} className="tw-h-6 tw-mr-1.5" />
        <div className="tw-font-medium">Enter your {getConnectionType(props.connectionType)} configuration:</div>
      </div>
      <form onSubmit={handleSubmit}>
        {inputs}
        {state.error && (
          <div className="tw-mt-4 tw-text-red-700 tw-p-2 tw-text-center tw-bg-red-50 tw-border tw-border-red-600 tw-rounded">
            {state.error}
          </div>
        )}
        <div className="tw-flex tw-flex-row tw-justify-start tw-w-100 tw-gap-5 tw-mt-12">
          <TestConnectionButton state={state} setState={setState} connectionType={props.connectionType} />
          <FormButton className="tw-w-full tw-h-10">
            {newDestinationMutation.isLoading ? <Loading /> : "Save"}
          </FormButton>
          <></>
        </div>
      </form>
      {newDestinationMutation.error && (
        <div className="tw-mt-4 tw-text-red-700 tw-p-2 tw-text-center tw-bg-red-50 tw-border tw-border-red-600 tw-rounded">
          {newDestinationMutation.error?.message || "Could not save destination"}
        </div>
      )}
    </div>
  );
};

const TestConnectionButton: React.FC<ConnectionConfigurationProps & { connectionType: ConnectionType }> = ({
  state,
  setState,
  connectionType,
}) => {
  const testConnectionMutation = useMutation(
    async () => {
      const payload: TestDataConnectionRequest = {
        display_name: state.displayName,
        connection_type: connectionType,
      };

      switch (connectionType) {
        case ConnectionType.BigQuery:
          payload.bigquery_config = {
            location: state.bigqueryConfig.location!.code,
            credentials: state.bigqueryConfig.credentials,
          };
          break;
        case ConnectionType.Snowflake:
          payload.snowflake_config = state.snowflakeConfig;
          break;
        case ConnectionType.Webhook:
          payload.webhook_config = state.webhookConfig;
          break;
        case ConnectionType.DynamoDb: {
          const config = DynamoDbConfigSchema.parse(state.dynamoDbConfig);
          payload.dynamodb_config = CreateDynamoDbConfigSchema.parse({
            access_key: config.accessKey,
            secret_key: config.secretKey,
            region: config.region.code,
          });
          break;
        }
        case ConnectionType.Redshift:
        case ConnectionType.MongoDb:
        case ConnectionType.Synapse:
        case ConnectionType.Postgres:
        case ConnectionType.MySQL:
          // TODO: throw error
          break;
      }

      await sendRequest(TestDataConnection, payload);
    },
    {
      onSuccess: () => {
        setState((state) => {
          return { ...state, error: undefined };
        });
      },
      onError: (e) => {
        const err = forceError(e);
        setState((state) => ({
          ...state,
          error: err?.message,
        }));
      },
    },
  );

  const handleTestConnection = () => {
    if (!validateAll(connectionType, state, setState)) {
      setState((state) => ({
        ...state,
        error: "Please fill out all required fields.",
      }));
      return;
    }

    testConnectionMutation.mutate();
  };

  return (
    <div className="tw-w-full">
      <Button
        className="tw-bg-slate-200 tw-text-slate-900 hover:tw-bg-slate-300 tw-border-slate-200 tw-w-full tw-h-10"
        onClick={handleTestConnection}
      >
        {testConnectionMutation.isLoading ? <Loading /> : "Test"}
      </Button>
      {!testConnectionMutation.isLoading && testConnectionMutation.isSuccess && (
        <div className="tw-text-green-500 tw-mt-1">Connection success!</div>
      )}
      {!testConnectionMutation.isLoading && (
        <ErrorDisplay error={testConnectionMutation.error} className="tw-text-red-500 tw-mt-1" />
      )}
    </div>
  );
};

type ConnectionConfigurationProps = {
  state: NewDestinationState;
  setState: React.Dispatch<React.SetStateAction<NewDestinationState>>;
};

const SnowflakeInputs: React.FC<ConnectionConfigurationProps> = (props) => {
  const state = props.state;
  return (
    <>
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-2 tw-mb-1">
        <span>Display Name</span>
        <Tooltip placement="right" label="Pick a name to help you identify this source in the future.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="displayName"
        value={state.displayName}
        setValue={(value) => {
          props.setState({ ...state, displayName: value });
        }}
        placeholder="Display Name"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip
          placement="right"
          label="You can choose your personal username or create a dedicated user for syncing."
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.snowflakeConfig.username}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, username: value } });
        }}
        placeholder="Username"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Password</span>
        <Tooltip placement="right" label="Password for the user specified above.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="password"
        type="password"
        value={state.snowflakeConfig.password}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, password: value } });
        }}
        placeholder="Password"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Database Name</span>
        <Tooltip placement="right" label="The Snowflake database to sync from.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="databaseName"
        value={state.snowflakeConfig.database_name}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, database_name: value } });
        }}
        placeholder="Database Name"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Warehouse Name</span>
        <Tooltip placement="right" label="The warehouse that will be used to run syncs in Snowflake.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="warehouseName"
        value={state.snowflakeConfig.warehouse_name}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, warehouse_name: value } });
        }}
        placeholder="Warehouse Name"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Role</span>
        <Tooltip placement="right" label="The role that will be used to run syncs.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="role"
        value={state.snowflakeConfig.role}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, role: value } });
        }}
        placeholder="Role"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Host</span>
        <Tooltip
          placement="right"
          label={
            <div className="tw-m-2">
              <span>This is your Snowflake URL. Format may differ based on Snowflake account age. For details, </span>
              <a
                className="tw-text-blue-400"
                target="_blank"
                rel="noreferrer"
                href="https://docs.snowflake.com/en/user-guide/admin-account-identifier.html"
              >
                visit the Snowflake docs.
              </a>
              <div className="tw-mt-2">
                <span>Example:</span>
                <div className="tw-mt-2 tw-w-full tw-bg-slate-900 tw-rounded-md tw-p-2">
                  abc123.us-east1.gcp.snowflakecomputing.com
                </div>
              </div>
            </div>
          }
          interactive
          maxWidth={500}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="host"
        value={state.snowflakeConfig.host}
        setValue={(value) => {
          props.setState({ ...state, snowflakeConfig: { ...state.snowflakeConfig, host: value } });
        }}
        placeholder="Host"
        className="tw-w-100"
      />
    </>
  );
};

const WebhookInputs: React.FC<ConnectionConfigurationProps> = (props) => {
  const state = props.state;

  return (
    <>
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-2 tw-mb-1">
        <span>Display Name</span>
        <Tooltip placement="right" label="Pick a name to help you identify this source in the future.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="displayName"
        value={state.displayName}
        setValue={(value) => {
          props.setState({ ...state, displayName: value });
        }}
        placeholder="Display Name"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>URL</span>
        <Tooltip
          placement="right"
          label="The URL that Fabra will send your customer's data to during syncs. Must use HTTPS."
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <div className="tw-flex tw-items-center tw-w-100 tw-border tw-border-slate-300 hover:tw-border-primary-hover focus:tw-border-primary tw-rounded-md tw-overflow-clip">
        <ValidatedInput
          className="tw-w-100 tw-border-0 tw-bg-transparent tw-rounded-none"
          id="URL"
          value={state.webhookConfig.url}
          setValue={(value) => {
            props.setState({ ...state, webhookConfig: { ...state.webhookConfig, url: value } });
          }}
          placeholder="URL"
        />
      </div>
    </>
  );
};

const BigQueryInputs: React.FC<ConnectionConfigurationProps> = (props) => {
  const state = props.state;
  return (
    <>
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-2 tw-mb-1">
        <span>Display Name</span>
        <Tooltip placement="right" label="Pick a name to help you identify this source in the future.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="displayName"
        value={state.displayName}
        setValue={(value) => {
          props.setState({ ...state, displayName: value });
        }}
        placeholder="Display Name"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Location</span>
        <Tooltip placement="right" label="The geographic location of your BigQuery dataset(s).">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <GoogleLocationSelector
        id="location"
        location={state.bigqueryConfig.location}
        setLocation={(value) => {
          props.setState({ ...state, bigqueryConfig: { ...state.bigqueryConfig, location: value } });
        }}
        placeholder="Location"
        className="tw-w-100 tw-mt-0"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Staging Bucket</span>
        <Tooltip
          placement="right"
          interactive
          label={
            <div>
              The Google Cloud Storage (GCS) bucket Fabra will use for temporarily staging data during syncs. Learn more{" "}
              <a className="tw-text-blue-400" target="_blank" rel="noreferrer" href="https://docs.fabra.io/staging">
                here
              </a>
              .
            </div>
          }
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="staging-bucket"
        value={state.staging_bucket}
        setValue={(value) => {
          props.setState({ ...state, staging_bucket: value });
        }}
        placeholder="Staging Bucket"
        className="tw-w-100"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Service Account Key</span>
        <Tooltip
          placement="right"
          label="This can be obtained in the Google Cloud web console by navigating to the IAM page and clicking on Service Accounts in the left sidebar. Then, find your service account in the list, go to its Keys tab, and click Add Key. Finally, click on Create new key and choose JSON. Fabra requires the BigQuery Data Editor, BigQuery Job User, Storage Object Creator, and Storage Object Viewer roles."
          interactive
          maxWidth={500}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        className="tw-h-24 tw-min-h-[40px] tw-max-h-80 tw-w-100"
        id="credentials"
        value={state.bigqueryConfig.credentials}
        setValue={(value) => {
          props.setState({ ...state, bigqueryConfig: { ...state.bigqueryConfig, credentials: value } });
        }}
        placeholder="Credentials (paste JSON here)"
        textarea={true}
      />
    </>
  );
};

type ConnectionTypeSelectorProps = {
  setConnectionType: (connectionType: ConnectionType) => void;
};

const ConnectionTypeSelector: React.FC<ConnectionTypeSelectorProps> = (props) => {
  const connectionButton =
    "tw-flex tw-flex-row tw-justify-center tw-items-center tw-py-5 tw-font-bold tw-w-56 tw-rounded-md tw-cursor-pointer tw-bg-white tw-text-slate-800 tw-border tw-border-slate-300 hover:tw-bg-slate-100 tw-tracking-[1px] tw-shadow-md tw-select-none";
  return (
    <>
      <div className="tw-text-center tw-mb-8">Choose one of our supported destinations for syncs:</div>
      <div className="tw-flex tw-flex-row tw-gap-5 tw-flex-wrap tw-justify-center">
        <button className={connectionButton} onClick={() => props.setConnectionType(ConnectionType.Snowflake)}>
          <ConnectionImage connectionType={ConnectionType.Snowflake} className="tw-h-6 tw-mr-1.5" />
          Snowflake
        </button>
        <button className={connectionButton} onClick={() => props.setConnectionType(ConnectionType.BigQuery)}>
          <ConnectionImage connectionType={ConnectionType.BigQuery} className="tw-h-6 tw-mr-1.5" />
          BigQuery
        </button>
        <button className={connectionButton} onClick={() => props.setConnectionType(ConnectionType.Webhook)}>
          <ConnectionImage connectionType={ConnectionType.Webhook} className="tw-h-6 tw-mr-1.5" />
          Webhook
        </button>
        <button className={connectionButton} onClick={() => props.setConnectionType(ConnectionType.DynamoDb)}>
          <ConnectionImage connectionType={ConnectionType.DynamoDb} className="tw-h-5 tw-mr-2" />
          DynamoDB
        </button>
      </div>
    </>
  );
};
