import "src/index.css";

import { ErrorBoundary } from "@highlight-run/react";
import { H } from "highlight.run";
import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { RouterProvider } from "react-router-dom";
import { router } from "src/app/App";
import { rudderanalytics } from "src/app/rudder";
import { createStore } from "src/root/model";
import { isProd } from "src/utils/env";

const store = createStore();

if (isProd()) {
  rudderanalytics.load("2DuH7iesuV4TtpwMqRvXqQttOvm", "https://fabranickbele.dataplane.rudderstack.com");
  H.init("7e3vw5g1", {
    tracingOrigins: ["localhost", "api.fabra.io"],
    networkRecording: {
      enabled: true,
      recordHeadersAndBody: true,
      networkHeadersToRedact: ["X-LINK-TOKEN"],
      networkBodyKeysToRedact: ["link_token"],
    },
  });
  window.Intercom("boot", {
    api_base: "https://api-iam.intercom.io",
    app_id: "pdc06iv8",
  });
} else {
  H.init("6gl37zg9", {
    tracingOrigins: ["localhost"],
    networkRecording: {
      enabled: true,
      recordHeadersAndBody: true,
      networkHeadersToRedact: ["X-LINK-TOKEN"],
      networkBodyKeysToRedact: ["link_token"],
    },
  });
}

const container = document.getElementById("root");
const root = createRoot(container!);

root.render(
  <React.StrictMode>
    <ErrorBoundary showDialog>
      <Provider store={store}>
        <RouterProvider router={router} />
      </Provider>
    </ErrorBoundary>
  </React.StrictMode>,
);
