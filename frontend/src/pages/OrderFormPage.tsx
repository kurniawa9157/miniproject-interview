import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { AxiosError } from 'axios'
import api from '@/lib/axios'
import type { Order } from '@/lib/types'
import {
  validateWhatsapp,
  validatePlate,
  validateFrameNumber,
  validateFile,
} from '@/lib/validation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { FileUpload } from '@/components/FileUpload'

interface FieldErrors {
  whatsapp?: string | null
  plate?: string | null
  frame?: string | null
  ktp?: string | null
  stnk?: string | null
}

export function OrderFormPage() {
  const navigate = useNavigate()

  const [whatsapp, setWhatsapp] = useState('')
  const [plate, setPlate] = useState('')
  const [frame, setFrame] = useState('')
  const [ktp, setKtp] = useState<File | null>(null)
  const [stnk, setStnk] = useState<File | null>(null)

  const [errors, setErrors] = useState<FieldErrors>({})
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const validateAll = (): boolean => {
    const next: FieldErrors = {
      whatsapp: validateWhatsapp(whatsapp),
      plate: validatePlate(plate),
      frame: validateFrameNumber(frame),
      ktp: validateFile(ktp, 'KTP'),
      stnk: validateFile(stnk, 'STNK'),
    }
    setErrors(next)
    return Object.values(next).every((e) => e === null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitError(null)

    if (!validateAll()) return

    setSubmitting(true)
    try {
      const form = new FormData()
      form.append('whatsapp', whatsapp.trim())
      form.append('plate_number', plate.trim().toUpperCase())
      form.append('frame_number', frame.trim().toUpperCase())
      form.append('ktp', ktp!)
      form.append('stnk', stnk!)

      const { data } = await api.post<{ data: Order }>('/api/orders', form)
      navigate(`/orders/${data.data.id}`, { replace: true })
    } catch (err) {
      const axiosErr = err as AxiosError<{ error: string }>
      setSubmitError(axiosErr.response?.data?.error ?? 'Gagal membuat order. Coba lagi.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="mx-auto max-w-2xl">
      <Card>
        <CardHeader>
          <CardTitle>Buat Order Perpanjangan STNK</CardTitle>
          <CardDescription>
            Lengkapi data dan upload dokumen. Order akan diproses setelah verifikasi admin.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-5" noValidate>
            <div>
              <Label htmlFor="whatsapp">Nomor WhatsApp</Label>
              <Input
                id="whatsapp"
                placeholder="08123456789"
                value={whatsapp}
                onChange={(e) => setWhatsapp(e.target.value)}
                className="mt-1.5"
              />
              {errors.whatsapp && <p className="mt-1 text-sm text-red-600">{errors.whatsapp}</p>}
            </div>

            <div>
              <Label htmlFor="plate">Nomor Plat Kendaraan</Label>
              <Input
                id="plate"
                placeholder="D 1234 ABC"
                value={plate}
                onChange={(e) => setPlate(e.target.value.toUpperCase())}
                className="mt-1.5"
              />
              {errors.plate && <p className="mt-1 text-sm text-red-600">{errors.plate}</p>}
            </div>

            <div>
              <Label htmlFor="frame">5 Digit Terakhir Nomor Rangka</Label>
              <Input
                id="frame"
                placeholder="AB123"
                maxLength={5}
                value={frame}
                onChange={(e) => setFrame(e.target.value.toUpperCase())}
                className="mt-1.5"
              />
              {errors.frame && <p className="mt-1 text-sm text-red-600">{errors.frame}</p>}
            </div>

            <div className="grid gap-5 sm:grid-cols-2">
              <FileUpload label="Foto KTP" file={ktp} onChange={setKtp} error={errors.ktp} />
              <FileUpload label="Foto STNK Lama" file={stnk} onChange={setStnk} error={errors.stnk} />
            </div>

            {submitError && (
              <div className="rounded-md bg-red-50 p-3 text-sm text-red-700">{submitError}</div>
            )}

            <Button type="submit" className="w-full" size="lg" disabled={submitting}>
              {submitting ? 'Mengirim...' : 'Submit Order'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
