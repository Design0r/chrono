import type React from "react";

type TitleSectionProps = {
  children?: React.ReactNode[] | React.ReactNode | undefined;
  title: string;
};

export function TitleSection({ children, title }: TitleSectionProps) {
  return (
    <div className="overflow-hidden my-4 p-3 rounded-xl flex flex-col items-left justify-center">
      <h1 className="text-xl">{title}</h1>
      <div className="border-b border-primary/25 mt-2 mb-8"></div>
      {children &&
        (Array.isArray(children) ? (
          <div>{...children}</div>
        ) : (
          <div>{children}</div>
        ))}
    </div>
  );
}
