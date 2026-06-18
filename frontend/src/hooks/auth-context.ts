import { createContext } from 'react'
import type { User } from '@/lib/types'

export interface AuthContextValue {
  user: User | null
  loading: boolean
  refetch: () => Promise<void>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | undefined>(undefined)
