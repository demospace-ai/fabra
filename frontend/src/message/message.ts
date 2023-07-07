import { CustomTheme } from "src/utils/theme";

export enum MessageType {
  IFrameReady = "fabra-iframe-ready",
  LinkToken = "fabra-link-token",
  Configure = "fabra-configure",
  Close = "fabra-window-close",
}

export type FabraMessage =
  | {
      messageType: MessageType.IFrameReady;
    }
  | {
      messageType: MessageType.LinkToken;
      linkToken: string;
    }
  | {
      messageType: MessageType.Close;
    }
  | {
      messageType: MessageType.Configure;
      theme: CustomTheme | undefined;
      useContainer: boolean;
      supportEmail: string | undefined;
      docsLink: string | undefined;
    };
