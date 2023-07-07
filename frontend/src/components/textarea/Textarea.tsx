import { useEffect, useRef } from "react";

export const ExpandingTextarea: React.FC<
  React.DetailedHTMLProps<React.TextareaHTMLAttributes<HTMLTextAreaElement>, HTMLTextAreaElement>
> = (props) => {
  const ref = useRef<HTMLTextAreaElement>(null);
  useEffect(() => {
    if (ref.current) {
      ref.current.style.height = "0px";
      ref.current.style.height = ref.current.scrollHeight + "px";
    }
  }, [props.value]);

  return <textarea {...props} ref={ref} />;
};
