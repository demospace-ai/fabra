import { ChevronRightIcon, PlusCircleIcon } from "@heroicons/react/24/outline";
import React, { useEffect } from "react";
import { Button } from "src/components/button/Button";
import { AddDatabase, InfoIcon } from "src/components/icons/Icons";
import { ConnectionImage } from "src/components/images/Connections";
import { EmptyTable } from "src/components/table/Table";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { SetupSyncProps, SyncSetupStep } from "src/connect/state";
import { Source } from "src/rpc/api";
import { useLinkSources } from "src/rpc/data";

export const Sources: React.FC<SetupSyncProps> = ({ linkToken, state, setState }) => {
  const { sources } = useLinkSources(linkToken);

  // Skip the selection step and go straight to setting up a new source if none exist.
  useEffect(() => {
    if (sources !== undefined) {
      if (sources.length === 0) {
        setState((state) => ({ ...state, step: SyncSetupStep.ChooseSourceType, skippedSourceSelection: true }));
      }
    }
  }, [sources, setState]);

  const setExistingSource = (source: Source) => {
    setState({
      ...state,
      source: source,
      step: SyncSetupStep.ChooseData,
      skippedSourceSetup: true,
      namespace: undefined,
      tableName: undefined,
    });
  };

  return (
    <div className="tw-w-full tw-px-20">
      <div className="tw-flex tw-flex-row tw-items-center">
        <span className="tw-text-2xl tw-font-semibold tw-text-slate-900">Your sources</span>
        <Tooltip placement="right" maxWidth="500px" label="These are the data sources you've setup previously.">
          <InfoIcon className="tw-ml-1.5 tw-h-3.5 tw-fill-slate-400" />
        </Tooltip>
        <Button
          className="tw-ml-auto tw-flex tw-flex-row tw-items-center tw-whitespace-nowrap tw-h-8"
          onClick={() => setState({ ...state, step: SyncSetupStep.ChooseSourceType })}
        >
          <PlusCircleIcon className="tw-h-5 tw-mr-2 tw-stroke-2" />
          <span className="tw-mr-1">New Source</span>
        </Button>
      </div>
      <div className="tw-text-left tw-mt-2 tw-text-slate-600">
        Setup additional syncs from your existing sources, or add a new source.
      </div>
      <div className="tw-flex tw-flex-row tw-items-center tw-justify-center tw-w-full tw-pb-4">
        <div className="tw-w-full">
          <div className="tw-mt-5 tw-flow-root tw-select-none">
            <div className="tw-inline-block tw-min-w-full tw-py-2 tw-align-middle">
              <div className="tw-overflow-auto tw-shadow tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-rounded-md">
                {sources ? (
                  <table className="tw-min-w-full tw-divide-y tw-divide-slate-200">
                    <thead className="tw-bg-slate-100">
                      <tr>
                        <th
                          scope="col"
                          className="tw-py-3.5 tw-pl-4 tw-pr-3 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                        >
                          Name
                        </th>
                        <th
                          scope="col"
                          className="tw-px-3 tw-py-3.5 tw-text-left tw-text-sm tw-font-semibold tw-text-slate-900"
                        >
                          Connection Type
                        </th>
                        <th scope="col" className="tw-relative tw-py-3.5 tw-pl-3">
                          <span className="tw-sr-only">Continue</span>
                        </th>
                      </tr>
                    </thead>
                    <tbody className="tw-divide-y tw-divide-slate-200 tw-bg-white">
                      {sources.length > 0 ? (
                        sources.map((source) => (
                          <tr
                            key={source.id}
                            className="tw-cursor-pointer hover:tw-bg-slate-50"
                            onClick={() => setExistingSource(source)}
                          >
                            <td className="tw-whitespace-nowrap tw-py-4 tw-pl-4 tw-pr-3 tw-text-sm tw-font-medium tw-text-slate-900 tw-flex tw-flex-row tw-items-center">
                              <ConnectionImage
                                connectionType={source.connection.connection_type}
                                className="tw-h-6 tw-mr-1.5"
                              />
                              {source.display_name}
                            </td>
                            <td className="tw-whitespace-nowrap tw-px-3 tw-py-4 tw-text-sm tw-text-slate-500">
                              {source.connection.connection_type}
                            </td>
                            <td className="tw-pr-4" align="right">
                              <ChevronRightIcon className="tw-h-4 tw-w-4 tw-text-slate-400" aria-hidden="true" />
                            </td>
                          </tr>
                        ))
                      ) : (
                        <tr>
                          <td className="tw-whitespace-nowrap tw-pt-10 tw-pb-16 tw-pl-12 tw-pr-3 tw-text-sm tw-font-medium tw-text-slate-900 tw-flex tw-flex-row tw-items-center">
                            <AddDatabase className="tw-h-16 tw-fill-slate-400" />
                            <div className="tw-flex tw-flex-col tw-ml-8">
                              <span className="tw-text-lg tw-font-semibold tw-text-slate-500">Add a source</span>
                              <span className="tw-text-slate-500">Start syncing your data by adding a source.</span>
                              <Button
                                className="tw-w-32 tw-px-0 tw-mt-4"
                                onClick={() => setState({ ...state, step: SyncSetupStep.ChooseSourceType })}
                              >
                                Add a Source
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
          </div>
        </div>
      </div>
    </div>
  );
};
