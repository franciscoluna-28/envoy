import { useValidateEnvironmentConnection } from "../../hooks/useMigrations";
import { useEffect, useState } from "react";

type Props = {
  environmentId: string;
};

export function ConnectionStatus({ environmentId }: Props) {
  const [isConnected, setIsConnected] = useState<boolean | null>(null);
  const { mutate: validateConnection } = useValidateEnvironmentConnection();

  useEffect(() => {
    if (environmentId) {
      validateConnection(environmentId, {
        onSuccess: () => {
          setIsConnected(true);
        },
        onError: () => {
          setIsConnected(false);
        },
      });
    }
  }, [environmentId, validateConnection]);

  const getStatusIcon = () => (
    <div className={`h-1.5 w-1.5 rounded-full ${
      isConnected === null 
        ? "bg-yellow-500 animate-pulse" 
        : isConnected 
        ? "bg-green-500 animate-pulse" 
        : "bg-red-500"
    }`} />
  );

  const getStatusText = () => {
    if (isConnected === null) return "Checking...";
    return isConnected ? "Connected" : "Disconnected";
  };

  const getCardClassName = () => {
    const baseClass = "border p-4 rounded-2xl bg-card flex flex-col gap-1.5";
    return isConnected === false ? `${baseClass} border-red-200 bg-red-50/30` : baseClass;
  };

  const getTextClassName = () => {
    return isConnected === false ? "text-red-700" : "text-stone-900";
  };

  return (
    <div className={getCardClassName()}>
      <span className="text-[10px] font-medium text-muted-foreground">
        Network Status
      </span>
      <div className="flex items-center gap-2">
        {getStatusIcon()}
        <span className={`text-sm font-medium capitalize tracking-tight ${getTextClassName()}`}>
          {getStatusText()}
        </span>
      </div>
    </div>
  );
}
