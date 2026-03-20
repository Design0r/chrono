import {useMutation, useQueryClient} from '@tanstack/react-query'
import {useEffect, useState} from 'react'
import {ChronoClient} from '../../api/chrono/client'
import type {Timestamp} from '../../types/response'
import {
	datetimeLocalToIso,
	isoToDatetimeLocal,
} from '../../lib/timestamp-utils'
import {useToast} from '../Toast'

export function EditTimestampModal({
	timestamp,
	onClose,
}: {
	timestamp: Timestamp
	onClose: () => void
}) {
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
