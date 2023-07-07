import classNames from "classnames";
import React from "react";
import { ResultRow, Schema } from "src/rpc/api";

type QueryResultsProps = {
  schema: Schema;
  results: ResultRow[];
};

const QueryResultsTable: React.FC<QueryResultsProps> = (props) => {
  // TODO: implement pagination to not render too many results
  return (
    <div className="tw-rounded-md tw-overflow-auto tw-max-h-full tw-flow-root" style={{ contain: "paint" }}>
      <table>
        <ResultsSchema schema={props.schema} />
        <tbody className="tw-py-2">
          {props.results.map((resultRow, index) => {
            return (
              <tr key={index}>
                <td key={-1} className={classNames("tw-px-3 tw-py-2 tw-text-right tw-bg-gray-100 tw-tabular-nums")}>
                  <div className="tw-h-5 tw-whitespace-nowrap">{index + 1}</div>
                </td>
                {resultRow.map((resultValue, valueIndex) => {
                  return (
                    <td
                      key={valueIndex}
                      className={classNames("tw-pl-3 tw-pr-5 tw-py-2 tw-text-left last:tw-w-full focus:tw-bg-blue-300")}
                    >
                      <div className="tw-h-5 tw-whitespace-nowrap">
                        {resultValue ? JSON.stringify(resultValue) : <span className="tw-text-gray-400">null</span>}
                      </div>
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

const ResultsSchema: React.FC<{ schema: Schema }> = ({ schema }) => {
  return (
    <thead className="tw-sticky tw-top-0">
      <tr>
        <th key={-1} scope="col" className="tw-pl-3 tw-pr-5 tw-py-2 tw-bg-gray-100"></th>
        {schema.map((fieldSchema, index) => {
          return (
            <th key={index} scope="col" className="tw-pl-3 tw-pr-5 tw-py-3 tw-text-left tw-bg-gray-100 ">
              <div className="tw-whitespace-nowrap">{fieldSchema.name}</div>
            </th>
          );
        })}
      </tr>
    </thead>
  );
};

export const MemoizedResultsTable = React.memo(QueryResultsTable);
