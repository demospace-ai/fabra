import { Combobox, Listbox, Transition } from "@headlessui/react";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/react/20/solid";
import { Fragment, HTMLInputTypeAttribute, InputHTMLAttributes, useRef, useState } from "react";
import { Modifier, usePopper } from "react-popper";
import { Loading } from "src/components/loading/Loading";
import { mergeClasses } from "src/utils/twmerge";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  placeholder?: string;
  setValue: (value: any) => void;
  className?: string;
  wrapperClass?: string;
  textarea?: boolean;
  type?: HTMLInputTypeAttribute;
  label?: string;
  disabled?: boolean;
}

export const Input: React.FC<InputProps> = (props) => {
  const { id, value, placeholder, disabled, setValue, className, wrapperClass, textarea, type, label, ...other } =
    props;
  const [focused, setFocused] = useState(false);
  let classes = [
    "tw-border tw-border-solid tw-border-slate-300 tw-rounded-md tw-py-2.5 tw-px-3 tw-w-full tw-box-border focus:tw-border-slate-700 tw-outline-none",
    !disabled && "hover:tw-border-slate-400",
    disabled && "tw-bg-slate-100 tw-text-slate-500 tw-cursor-not-allowed",
    className,
  ];

  const onKeydown = (event: React.KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    event.stopPropagation();
    if (event.key === "Escape") {
      event.currentTarget.blur();
    }
  };

  const showLabel =
    label !== undefined && (focused || (value !== undefined && (value !== "string" || value.length > 0)));

  return (
    <div className={mergeClasses("tw-relative", label && "tw-mt-4", wrapperClass)}>
      <Transition
        show={showLabel}
        enter="tw-transition tw-ease tw-duration-200 tw-transform"
        enterFrom="tw-translate-y-4 tw-opacity-10"
        enterTo="tw-translate-y-0"
        leave="tw-transition tw-ease tw-duration-200 tw-transform"
        leaveFrom="tw-translate-y-0"
        leaveTo="tw-translate-y-4 tw-opacity-10"
      >
        <label
          htmlFor="name"
          className="tw-absolute -tw-top-2 tw-left-2 -tw-mt-px tw-inline-block tw-bg-white tw-px-1 tw-text-xs tw-font-medium tw-text-gray-900"
        >
          {label}
        </label>
      </Transition>
      {textarea ? (
        <textarea
          id={id}
          name={id}
          autoComplete={id}
          placeholder={focused ? undefined : placeholder}
          className={mergeClasses(classes)}
          onKeyDown={onKeydown}
          onFocus={() => {
            setFocused(true);
          }}
          onChange={(e) => setValue(e.target.value)}
          onBlur={() => {
            setFocused(false);
          }}
          value={value ? value : ""}
          disabled={disabled}
        />
      ) : (
        <input
          type={type ? type : "text"}
          id={id}
          name={id}
          autoComplete={id}
          placeholder={focused ? undefined : placeholder}
          className={mergeClasses(classes)}
          onKeyDown={onKeydown}
          onFocus={() => {
            setFocused(true);
          }}
          onChange={(e) => {
            e.target.type === "number" ? setValue(parseInt(e.target.value)) : setValue(e.target.value);
          }}
          onBlur={() => {
            setFocused(false);
          }}
          value={value ? value : ""}
          disabled={disabled}
          {...other}
        />
      )}
    </div>
  );
};

export const ColorPicker: React.FC<InputProps> = (props) => {
  const { id, value, setValue, className, wrapperClass, textarea, type, label, ...other } = props;
  let classes = [
    "tw-border tw-border-solid tw-border-slate-300 tw-rounded-md tw-py-2.5 tw-px-3 tw-w-full tw-box-border focus:tw-border-slate-700 tw-outline-none",
    className,
  ];

  return (
    <div className={mergeClasses("tw-flex tw-items-center tw-relative", wrapperClass)}>
      <input className={mergeClasses(classes)} value={value} onChange={(e) => setValue(e.target.value)} {...other} />
      <input
        className="tw-w-6 tw-absolute tw-right-3"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        type="color"
      />
    </div>
  );
};

export const InputStyle =
  "tw-border tw-border-solid tw-border-slate-300 tw-rounded-md tw-py-2.5 tw-px-3 tw-w-full tw-box-border hover:tw-border-primary-hover focus:tw-border-primary tw-outline-none";

export const ValidatedInput: React.FC<InputProps> = (props) => {
  const { id, value, placeholder, setValue, className, textarea, type, label, ...other } = props;
  const [isValid, setIsValid] = useState(true);
  const [focused, setFocused] = useState(false);
  let classes = [
    "tw-border tw-border-solid tw-border-slate-300 tw-rounded-md tw-py-2.5 tw-px-3 tw-w-full tw-box-border hover:tw-border-primary-hover focus:tw-border-primary tw-outline-none",
    className,
  ];
  if (!isValid) {
    classes.push("tw-border-red-600");
  }

  const onKeydown = (event: React.KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    event.stopPropagation();
    if (event.key === "Escape") {
      event.currentTarget.blur();
    }
  };

  const validateNotEmpty = (value: string | number | readonly string[] | undefined): boolean => {
    if (typeof value === "number") {
      return true;
    }

    const valid = value !== undefined && value.length > 0;
    setIsValid(valid);
    return valid;
  };

  const showLabel =
    label !== undefined && (focused || (value !== undefined && (value !== "string" || value.length > 0)));

  return (
    <div className={mergeClasses("tw-relative", label && "tw-mt-4")}>
      <Transition
        show={showLabel}
        enter="tw-transition tw-ease tw-duration-200 tw-transform"
        enterFrom="tw-translate-y-4 tw-opacity-10"
        enterTo="tw-translate-y-0"
        leave="tw-transition tw-ease tw-duration-200 tw-transform"
        leaveFrom="tw-translate-y-0"
        leaveTo="tw-translate-y-4 tw-opacity-10"
      >
        <label
          htmlFor="name"
          className="tw-absolute -tw-top-2 tw-left-2 -tw-mt-px tw-inline-block tw-bg-white tw-px-1 tw-text-xs tw-font-medium tw-text-primary"
        >
          {label}
        </label>
      </Transition>
      {textarea ? (
        <textarea
          id={id}
          name={id}
          autoComplete={id}
          placeholder={focused ? undefined : placeholder}
          className={mergeClasses(classes)}
          onKeyDown={onKeydown}
          onFocus={() => {
            setIsValid(true);
            setFocused(true);
          }}
          onChange={(e) => setValue(e.target.value)}
          onBlur={() => {
            validateNotEmpty(value);
            setFocused(false);
          }}
          value={value ? value : ""}
        />
      ) : (
        <input
          type={type ? type : "text"}
          name={id}
          autoComplete={id}
          placeholder={focused ? undefined : placeholder}
          className={mergeClasses(classes)}
          onKeyDown={onKeydown}
          onFocus={() => {
            setIsValid(true);
            setFocused(true);
          }}
          onChange={(e) => {
            e.target.type === "number" ? setValue(parseInt(e.target.value)) : setValue(e.target.value);
          }}
          onBlur={() => {
            validateNotEmpty(value);
            setFocused(false);
          }}
          value={value ? value : ""}
          {...other}
        />
      )}
    </div>
  );
};

export type ValidatedDropdownInputProps = {
  id?: string;
  options: any[] | undefined;
  selected: any | undefined;
  setSelected: (option: any) => void;
  getElementForDisplay?: (option: any) => string | React.ReactElement;
  getElementForDropdown?: (option: any) => string | React.ReactElement;
  loading: boolean;
  placeholder: string;
  noOptionsString: string;
  className?: string;
  validated?: boolean;
  noCaret?: boolean;
  label?: string;
  dropdownHeight?: string;
  valid?: boolean;
  disabled?: boolean;
};

export const ValidatedDropdownInput: React.FC<ValidatedDropdownInputProps> = (props) => {
  const [computedValid, setComputedValid] = useState<boolean>(true);
  const [focused, setFocused] = useState(false);
  const [referenceElement, setReferenceElement] = useState<HTMLButtonElement | null>(null);
  const [popperElement, setPopperElement] = useState<HTMLDivElement | null>(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    strategy: "fixed",
    modifiers: [sameWidth],
  });

  const validateNotUndefined = (value: number | undefined): boolean => {
    const valid = value !== undefined;
    setComputedValid(valid);
    return valid;
  };

  const getElementForDisplay = props.getElementForDisplay ? props.getElementForDisplay : (value: any) => value;
  const getElementForDropdown = props.getElementForDropdown ? props.getElementForDropdown : getElementForDisplay;
  const showLabel = props.label !== undefined && (focused || props.selected !== undefined);

  // An undefined value will cause the input to be uncontrolled, so change to null
  const value = props.selected === undefined ? null : props.selected;

  // Allow explicitly passing the valid property to override the computed value
  const valid = props.valid !== undefined ? props.valid : computedValid;

  return (
    <Listbox
      value={value}
      disabled={props.disabled}
      onChange={(value) => {
        props.setSelected(value);
        setComputedValid(true);
      }}
    >
      <div className="tw-relative">
        <Transition
          show={showLabel}
          enter="tw-transition tw-ease tw-duration-200 tw-transform"
          enterFrom="tw-translate-y-4 tw-opacity-10"
          enterTo="tw-translate-y-0"
          leave="tw-transition tw-ease tw-duration-200 tw-transform"
          leaveFrom="tw-translate-y-0"
          leaveTo="tw-translate-y-4 tw-opacity-10"
        >
          <label
            htmlFor="name"
            className="tw-absolute -tw-top-2 tw-left-2 -tw-mt-px tw-inline-block tw-bg-white tw-px-1 tw-text-xs tw-font-medium tw-text-primary"
          >
            {props.label}
          </label>
        </Transition>
        <Listbox.Button
          ref={setReferenceElement}
          className={mergeClasses(
            "tw-flex tw-justify-center tw-items-center tw-w-96 tw-mt-5 tw-rounded-md tw-py-2.5 tw-px-3 tw-text-left tw-border tw-border-solid tw-border-slate-300 aria-expanded:tw-border-primary",
            !props.disabled && "hover:tw-border-primary-hover",
            props.disabled && "tw-bg-slate-100 tw-text-slate-400 tw-cursor-not-allowed",
            props.className,
            props.validated && !valid && "tw-border-red-600",
          )}
        >
          <div
            className={mergeClasses(
              "tw-inline-block tw-w-[calc(100%-20px)] tw-truncate tw-overflow-none",
              !props.selected && "tw-text-slate-400",
            )}
          >
            {value ? getElementForDisplay(props.selected) : props.placeholder}
          </div>
          {!props.noCaret && (
            <span className="tw-pointer-events-none pr-2">
              <ChevronUpDownIcon
                className="tw-inline tw-float-right tw-h-5 tw-w-5 tw-text-slate-400"
                aria-hidden="true"
              />
            </span>
          )}
        </Listbox.Button>
        <div className="tw-relative tw-z-10" ref={setPopperElement} style={styles.popper} {...attributes.popper}>
          <Transition
            as={Fragment}
            enter="tw-transition tw-ease-out tw-duration-100"
            enterFrom="tw-transform tw-opacity-0 tw-scale-95"
            enterTo="tw-transform tw-opacity-100 tw-scale-100"
            leave="tw-transition tw-ease-in tw-duration-100"
            leaveFrom="tw-transform tw-opacity-100 tw-scale-100"
            leaveTo="tw-transform tw-opacity-0 tw-scale-95"
            beforeEnter={() => setFocused(true)}
            afterLeave={() => {
              validateNotUndefined(props.selected);
              setFocused(false);
            }}
          >
            <Listbox.Options
              className={mergeClasses(
                "tw-absolute tw-z-20 tw-mt-1 tw-max-h-60 tw-min-w-full tw-overflow-auto tw-rounded-md tw-bg-white tw-py-1 tw-text-base tw-shadow-lg tw-ring-1 tw-ring-slate-900 tw-ring-opacity-5 focus:tw-outline-none sm:tw-text-sm",
                props.dropdownHeight,
              )}
            >
              <DropdownOptions
                loading={props.loading}
                options={props.options}
                noOptionsString={props.noOptionsString}
                getElementForDisplay={getElementForDropdown}
              />
            </Listbox.Options>
          </Transition>
        </div>
      </div>
    </Listbox>
  );
};

type DropdownOptionsProps = {
  loading: boolean;
  options: any[] | undefined;
  noOptionsString: string;
  getElementForDisplay: (value: any) => string | React.ReactElement;
};

const DropdownOptions: React.FC<DropdownOptionsProps> = (props) => {
  if (props.loading) {
    return (
      <div className="tw-p-2">
        <Loading className="tw-m-auto tw-block" />
      </div>
    );
  }

  if (props.options && props.options.length > 0) {
    return (
      <>
        {props.options!.map((option: any, index: number) => (
          <Listbox.Option
            key={index}
            value={option}
            className={({ active, selected }) =>
              `tw-relative tw-cursor-pointer tw-select-none tw-py-2.5 tw-pl-4 tw-pr-4 ${
                active || selected ? "tw-bg-slate-100 tw-text-slate-900" : "tw-text-slate-900"
              }`
            }
          >
            {({ selected }) => (
              <>
                <span className={`tw-block tw-truncate ${selected ? "tw-font-medium" : "tw-font-normal"}`}>
                  {props.getElementForDisplay(option)}
                </span>
                {selected ? (
                  <span className="tw-absolute tw-inset-y-0 tw-right-0 tw-flex tw-items-center tw-pr-3 tw-text-slate-600">
                    <CheckIcon className="tw-h-5 tw-w-5" aria-hidden="true" />
                  </span>
                ) : null}
              </>
            )}
          </Listbox.Option>
        ))}
      </>
    );
  } else {
    return <div className="tw-p-2 tw-pl-4 tw-select-none">{props.noOptionsString}</div>;
  }
};

export type ValidatedComboInputProps = {
  options: any[] | undefined;
  selected: any | undefined;
  setSelected: (option: any) => void;
  getElementForDisplay?: (option: any) => string | React.ReactElement;
  loading: boolean;
  placeholder: string;
  noOptionsString: string;
  className?: string;
  validated?: boolean;
  allowCustom?: boolean;
  label?: string;
  dropdownHeight?: string;
  nullable?: boolean;
  valid?: boolean;
  disabled?: boolean;
};

export const ValidatedComboInput: React.FC<ValidatedComboInputProps> = (props) => {
  const [computedValid, setComputedValid] = useState<boolean>(true);
  const [query, setQuery] = useState("");
  const [focused, setFocused] = useState(false);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [referenceElement, setReferenceElement] = useState<HTMLDivElement | null>(null);
  const [popperElement, setPopperElement] = useState<HTMLDivElement | null>(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    strategy: "fixed",
    modifiers: [sameWidth],
  });

  const getFilteredOptions = () => {
    if (query === "") {
      return props.options;
    }

    if (props.options) {
      return props.options.filter((option) => {
        const displayValue = props.getElementForDisplay ? props.getElementForDisplay(option) : option;
        if (typeof displayValue === "string") {
          return displayValue.toLowerCase().includes(query.toLowerCase());
        } else {
          return true; // TODO: figure out how to do this for elements
        }
      });
    }

    return [];
  };
  const filteredOptions = getFilteredOptions();
  const getElementForDisplay = props.getElementForDisplay ? props.getElementForDisplay : (value: any) => value;
  const showLabel = props.label !== undefined && (focused || props.selected !== undefined);

  const validateNotUndefined = (value: number | undefined): boolean => {
    const valid = value !== undefined;
    setComputedValid(valid);
    return valid;
  };

  // An undefined value will cause the input to be uncontrolled, so change to null
  const value = props.selected === undefined ? null : props.selected;

  // Allow explicitly passing the valid property to override the computed value
  const valid = props.valid !== undefined ? props.valid : computedValid;

  return (
    <Combobox
      value={value}
      disabled={props.disabled}
      onChange={(value: number) => {
        props.setSelected(value);
        setComputedValid(true);
      }}
    >
      <div className="tw-relative">
        <Transition
          show={showLabel}
          enter="tw-transition tw-ease tw-duration-200 tw-transform"
          enterFrom="tw-translate-y-4 tw-opacity-10"
          enterTo="tw-translate-y-0"
          leave="tw-transition tw-ease tw-duration-200 tw-transform"
          leaveFrom="tw-translate-y-0"
          leaveTo="tw-translate-y-4 tw-opacity-10"
        >
          <label
            htmlFor="name"
            className="tw-absolute -tw-top-2 tw-left-2 -tw-mt-px tw-inline-block tw-bg-white tw-px-1 tw-text-xs tw-font-medium tw-text-primary"
          >
            {props.label}
          </label>
        </Transition>
        <div
          ref={setReferenceElement}
          className={mergeClasses(
            "tw-flex tw-w-96 tw-mt-5 tw-rounded-md tw-bg-white tw-py-2.5 tw-px-3 tw-text-left tw-border tw-border-solid tw-border-slate-300 focus-within:!tw-border-primary tw-transition tw-duration-100",
            !props.disabled && "hover:tw-border-primary-hover",
            props.disabled && "tw-bg-slate-100 tw-cursor-not-allowed",
            props.className,
            props.validated && !valid && "tw-border-red-600",
          )}
        >
          <Combobox.Input
            className={mergeClasses(
              "tw-inline tw-bg-transparent tw-w-[calc(100%-20px)] tw-border-none tw-text-sm tw-leading-5 tw-text-slate-900 tw-outline-none tw-text-ellipsis tw-cursor-pointer focus:tw-cursor-text",
              props.disabled && "tw-bg-slate-100 tw-text-slate-400 tw-cursor-not-allowed",
            )}
            displayValue={(selected) => (selected ? getElementForDisplay(selected) : "")}
            onChange={(event) => setQuery(event.target.value)}
            placeholder={props.placeholder}
            onClick={() => buttonRef.current?.click()}
          ></Combobox.Input>
          <Combobox.Button className="tw-inline-block tw-h-full" ref={buttonRef}>
            <span className="tw-pointer-events-none pr-2">
              <ChevronUpDownIcon
                className="tw-inline tw-float-right tw-h-5 tw-w-5 tw-text-slate-400"
                aria-hidden="true"
              />
            </span>
          </Combobox.Button>
        </div>
        <div className="tw-relative tw-z-10" ref={setPopperElement} style={styles.popper} {...attributes.popper}>
          <Transition
            as={Fragment}
            enter="tw-transition tw-ease-out tw-duration-100"
            enterFrom="tw-transform tw-opacity-0 tw-scale-95"
            enterTo="tw-transform tw-opacity-100 tw-scale-100"
            leave="tw-transition tw-ease-in tw-duration-100"
            leaveFrom="tw-transform tw-opacity-100 tw-scale-100"
            leaveTo="tw-transform tw-opacity-0 tw-scale-95"
            beforeEnter={() => setFocused(true)}
            afterLeave={() => {
              validateNotUndefined(props.selected);
              setQuery("");
              setFocused(false);
            }}
          >
            <Combobox.Options
              className={mergeClasses(
                "tw-absolute tw-z-20 tw-mt-1 tw-min-w-full tw-max-h-60 tw-overflow-auto tw-rounded-md tw-bg-white tw-py-1 tw-text-base tw-shadow-lg tw-ring-1 tw-ring-slate-900 tw-ring-opacity-5 focus:tw-outline-none sm:tw-text-sm",
                props.dropdownHeight,
              )}
            >
              <ComboOptions
                loading={props.loading}
                options={filteredOptions}
                noOptionsString={props.noOptionsString}
                getElementForDisplay={getElementForDisplay}
                query={query}
                allowCustom={props.allowCustom}
              />
            </Combobox.Options>
          </Transition>
        </div>
      </div>
    </Combobox>
  );
};

type ComboOptionsProps = {
  loading: boolean;
  options: any[] | undefined;
  noOptionsString: string;
  getElementForDisplay: (value: any) => string | React.ReactElement;
  query: string;
  allowCustom?: boolean;
};

const ComboOptions: React.FC<ComboOptionsProps> = (props) => {
  if (props.loading) {
    return (
      <div className="tw-p-2">
        <Loading className="tw-m-auto tw-block" />
      </div>
    );
  }

  if (props.options && props.options.length > 0) {
    return (
      <>
        {props.allowCustom && props.query.length > 0 && (
          <Combobox.Option
            value={props.query}
            className={({ active, selected }) =>
              `tw-relative tw-cursor-pointer tw-select-none tw-py-2.5 tw-pl-4 tw-pr-4 ${
                active || selected ? "tw-bg-slate-100 tw-text-slate-900" : "tw-text-slate-900"
              }`
            }
          >
            Custom: "{props.query}"
          </Combobox.Option>
        )}
        {props.options!.map((option: any, index: number) => (
          <Combobox.Option
            key={index}
            value={option}
            className={({ active, selected }) =>
              `tw-relative tw-cursor-pointer tw-select-none tw-py-2.5 tw-pl-4 tw-pr-4 ${
                active || selected ? "tw-bg-slate-100 tw-text-slate-900" : "tw-text-slate-900"
              }`
            }
          >
            {({ selected }) => (
              <>
                <span className={`tw-block tw-truncate ${selected ? "tw-font-medium" : "tw-font-normal"}`}>
                  {props.getElementForDisplay(option)}
                </span>
                {selected ? (
                  <span className="tw-absolute tw-inset-y-0 tw-right-0 tw-flex tw-items-center tw-pr-3 tw-text-slate-600">
                    <CheckIcon className="tw-h-5 tw-w-5" aria-hidden="true" />
                  </span>
                ) : null}
              </>
            )}
          </Combobox.Option>
        ))}
      </>
    );
  } else {
    return (
      <>
        {props.allowCustom ? (
          props.allowCustom &&
          props.query.length > 0 && (
            <Combobox.Option
              value={props.query}
              className={({ active, selected }) =>
                `tw-relative tw-cursor-pointer tw-select-none tw-py-2.5 tw-pl-4 tw-pr-4 ${
                  active || selected ? "tw-bg-slate-100 tw-text-slate-900" : "tw-text-slate-900"
                }`
              }
            >
              Custom: "{props.query}"
            </Combobox.Option>
          )
        ) : (
          <div className="tw-p-2 tw-pl-4 tw-select-none">{props.noOptionsString}</div>
        )}
      </>
    );
  }
};

const sameWidth: Modifier<"sameWidth"> = {
  name: "sameWidth",
  enabled: true,
  phase: "beforeWrite",
  requires: ["computeStyles"],
  fn: ({ state }) => {
    state.styles.popper.width = `${state.rects.reference.width}px`;
  },
  effect: ({ state }) => {
    state.elements.popper.style.width = `${state.elements.reference.getBoundingClientRect().width}px`;
  },
};
