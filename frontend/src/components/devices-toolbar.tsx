import CreateDeviceDialog from "./create-device-dialog";
import { Switch } from "@/components/ui/switch";
import { useSearchParams } from "react-router-dom";

export default function DevicesToolbar() {
  const [searchParams, setSearchParams] = useSearchParams();
  const onlyActive = searchParams.get("is_active") === "true";

  const handleToggle = (value: boolean) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (value) {
        next.set("is_active", "true");
      } else {
        next.delete("is_active");
      }
      return next;
    });
  };

  return (
    <div className="flex items-center justify-between">
      <div
        className="flex items-center gap-2 cursor-pointer select-none"
        onClick={() => handleToggle(!onlyActive)}
      >
        <Switch checked={onlyActive} onCheckedChange={handleToggle} />
        <span className="text-sm">Только активные</span>
      </div>
      <CreateDeviceDialog />
    </div>
  );
}
