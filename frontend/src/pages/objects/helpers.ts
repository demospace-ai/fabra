import { sendRequest } from "src/rpc/ajax";
import {
  ConnectionType,
  CreateObject,
  CreateObjectRequest,
  DestinationSchema,
  FabraObject,
  FieldSchema,
  FieldType,
  FrequencyUnits,
  ObjectFieldSchema,
  SyncMode,
  TargetType,
  UpdateObject,
  UpdateObjectFields,
  UpdateObjectFieldsRequest,
  UpdateObjectFieldsResponse,
  UpdateObjectRequest,
  UpdateObjectResponse,
} from "src/rpc/api";
import { z } from "zod";

export enum Step {
  Initial = "Initial",
  ExistingFields = "ExistingFields",
  CreateFields = "CreateFields",
  Finalize = "Finalize",
  UnsupportedConnectionType = "UnsupportedConnectionType",
}

/** Destination setup form. */
export const DestinationSetupBaseSchema = z.object({
  displayName: z.string().min(1, { message: "Please enter a display name" }),
  destination: DestinationSchema,
});

export const DestinationSetupWebhookSchema = DestinationSetupBaseSchema.extend({
  connectionType: z.literal(ConnectionType.Webhook),
  targetType: z.literal(TargetType.Webhook),
  namespace: z.string().optional(),
  tableName: z.string().optional(),
});
export type DestinationSetupWebhookFormType = z.infer<typeof DestinationSetupWebhookSchema>;

export const DestinationSetupBigQuerySchema = DestinationSetupBaseSchema.extend({
  connectionType: z.literal(ConnectionType.BigQuery),
  targetType: z.enum([TargetType.SingleExisting]),
  namespace: z.string(),
  tableName: z.string(),
});
export type DestinationSetupBigQueryFormType = z.infer<typeof DestinationSetupBigQuerySchema>;

export const DestinationSetupDynamoDbSchema = DestinationSetupBaseSchema.extend({
  connectionType: z.literal(ConnectionType.DynamoDb),
  tableName: z.string(),
  targetType: z.enum([TargetType.SingleExisting]),
  namespace: z.string().optional(),
});
export type DestinationSetupDynamoDbFormType = z.infer<typeof DestinationSetupDynamoDbSchema>;

export const DestinationSetupUnsupportedSchema = DestinationSetupBaseSchema.extend({
  connectionType: z.enum([
    ConnectionType.MongoDb,
    ConnectionType.Postgres,
    ConnectionType.Redshift,
    ConnectionType.Snowflake,
    ConnectionType.Synapse,
  ]),
  tableName: z.string().optional(),
  targetType: z
    .enum([TargetType.SingleExisting, TargetType.SingleNew, TargetType.TablePerCustomer])
    .default(TargetType.SingleExisting),
  namespace: z.string().optional(),
});
export type DestinationSetupUnsupportedFormType = z.infer<typeof DestinationSetupUnsupportedSchema>;

export const DestinationSetupFormSchema = z
  .discriminatedUnion("connectionType", [
    DestinationSetupWebhookSchema,
    DestinationSetupBigQuerySchema,
    DestinationSetupDynamoDbSchema,
    DestinationSetupUnsupportedSchema,
  ])
  .superRefine((values, ctx) => {
    if (!SUPPORTED_CONNECTION_TYPES.includes(values.connectionType as SupportedConnectionType)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Unsupported connection type",
        path: ["connectionType"],
      });
    }
  });

export type DestinationSetupFormType = z.infer<typeof DestinationSetupFormSchema>;

/** Object fields form. */
export const ObjectFieldsSchema = z.object({
  objectFields: z
    .array(
      ObjectFieldSchema.partial({
        id: true,
        field_type: true,
      }).superRefine((values, ctx) => {
        if (!values.field_type) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Must set a field type",
            path: ["field_type"],
          });
        }
      }),
    )
    .min(1, { message: "Must have at least one field" }),
});

export type ObjectFieldsFormType = z.infer<typeof ObjectFieldsSchema>;

export const SUPPORTED_CONNECTION_TYPES = [
  ConnectionType.BigQuery,
  ConnectionType.DynamoDb,
  ConnectionType.Webhook,
] as const;
export type SupportedConnectionType = (typeof SUPPORTED_CONNECTION_TYPES)[number];

/** Finalize object form. */
export const FinalizeObjectFormSchema = z
  .object({
    connectionType: z.nativeEnum(ConnectionType),
    recurring: z.boolean(),
    cursorField: FieldSchema.optional(),
    primaryKey: FieldSchema.optional(),
    frequency: z
      .number({
        errorMap: (issue) => {
          if (issue.message === "Expected number, received nan") {
            return { message: "Please enter a valid number" };
          }
          return { message: issue.message ?? "Please enter a valid number" };
        },
      })
      .min(1)
      .optional(),
    frequencyUnits: z.nativeEnum(FrequencyUnits).optional(),
    syncMode: z.nativeEnum(SyncMode),
    endCustomerIdField: FieldSchema.optional(),
  })
  .superRefine((values, ctx) => {
    if (values.connectionType !== ConnectionType.Webhook) {
      if (!values.endCustomerIdField) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Please select an end customer ID field",
          path: ["endCustomerIdField"],
        });
      }
    }

    if (values.recurring) {
      if (!values.frequency) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Please enter a frequency",
          path: ["frequency"],
        });
      }
      if (!values.frequencyUnits) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Please select a frequency unit",
          path: ["frequencyUnits"],
        });
      }

      if (values.frequencyUnits === FrequencyUnits.Minutes && values.frequency && values.frequency < 30) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Minimum frequency is 30 minutes",
          path: ["frequency"],
        });
      }
    }

    if ([SyncMode.IncrementalAppend, SyncMode.IncrementalUpdate].includes(values.syncMode)) {
      if (!values.cursorField) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Please select a cursor field",
          path: ["cursorField"],
        });
      } else if (
        ![
          FieldType.Timestamp,
          FieldType.DatetimeTz,
          FieldType.DatetimeNtz,
          FieldType.Date,
          FieldType.Integer,
          FieldType.Number,
        ].includes(values.cursorField.type as FieldType)
      ) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Cursor field must be a timestamp, date, integer, or number",
          path: ["cursorField"],
        });
      }
    }

    if ([SyncMode.IncrementalUpdate].includes(values.syncMode)) {
      if (!values.primaryKey) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Please select a primary key",
          path: ["primaryKey"],
        });
      }
    }
  });

export type FinalizeObjectFormType = z.infer<typeof FinalizeObjectFormSchema>;

export type ObjectTargetOption = {
  type: TargetType;
  title: string;
  description: string;
};
export const objectTargetOptions: ObjectTargetOption[] = [
  {
    type: TargetType.SingleExisting,
    title: "Single Existing Table",
    description:
      "Data from all of your customers will be stored in a single existing table, with an extra ID column to distinguish between customers.",
  },
  // TODO
  // {
  //   type: TargetType.SingleNew,
  //   title: "Single New Table",
  //   description: "Data from all of your customers will be stored in a single new table, with an extra ID column to distinguish between customers."
  // },
  // {
  //   type: TargetType.TablePerCustomer,
  //   title: "Table Per Customer",
  //   description: "Data from each of your customers will be stored in a separate table in your destination. The name of the table will include the customer's ID as a suffix."
  // },
];

export const createNewObject = async (args: {
  destinationSetup: DestinationSetupFormType;
  objectFields: ObjectFieldsFormType;
  finalizeValues: FinalizeObjectFormType;
}) => {
  const { destinationSetup, objectFields, finalizeValues } = args;
  const payload: CreateObjectRequest = {
    display_name: destinationSetup.displayName,
    destination_id: destinationSetup.destination.id,
    target_type: destinationSetup.targetType,
    namespace: destinationSetup.namespace,
    table_name: destinationSetup.tableName,
    sync_mode: finalizeValues.syncMode,
    cursor_field: finalizeValues.cursorField?.name,
    primary_key: finalizeValues.primaryKey?.name,
    end_customer_id_field: finalizeValues.endCustomerIdField ? finalizeValues.endCustomerIdField.name : undefined,
    recurring: finalizeValues.recurring,
    frequency: finalizeValues.frequency,
    frequency_units: finalizeValues.frequencyUnits,
    object_fields: objectFields.objectFields,
  };
  const object = await sendRequest(CreateObject, payload);
  return object.object;
};

export const updateObject = async (args: {
  existingObject: FabraObject;
  destinationSetup: DestinationSetupFormType;
  objectFields: ObjectFieldsFormType;
  finalizeValues: FinalizeObjectFormType;
}) => {
  const { objectFields, destinationSetup, existingObject, finalizeValues } = args;
  if (!existingObject) {
    throw new Error("Cannot update object without existing object");
  }

  // For object field updates, we need to compute the change sets.
  // TODO: support adding and removing fields when updating objects
  const updatedFields = objectFields.objectFields.filter((field) =>
    existingObject.object_fields.find((existingField) => existingField.name === field.name),
  );

  const [updateObjectResponse, _] = await Promise.all([
    sendRequest<UpdateObjectRequest, UpdateObjectResponse>(UpdateObject, {
      objectID: Number(existingObject.id),
      display_name: destinationSetup.displayName,
      destination_id: destinationSetup.destination?.id,
      target_type: destinationSetup.targetType,
      namespace: destinationSetup.namespace,
      table_name: destinationSetup.tableName,
      sync_mode: finalizeValues.syncMode,
      cursor_field: finalizeValues.cursorField?.name,
      end_customer_id_field: finalizeValues.endCustomerIdField?.name,
      recurring: finalizeValues.recurring,
      frequency: finalizeValues.frequency,
      frequency_units: finalizeValues.frequencyUnits,
    }),
    sendRequest<UpdateObjectFieldsRequest, UpdateObjectFieldsResponse>(UpdateObjectFields, {
      objectID: Number(existingObject.id),
      object_fields: updatedFields as UpdateObjectFieldsRequest["object_fields"],
    }),
  ]);

  return updateObjectResponse.object;
};
