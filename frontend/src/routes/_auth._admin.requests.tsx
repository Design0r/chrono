import { createFileRoute } from "@tanstack/react-router";
import { RequestRow } from "../components/Requests";
import { useQuery } from "@tanstack/react-query";
import { ErrorPage } from "../components/ErrorPage";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";

export const Route = createFileRoute("/_auth/_admin/requests")({
  component: RouteComponent,
});

function RouteComponent() {
  const { chrono } = Route.useRouteContext();

  const { isPending, data, error, isError } = useQuery({
    queryKey: ["requests"],
    queryFn: () => chrono.requests.getRequests(),
  });

  if (isError) return <ErrorPage error={error} />;
  if (isPending) return <LoadingSpinnerPage />;

  return (
    <div className="p-2 my-2">
      <div className="overflow-x-auto bg-info/3">
        <table className="table table-zebra table-md ">
          <thead>
            <tr className="bg-base-200 text-base-content border-b-1.5 border-primary/65">
              <th>ID</th>
              <th>Name</th>
              <th>Event Type</th>
              <th>Number of Days</th>
              <th>Event Start Date</th>
              <th>Event End Date</th>
              <th>Conflicts</th>
              <th className="text-center">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data!.map((r) => (
              <RequestRow key={r.request.request_id} request={r} />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
