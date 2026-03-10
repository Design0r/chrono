import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {useEffect, useMemo, useState} from 'react'
import {ChronoClient} from '../api/chrono/client'
import type {User} from '../types/auth'
import type {Timestamp} from '../types/response'
import {useToast} from './Toast'

export function Timestamps({user}: {user: User}) {
	const chrono = useMemo(() => new ChronoClient(), [])

	const [timestamps, setTimestemps] = useState<Timestamp[]>([])
	const [paused, setPaused] = useState<boolean>(true)
	const {addToast, addErrorToast} = useToast()
	const [currTimer, setCurrTimer] = useState<Timestamp | null>(null)
	const [startTime, setStartTime] = useState<number>(Date.now())
	const [runningTimer, setRunningTimer] = useState<number>(0)

	const queryClient = useQueryClient()

	const latestTimestampQ = useQuery({
		queryKey: ['timestamps', 'latest'],
		queryFn: () => chrono.timestamps.getLatest(),
		staleTime: 1000 * 60 * 10, // 10min
		gcTime: 1000 * 60 * 20, // 20min
		retry: false,
	})

	const timestampsQ = useQuery({
		queryKey: ['timestamps'],
		queryFn: () => chrono.timestamps.getForToday(),
		staleTime: 1000 * 60 * 10, // 10min
		gcTime: 1000 * 60 * 20, // 20min
		retry: false,
	})

	const startMut = useMutation({
		mutationKey: ['timestamps', 'start'],
		mutationFn: () => chrono.timestamps.start(),
		onError: (e) => addErrorToast(e),
		onSuccess: (data) => {
			setCurrTimer(data)
			setPaused(false)
			setStartTime(Date.now())
			addToast('Started Timer', 'success')
		},
		retry: false,
	})

	const stopMut = useMutation({
		mutationKey: ['timestamps', 'stop'],
		mutationFn: (id: number) => chrono.timestamps.stop(id),
		onError: (e) => addErrorToast(e),
		onSuccess: () => {
			setCurrTimer(null)
			setStartTime(Date.now())
			setPaused(true)
			queryClient.invalidateQueries({queryKey: ['timestamps']})
			addToast('Stopped Timer', 'success')
		},
		retry: false,
	})

	useEffect(() => {
		if (timestampsQ.isError) return
		setTimestemps(timestampsQ.data || [])
	}, [timestampsQ.data, timestampsQ.isError])

	useEffect(() => {
		if (latestTimestampQ.isError) return
		const latest = latestTimestampQ.data
		if (!latest) return
		const hasEnded = latest.end_time !== null
		if (!hasEnded && latest.id !== currTimer?.id) {
			addToast('Resuming latest unfinished Timer', 'info')
			setCurrTimer(latest)
			setStartTime(Date.parse(latest.start_time))
			setPaused(false)
		} else setPaused(true)
	}, [latestTimestampQ.data, latestTimestampQ.isError])

	useEffect(() => {
		if (timestampsQ.isError) addErrorToast(timestampsQ.error)
	}, [timestampsQ.isError])

	const totalTime = secondsToCounter(durationFromTimestamps(timestamps) + runningTimer)

	return (
		<div className='mx-auto flex lg:flex-row flex-col w-full gap-4 justify-center '>
			<div className='flex items-center justify-center flex-1 rounded-2xl bg-base-200/60 border border-base-300/80 p-6 lg:p-8'>
				<div className='flex flex-col items-center justify-center gap-4'>
					<Timer
						paused={paused}
						startUnix={startTime}
						onUpdate={(seconds: number) => {
							setRunningTimer(seconds)
						}}
					/>
					<div className='flex mt-1 gap-3 justify-center items-center'>
						<button
							disabled={!paused}
							className='btn btn-lg btn-success w-22 rounded-full shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-filled'
							onClick={() => startMut.mutate()}
							title='Timer starten'
						>
							<span className='text-xl icon-filled scale-145'>play_arrow</span>
						</button>
						<button
							disabled={paused}
							className='btn btn-circle btn-lg btn-error shadow-md hover:shadow-lg transition-shadow disabled:opacity-40 disabled:shadow-none icon-outlined'
							onClick={() => {
								if (!currTimer) return
								stopMut.mutate(currTimer.id)
							}}
							title='Timer stoppen'
						>
							<span className='text-xl icon-filled scale-125'>stop</span>
						</button>
					</div>{' '}
					<div className='flex whitespace-nowrap items-center gap-2 mt-0 text-base-content/70'>
						<span
							className={`font-mono font-semibold ${paused ? 'text-success' : 'text-accent/80'}`}
						>
							{(() => {
								const f = formatCounter(totalTime)
								return (
									<>
										<span>{f.hours}</span>
										<span
											className={
												paused
													? 'text-success/90 animate-pulse'
													: 'text-accent/80 animate-pulse'
											}
										>
											:
										</span>
										<span>{f.minutes} h</span>
									</>
								)
							})()}
							<span
								className={`pl-2.5 ${paused ? 'text-success/50' : 'text-accent/60'}`}
							>
								Today
							</span>
						</span>
					</div>
				</div>
			</div>
			<div className='rounded-xl overflow-hidden w-full h-full'>
				<TimestampTable timestamps={timestamps} user={user} />
			</div>
		</div>
	)
}

export function durationFromTimestamps(timestamps: Timestamp[]): number {
	return timestamps
		.map((t) => {
			const start = new Date(t.start_time)
			const end = t.end_time && new Date(t.end_time)
			return end ? (end.getTime() - start.getTime()) / 1000 : 0
		})
		.reduce((acc, curr) => {
			return acc + curr
		}, 0)
}

/** ISO 8601 week number (1–53) for the given date. Week = Mon–Sun, week 1 = week with first Thursday. */
function getWeekNumber(date: Date): number {
	const d = new Date(date)
	d.setHours(0, 0, 0, 0)
	const day = d.getDay() || 7 // 1 = Mon, 7 = Sun
	d.setDate(d.getDate() + 4 - day) // Thursday of this week
	const jan1 = new Date(d.getFullYear(), 0, 1)
	return 1 + Math.floor((d.getTime() - jan1.getTime()) / 86400000 / 7)
}

/** Returns the Monday 00:00:00 (local) of the week containing the given date. Week = Mon–Sun. */
function getWeekMonday(date: Date): Date {
	const d = new Date(date)
	d.setHours(0, 0, 0, 0)
	const day = d.getDay() // 0 = Sun, 1 = Mon, ...
	const diff = day === 0 ? -6 : 1 - day
	d.setDate(d.getDate() + diff)
	return d
}

export type WeekGroup = {
	weekMonday: Date
	weekSunday: Date
	timestamps: Timestamp[]
	totalSeconds: number
}

/** Groups timestamps by calendar week (Mon–Sun), newest week first. */
export function groupTimestampsByWeek(timestamps: Timestamp[]): WeekGroup[] {
	const byWeek = new Map<number, Timestamp[]>()
	for (const t of timestamps) {
		const start = new Date(t.start_time)
		const weekMonday = getWeekMonday(start)
		const key = weekMonday.getTime()
		if (!byWeek.has(key)) byWeek.set(key, [])
		byWeek.get(key)!.push(t)
	}
	const weekSunday = (m: Date) => {
		const s = new Date(m)
		s.setDate(s.getDate() + 6)
		return s
	}
	const groups: WeekGroup[] = []
	byWeek.forEach((ts, key) => {
		const weekMonday = new Date(key)
		groups.push({
			weekMonday,
			weekSunday: weekSunday(weekMonday),
			timestamps: ts.sort(
				(a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
			),
			totalSeconds: durationFromTimestamps(ts),
		})
	})
	groups.sort((a, b) => b.weekMonday.getTime() - a.weekMonday.getTime())
	return groups
}

type TimeCounter = {
	hours: number
	minutes: number
	seconds: number
}

export function secondsToCounter(totalSeconds: number): TimeCounter {
	const seconds = Math.max(0, Math.floor(totalSeconds))
	const hours = Math.floor(seconds / 60 / 60)
	const minutes = Math.floor(seconds / 60) % 60
	const s = seconds % 60
	return {hours, minutes, seconds: s}
}

export function formatCounter(c: TimeCounter): {hours: string; minutes: string; seconds: string} {
	return {
		hours: String(c.hours).padStart(2, '0'),
		minutes: String(c.minutes).padStart(2, '0'),
		seconds: String(c.seconds).padStart(2, '0'),
	}
}

function Timer({
	startUnix,
	paused,
	onUpdate,
}: {
	startUnix: number
	paused: boolean
	onUpdate: (seconds: number) => void
}) {
	const [timer, setTimer] = useState<TimeCounter>(() => secondsToCounter(0))

	useEffect(() => {
		function tick() {
			const elapsedSeconds = (Date.now() - startUnix) / 1000
			onUpdate(elapsedSeconds)
			setTimer(secondsToCounter(elapsedSeconds))
		}

		tick()

		if (paused) return // no interval while paused

		const interval = setInterval(tick, 1000)
		return () => clearInterval(interval)
	}, [startUnix, paused])

	return (
		<div
			className={`font-mono text-3xl tabular-nums tracking-tight text-center transition-opacity duration-300 ${
				paused
					? 'text-base-content/60'
					: 'text-error *:even:text-error/70 *:even:animate-pulse'
			}`}
		>
			<span>{String(timer.hours).padStart(2, '0')}</span>
			<span className=''>:</span>
			<span>{String(timer.minutes).padStart(2, '0')}</span>
			<span className=''>:</span>
			<span>{String(timer.seconds).padStart(2, '0')}</span>
		</div>
	)
}

function formatWeekRange(monday: Date, sunday: Date): string {
	const fmt = (d: Date) =>
		d.toLocaleDateString('de-DE', {day: '2-digit', month: '2-digit', year: 'numeric'})
	return `${fmt(monday)} – ${fmt(sunday)}`
}

export function TimestampTableByWeek({timestamps, user}: {timestamps: Timestamp[]; user: User}) {
	const groups = groupTimestampsByWeek(timestamps)
	return (
		<div className='space-y-6'>
			{groups.map((g) => (
				<section key={g.weekMonday.getTime()}>
					<h3 className='text-base font-semibold mb-3.5 px-4'>
						{getWeekNumber(g.weekMonday)}. KW{' '}
						<span className='text-accent/50 font-normal pl-1 pr-1.5'>|</span>
						{formatWeekRange(g.weekMonday, g.weekSunday)}
					</h3>
					<TimestampTable
						timestamps={g.timestamps}
						user={user}
						footerTotalSeconds={g.totalSeconds}
					/>
				</section>
			))}
		</div>
	)
}

export function TimestampTable({
	timestamps,
	user,
	footerTotalSeconds,
}: {
	timestamps: Timestamp[]
	user: User
	footerTotalSeconds?: number
}) {
	const [modal, setModal] = useState<Timestamp | null>(null)

	const {addErrorToast} = useToast()

	const footerCounter =
		footerTotalSeconds !== undefined ? secondsToCounter(footerTotalSeconds) : null

	return (
		<>
			<table className='table bg-base-300/50 lg:rounded-none'>
				<thead>
					<tr className='text-accent/80 *:w-1/3 *:font-normal'>
						<th>Start</th>
						<th>End</th>
						<th>Duration</th>
					</tr>
				</thead>
				<tbody>
					{timestamps.map((t) => {
						const start = new Date(t.start_time)
						const end = t.end_time && new Date(t.end_time)
						const duration = end
							? secondsToCounter((end.getTime() - start.getTime()) / 1000)
							: {hours: 0, minutes: 0, seconds: 0}

						return (
							<tr
								onClick={() => {
									if (!user.is_superuser) {
										addErrorToast({
											name: 'Permission Error',
											message: 'Only for Admins',
										})

										return
									}
									setModal(t)
								}}
								key={t.id}
								className='hover:bg-base-300 *:font-extralight *:text-info/70 bg-base-200/40'
							>
								<td>
									<div className='flex gap-1.5 items-center'>
										<span className='text-info/50 w-6'>
											{start
												.toLocaleDateString('de-DE', {weekday: 'short'})
												.slice(0, 2)}
											.
										</span>
										<span className='text-info/50 w-6 md:w-8.5'>
											{start
												.toLocaleDateString('de-DE', {
													day: '2-digit',
													month: '2-digit',
												})
												.replaceAll('/', '.')
												.slice(0, -1)}
										</span>
										<span className='block md:border-l border-primary/20 pl-1.75 text-info/90'>
											{start.toLocaleTimeString('de-DE', {
												hour: '2-digit',
												minute: '2-digit',
											})}
										</span>
									</div>
								</td>
								<td>
									{end ? (
										<div className='flex gap-1.5 items-center'>
											<span className='text-info/50 w-6'>
												{end
													.toLocaleDateString('de-DE', {weekday: 'short'})
													.slice(0, 2)}
												.
											</span>
											<span className='text-info/50 w-6 md:w-8.5'>
												{' '}
												{end
													.toLocaleDateString('de-DE', {
														day: '2-digit',
														month: '2-digit',
													})
													.replaceAll('/', '.')
													.slice(0, -1)}
											</span>
											<span className='block md:border-l border-primary/20 pl-1.75 text-info/90'>
												{end.toLocaleTimeString('de-DE', {
													hour: '2-digit',
													minute: '2-digit',
												})}
											</span>
										</div>
									) : (
										'–'
									)}
								</td>
								<td className='*:text-info/90'>
									{(() => {
										const f = formatCounter(duration)
										return (
											<>
												<span>{f.hours}</span>
												<span>:</span>
												<span>{f.minutes}</span> h{/* <span>:</span> */}
												{/* <span>{f.seconds}</span> */}
											</>
										)
									})()}
								</td>
							</tr>
						)
					})}
				</tbody>
				{footerCounter !== null && (
					<tfoot>
						<tr className='bg-base-200/70 border-b-lg font-semibold'>
							<td></td>
							<td></td>
							<td className='text-primary'>
								{(() => {
									const f = formatCounter(footerCounter)
									return (
										<>
											<span>{f.hours}</span>
											<span>:</span>
											<span>{f.minutes}</span>{' '}
											<span className='text-primary/80'>h</span>
											{/* <span>{f.seconds}</span> */}
										</>
									)
								})()}
							</td>
						</tr>
					</tfoot>
				)}
			</table>
			{modal && user.is_superuser && (
				<EditModal timestamp={modal} onClose={() => setModal(null)} />
			)}
		</>
	)
}

// ISO ("2025-12-23T20:44:00Z") -> datetime-local ("2025-12-23T20:44")
export function isoToDatetimeLocal(iso: string) {
	const d = new Date(iso)
	const pad = (n: number) => String(n).padStart(2, '0')
	if (d.getSeconds() === 0) {
		d.setSeconds(new Date().getSeconds())
	}
	// datetime-local is *local time* by spec
	const yyyy = d.getFullYear()
	const mm = pad(d.getMonth() + 1)
	const dd = pad(d.getDate())
	const hh = pad(d.getHours())
	const min = pad(d.getMinutes())

	return `${yyyy}-${mm}-${dd}T${hh}:${min}`
}

export function isoToDateLocal(iso: string) {
	const d = new Date(iso)
	const pad = (n: number) => String(n).padStart(2, '0')

	const yyyy = d.getFullYear()
	const mm = pad(d.getMonth() + 1)
	const dd = pad(d.getDate())

	return `${yyyy}-${mm}-${dd}`
}

// datetime-local ("2025-12-23T21:44") -> ISO UTC ("2025-12-23T20:44:00Z")
export function datetimeLocalToIso(value: string) {
	const d = new Date(value)
	const iso = d.toISOString()
	const fixed = `${iso.split('.')[0]}Z`
	return fixed
}

export function EditModal({timestamp, onClose}: {timestamp: Timestamp; onClose: () => void}) {
	const queryClient = useQueryClient()
	const [startDate, setStartDate] = useState(isoToDatetimeLocal(timestamp.start_time))
	const [endDate, setEndDate] = useState<string | null>(
		timestamp.end_time ? isoToDatetimeLocal(timestamp.end_time) : null,
	)

	useEffect(() => {
		setStartDate(isoToDatetimeLocal(timestamp.start_time))
		if (timestamp.end_time) setEndDate(isoToDatetimeLocal(timestamp.end_time))
		else setEndDate(null)
	}, [timestamp.start_time, timestamp.end_time])

	const chrono = new ChronoClient()
	const {addToast, addErrorToast} = useToast()

	const mutation = useMutation({
		mutationKey: ['timestamps', timestamp.id],
		mutationFn: ({start, end}: {start: string; end: string | null}) =>
			chrono.timestamps.update({
				id: timestamp.id,
				user_id: timestamp.user_id,
				start_time: datetimeLocalToIso(start),
				end_time: end ? datetimeLocalToIso(end) : null,
			}),
		onError: (e) => addErrorToast(e),
		onSuccess: () => {
			addToast('Updated Timestamp', 'success')
			queryClient.invalidateQueries({queryKey: ['timestamps']})
			onClose()
		},
		retry: false,
	})

	return (
		<div className='fixed inset-0 z-50 flex text-white items-center justify-center p-4'>
			<button
				aria-label='Close modal'
				onClick={onClose}
				className='absolute inset-0 bg-black/50 backdrop-blur-sm'
			/>

			<div
				role='dialog'
				aria-modal='true'
				className='relative w-full max-w-lg rounded-2xl bg-base-100 shadow-2xl ring-1 ring-black/10'
			>
				<div className='flex items-center justify-between px-5 py-4 border-b border-black/10'>
					<h2 className='text-base font-semibold'>Edit Timestamp</h2>

					<button
						type='button'
						onClick={onClose}
						className='inline-flex h-9 w-9 items-center justify-center rounded-full hover:bg-slate-100 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-400'
					>
						<span className='icon-outlined text-[20px] leading-none'>close</span>
					</button>
				</div>

				<div className='px-5 py-4 space-y-4'>
					<div className='grid gap-4 sm:grid-cols-2'>
						<label className='space-y-2'>
							<span className='text-sm font-medium'>Start</span>
							<input
								type='datetime-local'
								className='input'
								value={startDate}
								onChange={(e) => setStartDate(e.target.value)}
							/>
						</label>

						<label className='space-y-2'>
							<span className='text-sm font-medium'>End</span>
							<input
								type='datetime-local'
								className='input'
								value={endDate || ''}
								onChange={(e) => setEndDate(e.target.value)}
							/>
						</label>
					</div>
				</div>

				<div className='flex items-center justify-end gap-2 px-5 py-4 border-t border-black/10'>
					<button type='button' onClick={onClose} className='btn btn-soft btn-error'>
						Cancel
					</button>
					<button
						type='button'
						className='btn btn-soft btn-success'
						onClick={() => mutation.mutate({start: startDate, end: endDate})}
						disabled={mutation.isPending}
					>
						{mutation.isPending ? 'Saving...' : 'Save'}
					</button>
				</div>
			</div>
		</div>
	)
}

export function TeamTimestamps({
	startDate,
	endDate,
	user: currUser,
}: {
	startDate?: string
	endDate?: string
	user: User
}) {
	const chrono = new ChronoClient()
	const currYear = startDate ? new Date(startDate).getFullYear() : new Date().getFullYear()

	const allTimestampsQ = useQuery({
		queryKey: ['timestamps', 'all', startDate, endDate],
		queryFn: () => chrono.timestamps.getAll(startDate, endDate),
		staleTime: 1000 * 60 * 1, // 1min
		gcTime: 1000 * 60 * 30, // 30min
		retry: false,
	})

	const usersQ = useQuery({
		queryKey: ['users'],
		queryFn: () => chrono.users.getUsers(),
		staleTime: 1000 * 60 * 1, // 1min
		gcTime: 1000 * 60 * 30, // 30min
		retry: false,
	})

	const allWorktimesQ = useQuery({
		queryKey: ['worktimes', 'all', startDate, endDate],
		queryFn: () => chrono.timestamps.getWorkHoursForAllUsers(currYear),
		staleTime: 1000 * 60 * 1, // 1min
		gcTime: 1000 * 60 * 30, // 30min
		retry: false,
	})

	const queries = [allTimestampsQ, usersQ, allWorktimesQ]
	const anyPending = queries.some((q) => q.isPending)
	const firstError = queries.find((q) => q.isError)?.error

	const usersMap = useMemo(() => {
		const users = (usersQ.data ?? []) as User[]
		return users.reduce(
			(map, u) => {
				map[u.id] = u
				return map
			},
			{} as Record<number, User>,
		)
	}, [usersQ.data])

	const timestampsMap = useMemo(() => {
		const timestamps = (allTimestampsQ.data ?? []) as Timestamp[]
		return timestamps.reduce(
			(map, ts) => {
				;(map[ts.user_id] ??= []).push(ts)
				return map
			},
			{} as Record<number, Timestamp[]>,
		)
	}, [allTimestampsQ.data])

	if (anyPending || firstError) return <></>

	const worktimes = allWorktimesQ.data!

	return (
		<>
			<hr />
			<h2>Team Timestamps</h2>
			<div className='w-full'>
				{Object.entries(timestampsMap).map(([k, v]) => {
					const user = usersMap[Number(k)]
					const counter = secondsToCounter(durationFromTimestamps(v))

					if (!user) return <div key={0}></div>

					const worktime = worktimes[user.id]
					const expectedCounter = secondsToCounter(worktime.expected * 3600)
					let overtime = (worktime.worked - worktime.expected) * 3600
					let overtimeLabel = 'Overtime'
					if (overtime < 0) {
						overtime *= -1
						overtimeLabel = 'Missing Time'
					}
					const overtimeCounter = secondsToCounter(overtime)

					return (
						<div key={user.id} className='my-4'>
							<details className='collapse border-base-300 border collapse-arrow'>
								<summary className='collapse-title bg-base-300/50 focus-within:bg-info/25 focus:text-white hover:bg-info/25 font-semibold '>
									{user.username}
								</summary>
								<div className='collapse-content px-0 pt-6 pb-0 bg-black/20 flex flex-col text-sm'>
									<h3 className='text-base font-semibold mb-3.5 px-4'>
										Overview
									</h3>{' '}
									<table className='bg-base-300/50 rounded-none table mb-12'>
										<thead>
											<tr className='text-accent/80 *:w-1/3 *:font-normal'>
												<th>Worked</th>
												<th>Expected</th>
												<th>{overtimeLabel}</th>
											</tr>
										</thead>
										<tbody>
											<tr className='hover:bg-base-300 text-info/70 bg-base-200/40'>
												<td className='*:text-info'>
													{(() => {
														const f = formatCounter(counter)
														return (
															<>
																<span>{f.hours}</span>
																<span>:</span>
																<span>{f.minutes}</span> h
																{/* <span>:</span> */}
																{/* <span>{f.seconds}</span> */}
															</>
														)
													})()}
												</td>
												<td className='*:text-info/90'>
													{(() => {
														const f = formatCounter(expectedCounter)
														return (
															<>
																<span>{f.hours}</span>
																<span>:</span>
																<span>{f.minutes}</span> h
																{/* <span>:</span> */}
																{/* <span>{f.seconds}</span> */}
															</>
														)
													})()}
												</td>
												<td className='*:text-info/90'>
													{(() => {
														const f = formatCounter(overtimeCounter)
														return (
															<>
																<span>{f.hours}</span>
																<span>:</span>
																<span>{f.minutes}</span> h
																{/* <span>:</span> */}
																{/* <span>{f.seconds}</span> */}
															</>
														)
													})()}
												</td>
											</tr>
										</tbody>
									</table>
									<TimestampTableByWeek timestamps={v} user={currUser} />
								</div>
							</details>
						</div>
					)
				})}
			</div>
		</>
	)
}
