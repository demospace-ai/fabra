import {
  ArrowPathIcon,
  ArrowTopRightOnSquareIcon,
  BellIcon,
  CircleStackIcon,
  CubeIcon,
  EyeIcon,
  HomeIcon,
  KeyIcon,
  MapIcon,
  UserPlusIcon,
} from "@heroicons/react/24/outline";
import classNames from "classnames";
import { NavLink } from "react-router-dom";
import { QuestionCircleIcon } from "src/components/icons/Icons";
import logo from "src/components/images/logo.svg";
import { useSelector } from "src/root/model";

export const NavigationBar: React.FC = () => {
  const isAuthenticated = useSelector((state) => state.login.authenticated);
  const organization = useSelector((state) => state.login.organization);

  // No navigation bar whatsoever for login page
  if (!isAuthenticated || !organization) {
    return <></>;
  }

  const route = "tw-inline-block tw-pt-[0.5px] tw-ml-2.5 tw-font-medium";
  const routeContainer =
    "tw-relative tw-flex tw-flex-row tw-h-9 tw-box-border tw-cursor-pointer tw-items-center tw-text-slate-800 tw-mt-0.5 tw-mb-0.5 tw-mx-2 tw-rounded-md tw-select-none";
  const navLink = "tw-w-full tw-h-full tw-pl-3 tw-rounded-md tw-flex tw-flex-row tw-items-center hover:tw-bg-slate-200";

  return (
    <>
      <div className="tw-min-w-[240px] tw-w-60 tw-h-full tw-flex tw-flex-col tw-box-border tw-border-r tw-border-solid tw-border-slate-200 tw-bg-white">
        <NavLink
          className="tw-py-4 tw-px-4 tw-flex tw-flex-row tw-h-16 tw-box-border tw-cursor-pointer tw-w-full tw-mb-4"
          to="/"
        >
          <img
            src={logo}
            className="tw-h-6 tw-w-6 tw-justify-center tw-items-center tw-rounded tw-flex tw-my-auto tw-select-none"
            alt="fabra logo"
          />
          <div className="tw-my-auto tw-ml-2.5 tw-max-w-[150px] tw-whitespace-nowrap tw-overflow-hidden tw-select-none tw-font-bold tw-font-[Montserrat] tw-text-2xl">
            fabra
          </div>
        </NavLink>
        <div className="tw-mx-4 tw-mb-2 tw-uppercase tw-text-xs tw-text-slate-500 tw-font-medium">Overview</div>
        <div className={routeContainer}>
          <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to={"/"}>
            <HomeIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Home</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to={"/syncs"}>
            <ArrowPathIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Syncs</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <NavLink
            className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")}
            to={"/notifications"}
          >
            <BellIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Notifications</div>
          </NavLink>
        </div>

        <div className="tw-my-5 tw-px-4">
          <div className="tw-border-b tw-border-solid tw-border-slate-300" />
        </div>
        <div className="tw-mx-4 tw-my-2 tw-uppercase tw-text-xs tw-text-slate-500 tw-font-medium">Develop</div>
        <div className={routeContainer}>
          <NavLink
            className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")}
            to={"/destinations"}
          >
            <CircleStackIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Destinations</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to={"/objects"}>
            <CubeIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Objects</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to={"/preview"}>
            <EyeIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Preview</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to={"/apikey"}>
            <KeyIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>API Keys</div>
          </NavLink>
        </div>
        <div className={routeContainer}>
          <a className={navLink} href="https://docs.fabra.io/" target="_blank" rel="noreferrer">
            <MapIcon className="tw-h-4" strokeWidth="2" />
            <div className={route}>Documentation</div>
            <ArrowTopRightOnSquareIcon className="tw-h-4 tw-ml-auto tw-mr-3" />
          </a>
        </div>
        <div id="bottomSection" className="tw-mt-auto tw-mb-4">
          <div className="tw-mx-4 tw-mb-2 tw-uppercase tw-text-xs tw-text-slate-500 tw-font-medium">Account</div>
          <div className={routeContainer}>
            <NavLink className={({ isActive }) => classNames(navLink, isActive && "tw-bg-slate-200")} to="/team">
              <UserPlusIcon className="tw-h-4 tw-ml-[1px] -tw-mr-[0.5px]" strokeWidth="2" />
              <div className={route}>Team</div>
            </NavLink>
          </div>
          <div className={routeContainer}>
            <a className={navLink} href="mailto:nick@fabra.io?subject=Help with Fabra">
              <QuestionCircleIcon className="tw-h-[18px] tw-mt-[1px]" strokeWidth="2" />
              <div className={route}>Help</div>
            </a>
          </div>
        </div>
      </div>
    </>
  );
};
