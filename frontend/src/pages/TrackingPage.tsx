import { useParams } from 'react-router-dom'

export function TrackingPage() {
  const { id } = useParams()
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Tracking Order {id}</h1>
      <p className="mt-2 text-gray-500">Halaman tracking akan dibangun di Phase 6.</p>
    </div>
  )
}
