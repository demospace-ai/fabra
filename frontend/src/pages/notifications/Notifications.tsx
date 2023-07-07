import { SparklesIcon } from "@heroicons/react/24/solid";
import { isProd } from "src/utils/env";

export const Notifications: React.FC = () => {
  return (
    <div className="tw-h-full tw-flex tw-justify-center">
      <div className="tw-mt-48 tw-w-full tw-max-w-sm tw-text-center">
        <div className="tw-flex tw-items-center tw-justify-center tw-gap-x-1">
          <h2 className="tw-text-xl tw-font-bold">Coming soon!</h2>
        </div>
        <div className="tw-mt-2 tw-text-center tw-text-base">
          Want to receive Slack or email notifications if a sync fails? Reach out for early access!
        </div>
        <button
          onClick={() => {
            if (isProd()) window.Intercom("showNewMessage", "I'd like early access to notifications.");
          }}
          className="tw-mt-4 tw-inline-flex tw-items-center tw-rounded-md tw-border tw-border-solid tw-border-slate-300 tw-px-3 tw-py-2 tw-text-sm tw-font-medium tw-shadow hover:tw-bg-slate-100"
        >
          <SparklesIcon className="tw-h-4 tw-mr-1.5 tw-fill-yellow-300" />
          Request Access
          <SparklesIcon className="tw-h-4 tw-ml-1.5 tw-fill-yellow-300" />
        </button>
      </div>
    </div>
  );
};
