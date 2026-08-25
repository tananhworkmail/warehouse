export const inventoryAgeCategories = [
  {
    key: "normal",
    color: "#22c55e",
    softColor: "#ecfdf5",
    labelKey: "ktp.inventoryAge.normal",
    shortLabelKey: "ktp.inventoryAge.normalShort",
  },
  {
    key: "caution",
    color: "#e2c20d",
    softColor: "#fefce8",
    labelKey: "ktp.inventoryAge.caution",
    shortLabelKey: "ktp.inventoryAge.cautionShort",
  },
  {
    key: "highRisk",
    color: "#f97c16",
    softColor: "#fff7ed",
    labelKey: "ktp.inventoryAge.highRisk",
    shortLabelKey: "ktp.inventoryAge.highRiskShort",
  },
  {
    key: "warning",
    color: "#ff0000",
    softColor: "#fff1f2",
    labelKey: "ktp.inventoryAge.warning",
    shortLabelKey: "ktp.inventoryAge.warningShort",
  },
];

export const unknownInventoryAgeCategory = {
  key: "unknown",
  color: "#94a3b8",
  softColor: "#f1f5f9",
  labelKey: "ktp.inventoryAge.unknown",
  shortLabelKey: "ktp.inventoryAge.unknown",
};

export const getInventoryAgeCategory = (days) => {
  if (days === null || days === undefined || days === "") {
    return unknownInventoryAgeCategory;
  }

  const numericDays = Number(days);
  if (!Number.isFinite(numericDays)) return unknownInventoryAgeCategory;
  if (numericDays <= 30) return inventoryAgeCategories[0];
  if (numericDays <= 60) return inventoryAgeCategories[1];
  if (numericDays <= 75) return inventoryAgeCategories[2];
  return inventoryAgeCategories[3];
};

export const normalizeInventoryAgeRow = (row = {}) => ({
  ddbh: String(row.ddbh || "").trim(),
  qty: Number(row.qty || 0),
  rackCode: String(row.rackCode || "").trim().toUpperCase(),
  lastInDate: row.lastInDate || "",
  storageAgeDays:
    row.storageAgeDays === null || row.storageAgeDays === undefined
      ? null
      : Number(row.storageAgeDays),
});
