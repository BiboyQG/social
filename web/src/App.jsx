import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import ConfirmPage from './route/ComfirmPage'
import SuccessPage from './route/SuccessPage'
import ErrorPage from './route/ErrorPage'
import './App.css'

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/confirm/:token" element={<ConfirmPage />} />
        <Route path="/success" element={<SuccessPage />} />
        <Route path="/error" element={<ErrorPage />} />
      </Routes>
    </Router>
  )
}

export default App
