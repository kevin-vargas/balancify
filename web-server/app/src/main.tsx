import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import Layout from "./components/Layout"
import "./index.css"
import Loading from './components/Loading'
import Config from './config';
const Csv = React.lazy(() => import('./pages/Csv'))
const Hello = React.lazy(() => import('./pages/Hello'))
const ErrorPage = React.lazy(() => import('./pages/Error'))
const Home = React.lazy(() => import('./pages/Home'))
import { QueryClient, QueryClientProvider } from 'react-query'
const queryClient = new QueryClient()

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "/csv",
        element: <Csv />
      },
      {
        path: "/hello",
        element: <Hello />,
      },
      {
        path: "/*",
        element: <Navigate to="/"/>,
      },
    ],
    errorElement: <ErrorPage />,
  },
], {basename: Config.uriPrefix});

const Main = () => {
  return (
    <React.StrictMode>
      <QueryClientProvider client={queryClient}>
        <React.Suspense fallback={<Loading />}>
          <RouterProvider router={router} />
        </React.Suspense>
      </QueryClientProvider>
    </React.StrictMode>
  )
}

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
          <Main />

)
