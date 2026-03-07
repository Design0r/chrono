import dayjs from 'dayjs'
import isoWeek from 'dayjs/plugin/isoWeek'
import {useMemo} from 'react'
import type {VacationGraphMonth} from '../types/response'
import {clamp} from '../utils/math'

dayjs.extend(isoWeek)

export function OverviewDay({day}: {day: VacationGraphMonth}) {
	const greens = [
		'#313745',
		'#a7f3d0',
		'#6ee7b7',
		'#34d399',
		'#10b981',
		'#059669',
		'#047857',
		'#065f46',
		'#064e3b',
	]
	const holidayColor = '#7C85FF'
	const color = day.is_holiday ? holidayColor : greens[clamp(day.count, 0, 8)]

	return (
		<div className='tooltip'>
			<div className='tooltip-content'>
				<div>{day.date}</div>
				<div>Count: {day.count}</div>
				<>
					{day.usernames?.map((u, i) => (
						<p key={i}>{u}</p>
					))}
				</>
			</div>
			<div
				className={`w-3 h-3 mx-0.25 rounded-full ${day.is_current_week ? 'border border-primary' : ''} ${day.is_current_day ? 'ring-2 ring-warning ring-offset-1 ring-offset-base-300' : ''}`}
				style={
					day.last_day_of_month
						? {boxShadow: '0 1.5rem 0 -0.125rem ', backgroundColor: color}
						: {backgroundColor: color}
				}
			></div>
		</div>
	)
}

export function VacationGraph({
	gaps,
	yearOffset,
	data,
}: {
	gaps: number[]
	yearOffset: number
	data: VacationGraphMonth[]
}) {
	const currWeek = dayjs().isoWeek()
	const currMonth = dayjs().month() + 1 // 1–12

	const cells = useMemo(() => {
		let week = 1
		const out: Array<number> = []

		for (const g of gaps) {
			out.push(week)
			for (let i = 0; i < g; i++) {
				week++
				out.push(week)
			}
			week++
		}

		return out
	}, [gaps])

	return (
		<div className='grid grid-cols-12 p-5 bg-base-300/50 rounded-2xl xl:overflow-x-hidden overflow-x-auto mb-12'>
			<div className='col-span-1' />
			<div className='col-span-11 grid grid-rows-1 grid-flow-col h-7 gap-1 text-accent/80'>
				{gaps.map((g, i) => (
					<div key={`top-${i}`} className='contents'>
						<p
							className={`h-3.5 w-3.5 text-center text-sm ${currMonth === i + 1 ? 'text-warning font-light' : 'font-light text-accent/80'}`}
						>
							{i + 1}
						</p>
						{Array.from({length: g}).map((_, j) => (
							<p key={`top-gap-${i}-${j}`} className='w-3.5 h-3.5 -z-10' />
						))}
					</div>
				))}
			</div>

			<div className='col-span-1' />
			<div className='col-span-11 grid grid-rows-1 grid-flow-col h-11 gap-1 text-xs text-base-content/40'>
				{cells.map((week) => (
					<p
						key={week}
						className={`w-3.5 h-3.5 text-center ${currWeek === week ? 'text-warning font-light' : 'font-light text-accent/40'}`}
					>
						{week}
					</p>
				))}
			</div>

			<div className='col-span-1 grid grid-rows-7 text-accent/45 pl-2.5 text-sm'>
				<p className='truncate'>Mon</p>
				<p className='truncate'>Tue</p>
				<p className='truncate'>Wed</p>
				<p className='truncate'>Thu</p>
				<p className='truncate'>Fri</p>
				<p className='opacity-50 truncate'>Sat</p>
				<p className='opacity-50 truncate'>Sun</p>
			</div>

			<div className='col-span-11 grid grid-rows-7 h-80 grid-flow-col gap-1'>
				<div className='contents'>
					{Array.from({length: yearOffset}).map((_, i) => (
						<p key={`yo-${i}`} />
					))}
					{data.map((d, i) => (
						<OverviewDay key={i} day={d} />
					))}
				</div>
			</div>
		</div>
	)
}
