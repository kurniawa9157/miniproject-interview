import { Link, Outlet, useNavigate } from 'react-router-dom'
import { LogOut } from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'

export function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login', { replace: true })
  }

  const homeLink = user?.is_admin ? '/admin' : '/orders/new'

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
          <Link to={homeLink} className="flex items-center gap-2">
            <span className="text-lg font-bold text-blue-600">JumpaPay</span>
            <span className="hidden text-sm text-gray-400 sm:inline">Perpanjangan STNK</span>
          </Link>

          <nav className="flex items-center gap-4">
            {user?.is_admin ? (
              <Link to="/admin" className="text-sm font-medium text-gray-600 hover:text-gray-900">
                Dashboard Admin
              </Link>
            ) : (
              <>
                <Link to="/orders/new" className="text-sm font-medium text-gray-600 hover:text-gray-900">
                  Buat Order
                </Link>
                <Link to="/orders" className="text-sm font-medium text-gray-600 hover:text-gray-900">
                  Order Saya
                </Link>
              </>
            )}

            {user && (
              <div className="flex items-center gap-2">
                {user.photo_url && (
                  <img src={user.photo_url} alt={user.name} className="h-8 w-8 rounded-full" />
                )}
                <span className="hidden text-sm text-gray-700 sm:inline">{user.name}</span>
                <Button variant="ghost" size="icon" onClick={handleLogout} title="Logout">
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            )}
          </nav>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
