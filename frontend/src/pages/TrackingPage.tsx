import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { AxiosError } from 'axios'
import api from '@/lib/axios'
import type { OrderWithLogs } from '@/lib/types'
import { STATUS_LABELS } from '@/lib/types'
import { formatDateTime } from '@/lib/format'
import { StatusBadge } from '@/components/StatusBadge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export function TrackingPage() {
  const { id } = useParams()
  const [order, setOrder] = useState<OrderWithLogs | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true
    setLoading(true)
    api
      .get<{ data: OrderWithLogs }>(`/api/orders/${id}`)
      .then(({ data }) => {
        if (active) setOrder(data.data)
      })
      .catch((err: AxiosError<{ error: string }>) => {
        if (active) setError(err.response?.data?.error ?? 'Order tidak ditemukan')
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [id])

  if (loading) {
    return <div className="py-12 text-center text-gray-500">Memuat data order...</div>
  }

  if (error || !order) {
    return (
      <div className="py-12 text-center">
        <p className="text-gray-600">{error ?? 'Order tidak ditemukan'}</p>
        <Link to="/orders" className="mt-2 inline-block text-blue-600 hover:underline">
          Kembali ke daftar order
        </Link>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div>
        <Link to="/orders" className="text-sm text-blue-600 hover:underline">
          &larr; Daftar Order
        </Link>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>{order.id}</CardTitle>
            <StatusBadge status={order.status} />
          </div>
        </CardHeader>
        <CardContent className="space-y-3">
          <DetailRow label="Layanan" value="Perpanjangan STNK" />
          <DetailRow label="Tanggal Order" value={formatDateTime(order.created_at)} />
          <DetailRow label="Nomor Plat" value={order.plate_number} />
          <DetailRow label="Nomor WhatsApp" value={order.whatsapp} />
          <DetailRow label="5 Digit Rangka" value={order.frame_number} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Timeline Status</CardTitle>
        </CardHeader>
        <CardContent>
          <ol className="relative space-y-6 border-l border-gray-200 pl-6">
            {order.status_logs.map((log, idx) => {
              const isLatest = idx === order.status_logs.length - 1
              return (
                <li key={log.id} className="relative">
                  <span
                    className={
                      'absolute -left-[31px] flex h-4 w-4 items-center justify-center rounded-full ring-4 ring-white ' +
                      (isLatest ? 'bg-blue-600' : 'bg-gray-300')
                    }
                  />
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-900">
                      {STATUS_LABELS[log.status]}
                    </span>
                    <span className="text-xs text-gray-500">{formatDateTime(log.created_at)}</span>
                  </div>
                </li>
              )
            })}
          </ol>
        </CardContent>
      </Card>
    </div>
  )
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between border-b border-gray-100 pb-2 last:border-0">
      <span className="text-sm text-gray-500">{label}</span>
      <span className="text-sm font-medium text-gray-900">{value}</span>
    </div>
  )
}
