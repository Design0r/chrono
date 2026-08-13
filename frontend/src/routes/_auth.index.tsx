import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ErrorPage } from "../components/ErrorPage";
import {
  LoadingSpinner,
  LoadingSpinnerPage,
} from "../components/LoadingSpinner";
import { StatCard, StatCardElement } from "../components/StatCard";
import { Timestamps } from "../components/Timestamps";
import { TitleSection } from "../components/TitleSection";
import { useToast } from "../components/Toast";
import { VacationGraph } from "../components/VacationGraph";
import type { UserWithVacation } from "../types/auth";
import type { WorkTime } from "../types/response";
import { dayOfYear, daysInYear } from "../utils/calendar";

export const Route = createFileRoute("/_auth/")({
  component: Home,
});

function Home() {
  const { chrono, auth } = Route.useRouteContext();
  const year = new Date().getFullYear();
  const { addErrorToast } = useToast();

  const userQ = useQuery({
    queryKey: ["user", auth.userId, "vacation", year],
    queryFn: () =>
      chrono.users.getUserById(auth.userId!, {
        year: year,
      }),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
    retry: false,
  });

  const vacationQ = useQuery({
    queryKey: ["vacationGraph"],
    queryFn: () => chrono.events.getVacationGraph(year),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
    retry: false,
  });

  const aworkQ = useQuery({
    queryKey: ["awork", year],
    queryFn: () => chrono.awork.getWorkTimesforYear(year),
    staleTime: 1000 * 60 * 60, // 1h
    gcTime: 1000 * 60 * 60 * 2, // 2h
    retry: false,
  });

  const worktimeQ = useQuery({
    queryKey: ["timestamps", "worked", year],
    queryFn: () => chrono.timestamps.getWorkHours(year),
    staleTime: 1000 * 60 * 60, // 1h
    gcTime: 1000 * 60 * 60 * 2, // 2h
    retry: false,
  });

  // awork ist optional: ohne verknüpfte awork-ID antwortet der Server mit 422.
  const awork: WorkTime | undefined = aworkQ.isError ? undefined : aworkQ.data;

  const queries = [userQ, vacationQ, worktimeQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const user = userQ.data! as UserWithVacation;
  const vacation = vacationQ.data!;
  const worktimes = worktimeQ.data! as WorkTime;

  const daysYear = daysInYear(year);
  const currDay = dayOfYear();
  const vacRemainingPercent =
    (user.vacation_remaining / user.vacation_days) * 100;
  const yearRemainingPercent = (currDay / daysYear) * 100;
  const vacTakenPercent =
    user.vacation_remaining > user.vacation_days
      ? (user.vacation_used / (user.vacation_used + user.vacation_remaining)) *
        100
      : 100 - vacRemainingPercent;

  const aworkWorkRemaining = awork ? awork.expected - awork.worked : 0;
  const overtimeDiff = worktimes.worked - worktimes.expected;

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle space-y-8">
      <div className="text-[48px] text-primary font-light mb-2">
        <span className="animate-pulse font-medium text-white pr-1"> Hej </span>
        {user?.username}
      </div>
      <TitleSection title="Timestamps">
        <Timestamps user={user} />
      </TitleSection>

      <TitleSection title="Your worktimes">
        <StatCard>
          <StatCardElement
            title={overtimeDiff >= 0 ? "Overtime" : "Time deficit"}
            subtitle={`${overtimeDiff >= 0 ? "" : ""}${overtimeDiff.toFixed(2)} hours`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl text-primary/85">
              {(Math.abs(overtimeDiff) / 8).toFixed(1)}{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Expected work"
            subtitle={`${worktimes.expected} hours`}
          >
            <span className=" pt-1.5 stat-value max-sm:text-2xl">
              {(worktimes.expected / 8).toFixed(1)}{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation taken"
            subtitle={`${worktimes.vacation} hours`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {(worktimes.vacation / 8).toFixed(1)}{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Holidays"
            subtitle={`${worktimes.holidays} hours`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {(worktimes.holidays / 8).toFixed(1)}{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
        </StatCard>
      </TitleSection>

      {aworkQ.isPending ? (
        <div className="skeleton h-30 flex justify-center w-full">
          <LoadingSpinner />
        </div>
      ) : (
        <>
          <TitleSection title="Your Awork worktimes">
            {awork ? (
              <StatCard>
                <StatCardElement
                  title="Work done this year"
                  subtitle={`${aworkWorkRemaining.toFixed(2)} h ${aworkWorkRemaining > 0 ? "remaining" : "over"}`}
                >
                  <span className="pt-1.5 stat-value max-sm:text-2xl">
                    {awork.worked.toFixed(2)}{" "}
                    <span className="text-accent/15">h</span>
                  </span>
                </StatCardElement>
                <StatCardElement
                  title="Expected work"
                  subtitle={`${(awork.expected / 8).toFixed(2)} days`}
                >
                  <span className="pt-1.5 stat-value max-sm:text-2xl">
                    {awork.expected} <span className="text-accent/15">h</span>
                  </span>
                </StatCardElement>
                <StatCardElement
                  title="Vacation taken"
                  subtitle={`${(awork.vacation / 8).toFixed(2)} days`}
                >
                  <span className="pt-1.5 stat-value max-sm:text-2xl">
                    {awork.vacation} <span className="text-accent/15">h</span>
                  </span>
                </StatCardElement>
                <StatCardElement
                  title="Holidays"
                  subtitle={`${(awork.holidays / 8).toFixed(2)} days`}
                >
                  <span className="pt-1.5 stat-value max-sm:text-2xl">
                    {awork.holidays} <span className="text-accent/15">h</span>
                  </span>
                </StatCardElement>
              </StatCard>
            ) : (
              <div className="mx-auto text-center space-y-2">
                <h1>Awork ID not configured</h1>
                <Link className="btn btn-soft btn-primary" to="/profile">
                  Change Now
                </Link>
              </div>
            )}
          </TitleSection>
        </>
      )}

      <TitleSection title="Your vacation">
        <StatCard>
          <StatCardElement
            title="Vacation remaining"
            subtitle={`${vacRemainingPercent.toFixed(2)}% remaining`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {user.vacation_remaining}{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation taken"
            subtitle={`${vacTakenPercent.toFixed(2)}% taken`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {user.vacation_used} <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation total"
            subtitle={`${user.vacation_days} days total`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {user.vacation_days} <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Vacation pending"
            subtitle={`${user.pending_events} event${user.pending_events > 1 ? "s" : ""} pending`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              <span className="animate-pulse text-primary">
                {user.pending_events}
              </span>{" "}
              <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
        </StatCard>
      </TitleSection>

      <TitleSection title="Year progession">
        <StatCard>
          <StatCardElement
            title="Days this year"
            subtitle={`${daysYear} days total`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {daysYear} <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days passed"
            subtitle={`${currDay} days passed`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {currDay} <span className="text-accent/15">d</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days completed"
            subtitle={`${(100 - yearRemainingPercent).toFixed(2)}% remaining`}
          >
            <span className="pt-1.5 stat-value max-sm:text-2xl">
              {yearRemainingPercent.toFixed(2)}{" "}
              <span className="text-accent/15">%</span>
            </span>
          </StatCardElement>
          <StatCardElement
            title="Days progress"
            subtitle={`${(100 - yearRemainingPercent).toFixed(2)}% remaining`}
          >
            <progress
              className="progress progress-primary my-5 h-4"
              value={yearRemainingPercent}
              max="100"
              role="progressbar"
            />
          </StatCardElement>
        </StatCard>
      </TitleSection>
      <TitleSection title="Team Vacation">
        <VacationGraph
          yearOffset={vacation.year_offset}
          gaps={vacation.month_gaps}
          data={vacation.vacation_data}
        />
      </TitleSection>
    </div>
  );
}
