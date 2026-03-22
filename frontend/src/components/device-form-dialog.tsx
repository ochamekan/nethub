import { useEffect } from "react";
import { useForm, Controller } from "react-hook-form";
import { IMaskInput } from "react-imask";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Loader2 } from "lucide-react";
import { hostnameMask, locationMask, ipMask } from "@/lib/input-filters";
import type { CreateDevicePayload } from "@/lib/api";
import { cn } from "@/lib/utils";

interface DeviceFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  defaultValues?: Partial<CreateDevicePayload>;
  isPending: boolean;
  isError: boolean;
  error: Error | null;
  submitLabel: string;
  onSubmit: (data: CreateDevicePayload) => void;
}

export function DeviceFormDialog({
  open,
  onOpenChange,
  title,
  defaultValues,
  isPending,
  isError,
  error,
  submitLabel,
  onSubmit,
}: DeviceFormDialogProps) {
  const {
    control,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateDevicePayload>({
    defaultValues: { is_active: true, ...defaultValues },
  });

  useEffect(() => {
    if (open) reset({ is_active: true, ...defaultValues });
  }, [open]);

  const inputClass = cn(
    "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1",
    "text-base shadow-xs outline-none transition-colors md:text-sm",
    "placeholder:text-muted-foreground",
    "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
    "disabled:cursor-not-allowed disabled:opacity-50",
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
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

          {isError && (
            <p className="text-xs text-destructive">
              {error instanceof Error ? error.message : "Неизвестная ошибка"}
            </p>
          )}

          <Button type="submit" disabled={isPending}>
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {isPending ? "Сохранение..." : submitLabel}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}

interface ConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: React.ReactNode;
  confirmLabel: string;
  isPending: boolean;
  onConfirm: () => void;
}

export function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel,
  isPending,
  onConfirm,
}: ConfirmDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription asChild>
            <div>{description}</div>
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-2 pt-2">
          <Button
            variant="destructive"
            disabled={isPending}
            onClick={onConfirm}
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {isPending ? "Удаление..." : confirmLabel}
          </Button>
          <Button
            variant="outline"
            disabled={isPending}
            onClick={() => onOpenChange(false)}
          >
            Отмена
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
