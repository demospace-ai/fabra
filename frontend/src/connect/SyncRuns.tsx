import { ChevronDownIcon, ChevronUpIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button } from "src/components/button/Button";
import { DotsLoading, Loading } from "src/components/loading/Loading";
import { useConnectShowToast } from "src/components/notifications/Notifications";
import { EmptyTable } from "src/components/table/Table";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { SyncDetails } from "src/pages/syncs/Sync";
import { sendLinkTokenRequest } from "src/rpc/ajax";
import { LinkRunSync, SyncRunStatus } from "src/rpc/api";
import { useLinkSync } from "src/rpc/data";
import { useMutation } from "src/utils/queryHelpers";
import { mergeClasses } from "src/utils/twmerge";

export const SyncRuns: React.FC<{
  linkToken: string;
  close: (() => void) | undefined;
}> = ({ linkToken, close }) => {
  return (
    <div className="tw-w-full tw-h-full tw-flex tw-flex-col">
      <Header close={close} />
      <SyncRunsList linkToken={linkToken} />
      <Footer />
    </div>
  );
};

const tableHeaderStyle =
  "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-px-4 sm:tw-pr-6 lg:tw-pr-8 tw-text-left tw-whitespace-nowrap";
const tableCellStyle =
  "tw-whitespace-nowrap tw-left tw-overflow-hidden tw-py-4 tw-pl-4 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

const SyncRunsList: React.FC<{ linkToken: string }> = ({ linkToken }) => {
  const { syncID } = useParams<{ syncID: string }>();
  const { sync, mutate } = useLinkSync(Number(syncID), linkToken);
  const showToast = useConnectShowToast();
  const [showDetails, setShowDetails] = useState<boolean>(false);

  const runSyncMutation = useMutation(
    async () => {
      await sendLinkTokenRequest(LinkRunSync, linkToken, { syncID: syncID! });
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

  const syncRuns = sync.sync_runs ?? [];

  return (
    <div className="tw-mt-4 tw-pb-16 tw-px-20 tw-flex tw-flex-col tw-overflow-auto">
      <div className="tw-flex tw-w-full tw-mb-2">
        <div className="tw-flex tw-flex-row tw-w-full tw-items-center tw-font-bold tw-text-xl tw-justify-between">
          Sync Runs • {sync.sync.display_name}
          <div className="tw-flex">
            <Button
              className="tw-ml-auto tw-h-8 tw-w-32 tw-text-sm tw-whitespace-nowrap"
              disabled={runSyncMutation.isLoading}
              onClick={() => runSyncMutation.mutate()}
            >
              {!runSyncMutation.isLoading ? "Run Sync" : <Loading light className="tw-w-4 tw-h-4" />}
            </Button>
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
      <div className="tw-shadow tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-rounded-md tw-overflow-auto tw-overscroll-contain tw-w-full">
        {sync ? (
          <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
            <thead className="tw-sticky tw-top-0 tw-bg-slate-100">
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
                          "tw-flex tw-justify-center tw-items-center tw-py-1 tw-px-2 tw-rounded tw-text-center tw-w-[110px] tw-border tw-text-xs tw-font-medium",
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
                        <div className="tw-overflow-hidden tw-text-ellipsis tw-max-w-[240px]">{syncRun.error}</div>
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

const Header: React.FC<{ close: (() => void) | undefined }> = ({ close }) => {
  return (
    <div
      className={classNames(
        "tw-flex tw-flex-row tw-items-center tw-w-full",
        close ? "tw-h-20 tw-min-h-[80px]" : "tw-h-10 tw-min-h-[48px]",
      )}
    >
      {close && (
        <button
          className="tw-absolute tw-flex tw-items-center t tw-right-10 tw-border-none tw-cursor-pointer tw-p-0"
          onClick={close}
        >
          <svg className="tw-h-6 tw-fill-slate-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="none">
            <path d="M5.1875 15.6875L4.3125 14.8125L9.125 10L4.3125 5.1875L5.1875 4.3125L10 9.125L14.8125 4.3125L15.6875 5.1875L10.875 10L15.6875 14.8125L14.8125 15.6875L10 10.875L5.1875 15.6875Z" />
          </svg>
        </button>
      )}
    </div>
  );
};

const Footer: React.FC = () => {
  const navigate = useNavigate();
  return (
    <div className="tw-flex tw-flex-row tw-w-full tw-h-20 tw-min-h-[80px] tw-border-t tw-border-slate-200 tw-mt-auto tw-items-center tw-px-20">
      <button
        className="tw-border tw-border-slate-300 tw-font-medium tw-rounded-md tw-w-32 tw-h-10 tw-select-none hover:tw-bg-slate-100"
        onClick={() => navigate("/")}
      >
        Back
      </button>
    </div>
  );
};

const getStatusStyle = (status: SyncRunStatus): string => {
  switch (status) {
    case SyncRunStatus.Running:
      return "tw-bg-sky-100 tw-border-sky-500 tw-text-sky-600";
    case SyncRunStatus.Completed:
      return "tw-bg-green-100 tw-border-green-500 tw-text-green-600";
    case SyncRunStatus.Failed:
      return "tw-bg-red-100 tw-border-red-500 tw-text-red-500";
    default:
      return "tw-bg-gray-100 tw-border-gray-500 tw-text-gray-500";
  }
};
