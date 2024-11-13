import { Link, useLocation } from 'react-router-dom'
import { Icon } from '../assets/Icon'

export const Header: React.FC = () => {
  const location = useLocation()

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    if (location.pathname === '/') {
      e.preventDefault()
      window.location.reload()
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
