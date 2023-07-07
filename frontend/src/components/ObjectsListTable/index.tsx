import { ChevronRightIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import { useNavigate } from "react-router-dom";
import { FabraObject } from "src/rpc/api";
import { mergeClasses } from "src/utils/twmerge";

interface ObjectsListTableProps {
  objects: FabraObject[];
  emptyTableComponent?: React.ReactNode;
}

const tableHeaderStyle = "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-pr-4 tw-pl-3 sm:tw-pr-6 lg:tw-pr-8 tw-text-left";
const tableCellStyle = "tw-whitespace-nowrap tw-px-3 tw-h-16 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export function ObjectsListTable({ objects, ...props }: ObjectsListTableProps) {
  const emptyTableComponent = props.emptyTableComponent ?? <DefaultEmptyTableRow />;
  const navigate = useNavigate();
  return (
    <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
      <thead className="tw-bg-slate-100 tw-text-slate-900">
        <tr>
          <th scope="col" className={tableHeaderStyle}>
            Name
          </th>
          <th scope="col" className={classNames(tableHeaderStyle, "tw-w-5")}></th>
        </tr>
      </thead>
      <tbody className="tw-divide-y tw-divide-slate-200">
        {objects.length > 0
          ? objects.map((object, index) => (
              <tr
                key={index}
                className="tw-cursor-pointer hover:tw-bg-slate-50"
                onClick={() => navigate(`/objects/${object.id}`)}
              >
                <td className={tableCellStyle}>{object.display_name}</td>
                <td className={mergeClasses(tableCellStyle, "tw-pr-5")}>
                  <ChevronRightIcon className="tw-ml-auto tw-h-4 tw-w-4 tw-text-slate-400" aria-hidden="true" />
                </td>
              </tr>
            ))
          : emptyTableComponent}
      </tbody>
    </table>
  );
}

function DefaultEmptyTableRow() {
  return (
    <tr>
      <td className={tableCellStyle}>No objects yet!</td>
    </tr>
  );
}
