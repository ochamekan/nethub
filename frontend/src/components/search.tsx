import { Input } from "./ui/input";
import { useSearchParams } from "react-router-dom";
import { useDebouncedCallback } from "use-debounce";

export default function Search() {
  const [searchParams, setSearchParams] = useSearchParams();

  const handleSearch = useDebouncedCallback((value: string) => {
    setSearchParams(value ? { search: value } : {});
  }, 400);

  return (
    <div className="w-full">
      <div className="w-full flex gap-2">
        <Input
          placeholder="Введите имя хоста"
          defaultValue={searchParams.get("search") ?? ""}
          onChange={(e) => handleSearch(e.target.value)}
        />
      </div>
    </div>
  );
}
