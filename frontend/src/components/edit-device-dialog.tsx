import { useEffect, useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { IMaskInput } from "react-imask";
import { updateDevice, type CreateDevicePayload } from "@/lib/api";
import { hostnameMask, locationMask, ipMask } from "@/lib/input-filters";
import type { Device } from "@/types/device";
import { Pencil } from "lucide-react";
import { cn } from "@/lib/utils";

interface EditDeviceDialogProps {
  device: Device;
}

export default function EditDeviceDialog({ device }: EditDeviceDialogProps) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const {
    control,
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

  const inputClass = cn(
    "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1",
    "text-base shadow-xs outline-none transition-colors md:text-sm",
    "placeholder:text-muted-foreground",
    "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
    "disabled:cursor-not-allowed disabled:opacity-50",
  );

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
            <Controller
              name="hostname"
              control={control}
              rules={{ required: "Хост обязателен" }}
              render={({ field }) => (
                <IMaskInput
                  {...hostnameMask}
                  value={field.value ?? ""}
                  onAccept={(val) => field.onChange(val)}
                  placeholder="edge-router-01"
                  className={inputClass}
                />
              )}
            />
            {errors.hostname && (
              <p className="text-xs text-destructive">
                {errors.hostname.message}
              </p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">IP-адрес</label>
            <Controller
              name="ip"
              control={control}
              rules={{ required: "IP-адрес обязателен" }}
              render={({ field }) => (
                <IMaskInput
                  {...ipMask}
                  value={field.value ?? ""}
                  onAccept={(val) => field.onChange(val)}
                  placeholder="192.168.1.1"
                  className={inputClass}
                />
              )}
            />
            {errors.ip && (
              <p className="text-xs text-destructive">{errors.ip.message}</p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">Расположение</label>
            <Controller
              name="location"
              control={control}
              rules={{ required: "Расположение обязательно" }}
              render={({ field }) => (
                <IMaskInput
                  {...locationMask}
                  value={field.value ?? ""}
                  onAccept={(val) => field.onChange(val)}
                  placeholder="Germany"
                  className={inputClass}
                />
              )}
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
