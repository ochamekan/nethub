import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { deleteDevice } from "@/lib/api";
import { Trash2 } from "lucide-react";

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
      setOpen(false);
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
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Удалить устройство</DialogTitle>
            <DialogDescription>
              Вы уверены, что хотите удалить{" "}
              <span className="font-medium text-foreground">{hostname}</span>?
              Это действие необратимо.
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-2 pt-2">
            <Button
              variant="destructive"
              disabled={mutation.isPending}
              onClick={() => mutation.mutate()}
            >
              {mutation.isPending ? "Удаление..." : "Удалить"}
            </Button>
            <Button
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={mutation.isPending}
            >
              Отмена
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
