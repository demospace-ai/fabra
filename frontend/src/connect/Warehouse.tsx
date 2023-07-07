import React from "react";
import { ConnectionImage } from "src/components/images/Connections";
import { INITIAL_SOURCE_STATE, SetupSyncProps, SyncSetupStep } from "src/connect/state";
import { ConnectionType } from "src/rpc/api";

export const WarehouseSelector: React.FC<SetupSyncProps> = (props) => {
  const connectionButton =
    "tw-flex tw-flex-row tw-justify-center tw-items-center tw-py-5 tw-font-medium tw-w-64 tw-rounded-md tw-cursor-pointer tw-bg-white tw-text-slate-800 tw-border tw-border-slate-300 hover:tw-bg-slate-100 tw-shadow tw-select-none";
  const onClick = (connectionType: ConnectionType) => {
    // Reset new source state to initial state when user selects a new connection type
    props.setState({
      ...props.state,
      connectionType: connectionType,
      step: SyncSetupStep.ConnectionDetails,
      skippedSourceSetup: false,
      newSourceState: INITIAL_SOURCE_STATE,
    });
  };

  return (
    <div className="tw-w-full tw-px-20">
      <div className="tw-text-left tw-mb-2 tw-text-2xl tw-font-semibold tw-text-slate-900">Add a new data source</div>
      <div className="tw-text-left tw-mb-10 tw-text-slate-600">
        Choose the data warehouse, database, or data lake to connect.
      </div>
      <div className="tw-flex tw-flex-row tw-gap-5 tw-flex-wrap tw-justify-start">
        <button className={connectionButton} onClick={() => onClick(ConnectionType.Snowflake)}>
          <ConnectionImage connectionType={ConnectionType.Snowflake} className="tw-h-6 tw-mr-1.5" />
          Snowflake
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.BigQuery)}>
          <ConnectionImage connectionType={ConnectionType.BigQuery} className="tw-h-6 tw-mr-2" />
          BigQuery
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.Redshift)}>
          <ConnectionImage connectionType={ConnectionType.Redshift} className="tw-h-6 tw-mr-2" />
          Redshift
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.MongoDb)}>
          <ConnectionImage connectionType={ConnectionType.MongoDb} className="tw-h-6 tw-mr-1" />
          MongoDB
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.Synapse)}>
          <ConnectionImage connectionType={ConnectionType.Synapse} className="tw-h-6 tw-mr-1.5" />
          Azure Synapse
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.Postgres)}>
          <ConnectionImage connectionType={ConnectionType.Postgres} className="tw-h-6 tw-mr-2" />
          Postgres
        </button>
        <button className={connectionButton} onClick={() => onClick(ConnectionType.MySQL)}>
          <ConnectionImage connectionType={ConnectionType.MySQL} className="tw-h-6 tw-mr-2" />
          MySQL
        </button>
      </div>
    </div>
  );
};
