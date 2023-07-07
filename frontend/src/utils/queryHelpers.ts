import { useState } from "react";
import { consumeError } from "./errors";

type AsyncFunction<Data = any, Args = any> = (variables?: Args) => Promise<Data>;

export function useMutation<Data = any, Args = any>(
  mutationFn: AsyncFunction<Data, Args>,
  opts: { onSuccess?: (data: Data) => void; onError?: (err: Error) => void } = {
    onSuccess: () => {},
    onError: () => {},
  },
) {
  const [error, setError] = useState<Error | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [data, setData] = useState<Data | undefined>();
  const [isSuccess, setIsSuccess] = useState(false);
  const [isFailed, setIsFailed] = useState(false);
  const mutate = async (variables?: Args) => {
    setIsLoading(true);
    try {
      const response = await mutationFn(variables);
      setIsSuccess(true);
      setData(response);
      setError(null);
      setIsFailed(false);
      opts.onSuccess?.(response);
    } catch (err) {
      consumeError(err);
      console.error(err);
      if (err instanceof Error) {
        setError(err);
        opts.onError?.(err);
      } else {
        const unknownError = new Error("Unknown error");
        setError(unknownError);
        opts.onError?.(unknownError);
      }
      setIsFailed(true);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    mutate,
    error,
    isLoading,
    isSuccess,
    isFailed,
    data,
    reset: () => {
      setIsLoading(false);
      setIsSuccess(false);
      setIsFailed(false);
      setData(undefined);
      setError(null);
    },
  };
}
