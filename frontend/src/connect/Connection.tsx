import React from "react";
import { Button } from "src/components/button/Button";
import { InfoIcon } from "src/components/icons/Icons";
import { ConnectionImage } from "src/components/images/Connections";
import { Input, ValidatedInput } from "src/components/input/Input";
import { Loading } from "src/components/loading/Loading";
import { GoogleLocationSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { FabraDisplayOptions } from "src/connect/ConnectApp";
import { NewSourceState, SetupSyncProps, SyncSetupStep, validateConnectionSetup } from "src/connect/state";
import { sendLinkTokenRequest } from "src/rpc/ajax";
import { ConnectionType, getConnectionType, TestDataConnection, TestDataConnectionRequest } from "src/rpc/api";
import { consumeError, forceError } from "src/utils/errors";
import { useMutation } from "src/utils/queryHelpers";
import { mergeClasses } from "src/utils/twmerge";

export const NewSourceConfiguration: React.FC<SetupSyncProps & FabraDisplayOptions> = (props) => {
  const state = props.state.newSourceState;

  // setState computes the NewSourceState using the provided function, then passes that new state to the parent setState
  const setNewSourceState = (getNewSourceState: (newSourceState: NewSourceState) => NewSourceState) => {
    props.setState((state) => {
      const newSourceState = getNewSourceState(state.newSourceState);
      return { ...state, newSourceState: newSourceState };
    });
  };

  const connectionType = props.state.connectionType;
  if (!connectionType) {
    // TODO: handle error, this should never happen
    return <></>;
  }

  if (props.state.newSourceState.sourceCreated) {
    props.setState((state) => ({ ...state, step: SyncSetupStep.ExistingSources }));
    return <Loading />;
  }

  let inputs: React.ReactElement;
  switch (connectionType) {
    case ConnectionType.Snowflake:
      inputs = <SnowflakeInputs state={props.state.newSourceState} setState={setNewSourceState} />;
      break;
    case ConnectionType.BigQuery:
      inputs = <BigQueryInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.Redshift:
      inputs = <RedshiftInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.MongoDb:
      inputs = <MongoDbInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.Synapse:
      inputs = <SynapseInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.Postgres:
      inputs = <PostgresInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.MySQL:
      inputs = <MySqlInputs state={state} setState={setNewSourceState} />;
      break;
    case ConnectionType.DynamoDb:
    case ConnectionType.Webhook:
      inputs = <>Unexpected</>;
      break;
  }

  return (
    <div className="tw-pl-20 tw-pr-[72px] tw-flex tw-flex-col tw-w-full">
      <div className="tw-flex tw-mb-2 tw-text-2xl tw-font-semibold tw-text-slate-900">
        <ConnectionImage connectionType={connectionType} className="tw-h-8 tw-mr-1.5" />
        Connect to {getConnectionType(connectionType)}
      </div>
      <div className="tw-flex tw-flex-row">
        <div className="tw-pb-16 tw-w-[500px] tw-mr-10">
          <div className="tw-mb-4 tw-text-slate-600">Provide the settings and credentials for your data source.</div>
          {inputs}
          <TestConnectionButton
            linkToken={props.linkToken}
            state={state}
            setState={setNewSourceState}
            connectionType={connectionType}
          />
          {state.error && (
            <div className="tw-mt-4 tw-text-red-700 tw-p-2 tw-text-center tw-bg-red-50 tw-border tw-border-red-600 tw-rounded">
              {state.error}
            </div>
          )}
        </div>
        <div className="tw-w-80 tw-ml-auto tw-text-xs tw-leading-5 tw-border-l tw-border-slate-200 tw-h-fit tw-py-2 tw-pl-8 tw-mr-10">
          <div className="">
            <div className="tw-text-[13px] tw-mb-1 tw-font-medium">Read our docs</div>
            Not sure where to start? Check out{" "}
            <a
              href={props.docsLink ? props.docsLink : "https://docs.fabra.io"}
              target="_blank"
              rel="noreferrer"
              className="tw-text-blue-500"
            >
              the docs
            </a>{" "}
            for step-by-step instructions.
          </div>
          <div className="tw-my-5 tw-py-5 tw-border-y tw-border-slate-200">
            <div className="tw-text-[13px] tw-mb-1 tw-font-medium">Allowed IPs</div>
            If your warehouse is behind a firewall/private network, please add the following static IP address:
            <ul className="tw-mt-1">
              <li>• 34.145.25.122</li>
            </ul>
          </div>
          <div>
            <div className="tw-text-[13px] tw-mb-1 tw-font-medium">Contact support</div>
            We"re here to help!{" "}
            <a
              href={props.supportEmail ? "mailto:" + props.supportEmail : "mailto:help@fabra.io"}
              className="tw-text-blue-500"
            >
              Reach out
            </a>{" "}
            if you feel stuck or have any questions.
          </div>
        </div>
      </div>
    </div>
  );
};

const TestConnectionButton: React.FC<{
  linkToken: string;
  state: NewSourceState;
  setState: React.Dispatch<(state: NewSourceState) => NewSourceState>;
  connectionType: ConnectionType;
}> = ({ linkToken, state, setState, connectionType }) => {
  const onClick = () => {
    testConnectionMutation.reset();
    if (!validateConnectionSetup(connectionType, state)) {
      setState((state) => {
        return { ...state, error: "Must fill out all required fields" };
      });
      return;
    } else {
      setState((state) => {
        return { ...state, error: undefined };
      });
    }
    testConnectionMutation.mutate();
  };

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
        case ConnectionType.MongoDb:
          payload.mongodb_config = state.mongodbConfig;
          break;
        case ConnectionType.Redshift:
          payload.redshift_config = state.redshiftConfig;
          break;
        case ConnectionType.Synapse:
          payload.synapse_config = state.synapseConfig;
          break;
        case ConnectionType.Postgres:
          payload.postgres_config = state.postgresConfig;
          break;
        case ConnectionType.MySQL:
          payload.mysql_config = state.mysqlConfig;
          break;
        case ConnectionType.DynamoDb:
          // TODO: throw error
          consumeError(new Error("DynamoDB is not supported as a source yet."));
          return;
        case ConnectionType.Webhook:
          // TODO: throw error
          consumeError(new Error("Webhook is not supported as a source."));
          return;
      }

      return await sendLinkTokenRequest(TestDataConnection, linkToken, payload);
    },
    {
      onError: (err) => {
        const error = forceError(err);
        setState((state) => {
          return {
            ...state,
            error: error?.message ?? "Failed",
          };
        });
      },
    },
  );

  const testColor = testConnectionMutation.isSuccess ? "tw-bg-green-700 hover:tw-bg-green-800" : null;
  return (
    <>
      <Button className={mergeClasses("tw-mt-8 tw-border-slate-200 tw-w-48 tw-h-10", testColor)} onClick={onClick}>
        {testConnectionMutation.isLoading ? <Loading /> : "Test"}
      </Button>
      {testConnectionMutation.isSuccess && (
        <div className="tw-mt-4 tw-w-48 tw-text-green-700 tw-p-2 tw-text-center tw-bg-green-50 tw-border tw-border-green-600 tw-rounded">
          Success!
        </div>
      )}
    </>
  );
};

type ConnectionConfigurationProps = {
  state: NewSourceState;
  setState: React.Dispatch<(state: NewSourceState) => NewSourceState>;
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.snowflakeConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, snowflakeConfig: { ...state.snowflakeConfig, username: value } }));
        }}
        placeholder="Username"
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
          props.setState((state) => ({ ...state, snowflakeConfig: { ...state.snowflakeConfig, password: value } }));
        }}
        placeholder="Password"
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
          props.setState((state) => ({
            ...state,
            snowflakeConfig: { ...state.snowflakeConfig, database_name: value },
          }));
        }}
        placeholder="Database Name"
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
          props.setState((state) => ({
            ...state,
            snowflakeConfig: { ...state.snowflakeConfig, warehouse_name: value },
          }));
        }}
        placeholder="Warehouse Name"
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
          props.setState((state) => ({ ...state, snowflakeConfig: { ...state.snowflakeConfig, role: value } }));
        }}
        placeholder="Role"
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
          props.setState((state) => ({ ...state, snowflakeConfig: { ...state.snowflakeConfig, host: value } }));
        }}
        placeholder="Host"
      />
    </>
  );
};

const RedshiftInputs: React.FC<ConnectionConfigurationProps> = (props) => {
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.redshiftConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, redshiftConfig: { ...state.redshiftConfig, username: value } }));
        }}
        placeholder="Username"
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
        value={state.redshiftConfig.password}
        setValue={(value) => {
          props.setState((state) => ({ ...state, redshiftConfig: { ...state.redshiftConfig, password: value } }));
        }}
        placeholder="Password"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Database Name</span>
        <Tooltip placement="right" label="The Redshift database to sync from.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="databaseName"
        value={state.redshiftConfig.database_name}
        setValue={(value) => {
          props.setState((state) => ({ ...state, redshiftConfig: { ...state.redshiftConfig, database_name: value } }));
        }}
        placeholder="Database Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Endpoint</span>
        <Tooltip
          placement="right"
          label={
            <div className="tw-m-2">
              <div>
                This is the URL for your Redshift data warehouse. For Redshift clusters, it can be found on the specific
                cluster page under "General Information" and should look like:
              </div>
              <div className="tw-mt-2 tw-w-full tw-bg-slate-900 tw-rounded-md tw-p-2">
                your-cluster.abc123.us-west-2.redshift.amazonaws.com
              </div>
              <div className="tw-mt-3">
                For Serverless Redshift,{" "}
                <a
                  className="tw-text-blue-400"
                  target="_blank"
                  rel="noreferrer"
                  href="https://docs.aws.amazon.com/redshift/latest/mgmt/serverless-connecting.html"
                >
                  visit the Redshift docs.
                </a>{" "}
                The following is the expected format for Serverless Redshift:
              </div>
              <div className="tw-mt-2 tw-w-full tw-bg-slate-900 tw-rounded-md tw-p-2">
                <span className="tw-italic">workgroup-name</span>.<span className="tw-italic">account-number</span>.
                <span className="tw-italic">aws-region</span>.redshift-serverless.amazonaws.com
              </div>
            </div>
          }
          interactive
          maxWidth={640}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="endpoint"
        value={state.redshiftConfig.endpoint}
        setValue={(value) => {
          props.setState((state) => ({ ...state, redshiftConfig: { ...state.redshiftConfig, endpoint: value } }));
        }}
        placeholder="Endpoint"
      />
    </>
  );
};

const SynapseInputs: React.FC<ConnectionConfigurationProps> = (props) => {
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.synapseConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, synapseConfig: { ...state.synapseConfig, username: value } }));
        }}
        placeholder="Username"
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
        value={state.synapseConfig.password}
        setValue={(value) => {
          props.setState((state) => ({ ...state, synapseConfig: { ...state.synapseConfig, password: value } }));
        }}
        placeholder="Password"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Database Name</span>
        <Tooltip placement="right" label="The Synapse database to sync from.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="databaseName"
        value={state.synapseConfig.database_name}
        setValue={(value) => {
          props.setState((state) => ({ ...state, synapseConfig: { ...state.synapseConfig, database_name: value } }));
        }}
        placeholder="Database Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Endpoint</span>
        <Tooltip
          placement="right"
          label={
            <div className="tw-m-2">
              <div>
                This is the URL for your Synapse data warehouse. It can be found on your Synapse Workspace Overview page
                under "Essentials". For Dedicated SQL pools it should look like:
              </div>
              <div className="tw-mt-2 tw-w-full tw-bg-slate-900 tw-rounded-md tw-p-2">abc123.sql.azuresynapse.net</div>
              <div className="tw-mt-3">For Serverless Synapse, it should look like:</div>
              <div className="tw-mt-2 tw-w-full tw-bg-slate-900 tw-rounded-md tw-p-2">
                abc123-ondemand.sql.azuresynapse.net
              </div>
            </div>
          }
          interactive
          maxWidth={640}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="endpoint"
        value={state.synapseConfig.endpoint}
        setValue={(value) => {
          props.setState((state) => ({ ...state, synapseConfig: { ...state.synapseConfig, endpoint: value } }));
        }}
        placeholder="Endpoint"
      />
    </>
  );
};

const PostgresInputs: React.FC<ConnectionConfigurationProps> = (props) => {
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.postgresConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, postgresConfig: { ...state.postgresConfig, username: value } }));
        }}
        placeholder="Username"
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
        value={state.postgresConfig.password}
        setValue={(value) => {
          props.setState((state) => ({ ...state, postgresConfig: { ...state.postgresConfig, password: value } }));
        }}
        placeholder="Password"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Database Name</span>
        <Tooltip placement="right" label="The Postgres database to sync from.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="databaseName"
        value={state.postgresConfig.database_name}
        setValue={(value) => {
          props.setState((state) => ({ ...state, postgresConfig: { ...state.postgresConfig, database_name: value } }));
        }}
        placeholder="Database Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Endpoint</span>
        <Tooltip
          placement="right"
          label="This is the endpoint for your Postgres database. It must be in the format <host>:<port>, i.e. 127.0.0.1:5432"
          interactive
          maxWidth={640}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="host"
        value={state.postgresConfig.endpoint}
        setValue={(value) => {
          props.setState((state) => ({ ...state, postgresConfig: { ...state.postgresConfig, endpoint: value } }));
        }}
        placeholder="Endpoint"
      />
    </>
  );
};

const MySqlInputs: React.FC<ConnectionConfigurationProps> = (props) => {
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.mysqlConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mysqlConfig: { ...state.mysqlConfig, username: value } }));
        }}
        placeholder="Username"
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
        value={state.mysqlConfig.password}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mysqlConfig: { ...state.mysqlConfig, password: value } }));
        }}
        placeholder="Password"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Database Name</span>
        <Tooltip placement="right" label="The MySQL database to sync from.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="databaseName"
        value={state.mysqlConfig.database_name}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mysqlConfig: { ...state.mysqlConfig, database_name: value } }));
        }}
        placeholder="Database Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Endpoint</span>
        <Tooltip
          placement="right"
          label="This is the endpoint for your MySQL database. It must be in the format <host>:<port>, i.e. 127.0.0.1:5432"
          interactive
          maxWidth={640}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="host"
        value={state.mysqlConfig.endpoint}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mysqlConfig: { ...state.mysqlConfig, endpoint: value } }));
        }}
        placeholder="Endpoint"
      />
    </>
  );
};

const MongoDbInputs: React.FC<ConnectionConfigurationProps> = (props) => {
  const state = props.state;
  return (
    <>
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Display Name</span>
        <Tooltip placement="right" label="Pick a name to help you identify this source in the future.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="displayName"
        value={state.displayName}
        setValue={(value) => {
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="MongoDB Source"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Username</span>
        <Tooltip placement="right" label="We recommend you create a dedicated user for syncing.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="username"
        value={state.mongodbConfig.username}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mongodbConfig: { ...state.mongodbConfig, username: value } }));
        }}
        placeholder="sync_user"
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
        value={state.mongodbConfig.password}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mongodbConfig: { ...state.mongodbConfig, password: value } }));
        }}
        placeholder="VerySecurePassword1"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Host</span>
        <Tooltip placement="right" label="The hostname of your MongoDB instance.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="host"
        value={state.mongodbConfig.host}
        setValue={(value) => {
          props.setState((state) => ({ ...state, mongodbConfig: { ...state.mongodbConfig, host: value } }));
        }}
        placeholder="mymongo.abc123.mongodb.net"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Connection Options (optional)</span>
        <Tooltip
          placement="right"
          interactive
          label={
            <div>
              Any additional options to apply to the MongoDB connection. See more{" "}
              <a
                className="tw-text-blue-400"
                target="_blank"
                rel="noreferrer"
                href="https://www.mongodb.com/docs/drivers/node/current/fundamentals/connection/connection-options/"
              >
                here
              </a>
              .
            </div>
          }
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <Input
        id="connectionOptions"
        value={state.mongodbConfig.connection_options}
        setValue={(value) => {
          props.setState((state) => ({
            ...state,
            mongodbConfig: { ...state.mongodbConfig, connection_options: value },
          }));
        }}
        placeholder="connectTimeoutMS=30&timeoutMS=1000"
      />
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
          props.setState((state) => ({ ...state, displayName: value }));
        }}
        placeholder="Display Name"
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
          props.setState((state) => ({ ...state, bigqueryConfig: { ...state.bigqueryConfig, location: value } }));
        }}
        placeholder="Location"
        className="tw-mt-0 tw-w-full"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Service Account Key</span>
        <Tooltip
          placement="right"
          label="This can be obtained in the Google Cloud web console by navigating to the IAM page and clicking on Service Accounts in the left sidebar. Then, find your service account in the list, go to its Keys tab, and click Add Key. Finally, click on Create new key and choose JSON."
          interactive
          maxWidth={500}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        className="tw-h-24 tw-min-h-[40px] tw-max-h-80"
        id="credentials"
        value={state.bigqueryConfig.credentials}
        setValue={(value) => {
          props.setState((state) => ({ ...state, bigqueryConfig: { ...state.bigqueryConfig, credentials: value } }));
        }}
        placeholder="Credentials (paste JSON here)"
        textarea={true}
      />
    </>
  );
};
