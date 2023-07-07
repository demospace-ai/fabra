import { sendLinkTokenRequest, sendRequest } from "src/rpc/ajax";
import {
  ConnectionType,
  GetAllUsers,
  GetAllUsersResponse,
  GetApiKey,
  GetDestination,
  GetDestinationResponse,
  GetDestinations,
  GetDestinationsResponse,
  GetFieldValues,
  GetFieldValuesRequest,
  GetFieldValuesResponse,
  GetNamespaces,
  GetNamespacesResponse,
  GetObject,
  GetObjectResponse,
  GetObjects,
  GetObjectsResponse,
  GetSchema,
  GetSchemaRequest,
  GetSchemaResponse,
  GetSourcesResponse,
  GetSync,
  GetSyncResponse,
  GetSyncs,
  GetSyncsResponse,
  GetTables,
  GetTablesResponse,
  LinkGetNamespaces,
  LinkGetSchema,
  LinkGetSources,
  LinkGetSync,
  LinkGetSyncs,
  LinkGetTables,
} from "src/rpc/api";
import useSWR, { Fetcher } from "swr";

export function useApiKey() {
  const fetcher: Fetcher<string, {}> = () => sendRequest(GetApiKey);
  const { data, error, isLoading, isValidating } = useSWR({ GetApiKey }, fetcher);
  return { apiKey: data, error, loading: isLoading || isValidating };
}

export function useSchema(
  connectionID: number | undefined,
  namespace?: string,
  tableName?: string,
  customJoin?: string,
) {
  const fetcher: Fetcher<GetSchemaResponse, GetSchemaRequest> = (request: GetSchemaRequest) =>
    sendRequest(GetSchema, request);
  const shouldFetch = connectionID && ((namespace && tableName) || customJoin);
  const { data, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetSchema, connectionID, namespace, tableName, customJoin } : null,
    fetcher,
  );
  return { schema: data?.schema, error, loading: isLoading || isValidating };
}

export function useUsers() {
  const fetcher: Fetcher<GetAllUsersResponse, {}> = () => sendRequest(GetAllUsers);
  const { data, mutate, error, isLoading, isValidating } = useSWR({ GetAllUsers }, fetcher);
  return { users: data?.users, mutate, error, loading: isLoading || isValidating };
}

export function useDestinations() {
  const fetcher: Fetcher<GetDestinationsResponse, {}> = () => sendRequest(GetDestinations);
  const { data, mutate, error, isLoading, isValidating } = useSWR({ GetDestinations }, fetcher);
  return { destinations: data?.destinations, mutate, error, loading: isLoading || isValidating };
}

export function useDestination(destinationID: number | undefined) {
  const shouldFetch = destinationID;
  const fetcher: Fetcher<GetDestinationResponse, { destinationID: number }> = (payload: { destinationID: number }) =>
    sendRequest(GetDestination, payload);
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetDestination, destinationID } : null,
    fetcher,
  );
  return { destination: data?.destination, mutate, error, loading: isLoading || isValidating };
}

export function useObjects(opts: { linkToken?: string; destinationID?: number } = {}) {
  const { linkToken, destinationID } = opts;
  let fetchFn;
  if (linkToken) {
    fetchFn = () => sendLinkTokenRequest(GetObjects, linkToken);
  } else {
    fetchFn = () => sendRequest(GetObjects, { destinationID });
  }

  const fetcher: Fetcher<GetObjectsResponse, {}> = fetchFn;
  const { data, mutate, error, isLoading, isValidating } = useSWR({ GetObjects, destinationID }, fetcher);
  return { objects: data?.objects, mutate, error, loading: isLoading || isValidating };
}

export function useObject(objectID: number | undefined, linkToken?: string) {
  let fetchFn;
  if (linkToken) {
    fetchFn = (payload: { objectID: number }) => sendLinkTokenRequest(GetObject, linkToken, payload);
  } else {
    fetchFn = (payload: { objectID: number }) => sendRequest(GetObject, payload);
  }

  const shouldFetch = objectID;
  const fetcher: Fetcher<GetObjectResponse, { objectID: number }> = fetchFn;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetObject, objectID } : null,
    fetcher,
  );
  return { object: data?.object, mutate, error, loading: isLoading || isValidating };
}

export function useNamespaces(connectionID: number | undefined) {
  const fetcher: Fetcher<GetNamespacesResponse, { connectionID: number }> = (payload: { connectionID: number }) =>
    sendRequest(GetNamespaces, payload);
  const shouldFetch = connectionID;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetNamespaces, connectionID } : null,
    fetcher,
  );
  return { namespaces: data?.namespaces, mutate, error, loading: isLoading || isValidating };
}

export function useTables({
  connectionID,
  connectionType,
  namespace,
}: {
  connectionID: number | undefined;
  connectionType: ConnectionType;
  namespace?: string | undefined;
}) {
  const fetcher: Fetcher<GetTablesResponse, { connectionID: number; namespace: string }> = (payload: {
    connectionID: number;
    namespace: string;
  }) => sendRequest(GetTables, payload);
  let shouldFetch = false;
  if (connectionType === ConnectionType.DynamoDb) {
    shouldFetch = !!connectionID;
  } else {
    shouldFetch = !!(connectionID && namespace);
  }

  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetTables, connectionID, namespace } : null,
    fetcher,
  );
  return { tables: data?.tables, mutate, error, loading: isLoading || isValidating };
}
export function useSyncs() {
  const fetcher: Fetcher<GetSyncsResponse, {}> = () => sendRequest(GetSyncs);
  const { data, mutate, error, isLoading, isValidating } = useSWR({ GetSyncs }, fetcher);
  return {
    syncs: data?.syncs,
    objects: data?.objects,
    sources: data?.sources,
    mutate,
    error,
    loading: isLoading || isValidating,
  };
}

export function useSync(syncID: number | undefined) {
  const shouldFetch = !!syncID;
  const fetcher: Fetcher<GetSyncResponse, { syncID: number }> = (payload: { syncID: number }) =>
    sendRequest(GetSync, payload);
  const { data, mutate, error, isLoading, isValidating } = useSWR(shouldFetch ? { GetSync, syncID } : null, fetcher, {
    refreshInterval: 1000,
  });
  return { sync: data, mutate, error, loading: isLoading || isValidating };
}

export function useLinkNamespaces(sourceID: number | undefined, linkToken: string) {
  const fetcher: Fetcher<GetNamespacesResponse, { sourceID: number }> = (payload: { sourceID: number }) =>
    sendLinkTokenRequest(LinkGetNamespaces, linkToken, payload);
  const shouldFetch = sourceID;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { LinkGetNamespaces, sourceID } : null,
    fetcher,
  );
  return { namespaces: data?.namespaces, mutate, error, loading: isLoading || isValidating };
}

export function useLinkTables(sourceID: number | undefined, namespace: string | undefined, linkToken: string) {
  const fetcher: Fetcher<GetTablesResponse, { sourceID: number; namespace: string }> = (payload: {
    sourceID: number;
    namespace: string;
  }) => sendLinkTokenRequest(LinkGetTables, linkToken, payload);
  const shouldFetch = sourceID && namespace;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { LinkGetTables, sourceID, namespace } : null,
    fetcher,
  );
  return { tables: data?.tables, mutate, error, loading: isLoading || isValidating };
}

export function useLinkSchema(
  sourceID: number | undefined,
  namespace: string | undefined,
  tableName: string | undefined,
  linkToken: string,
) {
  const fetcher: Fetcher<GetSchemaResponse, { sourceID: number; namespace: string; tableName: string }> = (payload: {
    sourceID: number;
    namespace: string;
    tableName: string;
  }) => sendLinkTokenRequest(LinkGetSchema, linkToken, payload);
  const shouldFetch = sourceID && namespace && tableName;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { LinkGetSchema, sourceID, namespace, tableName } : null,
    fetcher,
  );
  return { schema: data?.schema, mutate, error, loading: isLoading || isValidating };
}

export function useLinkSources(linkToken: string) {
  const fetcher: Fetcher<GetSourcesResponse, {}> = () => sendLinkTokenRequest(LinkGetSources, linkToken);
  const { data, mutate, error, isLoading, isValidating } = useSWR({ LinkGetSources }, fetcher);
  return { sources: data?.sources, mutate, error, loading: isLoading || isValidating };
}

export function useLinkSyncs(linkToken: string) {
  const fetcher: Fetcher<GetSyncsResponse, {}> = () => sendLinkTokenRequest(LinkGetSyncs, linkToken);
  const { data, mutate, error, isLoading, isValidating } = useSWR({ LinkGetSyncs }, fetcher);
  return {
    syncs: data?.syncs,
    objects: data?.objects,
    sources: data?.sources,
    mutate,
    error,
    loading: isLoading || isValidating,
  };
}

export function useLinkSync(syncID: number | undefined, linkToken: string) {
  const shouldFetch = syncID;
  const fetcher: Fetcher<GetSyncResponse, { syncID: number }> = (payload: { syncID: number }) =>
    sendLinkTokenRequest(LinkGetSync, linkToken, payload);
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { LinkGetSync: LinkGetSync, syncID } : null,
    fetcher,
  );
  return { sync: data, mutate, error, loading: isLoading || isValidating };
}

export function useFieldValues(
  connectionID: number | undefined,
  namespace: string | undefined,
  tableName: string | undefined,
  fieldName: string | undefined,
) {
  const fetcher: Fetcher<GetFieldValuesResponse, GetFieldValuesRequest> = (payload: GetFieldValuesRequest) =>
    sendRequest(GetFieldValues, payload);
  const shouldFetch = connectionID && namespace && tableName && fieldName;
  const { data, mutate, error, isLoading, isValidating } = useSWR(
    shouldFetch ? { GetFieldValues: GetFieldValues, connectionID, namespace, tableName, fieldName } : null,
    fetcher,
  );
  return { fieldValues: data?.field_values, mutate, error, loading: isLoading || isValidating };
}
