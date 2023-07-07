import { compile } from "path-to-regexp";
import { rudderanalytics } from "src/app/rudder";
import { IEndpoint } from "src/rpc/api";
import { isProd } from "src/utils/env";
import { HttpError } from "src/utils/errors";

const ROOT_DOMAIN = isProd() ? "https://api.fabra.io" : "http://localhost:8080";

export async function sendLinkTokenRequest<RequestType extends Record<string, any>, ResponseType>(
  endpoint: IEndpoint<RequestType, ResponseType>,
  linkToken: string,
  payload?: RequestType,
): Promise<ResponseType> {
  return sendRequest(endpoint, payload, [["X-LINK-TOKEN", linkToken]]);
}

export function getEndpointUrl<RequestType extends Record<string, any>, ResponseType>(
  endpoint: IEndpoint<RequestType, ResponseType>,
  payload?: RequestType,
  extraHeaders?: [string, string][],
): string {
  const toPath = compile(endpoint.path);
  const path = toPath(payload);

  const url = new URL(ROOT_DOMAIN + path);
  if (endpoint.queryParams && payload) {
    endpoint.queryParams.forEach((queryParam) => {
      const queryParamValue = payload[queryParam];
      if (queryParamValue) {
        url.searchParams.append(queryParam, queryParamValue);
      }
    });
  }

  return url.toString();
}

export async function sendRequest<RequestType extends Record<string, any>, ResponseType>(
  endpoint: IEndpoint<RequestType, ResponseType>,
  payload?: RequestType,
  extraHeaders?: [string, string][],
): Promise<ResponseType> {
  const toPath = compile(endpoint.path);
  const path = toPath(payload);

  const url = new URL(ROOT_DOMAIN + path);
  if (endpoint.queryParams && payload) {
    endpoint.queryParams.forEach((queryParam) => {
      const queryParamValue = payload[queryParam];
      if (queryParamValue) {
        url.searchParams.append(queryParam, queryParamValue);
      }
    });
  }

  const extraHeadersList = extraHeaders ? extraHeaders : [];
  const headers = new Headers([
    ["Content-Type", "application/json"],
    ["X-TIME-ZONE", Intl.DateTimeFormat().resolvedOptions().timeZone],
    ...extraHeadersList,
  ]);
  let options: RequestInit = {
    method: endpoint.method,
    headers: headers,
    credentials: "include",
  };

  if (["POST", "PATCH", "PUT"].includes(endpoint.method)) {
    options.body = JSON.stringify(payload);
  }

  const response = await fetch(url, options);

  if (!response.ok) {
    const errorMessage = response.statusText + ": " + (await response.text());
    throw new HttpError(response.status, response.statusText, errorMessage);
  }

  if (endpoint.track) {
    rudderanalytics.track(`${endpoint.name}`);
  }

  // TODO: clean this up
  // not all AJAX requests have a response. the ones that do will be formatted as JSON
  // so just catch any error from trying to fetch the json and do nothing with it
  if (endpoint.noJson) {
    return response.text() as ResponseType;
  }
  return response.json().catch(() => null);
}
