/* eslint-disable @typescript-eslint/no-explicit-any */
import { Link, useLocation } from 'react-router-dom'
import { Icon } from '../assets/Icon'
import { useStore } from '../services/store'

export const Header: React.FC = () => {
  const location = useLocation()
  const nMoviesHome = useStore((state: any) => state.nMoviesHome)
  const setNMoviesHome = useStore((state: any) => state.setNMoviesHome)

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    if (location.pathname === '/') {
      e.preventDefault()
      setNMoviesHome(nMoviesHome + 1)
    }
  }

  return (
    <header className='h-[20vh] flex items-center place-self-start'>
      <Link
        to='/'
        onClick={handleClick}
        className='p-0 m-0 flex items-center h-[15px] justify-center group'
      >
        <Icon className='w-[50px] text-secondary stroke-1' />
        <span className='text-[54px] group-hover:text-primary group-hover:underline transition-all duration-700 ease-in-out'>
          Movies Recomender
        </span>
      </Link>
    </header>
  )
}
