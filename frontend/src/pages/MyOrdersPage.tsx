import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '@/lib/axios'
import type { Order } from '@/lib/types'
import { formatDate } from '@/lib/format'
import { StatusBadge } from '@/components/StatusBadge'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export function MyOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true
    api
      .get<{ data: Order[] }>('/api/orders')
      .then(({ data }) => {
        if (active) setOrders(data.data)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Order Saya</h1>
        <Button asChild>
          <Link to="/orders/new">+ Buat Order</Link>
        </Button>
      </div>

      {loading ? (
        <div className="py-12 text-center text-gray-500">Memuat...</div>
      ) : orders.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center text-gray-500">
            <p>Belum ada order.</p>
            <Button asChild className="mt-4">
              <Link to="/orders/new">Buat Order Pertama</Link>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {orders.map((order) => (
            <Link key={order.id} to={`/orders/${order.id}`}>
              <Card className="transition-shadow hover:shadow-md">
                <CardContent className="flex items-center justify-between py-4">
                  <div>
                    <div className="font-semibold text-gray-900">{order.id}</div>
                    <div className="mt-1 text-sm text-gray-500">
                      {order.plate_number} &middot; {formatDate(order.created_at)}
                    </div>
                  </div>
                  <StatusBadge status={order.status} />
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
