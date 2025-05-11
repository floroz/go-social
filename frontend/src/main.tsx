import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Routes, Route } from 'react-router'; // Attempting import from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'; // Import React Query
import './index.css';
// App is no longer directly routed for `/` here, HomePage will be.
// If App.tsx is meant to be a global layout, it should be structured differently or wrapped by specific page routes.
// For now, we assume HomePage is the primary content for `/`.
// import App from './App.tsx'; 
import LoginPage from './pages/LoginPage.tsx'; // Import LoginPage
import SignupPage from './pages/SignupPage.tsx'; // Import SignupPage
import HomePage from './pages/HomePage.tsx'; // Import HomePage
import ProtectedRoute from './components/ProtectedRoute.tsx'; // Import ProtectedRoute

// Create a client
const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}> {/* Wrap with provider */}
      <BrowserRouter>
        <Routes>
          <Route 
            path="/" 
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            } 
          />
          <Route path="/login" element={<LoginPage />} /> 
          <Route path="/signup" element={<SignupPage />} />
        </Routes>
      </BrowserRouter>
  </QueryClientProvider>
  </StrictMode>,
);
