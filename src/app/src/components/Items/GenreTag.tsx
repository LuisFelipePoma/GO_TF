interface Props {
  genre: {
    id?: number
    name?: string
  }
  onClick: (genre: string) => void
  className: string
}

export const GenreTag: React.FC<Props> = ({ genre, onClick, className }) => {
  return (
    <span
      className={`
	 px-2 py-1 bg-secondary text-body-12 uppercase rounded-md hover:bg-tertiary transition-colors duration-300 ease-in-out cursor-pointer 
	 ${className}
	  `}
      onClick={() => onClick(genre.name!)}
    >
      {genre.name}
    </span>
  )
}
