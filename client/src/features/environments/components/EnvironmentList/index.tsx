import { useState } from 'react'
import type { Environment, EnviromentUpdateInput } from '@/features/types'
import { EmptyEnvironmentsState } from '../EmptyEnvironmentState'
import { EnvironmentCard } from '../EnvironmentCard'
import { UpdateEnvironmentModal } from '../UpdateEnvironmentModal'

interface EnvironmentListProps {
  environments: Environment[]
}

export function EnvironmentList({ environments }: EnvironmentListProps) {
  const [updateModalOpen, setUpdateModalOpen] = useState(false)
  const [selectedEnvironment, setSelectedEnvironment] = useState<EnviromentUpdateInput | null>(null)

  const handleUpdate = (env: Environment) => {
    setSelectedEnvironment({
      id: env.id as string,
      name: env.name || '',
    })
    setUpdateModalOpen(true)
  }

  const handleUpdated = () => {
    setUpdateModalOpen(false)
    setSelectedEnvironment(null)
  }

  if (environments.length === 0) {
    return <EmptyEnvironmentsState />
  }

  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {environments.map((env) => (
          <EnvironmentCard 
            key={env.id} 
            env={env}
            onUpdate={() => handleUpdate(env)}
          />
        ))}
      </div>

      {selectedEnvironment && (
        <UpdateEnvironmentModal
          environment={selectedEnvironment}
          open={updateModalOpen}
          onOpenChange={setUpdateModalOpen}
          onUpdated={handleUpdated}
        />
      )}
    </>
  )
}
