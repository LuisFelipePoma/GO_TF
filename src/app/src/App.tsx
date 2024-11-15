import React from 'react'
import { Routes, Route, BrowserRouter as Router } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import MovieInfo from './pages/MovieInfo'
import { AnimatePresence, motion } from 'framer-motion'

const App: React.FC = () => {
  return (
    <Router>
      <AnimatePresence mode='wait'>
        <Routes>
          <Route path='/' element={<Layout />}>
            <Route index element={<Home />} />
            <Route
              path='movie/:id'
              element={
                <motion.div
                  initial={{ opacity: 0, x: 100 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -100 }}
                  transition={{ duration: 0.5 }}
                >
                  <MovieInfo />
                </motion.div>
              }
            />
          </Route>
        </Routes>
      </AnimatePresence>
    </Router>
  )
}

export default App
