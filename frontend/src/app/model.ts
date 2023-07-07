import { ToastOptions } from "src/components/notifications/Notifications";

export type AppAction =
  | {
      type: "loading";
    }
  | {
      type: "done";
    }
  | {
      type: "forbidden";
    }
  | {
      type: "toast";
      toast?: ToastOptions;
    };

const INITIAL_APP_STATE: AppState = {
  loading: true,
  forbidden: false,
  toast: undefined,
};

export interface AppState {
  loading: boolean;
  forbidden: boolean;
  toast?: ToastOptions;
}

export function appReducer(state: AppState = INITIAL_APP_STATE, action: AppAction): AppState {
  switch (action.type) {
    case "loading":
      return {
        ...state,
        loading: true,
      };
    case "done":
      return {
        ...state,
        loading: false,
      };
    case "forbidden":
      return {
        ...state,
        forbidden: true,
      };
    case "toast":
      return {
        ...state,
        toast: action.toast,
      };
    default:
      return state;
  }
}
