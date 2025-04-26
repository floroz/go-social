import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
// Removed duplicate imports
import { BrowserRouter, Routes, Route } from 'react-router'; // Attempting import from 'react-router'
import './index.css';
import App from './App.tsx';
import LoginPage from './pages/LoginPage.tsx'; // Import LoginPage
import SignupPage from './pages/SignupPage.tsx'; // Import SignupPage

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} /> 
        <Route path="/login" element={<LoginPage />} /> 
        <Route path="/signup" element={<SignupPage />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
