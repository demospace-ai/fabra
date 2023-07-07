import { zodResolver } from "@hookform/resolvers/zod";
import { Checkbox } from "src/components/checkbox/Checkbox";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { FormError } from "src/components/FormError";
import { Button } from "src/components/button/Button";
import { Input } from "src/components/input/Input";
import { Loading } from "src/components/loading/Loading";
import { DestinationSetupFormType, ObjectFieldsFormType, ObjectFieldsSchema } from "src/pages/objects/helpers";
import { ObjectFieldInput } from "src/rpc/api";
import { useSchema } from "src/rpc/data";
import { mergeClasses } from "src/utils/twmerge";

interface ExistingObjectFieldsProps {
  isUpdate?: boolean;
  destinationSetupData: DestinationSetupFormType;
  initialFormState?: ObjectFieldsFormType;
  onComplete: (values: ObjectFieldsFormType) => void;
}

export const ExistingObjectFields: React.FC<ExistingObjectFieldsProps> = ({
  destinationSetupData,
  isUpdate = false,
  onComplete,
  initialFormState,
}) => {
  const schemaQuery = useSchema(
    destinationSetupData.destination.connection.id,
    "namespace" in destinationSetupData ? destinationSetupData.namespace : undefined,
    "tableName" in destinationSetupData ? destinationSetupData.tableName : undefined,
  );
  const initialObjectFields = initialFormState?.objectFields ?? [];

  return (
    <div>
      <div className="tw-w-full tw-text-center tw-mb-2 tw-font-bold tw-text-lg">
        {isUpdate ? "Update Object Fields" : "Object Fields"}
      </div>
      <div className="tw-text-center tw-mb-3">Provide customer-facing names and descriptions for each field.</div>
      {schemaQuery.loading && !schemaQuery.schema ? (
        <div className="tw-text-center">
          <h2 className="tw-mb-2">Loading fields...</h2>
          <Loading />
        </div>
      ) : schemaQuery.error ? (
        <div className="tw-text-center">
          <h2 className="tw-mb-2">Something went wrong</h2>
        </div>
      ) : (
        <ExistingObjectFieldsForm
          isUpdate={isUpdate}
          objectFields={
            initialObjectFields.length > 0
              ? initialObjectFields
              : schemaQuery.schema?.map((field) => {
                  return {
                    name: field.name,
                    type: field.type,
                    omit: false,
                    optional: false,
                    id: undefined,
                  };
                }) || []
          }
          onComplete={onComplete}
        />
      )}
    </div>
  );
};

const ExistingObjectFieldsForm: React.FC<{
  objectFields: ObjectFieldInput[];
  onComplete: (values: ObjectFieldsFormType) => void;
  isUpdate: boolean;
}> = ({ objectFields, isUpdate, onComplete }) => {
  const {
    control,
    formState: { errors },
    handleSubmit,
  } = useForm<ObjectFieldsFormType>({
    resolver: zodResolver(ObjectFieldsSchema),
    defaultValues: {
      objectFields,
    },
  });

  const { fields } = useFieldArray({
    name: "objectFields",
    control,
  });

  const onSubmit = handleSubmit((values) => {
    onComplete(values);
  });

  return (
    <form onSubmit={onSubmit} className="tw-w-full">
      <ul>
        {fields.map((objectField, i) => {
          return (
            <li key={objectField.id}>
              <div className={mergeClasses("tw-mt-5 tw-mb-7 tw-text-left")}>
                <h3 className="tw-text-base tw-font-semibold">{objectField.name}</h3>
                <div className="tw-flex tw-items-center tw-mt-2 tw-pb-1.5">
                  <span className="">Omit?</span>
                  <Controller
                    name={`objectFields.${i}.omit`}
                    control={control}
                    render={({ field }) => (
                      <Checkbox
                        className="tw-ml-2 tw-h-4 tw-w-4"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={isUpdate}
                      />
                    )}
                  />
                  <span className="tw-ml-4">Optional?</span>
                  <Controller
                    name={`objectFields.${i}.optional`}
                    control={control}
                    render={({ field }) => (
                      <Checkbox
                        className="tw-ml-2 tw-h-4 tw-w-4"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={isUpdate}
                      />
                    )}
                  />
                </div>
                <Controller
                  name={`objectFields.${i}.display_name`}
                  control={control}
                  render={({ field }) => (
                    <Input
                      value={field.value}
                      setValue={field.onChange}
                      placeholder="Display name (optional)"
                      className="tw-mb-2"
                    />
                  )}
                />
                <Controller
                  name={`objectFields.${i}.description`}
                  control={control}
                  render={({ field }) => (
                    <Input
                      value={field.value}
                      setValue={field.onChange}
                      placeholder="Description (optional)"
                      className="tw-mb-2"
                    />
                  )}
                />
              </div>
            </li>
          );
        })}
      </ul>
      <Button type="submit" className="tw-mt-6 tw-w-100 tw-h-10">
        Continue
      </Button>
      <FormError message={errors.objectFields?.message} />
    </form>
  );
};
