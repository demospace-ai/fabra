import { Loading } from "src/components/loading/Loading";

const tableHeaderStyle = "tw-sticky tw-top-0 tw-z-0 tw-h-12 sm:tw-pr-6 lg:tw-pr-8 tw-text-left tw-whitespace-nowrap";
const tableCellStyle =
  "tw-whitespace-nowrap tw-left tw-overflow-hidden tw-py-4 tw-pl-4 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export const EmptyTable: React.FC = () => {
  return (
    <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
      <thead className="tw-sticky tw-top-0 tw-bg-slate-100">
        <tr>
          <th scope="col" className={tableHeaderStyle}></th>
        </tr>
      </thead>
      <tbody className="tw-divide-y tw-divide-slate-200 tw-bg-white">
        <tr>
          <td className={tableCellStyle}>
            <Loading className="tw-my-10" />
          </td>
        </tr>
      </tbody>
    </table>
  );
};
