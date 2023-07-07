import { mergeClasses } from "src/utils/twmerge";

export function SectionLayout({ children, className }: { children: React.ReactNode; className?: string }) {
  return (
    <div
      className={mergeClasses(
        "tw-ring-1 tw-ring-black tw-ring-opacity-5 tw-bg-white tw-rounded-lg tw-overflow-x-auto tw-overscroll-contain tw-shadow-md",
        className,
      )}
    >
      {children}
    </div>
  );
}
