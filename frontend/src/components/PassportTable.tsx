import type { CommonPassport } from "../types";

const severityStyles: Record<string, string> = {
  critical: "bg-red-600/20 text-red-500",
  important: "bg-yellow-500/20 text-yellow-500",
  normal: "bg-green-500/20 text-green-500",
};

export function PassportTable({ data }: { data: CommonPassport[] }) {
  return (
    <div className="overflow-x-auto rounded-xl shadow p-4 bg-[#f9f9f9] !text-black">
      <table className="min-w-[900px] mx-auto border-collapse bg-[#f9f9f9] !text-black">
        <thead>
          <tr className="bg-gray-200 text-center text-sm uppercase tracking-wider !text-black">
            <th className="px-8 py-4">Сервис</th>
            <th className="px-8 py-4">Инфра</th>
            <th className="px-8 py-4">Host</th>
            <th className="px-8 py-4">Cluster</th>
            <th className="px-8 py-4">Namespace</th>
            <th className="px-8 py-4">Версия</th>
            <th className="px-8 py-4">Severity</th>
          </tr>
        </thead>
        <tbody>
          {data.map((item, idx) => (
            <tr
              key={idx}
              className={`${idx % 2 === 0 ? "bg-gray-50" : "bg-gray-100"} border-t border-gray-300 transition hover:bg-gray-200 !text-black`}
            >
              <td className="px-8 py-4 text-center font-semibold">{item.service_type}</td>
              <td className="px-8 py-4 text-center">{item.infrastructure.infra_type}</td>
              <td className="px-8 py-4 text-center">{item.infrastructure.host || "—"}</td>
              <td className="px-8 py-4 text-center">{item.infrastructure.cluster || "—"}</td>
              <td className="px-8 py-4 text-center">{item.infrastructure.namespace || "—"}</td>
              <td className="px-8 py-4 text-center">{item.version}</td>
              <td className="px-8 py-4 text-center">
                <span
                  className={`px-4 py-1 rounded text-sm font-semibold ${severityStyles[item.severity]}`}
                >
                  {item.severity}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
