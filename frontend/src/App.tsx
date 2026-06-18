import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from '@/hooks/AuthProvider'
import { ProtectedRoute } from '@/components/ProtectedRoute'
import { Layout } from '@/components/Layout'
import { LoginPage } from '@/pages/LoginPage'
import { OrderFormPage } from '@/pages/OrderFormPage'
import { MyOrdersPage } from '@/pages/MyOrdersPage'
import { TrackingPage } from '@/pages/TrackingPage'
import { AdminDashboard } from '@/pages/admin/AdminDashboard'
import { AdminOrderDetail } from '@/pages/admin/AdminOrderDetail'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />

          {/* Customer routes */}
          <Route element={<ProtectedRoute />}>
            <Route element={<Layout />}>
              <Route path="/orders/new" element={<OrderFormPage />} />
              <Route path="/orders" element={<MyOrdersPage />} />
              <Route path="/orders/:id" element={<TrackingPage />} />
            </Route>
          </Route>

          {/* Admin routes */}
          <Route element={<ProtectedRoute adminOnly />}>
            <Route element={<Layout />}>
              <Route path="/admin" element={<AdminDashboard />} />
              <Route path="/admin/orders/:id" element={<AdminOrderDetail />} />
            </Route>
          </Route>

          <Route path="/" element={<Navigate to="/orders/new" replace />} />
          <Route path="*" element={<Navigate to="/orders/new" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
