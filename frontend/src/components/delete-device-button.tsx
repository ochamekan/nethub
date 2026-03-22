import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/device-form-dialog";
import { deleteDevice } from "@/lib/api";
import { Trash2 } from "lucide-react";
import { toast } from "sonner";

interface DeleteDeviceButtonProps {
  id: number;
  hostname: string;
}

export default function DeleteDeviceButton({
  id,
  hostname,
}: DeleteDeviceButtonProps) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => deleteDevice(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      toast.success(`${hostname} удалён`);
      setOpen(false);
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : "Не удалось удалить устройство",
      );
    },
  });

  return (
    <>
      <Button
        variant="ghost"
        size="icon"
        onClick={() => setOpen(true)}
        className="hover:text-destructive"
      >
        <Trash2 className="h-4 w-4" />
      </Button>
      <ConfirmDialog
        open={open}
        onOpenChange={setOpen}
        title="Удалить устройство"
        description={
          <>
            Вы уверены, что хотите удалить{" "}
            <span className="font-medium text-foreground">{hostname}</span>? Это
            действие необратимо.
          </>
        }
        confirmLabel="Удалить"
        isPending={mutation.isPending}
        onConfirm={() => mutation.mutate()}
      />
    </>
  );
}
