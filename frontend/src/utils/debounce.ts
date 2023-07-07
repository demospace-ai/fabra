import { useState } from "react";

type Timer = ReturnType<typeof setTimeout>;

export function useDebounce<F extends (...args: any[]) => void>(func: F, delayMs: number) {
  const [timer, setTimer] = useState<Timer>(); //Create timer state

  const debouncedFunction = ((...args) => {
    const newTimer = setTimeout(() => {
      func(...args);
    }, delayMs);
    clearTimeout(timer);
    setTimer(newTimer);
  }) as F;

  return debouncedFunction;
}
