import { SVGAttributes } from "react";

interface DataPreviewProps extends SVGAttributes<SVGElement> {
  animate?: boolean;
}

export const DataPreview: React.FC<DataPreviewProps> = (props) => {
  return (
    <svg
      width="450"
      height="450"
      viewBox="0 0 450 450"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={props.className}
    >
      <circle
        cx="225.5"
        cy="225.5"
        r="146.5"
        stroke="var(--color-primary-hover)"
        strokeWidth="13"
        strokeLinecap="round"
        strokeDasharray="30 35"
        className={props.animate ? "tw-animate-spin-slow tw-origin-center" : ""}
      />
      <rect x="250" y="250" width="200" height="200" rx="12" fill="#94a3b8" />
      <rect x="268" y="270" width="116" height="25" rx="8" fill="white" />
      <rect x="294" y="304" width="116" height="25" rx="8" fill="white" />
      <rect x="268" y="338" width="56" height="25" rx="8" fill="white" />
      <rect x="268" y="372" width="84" height="25" rx="8" fill="white" />
      <rect x="326" y="406" width="65" height="25" rx="8" fill="white" />
      <rect width="200" height="200" rx="12" fill="var(--color-primary)" />
      <rect x="18" y="20" width="116" height="25" rx="8" fill="white" />
      <rect x="44" y="54" width="116" height="25" rx="8" fill="white" />
      <rect x="18" y="88" width="56" height="25" rx="8" fill="white" />
      <rect x="18" y="122" width="84" height="25" rx="8" fill="white" />
      <rect x="76" y="156" width="65" height="25" rx="8" fill="white" />
    </svg>
  );
};
