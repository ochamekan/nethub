import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { DeviceFormDialog } from "@/components/device-form-dialog";
import { createDevice, type CreateDevicePayload } from "@/lib/api";
import { toast } from "sonner";

export default function CreateDeviceDialog() {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: createDevice,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      toast.success("Устройство добавлено");
      setOpen(false);
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : "Не удалось создать устройство",
      );
    },
  });

  return (
    <>
      <Button onClick={() => setOpen(true)}>Добавить устройство</Button>
      <DeviceFormDialog
        open={open}
        onOpenChange={setOpen}
        title="Новое устройство"
        submitLabel="Создать"
        isPending={mutation.isPending}
        isError={mutation.isError}
        error={mutation.error}
        onSubmit={(data: CreateDevicePayload) => mutation.mutate(data)}
      />
    </>
  );
}
