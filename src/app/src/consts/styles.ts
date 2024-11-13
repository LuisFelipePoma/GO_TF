// Función para obtener las clases de estilo según el puntaje
export const averageStyles = (score: number) => {
  if (score >= 7) {
    return 'bg-green-500 text-white' // Verde más intenso y texto blanco para mejor contraste
  } else if (score >= 4) {
    return 'bg-yellow-100 text-yellow-800' // Mantener amarillo como está
  } else {
    return 'bg-red-500 text-white' // Rojo más intenso y texto blanco para mejor contraste
  }
}
