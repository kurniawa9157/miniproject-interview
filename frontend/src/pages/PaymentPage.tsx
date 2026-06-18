import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { AxiosError } from 'axios'
import api from '@/lib/axios'
import type { OrderWithLogs } from '@/lib/types'
import { formatRupiah } from '@/lib/format'
import { loadSnap, openSnap } from '@/lib/snap'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export function PaymentPage() {
  const { id } = useParams()
  const navigate = useNavigate()

  const [order, setOrder] = useState<OrderWithLogs | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [paying, setPaying] = useState(false)

  useEffect(() => {
    let active = true
    Promise.all([
      api.get<{ data: OrderWithLogs }>(`/api/orders/${id}`),
      api.get<{ client_key: string }>('/api/payment/config'),
    ])
      .then(async ([orderRes, configRes]) => {
        if (!active) return
        setOrder(orderRes.data.data)
        if (configRes.data.client_key) {
          await loadSnap(configRes.data.client_key)
        }
      })
      .catch(() => {
        if (active) setError('Gagal memuat data pembayaran')
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [id])

  const handlePay = async () => {
    setPaying(true)
    setError(null)
    try {
      const { data } = await api.post<{ token: string }>(`/api/orders/${id}/pay`)
      openSnap(data.token, {
        onSuccess: () => navigate(`/orders/${id}`, { replace: true }),
        onPending: () => navigate(`/orders/${id}`, { replace: true }),
        onError: () => {
          setError('Pembayaran gagal. Silakan coba lagi.')
          setPaying(false)
        },
        onClose: () => setPaying(false),
      })
    } catch (err) {
      const axiosErr = err as AxiosError<{ error: string }>
      setError(axiosErr.response?.data?.error ?? 'Gagal memulai pembayaran')
      setPaying(false)
    }
  }

  if (loading) {
    return <div className="py-12 text-center text-gray-500">Memuat...</div>
  }

  if (error && !order) {
    return (
      <div className="py-12 text-center">
        <p className="text-gray-600">{error}</p>
        <Link to="/orders" className="mt-2 inline-block text-blue-600 hover:underline">
          Kembali ke daftar order
        </Link>
      </div>
    )
  }

  if (!order) return null

  const alreadyPaid = order.payment_status === 'PAID'

  return (
    <div className="mx-auto max-w-md">
      <Card>
        <CardHeader className="text-center">
          <CardTitle>Konfirmasi Pembayaran</CardTitle>
          <CardDescription>Order {order.id}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-5">
          <div className="rounded-lg bg-gray-50 p-4">
            <div className="flex justify-between border-b border-gray-200 pb-2 text-sm">
              <span className="text-gray-500">Layanan</span>
              <span className="font-medium">Perpanjangan STNK</span>
            </div>
            <div className="flex justify-between border-b border-gray-200 py-2 text-sm">
              <span className="text-gray-500">Nomor Plat</span>
              <span className="font-medium">{order.plate_number}</span>
            </div>
            <div className="flex items-center justify-between pt-3">
              <span className="text-gray-700">Total</span>
              <span className="text-xl font-bold text-blue-600">{formatRupiah(order.amount)}</span>
            </div>
          </div>

          {error && (
            <div className="rounded-md bg-red-50 p-3 text-sm text-red-700">{error}</div>
          )}

          {alreadyPaid ? (
            <div className="rounded-md bg-green-50 p-3 text-center text-sm text-green-700">
              Pembayaran sudah lunas. Order sedang diproses.
            </div>
          ) : (
            <Button onClick={handlePay} className="w-full" size="lg" disabled={paying}>
              {paying ? 'Memproses...' : 'Bayar Sekarang'}
            </Button>
          )}

          <div className="text-center">
            <Link to={`/orders/${id}`} className="text-sm text-gray-500 hover:underline">
              Lewati, bayar nanti
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
