import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Home from './pages/Home';
import Header from './components/Header';
import './App.css';

// Создаем клиент для React Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 минут
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
        <QueryClientProvider client={queryClient}>
          <Router>
            <Header/>
            <Routes>
              <Route path="/" element={<Home/>}/>
            </Routes>
          </Router>
        </QueryClientProvider>
  );
}

export default App;