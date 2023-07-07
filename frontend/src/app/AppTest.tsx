import { screen } from "@testing-library/react";
import { App } from "src/app/App";
import { mockStore, renderWithProvider } from "src/test/utils";

test("renders loading when loading", async () => {
  const store = mockStore({
    app: {
      loading: true,
    },
  });

  renderWithProvider(<App />, store);

  const loadingElement = screen.getByTestId("loading");
  expect(loadingElement).toBeInTheDocument();
});

test("does not render loading when not loading", async () => {
  const store = mockStore({
    app: {
      loading: false,
    },
    login: {
      authenticated: false,
    },
  });

  renderWithProvider(<App />, store);

  const loadingElement = screen.queryByTestId("loading");
  expect(loadingElement).toBeNull();
});
