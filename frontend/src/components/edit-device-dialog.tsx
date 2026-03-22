import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { DeviceFormDialog } from "@/components/device-form-dialog";
import { updateDevice, type CreateDevicePayload } from "@/lib/api";
import type { Device } from "@/types/device";
import { Pencil } from "lucide-react";
import { toast } from "sonner";

interface EditDeviceDialogProps {
  device: Device;
}

export default function EditDeviceDialog({ device }: EditDeviceDialogProps) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (data: CreateDevicePayload) => updateDevice(device.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      toast.success("Устройство обновлено");
      setOpen(false);
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : "Не удалось обновить устройство",
      );
    },
  });

  return (
    <>
      <Button variant="ghost" size="icon" onClick={() => setOpen(true)}>
        <Pencil className="h-4 w-4" />
      </Button>
      <DeviceFormDialog
        open={open}
        onOpenChange={setOpen}
        title="Редактировать устройство"
        submitLabel="Сохранить"
        defaultValues={{
          hostname: device.hostname,
          ip: device.ip,
          location: device.location,
          is_active: device.is_active,
        }}
        isPending={mutation.isPending}
        isError={mutation.isError}
        error={mutation.error}
        onSubmit={(data: CreateDevicePayload) => mutation.mutate(data)}
      />
    </>
  );
}
