import { H } from "highlight.run";
import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { rudderanalytics } from "src/app/rudder";
import { useDispatch } from "src/root/model";
import { sendRequest } from "src/rpc/ajax";
import { Logout, Organization, SetOrganization, User } from "src/rpc/api";
import { isProd } from "src/utils/env";
import { consumeError } from "../../utils/errors";

export interface OrganizationArgs {
  organizationName?: string;
  organizationID?: number;
}

export function useSetOrganization() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const onLoginSuccess = useOnLoginSuccess();

  return useCallback(
    async (user: User, args: OrganizationArgs) => {
      const payload = { organization_name: args.organizationName, organization_id: args.organizationID };
      try {
        const response = await sendRequest(SetOrganization, payload);
        dispatch({
          type: "login.organizationSet",
          organization: response.organization,
        });

        onLoginSuccess(user, response.organization);
        navigate("/");
      } catch (e) {
        consumeError(e);
      }
    },
    [dispatch, navigate, onLoginSuccess],
  );
}

export function useOnLoginSuccess() {
  const navigate = useNavigate();

  return useCallback(
    async (user: User, organization: Organization | undefined) => {
      identifyUser(user);

      // If there's no organization, go to the login page so the user can set it
      if (!organization) {
        navigate("/login");
        return;
      }
    },
    [navigate],
  );
}

function identifyUser(user: User) {
  if (isProd()) {
    rudderanalytics.identify(user.id.toString(), {
      name: `${user.name}`,
      email: user.email,
    });

    H.identify(user.email, {
      id: user.id.toString(),
    });

    window.Intercom("boot", {
      api_base: "https://api-iam.intercom.io",
      app_id: "pdc06iv8",
      name: user.name,
      email: user.email,
      user_id: user.id,
      user_hash: user.intercom_hash,
    });
  }
}

export function useLogout() {
  const dispatch = useDispatch();

  return useCallback(async () => {
    rudderanalytics.reset();
    window.Intercom("shutdown");

    await sendRequest(Logout);
    dispatch({
      type: "login.logout",
    });
  }, [dispatch]);
}
