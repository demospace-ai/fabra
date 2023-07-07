import classNames from "classnames";
import React, { FormEvent, useEffect, useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Button, FormButton } from "src/components/button/Button";
import { GithubIcon } from "src/components/icons/Github";
import { GoogleIcon } from "src/components/icons/Google";
import longlogo from "src/components/images/long-logo.svg";
import mail from "src/components/images/mail.svg";
import { LogoLoading } from "src/components/loading/LogoLoading";
import { useSetOrganization } from "src/pages/login/actions";
import { useSelector } from "src/root/model";
import { getEndpointUrl } from "src/rpc/ajax";
import { OAuthProvider, OAuthRedirect } from "src/rpc/api";

export enum LoginStep {
  Start = 1,
  ValidateCode,
  Organization,
}

export const Login: React.FC<{ create?: boolean }> = ({ create }) => {
  const [loading, setLoading] = useState(false);
  const isAuthenticated = useSelector((state) => state.login.authenticated);
  const organization = useSelector((state) => state.login.organization);
  const navigate = useNavigate();

  // Use effect to navigate after render if authenticated
  useEffect(() => {
    let ignore = false;
    if (isAuthenticated && organization && !ignore) {
      navigate("/");
    }

    return () => {
      ignore = true;
    };
  }, [navigate, isAuthenticated, organization]);

  if (loading) {
    return <LogoLoading />;
  }

  let loginContent;
  if (!isAuthenticated) {
    loginContent = <StartContent create={create} />;
  } else if (!organization) {
    loginContent = <OrganizationInput setLoading={setLoading} />;
  }

  return (
    <div className="tw-flex tw-flex-row tw-h-full tw-bg-slate-100">
      <div className="tw-mt-56 tw-mb-auto tw-mx-auto tw-w-[400px]">
        <div className="tw-flex tw-flex-col tw-pt-12 tw-pb-10 tw-px-8 tw-rounded-lg tw-shadow-md tw-bg-white tw-items-center">
          <img src={longlogo} className="tw-h-8 tw-select-none tw-mb-4" alt="fabra logo" />
          <div className="tw-text-center tw-my-2">{loginContent}</div>
        </div>
        <div className="tw-text-xs tw-text-center tw-mt-4 tw-text-slate-800 tw-select-none">
          By continuing you agree to Fabra's{" "}
          <a className="tw-text-blue-500" href="https://fabra.io/terms" target="_blank" rel="noreferrer">
            Terms of Use
          </a>{" "}
          and{" "}
          <a className="tw-text-blue-500" href="https://fabra.io/privacy" target="_blank" rel="noreferrer">
            Privacy Policy
          </a>
          .
        </div>
      </div>
    </div>
  );
};

const StartContent: React.FC<{ create?: boolean }> = ({ create }) => {
  const loginError = useSelector((state) => state.login.error);
  return (
    <>
      {loginError && <div className="tw-text-red-500">{loginError?.toString()}</div>}
      <div className="tw-text-center tw-mb-6 tw-select-none">
        {create ? "Start your free 30-day trial of Fabra!" : "Sign in to continue to Fabra."}
      </div>
      <a
        className={classNames(
          "tw-flex tw-items-center tw-select-none tw-cursor-pointer tw-justify-center tw-mt-4 tw-h-10 tw-bg-slate-100 tw-border tw-border-slate-300 hover:tw-bg-slate-200 tw-transition-colors tw-font-medium tw-w-80 tw-text-slate-800 tw-rounded",
        )}
        href={getEndpointUrl(OAuthRedirect, { provider: OAuthProvider.Google })}
      >
        <GoogleIcon className="tw-mr-1.5 tw-h-[18px]" />
        Continue with Google
      </a>
      <a
        className={classNames(
          "tw-flex tw-items-center tw-select-none tw-cursor-pointer tw-justify-center tw-mt-4 tw-h-10 tw-bg-black hover:tw-bg-[#333333] tw-transition-colors tw-font-medium tw-w-80 tw-text-white tw-rounded",
        )}
        href={getEndpointUrl(OAuthRedirect, { provider: OAuthProvider.Github })}
      >
        <GithubIcon className="tw-mr-2" />
        Continue with Github
      </a>
      {create ? (
        <div className="tw-mt-5  tw-select-none">
          Already have an account?{" "}
          <NavLink className="tw-text-blue-500" to="/login">
            Sign in
          </NavLink>
        </div>
      ) : (
        <div className="tw-mt-5 tw-select-none">
          Need an account?{" "}
          <NavLink className="tw-text-blue-500" to="/signup">
            Sign up
          </NavLink>
        </div>
      )}
    </>
  );
};

type OrganizationInputProps = {
  setLoading: (loading: boolean) => void;
};

const OrganizationInput: React.FC<OrganizationInputProps> = (props) => {
  const user = useSelector((state) => state.login.user);
  const suggestedOrganizations = useSelector((state) => state.login.suggestedOrganizations);
  const [organizationInput, setOrganizationInput] = useState("");
  const setOrganization = useSetOrganization();
  const [isValid, setIsValid] = useState(true);
  const [overrideCreate, setOverrideCreate] = useState(false);

  const classes = ["tw-border tw-border-slate-400 tw-rounded-md tw-px-3 tw-py-2 tw-w-full tw-box-border"];
  if (!isValid) {
    classes.push("tw-border-red-500 tw-outline-none");
  }

  const validateOrganization = (): boolean => {
    const valid = organizationInput.length > 0;
    setIsValid(valid);
    return valid;
  };

  const onKeydown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    event.stopPropagation();
    if (event.key === "Escape") {
      event.currentTarget.blur();
    }
  };

  const createNewOrganization = async (e: FormEvent) => {
    e.preventDefault();
    props.setLoading(true);
    if (!validateOrganization()) {
      return;
    }

    await setOrganization(user!, { organizationName: organizationInput });
    props.setLoading(false);
  };

  const joinOrganization = async (organizationID: number) => {
    props.setLoading(true);
    // TODO how to specify positional arg with name
    await setOrganization(user!, { organizationID: organizationID });
  };

  if (!suggestedOrganizations || suggestedOrganizations.length === 0 || overrideCreate) {
    return (
      <form className="tw-mt-5" onSubmit={createNewOrganization}>
        <div className="tw-mb-5">Welcome, {user!.name}! Let's build out your team.</div>
        <input
          type="text"
          id="organization"
          name="organization"
          autoComplete="organization"
          placeholder="Organization Name"
          className={classNames(classes)}
          onKeyDown={onKeydown}
          onFocus={() => setIsValid(true)}
          onChange={(e) => setOrganizationInput(e.target.value)}
          onBlur={validateOrganization}
        />
        {!isValid && (
          <div className="tw-text-red-500 tw-mt-1 -tw-mb-1 tw-text-[15px] tw-text-left">
            Please enter a valid organization name.
          </div>
        )}
        <FormButton className="tw-mt-5 tw-h-10 tw-w-full">Continue</FormButton>
      </form>
    );
  }

  return (
    <div className="tw-mt-5">
      <div className="tw-mb-5">Welcome, {user!.name}! Join your team.</div>
      {suggestedOrganizations.map((suggestion, index) => (
        <li
          key={index}
          className="tw-border tw-border-black tw-rounded-md tw-list-none tw-p-8 tw-text-left tw-flex tw-flex-row"
        >
          <Button className="tw-inline-block tw-mr-8 tw-h-10 tw-w-1/2" onClick={() => joinOrganization(suggestion.id)}>
            Join
          </Button>
          <div className="tw-flex tw-flex-col tw-h-10 tw-w-1/2 tw-text-center tw-justify-center">
            <div className="tw-overflow-hidden tw-text-ellipsis tw-font-bold tw-text-lg">{suggestion.name}</div>
            {/* TODO: add team size */}
          </div>
        </li>
      ))}
      <div className="tw-my-5 tw-mx-0">or</div>
      <Button className="tw-w-full tw-h-10" onClick={() => setOverrideCreate(true)}>
        Create new organization
      </Button>
    </div>
  );
};

export const Unauthorized: React.FC = () => {
  return (
    <div className="tw-flex tw-flex-row tw-h-full tw-bg-slate-100">
      <div className="tw-mt-56 tw-mb-auto tw-mx-auto tw-w-[400px] tw-select-none">
        <div className="tw-flex tw-flex-col tw-pt-12 tw-pb-10 tw-px-8 tw-rounded-lg tw-shadow-md tw-bg-white tw-items-center">
          <img src={longlogo} className="tw-h-8 tw-mb-4" alt="fabra logo" />
          <div className="tw-text-center tw-my-2">
            <div className="tw-flex tw-flex-col tw-justify-center">
              <div>You must use a business account to access Fabra.</div>
              <div className="tw-mt-4">
                <NavLink className="tw-text-blue-500" to="/signup">
                  Try again
                </NavLink>{" "}
                with a different account or{" "}
                <a className="tw-text-blue-500" href="mailto:founders@fabra.io">
                  contact us
                </a>{" "}
                to get an account provisioned!
              </div>
              <img src={mail} alt="mail" className="tw-h-36 tw-mt-5" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
