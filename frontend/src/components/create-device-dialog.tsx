import { useState } from "react";
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
import { createDevice, type CreateDevicePayload } from "@/lib/api";

export default function CreateDeviceDialog() {
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
    defaultValues: { is_active: true },
  });

  const mutation = useMutation({
    mutationFn: createDevice,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      reset();
      setOpen(false);
    },
  });

  const onSubmit = (data: CreateDevicePayload) => mutation.mutate(data);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>Добавить устройство</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Новое устройство</DialogTitle>
        </DialogHeader>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-col gap-4 pt-2"
        >
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">Хост</label>
            <Input
              placeholder="edge-router-01"
              {...register("hostname", { required: "Хост обязателен" })}
            />
            {errors.hostname && (
              <p className="text-xs text-destructive">
                {errors.hostname.message}
              </p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">IP-адрес</label>
            <Input
              placeholder="192.168.1.1"
              {...register("ip", { required: "IP-адрес обязателен" })}
            />
            {errors.ip && (
              <p className="text-xs text-destructive">{errors.ip.message}</p>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium">Расположение</label>
            <Input
              placeholder="Germany"
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
            {mutation.isPending ? "Сохранение..." : "Создать"}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
