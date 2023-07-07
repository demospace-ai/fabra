import { useNavigate, useParams } from "react-router-dom";
import { Loading } from "src/components/loading/Loading";
import { NewObject } from "src/pages/objects/NewObject";
import { useDestination, useObject } from "src/rpc/data";

export const UpdateObject: React.FC = () => {
  const navigate = useNavigate();
  const { objectID } = useParams<{ objectID: string }>();
  const { object } = useObject(Number(objectID));
  const { destination } = useDestination(Number(object?.destination_id));

  if (!object || !destination) {
    return <Loading className="tw-mt-32" />;
  }

  return (
    <NewObject
      existingObject={object}
      existingDestination={destination}
      onComplete={() => {
        navigate(`/objects/${objectID}`);
      }}
    />
  );
};
