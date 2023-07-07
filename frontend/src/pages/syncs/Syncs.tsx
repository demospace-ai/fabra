import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { useNavigate } from "react-router-dom";
import { ConnectionImage } from "src/components/images/Connections";
import { EmptyTable } from "src/components/table/Table";
import { useSyncs } from "src/rpc/data";
import { mergeClasses } from "src/utils/twmerge";

const tableHeaderStyle = "tw-sticky tw-top-0 tw-z-0 tw-py-3.5 tw-pl-3 tw-text-left tw-whitespace-nowrap";
const tableCellStyle = "tw-whitespace-nowrap tw-left tw-pl-3 tw-min-w-[200px] tw-h-16 tw-text-sm tw-text-slate-800";

export const Syncs: React.FC = () => {
  const navigate = useNavigate();
  const { syncs, objects, sources } = useSyncs();
  const objectIdMap = new Map(objects?.map((object) => [object.id, object]));
  const sourceIdMap = new Map(sources?.map((source) => [source.id, source]));

  return (
    <div className="tw-py-5 tw-px-10 tw-h-full tw-overflow-scroll">
      <div className="tw-flex tw-w-full tw-mb-5 tw-mt-2 tw-h-[29px]">
        <div className="tw-flex tw-flex-col tw-justify-end tw-font-bold tw-text-lg">Syncs</div>
      </div>
      <div className="tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-bg-white tw-rounded-lg tw-overflow-x-auto tw-overscroll-contain tw-shadow-md">
        {syncs ? (
          <table className="tw-min-w-full tw-border-spacing-0 tw-divide-y tw-divide-slate-200">
            <thead className="tw-bg-slate-100 tw-text-slate-900">
              <tr>
                <th scope="col" className={tableHeaderStyle}>
                  Name
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  End Customer ID
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Object
                </th>
                <th scope="col" className={tableHeaderStyle}>
                  Source
                </th>
                <th scope="col" className={mergeClasses(tableHeaderStyle, "tw-w-5")}>
                  <span className="tw-sr-only">Continue</span>
                </th>
              </tr>
            </thead>
            <tbody className="tw-divide-y tw-divide-slate-200">
              {syncs!.length > 0 ? (
                syncs!.map((sync, index) => {
                  const object = objectIdMap.get(sync.object_id);
                  const source = sourceIdMap.get(sync.source_id);
                  return (
                    <tr
                      key={index}
                      className="tw-cursor-pointer hover:tw-bg-slate-50"
                      onClick={() => navigate(`/syncs/${sync.id}`)}
                    >
                      <td className={tableCellStyle}>{sync.display_name}</td>
                      <td className={tableCellStyle}>{sync.end_customer_id}</td>
                      <td className={tableCellStyle}>{object?.display_name}</td>
                      <td className={mergeClasses(tableCellStyle, "tw-flex tw-flex-row tw-items-center")}>
                        <ConnectionImage
                          connectionType={source!.connection.connection_type}
                          className="tw-h-6 tw-mr-1.5"
                        />
                        {source?.display_name}
                      </td>
                      <td className={mergeClasses(tableCellStyle, "tw-w-full tw-pr-5")}>
                        <ChevronRightIcon className="tw-ml-auto tw-h-4 tw-w-4 tw-text-slate-400" aria-hidden="true" />
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td className={tableCellStyle}>No syncs yet!</td>
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
