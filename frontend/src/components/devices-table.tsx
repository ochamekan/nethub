import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import type { Device } from "@/types/device";
import EditDeviceDialog from "./edit-device-dialog";
import DeleteDeviceButton from "./delete-device-button";
import DevicesTableSkeleton from "./device-table-skeleton";

interface DevicesTableProps {
  devices: Device[];
  isLoading?: boolean;
  isError?: boolean;
}

export default function DevicesTable({
  devices,
  isLoading,
  isError,
}: DevicesTableProps) {
  if (isLoading) {
    return (
      <div className="w-full">
        <DevicesTableSkeleton />
      </div>
    );
  }
  if (isError) {
    return (
      <div className="w-full py-10 text-center text-destructive">
        Не удалось загрузить устройства.
      </div>
    );
  }

  return (
    <div className="w-full">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Хост</TableHead>
            <TableHead>IP-адрес</TableHead>
            <TableHead>Расположение</TableHead>
            <TableHead>Статус</TableHead>
            <TableHead>Дата добавления</TableHead>
            <TableHead className="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {devices.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={6}
                className="text-center text-muted-foreground"
              >
                Устройства не найдены.
              </TableCell>
            </TableRow>
          ) : (
            devices.map((device) => (
              <TableRow
                key={device.id}
                className={device.is_deleted ? "opacity-40" : ""}
              >
                <TableCell className="font-medium">{device.hostname}</TableCell>
                <TableCell className="font-mono text-sm">{device.ip}</TableCell>
                <TableCell>{device.location}</TableCell>
                <TableCell>
                  <Badge
                    className={
                      device.is_active
                        ? "bg-emerald-500/15 text-emerald-600 border-emerald-500/20 hover:bg-emerald-500/15"
                        : "bg-rose-500/15 text-rose-600 border-rose-500/20 hover:bg-rose-500/15"
                    }
                    variant="outline"
                  >
                    {device.is_active ? "Активен" : "Неактивен"}
                  </Badge>{" "}
                </TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {new Date(device.created_at).toLocaleDateString("ru-RU")}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex items-center justify-end gap-1">
                    <EditDeviceDialog device={device} />
                    <DeleteDeviceButton
                      id={device.id}
                      hostname={device.hostname}
                    />
                  </div>
                </TableCell>{" "}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
