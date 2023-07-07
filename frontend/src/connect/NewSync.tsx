import { CheckIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import React, { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "src/components/button/Button";
import { ErrorDisplay } from "src/components/error/Error";
import { Loading } from "src/components/loading/Loading";
import { useConnectShowToast } from "src/components/notifications/Notifications";
import { FabraDisplayOptions } from "src/connect/ConnectApp";
import { NewSourceConfiguration } from "src/connect/Connection";
import { FinalizeSync } from "src/connect/Finalize";
import { ObjectSetup } from "src/connect/Object";
import { Sources } from "src/connect/Sources";
import {
  createNewSource,
  FieldMappingState,
  INITIAL_SETUP_STATE,
  SetupSyncState,
  SyncSetupStep,
  useCreateNewSync,
  validateObjectSetup,
} from "src/connect/state";
import { WarehouseSelector } from "src/connect/Warehouse";
import { useObject } from "src/rpc/data";
import { useMutation } from "src/utils/queryHelpers";

export const NewSync: React.FC<{ linkToken: string; close: (() => void) | undefined } & FabraDisplayOptions> = ({
  linkToken,
  close,
  supportEmail,
  docsLink,
}) => {
  const [state, setState] = useState<SetupSyncState>(INITIAL_SETUP_STATE);
  const [prevObject, setPrevObject] = useState<Object | undefined>(undefined);
  const { object } = useObject(state.object?.id, linkToken);
  const navigate = useNavigate();

  // Setup the initial values for the field mappings
  if (object && object !== prevObject) {
    setPrevObject(object);
    const fieldMappings: FieldMappingState[] = object
      ? object.object_fields
          .filter((objectField) => !objectField.omit)
          .map((objectField) => {
            return {
              sourceField: undefined,
              destinationField: objectField,
              expandedJson: false,
              jsonFields: [undefined],
            };
          })
      : [];
    setState((state) => ({ ...state, fieldMappings }));
  }

  const back = () => {
    if (state.step === SyncSetupStep.ExistingSources) {
      return navigate("/");
    }

    if (state.step === SyncSetupStep.ChooseSourceType && state.skippedSourceSelection) {
      return navigate("/");
    }

    let prevStep = state.step - 1;
    if (state.skippedSourceSetup && state.step === SyncSetupStep.ChooseData) {
      prevStep = SyncSetupStep.ExistingSources;
    }

    // clear errors here since it looks bad when they linger
    setState((state) => ({
      ...state,
      step: prevStep,
      error: undefined,
      newSourceState: { ...state.newSourceState, error: undefined },
    }));
  };

  return (
    <>
      <Header close={close} state={state} />
      <AppContent
        linkToken={linkToken}
        state={state}
        setState={setState}
        supportEmail={supportEmail}
        docsLink={docsLink}
      />
      <Footer back={back} linkToken={linkToken} state={state} setState={setState} />
    </>
  );
};

type AppContentProps = {
  linkToken: string;
  state: SetupSyncState;
  setState: React.Dispatch<React.SetStateAction<SetupSyncState>>;
} & FabraDisplayOptions;

const AppContent: React.FC<AppContentProps> = (props) => {
  const ref = useRef<HTMLDivElement>(null);
  // Scroll to the top on step change
  React.useEffect(() => {
    ref.current?.scrollTo(0, 0);
  }, [props.state.step]);

  let content: React.ReactNode;
  switch (props.state.step) {
    case SyncSetupStep.ExistingSources:
      content = <Sources linkToken={props.linkToken} state={props.state} setState={props.setState} />;
      break;
    case SyncSetupStep.ChooseSourceType:
      content = <WarehouseSelector linkToken={props.linkToken} state={props.state} setState={props.setState} />;
      break;
    case SyncSetupStep.ConnectionDetails:
      content = (
        <NewSourceConfiguration
          linkToken={props.linkToken}
          state={props.state}
          setState={props.setState}
          supportEmail={props.supportEmail}
          docsLink={props.docsLink}
        />
      );
      break;
    case SyncSetupStep.ChooseData:
      content = <ObjectSetup linkToken={props.linkToken} state={props.state} setState={props.setState} />;
      break;
    case SyncSetupStep.Finalize:
      content = <FinalizeSync linkToken={props.linkToken} state={props.state} setState={props.setState} />;
      break;
    default:
      // TODO: should never happen
      break;
  }

  return (
    <div
      ref={ref}
      className="tw-overflow-auto tw-w-full tw-h-full tw-flex tw-justify-center tw-pt-10 tw-bg-transparent"
    >
      {content}
    </div>
  );
};

const Header: React.FC<{ close: (() => void) | undefined; state: SetupSyncState }> = ({ close, state }) => {
  return (
    <div className="tw-flex tw-flex-row tw-items-center tw-w-full tw-h-20 tw-min-h-[80px] tw-border-b tw-border-slate-200">
      <div className="tw-flex tw-flex-row tw-gap-10 tw-justify-center tw-items-center tw-w-full">
        <StepBreadcrumb
          step={1}
          content="Select source"
          active={state.step <= SyncSetupStep.ChooseSourceType}
          complete={state.step > SyncSetupStep.ChooseSourceType}
        />
        <StepBreadcrumb
          step={2}
          content="Connect source"
          active={state.step === SyncSetupStep.ConnectionDetails}
          complete={state.step > SyncSetupStep.ConnectionDetails}
        />
        <StepBreadcrumb
          step={3}
          content="Define model"
          active={state.step === SyncSetupStep.ChooseData}
          complete={state.step > SyncSetupStep.ChooseData}
        />
        <StepBreadcrumb
          step={4}
          content="Finalize sync"
          active={state.step === SyncSetupStep.Finalize}
          complete={state.step > SyncSetupStep.Finalize}
        />
      </div>
      {close && (
        <button
          className="tw-absolute tw-flex tw-items-center t tw-right-10 tw-border-none tw-cursor-pointer tw-p-0"
          onClick={close}
        >
          <svg className="tw-h-6 tw-fill-slate-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="none">
            <path d="M5.1875 15.6875L4.3125 14.8125L9.125 10L4.3125 5.1875L5.1875 4.3125L10 9.125L14.8125 4.3125L15.6875 5.1875L10.875 10L15.6875 14.8125L14.8125 15.6875L10 10.875L5.1875 15.6875Z" />
          </svg>
        </button>
      )}
    </div>
  );
};

const StepBreadcrumb: React.FC<{ content: string; step: number; active: boolean; complete: boolean }> = ({
  step,
  content,
  active,
  complete,
}) => {
  return (
    <div className="tw-flex tw-flex-row tw-justify-center tw-items-center tw-select-none">
      <div
        className={classNames(
          "tw-rounded-md tw-h-[18px] tw-w-[18px] tw-flex tw-justify-center tw-items-center tw-text-[10px]",
          !active && !complete && "tw-bg-slate-200 tw-text-slate-900",
          active && "tw-bg-primary tw-text-primary-text",
          complete && "tw-bg-green-100 tw-text-green-800",
        )}
      >
        {complete ? <CheckIcon className="tw-h-3" /> : step}
      </div>
      <span className={classNames("tw-font-medium tw-pl-2", active && "tw-text-primary")}>{content}</span>
    </div>
  );
};

type FooterProps = {
  back: () => void;
  linkToken: string;
  state: SetupSyncState;
  setState: React.Dispatch<React.SetStateAction<SetupSyncState>>;
};

export const Footer: React.FC<FooterProps> = (props) => {
  let onClick = () => {};
  let continueText: string = "Continue";
  let showContinue = true;
  const createSync = useCreateNewSync();
  const showToast = useConnectShowToast();

  const createNewSourceMutation = useMutation(async () => {
    await createNewSource(props.linkToken, props.state, props.setState);
  });

  const createNewSyncMutation = useMutation(async () => {
    await createSync(props.linkToken, props.state, props.setState);
  });

  switch (props.state.step) {
    case SyncSetupStep.ExistingSources:
      showContinue = false;
      break;
    case SyncSetupStep.ChooseSourceType:
      showContinue = false;
      break;
    case SyncSetupStep.ConnectionDetails:
      onClick = () => {
        createNewSourceMutation.reset();
        createNewSourceMutation.mutate();
      };
      break;
    case SyncSetupStep.ChooseData:
      onClick = () => {
        if (validateObjectSetup(props.state, showToast)) {
          props.setState((state) => ({ ...state, step: props.state.step + 1 }));
        }
      };
      break;
    case SyncSetupStep.Finalize:
      continueText = "Create Sync";
      onClick = () => {
        createNewSyncMutation.reset();
        createNewSyncMutation.mutate();
      };
      break;
  }

  return (
    <div className="tw-w-full tw-min-h-[80px]">
      <div className="tw-flex tw-flex-row tw-w-full tw-h-full tw-px-20 tw-border-t tw-border-slate-200 tw-items-center tw-gap-x-2">
        <button
          className="tw-border tw-border-slate-300 tw-font-medium tw-rounded-md tw-w-32 tw-h-10 tw-select-none hover:tw-bg-slate-100"
          onClick={props.back}
        >
          Back
        </button>
        {showContinue && (
          <Button onClick={onClick} className="tw-border tw-w-36 tw-h-10 tw-ml-auto tw-select-none">
            {createNewSourceMutation.isLoading || createNewSyncMutation.isLoading ? <Loading light /> : continueText}
          </Button>
        )}
      </div>
      <div className="tw-flex tw-justify-end tw-mt-1">
        {createNewSourceMutation.error && props.state.step === SyncSetupStep.ConnectionDetails && (
          <ErrorDisplay error={createNewSourceMutation.error} className="tw-text-red-500" />
        )}
        {createNewSyncMutation.error && props.state.step === SyncSetupStep.Finalize && (
          <ErrorDisplay error={createNewSyncMutation.error} className="tw-text-red-500" />
        )}
      </div>
    </div>
  );
};
