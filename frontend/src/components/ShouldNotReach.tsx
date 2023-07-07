import { useEffect } from "react";
import { consumeError } from "src/utils/errors";

/** Use this to render things that should never be rendered. Like fallback states in switch statements. */
export function ShouldNotReach({ children, error }: { error?: Error; children: React.ReactNode }) {
  useEffect(() => {
    consumeError(error ?? new Error("Should not reach this component"));
  }, []);

  return <>{children}</>;
}
