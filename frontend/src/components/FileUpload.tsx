import { useRef, useState } from 'react'
import { Upload, X } from 'lucide-react'
import { cn } from '@/lib/utils'

interface FileUploadProps {
  label: string
  file: File | null
  onChange: (file: File | null) => void
  error?: string | null
}

export function FileUpload({ label, file, onChange, error }: FileUploadProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [preview, setPreview] = useState<string | null>(null)

  const handleSelect = (selected: File | null) => {
    onChange(selected)
    if (preview) URL.revokeObjectURL(preview)
    setPreview(selected ? URL.createObjectURL(selected) : null)
  }

  const handleClear = () => {
    handleSelect(null)
    if (inputRef.current) inputRef.current.value = ''
  }

  return (
    <div>
      <label className="mb-1.5 block text-sm font-medium text-gray-700">{label}</label>

      {preview ? (
        <div className="relative inline-block">
          <img
            src={preview}
            alt={`Preview ${label}`}
            className="h-40 w-full max-w-xs rounded-md border border-gray-200 object-cover"
          />
          <button
            type="button"
            onClick={handleClear}
            className="absolute -right-2 -top-2 rounded-full bg-red-500 p-1 text-white shadow hover:bg-red-600"
            title="Hapus"
          >
            <X className="h-4 w-4" />
          </button>
          <p className="mt-1 truncate text-xs text-gray-500">{file?.name}</p>
        </div>
      ) : (
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className={cn(
            'flex h-40 w-full max-w-xs flex-col items-center justify-center gap-2 rounded-md border-2 border-dashed bg-gray-50 text-gray-500 transition-colors hover:border-blue-400 hover:bg-blue-50',
            error ? 'border-red-400' : 'border-gray-300'
          )}
        >
          <Upload className="h-6 w-6" />
          <span className="text-sm">Klik untuk upload</span>
          <span className="text-xs text-gray-400">JPG / PNG, maks 2MB</span>
        </button>
      )}

      <input
        ref={inputRef}
        type="file"
        accept="image/jpeg,image/png"
        className="hidden"
        onChange={(e) => handleSelect(e.target.files?.[0] ?? null)}
      />

      {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
    </div>
  )
}
