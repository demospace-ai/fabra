import { Schema } from "src/rpc/api";

export type DateRange = {
  minDate: Date;
  maxDate: Date;
};

export const getDateStringInUTC = (d: Date): string => {
  const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  return months[d.getUTCMonth()] + " " + d.getUTCDate() + " " + d.getUTCFullYear();
};

export const formatSchema = (schema: Schema): Schema => {
  return schema.map((fieldSchema) => {
    return {
      name: fieldSchema.name.replaceAll("_", " "),
      type: fieldSchema.type,
    };
  });
};
