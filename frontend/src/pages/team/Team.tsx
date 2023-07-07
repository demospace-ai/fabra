import { PlusCircleIcon } from "@heroicons/react/20/solid";
import classNames from "classnames";
import { LinkButton } from "src/components/button/Button";
import { Loading } from "src/components/loading/Loading";
import { useUsers } from "src/rpc/data";

const tableHeaderStyle =
  "tw-sticky tw-top-0 tw-z-0 tw-border-b tw-border-slate-300 tw-py-3.5 tw-pr-4 tw-pl-3 backdrop-blur backdrop-filter sm:tw-pr-6 lg:tw-pr-8 tw-text-left";
const tableCellStyle = "tw-whitespace-nowrap tw-px-3 tw-py-4 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export const Team: React.FC = () => {
  const { users } = useUsers();
  return (
    <div className="tw-py-5 tw-px-10">
      <div className="tw-flex tw-w-full tw-mb-5 tw-mt-2">
        <div className="tw-flex tw-flex-col tw-justify-end tw-font-bold tw-text-lg">Users</div>
        <LinkButton
          className="tw-ml-auto tw-flex"
          href="mailto:nick@fabra.io?subject=Join me on Fabra!&body=Sign up here: https://app.fabra.io"
        >
          <div className="tw-flex tw-flex-col tw-justify-center tw-h-full">
            <PlusCircleIcon className="tw-h-4 tw-inline-block tw-mr-2" />
          </div>
          <div className="tw-flex tw-flex-col tw-justify-center tw-mr-0.5">Invite teammate</div>
        </LinkButton>
      </div>
      <div className="tw-border tw-border-solid tw-border-slate-300 tw-rounded-lg tw-max-h-64 tw-overflow-x-auto tw-overscroll-contain tw-bg-white tw-shadow-md">
        {users ? (
          <table className="tw-min-w-full tw-border-separate tw-border-spacing-0">
            <thead className="tw-bg-slate-600 tw-text-white">
              <tr>
                <th scope="col" className={tableHeaderStyle}>
                  Name
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Email
                </th>
                <th scope="col" className={classNames(tableHeaderStyle, "tw-w-5")}></th>
                <th scope="col" className={classNames(tableHeaderStyle, "tw-w-5")}></th>
              </tr>
            </thead>
            <tbody>
              {users!.length > 0 ? (
                users!.map((user, index) => (
                  <tr key={index} className="tw-border-b tw-border-solid tw-border-slate-200 last:tw-border-0">
                    <td className={tableCellStyle}>{user.name}</td>
                    <td className={tableCellStyle}>{user.email}</td>
                    <td className={tableCellStyle}>
                      <div
                        className="tw-cursor-pointer tw-font-medium tw-select-none tw-text-blue-700 hover:tw-text-blue-900"
                        onClick={() => null}
                      >
                        Delete
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td className={tableCellStyle}>No users yet!</td>
                </tr>
              )}
            </tbody>
          </table>
        ) : (
          <Loading className="tw-my-5" />
        )}
      </div>
    </div>
  );
};
