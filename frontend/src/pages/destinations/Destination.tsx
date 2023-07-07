import { useNavigate, useParams } from "react-router-dom";
import { BackButton } from "src/components/button/Button";
import { InfoIcon } from "src/components/icons/Icons";
import { ConnectionImage } from "src/components/images/Connections";
import { Loading } from "src/components/loading/Loading";
import { ObjectsListTable } from "src/components/ObjectsListTable";
import { AddObjectButton } from "src/components/ObjectsListTable/AddObjectsButton";
import { PrivateKey } from "src/components/privateKey/PrivateKey";
import { SectionLayout } from "src/components/SectionLayout";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { ConnectionType, getConnectionType } from "src/rpc/api";
import { useDestination, useObjects } from "src/rpc/data";

export const Destination: React.FC = () => {
  const { destinationID } = useParams<{ destinationID: string }>();
  const { destination } = useDestination(Number(destinationID));
  const objectsQuery = useObjects({ destinationID: Number(destinationID) });
  const navigate = useNavigate();
  const objects = objectsQuery.objects ?? [];

  if (!destination) {
    return <Loading />;
  }

  const onAddObjectClick = () => {
    navigate("/objects/new", { state: { destination } });
  };

  return (
    <div className="tw-py-5 tw-px-10 tw-h-full tw-overflow-scroll">
      <BackButton onClick={() => navigate("/destinations")} />
      <div className="tw-flex tw-w-full tw-mb-2 tw-mt-4">
        <div className="tw-flex tw-flex-row tw-items-center tw-font-bold tw-text-xl">{destination.display_name}</div>
      </div>
      <div className="tw-flex tw-flex-col tw-w-fit tw-flex-wrap tw-items-start tw-px-3 tw-pt-1 tw-pb-2 tw-mt-3 tw-mb-5 tw-bg-white tw-border tw-border-slate-200 tw-rounded-md">
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Destination ID: </span>
          {destination.id}
        </div>
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-1">
          <span className="tw-font-medium tw-whitespace-pre">Destination Type: </span>
          {getConnectionType(destination.connection.connection_type)}
          <ConnectionImage connectionType={destination.connection.connection_type} className="tw-h-5 tw-ml-1" />
        </div>
      </div>
      <div className="tw-font-bold tw-text-base">Configuration</div>
      {destination.connection.connection_type === ConnectionType.Webhook && (
        <>
          <div className="tw-flex tw-flex-row tw-items-center tw-mt-3 tw-mb-2">
            <span>Webhook Signing Key</span>
            <Tooltip
              placement="right"
              label="Use this signing key to verify the signature of incoming webhook requests from Fabra."
              interactive
              maxWidth={500}
            >
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
          <PrivateKey keyValue={destination.webhook_signing_key} />
        </>
      )}
      <section className="tw-flex tw-items-center tw-justify-between tw-mt-8">
        <div className="tw-text-base tw-font-bold">Objects</div>
        <AddObjectButton onClick={onAddObjectClick} />
      </section>
      <SectionLayout className="tw-mt-4">
        <ObjectsListTable objects={objects} />
      </SectionLayout>
    </div>
  );
};
