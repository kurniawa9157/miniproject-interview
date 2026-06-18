// Mirror of backend validation rules (internal/handler/order_handler.go).

export const MAX_FILE_SIZE = 2 * 1024 * 1024 // 2MB
export const ACCEPTED_IMAGE_TYPES = ['image/jpeg', 'image/png']

// WhatsApp: diawali 08 atau +62, total 10-15 digit angka.
const reWA = /^(\+62|08)\d{8,13}$/

// Plat Indonesia: huruf, spasi, angka, spasi, huruf (contoh: D 1234 ABC)
const rePlate = /^[A-Z]{1,2}\s\d{1,4}\s[A-Z]{1,3}$/

export function validateWhatsapp(value: string): string | null {
  const v = value.trim()
  if (!v) return 'Nomor WhatsApp wajib diisi'
  if (!reWA.test(v)) {
    return 'Format tidak valid. Contoh: 08123456789 atau +6281234567890'
  }
  return null
}

export function validatePlate(value: string): string | null {
  const v = value.trim().toUpperCase()
  if (!v) return 'Nomor plat wajib diisi'
  if (!rePlate.test(v)) {
    return 'Format tidak valid. Contoh: D 1234 ABC'
  }
  return null
}

export function validateFrameNumber(value: string): string | null {
  const v = value.trim()
  if (!v) return 'Nomor rangka wajib diisi'
  if (v.length !== 5) return 'Nomor rangka harus tepat 5 karakter'
  if (!/^[a-zA-Z0-9]{5}$/.test(v)) return 'Nomor rangka hanya boleh huruf dan angka'
  return null
}

export function validateFile(file: File | null, label: string): string | null {
  if (!file) return `${label} wajib diupload`
  if (!ACCEPTED_IMAGE_TYPES.includes(file.type)) {
    return `${label} harus berformat JPG atau PNG`
  }
  if (file.size > MAX_FILE_SIZE) {
    return `Ukuran ${label} melebihi 2MB`
  }
  return null
}
