import { CheckIcon } from "@heroicons/react/24/outline";
import * as RadixCheckbox from "@radix-ui/react-checkbox";
import { mergeClasses } from "src/utils/twmerge";

export const Checkbox: React.FC<{
  className: string;
  checked: boolean;
  disabled?: boolean;
  onCheckedChange: (checked: boolean) => void;
}> = ({ className, checked, disabled, onCheckedChange }) => {
  return (
    <RadixCheckbox.Root
      disabled={disabled}
      checked={checked}
      onCheckedChange={onCheckedChange}
      className={mergeClasses(
        "tw-bg-white tw-border-[1.2px] tw-border-slate-800 tw-rounded",
        checked && "tw-bg-slate-100",
        disabled && "tw-bg-gray-100 tw-border-gray-300",
        className,
      )}
    >
      <RadixCheckbox.Indicator>
        <CheckIcon className="tw-p-0.5 tw-stroke-[2]" />
      </RadixCheckbox.Indicator>
    </RadixCheckbox.Root>
  );
};
