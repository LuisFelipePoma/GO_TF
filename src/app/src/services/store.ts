import { create } from 'zustand'

export const useStore = create(set => ({
  backgroundPath: null,
  setBackgroundPath: (backgroundPath: string) => set({ backgroundPath }),
  nMoviesSearch: 21,
  setNMoviesSearch: (nMoviesSearch: number) => set({ nMoviesSearch }),
  nMoviesUser: 3,
  setNMoviesUser: (nMoviesUser: number) => set({ nMoviesUser }),
  nMoviesHome: 21,
  setNMoviesHome: (nMoviesHome: number) => set({ nMoviesHome }),
  forwardHistory: 0,
  setForwardHistory: (forwardHistory: number) => set({ forwardHistory })
}))
