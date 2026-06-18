import { useParams } from 'react-router-dom'

export function AdminOrderDetail() {
  const { id } = useParams()
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Detail Order {id}</h1>
      <p className="mt-2 text-gray-500">Detail order admin akan dibangun di Phase 7.</p>
    </div>
  )
}
