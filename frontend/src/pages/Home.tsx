import { useEffect, useState } from "react";
import type { CommonPassport } from "../types";
import { PassportTable } from "../components/PassportTable";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "/api";

export default function Home() {
  const [passports, setPassports] = useState<CommonPassport[]>([]);

  useEffect(() => {
    fetch(`${API_BASE_URL}/v1/passports`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(setPassports)
      .catch(console.error);
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 p-6">
      {/* Светлый контейнер для таблицы */}
      <div className="bg-[#f9f9f9] border border-gray-300 rounded-3xl shadow-lg p-10 max-w-6xl mx-auto text-black">
        <h1 className="text-3xl md:text-4xl font-extrabold mb-10 text-center flex items-center justify-center gap-3">
          <span>📋</span> Паспорта сервисов
        </h1>

        {/* Таблица */}
        <PassportTable data={passports} />
      </div>
    </div>
  );
}
