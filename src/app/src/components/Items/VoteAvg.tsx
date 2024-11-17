import { StarRating } from './StartRating'

interface Props {
  vote_average?: number | null
}

export const VoteAvg: React.FC<Props> = ({ vote_average }) => {
  return (
    <span
      className={`select-none flex gap-1 items-center px-4 py-2 rounded-full font-semibold 
					bg-secondary text-light
					shadow-md transition-colors duration-300 ease-in-out`}
    >
      <StarRating voteAverage={vote_average!} />
      {vote_average?.toFixed(1)}/10
    </span>
  )
}
