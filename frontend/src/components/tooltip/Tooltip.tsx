import Tippy, { TippyProps } from "@tippyjs/react";
import React from "react";

export interface TooltipProps extends TippyProps {
  children: React.ReactElement;
  label?: React.ReactElement | string;
}

export const Tooltip: React.FC<TooltipProps> = (props) => {
  const { label, ...other } = props;

  return (
    <>
      <Tippy content={props.label} delay={0} duration={100} {...other}>
        {props.children}
      </Tippy>
    </>
  );
};
