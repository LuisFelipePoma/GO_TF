export const Loading = () => {
  return (
    <div className='fixed inset-0 flex flex-col justify-center items-center bg-black bg-opacity-50'>
      <div className='rounded-full loader border-8 border-t-8 border-gray-200 h-32 w-32 animate-spin'></div>
      <p className='mt-4 text-white text-lg'>Cargando películas...</p>
    </div>
  )
}
