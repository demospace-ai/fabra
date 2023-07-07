import { forceError } from "src/utils/errors";

export const ErrorDisplay: React.FC<{ error: Error | unknown | string | null; className?: string }> = ({
  error,
  className,
}) => {
  const err = forceError(error);
  if (!err) {
    return null;
  }

  return <div className={className}>{err.message}</div>;
};
