import { Transition } from "@headlessui/react";
import { CheckCircleIcon, InformationCircleIcon, XCircleIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { Fragment } from "react";
import { useConnectDispatch } from "src/connect/model";
import { useDispatch } from "src/root/model";

export interface ToastOptions {
  type: "error" | "success" | "info";
  duration?: number;
  content: React.ReactNode;
}

type ToastProps = {
  content: React.ReactNode;
  show: boolean;
  close: () => void;
  duration?: number;
};

export type ShowToastFunction = (type: "success" | "error" | "info", content: string, duration?: number) => void;

export const getToastContentFromDetails = (toast?: ToastOptions) => {
  var toastContent = undefined;
  if (toast) {
    switch (toast.type) {
      case "error":
        toastContent = (
          <div className="tw-flex tw-flex-row tw-items-center tw-justify-start">
            <XCircleIcon className="tw-w-5 tw-h-5 tw-text-red-500 tw-stroke-2" />
            <p className="tw-ml-2 tw-text-sm tw-text-gray-900">{toast.content}</p>
          </div>
        );
        break;
      case "success":
        toastContent = (
          <div className="tw-flex tw-flex-row tw-items-center tw-justify-start">
            <CheckCircleIcon className="tw-w-5 tw-h-5 tw-text-green-500 tw-stroke-2" />
            <p className="tw-ml-2 tw-text-base tw-text-gray-900">{toast.content}</p>
          </div>
        );
        break;
      case "info":
        toastContent = (
          <div className="tw-flex tw-flex-row tw-items-center tw-justify-start">
            <InformationCircleIcon className="tw-w-5 tw-h-5 tw-text-yellow-500 tw-stroke-2" />
            <p className="tw-ml-2 tw-text-base tw-text-gray-900">{toast.content}</p>
          </div>
        );
        break;
    }
  }

  return toastContent;
};

export const useShowToast = (): ShowToastFunction => {
  const dispatch = useDispatch();
  return (type: "success" | "error" | "info", content: string, duration?: number) => {
    dispatch({ type: "toast", toast: { content, type, duration } });
  };
};

export const useConnectShowToast = (): ShowToastFunction => {
  const dispatch = useConnectDispatch();
  return (type: "success" | "error" | "info", content: string, duration?: number) => {
    dispatch({ type: "toast", toast: { content, type, duration } });
  };
};

export const Toast: React.FC<ToastProps> = ({ content, show, duration, close }) => {
  if (duration) {
    setTimeout(() => {
      close();
    }, duration);
  }

  return (
    <>
      {/* Global notification live region, render this permanently at the end of the document */}
      <div
        aria-live="assertive"
        className="tw-pointer-events-none tw-fixed tw-inset-0 tw-flex tw-items-end tw-px-4 tw-py-6 sm:tw-items-start sm:tw-p-6"
      >
        <div className="tw-flex tw-w-full tw-flex-col tw-items-center tw-space-y-4 sm:tw-items-end">
          {/* Notification panel, dynamically insert this into the live region when it needs to be displayed */}
          <Transition
            show={show}
            as={Fragment}
            enter="tw-transform tw-ease-out tw-duration-300 tw-transition"
            enterFrom="tw-translate-y-2 tw-opacity-0 sm:tw-translate-y-0 sm:tw-translate-x-2"
            enterTo="tw-translate-y-0 tw-opacity-100 sm:tw-translate-x-0"
            leave="tw-transition tw-ease-in tw-duration-100"
            leaveFrom="tw-opacity-100"
            leaveTo="tw-opacity-0"
          >
            <div className="tw-pointer-events-auto tw-w-full tw-max-w-sm tw-overflow-hidden tw-rounded-lg tw-bg-white tw-shadow-lg tw-ring-1 tw-ring-slate-900 tw-ring-opacity-5">
              <div className="tw-p-4">
                <div className="tw-flex tw-items-center">
                  <div className="tw-ml-3 tw-w-0 tw-flex-1 tw-pt-0.5">{content}</div>
                  <div className="tw-ml-4 tw-flex tw-flex-shrink-0">
                    <button
                      type="button"
                      className="tw-inline-flex tw-rounded-md tw-bg-white tw-text-gray-400 hover:tw-text-gray-500 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-indigo-500 focus:tw-ring-offset-2"
                      onClick={() => {
                        close();
                      }}
                    >
                      <span className="tw-sr-only">Close</span>
                      <XMarkIcon className="tw-h-5 tw-w-5" aria-hidden="true" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </>
  );
};
