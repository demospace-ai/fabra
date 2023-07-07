import { PlusCircleIcon } from "@heroicons/react/24/outline";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { Button, DeleteButton } from "src/components/button/Button";
import { Checkbox } from "src/components/checkbox/Checkbox";
import { FormError } from "src/components/FormError";
import { InfoIcon } from "src/components/icons/Icons";
import { Input } from "src/components/input/Input";
import { FieldTypeSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { ObjectFieldsFormType, ObjectFieldsSchema } from "src/pages/objects/helpers";

interface NewObjectFieldsProps {
  isUpdate?: boolean;
  initialFormState?: ObjectFieldsFormType;
  onComplete: (values: ObjectFieldsFormType) => void;
}

export function NewObjectFields({ initialFormState, onComplete, isUpdate = false }: NewObjectFieldsProps) {
  const {
    formState: { errors },
    control,
    handleSubmit,
  } = useForm<ObjectFieldsFormType>({
    resolver: zodResolver(ObjectFieldsSchema),
    defaultValues: initialFormState,
  });
  const { fields, append, remove } = useFieldArray({
    name: "objectFields",
    control,
  });

  return (
    <form
      className="tw-h-full tw-w-full tw-text-center"
      onSubmit={handleSubmit((values) => {
        onComplete(values);
      })}
    >
      <div className="tw-w-full tw-text-center tw-mb-2 tw-font-bold tw-text-lg">
        {isUpdate ? "Update Object Fields" : "Create Object Fields"}
      </div>
      <div className="tw-text-center tw-mb-3">Provide customer-facing names and descriptions for each field.</div>
      <div className="tw-w-full tw-px-24">
        <div>
          {fields.map((objectField, i) => (
            <div key={objectField.id} className="tw-mt-5 tw-mb-7 tw-text-left tw-p-4 tw-border tw-rounded-lg">
              <div className="tw-flex tw-items-center">
                <span className="tw-font-semibold tw-text-lg tw-grow">Field {i + 1}</span>
                <DeleteButton
                  className="tw-ml-auto tw-stroke-red-400 tw-p-2"
                  onClick={() => remove(i)}
                  disabled={isUpdate}
                />
              </div>
              <div className="tw-flex tw-items-center tw-mt-3">
                <span>Optional?</span>
                <Controller
                  name={`objectFields.${i}.optional`}
                  control={control}
                  render={({ field }) => (
                    <Checkbox
                      className="tw-ml-2 tw-h-4 tw-w-4"
                      checked={field.value}
                      disabled={isUpdate}
                      onCheckedChange={field.onChange}
                    />
                  )}
                />
              </div>
              <div className="tw-flex tw-w-full tw-items-center tw-mb-2">
                <div className="tw-w-full tw-mr-4">
                  <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
                    <span>Field Key</span>
                    <Tooltip
                      placement="right"
                      label="Choose a valid JSON key that will be used when sending this field to your webhook."
                    >
                      <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
                    </Tooltip>
                  </div>
                  <Controller
                    name={`objectFields.${i}.name`}
                    control={control}
                    defaultValue={objectField.name}
                    render={({ field }) => (
                      <Input
                        value={field.value}
                        disabled={isUpdate}
                        setValue={field.onChange}
                        placeholder="Field Key"
                      />
                    )}
                  />
                </div>
                <div>
                  <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
                    <span>Field Type</span>
                    <Tooltip placement="right" label="Choose the type for this field.">
                      <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
                    </Tooltip>
                  </div>
                  <Controller
                    name={`objectFields.${i}.field_type`}
                    control={control}
                    defaultValue={objectField.field_type}
                    render={({ field }) => (
                      <FieldTypeSelector
                        className="tw-w-48 tw-m-0"
                        disabled={isUpdate}
                        type={field.value}
                        setFieldType={(value) => field.onChange(value)}
                      />
                    )}
                  />
                </div>
              </div>
              <div className="tw-flex">
                <FormError className="tw-w-[405px]" message={errors.objectFields?.[i]?.name?.message} />
                <FormError message={errors.objectFields?.[i]?.field_type?.message} />
              </div>
              <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
                <span>Display Name</span>
                <Tooltip
                  placement="right"
                  label="Set a customer-facing name that your customers will see when setting up a sync."
                >
                  <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
                </Tooltip>
              </div>
              <Controller
                name={`objectFields.${i}.display_name`}
                control={control}
                defaultValue={objectField.display_name}
                render={({ field }) => (
                  <Input
                    className="tw-mb-2"
                    value={field.value}
                    setValue={field.onChange}
                    placeholder="Display Name (optional)"
                  />
                )}
              />
              <div className="tw-flex tw-flex-row tw-items-center tw-mt-2 tw-mb-1">
                <span>Description</span>
                <Tooltip
                  placement="right"
                  label="Add any extra information that will help your customers understand how to map their data to this object."
                >
                  <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
                </Tooltip>
              </div>
              <Controller
                name={`objectFields.${i}.description`}
                control={control}
                defaultValue={objectField.description}
                render={({ field }) => (
                  <Input
                    className="tw-mb-2"
                    value={field.value}
                    setValue={field.onChange}
                    placeholder="Description (optional)"
                  />
                )}
              />
            </div>
          ))}
          {/* No adding/removing fields on existing objects since this may break syncs */}
          {!isUpdate && (
            <Button
              className="tw-mt-7 tw-mx-auto tw-flex tw-items-center tw-mb-8"
              onClick={() =>
                append({
                  name: "",
                  field_type: undefined,
                  omit: false,
                  optional: false,
                })
              }
            >
              <PlusCircleIcon className="tw-h-5 tw-mr-1.5 tw-stroke-2" />
              Add Object Field
            </Button>
          )}
        </div>
      </div>
      <Button type="submit" className="tw-mt-8 tw-w-100 tw-h-10">
        Continue
      </Button>
      <FormError message={errors.objectFields?.message} />
    </form>
  );
}
