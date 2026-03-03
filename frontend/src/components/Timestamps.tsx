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
		<div className='flex flex-col space-y-4 lg:space-y-8'>
			<div className='mx-auto justify-center'>
				<div className='space-y-4'>
					<Timer
						paused={paused}
						startUnix={startTime}
						onUpdate={(seconds: number) => {
							setRunningTimer(seconds)
						}}
					/>
					<div className='justify-center flex space-x-2'>
						<button
							disabled={!paused}
							className='btn btn-soft btn-success icon-outlined'
							onClick={() => startMut.mutate()}
						>
							play_arrow
						</button>
						<button
							disabled={paused}
							className='btn btn-error btn-soft icon-outlined'
							onClick={() => {
								if (!currTimer) return
								stopMut.mutate(currTimer.id)
							}}
						>
							stop
						</button>
					</div>

					<p className='text-center text-xl'>
						Today:{' '}
						<span>
							{totalTime.hours}h {totalTime.minutes}m {totalTime.seconds}s
						</span>
					</p>
				</div>
			</div>

			<div className='overflow-x-auto rounded-box'>
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
		<div className='text-2xl text-center'>
			{timer.hours}h {timer.minutes}m {timer.seconds}s
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
					<h3 className='text-base font-semibold mb-2 px-4'>
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
			<table className='table bg-base-300'>
				<thead>
					<tr>
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
								className='hover:bg-base-300 bg-base-100'
							>
								<td>
									{new Date(t.start_time)
										.toLocaleString('de-DE', {
											dateStyle: 'medium',
											timeStyle: 'medium',
										})
										.replaceAll('/', '.')}
								</td>
								<td>
									{t.end_time &&
										new Date(t.end_time)
											.toLocaleString('de-DE', {
												dateStyle: 'medium',
												timeStyle: 'medium',
											})
											.replaceAll('/', '.')}
								</td>
								<td>
									{duration.hours}h {duration.minutes}m {duration.seconds}s
								</td>
							</tr>
						)
					})}
				</tbody>
				{footerCounter !== null && (
					<tfoot>
						<tr className='bg-base-200 font-semibold'>
							<td></td>
							<td></td>
							<td className='text-primary underline underline-offset-6'>
								{footerCounter.hours}h {footerCounter.minutes}m{' '}
								{footerCounter.seconds}s
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
						<div key={user.id} className='my-2'>
							<details className='collapse border-base-300 border collapse-arrow'>
								<summary className='collapse-title bg-base-300 focus:bg-primary/20 focus:text-white hover:bg-base-100 font-semibold '>
									{user.username}
								</summary>
								<div className='collapse-content px-0 flex flex-col text-sm'>
									<table className='bg-base-100 rounded-none table mb-12'>
										<thead>
											<tr>
												<th>Worked</th>
												<th>Expected</th>
												<th>{overtimeLabel}</th>
											</tr>
										</thead>
										<tbody>
											<tr className='*:text-xs *:lg:text-base'>
												<td>
													{counter.hours}
													<span className='font-normal px-1'>h</span>{' '}
													{'  '}
													{counter.minutes}
													<span className='font-normal px-1'>m</span>{' '}
													{'  '}
												</td>
												<td>
													{expectedCounter.hours}
													<span className='font-normal px-1'>h</span>{' '}
													{'  '}
													{expectedCounter.minutes}
													<span className='font-normal px-1'>m</span>{' '}
													{'  '}
												</td>
												<td>
													{overtimeCounter.hours}
													<span className='font-normal px-1'>h</span>{' '}
													{'  '}
													{overtimeCounter.minutes}
													<span className='font-normal px-1'>m</span>{' '}
													{'  '}
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
