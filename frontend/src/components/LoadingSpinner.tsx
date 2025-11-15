export function LoadingSpinner({
  size = "xl",
}: {
  size?: "xl" | "lg" | "md" | "sm" | "xs";
}) {
  return <span className={`loading loading-spinner loading-${size}`}></span>;
}

export function LoadingSpinnerPage() {
  return (
    <div className="fixed flex  inset-0 justify-center">
      <LoadingSpinner />
    </div>
  );
}
