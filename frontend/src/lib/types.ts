export interface User {
  id: string
  name: string
  email: string
  photo_url: string
  is_admin: boolean
  created_at: string
}

export type OrderStatus = 'PENDING' | 'IN_PROCESS' | 'DONE' | 'CANCELLED'

export interface OrderStatusLog {
  id: string
  order_id: string
  status: OrderStatus
  changed_by: string | null
  created_at: string
}

export type PaymentStatus = 'UNPAID' | 'PENDING' | 'PAID' | 'FAILED'

export interface Order {
  id: string
  user_id: string
  whatsapp: string
  plate_number: string
  frame_number: string
  ktp_url: string
  stnk_url: string
  status: OrderStatus
  amount: number
  payment_status: PaymentStatus
  created_at: string
  updated_at: string
}

export interface OrderWithLogs extends Order {
  status_logs: OrderStatusLog[]
  user_name: string
  user_email: string
}

export const STATUS_LABELS: Record<OrderStatus, string> = {
  PENDING: 'Menunggu Verifikasi',
  IN_PROCESS: 'Sedang Diproses',
  DONE: 'Selesai',
  CANCELLED: 'Dibatalkan',
}

// Valid transitions mirror the backend rules (model.ValidTransitions).
export const VALID_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  PENDING: ['IN_PROCESS', 'CANCELLED'],
  IN_PROCESS: ['DONE', 'CANCELLED'],
  DONE: [],
  CANCELLED: [],
}
