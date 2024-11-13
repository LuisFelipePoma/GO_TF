import { create } from 'zustand'

export const useStore = create(set => ({
  backgroundPath: null,
  setBackgroundPath: (backgroundPath: string) => set({ backgroundPath })
}))
