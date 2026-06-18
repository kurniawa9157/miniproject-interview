const SNAP_SRC = 'https://app.sandbox.midtrans.com/snap/snap.js'

interface SnapResult {
  order_id: string
  transaction_status: string
}

interface SnapCallbacks {
  onSuccess?: (result: SnapResult) => void
  onPending?: (result: SnapResult) => void
  onError?: (result: unknown) => void
  onClose?: () => void
}

interface SnapAPI {
  pay: (token: string, callbacks: SnapCallbacks) => void
}

declare global {
  interface Window {
    snap?: SnapAPI
  }
}

let loadingPromise: Promise<void> | null = null

// loadSnap injects the Snap.js script (once) with the given client key.
export function loadSnap(clientKey: string): Promise<void> {
  if (window.snap) return Promise.resolve()
  if (loadingPromise) return loadingPromise

  loadingPromise = new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = SNAP_SRC
    script.setAttribute('data-client-key', clientKey)
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Gagal memuat Midtrans Snap.js'))
    document.body.appendChild(script)
  })

  return loadingPromise
}

export function openSnap(token: string, callbacks: SnapCallbacks) {
  if (!window.snap) {
    callbacks.onError?.(new Error('Snap belum siap'))
    return
  }
  window.snap.pay(token, callbacks)
}
