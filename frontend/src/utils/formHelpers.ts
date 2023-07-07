import { FormEvent, SetStateAction, useState } from "react";

/**
 * Manages form state. API is similar to that of react-hook-form's useForm.
 * https://react-hook-form.com/get-started
 */
export function useForm<S = any>(
  initialState: S,
  opts: {
    validate: (state: S) => Record<string, string>;
  } = {
    validate: () => {
      return {};
    },
  },
) {
  const [state, setState] = useState<S>(initialState);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isDirty, setIsDirty] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const setStateDispatch: React.Dispatch<SetStateAction<S>> = (s) => {
    setIsDirty(true);
    setState(s);
  };

  const handleSubmit = (submitFn: (() => Promise<void>) | (() => void)) => {
    return (e: FormEvent) => {
      e.preventDefault();
      setIsSubmitted(true);
      const errors = opts.validate(state);
      setErrors(errors);
      if (Object.keys(errors).length > 0) {
        return;
      }
      submitFn();
    };
  };

  return {
    state,
    errors,
    setState: setStateDispatch,
    isDirty,
    isSubmitted,
    handleSubmit,
  };
}
