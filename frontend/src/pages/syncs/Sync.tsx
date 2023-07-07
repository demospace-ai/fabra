import { ChevronDownIcon, ChevronUpIcon } from "@heroicons/react/24/outline";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { BackButton } from "src/components/button/Button";
import { LongRightArrow } from "src/components/icons/Icons";
import { DotsLoading, Loading } from "src/components/loading/Loading";
import { useShowToast } from "src/components/notifications/Notifications";
import { EmptyTable } from "src/components/table/Table";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { sendRequest } from "src/rpc/ajax";
import { FieldMapping, FieldType, RunSync, Sync as SyncConfig, SyncRunStatus } from "src/rpc/api";
import { useSync } from "src/rpc/data";
import { useMutation } from "src/utils/queryHelpers";
import { mergeClasses } from "src/utils/twmerge";

const tableHeaderStyle =
  "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-px-4 sm:tw-pr-6 lg:tw-pr-8 tw-text-left tw-whitespace-nowrap";
const tableCellStyle =
  "tw-whitespace-nowrap tw-left tw-overflow-hidden tw-py-4 tw-pl-4 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export const Sync: React.FC = () => {
  const navigate = useNavigate();
  const showToast = useShowToast();
  const { syncID } = useParams<{ syncID: string }>();
  const { sync, mutate } = useSync(Number(syncID));
  const [showDetails, setShowDetails] = useState<boolean>(false);

  const runSyncMutation = useMutation(
    async () => {
      await sendRequest(RunSync, { syncID });
    },
    {
      onSuccess: () => {
        showToast("success", "Success! Sync will start shortly.", 2000);
        mutate();
        setTimeout(() => {
          runSyncMutation.reset();
        }, 2000);
      },
      onError: () => {
        showToast("error", "Failed to run sync.", 2000);
      },
    },
  );

  if (!sync) {
    return <Loading />;
  }

  const syncRuns = sync.sync_runs ? sync.sync_runs : [];

  return (
    <div className="tw-pt-5 tw-pb-24 tw-px-10 tw-h-full tw-w-full tw-overflow-scroll">
      <BackButton onClick={() => navigate("/syncs")} />
      <div className="tw-flex tw-w-full tw-mb-2 tw-mt-4">
        <div className="tw-flex tw-flex-row tw-w-full tw-items-center tw-justify-between">
          <div className="tw-font-bold tw-text-2xl">{sync.sync.display_name}</div>
          <div className="tw-flex">
            <button
              disabled={runSyncMutation.isLoading}
              className="tw-ml-auto tw-px-4 tw-py-1 tw-rounded-md tw-font-medium tw-text-base tw-bg-blue-600 hover:tw-bg-blue-500 tw-text-white disabled:tw-bg-gray-500"
              onClick={() => runSyncMutation.mutate()}
            >
              {runSyncMutation.isLoading ? <Loading light className="tw-w-4 tw-h-4" /> : "Run sync"}
            </button>
          </div>
        </div>
      </div>
      <div
        className="tw-flex tw-w-fit tw-items-center tw-mb-5 tw-cursor-pointer tw-text-blue-500 tw-select-none"
        onClick={() => setShowDetails(!showDetails)}
      >
        {showDetails ? (
          <>
            Collapse details <ChevronUpIcon className="tw-h-3" />
          </>
        ) : (
          <>
            Expand details <ChevronDownIcon className="tw-h-3" />
          </>
        )}
      </div>
      {showDetails && <SyncDetails sync={sync.sync} mappings={sync.field_mappings} />}
      <div className="tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-bg-white tw-rounded-lg tw-overflow-auto tw-shadow-md tw-w-full">
        {sync ? (
          <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
            <thead className="tw-bg-slate-100 tw-text-slate-900">
              <tr>
                <th scope="col" className={tableHeaderStyle}>
                  Status
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Started At
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Rows Synced
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Error
                </th>
              </tr>
            </thead>
            <tbody className="tw-divide-y tw-divide-slate-200 tw-bg-white">
              {syncRuns.length > 0 ? (
                syncRuns.map((syncRun, index) => (
                  <tr key={index} className="tw-cursor-pointer hover:tw-bg-slate-50" onClick={() => {}}>
                    <td className={tableCellStyle}>
                      <div
                        className={mergeClasses(
                          "tw-flex tw-justify-center tw-items-center tw-py-1 tw-px-2 tw-rounded tw-text-center tw-w-[110px] tw-text-xs tw-font-medium",
                          getStatusStyle(syncRun.status),
                        )}
                      >
                        {syncRun.status.toUpperCase()}{" "}
                        {syncRun.status === SyncRunStatus.Running && <DotsLoading className="tw-ml-1.5" />}
                      </div>
                    </td>
                    <td className={tableCellStyle}>
                      <div>
                        <div className="tw-font-medium tw-mb-0.5">{syncRun.started_at}</div>
                        {syncRun.duration && (
                          <div className="tw-text-xs tw-text-slate-500">Duration: {syncRun.duration}</div>
                        )}
                      </div>
                    </td>
                    <td className={tableCellStyle}>{syncRun.rows_written}</td>
                    <td className={tableCellStyle}>
                      <Tooltip
                        label={<div className="tw-m-2 tw-cursor-text tw-font-mono">{syncRun.error}</div>}
                        maxWidth={600}
                        interactive
                      >
                        <div className="tw-overflow-hidden tw-text-ellipsis tw-max-w-[450px]">{syncRun.error}</div>
                      </Tooltip>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td className={tableCellStyle}>No sync runs yet!</td>
                </tr>
              )}
            </tbody>
          </table>
        ) : (
          <EmptyTable />
        )}
      </div>
    </div>
  );
};

const getStatusStyle = (status: SyncRunStatus): string => {
  switch (status) {
    case SyncRunStatus.Running:
      return "tw-bg-sky-100 tw-border tw-border-solid tw-border-sky-500 tw-text-sky-600";
    case SyncRunStatus.Completed:
      return "tw-bg-green-100 tw-border tw-border-solid tw-border-green-500 tw-text-green-600";
    case SyncRunStatus.Failed:
      return "tw-bg-red-100 tw-border tw-border-solid tw-border-red-500 tw-text-red-500";
    default:
      return "tw-bg-gray-100 tw-border tw-border-solid tw-border-gray-500 tw-text-gray-500";
  }
};

export const SyncDetails: React.FC<{ sync: SyncConfig; mappings: FieldMapping[] }> = ({ sync, mappings }) => {
  return (
    <>
      <div className="tw-flex tw-flex-col tw-w-fit tw-flex-wrap tw-items-start tw-px-3 tw-pt-1 tw-pb-2 tw-mb-5 tw-bg-white tw-border tw-border-slate-200 tw-rounded-md">
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Source ID: </span>
          {sync.source_id}
        </div>
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Namespace: </span>
          {sync.namespace}
        </div>
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Table: </span>
          {sync.table_name}
        </div>
      </div>
      <div className="tw-font-semibold tw-text-base tw-mb-2">Field Mappings</div>
      <div className="tw-border tw-border-slate-200 tw-bg-white tw-w-fit tw-rounded-lg tw-divide-y tw-mb-8">
        {mappings.map((mapping) => (
          <div key={mapping.source_field_name} className="tw-flex tw-flex-row tw-p-3 tw-items-center">
            <span className="tw-w-32 tw-mr-4 tw-max-w-[128px] tw-font-medium tw-overflow-clip tw-text-ellipsis">
              {mapping.source_field_name}
            </span>
            <LongRightArrow className="tw-fill-slate-300 tw-h-2" />
            <MappedField name={mapping.destination_field_name} type={mapping.destination_field_type} />
          </div>
        ))}
      </div>
    </>
  );
};

const MappedField: React.FC<{
  name: string;
  type: FieldType;
}> = ({ name, type }) => {
  return (
    <div className="tw-ml-12 tw-mr-2">
      <div className="tw-flex tw-h-fit">
        <div className="tw-h-fit tw-border tw-border-slate-200 tw-rounded-md tw-px-2 tw-box-border tw-bg-slate-100 tw-flex tw-flex-row tw-items-center tw-text-slate-700 tw-font-mono tw-select-none">
          <div className="tw-w-fit">{name}</div>
        </div>
        <div className="tw-h-fit tw-ml-3 tw-lowercase tw-select-none tw-font-mono tw-text-slate-500 tw-flex">
          {type}
        </div>
      </div>
    </div>
  );
};
