// InputRange.tsx
import React, { useState, useEffect } from 'react'
import './InputRange.css'

interface Props {
  formSearch: number
  setFormSearch: React.Dispatch<React.SetStateAction<number>>
  description: string
  children: React.ReactNode
}

export const InputRange: React.FC<Props> = ({
  formSearch,
  setFormSearch,
  description,
  children
}) => {
  const [inputValue, setInputValue] = useState<string>(formSearch.toString())

  useEffect(() => {
    setInputValue(formSearch.toString())
  }, [formSearch])

  const handleDecrement = () => {
    if (formSearch > 1) {
      setFormSearch(formSearch - 1)
    }
  }

  const handleIncrement = () => {
    if (formSearch < 100) {
      setFormSearch(formSearch + 1)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    // Allow only numbers
    if (/^\d*$/.test(value)) {
      setInputValue(value)
    }
  }

  const handleBlur = () => {
    let value = Number(inputValue)
    if (isNaN(value) || value < 1) {
      value = 1
    } else if (value > 100) {
      value = 100
    }
    setFormSearch(value)
  }

  return (
    <div className='flex flex-col gap-2'>
      <label
        htmlFor='nMoviesHome'
        className='text-primary text-body-14 flex items-center  text-balance gap-5 p-1'
      >
        {children}
        <p className='text-light'>{description}</p>
      </label>
      <div className='input-number-container relative'>
        <input
          type='number'
          id='nMoviesSearch'
          min={1}
          max={100}
          value={inputValue}
          onChange={handleChange}
          onBlur={handleBlur}
          className='bg-dark text-light font-bold border-secondary outline-none p-2 rounded border w-full'
        />
        <div className='absolute inset-y-0 right-0 flex items-center px-2 space-x-1'>
          <button
            type='button'
            onClick={handleDecrement}
            className='text-gray hover:text-light focus:outline-none'
            aria-label='Decrement'
          >
            &minus;
          </button>
          <button
            type='button'
            onClick={handleIncrement}
            className='text-gray hover:text-light focus:outline-none'
            aria-label='Increment'
          >
            &#43;
          </button>
        </div>
      </div>
    </div>
  )
}
