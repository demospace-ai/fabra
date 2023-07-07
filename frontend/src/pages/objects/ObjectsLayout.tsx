import { Outlet } from "react-router-dom";

export function ObjectsLayout() {
  return (
    <div className="tw-h-full tw-py-5 tw-px-10 tw-overflow-auto">
      <Outlet />
    </div>
  );
}
