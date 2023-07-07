import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { BackButton } from "src/components/button/Button";
import { ShouldNotReach } from "src/components/ShouldNotReach";
import { Step } from "src/pages/objects/helpers";
import { DestinationSetup } from "src/pages/objects/NewObject/DestinationSetupStep";
import { ExistingObjectFields } from "src/pages/objects/NewObject/ExistingObjectFieldsStep";
import { Finalize } from "src/pages/objects/NewObject/FinalizeStep";
import { NewObjectFields } from "src/pages/objects/NewObject/NewObjectFieldsStep";
import { useStateMachine } from "src/pages/objects/NewObject/state";
import { Destination, FabraObject } from "src/rpc/api";

export type NewObjectProps = {
  existingObject?: FabraObject;
  existingDestination?: Destination;
  onComplete?: () => void;
};

export const NewObject: React.FC<NewObjectProps> = ({ existingDestination, existingObject, ...rest }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const maybeDestination = location.state?.destination as Destination | undefined;
  const onComplete = rest.onComplete
    ? rest.onComplete
    : () => {
        navigate("/objects");
      };
  const { advanceToObjectFields, advanceToFinalizeObject, state, back } = useStateMachine(
    {
      existingDestination,
      existingObject,
      maybeDestination,
    },
    onComplete,
  );

  let content: React.ReactElement;
  switch (state.step) {
    case Step.UnsupportedConnectionType: {
      content = (
        <ShouldNotReach error={new Error(`Found an unsupported connection type ${state.connectionType}`)}>
          <div className="tw-text-center">
            <h3 className="tw-font-bold tw-text-xl tw-mb-2">Unable to proceed</h3>
            <div className="tw-max-w-xs">
              {state.message ?? "Unsupported connection type. Please contact the Fabra team for more help."}
            </div>
          </div>
        </ShouldNotReach>
      );
      break;
    }
    case Step.Initial: {
      content = (
        <DestinationSetup
          isUpdate={!!existingObject}
          initialFormState={state.destinationSetup}
          onComplete={(values) => {
            advanceToObjectFields(values);
          }}
        />
      );
      // TODO: prompt if they want to exit here
      break;
    }
    case Step.ExistingFields: {
      const { destinationSetup } = state;
      content = (
        <ExistingObjectFields
          destinationSetupData={state.destinationSetup}
          isUpdate={!!existingObject}
          initialFormState={state.objectFields}
          onComplete={(values) => {
            advanceToFinalizeObject(destinationSetup, values);
          }}
        />
      );
      break;
    }
    case Step.CreateFields: {
      const { destinationSetup } = state;
      content = (
        <NewObjectFields
          initialFormState={state.objectFields}
          isUpdate={!!existingObject}
          onComplete={(values) => {
            advanceToFinalizeObject(destinationSetup, values);
          }}
        />
      );
      break;
    }
    case Step.Finalize: {
      content = (
        <Finalize
          isUpdate={!!existingObject}
          initialFormState={state.finalize}
          existingObject={existingObject}
          objectFields={state.objectFields}
          destinationSetup={state.destinationSetup}
          onComplete={onComplete}
        />
      );
      break;
    }
    default:
      content = (
        <ShouldNotReach error={new Error("Found an invalid step in NewObject")}>
          <div className="tw-text-center">
            <h3 className="tw-font-bold tw-text-xl tw-mb-2">Unable to proceed</h3>
            <div className="tw-max-w-xs">Unknown step. Please contact the Fabra team for more help.</div>
          </div>
        </ShouldNotReach>
      );
  }

  return (
    <div className="tw-flex tw-flex-col tw-mb-10">
      <BackButton onClick={back} />
      <div className="tw-flex tw-flex-col tw-w-[900px] tw-mt-8 tw-mb-24 tw-py-12 tw-px-10 tw-mx-auto tw-bg-white tw-rounded-lg tw-shadow-md tw-items-center">
        {content}
      </div>
    </div>
  );
};
