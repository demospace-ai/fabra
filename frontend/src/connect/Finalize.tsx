import { XMarkIcon } from "@heroicons/react/24/outline";
import React from "react";
import { Button } from "src/components/button/Button";
import { InfoIcon } from "src/components/icons/Icons";
import { ValidatedInput } from "src/components/input/Input";
import { Loading } from "src/components/loading/Loading";
import { LinkFieldSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { FieldMappingState, SetupSyncProps } from "src/connect/state";
import { Field, FieldType, ObjectField } from "src/rpc/api";
import { useObject } from "src/rpc/data";

export const FinalizeSync: React.FC<SetupSyncProps> = (props) => {
  return (
    <div className="tw-w-full tw-pl-20 tw-pr-[72px] tw-flex tw-flex-col">
      <div className="tw-text-left tw-mb-5 tw-text-2xl tw-font-bold tw-text-slate-900">
        Finalize your sync configuration
      </div>
      <div className="tw-w-[100%] tw-min-w-[400px] tw-h-full">
        <div className="tw-text-base tw-font-medium tw-mb-1 tw-text-slate-800">Display Name</div>
        <div className="tw-text-slate-600 tw-mb-4">Choose a name to help you identify this sync in the future.</div>
        <ValidatedInput
          id="display_name"
          className="tw-w-96"
          value={props.state.displayName}
          setValue={(value) => props.setState((state) => ({ ...props.state, displayName: value }))}
          placeholder="Display Name"
        />
        <div className="tw-text-base tw-font-medium tw-mt-9 tw-mb-1 tw-text-slate-800">Field Mapping</div>
        <div className="tw-text-slate-600 tw-mb-4">
          This is how your data will be mapped to the fields in the application.
        </div>
        <FieldMappings linkToken={props.linkToken} state={props.state} setState={props.setState} />
        {/* 
        <div className="tw-text-base tw-font-medium tw-mt-12 tw-mb-1 tw-text-slate-800">Additional Configuration</div>
        <div className="tw-text-slate-600">Configure additional settings for the sync.</div>
        <ValidatedInput id="frequency" className="tw-w-96" min={props.state.frequencyUnits === FrequencyUnits.Minutes ? 30 : 1} type="number" value={props.state.frequency} setValue={value => props.setState({ ...props.state, frequency: value })} placeholder="Sync Frequency" label="Sync Frequency" />
        <ValidatedDropdownInput className="tw-w-96" options={Object.values(FrequencyUnits)} selected={props.state.frequencyUnits} setSelected={value => props.setState({ ...props.state, frequencyUnits: value })} loading={false} placeholder="Frequency Units" noOptionsString="nil" label="Frequency Unit" getElementForDisplay={(value) => value.charAt(0).toUpperCase() + value.slice(1)} />
        */}
        {props.state.error && (
          <div className="tw-mt-4 tw-text-red-700 tw-p-2 tw-text-center tw-bg-red-50 tw-border tw-border-red-600 tw-rounded tw-max-w-3xl ">
            {props.state.error}
          </div>
        )}
        <div className="tw-pb-52"></div>
      </div>
    </div>
  );
};

const FieldMappings: React.FC<SetupSyncProps> = (props) => {
  const { object } = useObject(props.state.object?.id, props.linkToken);
  if (!object || !props.state.fieldMappings) {
    return <Loading />;
  }

  const updateFieldMapping = (newFieldMapping: FieldMappingState, index: number) => {
    props.setState((state) => {
      if (!state.fieldMappings) {
        // TODO: should not happen
        return state;
      }

      return {
        ...state,
        fieldMappings: state.fieldMappings.map((original, i) => {
          if (i === index) {
            return newFieldMapping;
          } else {
            return original;
          }
        }),
      };
    });
  };

  const updateJsonField = (newJsonField: Field, fieldMappingIdx: number, jsonIdx: number) => {
    const fieldMapping = props.state.fieldMappings![fieldMappingIdx];
    fieldMapping.jsonFields[jsonIdx] = newJsonField;
    updateFieldMapping(fieldMapping, fieldMappingIdx);
  };

  const removeJsonField = (fieldMappingIdx: number, jsonIdx: number) => {
    const fieldMapping = props.state.fieldMappings![fieldMappingIdx];
    fieldMapping.jsonFields.splice(jsonIdx, 1);
    updateFieldMapping(fieldMapping, fieldMappingIdx);
  };

  const addJsonField = (fieldMapping: FieldMappingState, fieldMappingIdx: number) => {
    updateFieldMapping(
      {
        ...fieldMapping,
        jsonFields: [...fieldMapping.jsonFields, undefined],
      },
      fieldMappingIdx!,
    );
  };

  return (
    <div className="tw-border tw-border-slate-200 tw-rounded-lg tw-max-w-3xl tw-divide-y tw-overflow-hidden">
      {object.object_fields.map((objectField) => {
        let fieldMappingIdx = props.state.fieldMappings?.findIndex(
          (fieldMapping) => fieldMapping.destinationField.id === objectField.id,
        );
        const fieldMapping = props.state.fieldMappings![fieldMappingIdx!];
        if (objectField.omit) {
          return;
        }

        if (fieldMapping.expandedJson) {
          return (
            <div className="tw-p-3 tw-flex tw-flex-col">
              <div key={objectField.name} className="tw-flex tw-flex-row tw-items-top tw-justify-between">
                <MappedField
                  objectField={objectField}
                  fieldMapping={fieldMapping}
                  fieldMappingIdx={fieldMappingIdx!}
                  updateFieldMapping={updateFieldMapping}
                />
                <div className="tw-flex tw-flex-col tw-gap-4">
                  {fieldMapping.jsonFields.map((jsonField, jsonIdx) => {
                    return (
                      <div className="tw-flex tw-items-center">
                        <LinkFieldSelector
                          className="tw-mt-0 tw-w-[360px] tw-flex"
                          field={jsonField}
                          setField={(value: Field) => {
                            updateJsonField(value, fieldMappingIdx!, jsonIdx);
                          }}
                          placeholder="Choose a field"
                          noOptionsString="No Fields Available!"
                          validated={true}
                          source={props.state.source}
                          namespace={props.state.namespace}
                          tableName={props.state.tableName}
                          linkToken={props.linkToken}
                        />
                        <XMarkIcon
                          className="tw-h-4 tw-ml-2 tw-text-slate-400 tw-cursor-pointer"
                          onClick={() => removeJsonField(fieldMappingIdx!, jsonIdx)}
                        />
                      </div>
                    );
                  })}
                </div>
              </div>
              <Button className="tw-ml-auto tw-mt-4" onClick={() => addJsonField(fieldMapping, fieldMappingIdx!)}>
                Add Field
              </Button>
            </div>
          );
        }

        return (
          <div className="tw-p-3">
            <div key={objectField.name} className="tw-flex tw-flex-row tw-items-center tw-justify-between">
              <MappedField
                objectField={objectField}
                fieldMapping={fieldMapping}
                fieldMappingIdx={fieldMappingIdx!}
                updateFieldMapping={updateFieldMapping}
              />
              <LinkFieldSelector
                className="tw-mt-0 tw-w-96 tw-flex"
                field={fieldMapping?.sourceField}
                setField={(value: Field) => {
                  updateFieldMapping({ ...fieldMapping, sourceField: value }, fieldMappingIdx!);
                }}
                placeholder="Choose a field"
                noOptionsString="No Fields Available!"
                validated={true}
                source={props.state.source}
                namespace={props.state.namespace}
                tableName={props.state.tableName}
                linkToken={props.linkToken}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
};

const MappedField: React.FC<{
  objectField: ObjectField;
  fieldMapping: FieldMappingState;
  fieldMappingIdx: number;
  updateFieldMapping: (fieldMapping: FieldMappingState, index: number) => void;
}> = ({ objectField, fieldMapping, fieldMappingIdx, updateFieldMapping }) => {
  return (
    <div className="tw-mr-2">
      <div className="tw-flex tw-h-fit">
        <div className="tw-h-fit tw-border tw-border-slate-200 tw-rounded-md tw-px-2 tw-box-border tw-bg-slate-100 tw-flex tw-flex-row tw-items-center tw-text-slate-700 tw-font-mono tw-select-none">
          <div className="tw-w-fit">{objectField.display_name ? objectField.display_name : objectField.name}</div>
          {objectField.description && (
            <Tooltip placement="top-start" label={objectField.description}>
              <InfoIcon className="tw-ml-2 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          )}
        </div>
        <div className="tw-h-fit tw-ml-3 tw-lowercase tw-select-none tw-font-mono tw-text-slate-500 tw-flex">
          {objectField.type}
          {!objectField.optional && <span className="tw-text-red-500">&nbsp;*</span>}
        </div>
      </div>
      {objectField.type === FieldType.Json &&
        (fieldMapping.expandedJson ? (
          <div
            className="tw-ml-1 tw-mt-2 tw-flex tw-items-center tw-cursor-pointer tw-text-xs tw-text-blue-600 tw-select-none"
            onClick={() => updateFieldMapping({ ...fieldMapping, expandedJson: false }, fieldMappingIdx!)}
          >
            Collapse JSON
            <Tooltip placement="top-start" label="Map a single field from your data source to this destination field.">
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
        ) : (
          <div
            className="tw-ml-1 tw-mt-2 tw-flex tw-items-center tw-cursor-pointer tw-text-xs tw-text-blue-600 tw-select-none"
            onClick={() => updateFieldMapping({ ...fieldMapping, expandedJson: true }, fieldMappingIdx!)}
          >
            Expand JSON
            <Tooltip
              placement="top-start"
              label="Map multiple fields from your data source to this destination JSON field."
            >
              <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
            </Tooltip>
          </div>
        ))}
    </div>
  );
};
