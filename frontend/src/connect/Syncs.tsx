import { ChevronRightIcon, PlusCircleIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "src/components/button/Button";
import { AddDatabase } from "src/components/icons/Icons";
import { ConnectionImage } from "src/components/images/Connections";
import { EmptyTable } from "src/components/table/Table";
import { LinkGetSyncs } from "src/rpc/api";
import { useLinkSyncs } from "src/rpc/data";
import { mergeClasses } from "src/utils/twmerge";
import { mutate } from "swr";

export const Syncs: React.FC<{ linkToken: string; close: (() => void) | undefined }> = ({ linkToken, close }) => {
  return (
    <div className="tw-w-full tw-h-full tw-flex tw-flex-col">
      <Header close={close} />
      <SyncList linkToken={linkToken} />
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

const SyncList: React.FC<{ linkToken: string }> = ({ linkToken }) => {
  const tableCellStyle = "tw-whitespace-nowrap tw-py-4 tw-pl-4 tw-pr-3 tw-text-sm tw-text-slate-900";
  const navigate = useNavigate();
  const { syncs, objects, sources } = useLinkSyncs(linkToken);
  const objectIdMap = new Map(objects?.map((object) => [object.id, object]));
  const sourceIdMap = new Map(sources?.map((source) => [source.id, source]));

  // Tell SWRs to refetch syncs whenever link token changes
  useEffect(() => {
    mutate({ LinkGetSyncs });
  }, [linkToken]);

  return (
    <div className="tw-mt-2 tw-px-20 tw-pb-16 tw-flex tw-flex-col tw-overflow-auto">
      <div className="tw-flex tw-w-full tw-mt-2">
        <div className="tw-flex-col">
          <div className="tw-flex tw-justify-start tw-font-bold tw-text-2xl">Your syncs</div>
          <div className="tw-mt-2 tw-text-slate-600">
            Setup a sync to connect your data source and map it to fields in the application.
          </div>
        </div>
        <Button
          className="tw-flex tw-flex-row tw-items-center tw-ml-auto tw-h-8 tw-whitespace-nowrap"
          onClick={() => navigate("/newsync")}
        >
          <PlusCircleIcon className="tw-h-5 tw-mr-2 tw-stroke-2" />
          <span className="tw-mr-1">New Sync</span>
        </Button>
      </div>
      <div className="tw-mt-10 tw-overflow-auto tw-shadow tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-rounded-md">
        {syncs ? (
          <table className="tw-min-w-full tw-divide-y tw-divide-slate-200">
            <thead className="tw-bg-slate-100">
              {syncs.length > 0 ? (
                <tr>
                  <th
                    scope="col"
                    className="tw-py-3.5 tw-pl-4 tw-pr-3 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                  >
                    Name
                  </th>
                  <th
                    scope="col"
                    className="tw-py-3.5 tw-pl-4 tw-pr-3 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                  >
                    Object
                  </th>
                  <th
                    scope="col"
                    className="tw-py-3.5 tw-pl-4 tw-pr-3 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                  >
                    Source
                  </th>
                  <th scope="col" className="tw-relative tw-py-3.5 tw-pl-3">
                    <span className="tw-sr-only">Continue</span>
                  </th>
                </tr>
              ) : (
                <tr>
                  <th
                    scope="col"
                    className="tw-py-3.5 tw-pl-4 tw-pr-3 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                  >
                    Name
                  </th>
                </tr>
              )}
            </thead>
            <tbody className="tw-divide-y tw-divide-slate-200 tw-bg-white">
              {syncs.length > 0 ? (
                syncs.map((sync) => {
                  const object = objectIdMap.get(sync.object_id);
                  const source = sourceIdMap.get(sync.source_id);
                  return (
                    <tr
                      key={sync.id}
                      className="tw-cursor-pointer hover:tw-bg-slate-50"
                      onClick={() => navigate(`/sync/${sync.id}`)}
                    >
                      <td className={tableCellStyle}>{sync.display_name}</td>
                      <td className={tableCellStyle}>{object?.display_name}</td>
                      <td className={mergeClasses(tableCellStyle, "tw-flex tw-items-center")}>
                        <ConnectionImage
                          connectionType={source!.connection.connection_type}
                          className="tw-h-6 tw-mr-1.5"
                        />
                        {source?.display_name}
                      </td>
                      <td className="tw-pr-4" align="right">
                        <ChevronRightIcon className="tw-h-4 tw-w-4 tw-text-slate-400" aria-hidden="true" />
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td className="tw-whitespace-nowrap tw-pt-10 tw-pb-16 tw-pl-12 tw-pr-3 tw-text-sm tw-font-medium tw-text-slate-900 tw-flex tw-flex-row tw-items-center">
                    <AddDatabase className="tw-h-16 tw-fill-slate-400" />
                    <div className="tw-flex tw-flex-col tw-ml-8">
                      <span className="tw-text-lg tw-font-semibold tw-text-slate-500">No syncs yet.</span>
                      <span className="tw-mt-1 tw-text-slate-500">
                        Setup a sync to connect your data and map it to fields in the application.
                      </span>
                      <Button className="tw-w-32 tw-px-0 tw-mt-4" onClick={() => navigate("/newsync")}>
                        Add a Sync
                      </Button>
                    </div>
                  </td>
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
