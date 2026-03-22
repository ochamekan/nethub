import { Search as SearchIcon } from "lucide-react";
import { Input } from "./ui/input";
import { useSearchParams } from "react-router-dom";
import { useDebouncedCallback } from "use-debounce";

export default function Search() {
  const [searchParams, setSearchParams] = useSearchParams();

  const handleSearch = useDebouncedCallback((value: string) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (value) {
        next.set("search", value);
      } else {
        next.delete("search");
      }
      return next;
    });
  }, 300);

  return (
    <div className="relative w-full">
      <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
      <Input
        placeholder="Поиск по хосту..."
        defaultValue={searchParams.get("search") ?? ""}
        onChange={(e) => handleSearch(e.target.value)}
        className="pl-9"
      />
    </div>
  );
}
