import { zodResolver } from "@hookform/resolvers/zod";
import { ChangeEvent } from "react";
import { Controller, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { Button } from "src/components/button/Button";
import { Checkbox } from "src/components/checkbox/Checkbox";
import { FormError } from "src/components/FormError";
import { InfoIcon } from "src/components/icons/Icons";
import { ValidatedComboInput, ValidatedDropdownInput, ValidatedInput } from "src/components/input/Input";
import { Loading } from "src/components/loading/Loading";
import { useShowToast } from "src/components/notifications/Notifications";
import { FieldSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import {
  createNewObject,
  DestinationSetupFormType,
  FinalizeObjectFormSchema,
  FinalizeObjectFormType,
  ObjectFieldsFormType,
  updateObject,
} from "src/pages/objects/helpers";
import { ConnectionType, FabraObject, Field, FrequencyUnits, GetObjects, ObjectField, SyncMode } from "src/rpc/api";
import { useMutation } from "src/utils/queryHelpers";
import { mutate } from "swr";

type FinalizeStepProps = {
  existingObject?: FabraObject;
  objectFields: ObjectFieldsFormType;
  destinationSetup: DestinationSetupFormType;
  isUpdate: boolean;
  initialFormState?: FinalizeObjectFormType;
  onComplete: () => void;
};

export const Finalize: React.FC<FinalizeStepProps> = ({
  objectFields,
  existingObject,
  isUpdate,
  destinationSetup,
  initialFormState,
}) => {
  const {
    control,
    formState: { errors },
    watch,
    setValue,
    setError,
    clearErrors,
    handleSubmit,
  } = useForm<FinalizeObjectFormType>({
    resolver: zodResolver(FinalizeObjectFormSchema),
    defaultValues: {
      ...initialFormState,
      connectionType: destinationSetup.connectionType,
      endCustomerIdField: initialFormState?.endCustomerIdField ? initialFormState.endCustomerIdField : undefined,
    },
  });
  const showToast = useShowToast();
  const navigate = useNavigate();
  const connectionType = destinationSetup.destination.connection.connection_type;

  const saveConfigurationMutation = useMutation(
    async (values) => {
      if (existingObject) {
        return await updateObject({ existingObject, objectFields, destinationSetup, finalizeValues: values });
      } else {
        return await createNewObject({ objectFields, destinationSetup, finalizeValues: values });
      }
    },
    {
      onSuccess: (object) => {
        showToast("success", isUpdate ? "Successfully updated object!" : "Successfully created object!", 4000);
        navigate(`/objects/${object.id}`);
        mutate({ GetObjects: GetObjects });
      },
      onError: (e) => {
        showToast("error", isUpdate ? "Failed to update object." : "Failed to create object.", 4000);
        // Sets the form-level error message.
        setError("root.createObject", { message: e.message });
      },
    },
  );

  const fields: Field[] = objectFields.objectFields
    .filter((field): field is ObjectField => Boolean(field.name && field.field_type && !field.omit && !field.optional))
    .map((field) => {
      return { name: field.name, type: field.field_type };
    });

  const syncMode = watch("syncMode");
  const recurring = watch("recurring");
  let recommendedCursor = <></>;
  switch (syncMode) {
    case SyncMode.IncrementalAppend:
      recommendedCursor = (
        <>
          For <span className="tw-px-1 tw-bg-black tw-font-mono">incremental_append</span> syncs, you should use an{" "}
          <span className="tw-px-1 tw-bg-black tw-font-mono">created_at</span> field.
        </>
      );
      break;
    case SyncMode.IncrementalUpdate:
      recommendedCursor = (
        <>
          For <span className="tw-px-1 tw-bg-black tw-font-mono">incremental_update</span> syncs, you should use an{" "}
          <span className="tw-px-1 tw-bg-black tw-font-mono">updated_at</span> field.
        </>
      );
      break;
    case SyncMode.FullOverwrite:
      break;
  }

  return (
    <form
      className="tw-flex tw-flex-col tw-w-100"
      onSubmit={handleSubmit((values) => {
        clearErrors("root.createObject");
        saveConfigurationMutation.mutate(values);
      })}
    >
      <div className="tw-w-full tw-text-center tw-mb-2 tw-font-bold tw-text-lg">Object Settings</div>
      <div className="tw-text-center tw-mb-3">Enter default settings for object syncs.</div>
      <Controller
        control={control}
        name="syncMode"
        render={({ field }) => (
          <>
            <SyncModeSelector
              value={field.value}
              onChange={(e) => {
                const value = e.target.value;
                if (value === SyncMode.FullOverwrite) {
                  setValue("cursorField", undefined);
                }
                field.onChange(value);
              }}
              disabled={isUpdate}
            />
          </>
        )}
      />
      <FormError message={errors.syncMode?.message} />
      {[SyncMode.IncrementalAppend, SyncMode.IncrementalUpdate].includes(syncMode) && (
        <>
          <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5 tw-mb-3">
            <span className="tw-font-medium">Cursor Field</span>
            <Tooltip
              placement="right"
              label={
                <>
                  Cursor field is usually a timestamp. This lets Fabra know what data has changed since the last sync.{" "}
                  {recommendedCursor}
                </>
              }
              maxWidth={400}
              interactive
            >
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
          <Controller
            control={control}
            name="cursorField"
            render={({ field }) => (
              <FieldSelector
                className="tw-mt-0 tw-w-100"
                field={field.value}
                setField={field.onChange}
                placeholder="Cursor Field"
                label="Cursor Field"
                noOptionsString="No Fields Available!"
                predefinedFields={fields}
                validated={true}
                valid={!errors.cursorField}
                disabled={!!existingObject}
              />
            )}
          />
          <FormError message={errors.cursorField?.message} />
        </>
      )}
      {[SyncMode.IncrementalUpdate].includes(syncMode) && (
        <>
          <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5 tw-mb-3">
            <span className="tw-font-medium">Primary Key</span>
            <Tooltip
              placement="right"
              label="Primary key is usually an ID field. This lets Fabra know which existing rows in the target to update when they change."
              maxWidth={400}
            >
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
          <Controller
            name="primaryKey"
            control={control}
            render={({ field }) => (
              <>
                <FieldSelector
                  className="tw-mt-0 tw-w-100"
                  field={field.value}
                  setField={field.onChange}
                  placeholder="Primary Key"
                  noOptionsString="No Fields Available!"
                  validated={true}
                  predefinedFields={fields}
                  disabled={!!existingObject}
                />
              </>
            )}
          />
          <FormError message={errors.primaryKey?.message} />
        </>
      )}
      {syncMode && (
        <>
          {connectionType === ConnectionType.DynamoDb && (
            <>
              <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5">
                <span className="tw-font-medium">End Customer ID</span>
              </div>
              <Controller
                name="endCustomerIdField"
                control={control}
                render={({ field }) => (
                  <ValidatedComboInput
                    className="tw-mt-3"
                    loading={false}
                    validated={true}
                    disabled={!!existingObject}
                    options={fields}
                    selected={field.value}
                    setSelected={field.onChange}
                    getElementForDisplay={(value: Field) => value.name}
                    noOptionsString={"No field available!"}
                    placeholder={"Choose field"}
                  />
                )}
              />
              <FormError message={errors.endCustomerIdField?.message} />
            </>
          )}
          {connectionType !== ConnectionType.Webhook && connectionType !== ConnectionType.DynamoDb && (
            <>
              <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5">
                <span className="tw-font-medium">End Customer ID</span>
              </div>
              <Controller
                control={control}
                name="endCustomerIdField"
                render={({ field }) => (
                  <FieldSelector
                    className="tw-mt-0 tw-w-100"
                    field={field.value}
                    setField={field.onChange}
                    placeholder="End Customer ID Field"
                    noOptionsString="No Fields Available!"
                    validated={true}
                    connection={destinationSetup.destination.connection}
                    namespace={"namespace" in destinationSetup ? destinationSetup.namespace : undefined}
                    tableName={"tableName" in destinationSetup ? destinationSetup.tableName : undefined}
                    disabled={!!existingObject}
                  />
                )}
              />
              <FormError message={errors.endCustomerIdField?.message} />
            </>
          )}
          <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-6">
            <span className="tw-font-medium">Recurring?</span>
            <Controller
              control={control}
              name="recurring"
              defaultValue={false}
              render={({ field }) => (
                <Checkbox
                  className="tw-ml-2 tw-h-4 tw-w-4"
                  checked={field.value}
                  onCheckedChange={() => field.onChange(!field.value)}
                />
              )}
            />
          </div>
          {recurring && (
            <>
              <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5 tw-mb-3">
                <span className="tw-font-medium">Frequency</span>
              </div>
              <Controller
                control={control}
                name="frequency"
                render={({ field }) => (
                  <ValidatedInput
                    id="frequency"
                    className="tw-w-100"
                    type="number"
                    value={field.value}
                    setValue={field.onChange}
                    placeholder="Sync Frequency"
                  />
                )}
              />
              <FormError message={errors.frequency?.message} />

              <div className="tw-w-full tw-flex tw-flex-row tw-items-center tw-mt-5 tw-mb-3">
                <span className="tw-font-medium">Frequency Units</span>
              </div>
              <Controller
                control={control}
                name="frequencyUnits"
                render={({ field }) => (
                  <ValidatedDropdownInput
                    className="tw-mt-0 tw-w-100"
                    options={Object.values(FrequencyUnits)}
                    selected={field.value}
                    setSelected={field.onChange}
                    loading={false}
                    placeholder="Frequency Units"
                    noOptionsString="nil"
                    getElementForDisplay={(value) => {
                      return value.charAt(0).toUpperCase() + value.slice(1);
                    }}
                  />
                )}
              />
              <FormError message={errors.frequencyUnits?.message} />
            </>
          )}
        </>
      )}
      <Button type="submit" className="tw-mt-10 tw-w-full tw-h-10">
        {saveConfigurationMutation.isLoading ? <Loading /> : existingObject ? "Update Object" : "Create Object"}
      </Button>
      <FormError message={errors.root?.createObject?.message} />
      <FormError message={errors.connectionType?.message} />
    </form>
  );
};

export const SyncModeSelector: React.FC<{
  value: SyncMode;
  onChange: (e: ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
}> = ({ value, onChange, disabled = false }) => {
  type SyncModeOption = {
    mode: SyncMode;
    title: string;
    description: string;
  };
  const syncModes: SyncModeOption[] = [
    {
      mode: SyncMode.FullOverwrite,
      title: "Full Overwrite",
      description: "Fabra will overwrite the entire target table on every sync.",
    },
    {
      mode: SyncMode.IncrementalAppend,
      title: "Incremental Append",
      description: "Fabra will append any new rows since the last sync to the existing target table.",
    },
    // TODO
    // {
    //   mode: SyncMode.IncrementalUpdate,
    //   title: "Incremental Update",
    //   description: "Fabra will add new rows and update any modified rows since the last sync."
    // },
  ];
  return (
    <div className="tw-mt-5">
      <label className="tw-font-medium">Sync Mode</label>
      <p className="tw-text-slate-600">How should Fabra load the data in your destination?</p>
      <fieldset className="tw-mt-4">
        <legend className="tw-sr-only">Sync Mode</legend>
        <div className="tw-space-y-4 tw-flex tw-flex-col">
          {syncModes.map((syncMode) => (
            <div key={String(syncMode.mode)} className="tw-flex tw-items-center">
              <input
                id={String(syncMode.mode)}
                name="syncmode"
                type="radio"
                disabled={disabled}
                checked={value === syncMode.mode}
                value={syncMode.mode}
                onChange={onChange}
                className="tw-h-4 tw-w-4 tw-border-slate-300 tw-text-indigo-600 focus:tw-ring-indigo-600 tw-cursor-pointer"
              />
              <div className="tw-flex tw-flex-row tw-items-center tw-ml-3 tw-leading-6">
                <label htmlFor={String(syncMode.mode)} className="tw-text-sm tw-cursor-pointer">
                  {syncMode.title}
                </label>
                <Tooltip label={syncMode.description} placement="top-start">
                  <InfoIcon className="tw-ml-1.5 tw-h-3 tw-fill-slate-400" />
                </Tooltip>
              </div>
            </div>
          ))}
        </div>
      </fieldset>
    </div>
  );
};
