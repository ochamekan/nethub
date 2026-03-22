import CreateDeviceDialog from "./create-device-dialog";
import { Switch } from "@/components/ui/switch";

interface DevicesToolbarProps {
  onlyActive: boolean;
  onOnlyActiveChange: (value: boolean) => void;
}

export default function DevicesToolbar({
  onlyActive,
  onOnlyActiveChange,
}: DevicesToolbarProps) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2">
        <Switch checked={onlyActive} onCheckedChange={onOnlyActiveChange} />
        <span className="text-sm text-muted-foreground">Только активные</span>
      </div>
      <CreateDeviceDialog />
    </div>
  );
}
