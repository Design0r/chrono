type StatCardElementProps = {
  title: string;
  skeleton?: boolean;
  stat?: string;
  statClassName?: string;
  subtitle: string;
  children?: React.ReactNode[] | React.ReactNode | undefined;
};

type StatCardProps = {
  children: React.ReactNode[] | React.ReactNode | undefined;
};

export function StatCard({ children }: StatCardProps) {
  return (
    <div className="stats bg-base-100 max-lg:stats-vertical grid grid-cols-2 grid-rows-2 lg:grid-rows-1 lg:grid-cols-4 w-full">
      {children &&
        (Array.isArray(children) ? <>{...children}</> : <>{children}</>)}
    </div>
  );
}

export function StatCardElement({
  title,
  stat,
  statClassName = "-mb-1 pt-1.5 stat-value max-sm:text-2xl text-primary",
  subtitle,
  skeleton = false,
  children,
}: StatCardElementProps) {
  return (
    <div className={`stat ${skeleton && "skeleton"}`}>
      <div className="stat-figure"></div>
      <div className="stat-title truncate text-accent/75 text-base ">
        {title}
      </div>
      {children ? <>{children}</> : <div className={statClassName}>{stat}</div>}
      <div className="stat-desc text-accent/30">{subtitle}</div>
    </div>
  );
}
