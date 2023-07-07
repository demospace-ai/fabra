import { PlusCircleIcon } from "@heroicons/react/20/solid";
import { ChevronRightIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import { useNavigate } from "react-router-dom";
import { Button } from "src/components/button/Button";
import { ConnectionImage } from "src/components/images/Connections";
import { EmptyTable } from "src/components/table/Table";
import { getConnectionType } from "src/rpc/api";
import { useDestinations } from "src/rpc/data";
import { mergeClasses } from "src/utils/twmerge";

const tableHeaderStyle = "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-pr-4 tw-pl-3 sm:tw-pr-6 lg:tw-pr-8 tw-text-left";
const tableCellStyle = "tw-whitespace-nowrap tw-px-3 tw-h-16 tw-text-sm tw-text-slate-800 tw-hidden sm:tw-table-cell";

export const Destinations: React.FC = () => {
  return (
    <div className="tw-py-5 tw-px-10 tw-h-full tw-overflow-scroll">
      <DestinationList />
    </div>
  );
};

const DestinationList: React.FC = () => {
  const { destinations } = useDestinations();
  const navigate = useNavigate();
  return (
    <>
      <div className="tw-flex tw-w-full tw-mb-5 tw-mt-2">
        <div className="tw-flex tw-flex-col tw-justify-end tw-font-bold tw-text-lg">Destinations</div>
        <Button
          className="tw-ml-auto tw-flex tw-justify-center tw-items-center"
          onClick={() => navigate("/destinations/new")}
        >
          <div className="tw-flex tw-flex-col tw-justify-center tw-h-full">
            <PlusCircleIcon className="tw-h-4 tw-inline-block tw-mr-2" />
          </div>
          <div className="tw-flex tw-flex-col tw-justify-center tw-mr-0.5">Add Destination</div>
        </Button>
      </div>
      <div className="tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-bg-white tw-rounded-lg tw-overflow-x-auto tw-overscroll-contain tw-shadow-md">
        {destinations ? (
          <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
            <thead className="tw-bg-slate-100 tw-text-slate-900">
              <tr>
                <th scope="col" className={tableHeaderStyle}>
                  Name
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Type
                </th>
                <th scope="col" className={classNames(tableHeaderStyle, "tw-w-5")}></th>
              </tr>
            </thead>
            <tbody className="tw-divide-y tw-divide-slate-200">
              {destinations!.length > 0 ? (
                destinations!.map((destination, index) => (
                  <tr
                    key={index}
                    className="tw-cursor-pointer hover:tw-bg-slate-50"
                    onClick={() => navigate(`/destinations/${destination.id}`)}
                  >
                    <td className={tableCellStyle}>{destination.display_name}</td>
                    <td className={tableCellStyle}>
                      <div className="tw-flex tw-items-center">
                        <ConnectionImage
                          connectionType={destination.connection.connection_type}
                          className="tw-h-6 tw-mr-1"
                        />
                        {getConnectionType(destination.connection.connection_type)}
                      </div>
                    </td>
                    <td className={mergeClasses(tableCellStyle, "tw-pr-5")}>
                      <ChevronRightIcon className="tw-ml-auto tw-h-4 tw-w-4 tw-text-slate-400" aria-hidden="true" />
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td className={tableCellStyle}>No destinations yet!</td>
                </tr>
              )}
            </tbody>
          </table>
        ) : (
          <EmptyTable />
        )}
      </div>
    </>
  );
};
