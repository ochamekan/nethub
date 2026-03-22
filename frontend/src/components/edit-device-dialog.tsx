import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { updateDevice, type CreateDevicePayload } from "@/lib/api";
import type { Device } from "@/types/device";
import { Pencil } from "lucide-react";

interface EditDeviceDialogProps {
  device: Device;
}

export default function EditDeviceDialog({ device }: EditDeviceDialogProps) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateDevicePayload>({
    defaultValues: {
      hostname: device.hostname,
      ip: device.ip,
      location: device.location,
      is_active: device.is_active,
    },
  });

  useEffect(() => {
    if (open) {
      reset({
        hostname: device.hostname,
        ip: device.ip,
        location: device.location,
        is_active: device.is_active,
      });
    }
  }, [open, device, reset]);

  const mutation = useMutation({
    mutationFn: (data: CreateDevicePayload) => updateDevice(device.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      setOpen(false);
    },
  });

  const onSubmit = (data: CreateDevicePayload) => mutation.mutate(data);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon">
          <Pencil className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Редактировать устройство</DialogTitle>
        </DialogHeader>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-4 pt-2"
        >
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">Хост</label>
            <Input {...register("hostname", { required: "Хост обязателен" })} />
            {errors.hostname && (
              <p className="text-xs text-destructive">
                {errors.hostname.message}
              </p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">IP-адрес</label>
            <Input {...register("ip", { required: "IP-адрес обязателен" })} />
            {errors.ip && (
              <p className="text-xs text-destructive">{errors.ip.message}</p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">Расположение</label>
            <Input
              {...register("location", {
                required: "Расположение обязательно",
              })}
            />
            {errors.location && (
              <p className="text-xs text-destructive">
                {errors.location.message}
              </p>
            )}
          </div>

          <div className="flex items-center justify-between">
            <label className="text-sm font-medium">Активно</label>
            <Switch
              checked={watch("is_active")}
              onCheckedChange={(val) => setValue("is_active", val)}
            />
          </div>

          {mutation.isError && (
            <p className="text-xs text-destructive">
              {mutation.error instanceof Error
                ? mutation.error.message
                : "Неизвестная ошибка"}
            </p>
          )}

          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Сохранение..." : "Сохранить"}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
