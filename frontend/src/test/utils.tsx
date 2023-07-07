import { render } from "@testing-library/react";
import { ReactElement } from "react";
import { Provider } from "react-redux";
import createMockStore, { MockStore } from "redux-mock-store";

export const mockStore = createMockStore();

export function renderWithProvider<Props>(node: ReactElement<Props>, store: MockStore) {
  render(<Provider store={store}>{node}</Provider>);
}
