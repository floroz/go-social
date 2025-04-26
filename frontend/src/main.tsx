import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router'; // Attempting import from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'; // Import React Query
import './index.css';
import App from './App.tsx';
import LoginPage from './pages/LoginPage.tsx'; // Import LoginPage
import SignupPage from './pages/SignupPage.tsx'; // Import SignupPage

// Create a client
const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}> {/* Wrap with provider */}
      <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} /> 
        <Route path="/login" element={<LoginPage />} /> 
        <Route path="/signup" element={<SignupPage />} />
      </Routes>
    </BrowserRouter>
  </QueryClientProvider>
  </StrictMode>,
);
