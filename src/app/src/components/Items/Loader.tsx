import './Loader.css'

export const Loader = () => {
  return (
    <article className='fixed inset-0 flex flex-col justify-center items-center bg-black bg-opacity-50 z-50'>
      <div className='loader '>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
        <div className='loader-square'></div>
      </div>
    </article>
  )
}
