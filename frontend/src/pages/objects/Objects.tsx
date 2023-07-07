import { useLocation, useNavigate } from "react-router-dom";
import { ObjectsListTable } from "src/components/ObjectsListTable";
import { AddObjectButton } from "src/components/ObjectsListTable/AddObjectsButton";
import { SectionLayout } from "src/components/SectionLayout";
import { EmptyTable } from "src/components/table/Table";
import { useObjects } from "src/rpc/data";

export const ObjectsList: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { objects } = useObjects();

  return (
    <>
      <div className="tw-flex tw-w-full tw-mb-5 tw-mt-2 tw-justify-between tw-items-center">
        <div className="tw-font-bold tw-text-lg">Objects</div>
        {location.pathname === "/objects" && <AddObjectButton onClick={() => navigate("/objects/new")} />}
      </div>
      <SectionLayout>{objects ? <ObjectsListTable objects={objects} /> : <EmptyTable />}</SectionLayout>
    </>
  );
};
