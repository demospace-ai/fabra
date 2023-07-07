import { useCallback, useState } from "react";
import {
  DestinationSetupBigQueryFormType,
  DestinationSetupDynamoDbFormType,
  DestinationSetupFormSchema,
  DestinationSetupFormType,
  DestinationSetupUnsupportedFormType,
  DestinationSetupWebhookFormType,
  FinalizeObjectFormSchema,
  FinalizeObjectFormType,
  ObjectFieldsFormType,
  ObjectFieldsSchema,
  Step,
} from "src/pages/objects/helpers";
import {
  ConnectionType,
  Destination,
  FabraObject,
  Field,
  getConnectionType,
  ObjectField,
  shouldCreateFields,
  TargetType,
} from "src/rpc/api";
import { z } from "zod";

const InitialStepSchema = z.object({
  step: z.literal(Step.Initial),
  destinationSetup: DestinationSetupFormSchema.optional(),
  objectFields: ObjectFieldsSchema.optional(),
  finalize: FinalizeObjectFormSchema.optional(),
});
type InitialStepSchema = z.infer<typeof InitialStepSchema>;

const CreatingObjectStepSchema = z.object({
  step: z.literal(Step.CreateFields),
  destinationSetup: DestinationSetupFormSchema,
  objectFields: ObjectFieldsSchema.optional(),
  finalize: FinalizeObjectFormSchema.optional(),
});
type CreatingObjectStepSchema = z.infer<typeof CreatingObjectStepSchema>;

const ExistingObjectFieldsStepSchema = z.object({
  step: z.literal(Step.ExistingFields),
  destinationSetup: DestinationSetupFormSchema,
  objectFields: ObjectFieldsSchema.optional(),
  finalize: FinalizeObjectFormSchema.optional(),
});
type ExistingObjectFieldsStepSchema = z.infer<typeof ExistingObjectFieldsStepSchema>;

const FinalizeStepSchema = z.object({
  step: z.literal(Step.Finalize),
  destinationSetup: DestinationSetupFormSchema,
  objectFields: ObjectFieldsSchema,
  finalize: FinalizeObjectFormSchema.optional(),
});
type FinalizeStepSchema = z.infer<typeof FinalizeStepSchema>;

/** If somehow we get to a destination whose connection type isn't supported yet. */
const UnsupportedConnectionTypeSchema = z.object({
  step: z.literal(Step.UnsupportedConnectionType),
  connectionType: z.enum([
    ConnectionType.Postgres,
    ConnectionType.Redshift,
    ConnectionType.Snowflake,
    ConnectionType.Synapse,
    ConnectionType.MongoDb,
    ConnectionType.MySQL,
  ]),
  message: z.string(),
});
type UnsupportedConnectionTypeSchema = z.infer<typeof UnsupportedConnectionTypeSchema>;

const StateSchema = z.discriminatedUnion("step", [
  InitialStepSchema,
  CreatingObjectStepSchema,
  FinalizeStepSchema,
  ExistingObjectFieldsStepSchema,
  UnsupportedConnectionTypeSchema,
]);
type StateSchemaType = z.infer<typeof StateSchema>;

type InitializeStateArgs = {
  existingObject: FabraObject | undefined;
  /** If user is updating an existing object. */
  existingDestination: Destination | undefined;
  /** Used for prefilling the destination field when creating an object from a destination details page. */
  maybeDestination: Destination | undefined;
};

function initializeState({
  existingObject,
  existingDestination,
  maybeDestination,
}: InitializeStateArgs): InitialStepSchema | UnsupportedConnectionTypeSchema {
  const connectionType =
    existingDestination?.connection.connection_type || maybeDestination?.connection.connection_type;
  if (!connectionType) {
    // There's no existing destination or maybeDestination, so we're creating a new object from scratch.
    return {
      step: Step.Initial,
    } as InitialStepSchema;
  }

  const destinationSetup = (() => {
    const base = {
      destination: existingDestination || maybeDestination,
      displayName: existingObject?.display_name ?? "",
      namespace: existingObject?.namespace,
      tableName: existingObject?.table_name,
    };
    switch (connectionType) {
      case ConnectionType.Webhook: {
        return {
          ...base,
          connectionType,
          targetType: TargetType.Webhook,
        } as DestinationSetupWebhookFormType;
      }
      case ConnectionType.BigQuery: {
        return {
          ...base,
          connectionType,
        } as DestinationSetupBigQueryFormType;
      }
      case ConnectionType.DynamoDb: {
        return {
          ...base,
          connectionType,
        } as DestinationSetupDynamoDbFormType;
      }
      case ConnectionType.Snowflake:
      case ConnectionType.Synapse:
      case ConnectionType.MySQL:
      case ConnectionType.MongoDb:
      case ConnectionType.Redshift:
      case ConnectionType.Postgres: {
        return {
          ...base,
          connectionType,
        } as DestinationSetupUnsupportedFormType;
      }
    }
  })();

  let objectFields: ObjectField[];
  let finalize: FinalizeObjectFormType | undefined;
  if (existingObject) {
    objectFields = existingObject.object_fields;
    const endCustomerIdField = objectFields.find((field) => field.name === existingObject.end_customer_id_field);
    let formCustomerIdField: Field | undefined;
    if (connectionType !== ConnectionType.Webhook) {
      if (!endCustomerIdField) {
        // This should never happen. Otherwise the server has a bug.
        // Maybe in the future we can return the full field in the API response instead of just the name.
        throw new Error("End Customer ID field not found");
      } else {
        formCustomerIdField = {
          name: existingObject.end_customer_id_field,
          type: endCustomerIdField.type!,
        };
      }
    }
    finalize = {
      connectionType,
      syncMode: existingObject.sync_mode,
      cursorField: objectFields.find((field) => field.name === existingObject.cursor_field),
      endCustomerIdField: formCustomerIdField,
      frequency: existingObject.frequency,
      frequencyUnits: existingObject.frequency_units,
      recurring: existingObject.recurring,
    };
  } else {
    finalize = undefined;
    objectFields = [];
  }

  const initialState: InitialStepSchema = {
    step: Step.Initial,
    destinationSetup,
    objectFields: {
      objectFields,
    },
    finalize,
  };

  return initialState;
}

export function useStateMachine(args: InitializeStateArgs, onComplete: () => void) {
  const [state, setState] = useState<StateSchemaType>(initializeState(args));
  const back = useCallback(() => {
    switch (state.step) {
      case Step.Initial: {
        onComplete();
        return;
      }
      case Step.CreateFields: {
        setState(
          (state) =>
            ({
              ...state,
              step: Step.Initial,
            } as InitialStepSchema),
        );
        return;
      }
      case Step.ExistingFields: {
        setState(
          (state) =>
            ({
              ...state,
              step: Step.Initial,
            } as InitialStepSchema),
        );
        return;
      }
      case Step.Finalize: {
        const createFields = shouldCreateFields(
          state.destinationSetup.connectionType,
          state.destinationSetup.targetType,
        );
        setState(
          (state) =>
            ({
              ...state,
              step: createFields ? Step.CreateFields : Step.ExistingFields,
            } as CreatingObjectStepSchema | ExistingObjectFieldsStepSchema),
        );
        return;
      }
      case Step.UnsupportedConnectionType: {
        setState(
          (state) =>
            ({
              ...state,
              step: Step.Initial,
            } as InitialStepSchema),
        );
        return;
      }
    }
  }, [state.step]);

  return {
    advanceToObjectFields: (destinationSetup: DestinationSetupFormType) => {
      const nextState = (() => {
        const connectionType = destinationSetup.connectionType;
        switch (connectionType) {
          case ConnectionType.BigQuery:
          case ConnectionType.Webhook:
          case ConnectionType.DynamoDb: {
            const createFields = shouldCreateFields(connectionType, destinationSetup.targetType);
            return {
              step: createFields ? Step.CreateFields : Step.ExistingFields,
              destinationSetup,
            } as CreatingObjectStepSchema;
          }
          default: {
            return {
              step: Step.UnsupportedConnectionType,
              connectionType: connectionType,
              message: `${getConnectionType(
                connectionType,
              )} destinations are not supported yet. Message the Fabra team about this!`,
            } as UnsupportedConnectionTypeSchema;
          }
        }
      })();
      setState((state) => {
        return { ...state, ...nextState };
      });
    },
    advanceToFinalizeObject: (destinationSetup: DestinationSetupFormType, objectFields: ObjectFieldsFormType) => {
      const nextState = {
        step: Step.Finalize,
        destinationSetup,
        objectFields,
      } as FinalizeStepSchema;
      setState((state) => ({ ...state, ...nextState }));
    },
    state,
    back,
  };
}
