import { Link } from "@tanstack/react-router";

export function ErrorPage({
  error,
}: {
  error: { name: string; message: string };
}) {
  return (
    <div className="fixed inset-0 gap-5 flex items-center flex-col justify-center text-center">
      <h1 className="text-4xl text-error">Error {error.name}</h1>
      <p className="text-xl">🚨 Whoops, looks like something went wrong 🚨</p>
      <p className="text-lg">{error.message}</p>
      <Link to="/" className="btn max-w-30 btn-primary">
        Home
      </Link>
    </div>
  );
}
