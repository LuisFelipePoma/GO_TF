import { averageStyles } from '../consts/styles'

interface Props {
  vote_average?: number | null
  className?: string
}

export const VoteAvg: React.FC<Props> = ({ vote_average, className = '' }) => {
  return (
    <span
      className={`select-none inline-flex items-center px-3 py-1 rounded-full font-semibold 
					${averageStyles(vote_average!)}
					shadow-md transition-colors duration-300 ease-in-out
					${className}`}
    >
      <svg
        className='w-4 h-4 mr-1 text-yellow-500'
        fill='currentColor'
        viewBox='0 0 20 20'
        xmlns='http://www.w3.org/2000/svg'
      >
        <path d='M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.286 3.97a1 1 0 00.95.69h4.178c.969 0 1.371 1.24.588 1.81l-3.385 2.455a1 1 0 00-.364 1.118l1.286 3.97c.3.921-.755 1.688-1.54 1.118L10 13.347l-3.385 2.455c-.784.57-1.838-.197-1.54-1.118l1.286-3.97a1 1 0 00-.364-1.118L2.98 9.397c-.783-.57-.38-1.81.588-1.81h4.178a1 1 0 00.95-.69l1.286-3.97z' />
      </svg>
      {vote_average?.toFixed(1)}/10
    </span>
  )
}
