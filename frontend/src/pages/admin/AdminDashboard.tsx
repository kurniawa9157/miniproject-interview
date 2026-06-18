import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '@/lib/axios'
import type { OrderStatus, OrderWithLogs } from '@/lib/types'
import { STATUS_LABELS } from '@/lib/types'
import { formatDate } from '@/lib/format'
import { StatusBadge } from '@/components/StatusBadge'
import { Card, CardContent } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const ALL = 'ALL'

export function AdminDashboard() {
  const navigate = useNavigate()
  const [orders, setOrders] = useState<OrderWithLogs[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<string>(ALL)

  useEffect(() => {
    let active = true
    setLoading(true)
    const query = filter === ALL ? '' : `?status=${filter}`
    api
      .get<{ data: OrderWithLogs[] }>(`/api/admin/orders${query}`)
      .then(({ data }) => {
        if (active) setOrders(data.data)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [filter])

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard Admin</h1>
        <div className="w-56">
          <Select value={filter} onValueChange={setFilter}>
            <SelectTrigger>
              <SelectValue placeholder="Filter status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={ALL}>Semua Status</SelectItem>
              {(Object.keys(STATUS_LABELS) as OrderStatus[]).map((s) => (
                <SelectItem key={s} value={s}>
                  {STATUS_LABELS[s]}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="py-12 text-center text-gray-500">Memuat...</div>
          ) : orders.length === 0 ? (
            <div className="py-12 text-center text-gray-500">Tidak ada order.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50 text-left text-gray-600">
                    <th className="px-4 py-3 font-medium">Order ID</th>
                    <th className="px-4 py-3 font-medium">Customer</th>
                    <th className="px-4 py-3 font-medium">Plat</th>
                    <th className="px-4 py-3 font-medium">WhatsApp</th>
                    <th className="px-4 py-3 font-medium">Tanggal</th>
                    <th className="px-4 py-3 font-medium">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((order) => (
                    <tr
                      key={order.id}
                      onClick={() => navigate(`/admin/orders/${order.id}`)}
                      className="cursor-pointer border-b border-gray-100 hover:bg-gray-50"
                    >
                      <td className="px-4 py-3 font-medium text-gray-900">{order.id}</td>
                      <td className="px-4 py-3 text-gray-700">{order.user_name}</td>
                      <td className="px-4 py-3 text-gray-700">{order.plate_number}</td>
                      <td className="px-4 py-3 text-gray-700">{order.whatsapp}</td>
                      <td className="px-4 py-3 text-gray-700">{formatDate(order.created_at)}</td>
                      <td className="px-4 py-3">
                        <StatusBadge status={order.status} />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
