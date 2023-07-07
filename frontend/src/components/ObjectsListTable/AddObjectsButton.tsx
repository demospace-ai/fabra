import { PlusCircleIcon } from "@heroicons/react/20/solid";
import { Button } from "src/components/button/Button";

export function AddObjectButton({ onClick }: { onClick: () => void }) {
  return (
    <Button className="tw-flex tw-gap-x-2 tw-justify-center tw-items-center" onClick={onClick}>
      <PlusCircleIcon className="tw-h-4 tw-w-4" />
      <div className="tw-mr-0.5">Add Object</div>
    </Button>
  );
}
