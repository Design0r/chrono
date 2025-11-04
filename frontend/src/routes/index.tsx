import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: AboutComponent,
});

function AboutComponent() {
  return <div className="text-white">Hello</div>;
}
