import { useEffect, useState } from "react";
import { Button } from "src/components/button/Button";
import { ErrorDisplay } from "src/components/error/Error";
import { InfoIcon } from "src/components/icons/Icons";
import { ColorPicker, Input } from "src/components/input/Input";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { FabraConnectOptions, initialize, open, updateTheme } from "src/initialize-internal";
import { sendRequest } from "src/rpc/ajax";
import { CreateLinkToken, CreateLinkTokenRequest } from "src/rpc/api";
import { useForm } from "src/utils/formHelpers";
import { useMutation } from "src/utils/queryHelpers";

export const Preview: React.FC = () => {
  const [currCustomerId, setCurrCustomerId] = useState<string>("");
  const [linkToken, setLinkToken] = useState<string | undefined>(undefined);

  const previewForm = useForm(
    {
      endCustomerID: "",
      baseColor: "#475569",
      hoverColor: "#1e293b",
      textColor: "#ffffff",
    },
    {
      validate: (values) => {
        const errors: Record<string, string> = {};
        for (const [key, value] of Object.entries(values)) {
          if (value === "") {
            errors[key] = "This field is required";
          }
        }
        return errors;
      },
    },
  );

  // Hack to update the colors of the active iFrame
  useEffect(() => {
    updateTheme({
      colors: {
        primary: {
          base: previewForm.state.baseColor,
          hover: previewForm.state.hoverColor,
          text: previewForm.state.textColor,
        },
      },
    });
  }, [previewForm.state.baseColor, previewForm.state.hoverColor, previewForm.state.textColor]);

  useFabraConnect({
    containerID: "fabra-container",
    customTheme: {
      colors: {
        primary: {
          base: previewForm.state.baseColor,
          hover: previewForm.state.hoverColor,
          text: previewForm.state.textColor,
        },
      },
    },
  });

  const openPreviewMutation = useMutation(
    async () => {
      const payload: CreateLinkTokenRequest = {
        end_customer_id: previewForm.state.endCustomerID,
      };
      const response = await sendRequest(CreateLinkToken, payload);
      return response.link_token;
    },
    {
      onSuccess: (link_token) => {
        setLinkToken(link_token);
        setCurrCustomerId(previewForm.state.endCustomerID);
        open(link_token);
      },
    },
  );

  const onSubmit = previewForm.handleSubmit(() => {
    openPreviewMutation.mutate();
  });

  return (
    <div className="tw-py-5 tw-px-10 tw-flex tw-w-full tw-h-full tw-flex-col xl:tw-flex-row">
      <div className="xl:tw-w-1/4 tw-mb-4 xl:tw-mb-0 xl:tw-mr-4">
        <div className="tw-flex tw-w-full tw-mt-2 tw-mb-3">
          <h2 className="tw-flex tw-flex-col tw-justify-end tw-font-bold tw-text-lg">Fabra Connect</h2>
        </div>
        <div>
          See what Fabra Connect looks like for your end customers. Enter a test end customer ID and click Preview.
        </div>
        <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1 tw-font-medium">
          <span>End Customer ID</span>
          <Tooltip
            placement="right"
            label="This can be any string. If you use an actual ID for one of your users, you can see what that user will see."
          >
            <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
          </Tooltip>
        </div>
        <form onSubmit={onSubmit}>
          <Input
            className="tw-flex-1"
            value={previewForm.state.endCustomerID}
            setValue={(v) => {
              previewForm.setState((s) => ({ ...s, endCustomerID: v }));
            }}
            placeholder="143"
          />
          {previewForm.errors.endCustomerID && (
            <div className="tw-text-red-500">{previewForm.errors.endCustomerID}</div>
          )}
          <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1 tw-font-medium">
            <span>Base Color</span>
          </div>
          <ColorPicker
            value={previewForm.state.baseColor}
            setValue={(v) => {
              previewForm.setState((s) => ({ ...s, baseColor: v }));
            }}
            placeholder="Base Color (optional)"
          />
          <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1 tw-font-medium">
            <span>Hover Color</span>
          </div>
          <ColorPicker
            value={previewForm.state.hoverColor}
            setValue={(v) => {
              previewForm.setState((s) => ({ ...s, hoverColor: v }));
            }}
            placeholder="Hover Color (optional)"
          />
          <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1 tw-font-medium">
            <span>Text Color</span>
          </div>
          <ColorPicker
            value={previewForm.state.textColor}
            setValue={(v) => {
              previewForm.setState((s) => ({ ...s, textColor: v }));
            }}
            placeholder="Text Color (optional)"
          />
          <div>
            <Button className="tw-px-4 tw-mt-6 tw-py-2 tw-w-full" type="submit">
              {currCustomerId !== previewForm.state.endCustomerID &&
              previewForm.state.endCustomerID.length > 0 &&
              linkToken
                ? "Change Test ID"
                : "Preview"}
            </Button>
            <ErrorDisplay error={openPreviewMutation.error} className="tw-text-red-500" />
          </div>
        </form>
      </div>
      <div
        id="fabra-container"
        className="tw-w-full tw-h-full tw-border tw-border-slate-200 tw-rounded-md tw-overflow-clip"
      />
    </div>
  );
};

// Slightly customized version of ReactFabraConnect to use local Connect code in development
const useFabraConnect = (options?: FabraConnectOptions) => {
  useEffect(() => {
    initialize(options);
  }, []);
};
