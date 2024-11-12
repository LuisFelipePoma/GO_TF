export const averageStyles = (average: number) => {
  if (average < 5) {
    return 'bg-red-400'
  } else if (average < 7) {
    return 'bg-yellow-400'
  } else {
    return 'bg-green-400'
  }
}
