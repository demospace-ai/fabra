import { ArrowRightIcon } from "@heroicons/react/24/outline";
import { useSelector } from "src/root/model";

const DAY_IN_MILLIS = 1000 * 60 * 60 * 24;

export const UpgradeBanner: React.FC = () => {
  const isAuthenticated = useSelector((state) => state.login.authenticated);
  const organization = useSelector((state) => state.login.organization);

  // No upgrade banner until they login and set an organization
  if (!isAuthenticated || !organization) {
    return <></>;
  }

  // If the free trial end value is not populated, it means the organization has upgraded
  if (!organization.free_trial_end) {
    return <></>;
  }

  const freeTrialEnd = Date.parse(organization.free_trial_end);
  const now = Date.now();
  const daysUntilTrialEnds = Math.round((freeTrialEnd - now) / DAY_IN_MILLIS);

  return (
    <a
      className="tw-flex tw-cursor-pointer tw-h-10 tw-justify-center tw-items-center tw-border-b tw-border-solid tw-text-white tw-bg-slate-900 hover:tw-bg-slate-800 tw-transition-colors tw-group"
      href="https://calendly.com/fabra-io/onboarding?month=2023-04"
      target="_blank"
      rel="noreferrer"
    >
      ⏱️ {daysUntilTrialEnds} days left in your free trial. Reach out to upgrade{" "}
      <ArrowRightIcon className="tw-relative tw-h-4 tw-left-1 tw-mt-[1px] group-hover:tw-left-2 tw-transition-all" />
    </a>
  );
};
