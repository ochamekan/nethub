import { useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getDevices } from "./lib/api";
import DevicesTable from "./components/devices-table";
import Search from "./components/search";
import Title from "./components/title";
import DevicesToolbar from "./components/devices-toolbar";
import ThemeToggle from "./components/theme-toggle";

export function App() {
  const [searchParams] = useSearchParams();

  const search = searchParams.get("search") ?? undefined;
  const onlyActive = searchParams.get("is_active") === "true";

  const {
    data: devices = [],
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["devices", search, onlyActive],
    queryFn: () => getDevices(search, onlyActive),
  });

  return (
    <div className="mx-auto my-0 w-full max-w-300 px-4 py-10">
      <div className="flex flex-col gap-10">
        <ThemeToggle />
        <Title />
        <Search />
        <DevicesToolbar />
        <DevicesTable
          devices={devices}
          isLoading={isLoading}
          isError={isError}
        />
      </div>
    </div>
  );
}

export default App;
