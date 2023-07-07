import { InfoIcon } from "src/components/icons/Icons";
import { ValidatedInput } from "src/components/input/Input";
import { AwsLocationSelector } from "src/components/selector/Selector";
import { Tooltip } from "src/components/tooltip/Tooltip";
import { NewDestinationState } from "src/pages/destinations/helpers";

export function DynamoDbInputs({
  setState,
  state,
}: {
  setState: (state: NewDestinationState) => void;
  state: NewDestinationState;
}) {
  return (
    <>
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-2 tw-mb-1">
        <span>Display Name</span>
        <Tooltip placement="right" label="Pick a name to help you identify this source in the future.">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="displayName"
        value={state.displayName}
        setValue={(value) => {
          setState({ ...state, displayName: value });
        }}
        placeholder="Display Name"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Region</span>
        <Tooltip placement="right" label="The geographic location of your BigQuery dataset(s).">
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <AwsLocationSelector
        id="location"
        location={state.dynamoDbConfig.region}
        setLocation={(value) => {
          setState({ ...state, dynamoDbConfig: { ...state.dynamoDbConfig, region: value } });
        }}
        placeholder="Location"
        className="tw-mt-0 tw-w-full"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Access Key ID</span>
        <Tooltip
          placement="right"
          interactive
          label={
            <div>
              You can create an access key for a particular user in the AWS console. Fabra requires Read and Write
              permissions to the DynamoDB table. Find more information{" "}
              <a
                className="tw-text-blue-300"
                href="https://docs.aws.amazon.com/IAM/latest/UserGuide/security-creds.html#access-keys-and-secret-access-keys"
                target="_blank"
                rel="noreferrer"
              >
                here
              </a>
              .
            </div>
          }
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="accessKey"
        value={state.dynamoDbConfig.accessKey}
        setValue={(value) => {
          setState({ ...state, dynamoDbConfig: { ...state.dynamoDbConfig, accessKey: value } });
        }}
        placeholder="Access Key"
      />
      <div className="tw-flex tw-flex-row tw-items-center tw-mt-4 tw-mb-1">
        <span>Secret Access Key</span>
        <Tooltip
          placement="right"
          label="This can be obtained in the AWS console after creating your access key."
          interactive
          maxWidth={500}
        >
          <InfoIcon className="tw-ml-1 tw-h-3 tw-fill-slate-400" />
        </Tooltip>
      </div>
      <ValidatedInput
        id="secretKey"
        value={state.dynamoDbConfig.secretKey}
        setValue={(value) => {
          setState({ ...state, dynamoDbConfig: { ...state.dynamoDbConfig, secretKey: value } });
        }}
        placeholder="Secret key"
        type="password"
      />
    </>
  );
}
