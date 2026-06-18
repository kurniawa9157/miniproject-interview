import { Badge } from '@/components/ui/badge'
import { STATUS_LABELS, type OrderStatus } from '@/lib/types'

const VARIANT_MAP: Record<OrderStatus, 'pending' | 'in_process' | 'done' | 'cancelled'> = {
  PENDING: 'pending',
  IN_PROCESS: 'in_process',
  DONE: 'done',
  CANCELLED: 'cancelled',
}

export function StatusBadge({ status }: { status: OrderStatus }) {
  return <Badge variant={VARIANT_MAP[status]}>{STATUS_LABELS[status]}</Badge>
}
