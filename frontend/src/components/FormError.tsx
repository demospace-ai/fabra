import { mergeClasses } from "src/utils/twmerge";

export function FormError({ message, className }: { message: string | undefined; className?: string }) {
  return <div className={mergeClasses("tw-text-red-500", className)}>{message}</div>;
}
