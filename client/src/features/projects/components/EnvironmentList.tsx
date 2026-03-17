import { EnvironmentCard } from './EnvironmentCard'
import { EmptyEnvironmentsState } from './EmptyEnvironmentsState'

interface EnvironmentListProps {
  environments: any[]
}

export function EnvironmentList({ environments }: EnvironmentListProps) {
  if (environments.length === 0) {
    return <EmptyEnvironmentsState />
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {environments.map((env) => (
        <EnvironmentCard key={env.id} env={env} />
      ))}
    </div>
  )
}
