import { EyeIcon, EyeSlashIcon, Square2StackIcon } from "@heroicons/react/24/outline";
import { useState } from "react";
import { Tooltip } from "src/components/tooltip/Tooltip";

export const PrivateKey: React.FC<{ keyValue: string | undefined }> = ({ keyValue }) => {
  const [visible, setVisible] = useState<boolean>(false);
  const [copyText, setCopyText] = useState<string>("Copy");
  const copy = () => {
    navigator.clipboard.writeText(keyValue ? keyValue : "");
    setCopyText("Copied!");
    setTimeout(() => setCopyText("Copy"), 1200);
  };

  return (
    <div className="tw-border tw-border-solid tw-border-slate-300 tw-rounded-lg tw-max-w-lg tw-overflow-x-auto tw-overscroll-contain tw-p-2 tw-bg-white">
      {visible ? (
        <div className="tw-flex tw-items-center">
          <EyeSlashIcon className="tw-h-4 tw-ml-1 tw-mr-2 tw-cursor-pointer" onClick={() => setVisible(false)} />
          {keyValue}
          <Tooltip label={copyText} placement="top" hideOnClick={false}>
            <Square2StackIcon className="tw-h-4 tw-ml-auto tw-mr-1 tw-cursor-pointer tw-outline-none" onClick={copy} />
          </Tooltip>
        </div>
      ) : (
        <div className="tw-flex tw-items-center">
          <EyeIcon className="tw-h-4 tw-cursor-pointer tw-ml-1 tw-mr-2" onClick={() => setVisible(true)} />
          •••••••••••••••••••••••••••••••••••••••••••••••••••
          <Tooltip label={copyText} placement="top" hideOnClick={false}>
            <Square2StackIcon className="tw-h-4 tw-ml-auto tw-mr-1 tw-cursor-pointer tw-outline-none" onClick={copy} />
          </Tooltip>
        </div>
      )}
    </div>
  );
};
