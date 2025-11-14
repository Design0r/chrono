export function LoadingSpinner() {
  return <span className="loading loading-spinner loading-xl"></span>;
}

export function LoadingSpinnerPage() {
  return (
    <div className="fixed flex  inset-0 justify-center">
      <LoadingSpinner />
    </div>
  );
}
