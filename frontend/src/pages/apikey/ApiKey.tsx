import { Callout } from "src/components/callouts/Callout";
import { PrivateKey } from "src/components/privateKey/PrivateKey";
import { useApiKey } from "src/rpc/data";

export const ApiKey: React.FC = () => {
  const { apiKey } = useApiKey();

  return (
    <div className="tw-py-5 tw-px-10">
      <div className="tw-flex tw-w-full tw-mt-2 tw-mb-3">
        <div className="tw-flex tw-flex-col tw-justify-end">
          <span className="tw-font-bold tw-text-lg">API Keys</span>
          <div className="tw-mt-2 tw-text-sm">
            Use this API key to authenticate your requests to the Fabra API.{" "}
            <a
              className="tw-text-blue-400"
              href="https://docs.fabra.io/concepts/authentication"
              target="_blank"
              rel="noreferrer"
            >
              Learn more.
            </a>
          </div>
        </div>
      </div>
      <Callout
        className="tw-mb-5 tw-mt-4 tw-max-w-lg"
        content="Never store this secret in plaintext!"
        tooltip="We recommend using a secrets manager to load this API key at runtime in your application."
      />
      <PrivateKey keyValue={apiKey} />
    </div>
  );
};
