import { ExclamationTriangleIcon, InformationCircleIcon } from "@heroicons/react/24/outline";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { mergeClasses } from "src/utils/twmerge";

export const Callout: React.FC<{ className?: string; tooltip?: string; content: string }> = (props) => {
  return (
    <div className={mergeClasses(props.className, "tw-rounded-md tw-bg-yellow-50 tw-p-3")}>
      <div className="tw-flex tw-items-center">
        <ExclamationTriangleIcon className="tw-h-5 tw-text-yellow-500" aria-hidden="true" />
        <div className="tw-ml-3 tw-text-sm tw-text-yellow-700">
          <p>{props.content}</p>
        </div>
        {props.tooltip && (
          <Tooltip content={props.tooltip}>
            <InformationCircleIcon className="tw-h-4 tw-text-slate-600 tw-ml-auto" />
          </Tooltip>
        )}
      </div>
    </div>
  );
};
