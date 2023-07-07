import { zodResolver } from "@hookform/resolvers/zod";
import { ChangeEvent } from "react";
import { Controller, useForm } from "react-hook-form";
import { Button } from "src/components/button/Button";
import { FormError } from "src/components/FormError";
import { InfoIcon } from "src/components/icons/Icons";
import { InputStyle } from "src/components/input/Input";
import { DestinationSelector, NamespaceSelector, TableSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import {
  DestinationSetupFormSchema,
  DestinationSetupFormType,
  objectTargetOptions,
  SupportedConnectionType,
  SUPPORTED_CONNECTION_TYPES,
} from "src/pages/objects/helpers";
import { Connection, ConnectionType, Destination, TargetType } from "src/rpc/api";
import { mergeClasses } from "src/utils/twmerge";

interface DestinationSetupProps {
  initialFormState?: DestinationSetupFormType;
  onComplete: (values: DestinationSetupFormType) => void;
  isUpdate?: boolean;
}

type UseFormReturn = ReturnType<typeof useForm<DestinationSetupFormType>>;
type Control = UseFormReturn["control"];
type Errors = UseFormReturn["formState"]["errors"];

function initializeFormState(initial: DestinationSetupProps["initialFormState"]): Partial<DestinationSetupFormType> {
  const connectionType = initial?.connectionType;
  if (connectionType && !SUPPORTED_CONNECTION_TYPES.includes(connectionType as SupportedConnectionType)) {
    return {};
  }

  if (connectionType === ConnectionType.BigQuery) {
    return {
      displayName: initial?.displayName,
      destination: initial?.destination,
      connectionType,
      targetType: TargetType.SingleExisting,
      namespace: initial?.namespace,
      tableName: initial?.tableName,
    };
  } else if (connectionType === ConnectionType.DynamoDb) {
    return {
      displayName: initial?.displayName,
      destination: initial?.destination,
      tableName: initial?.tableName,
      connectionType,
      targetType: TargetType.SingleExisting,
    };
  } else {
    return {
      displayName: initial?.displayName,
      destination: initial?.destination,
      connectionType: connectionType as ConnectionType.Webhook | undefined,
      targetType: TargetType.Webhook,
    };
  }
}

export const DestinationSetup: React.FC<DestinationSetupProps> = ({ isUpdate, onComplete, initialFormState }) => {
  const form = useForm<DestinationSetupFormType>({
    resolver: async (data, context, options) => {
      const errors = await zodResolver(DestinationSetupFormSchema)(data, context, options);
      return errors;
    },
    defaultValues: initializeFormState(initialFormState),
  });
  const {
    handleSubmit,
    register,
    formState: { errors },
    control,
    watch,
  } = form;

  const onSubmit = handleSubmit((values) => {
    onComplete(values);
  });

  const connectionType = watch("connectionType");
  const destination = watch("destination");
  const namespace = watch("namespace");

  return (
    <div>
      <h2 className="tw-mb-1 tw-font-bold tw-text-xl tw-text-center">{isUpdate ? "Update Object" : "New Object"}</h2>
      <p className="tw-text-center tw-mb-3">
        {isUpdate ? "Change your object's configuration." : "Enter your object configuration."}
      </p>
      <form onSubmit={onSubmit} className="tw-space-y-4">
        <div>
          <label className="tw-flex tw-items-center">
            <span className="tw-font-medium">Display Name</span>
            <Tooltip placement="right" label="Pick a name for this object that your customers will see.">
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </label>
          <input autoFocus className={InputStyle} {...register("displayName")} placeholder="My Object"></input>
          <FormError message={errors.displayName?.message} className="mt-1" />
        </div>

        <div>
          <label className="tw-font-medium">Destination</label>
          <Controller
            name="destination"
            control={control}
            render={({ field }) => {
              return (
                <DestinationSelector
                  className="tw-mt-0"
                  destination={field.value}
                  disabled={isUpdate}
                  setDestination={(d) => {
                    const connectionType = d.connection.connection_type;
                    field.onChange(d);
                    if (connectionType === ConnectionType.Webhook) {
                      form.setValue("targetType", TargetType.Webhook);
                    } else {
                      form.setValue("targetType", TargetType.SingleExisting);
                    }
                    form.setValue("connectionType", connectionType as SupportedConnectionType);
                    // Retrigger form validation for these fields since we're setting them manually.
                    form.trigger("targetType");
                    form.trigger("connectionType");
                  }}
                />
              );
            }}
          />
          <FormError message={errors.destination?.message || errors.connectionType?.message} />
        </div>
        {(connectionType === ConnectionType.BigQuery || connectionType === ConnectionType.DynamoDb) && (
          <ObjectTargetFieldset control={control} errors={errors} disabled={isUpdate} />
        )}
        {connectionType === ConnectionType.BigQuery && (
          <NamespaceField control={control} destination={destination} errors={errors} isUpdate={isUpdate} />
        )}
        {(connectionType === ConnectionType.BigQuery || connectionType === ConnectionType.DynamoDb) && (
          <TableField
            control={control}
            connection={destination.connection}
            errors={errors}
            isUpdate={isUpdate}
            namespace={namespace}
          />
        )}

        <Button type="submit" className="tw-w-full tw-py-2">
          Continue
        </Button>
      </form>
    </div>
  );
};

export function NamespaceField({
  control,
  destination,
  errors,
  isUpdate,
}: {
  control: Control;
  destination: Destination;
  errors: Errors;
  isUpdate?: boolean;
}) {
  return (
    <div>
      <label className="tw-font-medium">Namespace</label>
      <Controller
        name="namespace"
        control={control}
        render={({ field }) => {
          return (
            <NamespaceSelector
              className="tw-mt-0"
              validated={true}
              connection={destination.connection}
              namespace={field.value}
              disabled={isUpdate}
              setNamespace={field.onChange}
              noOptionsString="No Namespaces Available! (Choose a data source)"
            />
          );
        }}
      />
      <FormError message={errors.namespace?.message} className="mt-1" />
    </div>
  );
}

function TableField({
  namespace,
  control,
  connection,
  errors,
  isUpdate,
}: {
  namespace?: string | undefined;
  control: Control;
  connection: Connection;
  errors: Errors;
  isUpdate?: boolean;
}) {
  return (
    <div>
      <label className="tw-font-medium">Table</label>
      <Controller
        name="tableName"
        control={control}
        render={({ field }) => {
          return (
            <TableSelector
              className="tw-mt-0"
              connection={connection}
              namespace={namespace}
              tableName={field.value}
              disabled={isUpdate}
              setTableName={field.onChange}
              noOptionsString="No Tables Available! (Choose a namespace)"
              validated={true}
            />
          );
        }}
      />
      <FormError message={errors.tableName?.message} className="mt-1" />
    </div>
  );
}

function ObjectTargetFieldset({ control, errors, disabled }: { control: Control; errors: Errors; disabled?: boolean }) {
  return (
    <fieldset className="tw-my-4">
      <label className="tw-font-medium">Target</label>
      <p className="tw-text-slate-600">Where should Fabra load the data in your destination?</p>
      <legend className="tw-sr-only">Target</legend>
      <div className="tw-space-y-4 tw-mt-2">
        {objectTargetOptions.map((target) => (
          <Controller
            key={target.type}
            name="targetType"
            control={control}
            render={({ field }) => {
              return (
                <div key={String(target.type)} className="tw-flex tw-items-center">
                  <input
                    id={String(target.type)}
                    name="target"
                    type="radio"
                    value={target.type}
                    checked={field.value === target.type}
                    onChange={(e: ChangeEvent<HTMLInputElement>) => {
                      field.onChange(target.type);
                    }}
                    disabled={disabled}
                    className={mergeClasses(
                      "tw-h-4 tw-w-4 tw-border-slate-300 tw-text-indigo-600 focus:tw-ring-indigo-600 tw-cursor-pointer",
                      disabled ? "tw-cursor-not-allowed" : "tw-cursor-pointer",
                    )}
                  />
                  <div className="tw-flex tw-flex-row tw-items-center tw-ml-3 tw-leading-6">
                    <label htmlFor={String(target.type)} className="tw-text-sm tw-cursor-pointer">
                      {target.title}
                    </label>
                    <Tooltip label={target.description} placement="top-start">
                      <InfoIcon className="tw-ml-1.5 tw-h-3 tw-fill-slate-400" />
                    </Tooltip>
                  </div>
                </div>
              );
            }}
          />
        ))}
      </div>
      <FormError message={errors.targetType?.message} className="mt-1" />
    </fieldset>
  );
}
