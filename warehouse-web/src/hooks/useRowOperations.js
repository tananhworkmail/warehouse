const formatCellValue = (value, isDate) => {
  if (value === null || value === undefined) return "";

  if (isDate && typeof value === "string") {
    return value.replace("T", " ").replace(/\.\d{3}Z?$/, "");
  }

  return String(value);
};

/**
 * Keeps table-cell ordering deterministic while removing service fields that
 * are not meant to be displayed. Tracking tables receive row objects directly
 * from the scan API, so a small shared formatter is safer than mutating them.
 */
export function useRowOperations() {
  const createDisplayRow = (row, excludedFields = [], dateFields = []) => {
    if (!row || typeof row !== "object") return [];

    const excluded = new Set(excludedFields);
    const dates = new Set(dateFields);

    return Object.entries(row)
      .filter(([key]) => !excluded.has(key))
      .map(([key, value]) => formatCellValue(value, dates.has(key)));
  };

  return { createDisplayRow };
}
