import { useEffect, useState } from "react";

export const useLocalStorage = <T>(storageKey: string, fallbackState?: T): [T, (value: T) => void] => {
  const storedValue = localStorage.getItem(storageKey);
  const [value, setValue] = useState<T>(storedValue ? JSON.parse(storedValue) : fallbackState);

  useEffect(() => {
    if (value) {
      localStorage.setItem(storageKey, JSON.stringify(value));
    }
  }, [value, storageKey]);

  return [value, setValue];
};
