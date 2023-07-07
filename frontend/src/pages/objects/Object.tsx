import { useNavigate, useParams } from "react-router-dom";
import { BackButton, Button } from "src/components/button/Button";
import { Loading } from "src/components/loading/Loading";
import {
  needsCursorField,
  needsEndCustomerId,
  needsPrimaryKey,
  syncModeToString,
  TargetType,
  targetTypeToString,
  toReadableFrequency,
} from "src/rpc/api";
import { useObject } from "src/rpc/data";

const tableHeaderStyle =
  "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-px-4 sm:tw-pr-6 lg:tw-pr-8 tw-text-left tw-whitespace-nowrap";
const tableCellStyle =
  "tw-whitespace-nowrap tw-left tw-overflow-hidden tw-py-4 tw-pl-4 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export const Object: React.FC = () => {
  const navigate = useNavigate();
  const { objectID } = useParams<{ objectID: string }>();
  const { object } = useObject(Number(objectID));

  if (!object) {
    return <Loading />;
  }

  return (
    <div className="xl:tw-w-3/5 md:tw-w-full tw-flex tw-flex-col tw-mb-10">
      <BackButton onClick={() => navigate("/objects")} />
      <div className="tw-flex tw-flex-row tw-items-center tw-font-bold tw-text-xl tw-my-4">
        <span className="tw-grow">{object.display_name}</span>
        <Button
          className="tw-ml-auto tw-px-3 tw-py-1 tw-rounded-md tw-font-medium tw-text-base tw-bg-transparent hover:tw-bg-slate-100 tw-text-blue-600 tw-mr-2"
          onClick={() => {
            navigate("./update");
          }}
        >
          Edit
        </Button>
      </div>
      <div className="tw-flex tw-flex-col tw-flex-wrap tw-items-start tw-p-4 tw-mb-5 tw-bg-white tw-border tw-border-slate-200 tw-rounded-md">
        <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Destination ID: </span>
          {object.destination_id}
        </div>
        <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Target Type: </span>
          {targetTypeToString(object.target_type)}
        </div>
        {object.target_type !== TargetType.Webhook && (
          <>
            <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
              <span className="tw-font-medium tw-whitespace-pre">Namespace: </span>
              {object.namespace}
            </div>
            <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
              <span className="tw-font-medium tw-whitespace-pre">Table Name: </span>
              {object.table_name}
            </div>
          </>
        )}
        <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Sync Mode: </span>
          {syncModeToString(object.sync_mode)}
        </div>
        {needsCursorField(object.sync_mode) && (
          <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
            <span className="tw-font-medium tw-whitespace-pre">Cursor Field: </span>
            {object.cursor_field}
          </div>
        )}
        {needsPrimaryKey(object.sync_mode) && (
          <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
            <span className="tw-font-medium tw-whitespace-pre">Primary Key: </span>
            {object.primary_key}
          </div>
        )}
        {needsEndCustomerId(object.target_type) && (
          <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
            <span className="tw-font-medium tw-whitespace-pre">End Customer ID Field: </span>
            {object.end_customer_id_field}
          </div>
        )}
        <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Recurring: </span>
          {object.recurring ? "Yes" : "No"}
        </div>
        {object.recurring && (
          <div className="tw-flex tw-flex-row tw-items-center tw-text-base tw-mt-1">
            <span className="tw-font-medium tw-whitespace-pre">Frequency: </span>
            {toReadableFrequency(object.frequency, object.frequency_units)}
          </div>
        )}
      </div>
      <div className="tw-font-bold tw-text-base tw-mt-4 tw-mb-2">Object Fields</div>
      <div className="tw-border tw-border-solid tw-border-slate-200 tw-bg-white tw-rounded-lg tw-overflow-auto tw-overscroll-contain tw-shadow-md">
        <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
          <thead className="tw-bg-slate-100 tw-text-slate-900">
            <tr>
              <th scope="col" className={tableHeaderStyle}>
                Name
              </th>
              <th scope="col" className={tableHeaderStyle}>
                Type
              </th>
            </tr>
          </thead>
          <tbody className="tw-divide-y tw-divide-slate-200 tw-bg-white">
            {object.object_fields.map((objectField, index) => {
              return (
                <tr key={index} className="tw-cursor-pointer hover:tw-bg-slate-50">
                  <td className={tableCellStyle}>{objectField.name}</td>
                  <td className={tableCellStyle}>{objectField.type}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};
