import "src/index.css";

import { ErrorBoundary } from "@highlight-run/react";
import { H } from "highlight.run";
import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import { rudderanalytics } from "src/app/rudder";
import { ConnectApp } from "src/connect/ConnectApp";
import { createConnectStore } from "src/connect/model";
import { isProd } from "src/utils/env";

const store = createConnectStore();

// Initialize Highlight tracking for production
if (isProd()) {
  rudderanalytics.load("2NhFGB9sihlA85YlkD7p2kSuen4", "https://fabranickbele.dataplane.rudderstack.com");
  H.init("6gllomg9", {
    tracingOrigins: ["localhost", "api.fabra.io"],
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
        <BrowserRouter>
          <ConnectApp />
        </BrowserRouter>
      </Provider>
    </ErrorBoundary>
  </React.StrictMode>,
);
