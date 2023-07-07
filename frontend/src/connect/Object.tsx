import React, { useState } from "react";
import { Button } from "src/components/button/Button";
import { Checkbox } from "src/components/checkbox/Checkbox";
import { InfoIcon } from "src/components/icons/Icons";
import { DataPreview } from "src/components/images/DataPreview";
import { Loading } from "src/components/loading/Loading";
import { MemoizedResultsTable } from "src/components/queryResults/QueryResults";
import { ObjectSelector, SourceNamespaceSelector, SourceTableSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { SetupSyncProps } from "src/connect/state";
import { sendLinkTokenRequest } from "src/rpc/ajax";
import { FabraObject, LinkGetPreview, LinkGetPreviewRequest } from "src/rpc/api";
import { useMutation } from "src/utils/queryHelpers";

export const ObjectSetup: React.FC<SetupSyncProps> = (props) => {
  const setObject = (object: FabraObject) =>
    props.setState((state) => ({ ...state, object: object, fieldMappings: undefined }));
  const setNamespace = (namespace: string) =>
    props.setState((state) => ({ ...state, namespace: namespace, tableName: undefined }));
  const setTableName = (tableName: string) => props.setState((state) => ({ ...state, tableName: tableName }));
  const [limitPreview, setLimitPreview] = useState<boolean>(true);

  const fetchPreviewMutation = useMutation(async () => {
    const {
      state: { source, namespace, tableName },
    } = props;
    if (!source?.id || !namespace || !tableName) {
      throw new Error("missing required fields");
    }
    const payload: LinkGetPreviewRequest = {
      source_id: source.id,
      namespace,
      table_name: tableName,
    };

    const response = await sendLinkTokenRequest(LinkGetPreview, props.linkToken, payload);
    return response;
  });

  const previewData = fetchPreviewMutation.data?.data;
  const previewSchema = fetchPreviewMutation.data?.schema;

  return (
    <div className="tw-w-full tw-pl-20 tw-pr-[72px] tw-flex tw-flex-col">
      <div className="tw-mb-5 tw-text-2xl tw-font-semibold tw-text-slate-900">Define the data model to sync</div>
      <div className="tw-w-[50%] tw-min-w-[400px]">
        <div className="tw-text-base tw-font-medium tw-mb-1 tw-text-slate-800">Select object to create</div>
        <div className="tw-text-slate-600 tw-text-sm">
          This is the object that will be created from the data you define in this sync configuration.
        </div>
        <ObjectSelector object={props.state.object} setObject={setObject} linkToken={props.linkToken} />
        <div className="tw-text-base tw-font-medium tw-mt-8 tw-mb-1 tw-text-slate-800">Select a table to sync from</div>
        <div className="tw-text-slate-600 tw-text-sm">
          This is where the data will be pulled from in your own data warehouse.
        </div>
        <SourceNamespaceSelector
          namespace={props.state.namespace}
          setNamespace={setNamespace}
          linkToken={props.linkToken}
          source={props.state.source}
          dropdownHeight="tw-max-h-40"
        />
        <SourceTableSelector
          tableName={props.state.tableName}
          setTableName={setTableName}
          linkToken={props.linkToken}
          source={props.state.source}
          namespace={props.state.namespace}
          dropdownHeight="tw-max-h-40"
        />
        <div>
          <div className="tw-flex tw-flex-row tw-mt-4 tw-items-center">
            <Button className="tw-h-10 tw-w-32" onClick={() => fetchPreviewMutation.mutate()}>
              {fetchPreviewMutation.isLoading ? <Loading light /> : "Preview"}
            </Button>
            <Checkbox
              checked={limitPreview}
              onCheckedChange={() => setLimitPreview(!limitPreview)}
              className="tw-ml-4 tw-mr-2 tw-h-5 tw-w-5"
            />
            <span>Limit preview to 100 records</span>
            <Tooltip
              placement="right"
              label="Automatically add a LIMIT expression to the query to keep the number of rows fetched to 100."
            >
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
          {fetchPreviewMutation.error && (
            <div className="tw-text-red-500 tw-mt-1">{fetchPreviewMutation.error.message}</div>
          )}
        </div>
      </div>
      <div className="tw-mt-10 tw-h-full tw-max-h-[400px] tw-w-full tw-rounded-md tw-border tw-border-gray-200">
        {previewData && previewSchema ? (
          <MemoizedResultsTable schema={previewSchema} results={previewData} />
        ) : (
          <div className="tw-h-full tw-w-full tw-rounded-md tw-bg-slate-50 tw-justify-center tw-items-center tw-text-center tw-flex tw-flex-col">
            <DataPreview className="tw-h-20 tw-mb-3" animate={fetchPreviewMutation.isLoading} />
            <div className="tw-text-xl tw-font-semibold tw-mb-1 tw-text-slate-600">Preview your data</div>
            <div className="tw-text-slate-500">A preview of the resulting rows will appear here!</div>
          </div>
        )}
        <div className="tw-pb-24"></div>
      </div>
    </div>
  );
};
