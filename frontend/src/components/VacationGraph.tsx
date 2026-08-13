import dayjs from 'dayjs'
import isoWeek from 'dayjs/plugin/isoWeek'
import {useCallback, useEffect, useMemo, useRef, useState} from 'react'
import type {VacationGraphMonth} from '../types/response'
import {clamp} from '../utils/math'

dayjs.extend(isoWeek)

/** Maximale Tooltip-Breite (max-w-56), Basis für das Clamping am Viewport-Rand. */
const TOOLTIP_WIDTH = 224
/** Abstand zwischen Zelle und Tooltip. */
const TOOLTIP_GAP = 8
/** Platz über der Zelle, unter dem der Tooltip nach unten klappt. */
const TOOLTIP_FLIP_SPACE = 180

type TooltipState = {
	day: VacationGraphMonth
	x: number
	y: number
	below: boolean
}

type ShowTooltip = (day: VacationGraphMonth, rect: DOMRect) => void

export function OverviewDay({
	day,
	onShow,
	onHide,
}: {
	day: VacationGraphMonth
	onShow?: ShowTooltip
	onHide?: () => void
}) {
	const cellRef = useRef<HTMLDivElement>(null)

	// Grün → Gelb → Rot (OKLCH, 8 Stufen)
	const greens = [
		'oklch(0.273 0.028 265)', // Grün
		'oklch(0.40 0.035 230)',
		'oklch(0.50 0.07 185)', // Gelbgrün
		'oklch(0.70 0.09 150)', // Gelb
		'oklch(0.85 0.15 122)', // Orange-Gelb
		'oklch(0.87 0.19 100)', // Orange-Gelb
		'oklch(0.90 0.29 55)', // Orange-Gelb
		'oklch(0.75 0.23 40)', // Orange-Gelb
		'oklch(0.60 0.22 0)', // Orange-Gelb
	]
	const holidayColor = 'oklch(0.74 0.22 260)'
	const color = day.is_holiday ? holidayColor : greens[clamp(day.count, 0, greens.length - 1)]

	// Der Tooltip wird nicht hier gerendert: Graph-Container und TitleSection
	// kappen beide den Overflow, ein absolut positionierter Tooltip würde also
	// abgeschnitten. Die Zelle meldet stattdessen nur ihre Position, gezeichnet
	// wird der Tooltip vom Graphen selbst – fixed und damit außerhalb des Clips.
	return (
		<div
			className='relative'
			onPointerEnter={() => {
				const rect = cellRef.current?.getBoundingClientRect()
				if (rect) onShow?.(day, rect)
			}}
			onPointerLeave={() => onHide?.()}
		>
			<div
				ref={cellRef}
				className={`w-4 h-1/3 rounded-sm ${day.is_current_week ? '' : ''} ${day.is_current_day ? 'ring-2 ring-primary animate-pulse ring-offset-0.5 ring-offset-base-300' : ''} ${day.first_day_of_month ? 'after:content-[""] after:absolute after:-top-2.5 after:left-1.25 after:bg-info/30 after:w-1.5 after:h-0.75 after:rounded-xl' : ''}`}
				style={
					day.first_day_of_month
						? {
								backgroundColor: color,
							}
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
	const [tooltip, setTooltip] = useState<TooltipState | null>(null)

	const showTooltip = useCallback<ShowTooltip>((day, rect) => {
		const below = rect.top < TOOLTIP_FLIP_SPACE
		setTooltip({
			day,
			// Am Viewport-Rand einklemmen, damit der Tooltip nie aus dem Bild läuft.
			x: clamp(
				rect.left + rect.width / 2,
				TOOLTIP_WIDTH / 2 + TOOLTIP_GAP,
				window.innerWidth - TOOLTIP_WIDTH / 2 - TOOLTIP_GAP,
			),
			y: below ? rect.bottom + TOOLTIP_GAP : rect.top - TOOLTIP_GAP,
			below,
		})
	}, [])

	const hideTooltip = useCallback(() => setTooltip(null), [])

	// Fixed positioniert: beim Scrollen (Seite oder Graph) würde der Tooltip von
	// seiner Zelle weglaufen, deshalb blenden wir ihn dann aus.
	const tooltipOpen = tooltip !== null
	useEffect(() => {
		if (!tooltipOpen) return
		const hide = () => setTooltip(null)
		window.addEventListener('scroll', hide, true)
		window.addEventListener('resize', hide)
		return () => {
			window.removeEventListener('scroll', hide, true)
			window.removeEventListener('resize', hide)
		}
	}, [tooltipOpen])

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
		<>
			{/* Auf Touch-Geräten würde ein Drag über die Zellen Wochen- und
			    Monatszahlen markieren statt zu scrollen; select-none unterbindet
			    das, touch-callout das iOS-Menü beim langen Drücken. Mit Maus
			    bleibt der Text markierbar. */}
			<div className='grid grid-cols-[2.5rem_repeat(11,minmax(0,1fr))] p-5 bg-base-300/50 rounded-2xl xl:overflow-x-hidden overflow-x-auto mb-12 pointer-coarse:select-none pointer-coarse:[-webkit-touch-callout:none]'>
				<div className='col-span-1' />
				<div className='col-span-11 grid grid-rows-1 grid-flow-col h-7 gap-1 text-accent/80'>
					{gaps.map((g, i) => (
						<div key={`top-${i}`} className='contents'>
							<p
								className={`h-3.5 w-3.5 text-center text-sm ${currMonth === i + 1 ? 'text-primary font-light' : 'font-light text-accent/80'}`}
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
							className={`w-3.5 h-3.5 text-center ${currWeek === week ? 'text-primary font-light' : 'font-light text-accent/40'}`}
						>
							{week}
						</p>
					))}
				</div>

				<div className='col-span-1 grid grid-rows-7 text-accent/45 text-sm'>
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
							<OverviewDay key={i} day={d} onShow={showTooltip} onHide={hideTooltip} />
						))}
					</div>
				</div>
			</div>

			{tooltip && (
				<div
					role='tooltip'
					className='pointer-events-none fixed z-50 max-w-56 rounded-field bg-neutral px-2 py-1 text-center text-sm leading-tight text-neutral-content shadow-lg'
					style={{
						left: tooltip.x,
						top: tooltip.y,
						transform: `translate(-50%, ${tooltip.below ? '0' : '-100%'})`,
					}}
				>
					<div>{tooltip.day.date}</div>
					<div>Count: {tooltip.day.count}</div>
					{tooltip.day.usernames?.map((u, i) => (
						<p key={i}>{u}</p>
					))}
				</div>
			)}
		</>
	)
}
