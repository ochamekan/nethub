// ВНИМАНИЕ: в этом фрагменте есть несколько ошибок и плохих практик.
// Кандидату нужно:
// 1) Найти и описать проблемы.
// 2) Предложить, как переписать код лучше.

type Device = {
  id: number;
  hostname: string;
  ip: string;
};

// Имитация запроса к API
async function fetchDevices(): Promise<Device[]> {
  // Потенциальная проблема: игнорируются ошибки сети/HTTP-код
  // ИСПРАВЛЕНО: Обрабатываем ошибки
  const res = await fetch("/api/devices");
  if (!res.ok) {
    throw new Error(`HTTP error: ${res.status}`);
  }
  const devices = (await res.json()) as Device[];
  return devices;
}

// Глобальное состояние (антипаттерн для большинства приложений)
let devices: Device[] = [];
let isLoading = false;

// export async function loadAndFilterDevices(search: string)
// ИСПРАВЛЕНО: Нарушается принцип единной отвественности, делаем 2 отдельные ф-ии
export async function loadDevices() {
  // Потенциальная проблема: нет try/catch, при ошибке состояние "подвиснет"
  // ИСПРАВЛЕНО: Обрабатываем ошибки
  try {
    isLoading = true;
    const data = await fetchDevices();
    devices = data;
  } catch (error) {
    throw error;
  } finally {
    isLoading = false;
  }
}

// ИСПРАВЛЕНО: Теперь ф-ия чистая и ее можно тестировать, при одних и тех же данных ф-ия будет всегда возвращать один и тот же результат
export function filterDevices(search: string, devices: Device[]): Device[] {
  // Потенциальная проблема: сравнение без нормализации регистра и trim
  // ИСПРАВЛЕНО: Нормализируем
  const normalized = search.trim().toLowerCase();
  if (!normalized) return devices;
  return devices.filter((d) => d.hostname.toLowerCase().includes(normalized));
}

// Пример использования (упрощённо)
async function example() {
  const searchInput: HTMLInputElement | null =
    document.querySelector("#search");
  if (!searchInput) return;

  // Потенциальная проблема: нет debounce, каждый ввод символа может бить по API
  // ИСПРАВЛЕНО: Добавояем debounce
  const handleInput = debounce(async (searchVal: string) => {
    await loadDevices();
    const filteredList = filterDevices(searchVal, devices);
    console.log("Devices:", filteredList);
  }, 300);

  searchInput.addEventListener("input", (e: Event) =>
    handleInput((e.target as HTMLInputElement).value),
  );
}

function debounce<T extends (...args: never[]) => void>(fn: T, ms: number): T {
  let timer: ReturnType<typeof setTimeout>;
  return ((...args: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), ms);
  }) as T;
}

example();
